package game

import (
	"hash/crc32"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// test cases
var (
	cases = []string{
		"C2M_AccountLogon",
		"C2M_StartStageCombat",
		"C2M_TakeoffRune",
		"C2M_PutonRune",
		"C2M_QueryRunes",
		"C2M_DelRune",
		"C2M_AddRune",
		"MC_QueryTalents",
		"MC_AddTalent",
		"C2M_QueryTokens",
		"C2M_AddToken",
		"C2M_TakeoffEquip",
		"C2M_PutonEquip",
		"C2M_QueryItems",
		"C2M_UseItem",
		"C2M_DelItem",
		"C2M_AddItem",
		"C2M_QueryHeros",
		"C2M_DelHero",
		"C2M_AddHero",
		"C2M_ChangeLevel",
		"C2M_ChangeExp",
		"MC_SelectPlayer",
		"C2M_CreatePlayer",
		"C2M_QueryPlayerInfo",
		"C2M_AccountDisconnect",
		"MC_AccountConnected",
		"C2M_HeartBeat",
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
