package hero

import (
	"fmt"

	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/internal/attrcache"
	"github.com/east-eden/server/utils"
	"github.com/spf13/cast"
)

// Card 英雄卡牌。
//
// 采用 kof 的模型：一个英雄配置 id 对应一张卡，等级/星级/晋升/品质决定其属性，
// 不再有“每个英雄一个物品实例”的概念。
type Card struct {
	ConfigId    int32 `bson:"config_id" json:"config_id"`       // 英雄配置 id
	Level       int32 `bson:"level" json:"level"`               // 等级
	Star        int32 `bson:"star" json:"star"`                 // 星级
	PromotionLv int32 `bson:"promotion_lv" json:"promotion_lv"` // 晋升等级
	Exp         int32 `bson:"exp" json:"exp"`                   // 经验
	Rarity      int32 `bson:"rarity" json:"rarity"`             // 品质（可被晋升提升）

	Attrs       map[int32]float32 `bson:"attrs" json:"attrs"` // 战斗属性（每次按配置重算）
	BattlePower int64             `bson:"battle_power" json:"battle_power"`
}

// NewCard 创建一个 1 级卡牌
func NewCard(configId int32) (*Card, error) {
	entry, ok := config.GetHeroEntry(configId)
	if !ok {
		return nil, fmt.Errorf("hero entry<%d> not found", configId)
	}

	return &Card{
		ConfigId: configId,
		Level:    1,
		Star:     1,
		Rarity:   int32(entry.Rarity),
		Attrs:    make(map[int32]float32),
	}, nil
}

// Entry 返回英雄配置
func (c *Card) Entry() (*config.HeroEntry, bool) {
	return config.GetHeroEntry(c.ConfigId)
}

// attrField 把属性类型映射到 BattlePropertyEntry 上的公式字段
type attrField struct {
	Attr config.AttrType
	Get  func(*config.BattlePropertyEntry) string
}

// coreAttrFields 基础属性映射。
// kof 的属性表有 80+ 个字段，MVP 先覆盖战斗力与战斗结算实际用到的核心属性；
// 其余属性在对应玩法（战魂/神器/音律等）接入时再补。
var coreAttrFields = []attrField{
	{config.AttrType_HP, func(e *config.BattlePropertyEntry) string { return e.MAXHP }},
	{config.AttrType_MAXHP, func(e *config.BattlePropertyEntry) string { return e.MAXHP }},
	{config.AttrType_ATTACK, func(e *config.BattlePropertyEntry) string { return e.ATTACK }},
	{config.AttrType_DEFENCE, func(e *config.BattlePropertyEntry) string { return e.DEFENCE }},
	{config.AttrType_ATTACKSPEED, func(e *config.BattlePropertyEntry) string { return e.ATTACKSPEED }},
	{config.AttrType_MOVESPEED, func(e *config.BattlePropertyEntry) string { return e.MOVESPEED }},
	{config.AttrType_ATTACKRANGE, func(e *config.BattlePropertyEntry) string { return e.ATTACKRANGE }},
	{config.AttrType_CRITAL, func(e *config.BattlePropertyEntry) string { return e.CRITAL }},
	{config.AttrType_IGNCRITAL, func(e *config.BattlePropertyEntry) string { return e.IGNCRITAL }},
	{config.AttrType_DODGE, func(e *config.BattlePropertyEntry) string { return e.DODGE }},
	{config.AttrType_IGNDODGE, func(e *config.BattlePropertyEntry) string { return e.IGNDODGE }},
	{config.AttrType_PARRY, func(e *config.BattlePropertyEntry) string { return e.PARRY }},
	{config.AttrType_BREAKPARRY, func(e *config.BattlePropertyEntry) string { return e.BREAKPARRY }},
	{config.AttrType_ENDURE, func(e *config.BattlePropertyEntry) string { return e.ENDURE }},
	{config.AttrType_HITRATE, func(e *config.BattlePropertyEntry) string { return e.HITRATE }},
	{config.AttrType_MAXSP, func(e *config.BattlePropertyEntry) string { return e.MAXSP }},
	{config.AttrType_SPINITIAL, func(e *config.BattlePropertyEntry) string { return e.SPINITIAL }},
	{config.AttrType_SPRESERVES, func(e *config.BattlePropertyEntry) string { return e.SPRESERVES }},
}

// BasePropertyId 按品质档位与晋升等级取战斗属性表 id。
// kof 的映射是 PropertyIds[rarity - 初始品质][promotionLv]。
func (c *Card) BasePropertyId() (int32, bool) {
	entry, ok := c.Entry()
	if !ok {
		return 0, false
	}

	idx := c.Rarity - int32(entry.Rarity)
	if idx < 0 || idx >= int32(len(entry.PropertyIds)) {
		return 0, false
	}

	row := entry.PropertyIds[idx]
	if int(c.PromotionLv) < 0 || int(c.PromotionLv) >= len(row) {
		return 0, false
	}

	return row[c.PromotionLv], true
}

// battlePropertyAt 返回该等级命中的属性档位（最高的 Level <= card.Level）
func (c *Card) battlePropertyAt(propertyId int32) (*config.BattlePropertyEntry, bool) {
	var matched *config.BattlePropertyEntry
	for _, entry := range config.GetBattlePropertiesById(propertyId) {
		if c.Level >= entry.Level {
			matched = entry
			continue
		}

		break
	}

	return matched, matched != nil
}

// CalcBaseAttrs 按当前等级/品质/晋升重算基础属性
func (c *Card) CalcBaseAttrs(playerLevel int32) error {
	propertyId, ok := c.BasePropertyId()
	if !ok {
		return fmt.Errorf("hero<%d> has no valid property id (rarity=%d promotion=%d)", c.ConfigId, c.Rarity, c.PromotionLv)
	}

	entry, ok := c.battlePropertyAt(propertyId)
	if !ok {
		return fmt.Errorf("battle property<%d> level<%d> not found", propertyId, c.Level)
	}

	attrs := make(map[int32]float32, len(coreAttrFields))
	for _, field := range coreAttrFields {
		formula := field.Get(entry)
		if formula == "" {
			continue
		}

		attrs[int32(field.Attr)] = float32(attrcache.EvalWithVars(formula, c.Level, playerLevel))
	}

	c.Attrs = attrs
	return nil
}

// GetAttr 读取某项属性
func (c *Card) GetAttr(attrType config.AttrType) float32 {
	if c.Attrs == nil {
		return 0
	}

	return c.Attrs[int32(attrType)]
}

// applyAttrValues 把配置里的属性列表累加进卡牌属性
func (c *Card) applyAttrValues(values []*config.AttrValue) {
	for _, v := range values {
		if v == nil {
			continue
		}

		c.Attrs[int32(v.Type)] += v.Value
	}
}

// CalcStarAttrs 叠加星级属性。
//
// 口径与 kof 一致：先按职业取 HeroStar 表里的职业专属属性，
// 再叠加通用的 Attrs / ExtraAttrs。
func (c *Card) CalcStarAttrs() {
	entry, ok := c.Entry()
	if !ok {
		return
	}

	// HeroStar 的复合键是 "<星级>_<品质名>"，品质取枚举名字（如 Red）
	starEntry, ok := config.GetHeroStarEntry(
		cast.ToString(c.Star),
		entry.Rarity.String(),
	)
	if !ok {
		return
	}

	switch entry.Career {
	case config.Career_Attack:
		c.applyAttrValues(starEntry.AttackAttrs)
	case config.Career_Defense:
		c.applyAttrValues(starEntry.DefenseAttrs)
	case config.Career_Skill:
		c.applyAttrValues(starEntry.SkillAttrs)
	case config.Career_Auxiliary:
		c.applyAttrValues(starEntry.AuxiliaryAttrs)
	}

	c.applyAttrValues(starEntry.Attrs)
	c.applyAttrValues(starEntry.ExtraAttrs)
}

// CalcPower 计算战力：按 Attr 表的战力权重加权求和
func (c *Card) CalcPower() int64 {
	attrRows := config.GetAttrRows()

	var power float32
	for attrType, val := range c.Attrs {
		row, ok := attrRows[config.AttrType(attrType)]
		if !ok {
			continue
		}

		power += val * row.Value
	}

	c.BattlePower = int64(power)
	return c.BattlePower
}

// Recalc 重算完整属性链：基础属性 → 星级加成 → 百分比/最终属性折算 → 战力
func (c *Card) Recalc(playerLevel int32) error {
	if err := c.CalcBaseAttrs(playerLevel); err != nil {
		return err
	}

	c.CalcStarAttrs()

	c.Attrs = utils.CalcPercentAttribute(
		c.Attrs,
		int32(config.AttrType_PERCENT_BEGIN),
		int32(config.AttrType_PERCENT_END),
	)

	c.Attrs = utils.CalcFinalAttribute(
		c.Attrs,
		int32(config.AttrType_FINAL_BEGIN),
		int32(config.AttrType_END),
	)

	c.CalcPower()
	return nil
}
