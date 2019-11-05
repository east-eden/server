package battle

import (
	"context"
	"log"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Battle struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper

	db         *Datastore
	httpSrv    *HttpServer
	mi         *MicroService
	rpcHandler *RpcHandler
	pubSub     *PubSub
}

func New(opts *Options) (*Battle, error) {
	b := &Battle{
		opts: opts,
	}

	b.ctx, b.cancel = context.WithCancel(context.Background())
	b.db = NewDatastore(b)
	b.httpSrv = NewHttpServer(b)
	b.mi = NewMicroService(b)
	b.rpcHandler = NewRpcHandler(b)
	b.pubSub = NewPubSub(b)

	return b, nil
}

func (b *Battle) Main() error {

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Battle Main() error:", err)
			}
			exitCh <- err
		})
	}

	// battle run
	b.waitGroup.Wrap(func() {
		exitFunc(b.Run())
	})

	// database run
	b.waitGroup.Wrap(func() {
		exitFunc(b.db.Run())
	})

	// http server run
	b.waitGroup.Wrap(func() {
		exitFunc(b.httpSrv.Run())
	})

	// micro run
	b.waitGroup.Wrap(func() {
		exitFunc(b.mi.Run())
	})

	err := <-exitCh
	return err
}

func (b *Battle) Exit() {
	b.cancel()
	b.waitGroup.Wait()
}

func (b *Battle) Run() error {

	for {
		select {
		case <-b.ctx.Done():
			logger.Info("Battle context done...")
			return nil
		default:
		}

		// todo battle logic

		t := time.Now()
		d := time.Since(t)
		time.Sleep(time.Second - d)
	}
}

func (b *Battle) BattleResult() {
	b.pubSub.PubBattleResult(b.ctx, true)
}
