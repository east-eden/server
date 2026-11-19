package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var allianceConfigEntries *AllianceConfigEntries //AllianceConfig.json全局变量

// 属性表 联盟-拳皇/AllianceConfig
type AllianceConfigEntry struct {
	Test                                 int32                                      `json:"test,omitempty" bson:"test"`                                                                 //##测试数据
	AllianceNameMaxLimit                 int32                                      `json:"allianceNameMaxLimit,omitempty" bson:"allianceNameMaxLimit"`                                 //公会昵称字数最大值
	AllianceDescMaxLimit                 int32                                      `json:"allianceDescMaxLimit,omitempty" bson:"allianceDescMaxLimit"`                                 //公会介绍字数最大值(公告)
	CreateCost                           *IdCount                                   `json:"createCost,omitempty" bson:"createCost"`                                                     //创建公会需要消耗道具
	MaxDeclarationsWords                 int32                                      `json:"maxDeclarationsWords,omitempty" bson:"maxDeclarationsWords"`                                 //公会宣言最大字数
	AllianceIntroduceDesc                string                                     `json:"allianceIntroduceDesc,omitempty" bson:"allianceIntroduceDesc"`                               //公会默认介绍
	AllianceQrCodeUploadLimit            int32                                      `json:"allianceQrCodeUploadLimit,omitempty" bson:"allianceQrCodeUploadLimit"`                       //公会二维码每天上传次数限制
	AllianceQrCodeLoseTime               int32                                      `json:"allianceQrCodeLoseTime,omitempty" bson:"allianceQrCodeLoseTime"`                             //公会二维码失效时间，秒
	AllianceDefaultFlag                  int32                                      `json:"allianceDefaultFlag,omitempty" bson:"allianceDefaultFlag"`                                   //公会默认旗帜
	AllianceDefaultPattern               int32                                      `json:"allianceDefaultPattern,omitempty" bson:"allianceDefaultPattern"`                             //公会默认图标
	AllianceQrCodeCardUrl                string                                     `json:"allianceQrCodeCardUrl,omitempty" bson:"allianceQrCodeCardUrl"`                               //公会qrcode卡片链接
	AllianceApplyLimit                   int32                                      `json:"allianceApplyLimit,omitempty" bson:"allianceApplyLimit"`                                     //公会申请人数上限
	AllianceApplyExpire                  int32                                      `json:"allianceApplyExpire,omitempty" bson:"allianceApplyExpire"`                                   //公会申请时效（小时）
	BuildCost                            []*IdCount                                 `json:"buildCost,omitempty" bson:"buildCost"`                                                       //建设公会消耗
	DailyMaxBuildTimes                   int32                                      `json:"dailyMaxBuildTimes,omitempty" bson:"dailyMaxBuildTimes"`                                     //每日公会建设上限
	BuildAwards                          []*IdCount                                 `json:"buildAwards,omitempty" bson:"buildAwards"`                                                   //建设公会奖励
	AutomaticTransfeCD                   int32                                      `json:"automaticTransfeCD,omitempty" bson:"automaticTransfeCD"`                                     //会长不上线自动转让时间（天）
	ExitCD                               int32                                      `json:"exitCD,omitempty" bson:"exitCD"`                                                             //公会退出冷却时间（小时）
	ChangeNameCost                       *IdCount                                   `json:"changeNameCost,omitempty" bson:"changeNameCost"`                                             //修改公会名字消耗
	MaxJoinLevelOptions                  []int32                                    `json:"maxJoinLevelOptions,omitempty" bson:"maxJoinLevelOptions"`                                   //公会加入等级限制选项
	KickOutCD                            int32                                      `json:"kickOutCD,omitempty" bson:"kickOutCD"`                                                       //被踢出公会冷却时间（小时）
	KickWaitTime                         int32                                      `json:"kickWaitTime,omitempty" bson:"kickWaitTime"`                                                 //被踢时如果资源很多延时踢出时间(小时)
	DailyChallengeMaxFreeTimes           int32                                      `json:"dailyChallengeMaxFreeTimes,omitempty" bson:"dailyChallengeMaxFreeTimes"`                     //boss每日免费挑战次数
	DailyChallengeMaxCostTimes           int32                                      `json:"dailyChallengeMaxCostTimes,omitempty" bson:"dailyChallengeMaxCostTimes"`                     //boss每日消耗道具挑战次数
	DailyChallengeCostCount              []int32                                    `json:"dailyChallengeCostCount,omitempty" bson:"dailyChallengeCostCount"`                           //boss每日挑战次数消耗道具数量
	DailyChallengeCostId                 int32                                      `json:"dailyChallengeCostId,omitempty" bson:"dailyChallengeCostId"`                                 //boss挑战消耗
	BossStartLevelId                     int32                                      `json:"bossStartLevelId,omitempty" bson:"bossStartLevelId"`                                         //boss关卡初始id
	DailyChallengeAwards                 []*IdCount                                 `json:"dailyChallengeAwards,omitempty" bson:"dailyChallengeAwards"`                                 //boss每次攻击奖励道具
	AllianceLevelUpCost                  int32                                      `json:"allianceLevelUpCost,omitempty" bson:"allianceLevelUpCost"`                                   //公会升级消耗道具
	AllianceBargainGoods                 [][]int32                                  `json:"allianceBargainGoods,omitempty" bson:"allianceBargainGoods"`                                 //公会砍价商品库（掉落id，随机权重，消耗钻石）
	AllianceBargainDiscount              [][]float32                                `json:"allianceBargainDiscount,omitempty" bson:"allianceBargainDiscount"`                           //公会单次砍价折扣区间(人数,区间下限,区间上限)
	NewBargainDiscountServerOpenDays     int32                                      `json:"newBargainDiscountServerOpenDays,omitempty" bson:"newBargainDiscountServerOpenDays"`         //公会所在服开服第x天使用新的砍价区间
	NewAllianceBargainDiscount           [][]float32                                `json:"newAllianceBargainDiscount,omitempty" bson:"newAllianceBargainDiscount"`                     //开服第x天新的公会单次砍价折扣区间(人数,区间下限,区间上限)
	AllianceBossKillExp                  [][]int32                                  `json:"allianceBossKillExp,omitempty" bson:"allianceBossKillExp"`                                   //工会boss击杀后参与人数对应雕像经验
	AllianceTechFreeTimes                int32                                      `json:"allianceTechFreeTimes,omitempty" bson:"allianceTechFreeTimes"`                               //公会科技每天免费次数
	AllianceTechTicketCost               *IdCount                                   `json:"allianceTechTicketCost,omitempty" bson:"allianceTechTicketCost"`                             //公会科技每次抽取消耗代币
	AllianceTechFrameGot                 map[Rarity]int32                           `json:"allianceTechFrameGot,omitempty" bson:"allianceTechFrameGot"`                                 //公会科技抽中紫色、橙色、红色科技时获得对应碎片数量
	AllianceTechTicketGot                map[Rarity]int32                           `json:"allianceTechTicketGot,omitempty" bson:"allianceTechTicketGot"`                               //公会科技抽中紫色、橙色、红色科技时获得代币数量
	AllianceTechEquipNum                 map[int32]int32                            `json:"allianceTechEquipNum,omitempty" bson:"allianceTechEquipNum"`                                 //公会科技全部品质总等级和可装备科技数量
	AllianceTechTotalAttr                map[Rarity]*AllianceTechAttr               `json:"allianceTechTotalAttr,omitempty" bson:"allianceTechTotalAttr"`                               //公会科技X1品质总等级每X2级提升全体X3属性值X4
	TechRaritys                          []Rarity                                   `json:"techRaritys,omitempty" bson:"techRaritys"`                                                   //工会科技所有品质
	AlliaceTechTicketThreshold           int32                                      `json:"alliaceTechTicketThreshold,omitempty" bson:"alliaceTechTicketThreshold"`                     //工会科技持有阈值（大于阈值时踢人会有处理）
	AlliancePowerPositions               map[AlliancePowerType]AlliancePositionType `json:"alliancePowerPositions,omitempty" bson:"alliancePowerPositions"`                             //公会行为需要最低权限
	AllianceFuzzySearchNum               int32                                      `json:"allianceFuzzySearchNum,omitempty" bson:"allianceFuzzySearchNum"`                             //公会搜索模糊匹配查询数量
	AllianceCreateFirstDayVipLimitWechat int32                                      `json:"allianceCreateFirstDayVipLimitWechat,omitempty" bson:"allianceCreateFirstDayVipLimitWechat"` //开服首日创建公会vip等级限制微信服
	AllianceCreateFirstDayVipLimitOther  int32                                      `json:"allianceCreateFirstDayVipLimitOther,omitempty" bson:"allianceCreateFirstDayVipLimitOther"`   //开服首日创建公会vip等级限制其他服
	IsRemoveMemberNow                    bool                                       `json:"isRemoveMemberNow,omitempty" bson:"isRemoveMemberNow"`                                       //是否立即将玩家踢出公会（不管是否有公会科技秘术券,如果是false则走秘术券数量校验）
}

// AllianceConfig.json属性表集合
type AllianceConfigEntries struct {
	Rows   map[int32]*AllianceConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*AllianceConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("AllianceConfig.json", (*AllianceConfigEntries)(nil))
}

func (e *AllianceConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	allianceConfigEntries = &AllianceConfigEntries{
		Rows:   make(map[int32]*AllianceConfigEntry),
		Arrays: make([]*AllianceConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &AllianceConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return allianceConfigEntries, err
		}

		allianceConfigEntries.Rows[entry.Test] = entry
		allianceConfigEntries.Arrays = append(allianceConfigEntries.Arrays, entry)
	}

	sort.Slice(allianceConfigEntries.Arrays, func(i, j int) bool {
		return allianceConfigEntries.Arrays[i].Test < allianceConfigEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("allianceConfig.json", allianceConfigEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return allianceConfigEntries, nil

}

func (e *AllianceConfigEntries) GetArrays() any {
	return allianceConfigEntries.Arrays
}

func GetAllianceConfigEntry() *AllianceConfigEntry {
	entry := allianceConfigEntries.Rows[0]
	return entry
}

func GetAllianceConfigSize() int {
	return int(len(allianceConfigEntries.Rows))
}

func GetAllianceConfigRows() map[int32]*AllianceConfigEntry {
	return allianceConfigEntries.Rows
}

func GetAllianceConfigArrays() []*AllianceConfigEntry {
	return allianceConfigEntries.Arrays
}
