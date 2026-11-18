package config

// LD收藏表-拳皇/LDRoomEventEnum
type RoomEventId int

const (
	RoomEventId_BEGIN              RoomEventId = 1  //枚举开始
	RoomEventId_LivingTowardsDeath RoomEventId = 13 //向死而生
	RoomEventId_END                RoomEventId = 14 //枚举结束
)

func (e RoomEventId) String() string {

	if e == RoomEventId_LivingTowardsDeath {
		return "LivingTowardsDeath"
	}

	return ""

}

func GetRoomEventId(tp string) (RoomEventId, bool) {
	switch tp {
	case "LivingTowardsDeath":
		return RoomEventId_LivingTowardsDeath, true

	}
	return -1, false

}

func GetRoomEventId_ByAlias(tp string) (RoomEventId, bool) {
	switch tp {
	case "向死而生":
		return RoomEventId_LivingTowardsDeath, true

	}
	return -1, false

}
