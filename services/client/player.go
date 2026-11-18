package client

import (
	"bytes"
	"encoding/binary"

	"github.com/bits-and-blooms/bitset"
	pbGlobal "github.com/east-eden/server/proto/global"
	log "github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type Player struct {
	c               *Client
	Info            *pbGlobal.PlayerInfo
	Heros           map[int64]*pbGlobal.Hero
	Items           map[int64]*pbGlobal.Item
	Collections     map[int32]*pbGlobal.Collection
	HeroChips       map[int32]int32
	CollectionFrags map[int32]*pbGlobal.Fragment
	Chapters        map[int32]*pbGlobal.Chapter
	Stages          map[int32]*pbGlobal.Stage
	GuideInfo       *bitset.BitSet
	Quests          map[int32]*pbGlobal.Quest
	Towers          map[int32]*pbGlobal.Tower
}

func NewPlayer(ctx *cli.Context, c *Client) *Player {
	return &Player{
		c:               c,
		Heros:           make(map[int64]*pbGlobal.Hero),
		Items:           make(map[int64]*pbGlobal.Item),
		Collections:     make(map[int32]*pbGlobal.Collection),
		HeroChips:       make(map[int32]int32),
		CollectionFrags: make(map[int32]*pbGlobal.Fragment),
		Chapters:        make(map[int32]*pbGlobal.Chapter),
		Stages:          make(map[int32]*pbGlobal.Stage),
		Quests:          make(map[int32]*pbGlobal.Quest),
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

	for _, c := range msg.GetCollections() {
		p.Collections[c.TypeId] = c
	}

	for chipId, num := range msg.GetHeroChips() {
		p.HeroChips[chipId] = num
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

	for _, t := range msg.GetTowers() {
		p.Towers[t.Type] = t
	}
}
