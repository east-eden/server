package combat

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"e.coding.net/mmstudio/blade/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type GinServer struct {
	listenAddr string
	certPath   string
	keyPath    string
	ctx        context.Context
	cancel     context.CancelFunc
	c          *Combat
	e          *gin.Engine
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

	// start_combat
	s.e.POST("/start_combat", func(c *gin.Context) {
		var req struct {
			UserID   string `json:"userId"`
			UserName string `json:"userName"`
		}

		if err := c.Bind(&req); err != nil {
			log.Warn().Err(err).Msg("select_game_addr request bind failed")
			c.String(http.StatusBadRequest, "bad request:%s", err.Error())
			return
		}

		//if user, metadata := s.g.gs.SelectGame(req.UserID, req.UserName); user != nil {
		//h := gin.H{
		//"userId":     user.UserID,
		//"userName":   req.UserName,
		//"accountId":  user.AccountID,
		//"gameId":     metadata["gameId"],
		//"publicAddr": metadata["publicAddr"],
		//}
		//c.JSON(http.StatusOK, h)

		//logger.Info("select_game_addr calling with result:", h)
		//return
		//}

		c.String(http.StatusBadRequest, fmt.Sprintf("cannot find account by userid<%s>", req.UserID))
	})

}

func NewGinServer(ctx *cli.Context, c *Combat) *GinServer {
	s := &GinServer{
		c:          c,
		e:          gin.Default(),
		listenAddr: ctx.String("https_listen_addr"),
		certPath:   ctx.String("cert_path_release"),
		keyPath:    ctx.String("key_path_release"),
	}

	if ctx.Bool("debug") {
		s.certPath = ctx.String("cert_path_debug")
		s.keyPath = ctx.String("key_path_debug")
	}

	s.ctx, s.cancel = context.WithCancel(ctx.Context)
	s.setupRouter()

	return s
}

func (s *GinServer) Run() error {
	chExit := make(chan error)
	go func() {
		defer utils.CaptureException()
		err := s.e.RunTLS(
			s.listenAddr,
			s.certPath,
			s.keyPath,
		)

		chExit <- err
	}()

	select {
	case <-s.ctx.Done():
		log.Info().Msg("GinServer context done...")
		return nil
	case err := <-chExit:
		return err
	}
}
