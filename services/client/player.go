package client

import (
	"bytes"
	"encoding/binary"

	pbGlobal "e.coding.net/mmstudio/blade/server/proto/global"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/willf/bitset"
)

type Player struct {
	c               *Client
	Info            *pbGlobal.PlayerInfo
	Heros           map[int64]*pbGlobal.Hero
	Items           map[int64]*pbGlobal.Item
	Equips          map[int64]*pbGlobal.Equip
	Crystals        map[int64]*pbGlobal.Crystal
	Collections     map[int32]*pbGlobal.Collection
	HeroFrags       map[int32]*pbGlobal.Fragment
	CollectionFrags map[int32]*pbGlobal.Fragment
	Chapters        map[int32]*pbGlobal.Chapter
	Stages          map[int32]*pbGlobal.Stage
	GuideInfo       *bitset.BitSet
	Quests          map[int32]*pbGlobal.Quest
	Tokens          map[int32]*pbGlobal.Token
	Towers          map[int32]*pbGlobal.Tower
}

func NewPlayer(ctx *cli.Context, c *Client) *Player {
	return &Player{
		c:               c,
		Heros:           make(map[int64]*pbGlobal.Hero),
		Items:           make(map[int64]*pbGlobal.Item),
		Equips:          make(map[int64]*pbGlobal.Equip),
		Crystals:        make(map[int64]*pbGlobal.Crystal),
		Collections:     make(map[int32]*pbGlobal.Collection),
		HeroFrags:       make(map[int32]*pbGlobal.Fragment),
		CollectionFrags: make(map[int32]*pbGlobal.Fragment),
		Chapters:        make(map[int32]*pbGlobal.Chapter),
		Stages:          make(map[int32]*pbGlobal.Stage),
		Quests:          make(map[int32]*pbGlobal.Quest),
		Tokens:          make(map[int32]*pbGlobal.Token),
		Towers:          make(map[int32]*pbGlobal.Tower),
	}
}

func (p *Player) InitInfo(msg *pbGlobal.S2C_PlayerInitInfo) {
	p.Info = msg.GetInfo()

	for _, h := range msg.GetHeros() {
		p.Heros[h.Id] = h
	}

	for _, i := range msg.GetItems() {
		p.Items[i.Id] = i
	}

	for _, e := range msg.GetEquips() {
		p.Equips[e.Item.Id] = e
	}

	for _, c := range msg.GetCrystals() {
		p.Crystals[c.Item.Id] = c
	}

	for _, c := range msg.GetCollections() {
		p.Collections[c.TypeId] = c
	}

	for _, f := range msg.GetHeroFrags() {
		p.HeroFrags[f.Id] = f
	}

	for _, f := range msg.GetCollectionFrags() {
		p.CollectionFrags[f.Id] = f
	}

	for _, c := range msg.GetChapters() {
		p.Chapters[c.Id] = c
	}

	for _, s := range msg.GetStages() {
		p.Stages[s.Id] = s
	}

	if len(msg.GetGuideInfo())%8 != 0 {
		log.Fatal().Interface("guide_info", msg.GetGuideInfo()).Msg("invalid guide info")
	}

	buf := new(bytes.Buffer)
	guideInfo := msg.GetGuideInfo()
	guideBuffer := make([]uint64, len(guideInfo)/8)
	for n := range guideBuffer {
		buf.Reset()
		_ = binary.Write(buf, binary.LittleEndian, guideInfo[n*8:(n+1)*8])
		guideBuffer[n] = binary.LittleEndian.Uint64(buf.Bytes())
	}
	p.GuideInfo = bitset.From(guideBuffer)

	for _, q := range msg.GetQuests() {
		p.Quests[q.Id] = q
	}

	for _, t := range msg.GetTokens() {
		p.Tokens[int32(t.Type)] = t
	}

	for _, t := range msg.GetTowers() {
		p.Towers[t.Type] = t
	}
}
