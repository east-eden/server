package game

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yokaiio/yokai_server/internal/utils"
)

type Game struct {
	sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	opts      *Options
	waitGroup utils.WaitGroupWrapper
}

func New(opts *Options) (*Game, error) {
	g := &Game{
		opts: opts,
	}

	g.ctx, g.cancel = context.WithCancel(context.Background())

	return g, nil
}

// Main starts an instance of loki_conn and returns an
// error if there was a problem starting up.
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

	g.waitGroup.Wrap(func() {
		exitFunc(g.Run())
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
		case <-l.ctx.Done():
			return nil
		default:
		}

		//req := &PushRequest{
		//Streams: make([]*Stream, 0),
		//}

		//entry := &Entry{
		//TS:   time.Now().Format(time.RFC3339),
		//Line: "[info] heartbeat",
		//}

		//entries := make([]*Entry, 0)
		//entries = append(entries, entry)

		//labels := "{loki_conn=\"connection\"}"
		//req.Streams = append(req.Streams, &Stream{Labels: labels, Entries: entries})
		//reqJSON, err := json.Marshal(req)
		//if err != nil {
		//log.Println("marshal json error:", err)
		//d := time.Since(t)
		//time.Sleep(l.opts.Interval - d)
		//continue
		//}

		//request, err := http.NewRequest("POST", l.opts.URL, bytes.NewBuffer(reqJSON))
		//request.Header.Set("X-Custom-Header", "myvalue")
		//request.Header.Set("Content-Type", "application/json")

		//client := &http.Client{}
		//resp, err := client.Do(request)
		//if err != nil {
		//log.Println("http request with error:", err)
		//d := time.Since(t)
		//time.Sleep(l.opts.Interval - d)
		//continue
		//}

		//defer resp.Body.Close()

		//body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("response Status:", resp.Status, ", Body:", string(body))

		t := time.Now()
		d := time.Since(t)
		time.Sleep(time.Second - d)
	}
}
