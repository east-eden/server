// Package battle 负责与外部战斗服通信。
//
// 与 kof 的实现保持一致：服务器把战斗请求以 protobuf 序列化后 POST 给战斗服，
// 战斗服跑完模拟后回写结果，服务器再用结果做校验与结算。
package battle

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	pbGlobal "github.com/east-eden/server/proto/global"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// ServerBattleHTTPTimeout 请求战斗服的超时时间。
// 战斗服不返回时若没有超时，调用方会被永久阻塞。
const ServerBattleHTTPTimeout = 30 * time.Second

var (
	ErrBattleServerNotConfigured = errors.New("battle server url is empty")
	ErrBattleServerResponse      = errors.New("battle server response invalid")

	serverBattleClient = &http.Client{Timeout: ServerBattleHTTPTimeout}
)

// BuildServerBattleRequest 构建战斗请求
func BuildServerBattleRequest(
	battleType pbGlobal.EMBattleType,
	levelId int32,
	atkPlayer *pbGlobal.BattlePlayer,
	defPlayer *pbGlobal.BattlePlayer,
	extra *pbGlobal.BattleExtraData,
) *pbGlobal.ReqServerBattleStart {
	if extra == nil {
		extra = &pbGlobal.BattleExtraData{}
	}

	return &pbGlobal.ReqServerBattleStart{
		BattleId:   uuid.New().String(),
		BattleType: battleType,
		LevelId:    levelId,
		AtkPlayer:  atkPlayer,
		DefPlayer:  defPlayer,
		ExtraData:  extra,
	}
}

// StartServerBattle 把战斗请求发给战斗服并解析结果。
// baseURL 形如 "http://127.0.0.1:8088/"，内部会拼上 startBattle。
func StartServerBattle(baseURL string, req *pbGlobal.ReqServerBattleStart) (*pbGlobal.RespServerBattleResult, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, ErrBattleServerNotConfigured
	}

	if req == nil {
		return nil, fmt.Errorf("%w: nil request", ErrBattleServerResponse)
	}

	url := strings.TrimSuffix(baseURL, "/") + "/startBattle"

	body, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal battle request failed: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new battle request failed: %w", err)
	}

	q := httpReq.URL.Query()
	q.Add("battle_id", req.BattleId)
	httpReq.URL.RawQuery = q.Encode()
	httpReq.Header.Set("Content-Type", "application/octet-stream")

	resp, err := serverBattleClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("post battle request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read battle response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status=%d body=%s", ErrBattleServerResponse, resp.StatusCode, string(raw))
	}

	result := &pbGlobal.RespServerBattleResult{}
	if err := proto.Unmarshal(raw, result); err != nil {
		return nil, fmt.Errorf("unmarshal battle response failed: %w", err)
	}

	return result, nil
}
