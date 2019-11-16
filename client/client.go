package client

import (
	"context"
	"log"
	"sync"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/internal/utils"
)

type Client struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper
}

func New(opts *Options) (*Client, error) {
	c := &Client{
		opts: opts,
	}

	c.ctx, c.cancel = context.WithCancel(context.Background())

	return c, nil
}

func (c *Client) Main() error {

	exitCh := make(chan error)
	var once sync.Once
	exitFunc := func(err error) {
		once.Do(func() {
			if err != nil {
				log.Fatal("Client Main() error:", err)
			}
			exitCh <- err
		})
	}

	// client run
	c.waitGroup.Wrap(func() {
		exitFunc(c.Run())
	})

	err := <-exitCh
	return err
}

func (c *Client) Exit() {
	c.cancel()
	c.waitGroup.Wait()
}

func (c *Client) Run() error {

	for {
		select {
		case <-c.ctx.Done():
			logger.Info("Client context done...")
			return nil
		default:
		}

		// todo client logic

		t := time.Now()
		d := time.Since(t)
		time.Sleep(200*time.Millisecond - d)
	}
}
