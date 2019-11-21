package game

import (
	"context"
	"log"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
	pbClient "github.com/yokaiio/yokai_server/proto/client"
)

type Game struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper

	db         *Datastore
	httpSrv    *HttpServer
	tcpSrv     *TcpServer
	cm         *ClientManager
	pm         *PlayerManager
	mi         *MicroService
	rpcHandler *RpcHandler
	msgHandler *MsgHandler
	pubSub     *PubSub
}

func New(opts *Options) (*Game, error) {
	g := &Game{
		opts: opts,
	}

	g.ctx, g.cancel = context.WithCancel(context.Background())
	g.db = NewDatastore(g)
	g.httpSrv = NewHttpServer(g)
	g.tcpSrv = NewTcpServer(g)
	g.cm = NewClientManager(g)
	g.pm = NewPlayerManager(g)
	g.mi = NewMicroService(g)
	g.rpcHandler = NewRpcHandler(g)
	g.msgHandler = NewMsgHandler(g)
	g.pubSub = NewPubSub(g)

	return g, nil
}

func (g *Game) Main() error {

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Game Main() error:", err)
			}
			exitCh <- err
		})
	}

	// game run
	g.waitGroup.Wrap(func() {
		exitFunc(g.Run())
	})

	// database run
	g.waitGroup.Wrap(func() {
		exitFunc(g.db.Run())
	})

	// http server run
	g.waitGroup.Wrap(func() {
		exitFunc(g.httpSrv.Run())
	})

	// tcp server run
	g.waitGroup.Wrap(func() {
		g.tcpSrv.Run()
		g.tcpSrv.Exit()
	})

	// client mgr run
	g.waitGroup.Wrap(func() {
		g.cm.Main()
		g.cm.Exit()
	})

	// micro run
	g.waitGroup.Wrap(func() {
		exitFunc(g.mi.Run())
	})

	err := <-exitCh
	return err
}

func (g *Game) Exit() {
	g.cancel()
	g.waitGroup.Wait()
}

func (g *Game) Run() error {

	for {
		select {
		case <-g.ctx.Done():
			logger.Info("Game context done...")
			return nil
		default:
		}

		// todo game logic

		t := time.Now()
		d := time.Since(t)
		time.Sleep(time.Second - d)
	}
}

func (g *Game) StartBattle() {
	c := &pbClient.ClientInfo{Id: 12, Name: "game's client 12"}
	g.pubSub.PubStartBattle(g.ctx, c)
}
