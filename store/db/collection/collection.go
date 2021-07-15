package collection

import (
	"context"
	"runtime/debug"
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/server/utils"
	log "github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	WriteModelBufferSize = 1024 // BulkWrite buffer size
	SleepDuration        = 50 * time.Millisecond
	FlushInterval        = 2 * time.Second
	BulkWriteTimeout     = 10 * time.Second
)

type Collection struct {
	*mongo.Collection
	models             []mongo.WriteModel
	writeChan          chan mongo.WriteModel
	closeChan          chan bool
	exitChan           chan bool
	flushImmediateChan chan bool
	t                  *time.Timer
	once               sync.Once
}

func NewCollection(coll *mongo.Collection) *Collection {
	c := &Collection{
		Collection:         coll,
		models:             make([]mongo.WriteModel, 0, WriteModelBufferSize),
		writeChan:          make(chan mongo.WriteModel, WriteModelBufferSize),
		closeChan:          make(chan bool, 1),
		exitChan:           make(chan bool, 1),
		flushImmediateChan: make(chan bool, 1),
		t:                  time.NewTimer(FlushInterval),
	}

	c.run()
	return c
}

func (c *Collection) ResetFlushInterval(d time.Duration) {
	if c.t != nil && !c.t.Stop() {
		<-c.t.C
	}
	c.t.Reset(d)
}

func (c *Collection) Write(m mongo.WriteModel) {
	c.writeChan <- m
}

func (c *Collection) Flush() {
	c.flushImmediateChan <- true
}

func (c *Collection) Exit() {
	c.once.Do(func() {
		close(c.closeChan)
		<-c.exitChan
	})
}

func (c *Collection) run() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				log.Error().Msgf("catch exception:%v, panic recovered with stack:%s", err, stack)
			}

			c.exitChan <- true
		}()

		for {
			select {
			case <-c.closeChan:
				c.flush()
				return
			case model := <-c.writeChan:
				c.models = append(c.models, model)
			case <-c.t.C:
				c.flush()
			case <-c.flushImmediateChan:
				c.flush()
			default:
				if len(c.models) < WriteModelBufferSize {
					time.Sleep(SleepDuration)
				} else {
					c.flush()
				}
			}
		}
	}()
}

func (c *Collection) flush() {
	if len(c.models) <= 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), BulkWriteTimeout)
	defer cancel()

	res, err := c.Collection.BulkWrite(ctx, c.models)
	_ = utils.ErrCheck(err, "BulkWrite failed when Collection.Flush", c.models, res)
	c.models = c.models[:0]
}
