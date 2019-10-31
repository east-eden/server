package game

import (
	"context"
	"log"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Game struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper

	db      *Datastore
	httpSrv *HttpServer
	tcpSrv  *TcpServer
	cm      *ClientMgr
}

func New(opts *Options) (*Game, error) {
	g := &Game{
		opts: opts,
	}

	g.ctx, g.cancel = context.WithCancel(context.Background())
	g.db = NewDatastore(g)
	g.httpSrv = NewHttpServer(g)
	g.tcpSrv = NewTcpServer(g)
	g.cm = NewClientMgr(g)

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
		exitFunc(g.tcpSrv.Main())
	})

	// client mgr run
	g.waitGroup.Wrap(func() {
		err := g.cm.Main()
		if err != nil {
			g.cm.Exit()
		}
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
