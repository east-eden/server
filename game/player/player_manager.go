package player

type PlayerManager struct {
	mapPlayer map[int64]Player
}

func NewPlayerManager() PlayerManager {
	return &PlayerManager{
		mapPlayer: make(map[int64]Player, 0),
	}
}

func (m *PlayerManager) GetPlayerByID(id int64) (Player, bool) {
	return m.mapPlayer[id]
}
