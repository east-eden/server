package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var configEntries *ConfigEntries //Config.json全局变量

// 属性表 基础配置-拳皇/config
type ConfigEntry struct {
	Test                                  int32             `json:"test,omitempty" bson:"test"`                                                                   //##测试数据
	ConfigType                            string            `json:"configType,omitempty" bson:"configType"`                                                       //配置类型
	NewServerDay                          int32             `json:"newServerDay,omitempty" bson:"newServerDay"`                                                   //服务器开启几天内还算新服，1为当天
	InitGold                              int32             `json:"initGold,omitempty" bson:"initGold"`                                                           //玩家初始金币
	InitItems                             []*IdCount        `json:"initItems,omitempty" bson:"initItems"`                                                         //玩家初始道具[初始角色s]
	HighLevelSpeedUpMaxTimes              int32             `json:"highLevelSpeedUpMaxTimes,omitempty" bson:"highLevelSpeedUpMaxTimes"`                           //每日高级加速次数
	HighLevelSpeedUp                      int32             `json:"highLevelSpeedUp,omitempty" bson:"highLevelSpeedUp"`                                           //单次高级加速时长（秒）
	CameraLerpSpeed                       float32           `json:"cameraLerpSpeed,omitempty" bson:"cameraLerpSpeed"`                                             //相机缓动速度
	BuffAddEffeftBufftype                 []BuffType        `json:"buffAddEffeftBufftype,omitempty" bson:"buffAddEffeftBufftype"`                                 //状态加强属性影响的状态类型
	SpeedUpRate                           []float32         `json:"speedUpRate,omitempty" bson:"speedUpRate"`                                                     //加速实际倍率
	BgmFadeTime                           float32           `json:"bgmFadeTime,omitempty" bson:"bgmFadeTime"`                                                     //音乐切换时间(秒)
	WorldBossSweepTicketMax               int32             `json:"worldBossSweepTicketMax,omitempty" bson:"worldBossSweepTicketMax"`                             //世界boss扫荡券上限
	DailyWorldBossSweepTicketCost         []int32           `json:"dailyWorldBossSweepTicketCost,omitempty" bson:"dailyWorldBossSweepTicketCost"`                 //每日购买世界boss扫荡券消耗钻石
	WorldBossBattleTime                   int32             `json:"worldBossBattleTime,omitempty" bson:"worldBossBattleTime"`                                     //世界boss最大战斗时间(秒)
	FirstMainQuestId                      int32             `json:"firstMainQuestId,omitempty" bson:"firstMainQuestId"`                                           //首个主线任务id
	CameraScales                          []float32         `json:"cameraScales,omitempty" bson:"cameraScales"`                                                   //镜头比例三档对应缩放
	MonthCardDaliyAwardItems              []*IdCount        `json:"monthCardDaliyAwardItems,omitempty" bson:"monthCardDaliyAwardItems"`                           //月卡每日奖励列表
	ForeverCardDaliyAwardItems            []*IdCount        `json:"foreverCardDaliyAwardItems,omitempty" bson:"foreverCardDaliyAwardItems"`                       //终身卡每日奖励列表
	FishCardDailyAwardItems               []*IdCount        `json:"fishCardDailyAwardItems,omitempty" bson:"fishCardDailyAwardItems"`                             //钓鱼卡每日奖励列表
	SoulCardDailyAwardItems               []*IdCount        `json:"soulCardDailyAwardItems,omitempty" bson:"soulCardDailyAwardItems"`                             //战魂卡每日奖励列表
	SoulCardDailyAwardTimes               int32             `json:"soulCardDailyAwardTimes,omitempty" bson:"soulCardDailyAwardTimes"`                             //战魂卡奖励可领取次数
	DiamondCardDailyAwardItems            []*IdCount        `json:"diamondCardDailyAwardItems,omitempty" bson:"diamondCardDailyAwardItems"`                       //钻石卡每日奖励列表
	DiamondCardActiveReward               []*IdCount        `json:"diamondCardActiveReward,omitempty" bson:"diamondCardActiveReward"`                             //钻石卡激活奖励
	DiamondCardChargeAmount               int32             `json:"diamondCardChargeAmount,omitempty" bson:"diamondCardChargeAmount"`                             //钻石卡激活需要的充值数量
	WeekCardDays                          int32             `json:"weekCardDays,omitempty" bson:"weekCardDays"`                                                   //周卡持续时间（天）
	EnjoymentCardDays                     int32             `json:"enjoymentCardDays,omitempty" bson:"enjoymentCardDays"`                                         //畅玩卡持续时间（天）
	MonthCardDays                         int32             `json:"monthCardDays,omitempty" bson:"monthCardDays"`                                                 //月卡持续时间（天）
	FishCardDays                          int32             `json:"fishCardDays,omitempty" bson:"fishCardDays"`                                                   //钓鱼卡持续时间（天）
	SoulCardDays                          int32             `json:"soulCardDays,omitempty" bson:"soulCardDays"`                                                   //战魂卡持续时间（天）
	DiamondCardDays                       int32             `json:"diamondCardDays,omitempty" bson:"diamondCardDays"`                                             //钻石卡持续时间（天）
	WeekHeroRecruitIssues                 int32             `json:"weekHeroRecruitIssues,omitempty" bson:"weekHeroRecruitIssues"`                                 //招募冲刺特权持续期数
	WeekArtifactIssues                    int32             `json:"weekArtifactIssues,omitempty" bson:"weekArtifactIssues"`                                       //神器冲刺特权持续期数
	WeekAlliTechIssues                    int32             `json:"weekAlliTechIssues,omitempty" bson:"weekAlliTechIssues"`                                       //秘术冲刺特权持续期数
	TenDrawCardNum                        int32             `json:"tenDrawCardNum,omitempty" bson:"tenDrawCardNum"`                                               //十连特权卡总出红次数
	TenDrawCardLootId                     int32             `json:"tenDrawCardLootId,omitempty" bson:"tenDrawCardLootId"`                                         //十连特权卡掉落卡池ID
	ArenaTicketMaxNum                     int32             `json:"arenaTicketMaxNum,omitempty" bson:"arenaTicketMaxNum"`                                         //竞技场挑战券持有上限
	AfkNormalExtraLootId                  *IdInterval       `json:"afkNormalExtraLootId,omitempty" bson:"afkNormalExtraLootId"`                                   //离线普通额外掉落id,每间隔多少秒产出一次
	AfkQuickExtraLootId                   *IdInterval       `json:"afkQuickExtraLootId,omitempty" bson:"afkQuickExtraLootId"`                                     //离线快速额外掉落id,每间隔多少秒产出一次
	StagePassExtraLootId                  int32             `json:"stagePassExtraLootId,omitempty" bson:"stagePassExtraLootId"`                                   //主线关卡额外掉落id
	TowerPassExtraLootId                  int32             `json:"towerPassExtraLootId,omitempty" bson:"towerPassExtraLootId"`                                   //爬塔额外掉落id.通关即触发
	RoyalTrialExtraLootId                 int32             `json:"royalTrialExtraLootId,omitempty" bson:"royalTrialExtraLootId"`                                 //皇家试炼额外掉落id
	ArtsTrialExtraLootId                  int32             `json:"artsTrialExtraLootId,omitempty" bson:"artsTrialExtraLootId"`                                   //武道试炼额外掉落id
	OpenBoxLootId                         *IdCount          `json:"openBoxLootId,omitempty" bson:"openBoxLootId"`                                                 //每开箱多少次走一次该loot  (lootId,cnt)
	TimeTicketReduceTime                  int32             `json:"timeTicketReduceTime,omitempty" bson:"timeTicketReduceTime"`                                   //加速券每张缩短的时间(秒)
	FreeADRecruitCD                       int32             `json:"freeADRecruitCD,omitempty" bson:"freeADRecruitCD"`                                             //拳皇看广告招募CD(秒)
	MaxFreeRecruit                        int32             `json:"maxFreeRecruit,omitempty" bson:"maxFreeRecruit"`                                               //拳皇免费普通招募限制次数（每天0点刷新）
	FreerRecruitCD                        int32             `json:"freerRecruitCD,omitempty" bson:"freerRecruitCD"`                                               //拳皇免费普通招募CD时间（秒）
	MaxFreeADRecruit                      int32             `json:"maxFreeADRecruit,omitempty" bson:"maxFreeADRecruit"`                                           //拳皇普通（广告）招募限制次数
	RecruitADLootId                       int32             `json:"recruitADLootId,omitempty" bson:"recruitADLootId"`                                             //拳皇普通（广告）抽卡掉落卡池
	FirstADRecruitLoot                    int32             `json:"firstADRecruitLoot,omitempty" bson:"firstADRecruitLoot"`                                       //首次普通（广告）招募掉落id
	SeniorRecruitCost                     *IdCount          `json:"seniorRecruitCost,omitempty" bson:"seniorRecruitCost"`                                         //拳皇高级招募消耗
	FirstSeniorRecruitLoot                int32             `json:"firstSeniorRecruitLoot,omitempty" bson:"firstSeniorRecruitLoot"`                               //拳皇高级招募掉落id.首次
	FirstTenRecruitLoot                   int32             `json:"firstTenRecruitLoot,omitempty" bson:"firstTenRecruitLoot"`                                     //高级招募十连抽掉落id.首次
	SeniorRecruitLootIdList               []int32           `json:"seniorRecruitLootIdList,omitempty" bson:"seniorRecruitLootIdList"`                             //拳皇高级招募掉落卡池列表
	RecruitShowCardList                   [][]int32         `json:"recruitShowCardList,omitempty" bson:"recruitShowCardList"`                                     //招募卡牌动画icon配置
	RecruitPanttingRandomList             []int32           `json:"recruitPanttingRandomList,omitempty" bson:"recruitPanttingRandomList"`                         //每次招募完随机展示的英雄列表
	SeniorRecruitPeriod                   []int32           `json:"seniorRecruitPeriod,omitempty" bson:"seniorRecruitPeriod"`                                     //高级招募次数区间配置
	SeniorRecruitGuaranteeTimes           int32             `json:"seniorRecruitGuaranteeTimes,omitempty" bson:"seniorRecruitGuaranteeTimes"`                     //拳皇高级招募几次触发保底
	SeniorRecruitGuaranteeCycleList       []int32           `json:"seniorRecruitGuaranteeCycleList,omitempty" bson:"seniorRecruitGuaranteeCycleList"`             //拳皇高级招募保底触发循环次数数组
	SeniorRecruitSmallGuaranteeList       []int32           `json:"seniorRecruitSmallGuaranteeList,omitempty" bson:"seniorRecruitSmallGuaranteeList"`             //拳皇高级招募小保底掉落id列表；小保底必得_英雄
	SeniorRecruitSmallGuaranteeTimes      int32             `json:"seniorRecruitSmallGuaranteeTimes,omitempty" bson:"seniorRecruitSmallGuaranteeTimes"`           //拳皇高级招募小保底几次触发
	SeniorRecruitGuaranteeLootList        []int32           `json:"seniorRecruitGuaranteeLootList,omitempty" bson:"seniorRecruitGuaranteeLootList"`               //拳皇高级招募保底掉落id列表；循环保底必得_英雄.红
	SeniorRecruitLootIdList2              []int32           `json:"seniorRecruitLootIdList2,omitempty" bson:"seniorRecruitLootIdList2"`                           //拳皇高级招募掉落卡池列表
	SeniorRecruitPeriod2                  []int32           `json:"seniorRecruitPeriod2,omitempty" bson:"seniorRecruitPeriod2"`                                   //高级招募次数区间配置
	SeniorRecruitGuaranteeCycleList2      []int32           `json:"seniorRecruitGuaranteeCycleList2,omitempty" bson:"seniorRecruitGuaranteeCycleList2"`           //拳皇高级招募保底触发循环次数数组
	SeniorRecruitSmallGuaranteeList2      []int32           `json:"seniorRecruitSmallGuaranteeList2,omitempty" bson:"seniorRecruitSmallGuaranteeList2"`           //拳皇高级招募小保底掉落id列表；小保底必得_英雄
	SeniorRecruitGuaranteeLootList2       []int32           `json:"seniorRecruitGuaranteeLootList2,omitempty" bson:"seniorRecruitGuaranteeLootList2"`             //拳皇高级招募保底掉落id列表；循环保底必得_英雄.红
	TenTimesDiscount                      float32           `json:"tenTimesDiscount,omitempty" bson:"tenTimesDiscount"`                                           //高级招募十连抽钻石折扣
	TimeLimitRecruitCost                  *IdCount          `json:"timeLimitRecruitCost,omitempty" bson:"timeLimitRecruitCost"`                                   //拳皇限时招募消耗
	TimeLimitRecruitLootId                int32             `json:"timeLimitRecruitLootId,omitempty" bson:"timeLimitRecruitLootId"`                               //拳皇限时招募掉落id
	TimeLimitRecruitOnceScore             int32             `json:"timeLimitRecruitOnceScore,omitempty" bson:"timeLimitRecruitOnceScore"`                         //限时招募一次获得积分数
	SuperRecruitGuaranteeScore            int32             `json:"superRecruitGuaranteeScore,omitempty" bson:"superRecruitGuaranteeScore"`                       //拳皇招募大保底积分（达到积分后在库里随机一个英雄）适用于普通和高级招募
	RecruitOnceScore                      int32             `json:"recruitOnceScore,omitempty" bson:"recruitOnceScore"`                                           //拳皇（普通，高级）招募一次获得积分数
	SuperRecruitGuaranteeLootId           int32             `json:"superRecruitGuaranteeLootId,omitempty" bson:"superRecruitGuaranteeLootId"`                     //拳皇大保底掉落卡池
	FundRecruitGuaranteeScore             int32             `json:"fundRecruitGuaranteeScore,omitempty" bson:"fundRecruitGuaranteeScore"`                         //拳皇招募基金大保底积分（达到积分后在库里随机一个英雄）适用于普通和高级招募
	FundRecruitOnceScore                  int32             `json:"fundRecruitOnceScore,omitempty" bson:"fundRecruitOnceScore"`                                   //拳皇（普通，高级）招募一次获得基金大保底积分数
	FundRecruitGuaranteeLootId            int32             `json:"fundRecruitGuaranteeLootId,omitempty" bson:"fundRecruitGuaranteeLootId"`                       //拳皇基金大保底掉落卡池
	MaxDreamHeroCount                     int32             `json:"maxDreamHeroCount,omitempty" bson:"maxDreamHeroCount"`                                         //最大心愿英雄设定数量
	DreamHeroWeightSeq                    []int32           `json:"dreamHeroWeightSeq,omitempty" bson:"dreamHeroWeightSeq"`                                       //心愿英雄权重序列
	LifeCardStartDream                    int32             `json:"lifeCardStartDream,omitempty" bson:"lifeCardStartDream"`                                       //终身卡起始编号；即心愿权重序列从1开始
	DreamHeroTriggerLootId                int32             `json:"dreamHeroTriggerLootId,omitempty" bson:"dreamHeroTriggerLootId"`                               //触发许愿逻辑对应的掉落ID
	SuperRecruitUnlockVipLevel            int32             `json:"superRecruitUnlockVipLevel,omitempty" bson:"superRecruitUnlockVipLevel"`                       //解锁大保底招募vip等级
	DreamHeroOpenNum                      int32             `json:"dreamHeroOpenNum,omitempty" bson:"dreamHeroOpenNum"`                                           //心愿系统需要多少次高级招募开启
	CameraFollowDis                       float32           `json:"cameraFollowDis,omitempty" bson:"cameraFollowDis"`                                             //相机重新跟随距离
	SeventhEveningBoxGuaranteeLootId      int32             `json:"seventhEveningBoxGuaranteeLootId,omitempty" bson:"seventhEveningBoxGuaranteeLootId"`           //七夕盲盒保底掉落id
	DailyRecruitMaxTimes                  int32             `json:"dailyRecruitMaxTimes,omitempty" bson:"dailyRecruitMaxTimes"`                                   //每日招募最大次数
	AthenaBattleCountdown                 int32             `json:"athenaBattleCountdown,omitempty" bson:"athenaBattleCountdown"`                                 //雅呆娜战斗限时(秒)
	DouyinDailyCreateRoomTimes            int32             `json:"douyinDailyCreateRoomTimes,omitempty" bson:"douyinDailyCreateRoomTimes"`                       //抖音每日创建房间免费次数
	DouyinDailyJoinRoomTimes              int32             `json:"douyinDailyJoinRoomTimes,omitempty" bson:"douyinDailyJoinRoomTimes"`                           //抖音每日加入房间免费次数
	DouyinDailyAdGotCreateRoomTimes       int32             `json:"douyinDailyAdGotCreateRoomTimes,omitempty" bson:"douyinDailyAdGotCreateRoomTimes"`             //抖音每日广告获取创建房间次数
	DouyinDailyAdGotJoinRoomTimes         int32             `json:"douyinDailyAdGotJoinRoomTimes,omitempty" bson:"douyinDailyAdGotJoinRoomTimes"`                 //抖音每日广告获取加入房间次数
	MaidCardTakeDays                      int32             `json:"maidCardTakeDays,omitempty" bson:"maidCardTakeDays"`                                           //保姆卡生效天数
	MaidCardDailyLootId                   int32             `json:"maidCardDailyLootId,omitempty" bson:"maidCardDailyLootId"`                                     //保姆卡每日领奖掉落id
	ContributeInviteSuccessItems          []*IdCount        `json:"contributeInviteSuccessItems,omitempty" bson:"contributeInviteSuccessItems"`                   //贡献福利成功邀请奖励道具
	NewInviteValidScene                   []int32           `json:"newInviteValidScene,omitempty" bson:"newInviteValidScene"`                                     //新邀请玩家的有效场景值
	NewbieSkinDuration                    int32             `json:"newbieSkinDuration,omitempty" bson:"newbieSkinDuration"`                                       //新手皮肤体验时间(小时)
	SkinExperienceDuration                int32             `json:"skinExperienceDuration,omitempty" bson:"skinExperienceDuration"`                               //皮肤体验道具生效时间(小时)
	NewbieSkinList                        []int32           `json:"newbieSkinList,omitempty" bson:"newbieSkinList"`                                               //新手可选皮肤列表
	IosPlatformChangePunishTime           int32             `json:"iosPlatformChangePunishTime,omitempty" bson:"iosPlatformChangePunishTime"`                     //非首次登陆苹果平台充值惩罚持续时间(秒)
	ContributionRecordMax                 int32             `json:"contributionRecordMax,omitempty" bson:"contributionRecordMax"`                                 //赠送和收取记录的最大数量，独立使用
	ContributionSendRecycleLimit          int32             `json:"contributionSendRecycleLimit,omitempty" bson:"contributionSendRecycleLimit"`                   //贡献未回收的数量超过限制不让赠送
	FreeRecruitMaxCnt                     int32             `json:"freeRecruitMaxCnt,omitempty" bson:"freeRecruitMaxCnt"`                                         //触发免费招募的最大次数
	RecruitOneCost                        []*IdCount        `json:"recruitOneCost,omitempty" bson:"recruitOneCost"`                                               //单次招募消耗招募券及数量
	RecruitTenCost                        []*IdCount        `json:"recruitTenCost,omitempty" bson:"recruitTenCost"`                                               //十倍招募消耗招募券及数量
	EpicHeroGotStoneCnt                   int32             `json:"epicHeroGotStoneCnt,omitempty" bson:"epicHeroGotStoneCnt"`                                     //神话英雄分解得到神话石数量
	DailyFreeBlessTimes                   int32             `json:"dailyFreeBlessTimes,omitempty" bson:"dailyFreeBlessTimes"`                                     //每日免费祝福次数
	WatchAdAddFreeBlessTimes              int32             `json:"watchAdAddFreeBlessTimes,omitempty" bson:"watchAdAddFreeBlessTimes"`                           //每次广告获取免费祝福次数
	PvpBattleCountdown                    int32             `json:"pvpBattleCountdown,omitempty" bson:"pvpBattleCountdown"`                                       //pvp战斗限时(秒)
	DouyinAwardLootId                     int32             `json:"douyinAwardLootId,omitempty" bson:"douyinAwardLootId"`                                         //抖音用户领取奖励掉落id
	DouYinAddDesktopAwards                []*IdCount        `json:"douYinAddDesktopAwards,omitempty" bson:"douYinAddDesktopAwards"`                               //抖音添加桌面奖励
	SendContributionMaxLimit              int32             `json:"sendContributionMaxLimit,omitempty" bson:"sendContributionMaxLimit"`                           //赠送贡献最大值
	SendContributionOverTime              int32             `json:"sendContributionOverTime,omitempty" bson:"sendContributionOverTime"`                           //赠送贡献链接超时时间(秒)
	ContributionOutMax                    float32           `json:"contributionOutMax,omitempty" bson:"contributionOutMax"`                                       //贡献每日产出最大值
	ContributionOutInitCoe                float32           `json:"contributionOutInitCoe,omitempty" bson:"contributionOutInitCoe"`                               //初始产出的贡献期望
	ContributionOutOneMaxCoe              float32           `json:"contributionOutOneMaxCoe,omitempty" bson:"contributionOutOneMaxCoe"`                           //单次贡献上限系数
	FriendMatchFreeTimes                  int32             `json:"friendMatchFreeTimes,omitempty" bson:"friendMatchFreeTimes"`                                   //好友组队每日免费次数
	SingleMatchFreeTimes                  int32             `json:"singleMatchFreeTimes,omitempty" bson:"singleMatchFreeTimes"`                                   //单人匹配每日免费次数
	InviteeMatchFreeTimes                 int32             `json:"inviteeMatchFreeTimes,omitempty" bson:"inviteeMatchFreeTimes"`                                 //受邀匹配每日免费次数
	MutationFriendMatchFreeTimes          int32             `json:"mutationFriendMatchFreeTimes,omitempty" bson:"mutationFriendMatchFreeTimes"`                   //异变好友组队每日免费次数
	MutationSingleMatchFreeTimes          int32             `json:"mutationSingleMatchFreeTimes,omitempty" bson:"mutationSingleMatchFreeTimes"`                   //异变单人匹配每日免费次数
	FriendMatchRegenTimes                 int32             `json:"friendMatchRegenTimes,omitempty" bson:"friendMatchRegenTimes"`                                 //好友组队看广告后恢复次数值
	SingleMatchRegenTimes                 int32             `json:"singleMatchRegenTimes,omitempty" bson:"singleMatchRegenTimes"`                                 //单人匹配看广告后恢复次数值
	InviteeMatchRegenTimes                int32             `json:"inviteeMatchRegenTimes,omitempty" bson:"inviteeMatchRegenTimes"`                               //受邀者匹配看广告恢复次数值
	MutationFriendMatchRegenTimes         int32             `json:"mutationFriendMatchRegenTimes,omitempty" bson:"mutationFriendMatchRegenTimes"`                 //异变好友组队看广告后恢复次数值
	MutationSingleMatchRegenTimes         int32             `json:"mutationSingleMatchRegenTimes,omitempty" bson:"mutationSingleMatchRegenTimes"`                 //异变单人匹配看广告后恢复次数值
	FriendMatchEnterTimes                 int32             `json:"friendMatchEnterTimes,omitempty" bson:"friendMatchEnterTimes"`                                 //好友组队进入消耗次数
	SingleMatchEnterTimes                 int32             `json:"singleMatchEnterTimes,omitempty" bson:"singleMatchEnterTimes"`                                 //单人匹配进入消耗次数
	InviteeMatchEnterTimes                int32             `json:"inviteeMatchEnterTimes,omitempty" bson:"inviteeMatchEnterTimes"`                               //受邀者匹配进入消耗次数
	FriendMatchMaxTimes                   int32             `json:"friendMatchMaxTimes,omitempty" bson:"friendMatchMaxTimes"`                                     //好友组队进入最大次数
	SingleMatchMaxTimes                   int32             `json:"singleMatchMaxTimes,omitempty" bson:"singleMatchMaxTimes"`                                     //单人匹配进入最大次数
	InviteeMatchMaxTimes                  int32             `json:"inviteeMatchMaxTimes,omitempty" bson:"inviteeMatchMaxTimes"`                                   //受邀者匹配进入最大次数
	MutationFriendMatchEnterTimes         int32             `json:"mutationFriendMatchEnterTimes,omitempty" bson:"mutationFriendMatchEnterTimes"`                 //异变好友组队进入消耗次数
	MutationSingleMatchEnterTimes         int32             `json:"mutationSingleMatchEnterTimes,omitempty" bson:"mutationSingleMatchEnterTimes"`                 //异变单人匹配进入消耗次数
	MutationEventFriendMatchFreeTimes     int32             `json:"mutationEventFriendMatchFreeTimes,omitempty" bson:"mutationEventFriendMatchFreeTimes"`         //任务异变好友组队每日免费次数
	MutationEventFriendMatchRegenTimes    int32             `json:"mutationEventFriendMatchRegenTimes,omitempty" bson:"mutationEventFriendMatchRegenTimes"`       //任务异变好友组队看广告后恢复次数值
	MutationEventFriendMatchEnterTimes    int32             `json:"mutationEventFriendMatchEnterTimes,omitempty" bson:"mutationEventFriendMatchEnterTimes"`       //任务异变好友组队进入消耗次数
	SkinUnrepeatedTimes                   int32             `json:"skinUnrepeatedTimes,omitempty" bson:"skinUnrepeatedTimes"`                                     //前x次皮肤掉落去重
	CollectionUnrepeatedTimes             int32             `json:"collectionUnrepeatedTimes,omitempty" bson:"collectionUnrepeatedTimes"`                         //前x次藏品掉落去重
	BattleSkinMapUnrepeatedTimes          int32             `json:"battleSkinMapUnrepeatedTimes,omitempty" bson:"battleSkinMapUnrepeatedTimes"`                   //前x次战斗地图皮肤掉落去重
	BattleSkinRoleUnrepeatedTimes         int32             `json:"battleSkinRoleUnrepeatedTimes,omitempty" bson:"battleSkinRoleUnrepeatedTimes"`                 //前x次战斗角色皮肤掉落去重
	HighestWaveFinishGuide                int32             `json:"highestWaveFinishGuide,omitempty" bson:"highestWaveFinishGuide"`                               //完成新手引导后修改最高波数
	ChallengeWaveLoots                    map[int32]int32   `json:"challengeWaveLoots,omitempty" bson:"challengeWaveLoots"`                                       //挑战波次对应奖励掉落id
	ChallengeWaveHighest                  int32             `json:"challengeWaveHighest,omitempty" bson:"challengeWaveHighest"`                                   //挑战最高波次
	DailyFormationLegendNumProb           map[int32]float32 `json:"dailyFormationLegendNumProb,omitempty" bson:"dailyFormationLegendNumProb"`                     //每日阵容传说品质的数量对应的随机概率
	DailyFormationLegendCertainHeroId     int32             `json:"dailyFormationLegendCertainHeroId,omitempty" bson:"dailyFormationLegendCertainHeroId"`         //每日阵容随机传说品质时，必定有这个英雄
	HeroRarityLevelMaxWeights             [][]int32         `json:"heroRarityLevelMaxWeights,omitempty" bson:"heroRarityLevelMaxWeights"`                         //英雄品质对应各等级上限权重
	InviteMaxNum                          int32             `json:"inviteMaxNum,omitempty" bson:"inviteMaxNum"`                                                   //可以邀请好友上限
	InviteSuccLootIds                     []int32           `json:"inviteSuccLootIds,omitempty" bson:"inviteSuccLootIds"`                                         //邀请好友第奖励的掉落，格式为：首个好友邀请成功掉落ID,第二个掉落ID,第三个掉落ID。。。以此类推，没填够数量就按最后一个数据掉落
	DamageSpatterRate                     float32           `json:"damageSpatterRate,omitempty" bson:"damageSpatterRate"`                                         //伤害溅射衰减
	SpSpatterRate                         float32           `json:"spSpatterRate,omitempty" bson:"spSpatterRate"`                                                 //溅射回能衰减
	LDCoinOutMax                          float32           `json:"lDCoinOutMax,omitempty" bson:"lDCoinOutMax"`                                                   //LD副玩法每日代币产出最大值
	LDCoinOutInitCoe                      float32           `json:"lDCoinOutInitCoe,omitempty" bson:"lDCoinOutInitCoe"`                                           //LD副玩法初始产出的代币期望
	LDCoinOutOneMaxCoe                    float32           `json:"lDCoinOutOneMaxCoe,omitempty" bson:"lDCoinOutOneMaxCoe"`                                       //LD副玩法单次代币上限系数
	LDCoinHighYieldDay                    *LDHighCoin       `json:"lDCoinHighYieldDay,omitempty" bson:"lDCoinHighYieldDay"`                                       //星期五,限制2次,单次给200(暂时未使用)
	LDFormationRefreshCost                int32             `json:"lDFormationRefreshCost,omitempty" bson:"lDFormationRefreshCost"`                               //LD阵容钻石刷新消耗
	LDFormationRefreshMax                 int32             `json:"lDFormationRefreshMax,omitempty" bson:"lDFormationRefreshMax"`                                 //广告刷新阵容每天广告次数
	LDSpecialCardMax                      int32             `json:"lDSpecialCardMax,omitempty" bson:"lDSpecialCardMax"`                                           //广告获取道具卡每天广告次数
	AllianceBossServerBattlePerLevel      int32             `json:"allianceBossServerBattlePerLevel,omitempty" bson:"allianceBossServerBattlePerLevel"`           //公会boss每x级使用后端战斗验证
	StageSnapshotLevel                    int32             `json:"stageSnapshotLevel,omitempty" bson:"stageSnapshotLevel"`                                       //主线关卡存储阵容快照起始ID
	StageSnapshotPowerRange               []float32         `json:"stageSnapshotPowerRange,omitempty" bson:"stageSnapshotPowerRange"`                             //主线关卡推荐战力区间系数
	StageSnapshotCount                    int32             `json:"stageSnapshotCount,omitempty" bson:"stageSnapshotCount"`                                       //阵容存储数量
	StageSnapshotSimilarity               float32           `json:"stageSnapshotSimilarity,omitempty" bson:"stageSnapshotSimilarity"`                             //阵容筛选系数
	WechatServerGroupIdUseWrongProperty   int32             `json:"wechatServerGroupIdUseWrongProperty,omitempty" bson:"wechatServerGroupIdUseWrongProperty"`     //微信服X战区以下的所有服务器使用老装备突破计算方法
	WechatServerGroupIdArtsTrialFreeTimes int32             `json:"wechatServerGroupIdArtsTrialFreeTimes,omitempty" bson:"wechatServerGroupIdArtsTrialFreeTimes"` //微信服X战区一下的所有服务器vip+武道试炼免费次数
	WechatServerGroupIdUseOldSetting      int32             `json:"wechatServerGroupIdUseOldSetting,omitempty" bson:"wechatServerGroupIdUseOldSetting"`           //微信服X战区以下的所有服务器使用老的招募配置
	WechatRecruitLevelServerId            []int32           `json:"wechatRecruitLevelServerId,omitempty" bson:"wechatRecruitLevelServerId"`                       //微信服x服务器id开启招募升级
	TikTokServerIdUseOldSetting           []int32           `json:"tikTokServerIdUseOldSetting,omitempty" bson:"tikTokServerIdUseOldSetting"`                     //抖音服X服务器id使用老的招募配置（A）配置
	WechatHeroStarFinalMultiple           []int32           `json:"wechatHeroStarFinalMultiple,omitempty" bson:"wechatHeroStarFinalMultiple"`                     //微信服x服务器id开启英雄升星最后加成百分比算法
	WechatHeroLevelupOldPlanServerGroup   int32             `json:"wechatHeroLevelupOldPlanServerGroup,omitempty" bson:"wechatHeroLevelupOldPlanServerGroup"`     //微信服x战区以下(不含x)的所有服务器使用老的英雄升级方案
	WechatServerLoginForbidden            []int32           `json:"wechatServerLoginForbidden,omitempty" bson:"wechatServerLoginForbidden"`                       //微信服x服务器id不允许新玩家登录
	DouyinServerLoginForbidden            []int32           `json:"douyinServerLoginForbidden,omitempty" bson:"douyinServerLoginForbidden"`                       //抖音服x服务器id不允许新玩家登录
	UserNameMaxLimit                      int32             `json:"userNameMaxLimit,omitempty" bson:"userNameMaxLimit"`                                           //用户名最大程度（改名时候使用）
	ChangeNameCost                        *IdCount          `json:"changeNameCost,omitempty" bson:"changeNameCost"`                                               //改名消耗
	RecruitChipExp                        map[Rarity]int32  `json:"recruitChipExp,omitempty" bson:"recruitChipExp"`                                               //招募获得碎片品质对应招募经验
	RecruitHeroExp                        map[Rarity]int32  `json:"recruitHeroExp,omitempty" bson:"recruitHeroExp"`                                               //招募获得整卡品质对应招募经验
	UpRecruitReduceTime                   int32             `json:"upRecruitReduceTime,omitempty" bson:"upRecruitReduceTime"`                                     //加速券减少时间(秒)
	RecruitUpgradeAdVideoReduceTime       int32             `json:"recruitUpgradeAdVideoReduceTime,omitempty" bson:"recruitUpgradeAdVideoReduceTime"`             //广告加速减少时间(秒)
	AdCdRecruitLevelup                    int32             `json:"adCdRecruitLevelup,omitempty" bson:"adCdRecruitLevelup"`                                       //升级招募广告冷却时间（秒）
	InitStamina                           int32             `json:"initStamina,omitempty" bson:"initStamina"`                                                     //初始体力
	InitMaxStamina                        int32             `json:"initMaxStamina,omitempty" bson:"initMaxStamina"`                                               //初始体力上限
	StaminaRegen                          []int32           `json:"staminaRegen,omitempty" bson:"staminaRegen"`                                                   //每x秒恢复y点体力
	StaminaPunish                         []*AttrValue      `json:"staminaPunish,omitempty" bson:"staminaPunish"`                                                 //体力=0时，怪物攻击+X%，防御免伤比+X%
	IntervalVideoRewards                  []*IdCount        `json:"intervalVideoRewards,omitempty" bson:"intervalVideoRewards"`                                   //拳皇QTE胜利奖励
	QTELevel                              [][]int32         `json:"qTELevel,omitempty" bson:"qTELevel"`                                                           //拳皇QTE每日关卡
	SPINCRERESPOSTIVEFactor               float32           `json:"sPINCRERESPOSTIVEFactor,omitempty" bson:"sPINCRERESPOSTIVEFactor"`                             //攻击回怒最终乘区系数
	SPINCRERESNEGIATIVEFactor             float32           `json:"sPINCRERESNEGIATIVEFactor,omitempty" bson:"sPINCRERESNEGIATIVEFactor"`                         //受击回怒最终乘区系数
	AnchorServerGroupId                   int32             `json:"anchorServerGroupId,omitempty" bson:"anchorServerGroupId"`                                     //主播战区id
	AlipayFirstQuestAward                 []*IdCount        `json:"alipayFirstQuestAward,omitempty" bson:"alipayFirstQuestAward"`                                 //支付宝首任务奖励
	AlipayRepeatQuestAward                []*IdCount        `json:"alipayRepeatQuestAward,omitempty" bson:"alipayRepeatQuestAward"`                               //支付宝复访任务奖励
	AlipayRepeatQuestAwardDays            int32             `json:"alipayRepeatQuestAwardDays,omitempty" bson:"alipayRepeatQuestAwardDays"`                       //支付宝复访任务奖励天数
	PcDailyAward                          []*IdCount        `json:"pcDailyAward,omitempty" bson:"pcDailyAward"`                                                   //pc端每日奖励
	MeituanFirstQuestAward                []*IdCount        `json:"meituanFirstQuestAward,omitempty" bson:"meituanFirstQuestAward"`                               //美团首任务奖励
	MeituanRepeatQuestAward               []*IdCount        `json:"meituanRepeatQuestAward,omitempty" bson:"meituanRepeatQuestAward"`                             //美团复访任务奖励
	MeituanRepeatQuestAwardDays           int32             `json:"meituanRepeatQuestAwardDays,omitempty" bson:"meituanRepeatQuestAwardDays"`                     //美团复访任务奖励天数
	MaxPlayersPerAccount                  int32             `json:"maxPlayersPerAccount,omitempty" bson:"maxPlayersPerAccount"`                                   //每个账号最多创建多少个角色
}

// Config.json属性表集合
type ConfigEntries struct {
	Rows   map[int32]*ConfigEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ConfigEntry         `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Config.json", (*ConfigEntries)(nil))
}

func (e *ConfigEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	configEntries = &ConfigEntries{
		Rows:   make(map[int32]*ConfigEntry),
		Arrays: make([]*ConfigEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ConfigEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return configEntries, err
		}

		configEntries.Rows[entry.Test] = entry
		configEntries.Arrays = append(configEntries.Arrays, entry)
	}

	sort.Slice(configEntries.Arrays, func(i, j int) bool {
		return configEntries.Arrays[i].Test < configEntries.Arrays[j].Test
	})

	excel.AddArrayGetter("config.json", configEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return configEntries, nil

}

func (e *ConfigEntries) GetArrays() any {
	return configEntries.Arrays
}

func GetConfigEntry() *ConfigEntry {
	entry := configEntries.Rows[0]
	return entry
}

func GetConfigSize() int {
	return int(len(configEntries.Rows))
}

func GetConfigRows() map[int32]*ConfigEntry {
	return configEntries.Rows
}

func GetConfigArrays() []*ConfigEntry {
	return configEntries.Arrays
}
