package game

import (
	"hash/crc32"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// test cases
//
// 该用例保证 proto 消息名到 crc32 handler id 的映射保持稳定：
// 消息名一旦被改名，客户端与服务器之间的消息 id 会失配。
var (
	cases = []string{
		// account
		"C2S_WaitResponseMessage",
		"C2S_Ping",
		"C2S_AccountLogon",
		"C2S_HeartBeat",
		"C2S_AccountDisconnect",

		// player
		"C2S_CreatePlayer",
		"C2S_GmCmd",
		"C2S_GuidePass",
		"C2S_SaveBattleArray",

		// heros
		"C2S_DelHero",
		"C2S_HeroPromote",
		"C2S_HeroStarup",

		// fragment
		"C2S_HeroFragmentsCompose",
		"C2S_CollectionFragmentsCompose",

		// items & equips & crystals
		"C2S_DelItem",
		"C2S_UseItem",

		// collections
		"C2S_CollectionActive",
		"C2S_CollectionStarup",
		"C2S_CollectionWakeup",

		// chapter & stage
		"C2S_StageSweep",
		"C2S_StageChallenge",
		"C2S_BattleResult",
		"C2S_ActivitySignIn",
		"C2S_ActivityClaimTask",
		"C2S_ShopBuy",
		"C2S_Recruit",
		"C2S_BattlePassClaim",
		"C2S_BattlePassUnlock",
		"C2S_RechargeCreate",
		"C2S_FirstRechargeClaim",

		// tower
		"C2S_TowerChallenge",

		// quest
		"C2S_PlayerQuestReward",

		// rank
		"C2S_QueryRank",
	}
)

func TestMsgHandler(t *testing.T) {
	handler := NewMsgRegister(nil, nil, nil)

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := handler.r.GetHandler(crc32.ChecksumIEEE([]byte(name)))
			if err != nil {
				t.Fatalf("MsgHandler get handler failed: %v", err)
			}

			if diff := cmp.Diff(h.Name, name); diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
