package gate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
)

var users = make(map[string]string)

type GinServer struct {
	g  *Gate
	e  *gin.Engine
	wg utils.WaitGroupWrapper
}

// wrap http.HandlerFunc to gin.HandlerFunc
func ginHandlerWrapper(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
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

	// pprof
	s.e.GET("/debug/pprof", ginHandlerWrapper(pprof.Index))
	s.e.GET("/debug/cmdline", ginHandlerWrapper(pprof.Cmdline))
	s.e.GET("/debug/symbol", ginHandlerWrapper(pprof.Symbol))
	s.e.GET("/debug/profile", ginHandlerWrapper(pprof.Profile))
	s.e.GET("/debug/allocs", ginHandlerWrapper(pprof.Handler("allocs").ServeHTTP))
	s.e.GET("/debug/heap", ginHandlerWrapper(pprof.Handler("heap").ServeHTTP))
	s.e.GET("/debug/goroutine", ginHandlerWrapper(pprof.Handler("goroutine").ServeHTTP))
	s.e.GET("/debug/block", ginHandlerWrapper(pprof.Handler("block").ServeHTTP))
	s.e.GET("/debug/threadcreate", ginHandlerWrapper(pprof.Handler("threadcreate").ServeHTTP))

	// store_write
	s.e.POST("/store_write", func(c *gin.Context) {
		var req struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if c.Bind(&req) == nil {
			s.g.mi.StoreWrite(req.Key, req.Value)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		c.String(http.StatusBadRequest, "bad request")
	})

	// select_game_addr
	s.e.POST("/select_game_addr", func(c *gin.Context) {
		var req struct {
			UserID   string `json:"userId"`
			UserName string `json:"userName"`
		}

		if err := c.Bind(&req); err != nil {
			logger.Warn("select_game_addr request bind error:", err)
			c.String(http.StatusBadRequest, "bad request:%s", err.Error())
			return
		}

		if user, metadata := s.g.gs.SelectGame(req.UserID, req.UserName); user != nil {
			h := gin.H{
				"userId":     user.UserID,
				"userName":   req.UserName,
				"accountId":  user.AccountID,
				"gameId":     metadata["gameId"],
				"publicAddr": metadata["publicAddr"],
				"section":    metadata["section"],
			}
			c.JSON(http.StatusOK, h)

			logger.Info("select_game_addr calling with result:", h)
			return
		}

		c.String(http.StatusBadRequest, fmt.Sprintf("cannot find account by userid<%s>", req.UserID))
	})

	// pub_gate_result
	s.e.POST("/pub_gate_result", func(c *gin.Context) {
		s.g.GateResult()
		c.String(http.StatusOK, "status ok")
	})

	// update_player_exp
	s.e.POST("/update_player_exp", func(c *gin.Context) {
		var req struct {
			Id string `json:"id"`
		}

		if c.Bind(&req) == nil {
			id, err := strconv.ParseInt(req.Id, 10, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "request error")
				return
			}
			r, err := s.g.rpcHandler.CallUpdatePlayerExp(id)
			c.String(http.StatusOK, "UpdatePlayerExp result", r, err)
		}
	})

	// get_lite_account
	s.e.POST("/get_lite_account", func(c *gin.Context) {
		var req struct {
			AccountID string `json:"account_id"`
		}

		if c.Bind(&req) == nil {
			id, err := strconv.ParseInt(req.AccountID, 10, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "request error")
				return
			}

			rep, err := s.g.rpcHandler.CallGetRemoteLiteAccount(id)
			if err == nil {
				c.JSON(http.StatusOK, rep)
				return
			}

			c.String(http.StatusBadRequest, err.Error())
			return
		}

		c.String(http.StatusBadRequest, "request error")
	})
}

func NewGinServer(g *Gate, ctx *cli.Context) *GinServer {
	s := &GinServer{
		g: g,
		e: gin.Default(),
	}

	s.setupRouter()
	return s
}

func (s *GinServer) Main(ctx *cli.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("GinServer Run() error:", err)
			}
			exitCh <- err
		})
	}

	s.wg.Wrap(func() {
		exitFunc(s.Run(ctx))
	})

	s.wg.Wrap(func() {
		certPath := ctx.String("cert_path_release")
		keyPath := ctx.String("key_path_release")
		if ctx.Bool("debug") {
			certPath = ctx.String("cert_path_debug")
			keyPath = ctx.String("key_path_debug")
		}

		go func() {
			if err := s.e.RunTLS(ctx.String("https_listen_addr"), certPath, keyPath); err != nil {
				logger.Error("GinServer RunTLS error:", err)
			}
		}()

		select {
		case <-ctx.Done():
			return
		}
	})

	return <-exitCh
}

func (s *GinServer) Run(ctx *cli.Context) error {

	for {
		select {
		case <-ctx.Done():
			logger.Info("GinServer context done...")
			return nil
		}
	}

	return nil
}

func (s *GinServer) Exit(ctx context.Context) error {
	s.wg.Wait()
	return nil
}
