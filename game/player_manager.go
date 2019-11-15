package game

import "github.com/yokaiio/yokai_server/game/player"

type PlayerManager struct {
	g         *Game
	mapPlayer map[int64]player.Player
}

func NewPlayerManager(g *Game) *PlayerManager {
	return &PlayerManager{
		g:         g,
		mapPlayer: make(map[int64]player.Player, 0),
	}
}

func (m *PlayerManager) NewPlayer(id int64, name string) player.Player {
	p := player.NewPlayer(id, name)
	m.mapPlayer[id] = p
	return p
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	return m.mapPlayer[id]
}
