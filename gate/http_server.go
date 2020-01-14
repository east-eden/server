package gate

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"io/ioutil"
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
	g              *Gate
}

func NewHttpServer(g *Gate, c *cli.Context) *HttpServer {
	s := &HttpServer{
		g:              g,
		httpListenAddr: c.String("http_listen_addr"),
	}

	s.ctx, s.cancel = context.WithCancel(c)
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

	http.HandleFunc("/store_write", s.storeWrite)
	http.HandleFunc("/get_game", s.getGameAddr)
	http.HandleFunc("/pub_gate_result", s.pubGateResult)
	http.HandleFunc("/get_lite_account", s.getLiteAccount)

	// gate run
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

func (s *HttpServer) storeWrite(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	s.g.mi.StoreWrite(req.Key, req.Value)
}

func (s *HttpServer) getGameAddr(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	user, metadata := s.g.gs.SelectGame(req.UserID, req.UserName)
	if user == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("cannot find account by userid<%s>", req.UserID)))
		return
	}

	var resp map[string]string
	resp["user_id"] = user.UserID
	resp["user_name"] = req.UserName
	resp["account_id"] = user.AccountID
	resp["game_id"] = metadata["game_id"]
	resp["public_addr"] = metadata["public_addr"]
	resp["section"] = metadata["section"]
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *HttpServer) pubGateResult(w http.ResponseWriter, r *http.Request) {
	s.g.GateResult()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

func (s *HttpServer) getLiteAccount(w http.ResponseWriter, r *http.Request) {
	rep, err := s.g.rpcHandler.CallGetRemoteLiteAccount(281587826959645248)
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
