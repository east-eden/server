package game

import (
	"context"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/aviddiviner/gin-limit"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/yokaiio/yokai_server/logger"
	"github.com/yokaiio/yokai_server/utils"
)

var (
	httpReadTimeout           = time.Second * 5
	httpWriteTimeout          = time.Second * 5
	ginConcurrentRequestLimit = 1000
)

var (
	opsSelectGameCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "select_game_addr_ops_total",
		Help: "选择服务器操作总数",
	})
)

type GinServer struct {
	g         *Game
	router    *gin.Engine
	tlsRouter *gin.Engine
	wg        utils.WaitGroupWrapper
}

// wrap http.HandlerFunc to gin.HandlerFunc
func ginHandlerWrapper(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}

func (s *GinServer) setupHttpRouter() {
	s.router.Use(limit.MaxAllowed(ginConcurrentRequestLimit))
	s.router.Use(gin.LoggerWithWriter(logger.Logger))

	// pprof
	s.router.GET("/debug/pprof", ginHandlerWrapper(pprof.Index))
	s.router.GET("/debug/cmdline", ginHandlerWrapper(pprof.Cmdline))
	s.router.GET("/debug/symbol", ginHandlerWrapper(pprof.Symbol))
	s.router.GET("/debug/profile", ginHandlerWrapper(pprof.Profile))
	s.router.GET("/debug/trace", ginHandlerWrapper(pprof.Trace))
	s.router.GET("/debug/allocs", ginHandlerWrapper(pprof.Handler("allocs").ServeHTTP))
	s.router.GET("/debug/heap", ginHandlerWrapper(pprof.Handler("heap").ServeHTTP))
	s.router.GET("/debug/goroutine", ginHandlerWrapper(pprof.Handler("goroutine").ServeHTTP))
	s.router.GET("/debug/block", ginHandlerWrapper(pprof.Handler("block").ServeHTTP))
	s.router.GET("/debug/threadcreate", ginHandlerWrapper(pprof.Handler("threadcreate").ServeHTTP))

	// metrics
	s.router.GET("/metrics", ginHandlerWrapper(promhttp.Handler().ServeHTTP))
}

func (s *GinServer) setupHttpsRouter() {
	s.tlsRouter.Use(limit.MaxAllowed(ginConcurrentRequestLimit))
	s.tlsRouter.Use(gin.LoggerWithWriter(logger.Logger))
}

func NewGinServer(g *Game, ctx *cli.Context) *GinServer {
	s := &GinServer{
		g:         g,
		router:    gin.Default(),
		tlsRouter: gin.Default(),
	}

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Info().Msgf("[GIN-debug] %s %s %s %d", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	s.setupHttpRouter()
	s.setupHttpsRouter()
	return s
}

func (s *GinServer) Main(ctx *cli.Context) error {
	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal().Err(err).Msg("GinServer Run() failed")
			}
			exitCh <- err
		})
	}

	s.wg.Wrap(func() {
		exitFunc(s.Run(ctx))
	})

	// listen https
	go func() {
		defer utils.CaptureException()

		certPath := ctx.String("cert_path_release")
		keyPath := ctx.String("key_path_release")
		if ctx.Bool("debug") {
			certPath = ctx.String("cert_path_debug")
			keyPath = ctx.String("key_path_debug")
		}

		server := &http.Server{
			Addr:         ctx.String("https_listen_addr"),
			Handler:      s.tlsRouter,
			ReadTimeout:  httpReadTimeout,
			WriteTimeout: httpWriteTimeout,
		}

		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Error().Err(err).Msg("gin server ListenAndServeTLS return with error")
			exitCh <- err
		}
	}()

	// listen http
	go func() {
		defer utils.CaptureException()

		server := &http.Server{
			Addr:         ctx.String("http_listen_addr"),
			Handler:      s.router,
			ReadTimeout:  httpReadTimeout,
			WriteTimeout: httpWriteTimeout,
		}

		if err := server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("gin server ListenAndServe return with error")
			exitCh <- err
		}
	}()

	return <-exitCh
}

func (s *GinServer) Run(ctx *cli.Context) error {

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("gin server context done...")
			return nil
		}
	}

	return nil
}

func (s *GinServer) Exit(ctx context.Context) {
	s.wg.Wait()
	log.Info().Msg("gin server exit...")
}
