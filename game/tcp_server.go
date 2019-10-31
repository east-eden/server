package game

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/define"
)

// TcpCon with closed status
type TcpCon struct {
	sync.Mutex
	con    net.Conn
	closed bool
}

func NewTcpCon(con net.Conn) *TcpCon {
	return &TcpCon{con: con, closed: false}
}

func (c *TcpCon) Close() {
	if c.closed {
		return
	}

	c.Lock()
	defer c.Unlock()
	c.closed = true
	c.con.Close()
}

func (c *TcpCon) Write(b []byte) (n int, err error) {
	if c.closed {
		return 0, fmt.Errorf("connection closed, nothing will be write in")
	}

	return c.con.Write(b)
}

func (c *TcpCon) Closed() bool {
	return c.closed
}

type TcpServer struct {
	conns  map[*TcpCon]struct{}
	ln     net.Listener
	parser *MsgParser
	wp     *workerpool.WorkerPool
	mu     sync.Mutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTcpServer(g *Game) *TcpServer {
	s := &TcpServer{
		conns:  make(map[*TcpCon]struct{}),
		parser: NewMsgParser(g),
		wp:     workerpool.New(runtime.GOMAXPROCS(runtime.NumCPU())),
	}

	ln, err := net.Listen("tcp", g.opts.TCPListenAddr)
	if err != nil {
		logger.Fatal("NewTcpServer error", err)
		return nil
	}

	logger.Info("tcp server listening at ", g.opts.TCPListenAddr)

	s.ln = ln
	s.ctx, s.cancel = context.WithCancel(g.ctx)
	return s
}

func (s *TcpServer) Run() error {
	var tempDelay time.Duration
	for {
		conn, err := s.ln.Accept()
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
			return err
		}
		tempDelay = 0

		c := NewTcpCon(conn)

		s.mu.Lock()
		if len(s.conns) >= s.parser.g.opts.ClientConnectMax {
			s.mu.Unlock()
			c.Close()
			logger.WithFields(logger.Fields{
				"connections": len(s.conns),
			}).Warn("too many connections")
			continue
		}
		s.conns[c] = struct{}{}
		s.mu.Unlock()

		s.wg.Add(1)
		go func(con *TcpCon) {
			s.handleConnection(con)

			s.mu.Lock()
			delete(s.conns, con)
			s.mu.Unlock()

			s.wg.Done()
		}(c)
	}
}

func (s *TcpServer) Stop() {
	s.ln.Close()
	s.cancel()
	s.wg.Wait()

	s.mu.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mu.Unlock()
}

func (s *TcpServer) handleConnection(c *TcpCon) {
	defer c.Close()

	logger.Info("a new tcp connection with remote addr:", c.con.RemoteAddr().String())
	c.con.(*net.TCPConn).SetKeepAlive(true)
	c.con.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)

	for {
		select {
		case <-s.ctx.Done():
			logger.Print("tcp connection context done!")
			return
		default:
		}

		if c.Closed() {
			logger.Print("tcp connection closed:", c)
			return
		}

		// read len
		b := make([]byte, 4)
		if _, err := io.ReadFull(c.con, b); err != nil {
			logger.Info("one client connection was shut down:", err)
			return
		}

		var msgLen uint32
		msgLen = binary.LittleEndian.Uint32(b)

		// check len
		if msgLen > uint32(define.TCPReadBufMax) {
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
		if _, err := io.ReadFull(c.con, msgData); err != nil {
			logger.WithFields(logger.Fields{
				"error": err,
			}).Warn("tcp recv failed")
			continue
		}

		// add to worker pool
		c := c
		m := msgData
		p := s.parser
		s.wp.Submit(func() {
			p.ParserMessage(c, m)
		})
	}
}
