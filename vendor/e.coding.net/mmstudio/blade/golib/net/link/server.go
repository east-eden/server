package link

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/gorilla/websocket"
)

var defaultCloseTimeout = time.Duration(10) * time.Second

type Server struct {
	manager  *Manager
	listener net.Listener
	protocol Protocol
	handler  Handler
	spec     *Spec

	WebsocketUpGrader websocket.Upgrader
	HttpServer        *http.Server
}

type Handler interface {
	HandleSession(Session)
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(Session)

func (f HandlerFunc) HandleSession(sess Session) {
	f(sess)
}

func NewServerWithListener(listener net.Listener, spec *Spec, protocol Protocol, handler Handler) *Server {
	srv := &Server{
		listener: listener,
		manager:  NewManager(),
		protocol: protocol,
		handler:  handler,
		spec:     spec,
	}
	if srv.isWebSocket() {
		srv.initWebsocketUpGrader()
	}
	return srv
}

// Deprecated, use NewServerWithListener instead
func NewServer(spec *Spec, protocol Protocol, handler Handler) *Server {
	s := &Server{
		manager:  NewManager(),
		protocol: protocol,
		handler:  handler,
		spec:     spec,
	}

	addr, err := net.ResolveTCPAddr(spec.Proto, spec.Address)
	if err != nil {
		return nil
	}
	if s.isWebSocket() {
		// websocket proto 使用 tcp 的network
		s.listener, err = net.ListenTCP("tcp", addr)
		if err != nil {
			return nil
		}
		s.initWebsocketUpGrader()
	} else {
		s.listener, err = net.ListenTCP(spec.Proto, addr)
		if err != nil {
			return nil
		}
	}
	return s
}

func (server *Server) Listener() net.Listener     { return server.listener }
func (server *Server) SetHandler(handler Handler) { server.handler = handler }
func (server *Server) Handler() Handler           { return server.handler }
func (server *Server) Protocol() Protocol         { return server.protocol }

func (server *Server) Serve() error {
	if server.isWebSocket() {
		mux := http.NewServeMux()
		mux.HandleFunc(server.spec.WebSocketPattern, server.newWebConn)
		server.HttpServer = &http.Server{
			Handler:        mux,
			ReadTimeout:    server.spec.ReadTimeout,
			WriteTimeout:   server.spec.WriteTimeout,
			IdleTimeout:    server.spec.IdleTimeout,
			MaxHeaderBytes: server.spec.TCPReadBuffer,
		}

		var err error
		if server.spec.TLSCertFile != "" &&
			server.spec.TLSKeyFile != "" {
			err = server.HttpServer.ServeTLS(server.listener, server.spec.TLSCertFile, server.spec.TLSKeyFile)
		} else {
			err = server.HttpServer.Serve(server.listener)
		}
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		} else if strings.Contains(err.Error(), "use of closed network connection") {
			// 兼容低版本，go 1.16.3以上可以使用这个判断 errors.Is(err, net.ErrClosed)
			return nil
		}
		return err
	}
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		go func() {
			trans, err := server.protocol.NewTransporter(TCPConnSetup(conn.(*net.TCPConn), server.spec))
			if err != nil {
				_ = conn.Close()
				return
			}
			session := server.manager.NewSession(trans, server.spec)
			server.handler.HandleSession(session)
		}()
	}
}

func (server *Server) GetSession(sessionID uint64) (Session, bool) {
	return server.manager.GetSession(sessionID)
}

func (server *Server) GetManager() *Manager {
	return server.manager
}

func (server *Server) Stop(ctx context.Context) error {
	// step 1: close listeners, notify client reconnect
	if server.listener != nil {
		_ = server.listener.Close()
	}
	// step 2: wait and close connections
	done := make(chan struct{})
	go func() {
		server.manager.Dispose()
		close(done)
	}()

	if ctx == nil {
		ctx = context.Background()
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, time.Now().Add(defaultCloseTimeout))
		defer cancel()
	}
	select {
	case <-ctx.Done():
	case <-done:
	}
	return nil
}

func (server *Server) isWebSocket() bool {
	return strings.ToUpper(server.spec.Proto) == "WS" ||
		strings.ToUpper(server.spec.Proto) == "WEBSOCKET"
}

func (server *Server) initWebsocketUpGrader() {
	server.WebsocketUpGrader = websocket.Upgrader{
		HandshakeTimeout: server.spec.AcceptTimeout,
		ReadBufferSize:   server.spec.TCPReadBuffer,
		WriteBufferSize:  server.spec.TCPWriteBuffer,
	}
}

func (server *Server) newWebConn(w http.ResponseWriter, r *http.Request) {
	conn, err := server.WebsocketUpGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("remote", r.RemoteAddr).Msg("web conn upgrade fail")
		return
	}

	go func() {
		trans, err := server.protocol.NewTransporter(NewWebSocketConn(conn))
		if err != nil {
			_ = conn.Close()
			return
		}
		session := server.manager.NewSession(trans, server.spec)
		server.handler.HandleSession(session)
	}()
}
