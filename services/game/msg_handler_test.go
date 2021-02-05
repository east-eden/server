package game

import (
	"hash/crc32"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// test cases
var (
	cases = []string{
		"C2S_AccountLogon",
		"C2S_StartStageCombat",
		"C2S_TakeoffRune",
		"C2S_PutonRune",
		"C2S_QueryRunes",
		"C2S_DelRune",
		"C2S_AddRune",
		"C2S_QueryTalents",
		"C2S_AddTalent",
		"C2S_QueryTokens",
		"C2S_AddToken",
		"C2S_TakeoffEquip",
		"C2S_PutonEquip",
		"C2S_QueryItems",
		"C2S_UseItem",
		"C2S_DelItem",
		"C2S_AddItem",
		"C2S_QueryHeros",
		"C2S_DelHero",
		"C2S_AddHero",
		"C2S_ChangeLevel",
		"C2S_ChangeExp",
		"MC_SelectPlayer",
		"C2S_CreatePlayer",
		"C2S_QueryPlayerInfo",
		"C2S_AccountDisconnect",
		"MC_AccountConnected",
		"C2S_HeartBeat",
	}
)

func TestMsgHandler(t *testing.T) {
	handler := NewMsgHandler(nil)

	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			h, err := handler.r.GetHandler(crc32.ChecksumIEEE([]byte(name)))
			if err != nil {
				t.Errorf("MsgHandler get handler failed:%v", err)
			}

			diff := cmp.Diff(h.Name, name)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}

}
