package gate

import (
	"e.coding.net/mmstudio/blade/gate/msg"
	"e.coding.net/mmstudio/blade/golib/net/link"
	"github.com/rs/zerolog/log"
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"runtime/debug"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 10,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

func (g *Gate) webServe() {
	if !g.spec.EnableWebSocket {
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc(g.spec.PathForWebClient, g.newWebConn)
	g.wsProto = link.NewWsProtocol()
	g.wsServer = &http.Server{Addr: fmt.Sprintf(":%d", g.spec.PortForWebClient), Handler: mux}
	err := g.wsServer.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("start web server fail")
		}
	}
}

func (g *Gate) newWebConn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("remote", r.RemoteAddr).Msg("web conn upgrade fail")
		return
	}
	go g.handleWsConn(link.NewWebSocketConn(c))
}

func (g *Gate) handleWsConn(conn net.Conn) {
	backend, err := g.handshake(conn, true)
	if err != nil {
		_ = conn.Close()
		return
	}

	frontTrans, err := g.wsProto.NewTransporter(conn)
	if err != nil {
		_ = conn.Close()
		_ = backend.Close()
		return
	}
	backendTrans, err := g.spec.TransferProvider(backend)
	if err != nil {
		log.Warn().Err(err).Msg("get transfer for backend fail")
		resp := msg.HandshakeResp{
			Code: msg.ErrorCode_ServiceUnavailable,
			Desc: fmt.Sprintf("backend transfer fail:%s", err.Error()),
		}
		if data, err := g.spec.MessageCodec.Marshal(&resp); err == nil {
			log.Debug().Str("resp", resp.String()).Str("remote", conn.RemoteAddr().String()).Msg("response handshake")

			_ = frontTrans.Send(data)
		}
		_ = conn.Close()
		_ = backend.Close()
		return
	}
	go g.pipeWs(backend, backendTrans, frontTrans)
	g.pipeWs(conn, frontTrans, backendTrans)
}

func (g *Gate) pipeWs(destConn net.Conn, dest, src link.Transporter) {
	defer func() {
		if r := recover(); r != nil {
			log.Warn().Interface("recover", r).Str("stack", string(debug.Stack())).Msg("recover in pipe")
		}
	}()
	defer func() { _ = destConn.Close() }()

	for {
		select {
		case <-g.stopChan:
			return
		default:
		}
		_ = src.SetReadDeadline(time.Now().Add(g.spec.ConnReadTimeout))
		r, er := src.Receive()
		if er != nil {
			if isTemporaryErr(er) {
				continue
			} else {
				log.Error().Err(er).Msg("srcConn read fail")
				break
			}
		}
		for i := 0; i < 3; i++ {
			// retry 3 times
			ew := dest.Send(r)
			if ew != nil {
				if isTemporaryErr(ew) {
					continue
				} else {
					log.Error().Err(ew).Msg("destConn write fail")
					return
				}
			} else {
				break
			}
		}
	}
}
