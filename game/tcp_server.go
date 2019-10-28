package server

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hellodudu/Ultimate/iface"
	"github.com/hellodudu/Ultimate/utils/global"
	"github.com/hellodudu/Ultimate/utils/task"
	logger "github.com/sirupsen/logrus"
)

var tcpReadBufMax = 1024 * 1024 * 2

// TcpCon with closed status
type TCPCon struct {
	sync.Mutex
	con    net.Conn
	closed bool
}

func NewTCPCon(con net.Conn) *TCPCon {
	return &TCPCon{con: con, closed: false}
}

func (c *TCPCon) Close() {
	if c.closed {
		return
	}

	c.Lock()
	defer c.Unlock()
	c.closed = true
	c.con.Close()
}

func (c *TCPCon) Write(b []byte) (n int, err error) {
	if c.closed {
		return 0, fmt.Errorf("connection closed, nothing will be write in")
	}

	return c.con.Write(b)
}

func (c *TCPCon) Closed() bool {
	return c.closed
}

type TCPServer struct {
	conns      map[*TCPCon]struct{}
	ln         net.Listener
	parser     iface.IMsgParser
	dispatcher iface.IDispatcher
	mutexConns sync.Mutex
	wgConns    sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewTcpServer(parser iface.IMsgParser, dispatcher iface.IDispatcher) (*TCPServer, error) {
	s := &TCPServer{
		conns:      make(map[*TCPCon]struct{}),
		parser:     parser,
		dispatcher: dispatcher,
	}

	addr, err := global.GetIniMgr().GetIniValue("../config/ultimate.ini", "listen", "TcpListenAddr")
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	logger.Info("tcp server listening at ", addr)

	s.ln = ln
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s, nil
}

func (server *TCPServer) Run() {
	var tempDelay time.Duration
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				logger.WithFields(logger.Fields{
					"error":         err,
					"retry_seconds": tempDelay,
				}).Warn("accept error")

				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		connection := NewTCPCon(conn)

		server.mutexConns.Lock()
		if len(server.conns) >= 5000 {
			server.mutexConns.Unlock()
			connection.Close()
			logger.WithFields(logger.Fields{
				"connections": len(server.conns),
			}).Warn("too many connections")
			continue
		}
		server.conns[connection] = struct{}{}
		server.mutexConns.Unlock()

		server.wgConns.Add(1)
		go func(c *TCPCon) {
			server.handleTCPConnection(c)

			server.mutexConns.Lock()
			delete(server.conns, c)
			server.mutexConns.Unlock()

			server.wgConns.Done()
		}(connection)
	}
}

func (server *TCPServer) Stop() {
	server.ln.Close()
	server.cancel()
	server.wgConns.Wait()

	server.mutexConns.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.conns = nil
	server.mutexConns.Unlock()
}

func (server *TCPServer) handleTCPConnection(connection *TCPCon) {
	defer connection.Close()

	logger.Info("a new tcp connection with remote addr:", connection.con.RemoteAddr().String())
	connection.con.(*net.TCPConn).SetKeepAlive(true)
	connection.con.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)

	for {
		select {
		case <-server.ctx.Done():
			logger.Print("tcp connection context done!")
			return
		default:
		}

		if connection.Closed() {
			logger.Print("tcp connection closed:", connection)
			return
		}

		// read len
		b := make([]byte, 4)
		if _, err := io.ReadFull(connection.con, b); err != nil {
			logger.Info("one client connection was shut down:", err)
			return
		}

		var msgLen uint32
		msgLen = binary.LittleEndian.Uint32(b)

		// check len
		if msgLen > uint32(tcpReadBufMax) {
			logger.WithFields(logger.Fields{
				"error":  "message too long",
				"length": msgLen,
			}).Warn("tcp recv failed")
			continue
		} else if msgLen < 4 {
			logger.WithFields(logger.Fields{
				"error":  "message too short",
				"length": msgLen,
			}).Warn("tcp recv failed")
			continue
		}

		// data
		msgData := make([]byte, msgLen)
		if _, err := io.ReadFull(connection.con, msgData); err != nil {
			logger.WithFields(logger.Fields{
				"error": err,
			}).Warn("tcp recv failed")
			continue
		}

		server.dispatcher.AddTask(&task.TaskReqInfo{
			Con:  connection,
			Data: msgData,
			CB:   server.parser.ParserMessage,
		})
	}
}
