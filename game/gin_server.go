package game

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var upgrader = websocket.Upgrader{}

type GinServer struct {
	listenAddr string
	certPath   string
	keyPath    string
	ctx        context.Context
	cancel     context.CancelFunc
	g          *Game
	e          *gin.Engine

	socks  map[*websocket.Conn]struct{}
	chSend chan string
}

// timeout middleware wraps the request context with a timeout
func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func timedHandler(duration time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// get the underlying request context
		ctx := c.Request.Context()

		// create the response data type to use as a channel type
		type responseData struct {
			status int
			body   map[string]interface{}
		}

		// create a done channel to tell the request it's done
		doneChan := make(chan responseData)

		// here you put the actual work needed for the request
		// and then send the doneChan with the status and body
		// to finish the request by writing the response
		go func() {
			time.Sleep(duration)
			doneChan <- responseData{
				status: 200,
				body:   gin.H{"hello": "world"},
			}
		}()

		// non-blocking select on two channels see if the request
		// times out or finishes
		select {

		// if the context is done it timed out or was cancelled
		// so don't return anything
		case <-ctx.Done():
			return

			// if the request finished then finish the request by
			// writing the response
		case res := <-doneChan:
			c.JSON(res.status, res.body)
		}
	}
}

func (s *GinServer) setupRouter() {
	// Disable Console Color
	// gin.DisableConsoleColor()
	s.e.Use(timeoutMiddleware(time.Second * 120))

	// websocket
	s.e.GET("/ws", s.WebSocketHandler)

	// test websocket push
	s.e.GET("/ws_push", func(c *gin.Context) {
		s.chSend <- "push test!!!"
	})
}

func NewGinServer(g *Game, c *cli.Context) *GinServer {
	s := &GinServer{
		g:          g,
		e:          gin.Default(),
		listenAddr: c.String("https_listen_addr"),
		certPath:   c.String("cert_path_release"),
		keyPath:    c.String("key_path_release"),

		socks:  make(map[*websocket.Conn]struct{}),
		chSend: make(chan string, 10),
	}

	if c.Bool("debug") {
		s.certPath = c.String("cert_path_debug")
		s.keyPath = c.String("key_path_debug")
	}

	s.ctx, s.cancel = context.WithCancel(c)
	s.setupRouter()

	return s
}

func (s *GinServer) Run() error {
	chExit := make(chan error)
	go func() {
		err := s.e.RunTLS(
			s.listenAddr,
			s.certPath,
			s.keyPath,
		)

		chExit <- err
	}()

	select {
	case <-s.ctx.Done():
		break
	case err := <-chExit:
		return err
	}

	logger.Info("GinServer context done...")
	return nil
}

func (s *GinServer) WebSocketHandler(c *gin.Context) {
	ws, error := upgrader.Upgrade(c.Writer, c.Request, nil)
	if error != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	s.socks[ws] = struct{}{}

	go s.WebSocketSend(ws)
	go s.WebSocketRead(ws)
}

func (s *GinServer) WebSocketRead(ws *websocket.Conn) {
	defer ws.Close()

	ws.SetReadDeadline(time.Now().Add(time.Minute))

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		log.Println("websocket recv:", message)
	}
}

func (s *GinServer) WebSocketSend(ws *websocket.Conn) {
	defer ws.Close()

	for {
		select {
		case msg := <-s.chSend:
			for sock, _ := range s.socks {
				sock.SetWriteDeadline(time.Now().Add(time.Second * 10))
				sock.WriteMessage(websocket.TextMessage, []byte(msg))
			}
		}
	}
}
