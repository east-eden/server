package battle

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	pbGlobal "github.com/east-eden/server/proto/global"
	"google.golang.org/protobuf/proto"
)

// TestStartServerBattle 起一个假战斗服，验证请求/响应 protobuf 往返
func TestStartServerBattle(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/startBattle" {
			http.NotFound(w, r)
			return
		}

		if got := r.URL.Query().Get("battle_id"); got == "" {
			t.Errorf("battle_id query missing")
		}

		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body failed: %v", err)
			return
		}

		var req pbGlobal.ReqServerBattleStart
		if err := proto.Unmarshal(raw, &req); err != nil {
			t.Errorf("unmarshal request failed: %v", err)
			return
		}

		resp := &pbGlobal.RespServerBattleResult{
			BattleId: req.BattleId,
			Win:      true,
			Winner:   req.AtkPlayer.GetId(),
			TotalDmg: 1234.5,
		}

		out, _ := proto.Marshal(resp)
		_, _ = w.Write(out)
	}))
	defer srv.Close()

	req := BuildServerBattleRequest(
		pbGlobal.EMBattleType_BattleType_Stage,
		1,
		&pbGlobal.BattlePlayer{Id: 1001, Name: "atk"},
		nil,
		&pbGlobal.BattleExtraData{BattleIndex: 1},
	)

	result, err := StartServerBattle(srv.URL, req)
	if err != nil {
		t.Fatalf("StartServerBattle failed: %v", err)
	}

	if !result.Win {
		t.Fatal("expect win")
	}

	if result.BattleId != req.BattleId {
		t.Fatalf("battle id mismatch: %q vs %q", result.BattleId, req.BattleId)
	}

	if result.Winner != 1001 {
		t.Fatalf("winner = %d, want 1001", result.Winner)
	}
}

func TestStartServerBattleEmptyURL(t *testing.T) {
	if _, err := StartServerBattle("", &pbGlobal.ReqServerBattleStart{}); err == nil {
		t.Fatal("expect error for empty battle server url")
	}
}
