package game

import (
	"context"
	"encoding/json"
	"expvar"
	"net/http"
	"runtime"
	"time"

	logger "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
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
	httpListenAddr string
	ctx            context.Context
	cancel         context.CancelFunc
	g              *Game
}

func NewHttpServer(g *Game, ctx *cli.Context) *HttpServer {
	s := &HttpServer{
		g:              g,
		httpListenAddr: ctx.String("http_listen_addr"),
	}

	s.ctx, s.cancel = context.WithCancel(ctx)
	logger.Info("HttpServer listening at ", s.httpListenAddr)
	return s
}

func (s *HttpServer) Run() error {

	expvar.Publish("ticktime", expvar.Func(calculateUptime))
	expvar.Publish("version", expvar.Func(currentGoVersion))
	expvar.Publish("cores", expvar.Func(getNumCPUs))
	expvar.Publish("os", expvar.Func(getGoOS))
	expvar.Publish("goroutine", expvar.Func(getNumGoroutins))
	expvar.Publish("gcpause", expvar.Func(getLastGCPauseTime))

	http.HandleFunc("/pub_start_battle", s.pubStartBattle)
	http.HandleFunc("/get_battle_status", s.getBattleStatus)

	// game run
	chExit := make(chan error)
	go func() {
		err := http.ListenAndServe(s.httpListenAddr, nil)
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

func (s *HttpServer) pubStartBattle(w http.ResponseWriter, r *http.Request) {
	s.g.StartBattle()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (s *HttpServer) getBattleStatus(w http.ResponseWriter, r *http.Request) {
	rep, err := s.g.rpcHandler.GetBattleStatus()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("error"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rep)
}
