package fkqst

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/east-eden/server/utils"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
)

// Adapter 提供骑士团对接所需的配置项（地址、API 凭据等）。
// kof 里由平台适配器实现，这里抽象成接口，避免依赖具体平台 SDK。
type Adapter interface {
	// GetConfigs 返回配置值，第一个参数是 key，其余为可选策略参数
	GetConfigs(key string, opts ...string) string
}

// 骑士团角色信息
type PlayerInfo struct {
	PlayerId  string `json:"player_id"`
	PlayerUid int64  `json:"player_uid"`
	AccountId string `json:"account_id"`
	PlayName  string `json:"play_name"`
	Avatar    string `json:"avatar"`
	GuildId   string `json:"guild_id"`
	GuildUid  int64  `json:"guild_uid"`
	GuildName string `json:"guild_name"`
	ServerId  int32  `json:"server_id"`
}

// 请求骑士团角色信息
func QueryPlayerInfo(adapter Adapter, strategy string, playerId string) *PlayerInfo {
	if adapter == nil {
		return nil
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id": playerId,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/player-info"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return nil
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", playerId).
		Msg("fkqst api: prepare call query player info")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return nil
	}

	var response struct {
		ErrCode int            `json:"err_code"`
		ErrMsg  string         `json:"err_msg"`
		Data    map[string]any `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return nil
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: query player info response data is nil")
		return nil
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: query player info response")

	return &PlayerInfo{
		AccountId: cast.ToString(response.Data["account_id"]),
		PlayerId:  cast.ToString(response.Data["player_id"]),
		PlayName:  cast.ToString(response.Data["play_name"]),
		PlayerUid: cast.ToInt64(response.Data["player_uid"]),
		ServerId:  cast.ToInt32(response.Data["server_id"]),
	}
}

type PlayerItem struct {
	Num     int32 `json:"num"`
	FreeNum int32 `json:"free_num"`
}

type PlayerItems struct {
	GoodsNum map[int32]*PlayerItem `json:"goods_num"`
}

// 请求骑士团道具
func QueryPlayerItems(adapter Adapter, strategy string, playerId string, itemIds []int32) *PlayerItems {
	if adapter == nil {
		return nil
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id": playerId,
		"goods_ids": itemIds,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/goods-info"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return nil
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", playerId).
		Msg("fkqst api: prepare call query player items")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return nil
	}

	var response struct {
		ErrCode int          `json:"err_code"`
		ErrMsg  string       `json:"err_msg"`
		Data    *PlayerItems `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return nil
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: query player items response data is nil")
		return nil
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: query player items response")

	return response.Data
}

// 扣除骑士团道具
func CostPlayerItems(adapter Adapter, strategy string, uid int64, playerId string, itemId int32, num int32) (*PlayerItem, int32, bool) {
	if adapter == nil {
		return nil, 0, false
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id":         playerId,
		"goods_id":          itemId,
		"num":               num,
		"daidai_player_uid": uid,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/deduct-goods"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return nil, 0, false
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", playerId).
		Msg("fkqst api: prepare call cost player items")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return nil, 0, false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return nil, 0, false
	}

	var response struct {
		ErrCode int            `json:"err_code"`
		ErrMsg  string         `json:"err_msg"`
		Data    map[string]any `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return nil, 0, false
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: cost player items response data is nil")
		return nil, 0, false
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: cost player items response")

	return &PlayerItem{
		Num:     cast.ToInt32(response.Data["goods_num"]),
		FreeNum: cast.ToInt32(response.Data["ext_goods_num"]),
	}, cast.ToInt32(response.Data["deduct_type"]), true
}

type PlayerTask struct {
	TaskType int32 `json:"task_type"`
	Progress int32 `json:"progress"`
}

// 同步任务进度
func SyncPlayerTask(adapter Adapter, strategy string, playerId string, taskType int32, progress int32) *PlayerTask {
	if adapter == nil {
		return nil
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id": playerId,
		"task_type": taskType,
		"progress":  progress,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/task-progress"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return nil
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", playerId).
		Msg("fkqst api: prepare call sync player task")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return nil
	}

	var response struct {
		ErrCode int         `json:"err_code"`
		ErrMsg  string      `json:"err_msg"`
		Data    *PlayerTask `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return nil
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: sync player task response data is nil")
		return nil
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: sync player task response")

	return response.Data
}

type PlayerAward struct {
	PlayerId string          `json:"player_id"`
	Status   int32           `json:"status"`
	Awards   map[int32]int32 `json:"awards"`
}

// 同步结算奖励
func SyncBattleResult(adapter Adapter, strategy string, uid int64, reqId string, playerId string, passLevel int32, pass bool) bool {
	if adapter == nil {
		return false
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"req_id":            reqId,
		"player_id":         playerId,
		"pass_level":        passLevel,
		"battle_mode":       0,
		"is_pass":           pass,
		"daidai_player_uid": uid,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/daidai-battle-result"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return false
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", playerId).
		Msg("fkqst api: prepare call sync battle result")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return false
	}

	var response struct {
		ErrCode int          `json:"err_code"`
		ErrMsg  string       `json:"err_msg"`
		Data    *PlayerAward `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return false
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: sync player battle result response data is nil")
		return false
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: sync player battle result response")
	return true
}

// 绑定骑士团账号
func BindAccount(adapter Adapter, strategy string, qstPlayerId string, qstAccountId string, playerId int64, accountId int64, uid int64) bool {
	if adapter == nil {
		return false
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id":         qstPlayerId,
		"account_id":        qstAccountId,
		"daidai_account_id": accountId,
		"daidai_player_id":  playerId,
		"daidai_player_uid": uid,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/bind-daidai-account"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return false
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", qstPlayerId).
		Str("qstAccountId", qstAccountId).
		Int64("playerId", playerId).
		Int64("accountId", accountId).
		Int64("uid", uid).
		Msg("fkqst api: prepare call bind account")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return false
	}

	var response struct {
		ErrCode int            `json:"err_code"`
		ErrMsg  string         `json:"err_msg"`
		Data    map[string]any `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return false
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: bind account response data is nil")
		return false
	}

	if response.ErrCode != 0 {
		log.Warn().Str("errMsg", response.ErrMsg).Msg("fkqst api: bind account response error code is nil")
		return false
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: bind account response")

	return true
}

// 解绑骑士团账号
func UnbindAccount(adapter Adapter, strategy string, qstPlayerId string, qstAccountId string, playerId int64, accountId int64, uid int64) bool {
	if adapter == nil {
		return false
	}

	now := utils.NowInLocation().UnixMilli()
	content := map[string]interface{}{
		"player_id":         qstPlayerId,
		"account_id":        qstAccountId,
		"daidai_account_id": accountId,
		"daidai_player_id":  playerId,
		"daidai_player_uid": uid,
	}
	contentStr := utils.JsonStringfy(content)
	url := adapter.GetConfigs("fkqstCallbackUrl", strategy) + "open/unbind-daidai-account"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(contentStr)))
	if !utils.ErrCheck(err, "get request failed") {
		return false
	}

	// 计算签名
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s.%s.%s.%s", adapter.GetConfigs("fkqstDaiDaiApiId"), cast.ToString(now), contentStr, adapter.GetConfigs("fkqstDaiDaiApiToken", strategy))))
	sign := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GameId", adapter.GetConfigs("fkqstDaiDaiApiId"))
	req.Header.Set("TS", cast.ToString(now))
	req.Header.Set("Sign", sign)

	log.Info().
		Str("url", url).
		Str("qstPlayerId", qstPlayerId).
		Str("qstAccountId", qstAccountId).
		Int64("playerId", playerId).
		Int64("accountId", accountId).
		Int64("uid", uid).
		Msg("fkqst api: prepare call unbind account")

	resp, err := http.DefaultClient.Do(req)
	if !utils.ErrCheck(err, "fkqst api: do request failed") {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if !utils.ErrCheck(err, "fkqst api: read response body failed") {
		return false
	}

	var response struct {
		ErrCode int            `json:"err_code"`
		ErrMsg  string         `json:"err_msg"`
		Data    map[string]any `json:"data"`
	}

	err = json.Unmarshal(body, &response)
	if !utils.ErrCheck(err, "fkqst api: unmarshal response body failed", string(body)) {
		return false
	}

	if response.Data == nil {
		log.Error().Msg("fkqst api: unbind account response data is nil")
		return false
	}

	if response.ErrCode != 0 {
		log.Error().Str("errMsg", response.ErrMsg).Msg("fkqst api: unbind account response error code is nil")
		return false
	}

	log.Info().
		Interface("response", response).
		Msg("fkqst api: unbind account response")

	return true
}
