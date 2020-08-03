package client

import (
	"context"
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/utils"
)

type GinServer struct {
	engine *gin.Engine
	wg     utils.WaitGroupWrapper
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

func (s *GinServer) setupHttpRouter() {
	s.engine.Use(timeoutMiddleware(time.Second * 5))

	// pprof
	s.engine.GET("/debug/pprof", ginHandlerWrapper(pprof.Index))
	s.engine.GET("/debug/cmdline", ginHandlerWrapper(pprof.Cmdline))
	s.engine.GET("/debug/symbol", ginHandlerWrapper(pprof.Symbol))
	s.engine.GET("/debug/profile", ginHandlerWrapper(pprof.Profile))
	s.engine.GET("/debug/trace", ginHandlerWrapper(pprof.Trace))
	s.engine.GET("/debug/allocs", ginHandlerWrapper(pprof.Handler("allocs").ServeHTTP))
	s.engine.GET("/debug/heap", ginHandlerWrapper(pprof.Handler("heap").ServeHTTP))
	s.engine.GET("/debug/goroutine", ginHandlerWrapper(pprof.Handler("goroutine").ServeHTTP))
	s.engine.GET("/debug/block", ginHandlerWrapper(pprof.Handler("block").ServeHTTP))
	s.engine.GET("/debug/threadcreate", ginHandlerWrapper(pprof.Handler("threadcreate").ServeHTTP))

	// metrics
	s.engine.GET("/metrics", ginHandlerWrapper(promhttp.Handler().ServeHTTP))
}

func NewGinServer(ctx *cli.Context) *GinServer {
	s := &GinServer{
		engine: gin.Default(),
	}

	s.setupHttpRouter()
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

	// listen http
	go func() {
		if err := s.engine.Run(ctx.String("http_listen_addr")); err != nil {
			logger.Error("GinServer Run error:", err)
			exitCh <- err
		}
	}()

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

func (s *GinServer) Exit(ctx context.Context) {
	s.wg.Wait()
	logger.Info("gin server exit...")
}