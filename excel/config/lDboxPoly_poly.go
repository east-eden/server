package config

import (
	"github.com/east-eden/server/excel"
)

// LD宝箱及掉落-拳皇/LDboxPoly
type LDRoleRaceLoot struct {
	Id     RoleRaceType `json:"id,omitempty" bson:"id"`         //角色种族
	Weight int32        `json:"weight,omitempty" bson:"weight"` //权重
}

func init() {
	excel.AddEntryNewFunc("*LDRoleRaceLoot", newLDRoleRaceLoot)

}

// 掉落接口
func (e *LDRoleRaceLoot) GetId() RoleRaceType {
	return RoleRaceType(e.Id)
}

// 掉落接口
func (e *LDRoleRaceLoot) GetWeight() int32 {
	return int32(e.Weight)
}

func newLDRoleRaceLoot() any {
	return &LDRoleRaceLoot{}
}
