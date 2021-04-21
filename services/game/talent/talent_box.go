package talent

import (
	"errors"

	"bitbucket.org/funplus/server/define"
	"bitbucket.org/funplus/server/excel/auto"
	"bitbucket.org/funplus/server/utils"
)

var (
	ErrInvalidStar     = errors.New("invalid star")
	ErrInvalidTalentId = errors.New("invalid talent id")
)

type TalentOwner interface {
	GetTypeId() int32
}

// 天赋管理
type TalentBox struct {
	owner      TalentOwner                         `bson:"-" json:"-"`
	TalentList [define.Hero_Max_Starup_Times]int32 `bson:"talent_list" json:"talent_list"`
}

func NewTalentBox(owner TalentOwner) *TalentBox {
	m := &TalentBox{
		owner: owner,
	}

	for k := range m.TalentList {
		m.TalentList[k] = -1
	}

	return m
}

func (tb *TalentBox) GetTalentByStar(star int32) int32 {
	if !utils.BetweenInt32(star, 0, define.Hero_Max_Starup_Times) {
		return -1
	}

	return tb.TalentList[star]
}

func (tb *TalentBox) ChooseTalent(talentId int32) error {
	talentEntry, ok := auto.GetTalentEntry(talentId)
	if !ok {
		return ErrInvalidTalentId
	}

	if talentEntry.HeroTypeId != tb.owner.GetTypeId() {
		return ErrInvalidTalentId
	}

	tb.TalentList[talentEntry.Star-1] = talentId
	return nil
}
