package config

import (
	"github.com/east-eden/server/excel"
)

// 基础配置-拳皇/basePoly
type RandomId struct {
	Id     int32 `json:"id,omitempty" bson:"id"`         //
	Weight int32 `json:"weight,omitempty" bson:"weight"` //
}

// 基础配置-拳皇/basePoly
type RandomStr struct {
	Id     string `json:"id,omitempty" bson:"id"`         //
	Weight int32  `json:"weight,omitempty" bson:"weight"` //
}

// 基础配置-拳皇/basePoly
type IdCount struct {
	Id    int32 `json:"id,omitempty" bson:"id"`       //
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// 基础配置-拳皇/basePoly
type GoodsIdCount struct {
	Id    int32 `json:"id,omitempty" bson:"id"`       //
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// 基础配置-拳皇/basePoly
type IdCountCalc struct {
	Id    int32  `json:"id,omitempty" bson:"id"`       //数量公式
	Value string `json:"value,omitempty" bson:"value"` //
}

// 基础配置-拳皇/basePoly
type IdLevel struct {
	Id    int32 `json:"id,omitempty" bson:"id"`       //
	Count int32 `json:"count,omitempty" bson:"count"` //
}

// 基础配置-拳皇/basePoly
type BundleRes struct {
	Name   string `json:"name,omitempty" bson:"name"`     //
	Bundle string `json:"bundle,omitempty" bson:"bundle"` //
}

// 基础配置-拳皇/basePoly
type SystemOpenCondi struct {
	Type SysOpenRuleExType `json:"type,omitempty" bson:"type"` //
	Data string            `json:"data,omitempty" bson:"data"` //
}

// 基础配置-拳皇/basePoly
type AttrValue struct {
	Type  AttrType `json:"type,omitempty" bson:"type"`   //
	Value float32  `json:"value,omitempty" bson:"value"` //
}

// 基础配置-拳皇/basePoly
type AttrValueCalc struct {
	Type  AttrType `json:"type,omitempty" bson:"type"`   //属性公式
	Value string   `json:"value,omitempty" bson:"value"` //
}

// 基础配置-拳皇/basePoly
type RandomAttr struct {
	Id     AttrType `json:"id,omitempty" bson:"id"`         //
	Weight int32    `json:"weight,omitempty" bson:"weight"` //
}

// 基础配置-拳皇/basePoly
type RandomAttrValue struct {
	Id       AttrType `json:"id,omitempty" bson:"id"`             //属性枚举
	MinValue float32  `json:"minValue,omitempty" bson:"minValue"` //属性最小值
	MaxValue float32  `json:"maxValue,omitempty" bson:"maxValue"` //属性最大值
	Weight   int32    `json:"weight,omitempty" bson:"weight"`     //权重
}

// 基础配置-拳皇/basePoly
type IdPair struct {
	Id1 int32 `json:"id1,omitempty" bson:"id1"` //
	Id2 int32 `json:"id2,omitempty" bson:"id2"` //
}

// 基础配置-拳皇/basePoly
type IdInterval struct {
	Id       int32 `json:"id,omitempty" bson:"id"`             //id
	Interval int32 `json:"interval,omitempty" bson:"interval"` //间隔
}

// 基础配置-拳皇/basePoly
type ValueWeight struct {
	Id     int32 `json:"id,omitempty" bson:"id"`         //值
	Weight int32 `json:"weight,omitempty" bson:"weight"` //权重
}

// 基础配置-拳皇/basePoly
type SpecialAttr struct {
	Id      int32   `json:"id,omitempty" bson:"id"`           //道具ID
	Percent float32 `json:"percent,omitempty" bson:"percent"` //数量百分比加成
	Value   int32   `json:"value,omitempty" bson:"value"`     //数量增加数值
	Rand    int32   `json:"rand,omitempty" bson:"rand"`       //随机概率(随机到才给道具)
}

// 基础配置-拳皇/basePoly
type ChapterLevel struct {
	Chapter int32 `json:"chapter,omitempty" bson:"chapter"` //章节
	Level   int32 `json:"level,omitempty" bson:"level"`     //小关
}

// 基础配置-拳皇/basePoly
type AllianceTechAttr struct {
	PerLevel int32    `json:"perLevel,omitempty" bson:"perLevel"` //每多少级
	Type     AttrType `json:"type,omitempty" bson:"type"`         //属性类型
	Value    float32  `json:"value,omitempty" bson:"value"`       //属性值
}

// 基础配置-拳皇/basePoly
type ValuePair struct {
	Value1 int32 `json:"value1,omitempty" bson:"value1"` //
	Value2 int32 `json:"value2,omitempty" bson:"value2"` //
}

// 基础配置-拳皇/basePoly
type StringVal struct {
	Id    string `json:"id,omitempty" bson:"id"`       //
	Value int32  `json:"value,omitempty" bson:"value"` //
}

// 基础配置-拳皇/basePoly
type AttrPair struct {
	AttrType AttrType `json:"attrType,omitempty" bson:"attrType"` //
	Value    float32  `json:"value,omitempty" bson:"value"`       //
}

// 基础配置-拳皇/basePoly
type LDAttrPair struct {
	AttrType LDAttrType `json:"attrType,omitempty" bson:"attrType"` //
	Value    float32    `json:"value,omitempty" bson:"value"`       //
}

// 基础配置-拳皇/basePoly
type CrossArenaBet struct {
	CountLimit int32 `json:"countLimit,omitempty" bson:"countLimit"` //每个进度可投注场次上限
	AmountMax  int32 `json:"amountMax,omitempty" bson:"amountMax"`   //每个阶段可投注币上限
	AmountMin  int32 `json:"amountMin,omitempty" bson:"amountMin"`   //每个阶段可投注币下限
	RandCount  int32 `json:"randCount,omitempty" bson:"randCount"`   //每个阶段随机的场次
}

// 基础配置-拳皇/basePoly
type CrossArenaOddsParam struct {
	X       float32 `json:"x,omitempty" bson:"x"`             //
	M       float32 `json:"m,omitempty" bson:"m"`             //
	N       float32 `json:"n,omitempty" bson:"n"`             //
	Z       float32 `json:"z,omitempty" bson:"z"`             //
	MinOdds float32 `json:"minOdds,omitempty" bson:"minOdds"` //
	MaxOdds float32 `json:"maxOdds,omitempty" bson:"maxOdds"` //
}

// 基础配置-拳皇/basePoly
type ValueRarity struct {
	Value  int32  `json:"value,omitempty" bson:"value"`   //
	Rarity Rarity `json:"rarity,omitempty" bson:"rarity"` //
}

// 基础配置-拳皇/basePoly
type NestPos struct {
	Pos1 int32 `json:"pos1,omitempty" bson:"pos1"` //
	Pos2 int32 `json:"pos2,omitempty" bson:"pos2"` //
}

// 基础配置-拳皇/basePoly
type LDHighCoin struct {
	Day      int32 `json:"day,omitempty" bson:"day"`           //天
	LimitCnt int32 `json:"limitCnt,omitempty" bson:"limitCnt"` //限制次数
	CoinNum  int32 `json:"coinNum,omitempty" bson:"coinNum"`   //给予的代币数量
}

// 基础配置-拳皇/basePoly
type SuitAttrValue struct {
	AttrType  AttrType `json:"attrType,omitempty" bson:"attrType"`   //属性类型
	AttrValue float32  `json:"attrValue,omitempty" bson:"attrValue"` //属性值
	SuitNum   int32    `json:"suitNum,omitempty" bson:"suitNum"`     //套装件数
}

// 基础配置-拳皇/basePoly
type ValueProbability struct {
	Value       int32   `json:"value,omitempty" bson:"value"`             //值
	Probability float32 `json:"probability,omitempty" bson:"probability"` //概率
}

// 基础配置-拳皇/basePoly
type DamageAwardId struct {
	Damage  float32 `json:"damage,omitempty" bson:"damage"`   //伤害
	AwardId int32   `json:"awardId,omitempty" bson:"awardId"` //奖励ID
}

// 基础配置-拳皇/basePoly
type RandomValueWeight struct {
	MinVal float32 `json:"minVal,omitempty" bson:"minVal"` //最小值
	MaxVal float32 `json:"maxVal,omitempty" bson:"maxVal"` //最大值
	Weight int32   `json:"weight,omitempty" bson:"weight"` //权重
}

// 基础配置-拳皇/basePoly
type SkeletonRes struct {
	Name      string `json:"name,omitempty" bson:"name"`           //资源名
	Animation string `json:"animation,omitempty" bson:"animation"` //动画名
	Bundle    string `json:"bundle,omitempty" bson:"bundle"`       //
}

// 基础配置-拳皇/basePoly
type FloatPair struct {
	Val1 float32 `json:"val1,omitempty" bson:"val1"` //
	Val2 float32 `json:"val2,omitempty" bson:"val2"` //
}

// 基础配置-拳皇/basePoly
type IdFloatVal struct {
	Id  int32   `json:"id,omitempty" bson:"id"`   //
	Val float32 `json:"val,omitempty" bson:"val"` //
}

// 基础配置-拳皇/basePoly
type SoulCustomHideAttr struct {
	Id     int32   `json:"id,omitempty" bson:"id"`         //
	Rarity Rarity  `json:"rarity,omitempty" bson:"rarity"` //
	Value  float32 `json:"value,omitempty" bson:"value"`   //
}

// 基础配置-拳皇/basePoly
type SoulStarUpMaterial struct {
	SoulType HeroSoulType `json:"soulType,omitempty" bson:"soulType"` //
	Star     int32        `json:"star,omitempty" bson:"star"`         //
	Count    int32        `json:"count,omitempty" bson:"count"`       //
}

// 基础配置-拳皇/basePoly
type SoulSkillRequire struct {
	Num     int32 `json:"num,omitempty" bson:"num"`         //
	Star    int32 `json:"star,omitempty" bson:"star"`       //
	SkillId int32 `json:"skillId,omitempty" bson:"skillId"` //
}

// 基础配置-拳皇/basePoly
type ItemLoot struct {
	LootId  int32 `json:"lootId,omitempty" bson:"lootId"`   //
	ItemId  int32 `json:"itemId,omitempty" bson:"itemId"`   //
	ItemNum int32 `json:"itemNum,omitempty" bson:"itemNum"` //
	Weight  int32 `json:"weight,omitempty" bson:"weight"`   //
}

func init() {
	excel.AddEntryNewFunc("*RandomId", newRandomId)
	excel.AddEntryNewFunc("*RandomStr", newRandomStr)
	excel.AddEntryNewFunc("*IdCount", newIdCount)
	excel.AddEntryNewFunc("*GoodsIdCount", newGoodsIdCount)
	excel.AddEntryNewFunc("*IdCountCalc", newIdCountCalc)
	excel.AddEntryNewFunc("*IdLevel", newIdLevel)
	excel.AddEntryNewFunc("*BundleRes", newBundleRes)
	excel.AddEntryNewFunc("*SystemOpenCondi", newSystemOpenCondi)
	excel.AddEntryNewFunc("*AttrValue", newAttrValue)
	excel.AddEntryNewFunc("*AttrValueCalc", newAttrValueCalc)
	excel.AddEntryNewFunc("*RandomAttr", newRandomAttr)
	excel.AddEntryNewFunc("*RandomAttrValue", newRandomAttrValue)
	excel.AddEntryNewFunc("*IdPair", newIdPair)
	excel.AddEntryNewFunc("*IdInterval", newIdInterval)
	excel.AddEntryNewFunc("*ValueWeight", newValueWeight)
	excel.AddEntryNewFunc("*SpecialAttr", newSpecialAttr)
	excel.AddEntryNewFunc("*ChapterLevel", newChapterLevel)
	excel.AddEntryNewFunc("*AllianceTechAttr", newAllianceTechAttr)
	excel.AddEntryNewFunc("*ValuePair", newValuePair)
	excel.AddEntryNewFunc("*StringVal", newStringVal)
	excel.AddEntryNewFunc("*AttrPair", newAttrPair)
	excel.AddEntryNewFunc("*LDAttrPair", newLDAttrPair)
	excel.AddEntryNewFunc("*CrossArenaBet", newCrossArenaBet)
	excel.AddEntryNewFunc("*CrossArenaOddsParam", newCrossArenaOddsParam)
	excel.AddEntryNewFunc("*ValueRarity", newValueRarity)
	excel.AddEntryNewFunc("*NestPos", newNestPos)
	excel.AddEntryNewFunc("*LDHighCoin", newLDHighCoin)
	excel.AddEntryNewFunc("*SuitAttrValue", newSuitAttrValue)
	excel.AddEntryNewFunc("*ValueProbability", newValueProbability)
	excel.AddEntryNewFunc("*DamageAwardId", newDamageAwardId)
	excel.AddEntryNewFunc("*RandomValueWeight", newRandomValueWeight)
	excel.AddEntryNewFunc("*SkeletonRes", newSkeletonRes)
	excel.AddEntryNewFunc("*FloatPair", newFloatPair)
	excel.AddEntryNewFunc("*IdFloatVal", newIdFloatVal)
	excel.AddEntryNewFunc("*SoulCustomHideAttr", newSoulCustomHideAttr)
	excel.AddEntryNewFunc("*SoulStarUpMaterial", newSoulStarUpMaterial)
	excel.AddEntryNewFunc("*SoulSkillRequire", newSoulSkillRequire)
	excel.AddEntryNewFunc("*ItemLoot", newItemLoot)

}

// 掉落接口
func (e *RandomId) GetId() int32 {
	return int32(e.Id)
}

// 掉落接口
func (e *RandomId) GetWeight() int32 {
	return int32(e.Weight)
}

// 掉落接口
func (e *RandomStr) GetId() string {
	return string(e.Id)
}

// 掉落接口
func (e *RandomStr) GetWeight() int32 {
	return int32(e.Weight)
}

// 掉落接口
func (e *RandomAttr) GetId() int32 {
	return int32(e.Id)
}

// 掉落接口
func (e *RandomAttr) GetWeight() int32 {
	return int32(e.Weight)
}

// 掉落接口
func (e *RandomAttrValue) GetId() AttrType {
	return AttrType(e.Id)
}

// 掉落接口
func (e *RandomAttrValue) GetWeight() int32 {
	return int32(e.Weight)
}

// 掉落接口
func (e *ValueWeight) GetId() int32 {
	return int32(e.Id)
}

// 掉落接口
func (e *ValueWeight) GetWeight() int32 {
	return int32(e.Weight)
}

func newRandomId() any {
	return &RandomId{}
}

func newRandomStr() any {
	return &RandomStr{}
}

func newIdCount() any {
	return &IdCount{}
}

func newGoodsIdCount() any {
	return &GoodsIdCount{}
}

func newIdCountCalc() any {
	return &IdCountCalc{}
}

func newIdLevel() any {
	return &IdLevel{}
}

func newBundleRes() any {
	return &BundleRes{}
}

func newSystemOpenCondi() any {
	return &SystemOpenCondi{}
}

func newAttrValue() any {
	return &AttrValue{}
}

func newAttrValueCalc() any {
	return &AttrValueCalc{}
}

func newRandomAttr() any {
	return &RandomAttr{}
}

func newRandomAttrValue() any {
	return &RandomAttrValue{}
}

func newIdPair() any {
	return &IdPair{}
}

func newIdInterval() any {
	return &IdInterval{}
}

func newValueWeight() any {
	return &ValueWeight{}
}

func newSpecialAttr() any {
	return &SpecialAttr{}
}

func newChapterLevel() any {
	return &ChapterLevel{}
}

func newAllianceTechAttr() any {
	return &AllianceTechAttr{}
}

func newValuePair() any {
	return &ValuePair{}
}

func newStringVal() any {
	return &StringVal{}
}

func newAttrPair() any {
	return &AttrPair{}
}

func newLDAttrPair() any {
	return &LDAttrPair{}
}

func newCrossArenaBet() any {
	return &CrossArenaBet{}
}

func newCrossArenaOddsParam() any {
	return &CrossArenaOddsParam{}
}

func newValueRarity() any {
	return &ValueRarity{}
}

func newNestPos() any {
	return &NestPos{}
}

func newLDHighCoin() any {
	return &LDHighCoin{}
}

func newSuitAttrValue() any {
	return &SuitAttrValue{}
}

func newValueProbability() any {
	return &ValueProbability{}
}

func newDamageAwardId() any {
	return &DamageAwardId{}
}

func newRandomValueWeight() any {
	return &RandomValueWeight{}
}

func newSkeletonRes() any {
	return &SkeletonRes{}
}

func newFloatPair() any {
	return &FloatPair{}
}

func newIdFloatVal() any {
	return &IdFloatVal{}
}

func newSoulCustomHideAttr() any {
	return &SoulCustomHideAttr{}
}

func newSoulStarUpMaterial() any {
	return &SoulStarUpMaterial{}
}

func newSoulSkillRequire() any {
	return &SoulSkillRequire{}
}

func newItemLoot() any {
	return &ItemLoot{}
}
