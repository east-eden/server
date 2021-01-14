package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
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

	// consul service discovery
	s.engine.GET("/consul_get", func(c *gin.Context) {
		c.String(http.StatusOK, "consul_get success")
	})

	s.engine.GET("/watch_key", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Err(err).Msg("Error reading body")
			return
		}

		var objmap map[string]interface{}
		err = json.Unmarshal(body, &objmap)
		if err != nil {
			log.Err(err).Str("body", string(body)).Msg("unmarshal json failed")
			return
		}

		log.Info().Interface("body", objmap).Msg("watch_key success!")
		c.String(http.StatusOK, "watch_key success")
	})

	s.engine.GET("/watch_service", func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Err(err).Msg("Error reading body")
			return
		}

		var objmap []interface{}
		err = json.Unmarshal(body, &objmap)
		if err != nil {
			log.Err(err).Str("body", string(body)).Msg("unmarshal json failed")
			return
		}

		log.Info().Interface("body", objmap).Msg("watch_service success!")
		c.String(http.StatusOK, "watch_service success")
	})
}

func NewGinServer(ctx *cli.Context) *GinServer {
	s := &GinServer{
		engine: gin.Default(),
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info().Msgf("[GIN-debug] %s %s %s %d", httpMethod, absolutePath, handlerName, nuHandlers)
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
				log.Fatal().Err(err).Msg("GinServer Run() error")
			}
			exitCh <- err
		})
	}

	s.wg.Wrap(func() {
		exitFunc(s.Run(ctx))
	})

	// listen http
	go func() {
		defer utils.CaptureException()

		if err := s.engine.Run(ctx.String("http_listen_addr")); err != nil {
			log.Error().Err(err).Msg("GinServer Run error")
			exitCh <- err
		}
	}()

	return <-exitCh
}

func (s *GinServer) Run(ctx *cli.Context) error {
	<-ctx.Done()
	log.Info().Msg("GinServer context done...")
	return nil
}

func (s *GinServer) Exit(ctx context.Context) {
	s.wg.Wait()
	log.Info().Msg("gin server exit...")
}
