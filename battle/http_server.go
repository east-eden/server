package battle

import (
	"context"
	"encoding/json"
	"expvar"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
)

var startTime = time.Now()
var lastGCPause uint32

func calculateUptime() interface{} {
	return time.Since(startTime).String()
}

func currentGoVersion() interface{} {
	return runtime.Version()
}

func getNumCPUs() interface{} {
	return runtime.NumCPU()
}

func getGoOS() interface{} {
	return runtime.GOOS
}

func getNumGoroutins() interface{} {
	return runtime.NumGoroutine()
}

func getLastGCPauseTime() interface{} {
	var gcPause uint64
	ms := new(runtime.MemStats)

	statString := expvar.Get("memstats").String()
	if statString != "" {
		json.Unmarshal([]byte(statString), ms)

		if lastGCPause == 0 || lastGCPause != ms.NumGC {
			gcPause = ms.PauseNs[(ms.NumGC+255)%256]
			lastGCPause = ms.NumGC
		}
	}

	return gcPause
}

type HttpServer struct {
	ctx    context.Context
	cancel context.CancelFunc
	b      *Battle
}

func NewHttpServer(b *Battle) *HttpServer {
	s := &HttpServer{
		b: b,
	}

	s.ctx, s.cancel = context.WithCancel(b.ctx)
	logger.Info("HttpServer listening at ", s.b.opts.HTTPListenAddr)
	return s
}

func (s *HttpServer) Run() error {

	expvar.Publish("ticktime", expvar.Func(calculateUptime))
	expvar.Publish("version", expvar.Func(currentGoVersion))
	expvar.Publish("cores", expvar.Func(getNumCPUs))
	expvar.Publish("os", expvar.Func(getGoOS))
	expvar.Publish("goroutine", expvar.Func(getNumGoroutins))
	expvar.Publish("gcpause", expvar.Func(getLastGCPauseTime))

	http.HandleFunc("/pub_battle_result", s.pubBattleResult)
	http.Handle("/metrics", promhttp.Handler())

	// battle run
	chExit := make(chan error)
	go func() {
		err := http.ListenAndServe(s.b.opts.HTTPListenAddr, nil)
		chExit <- err
	}()

	select {
	case <-s.ctx.Done():
		break
	case err := <-chExit:
		return err
	}

	logger.Info("HttpServer context done...")
	return nil
}

func (s *HttpServer) pubBattleResult(w http.ResponseWriter, r *http.Request) {
	s.b.BattleResult()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode([]byte("success"))
}
