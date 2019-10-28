package server

import (
	"context"
	"encoding/json"
	"expvar"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/hellodudu/Ultimate/iface"
	"github.com/hellodudu/Ultimate/utils/global"
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

func (s *HttpServer) getArenaPlayerDataNum() interface{} {
	return s.gm.GetArenaDataNum()
}

func (s *HttpServer) getArenaRecordNum() interface{} {
	return s.gm.GetArenaRecordNum()
}

type HttpServer struct {
	ctx    context.Context
	cancel context.CancelFunc
	gm     iface.IGameMgr
}

func NewHttpServer(gm iface.IGameMgr) *HttpServer {
	s := &HttpServer{
		gm: gm,
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

func (s *HttpServer) Run() {

	expvar.Publish("ticktime", expvar.Func(calculateUptime))
	expvar.Publish("version", expvar.Func(currentGoVersion))
	expvar.Publish("cores", expvar.Func(getNumCPUs))
	expvar.Publish("os", expvar.Func(getGoOS))
	expvar.Publish("goroutine", expvar.Func(getNumGoroutins))
	expvar.Publish("gcpause", expvar.Func(getLastGCPauseTime))
	expvar.Publish("arena_player_data_num", expvar.Func(s.getArenaPlayerDataNum))
	expvar.Publish("arena_record_num", expvar.Func(s.getArenaRecordNum))

	http.HandleFunc("/arena_matching_list", s.arenaMatchingListHandler)
	http.HandleFunc("/arena_record_req_list", s.arenaRecordReqListHandler)
	http.HandleFunc("/arena_get_record", s.arenaGetRecordHandler)
	http.HandleFunc("/arena_rank_list", s.arenaGetRankListHandler)
	http.HandleFunc("/arena_api_request_rank", s.arenaAPIRequestRankHandler)
	http.HandleFunc("/arena_save_champion", s.arenaSaveChampion)
	http.HandleFunc("/arena_weekend", s.arenaWeekEnd)
	http.HandleFunc("/player_info", s.getPlayerInfoHandler)
	http.HandleFunc("/guild_info", s.getGuildInfoHandler)

	http.Handle("/metrics", promhttp.Handler())

	addr, err := global.GetIniMgr().GetIniValue("../config/ultimate.ini", "listen", "HttpListenAddr")
	if err != nil {
		logger.Error("cannot read ini HttpListenAddr!")
		return
	}

	logger.Error(http.ListenAndServe(addr, nil))

}

func (s *HttpServer) arenaMatchingListHandler(w http.ResponseWriter, r *http.Request) {
	ids, err := s.gm.GetArenaMatchingList()
	if err != nil {
		return
	}

	var resp struct {
		ID []int64 `json:"id"`
	}

	resp.ID = ids

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp.ID)
}

func (s *HttpServer) arenaRecordReqListHandler(w http.ResponseWriter, r *http.Request) {
	list, err := s.gm.GetArenaRecordReqList()
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(list)
}

func (s *HttpServer) arenaGetRecordHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if r, err := s.gm.GetArenaRecord(req.ID); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(r)
	}
}

func (s *HttpServer) arenaGetRankListHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		Page int `json:"page"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if data, err := s.gm.GetArenaRankList(req.Page); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(data)
	}

}

func (s *HttpServer) arenaAPIRequestRankHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		ID   int64 `json:"id"`
		Page int   `json:"page"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if data, err := s.gm.ArenaAPIRequestRank(req.ID, req.Page); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *HttpServer) arenaSaveChampion(w http.ResponseWriter, r *http.Request) {

	if err := s.gm.ArenaSaveChampion(); err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("success"))
}

func (s *HttpServer) arenaWeekEnd(w http.ResponseWriter, r *http.Request) {

	if err := s.gm.ArenaWeekEnd(); err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("success"))
}

func (s *HttpServer) getPlayerInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	info, err := s.gm.GetPlayerInfoByID(req.ID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
			"id":    req.ID,
		}).Warn("cannot find player info by id")

		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(info)
}

func (s *HttpServer) getGuildInfoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	info, err := s.gm.GetGuildInfoByID(req.ID)
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err,
			"id":    req.ID,
		}).Warn("cannot find guild info by id")
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(info)
}
