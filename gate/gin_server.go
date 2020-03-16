package gate

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var users = make(map[string]string)

type GinServer struct {
	httpListenAddr string
	ctx            context.Context
	cancel         context.CancelFunc
	g              *Gate
	e              *gin.Engine
}

func (s *GinServer) setupRouter() {
	// Disable Console Color
	// gin.DisableConsoleColor()

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
			UserID   string `json:"user_id"`
			UserName string `json:"user_name"`
		}

		if c.Bind(&req) != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		if user, metadata := s.g.gs.SelectGame(req.UserID, req.UserName); user != nil {
			c.JSON(http.StatusOK, gin.H{
				"user_id":     req.UserID,
				"user_name":   req.UserName,
				"account_id":  strconv.FormatInt(user.AccountID, 10),
				"game_id":     metadata["game_id"],
				"public_addr": metadata["public_addr"],
				"section":     metadata["section"],
			})
			return
		}

		c.String(http.StatusBadRequest, fmt.Sprintf("cannot find account by userid<%s>", req.UserID))
	})

	// pub_gate_result
	s.e.POST("/pub_gate_result", func(c *gin.Context) {
		s.g.GateResult()
		c.String(http.StatusOK, "status ok")
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

func NewGinServer(g *Gate, c *cli.Context) *GinServer {
	s := &GinServer{
		g:              g,
		e:              gin.Default(),
		httpListenAddr: c.String("http_listen_addr"),
	}

	s.ctx, s.cancel = context.WithCancel(c)
	s.setupRouter()

	return s
}

func (s *GinServer) Run() error {
	chExit := make(chan error)
	go func() {
		err := s.e.Run(s.httpListenAddr)
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
