package game

import (
	logger "github.com/sirupsen/logrus"
	"github.com/yokaiio/yokai_server/game/player"
)

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

	// add hero
	hero := p.HeroManager().NewHero(1)
	item := p.ItemManager().NewItem(1)

	heroEntry := hero.Entry()
	itemEntry := item.Entry()

	logger.Println("heroEntry:", heroEntry)
	logger.Println("itemEntry:", itemEntry)
	return p
}

func (m *PlayerManager) GetPlayerByID(id int64) player.Player {
	return m.mapPlayer[id]
}
