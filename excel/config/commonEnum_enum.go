package config

// 基础配置-拳皇/commonEnum
type EMPlatform int

const (
	EMPlatform_BEGIN   EMPlatform = 0 //枚举开始
	EMPlatform_All     EMPlatform = 0 //全平台
	EMPlatform_Wechat  EMPlatform = 1 //微信
	EMPlatform_DouYin  EMPlatform = 2 //抖音
	EMPlatform_AliPay  EMPlatform = 3 //支付宝
	EMPlatform_MeiTuan EMPlatform = 4 //美团
	EMPlatform_END     EMPlatform = 5 //枚举结束
)

// 基础配置-拳皇/commonEnum
type GoodsType int

const (
	GoodsType_BEGIN          GoodsType = 1  //枚举开始
	GoodsType_Value          GoodsType = 1  //数值
	GoodsType_Hero           GoodsType = 2  //英雄
	GoodsType_Shop           GoodsType = 4  //商店
	GoodsType_Tech           GoodsType = 5  //科技
	GoodsType_BattleSkin     GoodsType = 6  //关卡皮肤
	GoodsType_Level          GoodsType = 7  //关卡道具
	GoodsType_FashionClothes GoodsType = 8  //装扮
	GoodsType_Pet            GoodsType = 9  //宠物
	GoodsType_END            GoodsType = 10 //枚举结束
)

// 基础配置-拳皇/commonEnum
type GoodsSubType int

const (
	GoodsSubType_BEGIN                        GoodsSubType = 0   //枚举开始
	GoodsSubType_Empty                        GoodsSubType = 0   //空掉落
	GoodsSubType_Gold                         GoodsSubType = 1   //金币
	GoodsSubType_Diamond                      GoodsSubType = 2   //钻石
	GoodsSubType_Exp                          GoodsSubType = 3   //玩家经验
	GoodsSubType_HeroExp                      GoodsSubType = 4   //英雄经验
	GoodsSubType_HeroBadgePoint               GoodsSubType = 5   //英雄图鉴点数
	GoodsSubType_TowerTicket                  GoodsSubType = 6   //爬塔挑战券
	GoodsSubType_ManorResTicket               GoodsSubType = 7   //领地资源刷新券
	GoodsSubType_MonthCardPoint               GoodsSubType = 8   //月卡积分
	GoodsSubType_ArenaTicket                  GoodsSubType = 9   //竞技场挑战券
	GoodsSubType_EquipStone                   GoodsSubType = 10  //装备强化石
	GoodsSubType_EnjoymentCard                GoodsSubType = 11  //畅玩卡
	GoodsSubType_VipRight                     GoodsSubType = 12  //Vip特权
	GoodsSubType_CashCharge                   GoodsSubType = 13  //现金充值
	GoodsSubType_GameCartridge                GoodsSubType = 14  //游戏卡带
	GoodsSubType_WorldBossSweepTicket         GoodsSubType = 15  //世界boss扫荡券
	GoodsSubType_MonthCard                    GoodsSubType = 16  //月卡
	GoodsSubType_LifetimeCard                 GoodsSubType = 17  //终身卡
	GoodsSubType_HighLevelBattlePass          GoodsSubType = 18  //高级等级通行证
	GoodsSubType_AllianceDailyTaskProgress    GoodsSubType = 19  //联盟每日任务进度
	GoodsSubType_AllianceExp                  GoodsSubType = 20  //联盟经验
	GoodsSubType_VipExp                       GoodsSubType = 21  //Vip经验
	GoodsSubType_PromotionExp                 GoodsSubType = 22  //佣兵团晋升经验
	GoodsSubType_HighGrowFund1BattlePass      GoodsSubType = 23  //高级成长基金1通行证
	GoodsSubType_HighGrowFund2BattlePass      GoodsSubType = 24  //高级成长基金2通行证
	GoodsSubType_SeniorRecruitTicket          GoodsSubType = 25  //高级招募券
	GoodsSubType_HighChallengeFund1BattlePass GoodsSubType = 26  //高级挑战基金1通行证
	GoodsSubType_HighChallengeFund2BattlePass GoodsSubType = 27  //高级挑战基金2通行证
	GoodsSubType_HighHonorFund1BattlePass     GoodsSubType = 28  //高级荣誉基金1通行证
	GoodsSubType_HighHonorFund2BattlePass     GoodsSubType = 29  //高级荣誉基金2通行证
	GoodsSubType_DailyActiveScore             GoodsSubType = 30  //日常活跃积分
	GoodsSubType_WeeklyActiveScore            GoodsSubType = 31  //周常活跃积分
	GoodsSubType_Astrology                    GoodsSubType = 32  //神器占星
	GoodsSubType_BoxGift                      GoodsSubType = 33  //格斗家收藏
	GoodsSubType_AllianceTechFragment         GoodsSubType = 34  //公会科技碎片
	GoodsSubType_AllianceTechTicket           GoodsSubType = 35  //公会科技能量
	GoodsSubType_HeroWeekCard                 GoodsSubType = 36  //英雄周卡
	GoodsSubType_FishWeekCard                 GoodsSubType = 37  //钓鱼周卡
	GoodsSubType_TreasureWeekCard             GoodsSubType = 38  //宝物周卡
	GoodsSubType_ActivityPoint                GoodsSubType = 39  //活动积分
	GoodsSubType_AfkExpPackage                GoodsSubType = 40  //离线英雄经验包
	GoodsSubType_AfkGoldPackage               GoodsSubType = 41  //离线金币包
	GoodsSubType_WeekHeroRecruitPoint         GoodsSubType = 42  //招募周积分
	GoodsSubType_WeekArtifactDrawPoint        GoodsSubType = 43  //神器周积分
	GoodsSubType_SevenDayTaskScore            GoodsSubType = 44  //七日任务积分
	GoodsSubType_Box                          GoodsSubType = 45  //格斗家收藏盒
	GoodsSubType_FishEnergy                   GoodsSubType = 46  //钓鱼体力
	GoodsSubType_FishSilver                   GoodsSubType = 47  //钓鱼用银币
	GoodsSubType_FishPearl                    GoodsSubType = 48  //钓鱼用珍珠
	GoodsSubType_FishBlackPearl               GoodsSubType = 49  //钓鱼用黑珍珠
	GoodsSubType_EquipStoneBox                GoodsSubType = 50  //装备强化石自选箱
	GoodsSubType_FishOrderRefreshTicket       GoodsSubType = 51  //钓鱼订单刷新券
	GoodsSubType_FishOrderResetTicket         GoodsSubType = 52  //钓鱼订单重置券
	GoodsSubType_ReduceTimeTicket             GoodsSubType = 53  //加速券
	GoodsSubType_FishBait                     GoodsSubType = 54  //鱼饵
	GoodsSubType_FishGalaPoint                GoodsSubType = 55  //鱼获宴会积分
	GoodsSubType_BountyChallengeScore         GoodsSubType = 56  //赏金值
	GoodsSubType_SoulPotion                   GoodsSubType = 57  //战魂合成药剂
	GoodsSubType_AllianceCheerExp             GoodsSubType = 58  //跨服竞技助威经验
	GoodsSubType_CrossArenaTicket             GoodsSubType = 59  //跨服竞技门票
	GoodsSubType_CrossArenaMedal              GoodsSubType = 60  //跨服竞技勋章
	GoodsSubType_Contribute                   GoodsSubType = 61  //贡献
	GoodsSubType_ContributeChallengeTicket    GoodsSubType = 62  //贡献挑战券
	GoodsSubType_AthenaSweepTicket            GoodsSubType = 63  //亚呆娜扫荡券
	GoodsSubType_PlaygroundTicket             GoodsSubType = 64  //乐园券
	GoodsSubType_RoomCard                     GoodsSubType = 65  //房卡
	GoodsSubType_HeroSkin                     GoodsSubType = 66  //英雄皮肤
	GoodsSubType_Collection                   GoodsSubType = 67  //藏品
	GoodsSubType_Coin                         GoodsSubType = 68  //硬币
	GoodsSubType_LuckyStone                   GoodsSubType = 69  //幸运石
	GoodsSubType_JZYFish                      GoodsSubType = 70  //姜子牙钓鱼
	GoodsSubType_ChooseSkinExperience         GoodsSubType = 71  //自选皮肤体验
	GoodsSubType_Avatar                       GoodsSubType = 72  //头像
	GoodsSubType_AvatarFrame                  GoodsSubType = 73  //头像框
	GoodsSubType_Title                        GoodsSubType = 74  //称号
	GoodsSubType_FormationRefresh             GoodsSubType = 75  //阵容刷新券
	GoodsSubType_ChatBuble                    GoodsSubType = 76  //聊天气泡
	GoodsSubType_NewbieHeroRecruitPoint       GoodsSubType = 77  //新手招募积分
	GoodsSubType_AwakenStone                  GoodsSubType = 78  //觉醒石
	GoodsSubType_ArtifactBattlePass           GoodsSubType = 79  //神器基金通行证
	GoodsSubType_CollectBattlePass            GoodsSubType = 80  //收藏基金通信证
	GoodsSubType_NestBattlePass               GoodsSubType = 81  //音巢基金通信证
	GoodsSubType_RecruitBattlePass            GoodsSubType = 82  //招募基金通信证
	GoodsSubType_WeekAllianceTechPoint        GoodsSubType = 83  //公会科技周积分
	GoodsSubType_Treasure                     GoodsSubType = 84  //宝箱
	GoodsSubType_LDCoin                       GoodsSubType = 85  //LD代币
	GoodsSubType_LDSingleMatchCnt             GoodsSubType = 86  //LD单人匹配次数
	GoodsSubType_LDFriendMatchCnt             GoodsSubType = 87  //LD好友匹配次数
	GoodsSubType_LDInviteeMatchCnt            GoodsSubType = 88  //LD受邀匹配次数
	GoodsSubType_AllianceStatueExp            GoodsSubType = 89  //公会雕像经验
	GoodsSubType_StarCoin                     GoodsSubType = 90  //星钻
	GoodsSubType_ArtifactPassPoint            GoodsSubType = 91  //神器基金积分
	GoodsSubType_FameMedal                    GoodsSubType = 92  //勋章
	GoodsSubType_ManorStone                   GoodsSubType = 93  //石料
	GoodsSubType_ManorWood                    GoodsSubType = 94  //木料
	GoodsSubType_ManorFood                    GoodsSubType = 95  //食物
	GoodsSubType_AllianceTech                 GoodsSubType = 96  //秘术
	GoodsSubType_BoxGiftExp                   GoodsSubType = 97  //收藏经验
	GoodsSubType_AllianceDuelScore            GoodsSubType = 98  //公会对决积分
	GoodsSubType_BattlePassScore              GoodsSubType = 99  //通用战令积分
	GoodsSubType_WorldBossCoin                GoodsSubType = 100 //世界boss代币
	GoodsSubType_FishCard                     GoodsSubType = 101 //钓鱼权益卡
	GoodsSubType_UpRecruitTicket              GoodsSubType = 102 //Up招募券
	GoodsSubType_WishRecruitTicket            GoodsSubType = 103 //心愿招募券
	GoodsSubType_AutoFishTicket               GoodsSubType = 104 //自动钓鱼券
	GoodsSubType_ArtsTrialTimes               GoodsSubType = 105 //武道试炼次数
	GoodsSubType_BoxGiftBlessExp              GoodsSubType = 106 //祝福经验
	GoodsSubType_NestChallengeTicket          GoodsSubType = 107 //音巢挑战券
	GoodsSubType_NestChallengeKey             GoodsSubType = 108 //音巢挑战密室钥匙
	GoodsSubType_NestSweepTicket              GoodsSubType = 109 //音巢扫荡券
	GoodsSubType_AllianceBossAttackTimes      GoodsSubType = 110 //联盟boss攻击次数
	GoodsSubType_SpSweepTicket                GoodsSubType = 111 //南镇扫荡券
	GoodsSubType_WeekHeroRecruitCard          GoodsSubType = 112 //招募冲刺特权卡
	GoodsSubType_WeekArtifactDrawCard         GoodsSubType = 113 //神器冲刺特权卡
	GoodsSubType_WeekAllianceTechCard         GoodsSubType = 114 //秘术冲刺特权卡
	GoodsSubType_RecruitExp                   GoodsSubType = 115 //招募经验
	GoodsSubType_FightingPlaceTicket          GoodsSubType = 116 //格斗圣殿挑战券
	GoodsSubType_TenDrawCard                  GoodsSubType = 117 //十连特权卡
	GoodsSubType_Stamina                      GoodsSubType = 118 //体力道具
	GoodsSubType_StaminaMax                   GoodsSubType = 119 //体力上限道具
	GoodsSubType_FameForgeTicket              GoodsSubType = 120 //名人堂锻造券
	GoodsSubType_TimesTicket                  GoodsSubType = 121 //限次门票
	GoodsSubType_FameGem                      GoodsSubType = 122 //名人堂宝石
	GoodsSubType_DailyTopicTicket             GoodsSubType = 123 //每日话题投票券
	GoodsSubType_DailyTopicScore              GoodsSubType = 124 //每日话题积分
	GoodsSubType_FlipCardTicket               GoodsSubType = 125 //翻翻乐翻牌道具
	GoodsSubType_FlipCardPackageTicket        GoodsSubType = 126 //翻翻乐礼包抽奖券
	GoodsSubType_CampBattleShopToken          GoodsSubType = 127 //阵营战商店代币
	GoodsSubType_FlipCardBattlePass           GoodsSubType = 128 //翻翻乐通行证
	GoodsSubType_FlipCardBattlePassPoint      GoodsSubType = 129 //翻翻乐通行证积分
	GoodsSubType_FlipCardProgressScore        GoodsSubType = 130 //翻翻乐进度积分
	GoodsSubType_FlipCardTurntableTicket      GoodsSubType = 131 //翻翻乐转盘券
	GoodsSubType_DestinyFetterLevel           GoodsSubType = 132 //宿命羁绊升级道具
	GoodsSubType_MultipleProp                 GoodsSubType = 133 //多倍道具
	GoodsSubType_DynamicHero                  GoodsSubType = 134 //动态立绘
	GoodsSubType_UseTreasure                  GoodsSubType = 135 //使用宝箱
	GoodsSubType_OneYuanCard                  GoodsSubType = 136 //一元权益卡
	GoodsSubType_UpPuzzle                     GoodsSubType = 137 //Up养成拼图
	GoodsSubType_PatrolRefreshTicket          GoodsSubType = 138 //探秘刷新券
	GoodsSubType_RedEnvelope                  GoodsSubType = 139 //红包
	GoodsSubType_BoxChipCompose               GoodsSubType = 140 //箱子碎片合成
	GoodsSubType_VeteranSevenDayTaskScore     GoodsSubType = 141 //老手七日任务积分
	GoodsSubType_VeteranHeroRecruitPoint      GoodsSubType = 142 //老手招募积分
	GoodsSubType_HeroStarPass                 GoodsSubType = 143 //英雄星级通行证
	GoodsSubType_SecretArtsPass               GoodsSubType = 144 //秘术通行证
	GoodsSubType_HeroSoulPass                 GoodsSubType = 145 //八门通行证
	GoodsSubType_HeroArtifactUniversalChip    GoodsSubType = 146 //英雄神器万能碎片
	GoodsSubType_UpSignPass                   GoodsSubType = 147 //Up招募签到战令
	GoodsSubType_WishSignPass                 GoodsSubType = 148 //自选招募签到战令
	GoodsSubType_BoxGiftCoin                  GoodsSubType = 149 //格斗家礼物等级代币
	GoodsSubType_AstrologyActivity            GoodsSubType = 150 //占星限时活动
	GoodsSubType_FashionBox                   GoodsSubType = 151 //时装箱
	GoodsSubType_FashionBoxKey                GoodsSubType = 152 //时装箱钥匙
	GoodsSubType_Strong108Ticket              GoodsSubType = 153 //最强108挑战券
	GoodsSubType_Strong108Token               GoodsSubType = 154 //最强108代币
	GoodsSubType_FourSeasonsCoin              GoodsSubType = 155 //四季活动徽章
	GoodsSubType_ExchangeCoin                 GoodsSubType = 156 //兑换币
	GoodsSubType_CommonItem                   GoodsSubType = 157 //普通道具
	GoodsSubType_RichProgress                 GoodsSubType = 158 //存钱罐进度
	GoodsSubType_GameProficiency              GoodsSubType = 159 //街机厅游戏熟练度
	GoodsSubType_NationalCoin                 GoodsSubType = 160 //国庆徽章
	GoodsSubType_TeamLeagueScore              GoodsSubType = 161 //公会团赛积分
	GoodsSubType_SoulCard                     GoodsSubType = 162 //战魂权益卡
	GoodsSubType_DiamondCard                  GoodsSubType = 163 //钻石权益卡
	GoodsSubType_ArcadeAllianceScore          GoodsSubType = 164 //街机厅公会竞技值
	GoodsSubType_Strongest108Score            GoodsSubType = 165 //最强108积分
	GoodsSubType_MegaNurtureChip              GoodsSubType = 166 //Mega养成拼图
	GoodsSubType_RugalSignPass                GoodsSubType = 167 //卢卡尔签到战令
	GoodsSubType_SoulRecastTicket             GoodsSubType = 168 //战魂重铸券
	GoodsSubType_ManorSuperRefreshCard        GoodsSubType = 169 //超刷卡
	GoodsSubType_ManorBattlePowerCard         GoodsSubType = 170 //战力卡
	GoodsSubType_ManorKickCard                GoodsSubType = 171 //踢飞卡
	GoodsSubType_ManorJumpInCard              GoodsSubType = 172 //插队卡
	GoodsSubType_ManorPoopCard                GoodsSubType = 173 //便便卡
	GoodsSubType_ArcadeTicket                 GoodsSubType = 174 //街机厅挑战券
	GoodsSubType_TeamStageTicket              GoodsSubType = 175 //组队关卡门票
	GoodsSubType_TeamLeagueTaskCnt            GoodsSubType = 176 //团赛任务次数
	GoodsSubType_PatrolSupportTimes           GoodsSubType = 177 //探秘协助次数
	GoodsSubType_PatrolDailyTimes             GoodsSubType = 178 //探秘次数
	GoodsSubType_ManorMeteor                  GoodsSubType = 179 //陨石
	GoodsSubType_CrazyFruitTicket             GoodsSubType = 180 //水果机抽奖券
	GoodsSubType_SnowHighPass                 GoodsSubType = 181 //雪人高级战令
	GoodsSubType_TradeCoin                    GoodsSubType = 182 //贸易币
	GoodsSubType_DaSheBackSignPass            GoodsSubType = 183 //大蛇返场签到战令
	GoodsSubType_FashionFragment              GoodsSubType = 184 //时装碎片
	GoodsSubType_SpringFestivalHighPass       GoodsSubType = 185 //春节高级战令
	GoodsSubType_SpringPassScore              GoodsSubType = 186 //春节战令积分
	GoodsSubType_Crazy8SignPass               GoodsSubType = 187 //疯八招募签到战令
	GoodsSubType_RugalBackSignPass            GoodsSubType = 188 //卢卡尔返场签到战令
	GoodsSubType_CommonBattlePass             GoodsSubType = 189 //通用签到战令
	GoodsSubType_ArcadeAllianceScoreItem      GoodsSubType = 190 //街机厅公会竞技值道具
	GoodsSubType_TradeMedal                   GoodsSubType = 191 //贸易银奖章
	GoodsSubType_TradeBigMedal                GoodsSubType = 192 //贸易金奖章
	GoodsSubType_KingWarWeekMedal             GoodsSubType = 193 //皇城周赛勋章
	GoodsSubType_KingWarFinalMedal            GoodsSubType = 194 //皇城决赛勋章
	GoodsSubType_MercenaryToken               GoodsSubType = 195 //皇城佣兵令
	GoodsSubType_PetSoul                      GoodsSubType = 196 //宠物碎片
	GoodsSubType_PetSoulTicket                GoodsSubType = 197 //远古战魂券
	GoodsSubType_PetEgg                       GoodsSubType = 198 //宠物蛋
	GoodsSubType_FreeReviveTimes              GoodsSubType = 199 //皇城免费复活次数
	GoodsSubType_ReviveTicket                 GoodsSubType = 200 //皇城复活券
	GoodsSubType_KingWarSkin                  GoodsSubType = 201 //皇城皮肤
	GoodsSubType_AniAct                       GoodsSubType = 202 //周年庆道具
	GoodsSubType_FashionBoxBind               GoodsSubType = 203 //绑定时装箱
	GoodsSubType_Purifier                     GoodsSubType = 204 //抗扰剂
	GoodsSubType_Transformation               GoodsSubType = 205 //装扮卡
	GoodsSubType_HeroCard                     GoodsSubType = 1   //英雄卡牌
	GoodsSubType_HeroChip                     GoodsSubType = 2   //英雄碎片
	GoodsSubType_HeroSoul                     GoodsSubType = 3   //英雄战魂
	GoodsSubType_HeroArtifact                 GoodsSubType = 4   //英雄神器
	GoodsSubType_HeroArtifactChip             GoodsSubType = 5   //英雄神器碎片
	GoodsSubType_HeroArtifactBreak            GoodsSubType = 6   //英雄神器突破石
	GoodsSubType_HeroArtifactStone            GoodsSubType = 7   //神器神格强化石
	GoodsSubType_HeroRandomSoul               GoodsSubType = 8   //品质随机战魂
	GoodsSubType_HeroMagic                    GoodsSubType = 9   //英雄自选
	GoodsSubType_HeroChipMagic                GoodsSubType = 10  //英雄碎片自选
	GoodsSubType_HeroChipRandom               GoodsSubType = 11  //英雄碎片随机
	GoodsSubType_HeroArtifactMagic            GoodsSubType = 12  //英雄神器自选
	GoodsSubType_HeroArtifactChipMagic        GoodsSubType = 13  //英雄神器碎片自选
	GoodsSubType_HeroArtifactChipRandom       GoodsSubType = 14  //英雄神器碎片随机
	GoodsSubType_HeroArtifactCompose          GoodsSubType = 15  //英雄神器合成
	GoodsSubType_HeroRandom                   GoodsSubType = 16  //英雄随机
	GoodsSubType_HeroUniversalChip            GoodsSubType = 17  //英雄万能碎片
	GoodsSubType_HeroRarityStone              GoodsSubType = 18  //英雄升品石
	GoodsSubType_ItemMagic                    GoodsSubType = 19  //道具自选
	GoodsSubType_ItemRandom                   GoodsSubType = 20  //道具随机
	GoodsSubType_NestChip                     GoodsSubType = 21  //音巢芯片
	GoodsSubType_NestRefineAttrLock           GoodsSubType = 22  //音巢芯片属性锁定道具
	GoodsSubType_HeroSoulRefining             GoodsSubType = 23  //战魂洗练道具
	GoodsSubType_DreamHeroRandom              GoodsSubType = 24  //许愿英雄随机
	GoodsSubType_HeroRights                   GoodsSubType = 25  //英雄权益
	GoodsSubType_NestRefineAttrRefine         GoodsSubType = 26  //音巢芯片属性洗练道具
	GoodsSubType_NestStarup                   GoodsSubType = 27  //芯片升星消耗道具ID
	GoodsSubType_HeroSwapStone                GoodsSubType = 28  //英雄置换石
	GoodsSubType_BindHeroSoul                 GoodsSubType = 29  //绑定战魂
	GoodsSubType_BindHeroRandomSoul           GoodsSubType = 30  //绑定品质随机战魂
	GoodsSubType_ShopGood                     GoodsSubType = 1   //商品
	GoodsSubType_SkinMap                      GoodsSubType = 1   //地图
	GoodsSubType_SkinRole                     GoodsSubType = 2   //角色
	GoodsSubType_LevelPositive                GoodsSubType = 1   //主动道具
	GoodsSubType_LevelPassive                 GoodsSubType = 2   //被动道具
	GoodsSubType_FashionAvatar                GoodsSubType = 1   //时装头像
	GoodsSubType_FashionAvatarFrame           GoodsSubType = 2   //时装头像框
	GoodsSubType_FashionChatBubble            GoodsSubType = 3   //时装聊天气泡
	GoodsSubType_FashionCutscene              GoodsSubType = 4   //PVP过场动画
	GoodsSubType_FashionDamageEffects         GoodsSubType = 5   //伤害数字
	GoodsSubType_FashionDynamicBackground     GoodsSubType = 6   //个人信息动态背景
	GoodsSubType_FashionWing                  GoodsSubType = 7   //翅膀
	GoodsSubType_FashionKillSpecialEffects    GoodsSubType = 8   //必杀特效
	GoodsSubType_FashionFinalEffects          GoodsSubType = 9   //终结特效
	GoodsSubType_FashionEntryAction           GoodsSubType = 10  //入场动作
	GoodsSubType_FashionStandbyAction         GoodsSubType = 11  //待机动作
	GoodsSubType_FashionVerticalPainting      GoodsSubType = 12  //时装立绘
	GoodsSubType_FashionStandinPuppet         GoodsSubType = 13  //替身人偶
	GoodsSubType_FashionableDress             GoodsSubType = 14  //时装
	GoodsSubType_FashionManorSkin             GoodsSubType = 15  //领地皮肤
	GoodsSubType_FashionHomeSkin              GoodsSubType = 16  //主页皮肤
	GoodsSubType_FashionAoYi                  GoodsSubType = 17  //奥义技
	GoodsSubType_FashionAidSkill              GoodsSubType = 18  //援助技
	GoodsSubType_FashionHeroPart1             GoodsSubType = 19  //变身1
	GoodsSubType_FashionHeroPart2             GoodsSubType = 20  //变身2
	GoodsSubType_FashionHeroPart3             GoodsSubType = 21  //变身3
	GoodsSubType_FashionHeroPart4             GoodsSubType = 22  //变身4
	GoodsSubType_Slime                        GoodsSubType = 1   //史莱姆
	GoodsSubType_GrayPet                      GoodsSubType = 10  //灰宠物
	GoodsSubType_GreenPet                     GoodsSubType = 11  //绿宠物
	GoodsSubType_BluePet                      GoodsSubType = 12  //蓝宠物
	GoodsSubType_PurplePet                    GoodsSubType = 13  //紫宠物
	GoodsSubType_OrangePet                    GoodsSubType = 14  //金宠物
	GoodsSubType_GoldPet                      GoodsSubType = 15  //橙宠物
	GoodsSubType_RedPet                       GoodsSubType = 16  //红宠物
	GoodsSubType_FullColourPet                GoodsSubType = 17  //彩宠物
	GoodsSubType_RandomPet                    GoodsSubType = 18  //随机宠物
	GoodsSubType_PetMagic                     GoodsSubType = 19  //自选宠物
	GoodsSubType_SpiritDew                    GoodsSubType = 206 //甘露
	GoodsSubType_END                          GoodsSubType = 207 //枚举结束
)

// 基础配置-拳皇/commonEnum
type Rarity int

const (
	Rarity_BEGIN    Rarity = 0   //枚举开始
	Rarity_Gray     Rarity = 0   //灰
	Rarity_Green    Rarity = 1   //绿
	Rarity_Blue     Rarity = 2   //蓝
	Rarity_Purple   Rarity = 3   //紫
	Rarity_Orange   Rarity = 4   //橙
	Rarity_Red      Rarity = 5   //红
	Rarity_Rainbow  Rarity = 6   //彩
	Rarity_Rainbow1 Rarity = 7   //彩1
	Rarity_Rainbow2 Rarity = 8   //彩2
	Rarity_Rainbow3 Rarity = 9   //彩3
	Rarity_Rainbow4 Rarity = 10  //彩4
	Rarity_Rainbow5 Rarity = 11  //彩5
	Rarity_Rainbow6 Rarity = 12  //彩6
	Rarity_Yellow   Rarity = 99  //黄
	Rarity_END      Rarity = 100 //枚举结束
)

// 基础配置-拳皇/commonEnum
type LDRarity int

const (
	LDRarity_BEGIN     LDRarity = 1 //枚举开始
	LDRarity_Good      LDRarity = 1 //优秀
	LDRarity_Excellent LDRarity = 2 //精良
	LDRarity_Rare      LDRarity = 3 //稀有
	LDRarity_Legend    LDRarity = 4 //传奇
	LDRarity_Epic      LDRarity = 5 //史诗
	LDRarity_END       LDRarity = 6 //枚举结束
)

// 基础配置-拳皇/commonEnum
type HeroSoulType int

const (
	HeroSoulType_BEGIN     HeroSoulType = 0 //枚举开始
	HeroSoulType_Effect    HeroSoulType = 0 //天赋
	HeroSoulType_Attack    HeroSoulType = 1 //攻击
	HeroSoulType_Energy    HeroSoulType = 2 //能量
	HeroSoulType_Survive   HeroSoulType = 3 //生命
	HeroSoulType_Universal HeroSoulType = 4 //万能
	HeroSoulType_NoLimit   HeroSoulType = 5 //任意
	HeroSoulType_END       HeroSoulType = 6 //枚举结束
)

// 基础配置-拳皇/commonEnum
type HeroSoulGrade int

const (
	HeroSoulGrade_BEGIN  HeroSoulGrade = 0 //枚举开始
	HeroSoulGrade_GradeD HeroSoulGrade = 0 //D
	HeroSoulGrade_GradeC HeroSoulGrade = 1 //C
	HeroSoulGrade_GradeB HeroSoulGrade = 2 //B
	HeroSoulGrade_GradeA HeroSoulGrade = 3 //A
	HeroSoulGrade_GradeS HeroSoulGrade = 4 //S
	HeroSoulGrade_END    HeroSoulGrade = 5 //枚举结束
)

// 基础配置-拳皇/commonEnum
type HeroSoulEffectType int

const (
	HeroSoulEffectType_BEGIN                HeroSoulEffectType = 1 //枚举开始
	HeroSoulEffectType_ComposeProbAdd       HeroSoulEffectType = 1 //合成概率增加
	HeroSoulEffectType_BaseAttrBonus        HeroSoulEffectType = 2 //基础属性加成
	HeroSoulEffectType_UniversalTag         HeroSoulEffectType = 3 //万能词条
	HeroSoulEffectType_ExtraSkill           HeroSoulEffectType = 4 //额外技能
	HeroSoulEffectType_AsFoodProbAdd        HeroSoulEffectType = 5 //作为狗粮概率增加
	HeroSoulEffectType_UniversalTagAndBonus HeroSoulEffectType = 6 //万能词条及属性加成
	HeroSoulEffectType_END                  HeroSoulEffectType = 7 //枚举结束
)

// 基础配置-拳皇/commonEnum
type GoodsConfigId int

const (
	GoodsConfigId_BEGIN                     GoodsConfigId = 1         //枚举开始
	GoodsConfigId_Gold                      GoodsConfigId = 10100000  //金币
	GoodsConfigId_Diamond                   GoodsConfigId = 10200000  //钻石
	GoodsConfigId_Exp                       GoodsConfigId = 10300000  //玩家经验
	GoodsConfigId_HeroExp                   GoodsConfigId = 10400000  //英雄经验
	GoodsConfigId_HeroAtlasBadgePoint       GoodsConfigId = 10500000  //英雄图鉴点数
	GoodsConfigId_TowerTicket               GoodsConfigId = 10600000  //爬塔挑战券
	GoodsConfigId_ManorResTicket            GoodsConfigId = 10700000  //领地资源刷新券
	GoodsConfigId_MonthCardPoint            GoodsConfigId = 10800000  //月卡积分
	GoodsConfigId_ArenaTicket               GoodsConfigId = 10900000  //竞技场挑战券
	GoodsConfigId_EnjoymentCard             GoodsConfigId = 11100000  //畅玩卡
	GoodsConfigId_CashCharge                GoodsConfigId = 11300000  //现金充值
	GoodsConfigId_WorldBossSweepTicket      GoodsConfigId = 11500000  //世界boss扫荡券
	GoodsConfigId_HeroSoulRefiningExpend    GoodsConfigId = 22300000  //英雄战魂洗练消耗道具
	GoodsConfigId_AfkPayExpend              GoodsConfigId = 10200000  //付费离线挂机消耗道具
	GoodsConfigId_VipExp                    GoodsConfigId = 12100000  //VIP经验
	GoodsConfigId_ArtifactBreakExpend       GoodsConfigId = 20600000  //神器突破消耗
	GoodsConfigId_AstrologyItem             GoodsConfigId = 10100000  //神器占星货币
	GoodsConfigId_HolyCrystal               GoodsConfigId = 13200000  //占星升档道具
	GoodsConfigId_AllianceTechUniversal     GoodsConfigId = 13400204  //公会科技万能碎片
	GoodsConfigId_AllianceTechTicket        GoodsConfigId = 13500000  //公会科技能量
	GoodsConfigId_BoxGiftExp                GoodsConfigId = 10100012  //收藏宝箱经验
	GoodsConfigId_HeroTicket                GoodsConfigId = 12500001  //英雄招募券
	GoodsConfigId_ArtifactTicket            GoodsConfigId = 10100004  //占星之光
	GoodsConfigId_Box                       GoodsConfigId = 14500000  //格斗家收藏盒
	GoodsConfigId_FishSliver                GoodsConfigId = 14700000  //钓鱼银币
	GoodsConfigId_RedHeroMagicChip          GoodsConfigId = 21700004  //英雄万能碎片
	GoodsConfigId_ReduceTimeTicket          GoodsConfigId = 15300001  //加速券
	GoodsConfigId_BountyChallengeScore      GoodsConfigId = 15600000  //赏金值
	GoodsConfigId_CrossArenaTicket          GoodsConfigId = 15900000  //跨服竞技场门票
	GoodsConfigId_Contribution              GoodsConfigId = 16100000  //贡献
	GoodsConfigId_ContributeChallengeTicket GoodsConfigId = 16200000  //贡献挑战入场券
	GoodsConfigId_AthenaSweepTicket         GoodsConfigId = 16300000  //亚呆娜扫荡券
	GoodsConfigId_PlaygroundTicket          GoodsConfigId = 16400000  //乐园券
	GoodsConfigId_RoomCard                  GoodsConfigId = 16500000  //房卡
	GoodsConfigId_Coin                      GoodsConfigId = 16800000  //硬币
	GoodsConfigId_LuckyStone                GoodsConfigId = 16900000  //幸运石
	GoodsConfigId_RebornCard                GoodsConfigId = 70800000  //复活卡
	GoodsConfigId_CleanCard                 GoodsConfigId = 70900000  //清除卡
	GoodsConfigId_RandomCard                GoodsConfigId = 71000000  //随机卡
	GoodsConfigId_GoldCard                  GoodsConfigId = 71100000  //金币卡
	GoodsConfigId_RoomHeroCard              GoodsConfigId = 71200000  //英雄卡
	GoodsConfigId_GridCard                  GoodsConfigId = 71300000  //格子卡
	GoodsConfigId_CopyCard                  GoodsConfigId = 71300001  //复制卡
	GoodsConfigId_FundCard                  GoodsConfigId = 71300002  //基金卡
	GoodsConfigId_LoongCard                 GoodsConfigId = 71300003  //龙血卡
	GoodsConfigId_LegendCard                GoodsConfigId = 71300004  //传奇卡
	GoodsConfigId_GoldHammer                GoodsConfigId = 71500020  //黄金锤子
	GoodsConfigId_JZYGoneFish               GoodsConfigId = 11700001  //姜子牙跑鱼
	GoodsConfigId_JZYGoldFish               GoodsConfigId = 11700002  //姜子牙金鱼
	GoodsConfigId_ChooseSkinExperience      GoodsConfigId = 17100000  //自选皮肤体验
	GoodsConfigId_FormationRefresh          GoodsConfigId = 17200000  //阵容刷新券
	GoodsConfigId_NestTicket                GoodsConfigId = 22200002  //音巢挑战门票
	GoodsConfigId_AllianceBossExp           GoodsConfigId = 18900000  //公会雕像经验
	GoodsConfigId_LDCoin                    GoodsConfigId = 18500000  //LD玩法代币
	GoodsConfigId_LDSingleMatchCnt          GoodsConfigId = 18600000  //LD单人匹配次数
	GoodsConfigId_LDFriendMatchCnt          GoodsConfigId = 18700000  //LD好友匹配次数
	GoodsConfigId_LDInviteeMatchCnt         GoodsConfigId = 18800000  //LD受邀匹配次数
	GoodsConfigId_StarCoin                  GoodsConfigId = 19000000  //星钻
	GoodsConfigId_InvalidItem               GoodsConfigId = 99999999  //非法道具
	GoodsConfigId_ArtifactStone             GoodsConfigId = 21500000  //神器源石
	GoodsConfigId_AllianceFightScore        GoodsConfigId = 19800000  //公会对决积分
	GoodsConfigId_AllianceContribute        GoodsConfigId = 10100013  //公会贡献道具
	GoodsConfigId_SpTDCoin                  GoodsConfigId = 10100016  //南镇金币
	GoodsConfigId_AutoFishTicket            GoodsConfigId = 110400000 //自动钓鱼券
	GoodsConfigId_BoxGiftBlessExp           GoodsConfigId = 10100017  //祝福经验
	GoodsConfigId_SpTDCurrency              GoodsConfigId = 10100021  //南镇货币
	GoodsConfigId_SpTDSweep                 GoodsConfigId = 10100022  //南镇扫荡券
	GoodsConfigId_RecruitExp                GoodsConfigId = 111500000 //招募经验
	GoodsConfigId_FightingPlaceTicket       GoodsConfigId = 116000000 //格斗圣殿挑战券
	GoodsConfigId_DailyTopicTicket          GoodsConfigId = 123000000 //每日话题投票券
	GoodsConfigId_DailyTopicScore           GoodsConfigId = 124000000 //每日话题积分
	GoodsConfigId_CampBattleShopToken       GoodsConfigId = 112700000 //阵营战商店代币
	GoodsConfigId_SeasonsTicket             GoodsConfigId = 12500002  //四季活动心愿召唤券
	GoodsConfigId_SeasonsExchangeCoin       GoodsConfigId = 10100003  //庆典兑换币
	GoodsConfigId_SpringFestivalCoin        GoodsConfigId = 10100023  //春芽徽章
	GoodsConfigId_WinterFestivalCoin        GoodsConfigId = 10100032  //冬日徽章
	GoodsConfigId_SummerFestivalCoin        GoodsConfigId = 10100033  //夏荷徽章
	GoodsConfigId_FlipCardBoxTicket         GoodsConfigId = 112500001 //翻翻乐-宝匣抽奖券
	GoodsConfigId_FlipCardBattlePassPoint   GoodsConfigId = 112900001 //翻翻乐战令积分
	GoodsConfigId_FlipFlipCardProgressPoint GoodsConfigId = 113000001 //翻翻乐进度积分
	GoodsConfigId_FlipCardTurntableTicket   GoodsConfigId = 113100001 //翻翻乐转盘券
	GoodsConfigId_ArtifactMagicChip         GoodsConfigId = 214600001 //神器万能碎片
	GoodsConfigId_UpSignBattlePassScore     GoodsConfigId = 19900001  //Up招募签到战令积分
	GoodsConfigId_WishSignBattlePassScore   GoodsConfigId = 19900002  //自选招募签到战令积分
	GoodsConfigId_AstrologyActivity         GoodsConfigId = 12500005  //神器祈福活动道具
	GoodsConfigId_BoxGiftExtraCoin          GoodsConfigId = 149000000 //格斗家礼物代币
	GoodsConfigId_UniversalHeroSoul         GoodsConfigId = 20300045  //万能本源
	GoodsConfigId_RecruitPrivilageCard      GoodsConfigId = 111200000 //招募冲刺特权卡
	GoodsConfigId_ArtifactPrivilageCard     GoodsConfigId = 111300000 //神器冲刺特权卡
	GoodsConfigId_MiShuPrivilageCard        GoodsConfigId = 111400000 //秘术冲刺特权卡
	GoodsConfigId_TeamLeagueScore           GoodsConfigId = 116100000 //团赛积分
	GoodsConfigId_PromotionTicket           GoodsConfigId = 21800000  //升品凭证
	GoodsConfigId_SRTicket                  GoodsConfigId = 21800001  //SR凭证
	GoodsConfigId_SSRTicket                 GoodsConfigId = 21800002  //SSR凭证
	GoodsConfigId_URTicket                  GoodsConfigId = 21800003  //UR凭证
	GoodsConfigId_ManorSuperRefreshCard     GoodsConfigId = 116900000 //超刷卡
	GoodsConfigId_ManorBattlePowerCard      GoodsConfigId = 117000000 //战力卡
	GoodsConfigId_ManorJumpQueueCard        GoodsConfigId = 117200000 //插队卡
	GoodsConfigId_ManorKickCard             GoodsConfigId = 117100000 //踢飞卡
	GoodsConfigId_ManorRubbishCard          GoodsConfigId = 117300000 //便便卡
	GoodsConfigId_TeamStageTicket           GoodsConfigId = 117500000 //组队关卡门票
	GoodsConfigId_ManorMeteor1              GoodsConfigId = 10110001  //低级陨石
	GoodsConfigId_ManorMeteor2              GoodsConfigId = 10110002  //中级陨石
	GoodsConfigId_ManorMeteor3              GoodsConfigId = 10110003  //高级陨石
	GoodsConfigId_ManorMeteor4              GoodsConfigId = 10110004  //终级陨石
	GoodsConfigId_FashionChip               GoodsConfigId = 11840001  //装扮碎片
	GoodsConfigId_SpringFestivalNestTicket  GoodsConfigId = 115000111 //春节音巢入场券
	GoodsConfigId_ReviveTicket              GoodsConfigId = 120000000 //皇城复活券
	GoodsConfigId_Purifier                  GoodsConfigId = 10100030  //抗扰剂
	GoodsConfigId_SpiritDew                 GoodsConfigId = 120600000 //灵池甘露
	GoodsConfigId_END                       GoodsConfigId = 214600002 //枚举结束
)

// 基础配置-拳皇/commonEnum
type RecruitType int

const (
	RecruitType_BEGIN     RecruitType = 1  //枚举开始
	RecruitType_Free      RecruitType = 1  //免费招募
	RecruitType_AD        RecruitType = 2  //广告招募
	RecruitType_Senior    RecruitType = 3  //高级招募
	RecruitType_LimitTime RecruitType = 4  //限时招募
	RecruitType_Super     RecruitType = 5  //大保底
	RecruitType_Up        RecruitType = 6  //Up招募
	RecruitType_Wish      RecruitType = 7  //心愿招募
	RecruitType_UpSuper   RecruitType = 8  //Up超级招募
	RecruitType_Seasons   RecruitType = 9  //四季活动招募
	RecruitType_FundSuper RecruitType = 10 //基金大保底
	RecruitType_WishSuper RecruitType = 11 //自选超级招募
	RecruitType_END       RecruitType = 12 //枚举结束
)

func (e EMPlatform) String() string {

	if e == EMPlatform_All {
		return "All"
	}
	if e == EMPlatform_Wechat {
		return "Wechat"
	}
	if e == EMPlatform_DouYin {
		return "DouYin"
	}
	if e == EMPlatform_AliPay {
		return "AliPay"
	}
	if e == EMPlatform_MeiTuan {
		return "MeiTuan"
	}

	return ""

}

func GetEMPlatform(tp string) (EMPlatform, bool) {
	switch tp {
	case "All":
		return EMPlatform_All, true
	case "Wechat":
		return EMPlatform_Wechat, true
	case "DouYin":
		return EMPlatform_DouYin, true
	case "AliPay":
		return EMPlatform_AliPay, true
	case "MeiTuan":
		return EMPlatform_MeiTuan, true

	}
	return -1, false

}

func GetEMPlatform_ByAlias(tp string) (EMPlatform, bool) {
	switch tp {
	case "全平台":
		return EMPlatform_All, true
	case "微信":
		return EMPlatform_Wechat, true
	case "抖音":
		return EMPlatform_DouYin, true
	case "支付宝":
		return EMPlatform_AliPay, true
	case "美团":
		return EMPlatform_MeiTuan, true

	}
	return -1, false

}

func (e GoodsType) String() string {

	if e == GoodsType_Value {
		return "Value"
	}
	if e == GoodsType_Hero {
		return "Hero"
	}
	if e == GoodsType_Shop {
		return "Shop"
	}
	if e == GoodsType_Tech {
		return "Tech"
	}
	if e == GoodsType_BattleSkin {
		return "BattleSkin"
	}
	if e == GoodsType_Level {
		return "Level"
	}
	if e == GoodsType_FashionClothes {
		return "FashionClothes"
	}
	if e == GoodsType_Pet {
		return "Pet"
	}

	return ""

}

func GetGoodsType(tp string) (GoodsType, bool) {
	switch tp {
	case "Value":
		return GoodsType_Value, true
	case "Hero":
		return GoodsType_Hero, true
	case "Shop":
		return GoodsType_Shop, true
	case "Tech":
		return GoodsType_Tech, true
	case "BattleSkin":
		return GoodsType_BattleSkin, true
	case "Level":
		return GoodsType_Level, true
	case "FashionClothes":
		return GoodsType_FashionClothes, true
	case "Pet":
		return GoodsType_Pet, true

	}
	return -1, false

}

func GetGoodsType_ByAlias(tp string) (GoodsType, bool) {
	switch tp {
	case "数值":
		return GoodsType_Value, true
	case "英雄":
		return GoodsType_Hero, true
	case "商店":
		return GoodsType_Shop, true
	case "科技":
		return GoodsType_Tech, true
	case "关卡皮肤":
		return GoodsType_BattleSkin, true
	case "关卡道具":
		return GoodsType_Level, true
	case "装扮":
		return GoodsType_FashionClothes, true
	case "宠物":
		return GoodsType_Pet, true

	}
	return -1, false

}

func (e GoodsSubType) String() string {

	if e == GoodsSubType_Empty {
		return "Empty"
	}
	if e == GoodsSubType_Gold {
		return "Gold"
	}
	if e == GoodsSubType_Diamond {
		return "Diamond"
	}
	if e == GoodsSubType_Exp {
		return "Exp"
	}
	if e == GoodsSubType_HeroExp {
		return "HeroExp"
	}
	if e == GoodsSubType_HeroBadgePoint {
		return "HeroBadgePoint"
	}
	if e == GoodsSubType_TowerTicket {
		return "TowerTicket"
	}
	if e == GoodsSubType_ManorResTicket {
		return "ManorResTicket"
	}
	if e == GoodsSubType_MonthCardPoint {
		return "MonthCardPoint"
	}
	if e == GoodsSubType_ArenaTicket {
		return "ArenaTicket"
	}
	if e == GoodsSubType_EquipStone {
		return "EquipStone"
	}
	if e == GoodsSubType_EnjoymentCard {
		return "EnjoymentCard"
	}
	if e == GoodsSubType_VipRight {
		return "VipRight"
	}
	if e == GoodsSubType_CashCharge {
		return "CashCharge"
	}
	if e == GoodsSubType_GameCartridge {
		return "GameCartridge"
	}
	if e == GoodsSubType_WorldBossSweepTicket {
		return "WorldBossSweepTicket"
	}
	if e == GoodsSubType_MonthCard {
		return "MonthCard"
	}
	if e == GoodsSubType_LifetimeCard {
		return "LifetimeCard"
	}
	if e == GoodsSubType_HighLevelBattlePass {
		return "HighLevelBattlePass"
	}
	if e == GoodsSubType_AllianceDailyTaskProgress {
		return "AllianceDailyTaskProgress"
	}
	if e == GoodsSubType_AllianceExp {
		return "AllianceExp"
	}
	if e == GoodsSubType_VipExp {
		return "VipExp"
	}
	if e == GoodsSubType_PromotionExp {
		return "PromotionExp"
	}
	if e == GoodsSubType_HighGrowFund1BattlePass {
		return "HighGrowFund1BattlePass"
	}
	if e == GoodsSubType_HighGrowFund2BattlePass {
		return "HighGrowFund2BattlePass"
	}
	if e == GoodsSubType_SeniorRecruitTicket {
		return "SeniorRecruitTicket"
	}
	if e == GoodsSubType_HighChallengeFund1BattlePass {
		return "HighChallengeFund1BattlePass"
	}
	if e == GoodsSubType_HighChallengeFund2BattlePass {
		return "HighChallengeFund2BattlePass"
	}
	if e == GoodsSubType_HighHonorFund1BattlePass {
		return "HighHonorFund1BattlePass"
	}
	if e == GoodsSubType_HighHonorFund2BattlePass {
		return "HighHonorFund2BattlePass"
	}
	if e == GoodsSubType_DailyActiveScore {
		return "DailyActiveScore"
	}
	if e == GoodsSubType_WeeklyActiveScore {
		return "WeeklyActiveScore"
	}
	if e == GoodsSubType_Astrology {
		return "Astrology"
	}
	if e == GoodsSubType_BoxGift {
		return "BoxGift"
	}
	if e == GoodsSubType_AllianceTechFragment {
		return "AllianceTechFragment"
	}
	if e == GoodsSubType_AllianceTechTicket {
		return "AllianceTechTicket"
	}
	if e == GoodsSubType_HeroWeekCard {
		return "HeroWeekCard"
	}
	if e == GoodsSubType_FishWeekCard {
		return "FishWeekCard"
	}
	if e == GoodsSubType_TreasureWeekCard {
		return "TreasureWeekCard"
	}
	if e == GoodsSubType_ActivityPoint {
		return "ActivityPoint"
	}
	if e == GoodsSubType_AfkExpPackage {
		return "AfkExpPackage"
	}
	if e == GoodsSubType_AfkGoldPackage {
		return "AfkGoldPackage"
	}
	if e == GoodsSubType_WeekHeroRecruitPoint {
		return "WeekHeroRecruitPoint"
	}
	if e == GoodsSubType_WeekArtifactDrawPoint {
		return "WeekArtifactDrawPoint"
	}
	if e == GoodsSubType_SevenDayTaskScore {
		return "SevenDayTaskScore"
	}
	if e == GoodsSubType_Box {
		return "Box"
	}
	if e == GoodsSubType_FishEnergy {
		return "FishEnergy"
	}
	if e == GoodsSubType_FishSilver {
		return "FishSilver"
	}
	if e == GoodsSubType_FishPearl {
		return "FishPearl"
	}
	if e == GoodsSubType_FishBlackPearl {
		return "FishBlackPearl"
	}
	if e == GoodsSubType_EquipStoneBox {
		return "EquipStoneBox"
	}
	if e == GoodsSubType_FishOrderRefreshTicket {
		return "FishOrderRefreshTicket"
	}
	if e == GoodsSubType_FishOrderResetTicket {
		return "FishOrderResetTicket"
	}
	if e == GoodsSubType_ReduceTimeTicket {
		return "ReduceTimeTicket"
	}
	if e == GoodsSubType_FishBait {
		return "FishBait"
	}
	if e == GoodsSubType_FishGalaPoint {
		return "FishGalaPoint"
	}
	if e == GoodsSubType_BountyChallengeScore {
		return "BountyChallengeScore"
	}
	if e == GoodsSubType_SoulPotion {
		return "SoulPotion"
	}
	if e == GoodsSubType_AllianceCheerExp {
		return "AllianceCheerExp"
	}
	if e == GoodsSubType_CrossArenaTicket {
		return "CrossArenaTicket"
	}
	if e == GoodsSubType_CrossArenaMedal {
		return "CrossArenaMedal"
	}
	if e == GoodsSubType_Contribute {
		return "Contribute"
	}
	if e == GoodsSubType_ContributeChallengeTicket {
		return "ContributeChallengeTicket"
	}
	if e == GoodsSubType_AthenaSweepTicket {
		return "AthenaSweepTicket"
	}
	if e == GoodsSubType_PlaygroundTicket {
		return "PlaygroundTicket"
	}
	if e == GoodsSubType_RoomCard {
		return "RoomCard"
	}
	if e == GoodsSubType_HeroSkin {
		return "HeroSkin"
	}
	if e == GoodsSubType_Collection {
		return "Collection"
	}
	if e == GoodsSubType_Coin {
		return "Coin"
	}
	if e == GoodsSubType_LuckyStone {
		return "LuckyStone"
	}
	if e == GoodsSubType_JZYFish {
		return "JZYFish"
	}
	if e == GoodsSubType_ChooseSkinExperience {
		return "ChooseSkinExperience"
	}
	if e == GoodsSubType_Avatar {
		return "Avatar"
	}
	if e == GoodsSubType_AvatarFrame {
		return "AvatarFrame"
	}
	if e == GoodsSubType_Title {
		return "Title"
	}
	if e == GoodsSubType_FormationRefresh {
		return "FormationRefresh"
	}
	if e == GoodsSubType_ChatBuble {
		return "ChatBuble"
	}
	if e == GoodsSubType_NewbieHeroRecruitPoint {
		return "NewbieHeroRecruitPoint"
	}
	if e == GoodsSubType_AwakenStone {
		return "AwakenStone"
	}
	if e == GoodsSubType_ArtifactBattlePass {
		return "ArtifactBattlePass"
	}
	if e == GoodsSubType_CollectBattlePass {
		return "CollectBattlePass"
	}
	if e == GoodsSubType_NestBattlePass {
		return "NestBattlePass"
	}
	if e == GoodsSubType_RecruitBattlePass {
		return "RecruitBattlePass"
	}
	if e == GoodsSubType_WeekAllianceTechPoint {
		return "WeekAllianceTechPoint"
	}
	if e == GoodsSubType_Treasure {
		return "Treasure"
	}
	if e == GoodsSubType_LDCoin {
		return "LDCoin"
	}
	if e == GoodsSubType_LDSingleMatchCnt {
		return "LDSingleMatchCnt"
	}
	if e == GoodsSubType_LDFriendMatchCnt {
		return "LDFriendMatchCnt"
	}
	if e == GoodsSubType_LDInviteeMatchCnt {
		return "LDInviteeMatchCnt"
	}
	if e == GoodsSubType_AllianceStatueExp {
		return "AllianceStatueExp"
	}
	if e == GoodsSubType_StarCoin {
		return "StarCoin"
	}
	if e == GoodsSubType_ArtifactPassPoint {
		return "ArtifactPassPoint"
	}
	if e == GoodsSubType_FameMedal {
		return "FameMedal"
	}
	if e == GoodsSubType_ManorStone {
		return "ManorStone"
	}
	if e == GoodsSubType_ManorWood {
		return "ManorWood"
	}
	if e == GoodsSubType_ManorFood {
		return "ManorFood"
	}
	if e == GoodsSubType_AllianceTech {
		return "AllianceTech"
	}
	if e == GoodsSubType_BoxGiftExp {
		return "BoxGiftExp"
	}
	if e == GoodsSubType_AllianceDuelScore {
		return "AllianceDuelScore"
	}
	if e == GoodsSubType_BattlePassScore {
		return "BattlePassScore"
	}
	if e == GoodsSubType_WorldBossCoin {
		return "WorldBossCoin"
	}
	if e == GoodsSubType_FishCard {
		return "FishCard"
	}
	if e == GoodsSubType_UpRecruitTicket {
		return "UpRecruitTicket"
	}
	if e == GoodsSubType_WishRecruitTicket {
		return "WishRecruitTicket"
	}
	if e == GoodsSubType_AutoFishTicket {
		return "AutoFishTicket"
	}
	if e == GoodsSubType_ArtsTrialTimes {
		return "ArtsTrialTimes"
	}
	if e == GoodsSubType_BoxGiftBlessExp {
		return "BoxGiftBlessExp"
	}
	if e == GoodsSubType_NestChallengeTicket {
		return "NestChallengeTicket"
	}
	if e == GoodsSubType_NestChallengeKey {
		return "NestChallengeKey"
	}
	if e == GoodsSubType_NestSweepTicket {
		return "NestSweepTicket"
	}
	if e == GoodsSubType_AllianceBossAttackTimes {
		return "AllianceBossAttackTimes"
	}
	if e == GoodsSubType_SpSweepTicket {
		return "SpSweepTicket"
	}
	if e == GoodsSubType_WeekHeroRecruitCard {
		return "WeekHeroRecruitCard"
	}
	if e == GoodsSubType_WeekArtifactDrawCard {
		return "WeekArtifactDrawCard"
	}
	if e == GoodsSubType_WeekAllianceTechCard {
		return "WeekAllianceTechCard"
	}
	if e == GoodsSubType_RecruitExp {
		return "RecruitExp"
	}
	if e == GoodsSubType_FightingPlaceTicket {
		return "FightingPlaceTicket"
	}
	if e == GoodsSubType_TenDrawCard {
		return "TenDrawCard"
	}
	if e == GoodsSubType_Stamina {
		return "Stamina"
	}
	if e == GoodsSubType_StaminaMax {
		return "StaminaMax"
	}
	if e == GoodsSubType_FameForgeTicket {
		return "FameForgeTicket"
	}
	if e == GoodsSubType_TimesTicket {
		return "TimesTicket"
	}
	if e == GoodsSubType_FameGem {
		return "FameGem"
	}
	if e == GoodsSubType_DailyTopicTicket {
		return "DailyTopicTicket"
	}
	if e == GoodsSubType_DailyTopicScore {
		return "DailyTopicScore"
	}
	if e == GoodsSubType_FlipCardTicket {
		return "FlipCardTicket"
	}
	if e == GoodsSubType_FlipCardPackageTicket {
		return "FlipCardPackageTicket"
	}
	if e == GoodsSubType_CampBattleShopToken {
		return "CampBattleShopToken"
	}
	if e == GoodsSubType_FlipCardBattlePass {
		return "FlipCardBattlePass"
	}
	if e == GoodsSubType_FlipCardBattlePassPoint {
		return "FlipCardBattlePassPoint"
	}
	if e == GoodsSubType_FlipCardProgressScore {
		return "FlipCardProgressScore"
	}
	if e == GoodsSubType_FlipCardTurntableTicket {
		return "FlipCardTurntableTicket"
	}
	if e == GoodsSubType_DestinyFetterLevel {
		return "DestinyFetterLevel"
	}
	if e == GoodsSubType_MultipleProp {
		return "MultipleProp"
	}
	if e == GoodsSubType_DynamicHero {
		return "DynamicHero"
	}
	if e == GoodsSubType_UseTreasure {
		return "UseTreasure"
	}
	if e == GoodsSubType_OneYuanCard {
		return "OneYuanCard"
	}
	if e == GoodsSubType_UpPuzzle {
		return "UpPuzzle"
	}
	if e == GoodsSubType_PatrolRefreshTicket {
		return "PatrolRefreshTicket"
	}
	if e == GoodsSubType_RedEnvelope {
		return "RedEnvelope"
	}
	if e == GoodsSubType_BoxChipCompose {
		return "BoxChipCompose"
	}
	if e == GoodsSubType_VeteranSevenDayTaskScore {
		return "VeteranSevenDayTaskScore"
	}
	if e == GoodsSubType_VeteranHeroRecruitPoint {
		return "VeteranHeroRecruitPoint"
	}
	if e == GoodsSubType_HeroStarPass {
		return "HeroStarPass"
	}
	if e == GoodsSubType_SecretArtsPass {
		return "SecretArtsPass"
	}
	if e == GoodsSubType_HeroSoulPass {
		return "HeroSoulPass"
	}
	if e == GoodsSubType_HeroArtifactUniversalChip {
		return "HeroArtifactUniversalChip"
	}
	if e == GoodsSubType_UpSignPass {
		return "UpSignPass"
	}
	if e == GoodsSubType_WishSignPass {
		return "WishSignPass"
	}
	if e == GoodsSubType_BoxGiftCoin {
		return "BoxGiftCoin"
	}
	if e == GoodsSubType_AstrologyActivity {
		return "AstrologyActivity"
	}
	if e == GoodsSubType_FashionBox {
		return "FashionBox"
	}
	if e == GoodsSubType_FashionBoxKey {
		return "FashionBoxKey"
	}
	if e == GoodsSubType_Strong108Ticket {
		return "Strong108Ticket"
	}
	if e == GoodsSubType_Strong108Token {
		return "Strong108Token"
	}
	if e == GoodsSubType_FourSeasonsCoin {
		return "FourSeasonsCoin"
	}
	if e == GoodsSubType_ExchangeCoin {
		return "ExchangeCoin"
	}
	if e == GoodsSubType_CommonItem {
		return "CommonItem"
	}
	if e == GoodsSubType_RichProgress {
		return "RichProgress"
	}
	if e == GoodsSubType_GameProficiency {
		return "GameProficiency"
	}
	if e == GoodsSubType_NationalCoin {
		return "NationalCoin"
	}
	if e == GoodsSubType_TeamLeagueScore {
		return "TeamLeagueScore"
	}
	if e == GoodsSubType_SoulCard {
		return "SoulCard"
	}
	if e == GoodsSubType_DiamondCard {
		return "DiamondCard"
	}
	if e == GoodsSubType_ArcadeAllianceScore {
		return "ArcadeAllianceScore"
	}
	if e == GoodsSubType_Strongest108Score {
		return "Strongest108Score"
	}
	if e == GoodsSubType_MegaNurtureChip {
		return "MegaNurtureChip"
	}
	if e == GoodsSubType_RugalSignPass {
		return "RugalSignPass"
	}
	if e == GoodsSubType_SoulRecastTicket {
		return "SoulRecastTicket"
	}
	if e == GoodsSubType_ManorSuperRefreshCard {
		return "ManorSuperRefreshCard"
	}
	if e == GoodsSubType_ManorBattlePowerCard {
		return "ManorBattlePowerCard"
	}
	if e == GoodsSubType_ManorKickCard {
		return "ManorKickCard"
	}
	if e == GoodsSubType_ManorJumpInCard {
		return "ManorJumpInCard"
	}
	if e == GoodsSubType_ManorPoopCard {
		return "ManorPoopCard"
	}
	if e == GoodsSubType_ArcadeTicket {
		return "ArcadeTicket"
	}
	if e == GoodsSubType_TeamStageTicket {
		return "TeamStageTicket"
	}
	if e == GoodsSubType_TeamLeagueTaskCnt {
		return "TeamLeagueTaskCnt"
	}
	if e == GoodsSubType_PatrolSupportTimes {
		return "PatrolSupportTimes"
	}
	if e == GoodsSubType_PatrolDailyTimes {
		return "PatrolDailyTimes"
	}
	if e == GoodsSubType_ManorMeteor {
		return "ManorMeteor"
	}
	if e == GoodsSubType_CrazyFruitTicket {
		return "CrazyFruitTicket"
	}
	if e == GoodsSubType_SnowHighPass {
		return "SnowHighPass"
	}
	if e == GoodsSubType_TradeCoin {
		return "TradeCoin"
	}
	if e == GoodsSubType_DaSheBackSignPass {
		return "DaSheBackSignPass"
	}
	if e == GoodsSubType_FashionFragment {
		return "FashionFragment"
	}
	if e == GoodsSubType_SpringFestivalHighPass {
		return "SpringFestivalHighPass"
	}
	if e == GoodsSubType_SpringPassScore {
		return "SpringPassScore"
	}
	if e == GoodsSubType_Crazy8SignPass {
		return "Crazy8SignPass"
	}
	if e == GoodsSubType_RugalBackSignPass {
		return "RugalBackSignPass"
	}
	if e == GoodsSubType_CommonBattlePass {
		return "CommonBattlePass"
	}
	if e == GoodsSubType_ArcadeAllianceScoreItem {
		return "ArcadeAllianceScoreItem"
	}
	if e == GoodsSubType_TradeMedal {
		return "TradeMedal"
	}
	if e == GoodsSubType_TradeBigMedal {
		return "TradeBigMedal"
	}
	if e == GoodsSubType_KingWarWeekMedal {
		return "KingWarWeekMedal"
	}
	if e == GoodsSubType_KingWarFinalMedal {
		return "KingWarFinalMedal"
	}
	if e == GoodsSubType_MercenaryToken {
		return "MercenaryToken"
	}
	if e == GoodsSubType_PetSoul {
		return "PetSoul"
	}
	if e == GoodsSubType_PetSoulTicket {
		return "PetSoulTicket"
	}
	if e == GoodsSubType_PetEgg {
		return "PetEgg"
	}
	if e == GoodsSubType_FreeReviveTimes {
		return "FreeReviveTimes"
	}
	if e == GoodsSubType_ReviveTicket {
		return "ReviveTicket"
	}
	if e == GoodsSubType_KingWarSkin {
		return "KingWarSkin"
	}
	if e == GoodsSubType_AniAct {
		return "AniAct"
	}
	if e == GoodsSubType_FashionBoxBind {
		return "FashionBoxBind"
	}
	if e == GoodsSubType_Purifier {
		return "Purifier"
	}
	if e == GoodsSubType_Transformation {
		return "Transformation"
	}
	if e == GoodsSubType_HeroCard {
		return "HeroCard"
	}
	if e == GoodsSubType_HeroChip {
		return "HeroChip"
	}
	if e == GoodsSubType_HeroSoul {
		return "HeroSoul"
	}
	if e == GoodsSubType_HeroArtifact {
		return "HeroArtifact"
	}
	if e == GoodsSubType_HeroArtifactChip {
		return "HeroArtifactChip"
	}
	if e == GoodsSubType_HeroArtifactBreak {
		return "HeroArtifactBreak"
	}
	if e == GoodsSubType_HeroArtifactStone {
		return "HeroArtifactStone"
	}
	if e == GoodsSubType_HeroRandomSoul {
		return "HeroRandomSoul"
	}
	if e == GoodsSubType_HeroMagic {
		return "HeroMagic"
	}
	if e == GoodsSubType_HeroChipMagic {
		return "HeroChipMagic"
	}
	if e == GoodsSubType_HeroChipRandom {
		return "HeroChipRandom"
	}
	if e == GoodsSubType_HeroArtifactMagic {
		return "HeroArtifactMagic"
	}
	if e == GoodsSubType_HeroArtifactChipMagic {
		return "HeroArtifactChipMagic"
	}
	if e == GoodsSubType_HeroArtifactChipRandom {
		return "HeroArtifactChipRandom"
	}
	if e == GoodsSubType_HeroArtifactCompose {
		return "HeroArtifactCompose"
	}
	if e == GoodsSubType_HeroRandom {
		return "HeroRandom"
	}
	if e == GoodsSubType_HeroUniversalChip {
		return "HeroUniversalChip"
	}
	if e == GoodsSubType_HeroRarityStone {
		return "HeroRarityStone"
	}
	if e == GoodsSubType_ItemMagic {
		return "ItemMagic"
	}
	if e == GoodsSubType_ItemRandom {
		return "ItemRandom"
	}
	if e == GoodsSubType_NestChip {
		return "NestChip"
	}
	if e == GoodsSubType_NestRefineAttrLock {
		return "NestRefineAttrLock"
	}
	if e == GoodsSubType_HeroSoulRefining {
		return "HeroSoulRefining"
	}
	if e == GoodsSubType_DreamHeroRandom {
		return "DreamHeroRandom"
	}
	if e == GoodsSubType_HeroRights {
		return "HeroRights"
	}
	if e == GoodsSubType_NestRefineAttrRefine {
		return "NestRefineAttrRefine"
	}
	if e == GoodsSubType_NestStarup {
		return "NestStarup"
	}
	if e == GoodsSubType_HeroSwapStone {
		return "HeroSwapStone"
	}
	if e == GoodsSubType_BindHeroSoul {
		return "BindHeroSoul"
	}
	if e == GoodsSubType_BindHeroRandomSoul {
		return "BindHeroRandomSoul"
	}
	if e == GoodsSubType_ShopGood {
		return "ShopGood"
	}
	if e == GoodsSubType_SkinMap {
		return "SkinMap"
	}
	if e == GoodsSubType_SkinRole {
		return "SkinRole"
	}
	if e == GoodsSubType_LevelPositive {
		return "LevelPositive"
	}
	if e == GoodsSubType_LevelPassive {
		return "LevelPassive"
	}
	if e == GoodsSubType_FashionAvatar {
		return "FashionAvatar"
	}
	if e == GoodsSubType_FashionAvatarFrame {
		return "FashionAvatarFrame"
	}
	if e == GoodsSubType_FashionChatBubble {
		return "FashionChatBubble"
	}
	if e == GoodsSubType_FashionCutscene {
		return "FashionCutscene"
	}
	if e == GoodsSubType_FashionDamageEffects {
		return "FashionDamageEffects"
	}
	if e == GoodsSubType_FashionDynamicBackground {
		return "FashionDynamicBackground"
	}
	if e == GoodsSubType_FashionWing {
		return "FashionWing"
	}
	if e == GoodsSubType_FashionKillSpecialEffects {
		return "FashionKillSpecialEffects"
	}
	if e == GoodsSubType_FashionFinalEffects {
		return "FashionFinalEffects"
	}
	if e == GoodsSubType_FashionEntryAction {
		return "FashionEntryAction"
	}
	if e == GoodsSubType_FashionStandbyAction {
		return "FashionStandbyAction"
	}
	if e == GoodsSubType_FashionVerticalPainting {
		return "FashionVerticalPainting"
	}
	if e == GoodsSubType_FashionStandinPuppet {
		return "FashionStandinPuppet"
	}
	if e == GoodsSubType_FashionableDress {
		return "FashionableDress"
	}
	if e == GoodsSubType_FashionManorSkin {
		return "FashionManorSkin"
	}
	if e == GoodsSubType_FashionHomeSkin {
		return "FashionHomeSkin"
	}
	if e == GoodsSubType_FashionAoYi {
		return "FashionAoYi"
	}
	if e == GoodsSubType_FashionAidSkill {
		return "FashionAidSkill"
	}
	if e == GoodsSubType_FashionHeroPart1 {
		return "FashionHeroPart1"
	}
	if e == GoodsSubType_FashionHeroPart2 {
		return "FashionHeroPart2"
	}
	if e == GoodsSubType_FashionHeroPart3 {
		return "FashionHeroPart3"
	}
	if e == GoodsSubType_FashionHeroPart4 {
		return "FashionHeroPart4"
	}
	if e == GoodsSubType_Slime {
		return "Slime"
	}
	if e == GoodsSubType_GrayPet {
		return "GrayPet"
	}
	if e == GoodsSubType_GreenPet {
		return "GreenPet"
	}
	if e == GoodsSubType_BluePet {
		return "BluePet"
	}
	if e == GoodsSubType_PurplePet {
		return "PurplePet"
	}
	if e == GoodsSubType_OrangePet {
		return "OrangePet"
	}
	if e == GoodsSubType_GoldPet {
		return "GoldPet"
	}
	if e == GoodsSubType_RedPet {
		return "RedPet"
	}
	if e == GoodsSubType_FullColourPet {
		return "FullColourPet"
	}
	if e == GoodsSubType_RandomPet {
		return "RandomPet"
	}
	if e == GoodsSubType_PetMagic {
		return "PetMagic"
	}
	if e == GoodsSubType_SpiritDew {
		return "SpiritDew"
	}

	return ""

}

func GetGoodsSubType(tp string) (GoodsSubType, bool) {
	switch tp {
	case "Empty":
		return GoodsSubType_Empty, true
	case "Gold":
		return GoodsSubType_Gold, true
	case "Diamond":
		return GoodsSubType_Diamond, true
	case "Exp":
		return GoodsSubType_Exp, true
	case "HeroExp":
		return GoodsSubType_HeroExp, true
	case "HeroBadgePoint":
		return GoodsSubType_HeroBadgePoint, true
	case "TowerTicket":
		return GoodsSubType_TowerTicket, true
	case "ManorResTicket":
		return GoodsSubType_ManorResTicket, true
	case "MonthCardPoint":
		return GoodsSubType_MonthCardPoint, true
	case "ArenaTicket":
		return GoodsSubType_ArenaTicket, true
	case "EquipStone":
		return GoodsSubType_EquipStone, true
	case "EnjoymentCard":
		return GoodsSubType_EnjoymentCard, true
	case "VipRight":
		return GoodsSubType_VipRight, true
	case "CashCharge":
		return GoodsSubType_CashCharge, true
	case "GameCartridge":
		return GoodsSubType_GameCartridge, true
	case "WorldBossSweepTicket":
		return GoodsSubType_WorldBossSweepTicket, true
	case "MonthCard":
		return GoodsSubType_MonthCard, true
	case "LifetimeCard":
		return GoodsSubType_LifetimeCard, true
	case "HighLevelBattlePass":
		return GoodsSubType_HighLevelBattlePass, true
	case "AllianceDailyTaskProgress":
		return GoodsSubType_AllianceDailyTaskProgress, true
	case "AllianceExp":
		return GoodsSubType_AllianceExp, true
	case "VipExp":
		return GoodsSubType_VipExp, true
	case "PromotionExp":
		return GoodsSubType_PromotionExp, true
	case "HighGrowFund1BattlePass":
		return GoodsSubType_HighGrowFund1BattlePass, true
	case "HighGrowFund2BattlePass":
		return GoodsSubType_HighGrowFund2BattlePass, true
	case "SeniorRecruitTicket":
		return GoodsSubType_SeniorRecruitTicket, true
	case "HighChallengeFund1BattlePass":
		return GoodsSubType_HighChallengeFund1BattlePass, true
	case "HighChallengeFund2BattlePass":
		return GoodsSubType_HighChallengeFund2BattlePass, true
	case "HighHonorFund1BattlePass":
		return GoodsSubType_HighHonorFund1BattlePass, true
	case "HighHonorFund2BattlePass":
		return GoodsSubType_HighHonorFund2BattlePass, true
	case "DailyActiveScore":
		return GoodsSubType_DailyActiveScore, true
	case "WeeklyActiveScore":
		return GoodsSubType_WeeklyActiveScore, true
	case "Astrology":
		return GoodsSubType_Astrology, true
	case "BoxGift":
		return GoodsSubType_BoxGift, true
	case "AllianceTechFragment":
		return GoodsSubType_AllianceTechFragment, true
	case "AllianceTechTicket":
		return GoodsSubType_AllianceTechTicket, true
	case "HeroWeekCard":
		return GoodsSubType_HeroWeekCard, true
	case "FishWeekCard":
		return GoodsSubType_FishWeekCard, true
	case "TreasureWeekCard":
		return GoodsSubType_TreasureWeekCard, true
	case "ActivityPoint":
		return GoodsSubType_ActivityPoint, true
	case "AfkExpPackage":
		return GoodsSubType_AfkExpPackage, true
	case "AfkGoldPackage":
		return GoodsSubType_AfkGoldPackage, true
	case "WeekHeroRecruitPoint":
		return GoodsSubType_WeekHeroRecruitPoint, true
	case "WeekArtifactDrawPoint":
		return GoodsSubType_WeekArtifactDrawPoint, true
	case "SevenDayTaskScore":
		return GoodsSubType_SevenDayTaskScore, true
	case "Box":
		return GoodsSubType_Box, true
	case "FishEnergy":
		return GoodsSubType_FishEnergy, true
	case "FishSilver":
		return GoodsSubType_FishSilver, true
	case "FishPearl":
		return GoodsSubType_FishPearl, true
	case "FishBlackPearl":
		return GoodsSubType_FishBlackPearl, true
	case "EquipStoneBox":
		return GoodsSubType_EquipStoneBox, true
	case "FishOrderRefreshTicket":
		return GoodsSubType_FishOrderRefreshTicket, true
	case "FishOrderResetTicket":
		return GoodsSubType_FishOrderResetTicket, true
	case "ReduceTimeTicket":
		return GoodsSubType_ReduceTimeTicket, true
	case "FishBait":
		return GoodsSubType_FishBait, true
	case "FishGalaPoint":
		return GoodsSubType_FishGalaPoint, true
	case "BountyChallengeScore":
		return GoodsSubType_BountyChallengeScore, true
	case "SoulPotion":
		return GoodsSubType_SoulPotion, true
	case "AllianceCheerExp":
		return GoodsSubType_AllianceCheerExp, true
	case "CrossArenaTicket":
		return GoodsSubType_CrossArenaTicket, true
	case "CrossArenaMedal":
		return GoodsSubType_CrossArenaMedal, true
	case "Contribute":
		return GoodsSubType_Contribute, true
	case "ContributeChallengeTicket":
		return GoodsSubType_ContributeChallengeTicket, true
	case "AthenaSweepTicket":
		return GoodsSubType_AthenaSweepTicket, true
	case "PlaygroundTicket":
		return GoodsSubType_PlaygroundTicket, true
	case "RoomCard":
		return GoodsSubType_RoomCard, true
	case "HeroSkin":
		return GoodsSubType_HeroSkin, true
	case "Collection":
		return GoodsSubType_Collection, true
	case "Coin":
		return GoodsSubType_Coin, true
	case "LuckyStone":
		return GoodsSubType_LuckyStone, true
	case "JZYFish":
		return GoodsSubType_JZYFish, true
	case "ChooseSkinExperience":
		return GoodsSubType_ChooseSkinExperience, true
	case "Avatar":
		return GoodsSubType_Avatar, true
	case "AvatarFrame":
		return GoodsSubType_AvatarFrame, true
	case "Title":
		return GoodsSubType_Title, true
	case "FormationRefresh":
		return GoodsSubType_FormationRefresh, true
	case "ChatBuble":
		return GoodsSubType_ChatBuble, true
	case "NewbieHeroRecruitPoint":
		return GoodsSubType_NewbieHeroRecruitPoint, true
	case "AwakenStone":
		return GoodsSubType_AwakenStone, true
	case "ArtifactBattlePass":
		return GoodsSubType_ArtifactBattlePass, true
	case "CollectBattlePass":
		return GoodsSubType_CollectBattlePass, true
	case "NestBattlePass":
		return GoodsSubType_NestBattlePass, true
	case "RecruitBattlePass":
		return GoodsSubType_RecruitBattlePass, true
	case "WeekAllianceTechPoint":
		return GoodsSubType_WeekAllianceTechPoint, true
	case "Treasure":
		return GoodsSubType_Treasure, true
	case "LDCoin":
		return GoodsSubType_LDCoin, true
	case "LDSingleMatchCnt":
		return GoodsSubType_LDSingleMatchCnt, true
	case "LDFriendMatchCnt":
		return GoodsSubType_LDFriendMatchCnt, true
	case "LDInviteeMatchCnt":
		return GoodsSubType_LDInviteeMatchCnt, true
	case "AllianceStatueExp":
		return GoodsSubType_AllianceStatueExp, true
	case "StarCoin":
		return GoodsSubType_StarCoin, true
	case "ArtifactPassPoint":
		return GoodsSubType_ArtifactPassPoint, true
	case "FameMedal":
		return GoodsSubType_FameMedal, true
	case "ManorStone":
		return GoodsSubType_ManorStone, true
	case "ManorWood":
		return GoodsSubType_ManorWood, true
	case "ManorFood":
		return GoodsSubType_ManorFood, true
	case "AllianceTech":
		return GoodsSubType_AllianceTech, true
	case "BoxGiftExp":
		return GoodsSubType_BoxGiftExp, true
	case "AllianceDuelScore":
		return GoodsSubType_AllianceDuelScore, true
	case "BattlePassScore":
		return GoodsSubType_BattlePassScore, true
	case "WorldBossCoin":
		return GoodsSubType_WorldBossCoin, true
	case "FishCard":
		return GoodsSubType_FishCard, true
	case "UpRecruitTicket":
		return GoodsSubType_UpRecruitTicket, true
	case "WishRecruitTicket":
		return GoodsSubType_WishRecruitTicket, true
	case "AutoFishTicket":
		return GoodsSubType_AutoFishTicket, true
	case "ArtsTrialTimes":
		return GoodsSubType_ArtsTrialTimes, true
	case "BoxGiftBlessExp":
		return GoodsSubType_BoxGiftBlessExp, true
	case "NestChallengeTicket":
		return GoodsSubType_NestChallengeTicket, true
	case "NestChallengeKey":
		return GoodsSubType_NestChallengeKey, true
	case "NestSweepTicket":
		return GoodsSubType_NestSweepTicket, true
	case "AllianceBossAttackTimes":
		return GoodsSubType_AllianceBossAttackTimes, true
	case "SpSweepTicket":
		return GoodsSubType_SpSweepTicket, true
	case "WeekHeroRecruitCard":
		return GoodsSubType_WeekHeroRecruitCard, true
	case "WeekArtifactDrawCard":
		return GoodsSubType_WeekArtifactDrawCard, true
	case "WeekAllianceTechCard":
		return GoodsSubType_WeekAllianceTechCard, true
	case "RecruitExp":
		return GoodsSubType_RecruitExp, true
	case "FightingPlaceTicket":
		return GoodsSubType_FightingPlaceTicket, true
	case "TenDrawCard":
		return GoodsSubType_TenDrawCard, true
	case "Stamina":
		return GoodsSubType_Stamina, true
	case "StaminaMax":
		return GoodsSubType_StaminaMax, true
	case "FameForgeTicket":
		return GoodsSubType_FameForgeTicket, true
	case "TimesTicket":
		return GoodsSubType_TimesTicket, true
	case "FameGem":
		return GoodsSubType_FameGem, true
	case "DailyTopicTicket":
		return GoodsSubType_DailyTopicTicket, true
	case "DailyTopicScore":
		return GoodsSubType_DailyTopicScore, true
	case "FlipCardTicket":
		return GoodsSubType_FlipCardTicket, true
	case "FlipCardPackageTicket":
		return GoodsSubType_FlipCardPackageTicket, true
	case "CampBattleShopToken":
		return GoodsSubType_CampBattleShopToken, true
	case "FlipCardBattlePass":
		return GoodsSubType_FlipCardBattlePass, true
	case "FlipCardBattlePassPoint":
		return GoodsSubType_FlipCardBattlePassPoint, true
	case "FlipCardProgressScore":
		return GoodsSubType_FlipCardProgressScore, true
	case "FlipCardTurntableTicket":
		return GoodsSubType_FlipCardTurntableTicket, true
	case "DestinyFetterLevel":
		return GoodsSubType_DestinyFetterLevel, true
	case "MultipleProp":
		return GoodsSubType_MultipleProp, true
	case "DynamicHero":
		return GoodsSubType_DynamicHero, true
	case "UseTreasure":
		return GoodsSubType_UseTreasure, true
	case "OneYuanCard":
		return GoodsSubType_OneYuanCard, true
	case "UpPuzzle":
		return GoodsSubType_UpPuzzle, true
	case "PatrolRefreshTicket":
		return GoodsSubType_PatrolRefreshTicket, true
	case "RedEnvelope":
		return GoodsSubType_RedEnvelope, true
	case "BoxChipCompose":
		return GoodsSubType_BoxChipCompose, true
	case "VeteranSevenDayTaskScore":
		return GoodsSubType_VeteranSevenDayTaskScore, true
	case "VeteranHeroRecruitPoint":
		return GoodsSubType_VeteranHeroRecruitPoint, true
	case "HeroStarPass":
		return GoodsSubType_HeroStarPass, true
	case "SecretArtsPass":
		return GoodsSubType_SecretArtsPass, true
	case "HeroSoulPass":
		return GoodsSubType_HeroSoulPass, true
	case "HeroArtifactUniversalChip":
		return GoodsSubType_HeroArtifactUniversalChip, true
	case "UpSignPass":
		return GoodsSubType_UpSignPass, true
	case "WishSignPass":
		return GoodsSubType_WishSignPass, true
	case "BoxGiftCoin":
		return GoodsSubType_BoxGiftCoin, true
	case "AstrologyActivity":
		return GoodsSubType_AstrologyActivity, true
	case "FashionBox":
		return GoodsSubType_FashionBox, true
	case "FashionBoxKey":
		return GoodsSubType_FashionBoxKey, true
	case "Strong108Ticket":
		return GoodsSubType_Strong108Ticket, true
	case "Strong108Token":
		return GoodsSubType_Strong108Token, true
	case "FourSeasonsCoin":
		return GoodsSubType_FourSeasonsCoin, true
	case "ExchangeCoin":
		return GoodsSubType_ExchangeCoin, true
	case "CommonItem":
		return GoodsSubType_CommonItem, true
	case "RichProgress":
		return GoodsSubType_RichProgress, true
	case "GameProficiency":
		return GoodsSubType_GameProficiency, true
	case "NationalCoin":
		return GoodsSubType_NationalCoin, true
	case "TeamLeagueScore":
		return GoodsSubType_TeamLeagueScore, true
	case "SoulCard":
		return GoodsSubType_SoulCard, true
	case "DiamondCard":
		return GoodsSubType_DiamondCard, true
	case "ArcadeAllianceScore":
		return GoodsSubType_ArcadeAllianceScore, true
	case "Strongest108Score":
		return GoodsSubType_Strongest108Score, true
	case "MegaNurtureChip":
		return GoodsSubType_MegaNurtureChip, true
	case "RugalSignPass":
		return GoodsSubType_RugalSignPass, true
	case "SoulRecastTicket":
		return GoodsSubType_SoulRecastTicket, true
	case "ManorSuperRefreshCard":
		return GoodsSubType_ManorSuperRefreshCard, true
	case "ManorBattlePowerCard":
		return GoodsSubType_ManorBattlePowerCard, true
	case "ManorKickCard":
		return GoodsSubType_ManorKickCard, true
	case "ManorJumpInCard":
		return GoodsSubType_ManorJumpInCard, true
	case "ManorPoopCard":
		return GoodsSubType_ManorPoopCard, true
	case "ArcadeTicket":
		return GoodsSubType_ArcadeTicket, true
	case "TeamStageTicket":
		return GoodsSubType_TeamStageTicket, true
	case "TeamLeagueTaskCnt":
		return GoodsSubType_TeamLeagueTaskCnt, true
	case "PatrolSupportTimes":
		return GoodsSubType_PatrolSupportTimes, true
	case "PatrolDailyTimes":
		return GoodsSubType_PatrolDailyTimes, true
	case "ManorMeteor":
		return GoodsSubType_ManorMeteor, true
	case "CrazyFruitTicket":
		return GoodsSubType_CrazyFruitTicket, true
	case "SnowHighPass":
		return GoodsSubType_SnowHighPass, true
	case "TradeCoin":
		return GoodsSubType_TradeCoin, true
	case "DaSheBackSignPass":
		return GoodsSubType_DaSheBackSignPass, true
	case "FashionFragment":
		return GoodsSubType_FashionFragment, true
	case "SpringFestivalHighPass":
		return GoodsSubType_SpringFestivalHighPass, true
	case "SpringPassScore":
		return GoodsSubType_SpringPassScore, true
	case "Crazy8SignPass":
		return GoodsSubType_Crazy8SignPass, true
	case "RugalBackSignPass":
		return GoodsSubType_RugalBackSignPass, true
	case "CommonBattlePass":
		return GoodsSubType_CommonBattlePass, true
	case "ArcadeAllianceScoreItem":
		return GoodsSubType_ArcadeAllianceScoreItem, true
	case "TradeMedal":
		return GoodsSubType_TradeMedal, true
	case "TradeBigMedal":
		return GoodsSubType_TradeBigMedal, true
	case "KingWarWeekMedal":
		return GoodsSubType_KingWarWeekMedal, true
	case "KingWarFinalMedal":
		return GoodsSubType_KingWarFinalMedal, true
	case "MercenaryToken":
		return GoodsSubType_MercenaryToken, true
	case "PetSoul":
		return GoodsSubType_PetSoul, true
	case "PetSoulTicket":
		return GoodsSubType_PetSoulTicket, true
	case "PetEgg":
		return GoodsSubType_PetEgg, true
	case "FreeReviveTimes":
		return GoodsSubType_FreeReviveTimes, true
	case "ReviveTicket":
		return GoodsSubType_ReviveTicket, true
	case "KingWarSkin":
		return GoodsSubType_KingWarSkin, true
	case "AniAct":
		return GoodsSubType_AniAct, true
	case "FashionBoxBind":
		return GoodsSubType_FashionBoxBind, true
	case "Purifier":
		return GoodsSubType_Purifier, true
	case "Transformation":
		return GoodsSubType_Transformation, true
	case "HeroCard":
		return GoodsSubType_HeroCard, true
	case "HeroChip":
		return GoodsSubType_HeroChip, true
	case "HeroSoul":
		return GoodsSubType_HeroSoul, true
	case "HeroArtifact":
		return GoodsSubType_HeroArtifact, true
	case "HeroArtifactChip":
		return GoodsSubType_HeroArtifactChip, true
	case "HeroArtifactBreak":
		return GoodsSubType_HeroArtifactBreak, true
	case "HeroArtifactStone":
		return GoodsSubType_HeroArtifactStone, true
	case "HeroRandomSoul":
		return GoodsSubType_HeroRandomSoul, true
	case "HeroMagic":
		return GoodsSubType_HeroMagic, true
	case "HeroChipMagic":
		return GoodsSubType_HeroChipMagic, true
	case "HeroChipRandom":
		return GoodsSubType_HeroChipRandom, true
	case "HeroArtifactMagic":
		return GoodsSubType_HeroArtifactMagic, true
	case "HeroArtifactChipMagic":
		return GoodsSubType_HeroArtifactChipMagic, true
	case "HeroArtifactChipRandom":
		return GoodsSubType_HeroArtifactChipRandom, true
	case "HeroArtifactCompose":
		return GoodsSubType_HeroArtifactCompose, true
	case "HeroRandom":
		return GoodsSubType_HeroRandom, true
	case "HeroUniversalChip":
		return GoodsSubType_HeroUniversalChip, true
	case "HeroRarityStone":
		return GoodsSubType_HeroRarityStone, true
	case "ItemMagic":
		return GoodsSubType_ItemMagic, true
	case "ItemRandom":
		return GoodsSubType_ItemRandom, true
	case "NestChip":
		return GoodsSubType_NestChip, true
	case "NestRefineAttrLock":
		return GoodsSubType_NestRefineAttrLock, true
	case "HeroSoulRefining":
		return GoodsSubType_HeroSoulRefining, true
	case "DreamHeroRandom":
		return GoodsSubType_DreamHeroRandom, true
	case "HeroRights":
		return GoodsSubType_HeroRights, true
	case "NestRefineAttrRefine":
		return GoodsSubType_NestRefineAttrRefine, true
	case "NestStarup":
		return GoodsSubType_NestStarup, true
	case "HeroSwapStone":
		return GoodsSubType_HeroSwapStone, true
	case "BindHeroSoul":
		return GoodsSubType_BindHeroSoul, true
	case "BindHeroRandomSoul":
		return GoodsSubType_BindHeroRandomSoul, true
	case "ShopGood":
		return GoodsSubType_ShopGood, true
	case "SkinMap":
		return GoodsSubType_SkinMap, true
	case "SkinRole":
		return GoodsSubType_SkinRole, true
	case "LevelPositive":
		return GoodsSubType_LevelPositive, true
	case "LevelPassive":
		return GoodsSubType_LevelPassive, true
	case "FashionAvatar":
		return GoodsSubType_FashionAvatar, true
	case "FashionAvatarFrame":
		return GoodsSubType_FashionAvatarFrame, true
	case "FashionChatBubble":
		return GoodsSubType_FashionChatBubble, true
	case "FashionCutscene":
		return GoodsSubType_FashionCutscene, true
	case "FashionDamageEffects":
		return GoodsSubType_FashionDamageEffects, true
	case "FashionDynamicBackground":
		return GoodsSubType_FashionDynamicBackground, true
	case "FashionWing":
		return GoodsSubType_FashionWing, true
	case "FashionKillSpecialEffects":
		return GoodsSubType_FashionKillSpecialEffects, true
	case "FashionFinalEffects":
		return GoodsSubType_FashionFinalEffects, true
	case "FashionEntryAction":
		return GoodsSubType_FashionEntryAction, true
	case "FashionStandbyAction":
		return GoodsSubType_FashionStandbyAction, true
	case "FashionVerticalPainting":
		return GoodsSubType_FashionVerticalPainting, true
	case "FashionStandinPuppet":
		return GoodsSubType_FashionStandinPuppet, true
	case "FashionableDress":
		return GoodsSubType_FashionableDress, true
	case "FashionManorSkin":
		return GoodsSubType_FashionManorSkin, true
	case "FashionHomeSkin":
		return GoodsSubType_FashionHomeSkin, true
	case "FashionAoYi":
		return GoodsSubType_FashionAoYi, true
	case "FashionAidSkill":
		return GoodsSubType_FashionAidSkill, true
	case "FashionHeroPart1":
		return GoodsSubType_FashionHeroPart1, true
	case "FashionHeroPart2":
		return GoodsSubType_FashionHeroPart2, true
	case "FashionHeroPart3":
		return GoodsSubType_FashionHeroPart3, true
	case "FashionHeroPart4":
		return GoodsSubType_FashionHeroPart4, true
	case "Slime":
		return GoodsSubType_Slime, true
	case "GrayPet":
		return GoodsSubType_GrayPet, true
	case "GreenPet":
		return GoodsSubType_GreenPet, true
	case "BluePet":
		return GoodsSubType_BluePet, true
	case "PurplePet":
		return GoodsSubType_PurplePet, true
	case "OrangePet":
		return GoodsSubType_OrangePet, true
	case "GoldPet":
		return GoodsSubType_GoldPet, true
	case "RedPet":
		return GoodsSubType_RedPet, true
	case "FullColourPet":
		return GoodsSubType_FullColourPet, true
	case "RandomPet":
		return GoodsSubType_RandomPet, true
	case "PetMagic":
		return GoodsSubType_PetMagic, true
	case "SpiritDew":
		return GoodsSubType_SpiritDew, true

	}
	return -1, false

}

func GetGoodsSubType_ByAlias(tp string) (GoodsSubType, bool) {
	switch tp {
	case "空掉落":
		return GoodsSubType_Empty, true
	case "金币":
		return GoodsSubType_Gold, true
	case "钻石":
		return GoodsSubType_Diamond, true
	case "玩家经验":
		return GoodsSubType_Exp, true
	case "英雄经验":
		return GoodsSubType_HeroExp, true
	case "英雄图鉴点数":
		return GoodsSubType_HeroBadgePoint, true
	case "爬塔挑战券":
		return GoodsSubType_TowerTicket, true
	case "领地资源刷新券":
		return GoodsSubType_ManorResTicket, true
	case "月卡积分":
		return GoodsSubType_MonthCardPoint, true
	case "竞技场挑战券":
		return GoodsSubType_ArenaTicket, true
	case "装备强化石":
		return GoodsSubType_EquipStone, true
	case "畅玩卡":
		return GoodsSubType_EnjoymentCard, true
	case "Vip特权":
		return GoodsSubType_VipRight, true
	case "现金充值":
		return GoodsSubType_CashCharge, true
	case "游戏卡带":
		return GoodsSubType_GameCartridge, true
	case "世界boss扫荡券":
		return GoodsSubType_WorldBossSweepTicket, true
	case "月卡":
		return GoodsSubType_MonthCard, true
	case "终身卡":
		return GoodsSubType_LifetimeCard, true
	case "高级等级通行证":
		return GoodsSubType_HighLevelBattlePass, true
	case "联盟每日任务进度":
		return GoodsSubType_AllianceDailyTaskProgress, true
	case "联盟经验":
		return GoodsSubType_AllianceExp, true
	case "Vip经验":
		return GoodsSubType_VipExp, true
	case "佣兵团晋升经验":
		return GoodsSubType_PromotionExp, true
	case "高级成长基金1通行证":
		return GoodsSubType_HighGrowFund1BattlePass, true
	case "高级成长基金2通行证":
		return GoodsSubType_HighGrowFund2BattlePass, true
	case "高级招募券":
		return GoodsSubType_SeniorRecruitTicket, true
	case "高级挑战基金1通行证":
		return GoodsSubType_HighChallengeFund1BattlePass, true
	case "高级挑战基金2通行证":
		return GoodsSubType_HighChallengeFund2BattlePass, true
	case "高级荣誉基金1通行证":
		return GoodsSubType_HighHonorFund1BattlePass, true
	case "高级荣誉基金2通行证":
		return GoodsSubType_HighHonorFund2BattlePass, true
	case "日常活跃积分":
		return GoodsSubType_DailyActiveScore, true
	case "周常活跃积分":
		return GoodsSubType_WeeklyActiveScore, true
	case "神器占星":
		return GoodsSubType_Astrology, true
	case "格斗家收藏":
		return GoodsSubType_BoxGift, true
	case "公会科技碎片":
		return GoodsSubType_AllianceTechFragment, true
	case "公会科技能量":
		return GoodsSubType_AllianceTechTicket, true
	case "英雄周卡":
		return GoodsSubType_HeroWeekCard, true
	case "钓鱼周卡":
		return GoodsSubType_FishWeekCard, true
	case "宝物周卡":
		return GoodsSubType_TreasureWeekCard, true
	case "活动积分":
		return GoodsSubType_ActivityPoint, true
	case "离线英雄经验包":
		return GoodsSubType_AfkExpPackage, true
	case "离线金币包":
		return GoodsSubType_AfkGoldPackage, true
	case "招募周积分":
		return GoodsSubType_WeekHeroRecruitPoint, true
	case "神器周积分":
		return GoodsSubType_WeekArtifactDrawPoint, true
	case "七日任务积分":
		return GoodsSubType_SevenDayTaskScore, true
	case "格斗家收藏盒":
		return GoodsSubType_Box, true
	case "钓鱼体力":
		return GoodsSubType_FishEnergy, true
	case "钓鱼用银币":
		return GoodsSubType_FishSilver, true
	case "钓鱼用珍珠":
		return GoodsSubType_FishPearl, true
	case "钓鱼用黑珍珠":
		return GoodsSubType_FishBlackPearl, true
	case "装备强化石自选箱":
		return GoodsSubType_EquipStoneBox, true
	case "钓鱼订单刷新券":
		return GoodsSubType_FishOrderRefreshTicket, true
	case "钓鱼订单重置券":
		return GoodsSubType_FishOrderResetTicket, true
	case "加速券":
		return GoodsSubType_ReduceTimeTicket, true
	case "鱼饵":
		return GoodsSubType_FishBait, true
	case "鱼获宴会积分":
		return GoodsSubType_FishGalaPoint, true
	case "赏金值":
		return GoodsSubType_BountyChallengeScore, true
	case "战魂合成药剂":
		return GoodsSubType_SoulPotion, true
	case "跨服竞技助威经验":
		return GoodsSubType_AllianceCheerExp, true
	case "跨服竞技门票":
		return GoodsSubType_CrossArenaTicket, true
	case "跨服竞技勋章":
		return GoodsSubType_CrossArenaMedal, true
	case "贡献":
		return GoodsSubType_Contribute, true
	case "贡献挑战券":
		return GoodsSubType_ContributeChallengeTicket, true
	case "亚呆娜扫荡券":
		return GoodsSubType_AthenaSweepTicket, true
	case "乐园券":
		return GoodsSubType_PlaygroundTicket, true
	case "房卡":
		return GoodsSubType_RoomCard, true
	case "英雄皮肤":
		return GoodsSubType_HeroSkin, true
	case "藏品":
		return GoodsSubType_Collection, true
	case "硬币":
		return GoodsSubType_Coin, true
	case "幸运石":
		return GoodsSubType_LuckyStone, true
	case "姜子牙钓鱼":
		return GoodsSubType_JZYFish, true
	case "自选皮肤体验":
		return GoodsSubType_ChooseSkinExperience, true
	case "头像":
		return GoodsSubType_Avatar, true
	case "头像框":
		return GoodsSubType_AvatarFrame, true
	case "称号":
		return GoodsSubType_Title, true
	case "阵容刷新券":
		return GoodsSubType_FormationRefresh, true
	case "聊天气泡":
		return GoodsSubType_ChatBuble, true
	case "新手招募积分":
		return GoodsSubType_NewbieHeroRecruitPoint, true
	case "觉醒石":
		return GoodsSubType_AwakenStone, true
	case "神器基金通行证":
		return GoodsSubType_ArtifactBattlePass, true
	case "收藏基金通信证":
		return GoodsSubType_CollectBattlePass, true
	case "音巢基金通信证":
		return GoodsSubType_NestBattlePass, true
	case "招募基金通信证":
		return GoodsSubType_RecruitBattlePass, true
	case "公会科技周积分":
		return GoodsSubType_WeekAllianceTechPoint, true
	case "宝箱":
		return GoodsSubType_Treasure, true
	case "LD代币":
		return GoodsSubType_LDCoin, true
	case "LD单人匹配次数":
		return GoodsSubType_LDSingleMatchCnt, true
	case "LD好友匹配次数":
		return GoodsSubType_LDFriendMatchCnt, true
	case "LD受邀匹配次数":
		return GoodsSubType_LDInviteeMatchCnt, true
	case "公会雕像经验":
		return GoodsSubType_AllianceStatueExp, true
	case "星钻":
		return GoodsSubType_StarCoin, true
	case "神器基金积分":
		return GoodsSubType_ArtifactPassPoint, true
	case "勋章":
		return GoodsSubType_FameMedal, true
	case "石料":
		return GoodsSubType_ManorStone, true
	case "木料":
		return GoodsSubType_ManorWood, true
	case "食物":
		return GoodsSubType_ManorFood, true
	case "秘术":
		return GoodsSubType_AllianceTech, true
	case "收藏经验":
		return GoodsSubType_BoxGiftExp, true
	case "公会对决积分":
		return GoodsSubType_AllianceDuelScore, true
	case "通用战令积分":
		return GoodsSubType_BattlePassScore, true
	case "世界boss代币":
		return GoodsSubType_WorldBossCoin, true
	case "钓鱼权益卡":
		return GoodsSubType_FishCard, true
	case "Up招募券":
		return GoodsSubType_UpRecruitTicket, true
	case "心愿招募券":
		return GoodsSubType_WishRecruitTicket, true
	case "自动钓鱼券":
		return GoodsSubType_AutoFishTicket, true
	case "武道试炼次数":
		return GoodsSubType_ArtsTrialTimes, true
	case "祝福经验":
		return GoodsSubType_BoxGiftBlessExp, true
	case "音巢挑战券":
		return GoodsSubType_NestChallengeTicket, true
	case "音巢挑战密室钥匙":
		return GoodsSubType_NestChallengeKey, true
	case "音巢扫荡券":
		return GoodsSubType_NestSweepTicket, true
	case "联盟boss攻击次数":
		return GoodsSubType_AllianceBossAttackTimes, true
	case "南镇扫荡券":
		return GoodsSubType_SpSweepTicket, true
	case "招募冲刺特权卡":
		return GoodsSubType_WeekHeroRecruitCard, true
	case "神器冲刺特权卡":
		return GoodsSubType_WeekArtifactDrawCard, true
	case "秘术冲刺特权卡":
		return GoodsSubType_WeekAllianceTechCard, true
	case "招募经验":
		return GoodsSubType_RecruitExp, true
	case "格斗圣殿挑战券":
		return GoodsSubType_FightingPlaceTicket, true
	case "十连特权卡":
		return GoodsSubType_TenDrawCard, true
	case "体力道具":
		return GoodsSubType_Stamina, true
	case "体力上限道具":
		return GoodsSubType_StaminaMax, true
	case "名人堂锻造券":
		return GoodsSubType_FameForgeTicket, true
	case "限次门票":
		return GoodsSubType_TimesTicket, true
	case "名人堂宝石":
		return GoodsSubType_FameGem, true
	case "每日话题投票券":
		return GoodsSubType_DailyTopicTicket, true
	case "每日话题积分":
		return GoodsSubType_DailyTopicScore, true
	case "翻翻乐翻牌道具":
		return GoodsSubType_FlipCardTicket, true
	case "翻翻乐礼包抽奖券":
		return GoodsSubType_FlipCardPackageTicket, true
	case "阵营战商店代币":
		return GoodsSubType_CampBattleShopToken, true
	case "翻翻乐通行证":
		return GoodsSubType_FlipCardBattlePass, true
	case "翻翻乐通行证积分":
		return GoodsSubType_FlipCardBattlePassPoint, true
	case "翻翻乐进度积分":
		return GoodsSubType_FlipCardProgressScore, true
	case "翻翻乐转盘券":
		return GoodsSubType_FlipCardTurntableTicket, true
	case "宿命羁绊升级道具":
		return GoodsSubType_DestinyFetterLevel, true
	case "多倍道具":
		return GoodsSubType_MultipleProp, true
	case "动态立绘":
		return GoodsSubType_DynamicHero, true
	case "使用宝箱":
		return GoodsSubType_UseTreasure, true
	case "一元权益卡":
		return GoodsSubType_OneYuanCard, true
	case "Up养成拼图":
		return GoodsSubType_UpPuzzle, true
	case "探秘刷新券":
		return GoodsSubType_PatrolRefreshTicket, true
	case "红包":
		return GoodsSubType_RedEnvelope, true
	case "箱子碎片合成":
		return GoodsSubType_BoxChipCompose, true
	case "老手七日任务积分":
		return GoodsSubType_VeteranSevenDayTaskScore, true
	case "老手招募积分":
		return GoodsSubType_VeteranHeroRecruitPoint, true
	case "英雄星级通行证":
		return GoodsSubType_HeroStarPass, true
	case "秘术通行证":
		return GoodsSubType_SecretArtsPass, true
	case "八门通行证":
		return GoodsSubType_HeroSoulPass, true
	case "英雄神器万能碎片":
		return GoodsSubType_HeroArtifactUniversalChip, true
	case "Up招募签到战令":
		return GoodsSubType_UpSignPass, true
	case "自选招募签到战令":
		return GoodsSubType_WishSignPass, true
	case "格斗家礼物等级代币":
		return GoodsSubType_BoxGiftCoin, true
	case "占星限时活动":
		return GoodsSubType_AstrologyActivity, true
	case "时装箱":
		return GoodsSubType_FashionBox, true
	case "时装箱钥匙":
		return GoodsSubType_FashionBoxKey, true
	case "最强108挑战券":
		return GoodsSubType_Strong108Ticket, true
	case "最强108代币":
		return GoodsSubType_Strong108Token, true
	case "四季活动徽章":
		return GoodsSubType_FourSeasonsCoin, true
	case "兑换币":
		return GoodsSubType_ExchangeCoin, true
	case "普通道具":
		return GoodsSubType_CommonItem, true
	case "存钱罐进度":
		return GoodsSubType_RichProgress, true
	case "街机厅游戏熟练度":
		return GoodsSubType_GameProficiency, true
	case "国庆徽章":
		return GoodsSubType_NationalCoin, true
	case "公会团赛积分":
		return GoodsSubType_TeamLeagueScore, true
	case "战魂权益卡":
		return GoodsSubType_SoulCard, true
	case "钻石权益卡":
		return GoodsSubType_DiamondCard, true
	case "街机厅公会竞技值":
		return GoodsSubType_ArcadeAllianceScore, true
	case "最强108积分":
		return GoodsSubType_Strongest108Score, true
	case "Mega养成拼图":
		return GoodsSubType_MegaNurtureChip, true
	case "卢卡尔签到战令":
		return GoodsSubType_RugalSignPass, true
	case "战魂重铸券":
		return GoodsSubType_SoulRecastTicket, true
	case "超刷卡":
		return GoodsSubType_ManorSuperRefreshCard, true
	case "战力卡":
		return GoodsSubType_ManorBattlePowerCard, true
	case "踢飞卡":
		return GoodsSubType_ManorKickCard, true
	case "插队卡":
		return GoodsSubType_ManorJumpInCard, true
	case "便便卡":
		return GoodsSubType_ManorPoopCard, true
	case "街机厅挑战券":
		return GoodsSubType_ArcadeTicket, true
	case "组队关卡门票":
		return GoodsSubType_TeamStageTicket, true
	case "团赛任务次数":
		return GoodsSubType_TeamLeagueTaskCnt, true
	case "探秘协助次数":
		return GoodsSubType_PatrolSupportTimes, true
	case "探秘次数":
		return GoodsSubType_PatrolDailyTimes, true
	case "陨石":
		return GoodsSubType_ManorMeteor, true
	case "水果机抽奖券":
		return GoodsSubType_CrazyFruitTicket, true
	case "雪人高级战令":
		return GoodsSubType_SnowHighPass, true
	case "贸易币":
		return GoodsSubType_TradeCoin, true
	case "大蛇返场签到战令":
		return GoodsSubType_DaSheBackSignPass, true
	case "时装碎片":
		return GoodsSubType_FashionFragment, true
	case "春节高级战令":
		return GoodsSubType_SpringFestivalHighPass, true
	case "春节战令积分":
		return GoodsSubType_SpringPassScore, true
	case "疯八招募签到战令":
		return GoodsSubType_Crazy8SignPass, true
	case "卢卡尔返场签到战令":
		return GoodsSubType_RugalBackSignPass, true
	case "通用签到战令":
		return GoodsSubType_CommonBattlePass, true
	case "街机厅公会竞技值道具":
		return GoodsSubType_ArcadeAllianceScoreItem, true
	case "贸易银奖章":
		return GoodsSubType_TradeMedal, true
	case "贸易金奖章":
		return GoodsSubType_TradeBigMedal, true
	case "皇城周赛勋章":
		return GoodsSubType_KingWarWeekMedal, true
	case "皇城决赛勋章":
		return GoodsSubType_KingWarFinalMedal, true
	case "皇城佣兵令":
		return GoodsSubType_MercenaryToken, true
	case "宠物碎片":
		return GoodsSubType_PetSoul, true
	case "远古战魂券":
		return GoodsSubType_PetSoulTicket, true
	case "宠物蛋":
		return GoodsSubType_PetEgg, true
	case "皇城免费复活次数":
		return GoodsSubType_FreeReviveTimes, true
	case "皇城复活券":
		return GoodsSubType_ReviveTicket, true
	case "皇城皮肤":
		return GoodsSubType_KingWarSkin, true
	case "周年庆道具":
		return GoodsSubType_AniAct, true
	case "绑定时装箱":
		return GoodsSubType_FashionBoxBind, true
	case "抗扰剂":
		return GoodsSubType_Purifier, true
	case "装扮卡":
		return GoodsSubType_Transformation, true
	case "英雄卡牌":
		return GoodsSubType_HeroCard, true
	case "英雄碎片":
		return GoodsSubType_HeroChip, true
	case "英雄战魂":
		return GoodsSubType_HeroSoul, true
	case "英雄神器":
		return GoodsSubType_HeroArtifact, true
	case "英雄神器碎片":
		return GoodsSubType_HeroArtifactChip, true
	case "英雄神器突破石":
		return GoodsSubType_HeroArtifactBreak, true
	case "神器神格强化石":
		return GoodsSubType_HeroArtifactStone, true
	case "品质随机战魂":
		return GoodsSubType_HeroRandomSoul, true
	case "英雄自选":
		return GoodsSubType_HeroMagic, true
	case "英雄碎片自选":
		return GoodsSubType_HeroChipMagic, true
	case "英雄碎片随机":
		return GoodsSubType_HeroChipRandom, true
	case "英雄神器自选":
		return GoodsSubType_HeroArtifactMagic, true
	case "英雄神器碎片自选":
		return GoodsSubType_HeroArtifactChipMagic, true
	case "英雄神器碎片随机":
		return GoodsSubType_HeroArtifactChipRandom, true
	case "英雄神器合成":
		return GoodsSubType_HeroArtifactCompose, true
	case "英雄随机":
		return GoodsSubType_HeroRandom, true
	case "英雄万能碎片":
		return GoodsSubType_HeroUniversalChip, true
	case "英雄升品石":
		return GoodsSubType_HeroRarityStone, true
	case "道具自选":
		return GoodsSubType_ItemMagic, true
	case "道具随机":
		return GoodsSubType_ItemRandom, true
	case "音巢芯片":
		return GoodsSubType_NestChip, true
	case "音巢芯片属性锁定道具":
		return GoodsSubType_NestRefineAttrLock, true
	case "战魂洗练道具":
		return GoodsSubType_HeroSoulRefining, true
	case "许愿英雄随机":
		return GoodsSubType_DreamHeroRandom, true
	case "英雄权益":
		return GoodsSubType_HeroRights, true
	case "音巢芯片属性洗练道具":
		return GoodsSubType_NestRefineAttrRefine, true
	case "芯片升星消耗道具ID":
		return GoodsSubType_NestStarup, true
	case "英雄置换石":
		return GoodsSubType_HeroSwapStone, true
	case "绑定战魂":
		return GoodsSubType_BindHeroSoul, true
	case "绑定品质随机战魂":
		return GoodsSubType_BindHeroRandomSoul, true
	case "商品":
		return GoodsSubType_ShopGood, true
	case "地图":
		return GoodsSubType_SkinMap, true
	case "角色":
		return GoodsSubType_SkinRole, true
	case "主动道具":
		return GoodsSubType_LevelPositive, true
	case "被动道具":
		return GoodsSubType_LevelPassive, true
	case "时装头像":
		return GoodsSubType_FashionAvatar, true
	case "时装头像框":
		return GoodsSubType_FashionAvatarFrame, true
	case "时装聊天气泡":
		return GoodsSubType_FashionChatBubble, true
	case "PVP过场动画":
		return GoodsSubType_FashionCutscene, true
	case "伤害数字":
		return GoodsSubType_FashionDamageEffects, true
	case "个人信息动态背景":
		return GoodsSubType_FashionDynamicBackground, true
	case "翅膀":
		return GoodsSubType_FashionWing, true
	case "必杀特效":
		return GoodsSubType_FashionKillSpecialEffects, true
	case "终结特效":
		return GoodsSubType_FashionFinalEffects, true
	case "入场动作":
		return GoodsSubType_FashionEntryAction, true
	case "待机动作":
		return GoodsSubType_FashionStandbyAction, true
	case "时装立绘":
		return GoodsSubType_FashionVerticalPainting, true
	case "替身人偶":
		return GoodsSubType_FashionStandinPuppet, true
	case "时装":
		return GoodsSubType_FashionableDress, true
	case "领地皮肤":
		return GoodsSubType_FashionManorSkin, true
	case "主页皮肤":
		return GoodsSubType_FashionHomeSkin, true
	case "奥义技":
		return GoodsSubType_FashionAoYi, true
	case "援助技":
		return GoodsSubType_FashionAidSkill, true
	case "变身1":
		return GoodsSubType_FashionHeroPart1, true
	case "变身2":
		return GoodsSubType_FashionHeroPart2, true
	case "变身3":
		return GoodsSubType_FashionHeroPart3, true
	case "变身4":
		return GoodsSubType_FashionHeroPart4, true
	case "史莱姆":
		return GoodsSubType_Slime, true
	case "灰宠物":
		return GoodsSubType_GrayPet, true
	case "绿宠物":
		return GoodsSubType_GreenPet, true
	case "蓝宠物":
		return GoodsSubType_BluePet, true
	case "紫宠物":
		return GoodsSubType_PurplePet, true
	case "金宠物":
		return GoodsSubType_OrangePet, true
	case "橙宠物":
		return GoodsSubType_GoldPet, true
	case "红宠物":
		return GoodsSubType_RedPet, true
	case "彩宠物":
		return GoodsSubType_FullColourPet, true
	case "随机宠物":
		return GoodsSubType_RandomPet, true
	case "自选宠物":
		return GoodsSubType_PetMagic, true
	case "甘露":
		return GoodsSubType_SpiritDew, true

	}
	return -1, false

}

func (e Rarity) String() string {

	if e == Rarity_Gray {
		return "Gray"
	}
	if e == Rarity_Green {
		return "Green"
	}
	if e == Rarity_Blue {
		return "Blue"
	}
	if e == Rarity_Purple {
		return "Purple"
	}
	if e == Rarity_Orange {
		return "Orange"
	}
	if e == Rarity_Red {
		return "Red"
	}
	if e == Rarity_Rainbow {
		return "Rainbow"
	}
	if e == Rarity_Rainbow1 {
		return "Rainbow1"
	}
	if e == Rarity_Rainbow2 {
		return "Rainbow2"
	}
	if e == Rarity_Rainbow3 {
		return "Rainbow3"
	}
	if e == Rarity_Rainbow4 {
		return "Rainbow4"
	}
	if e == Rarity_Rainbow5 {
		return "Rainbow5"
	}
	if e == Rarity_Rainbow6 {
		return "Rainbow6"
	}
	if e == Rarity_Yellow {
		return "Yellow"
	}

	return ""

}

func GetRarity(tp string) (Rarity, bool) {
	switch tp {
	case "Gray":
		return Rarity_Gray, true
	case "Green":
		return Rarity_Green, true
	case "Blue":
		return Rarity_Blue, true
	case "Purple":
		return Rarity_Purple, true
	case "Orange":
		return Rarity_Orange, true
	case "Red":
		return Rarity_Red, true
	case "Rainbow":
		return Rarity_Rainbow, true
	case "Rainbow1":
		return Rarity_Rainbow1, true
	case "Rainbow2":
		return Rarity_Rainbow2, true
	case "Rainbow3":
		return Rarity_Rainbow3, true
	case "Rainbow4":
		return Rarity_Rainbow4, true
	case "Rainbow5":
		return Rarity_Rainbow5, true
	case "Rainbow6":
		return Rarity_Rainbow6, true
	case "Yellow":
		return Rarity_Yellow, true

	}
	return -1, false

}

func GetRarity_ByAlias(tp string) (Rarity, bool) {
	switch tp {
	case "灰":
		return Rarity_Gray, true
	case "绿":
		return Rarity_Green, true
	case "蓝":
		return Rarity_Blue, true
	case "紫":
		return Rarity_Purple, true
	case "橙":
		return Rarity_Orange, true
	case "红":
		return Rarity_Red, true
	case "彩":
		return Rarity_Rainbow, true
	case "彩1":
		return Rarity_Rainbow1, true
	case "彩2":
		return Rarity_Rainbow2, true
	case "彩3":
		return Rarity_Rainbow3, true
	case "彩4":
		return Rarity_Rainbow4, true
	case "彩5":
		return Rarity_Rainbow5, true
	case "彩6":
		return Rarity_Rainbow6, true
	case "黄":
		return Rarity_Yellow, true

	}
	return -1, false

}

func (e LDRarity) String() string {

	if e == LDRarity_Good {
		return "Good"
	}
	if e == LDRarity_Excellent {
		return "Excellent"
	}
	if e == LDRarity_Rare {
		return "Rare"
	}
	if e == LDRarity_Legend {
		return "Legend"
	}
	if e == LDRarity_Epic {
		return "Epic"
	}

	return ""

}

func GetLDRarity(tp string) (LDRarity, bool) {
	switch tp {
	case "Good":
		return LDRarity_Good, true
	case "Excellent":
		return LDRarity_Excellent, true
	case "Rare":
		return LDRarity_Rare, true
	case "Legend":
		return LDRarity_Legend, true
	case "Epic":
		return LDRarity_Epic, true

	}
	return -1, false

}

func GetLDRarity_ByAlias(tp string) (LDRarity, bool) {
	switch tp {
	case "优秀":
		return LDRarity_Good, true
	case "精良":
		return LDRarity_Excellent, true
	case "稀有":
		return LDRarity_Rare, true
	case "传奇":
		return LDRarity_Legend, true
	case "史诗":
		return LDRarity_Epic, true

	}
	return -1, false

}

func (e HeroSoulType) String() string {

	if e == HeroSoulType_Effect {
		return "Effect"
	}
	if e == HeroSoulType_Attack {
		return "Attack"
	}
	if e == HeroSoulType_Energy {
		return "Energy"
	}
	if e == HeroSoulType_Survive {
		return "Survive"
	}
	if e == HeroSoulType_Universal {
		return "Universal"
	}
	if e == HeroSoulType_NoLimit {
		return "NoLimit"
	}

	return ""

}

func GetHeroSoulType(tp string) (HeroSoulType, bool) {
	switch tp {
	case "Effect":
		return HeroSoulType_Effect, true
	case "Attack":
		return HeroSoulType_Attack, true
	case "Energy":
		return HeroSoulType_Energy, true
	case "Survive":
		return HeroSoulType_Survive, true
	case "Universal":
		return HeroSoulType_Universal, true
	case "NoLimit":
		return HeroSoulType_NoLimit, true

	}
	return -1, false

}

func GetHeroSoulType_ByAlias(tp string) (HeroSoulType, bool) {
	switch tp {
	case "天赋":
		return HeroSoulType_Effect, true
	case "攻击":
		return HeroSoulType_Attack, true
	case "能量":
		return HeroSoulType_Energy, true
	case "生命":
		return HeroSoulType_Survive, true
	case "万能":
		return HeroSoulType_Universal, true
	case "任意":
		return HeroSoulType_NoLimit, true

	}
	return -1, false

}

func (e HeroSoulGrade) String() string {

	if e == HeroSoulGrade_GradeD {
		return "GradeD"
	}
	if e == HeroSoulGrade_GradeC {
		return "GradeC"
	}
	if e == HeroSoulGrade_GradeB {
		return "GradeB"
	}
	if e == HeroSoulGrade_GradeA {
		return "GradeA"
	}
	if e == HeroSoulGrade_GradeS {
		return "GradeS"
	}

	return ""

}

func GetHeroSoulGrade(tp string) (HeroSoulGrade, bool) {
	switch tp {
	case "GradeD":
		return HeroSoulGrade_GradeD, true
	case "GradeC":
		return HeroSoulGrade_GradeC, true
	case "GradeB":
		return HeroSoulGrade_GradeB, true
	case "GradeA":
		return HeroSoulGrade_GradeA, true
	case "GradeS":
		return HeroSoulGrade_GradeS, true

	}
	return -1, false

}

func GetHeroSoulGrade_ByAlias(tp string) (HeroSoulGrade, bool) {
	switch tp {
	case "D":
		return HeroSoulGrade_GradeD, true
	case "C":
		return HeroSoulGrade_GradeC, true
	case "B":
		return HeroSoulGrade_GradeB, true
	case "A":
		return HeroSoulGrade_GradeA, true
	case "S":
		return HeroSoulGrade_GradeS, true

	}
	return -1, false

}

func (e HeroSoulEffectType) String() string {

	if e == HeroSoulEffectType_ComposeProbAdd {
		return "ComposeProbAdd"
	}
	if e == HeroSoulEffectType_BaseAttrBonus {
		return "BaseAttrBonus"
	}
	if e == HeroSoulEffectType_UniversalTag {
		return "UniversalTag"
	}
	if e == HeroSoulEffectType_ExtraSkill {
		return "ExtraSkill"
	}
	if e == HeroSoulEffectType_AsFoodProbAdd {
		return "AsFoodProbAdd"
	}
	if e == HeroSoulEffectType_UniversalTagAndBonus {
		return "UniversalTagAndBonus"
	}

	return ""

}

func GetHeroSoulEffectType(tp string) (HeroSoulEffectType, bool) {
	switch tp {
	case "ComposeProbAdd":
		return HeroSoulEffectType_ComposeProbAdd, true
	case "BaseAttrBonus":
		return HeroSoulEffectType_BaseAttrBonus, true
	case "UniversalTag":
		return HeroSoulEffectType_UniversalTag, true
	case "ExtraSkill":
		return HeroSoulEffectType_ExtraSkill, true
	case "AsFoodProbAdd":
		return HeroSoulEffectType_AsFoodProbAdd, true
	case "UniversalTagAndBonus":
		return HeroSoulEffectType_UniversalTagAndBonus, true

	}
	return -1, false

}

func GetHeroSoulEffectType_ByAlias(tp string) (HeroSoulEffectType, bool) {
	switch tp {
	case "合成概率增加":
		return HeroSoulEffectType_ComposeProbAdd, true
	case "基础属性加成":
		return HeroSoulEffectType_BaseAttrBonus, true
	case "万能词条":
		return HeroSoulEffectType_UniversalTag, true
	case "额外技能":
		return HeroSoulEffectType_ExtraSkill, true
	case "作为狗粮概率增加":
		return HeroSoulEffectType_AsFoodProbAdd, true
	case "万能词条及属性加成":
		return HeroSoulEffectType_UniversalTagAndBonus, true

	}
	return -1, false

}

func (e GoodsConfigId) String() string {

	if e == GoodsConfigId_Gold {
		return "Gold"
	}
	if e == GoodsConfigId_Diamond {
		return "Diamond"
	}
	if e == GoodsConfigId_Exp {
		return "Exp"
	}
	if e == GoodsConfigId_HeroExp {
		return "HeroExp"
	}
	if e == GoodsConfigId_HeroAtlasBadgePoint {
		return "HeroAtlasBadgePoint"
	}
	if e == GoodsConfigId_TowerTicket {
		return "TowerTicket"
	}
	if e == GoodsConfigId_ManorResTicket {
		return "ManorResTicket"
	}
	if e == GoodsConfigId_MonthCardPoint {
		return "MonthCardPoint"
	}
	if e == GoodsConfigId_ArenaTicket {
		return "ArenaTicket"
	}
	if e == GoodsConfigId_EnjoymentCard {
		return "EnjoymentCard"
	}
	if e == GoodsConfigId_CashCharge {
		return "CashCharge"
	}
	if e == GoodsConfigId_WorldBossSweepTicket {
		return "WorldBossSweepTicket"
	}
	if e == GoodsConfigId_HeroSoulRefiningExpend {
		return "HeroSoulRefiningExpend"
	}
	if e == GoodsConfigId_AfkPayExpend {
		return "AfkPayExpend"
	}
	if e == GoodsConfigId_VipExp {
		return "VipExp"
	}
	if e == GoodsConfigId_ArtifactBreakExpend {
		return "ArtifactBreakExpend"
	}
	if e == GoodsConfigId_AstrologyItem {
		return "AstrologyItem"
	}
	if e == GoodsConfigId_HolyCrystal {
		return "HolyCrystal"
	}
	if e == GoodsConfigId_AllianceTechUniversal {
		return "AllianceTechUniversal"
	}
	if e == GoodsConfigId_AllianceTechTicket {
		return "AllianceTechTicket"
	}
	if e == GoodsConfigId_BoxGiftExp {
		return "BoxGiftExp"
	}
	if e == GoodsConfigId_HeroTicket {
		return "HeroTicket"
	}
	if e == GoodsConfigId_ArtifactTicket {
		return "ArtifactTicket"
	}
	if e == GoodsConfigId_Box {
		return "Box"
	}
	if e == GoodsConfigId_FishSliver {
		return "FishSliver"
	}
	if e == GoodsConfigId_RedHeroMagicChip {
		return "RedHeroMagicChip"
	}
	if e == GoodsConfigId_ReduceTimeTicket {
		return "ReduceTimeTicket"
	}
	if e == GoodsConfigId_BountyChallengeScore {
		return "BountyChallengeScore"
	}
	if e == GoodsConfigId_CrossArenaTicket {
		return "CrossArenaTicket"
	}
	if e == GoodsConfigId_Contribution {
		return "Contribution"
	}
	if e == GoodsConfigId_ContributeChallengeTicket {
		return "ContributeChallengeTicket"
	}
	if e == GoodsConfigId_AthenaSweepTicket {
		return "AthenaSweepTicket"
	}
	if e == GoodsConfigId_PlaygroundTicket {
		return "PlaygroundTicket"
	}
	if e == GoodsConfigId_RoomCard {
		return "RoomCard"
	}
	if e == GoodsConfigId_Coin {
		return "Coin"
	}
	if e == GoodsConfigId_LuckyStone {
		return "LuckyStone"
	}
	if e == GoodsConfigId_RebornCard {
		return "RebornCard"
	}
	if e == GoodsConfigId_CleanCard {
		return "CleanCard"
	}
	if e == GoodsConfigId_RandomCard {
		return "RandomCard"
	}
	if e == GoodsConfigId_GoldCard {
		return "GoldCard"
	}
	if e == GoodsConfigId_RoomHeroCard {
		return "RoomHeroCard"
	}
	if e == GoodsConfigId_GridCard {
		return "GridCard"
	}
	if e == GoodsConfigId_CopyCard {
		return "CopyCard"
	}
	if e == GoodsConfigId_FundCard {
		return "FundCard"
	}
	if e == GoodsConfigId_LoongCard {
		return "LoongCard"
	}
	if e == GoodsConfigId_LegendCard {
		return "LegendCard"
	}
	if e == GoodsConfigId_GoldHammer {
		return "GoldHammer"
	}
	if e == GoodsConfigId_JZYGoneFish {
		return "JZYGoneFish"
	}
	if e == GoodsConfigId_JZYGoldFish {
		return "JZYGoldFish"
	}
	if e == GoodsConfigId_ChooseSkinExperience {
		return "ChooseSkinExperience"
	}
	if e == GoodsConfigId_FormationRefresh {
		return "FormationRefresh"
	}
	if e == GoodsConfigId_NestTicket {
		return "NestTicket"
	}
	if e == GoodsConfigId_AllianceBossExp {
		return "AllianceBossExp"
	}
	if e == GoodsConfigId_LDCoin {
		return "LDCoin"
	}
	if e == GoodsConfigId_LDSingleMatchCnt {
		return "LDSingleMatchCnt"
	}
	if e == GoodsConfigId_LDFriendMatchCnt {
		return "LDFriendMatchCnt"
	}
	if e == GoodsConfigId_LDInviteeMatchCnt {
		return "LDInviteeMatchCnt"
	}
	if e == GoodsConfigId_StarCoin {
		return "StarCoin"
	}
	if e == GoodsConfigId_InvalidItem {
		return "InvalidItem"
	}
	if e == GoodsConfigId_ArtifactStone {
		return "ArtifactStone"
	}
	if e == GoodsConfigId_AllianceFightScore {
		return "AllianceFightScore"
	}
	if e == GoodsConfigId_AllianceContribute {
		return "AllianceContribute"
	}
	if e == GoodsConfigId_SpTDCoin {
		return "SpTDCoin"
	}
	if e == GoodsConfigId_AutoFishTicket {
		return "AutoFishTicket"
	}
	if e == GoodsConfigId_BoxGiftBlessExp {
		return "BoxGiftBlessExp"
	}
	if e == GoodsConfigId_SpTDCurrency {
		return "SpTDCurrency"
	}
	if e == GoodsConfigId_SpTDSweep {
		return "SpTDSweep"
	}
	if e == GoodsConfigId_RecruitExp {
		return "RecruitExp"
	}
	if e == GoodsConfigId_FightingPlaceTicket {
		return "FightingPlaceTicket"
	}
	if e == GoodsConfigId_DailyTopicTicket {
		return "DailyTopicTicket"
	}
	if e == GoodsConfigId_DailyTopicScore {
		return "DailyTopicScore"
	}
	if e == GoodsConfigId_CampBattleShopToken {
		return "CampBattleShopToken"
	}
	if e == GoodsConfigId_SeasonsTicket {
		return "SeasonsTicket"
	}
	if e == GoodsConfigId_SeasonsExchangeCoin {
		return "SeasonsExchangeCoin"
	}
	if e == GoodsConfigId_SpringFestivalCoin {
		return "SpringFestivalCoin"
	}
	if e == GoodsConfigId_WinterFestivalCoin {
		return "WinterFestivalCoin"
	}
	if e == GoodsConfigId_SummerFestivalCoin {
		return "SummerFestivalCoin"
	}
	if e == GoodsConfigId_FlipCardBoxTicket {
		return "FlipCardBoxTicket"
	}
	if e == GoodsConfigId_FlipCardBattlePassPoint {
		return "FlipCardBattlePassPoint"
	}
	if e == GoodsConfigId_FlipFlipCardProgressPoint {
		return "FlipFlipCardProgressPoint"
	}
	if e == GoodsConfigId_FlipCardTurntableTicket {
		return "FlipCardTurntableTicket"
	}
	if e == GoodsConfigId_ArtifactMagicChip {
		return "ArtifactMagicChip"
	}
	if e == GoodsConfigId_UpSignBattlePassScore {
		return "UpSignBattlePassScore"
	}
	if e == GoodsConfigId_WishSignBattlePassScore {
		return "WishSignBattlePassScore"
	}
	if e == GoodsConfigId_AstrologyActivity {
		return "AstrologyActivity"
	}
	if e == GoodsConfigId_BoxGiftExtraCoin {
		return "BoxGiftExtraCoin"
	}
	if e == GoodsConfigId_UniversalHeroSoul {
		return "UniversalHeroSoul"
	}
	if e == GoodsConfigId_RecruitPrivilageCard {
		return "RecruitPrivilageCard"
	}
	if e == GoodsConfigId_ArtifactPrivilageCard {
		return "ArtifactPrivilageCard"
	}
	if e == GoodsConfigId_MiShuPrivilageCard {
		return "MiShuPrivilageCard"
	}
	if e == GoodsConfigId_TeamLeagueScore {
		return "TeamLeagueScore"
	}
	if e == GoodsConfigId_PromotionTicket {
		return "PromotionTicket"
	}
	if e == GoodsConfigId_SRTicket {
		return "SRTicket"
	}
	if e == GoodsConfigId_SSRTicket {
		return "SSRTicket"
	}
	if e == GoodsConfigId_URTicket {
		return "URTicket"
	}
	if e == GoodsConfigId_ManorSuperRefreshCard {
		return "ManorSuperRefreshCard"
	}
	if e == GoodsConfigId_ManorBattlePowerCard {
		return "ManorBattlePowerCard"
	}
	if e == GoodsConfigId_ManorJumpQueueCard {
		return "ManorJumpQueueCard"
	}
	if e == GoodsConfigId_ManorKickCard {
		return "ManorKickCard"
	}
	if e == GoodsConfigId_ManorRubbishCard {
		return "ManorRubbishCard"
	}
	if e == GoodsConfigId_TeamStageTicket {
		return "TeamStageTicket"
	}
	if e == GoodsConfigId_ManorMeteor1 {
		return "ManorMeteor1"
	}
	if e == GoodsConfigId_ManorMeteor2 {
		return "ManorMeteor2"
	}
	if e == GoodsConfigId_ManorMeteor3 {
		return "ManorMeteor3"
	}
	if e == GoodsConfigId_ManorMeteor4 {
		return "ManorMeteor4"
	}
	if e == GoodsConfigId_FashionChip {
		return "FashionChip"
	}
	if e == GoodsConfigId_SpringFestivalNestTicket {
		return "SpringFestivalNestTicket"
	}
	if e == GoodsConfigId_ReviveTicket {
		return "ReviveTicket"
	}
	if e == GoodsConfigId_Purifier {
		return "Purifier"
	}
	if e == GoodsConfigId_SpiritDew {
		return "SpiritDew"
	}

	return ""

}

func GetGoodsConfigId(tp string) (GoodsConfigId, bool) {
	switch tp {
	case "Gold":
		return GoodsConfigId_Gold, true
	case "Diamond":
		return GoodsConfigId_Diamond, true
	case "Exp":
		return GoodsConfigId_Exp, true
	case "HeroExp":
		return GoodsConfigId_HeroExp, true
	case "HeroAtlasBadgePoint":
		return GoodsConfigId_HeroAtlasBadgePoint, true
	case "TowerTicket":
		return GoodsConfigId_TowerTicket, true
	case "ManorResTicket":
		return GoodsConfigId_ManorResTicket, true
	case "MonthCardPoint":
		return GoodsConfigId_MonthCardPoint, true
	case "ArenaTicket":
		return GoodsConfigId_ArenaTicket, true
	case "EnjoymentCard":
		return GoodsConfigId_EnjoymentCard, true
	case "CashCharge":
		return GoodsConfigId_CashCharge, true
	case "WorldBossSweepTicket":
		return GoodsConfigId_WorldBossSweepTicket, true
	case "HeroSoulRefiningExpend":
		return GoodsConfigId_HeroSoulRefiningExpend, true
	case "AfkPayExpend":
		return GoodsConfigId_AfkPayExpend, true
	case "VipExp":
		return GoodsConfigId_VipExp, true
	case "ArtifactBreakExpend":
		return GoodsConfigId_ArtifactBreakExpend, true
	case "AstrologyItem":
		return GoodsConfigId_AstrologyItem, true
	case "HolyCrystal":
		return GoodsConfigId_HolyCrystal, true
	case "AllianceTechUniversal":
		return GoodsConfigId_AllianceTechUniversal, true
	case "AllianceTechTicket":
		return GoodsConfigId_AllianceTechTicket, true
	case "BoxGiftExp":
		return GoodsConfigId_BoxGiftExp, true
	case "HeroTicket":
		return GoodsConfigId_HeroTicket, true
	case "ArtifactTicket":
		return GoodsConfigId_ArtifactTicket, true
	case "Box":
		return GoodsConfigId_Box, true
	case "FishSliver":
		return GoodsConfigId_FishSliver, true
	case "RedHeroMagicChip":
		return GoodsConfigId_RedHeroMagicChip, true
	case "ReduceTimeTicket":
		return GoodsConfigId_ReduceTimeTicket, true
	case "BountyChallengeScore":
		return GoodsConfigId_BountyChallengeScore, true
	case "CrossArenaTicket":
		return GoodsConfigId_CrossArenaTicket, true
	case "Contribution":
		return GoodsConfigId_Contribution, true
	case "ContributeChallengeTicket":
		return GoodsConfigId_ContributeChallengeTicket, true
	case "AthenaSweepTicket":
		return GoodsConfigId_AthenaSweepTicket, true
	case "PlaygroundTicket":
		return GoodsConfigId_PlaygroundTicket, true
	case "RoomCard":
		return GoodsConfigId_RoomCard, true
	case "Coin":
		return GoodsConfigId_Coin, true
	case "LuckyStone":
		return GoodsConfigId_LuckyStone, true
	case "RebornCard":
		return GoodsConfigId_RebornCard, true
	case "CleanCard":
		return GoodsConfigId_CleanCard, true
	case "RandomCard":
		return GoodsConfigId_RandomCard, true
	case "GoldCard":
		return GoodsConfigId_GoldCard, true
	case "RoomHeroCard":
		return GoodsConfigId_RoomHeroCard, true
	case "GridCard":
		return GoodsConfigId_GridCard, true
	case "CopyCard":
		return GoodsConfigId_CopyCard, true
	case "FundCard":
		return GoodsConfigId_FundCard, true
	case "LoongCard":
		return GoodsConfigId_LoongCard, true
	case "LegendCard":
		return GoodsConfigId_LegendCard, true
	case "GoldHammer":
		return GoodsConfigId_GoldHammer, true
	case "JZYGoneFish":
		return GoodsConfigId_JZYGoneFish, true
	case "JZYGoldFish":
		return GoodsConfigId_JZYGoldFish, true
	case "ChooseSkinExperience":
		return GoodsConfigId_ChooseSkinExperience, true
	case "FormationRefresh":
		return GoodsConfigId_FormationRefresh, true
	case "NestTicket":
		return GoodsConfigId_NestTicket, true
	case "AllianceBossExp":
		return GoodsConfigId_AllianceBossExp, true
	case "LDCoin":
		return GoodsConfigId_LDCoin, true
	case "LDSingleMatchCnt":
		return GoodsConfigId_LDSingleMatchCnt, true
	case "LDFriendMatchCnt":
		return GoodsConfigId_LDFriendMatchCnt, true
	case "LDInviteeMatchCnt":
		return GoodsConfigId_LDInviteeMatchCnt, true
	case "StarCoin":
		return GoodsConfigId_StarCoin, true
	case "InvalidItem":
		return GoodsConfigId_InvalidItem, true
	case "ArtifactStone":
		return GoodsConfigId_ArtifactStone, true
	case "AllianceFightScore":
		return GoodsConfigId_AllianceFightScore, true
	case "AllianceContribute":
		return GoodsConfigId_AllianceContribute, true
	case "SpTDCoin":
		return GoodsConfigId_SpTDCoin, true
	case "AutoFishTicket":
		return GoodsConfigId_AutoFishTicket, true
	case "BoxGiftBlessExp":
		return GoodsConfigId_BoxGiftBlessExp, true
	case "SpTDCurrency":
		return GoodsConfigId_SpTDCurrency, true
	case "SpTDSweep":
		return GoodsConfigId_SpTDSweep, true
	case "RecruitExp":
		return GoodsConfigId_RecruitExp, true
	case "FightingPlaceTicket":
		return GoodsConfigId_FightingPlaceTicket, true
	case "DailyTopicTicket":
		return GoodsConfigId_DailyTopicTicket, true
	case "DailyTopicScore":
		return GoodsConfigId_DailyTopicScore, true
	case "CampBattleShopToken":
		return GoodsConfigId_CampBattleShopToken, true
	case "SeasonsTicket":
		return GoodsConfigId_SeasonsTicket, true
	case "SeasonsExchangeCoin":
		return GoodsConfigId_SeasonsExchangeCoin, true
	case "SpringFestivalCoin":
		return GoodsConfigId_SpringFestivalCoin, true
	case "WinterFestivalCoin":
		return GoodsConfigId_WinterFestivalCoin, true
	case "SummerFestivalCoin":
		return GoodsConfigId_SummerFestivalCoin, true
	case "FlipCardBoxTicket":
		return GoodsConfigId_FlipCardBoxTicket, true
	case "FlipCardBattlePassPoint":
		return GoodsConfigId_FlipCardBattlePassPoint, true
	case "FlipFlipCardProgressPoint":
		return GoodsConfigId_FlipFlipCardProgressPoint, true
	case "FlipCardTurntableTicket":
		return GoodsConfigId_FlipCardTurntableTicket, true
	case "ArtifactMagicChip":
		return GoodsConfigId_ArtifactMagicChip, true
	case "UpSignBattlePassScore":
		return GoodsConfigId_UpSignBattlePassScore, true
	case "WishSignBattlePassScore":
		return GoodsConfigId_WishSignBattlePassScore, true
	case "AstrologyActivity":
		return GoodsConfigId_AstrologyActivity, true
	case "BoxGiftExtraCoin":
		return GoodsConfigId_BoxGiftExtraCoin, true
	case "UniversalHeroSoul":
		return GoodsConfigId_UniversalHeroSoul, true
	case "RecruitPrivilageCard":
		return GoodsConfigId_RecruitPrivilageCard, true
	case "ArtifactPrivilageCard":
		return GoodsConfigId_ArtifactPrivilageCard, true
	case "MiShuPrivilageCard":
		return GoodsConfigId_MiShuPrivilageCard, true
	case "TeamLeagueScore":
		return GoodsConfigId_TeamLeagueScore, true
	case "PromotionTicket":
		return GoodsConfigId_PromotionTicket, true
	case "SRTicket":
		return GoodsConfigId_SRTicket, true
	case "SSRTicket":
		return GoodsConfigId_SSRTicket, true
	case "URTicket":
		return GoodsConfigId_URTicket, true
	case "ManorSuperRefreshCard":
		return GoodsConfigId_ManorSuperRefreshCard, true
	case "ManorBattlePowerCard":
		return GoodsConfigId_ManorBattlePowerCard, true
	case "ManorJumpQueueCard":
		return GoodsConfigId_ManorJumpQueueCard, true
	case "ManorKickCard":
		return GoodsConfigId_ManorKickCard, true
	case "ManorRubbishCard":
		return GoodsConfigId_ManorRubbishCard, true
	case "TeamStageTicket":
		return GoodsConfigId_TeamStageTicket, true
	case "ManorMeteor1":
		return GoodsConfigId_ManorMeteor1, true
	case "ManorMeteor2":
		return GoodsConfigId_ManorMeteor2, true
	case "ManorMeteor3":
		return GoodsConfigId_ManorMeteor3, true
	case "ManorMeteor4":
		return GoodsConfigId_ManorMeteor4, true
	case "FashionChip":
		return GoodsConfigId_FashionChip, true
	case "SpringFestivalNestTicket":
		return GoodsConfigId_SpringFestivalNestTicket, true
	case "ReviveTicket":
		return GoodsConfigId_ReviveTicket, true
	case "Purifier":
		return GoodsConfigId_Purifier, true
	case "SpiritDew":
		return GoodsConfigId_SpiritDew, true

	}
	return -1, false

}

func GetGoodsConfigId_ByAlias(tp string) (GoodsConfigId, bool) {
	switch tp {
	case "金币":
		return GoodsConfigId_Gold, true
	case "钻石":
		return GoodsConfigId_Diamond, true
	case "玩家经验":
		return GoodsConfigId_Exp, true
	case "英雄经验":
		return GoodsConfigId_HeroExp, true
	case "英雄图鉴点数":
		return GoodsConfigId_HeroAtlasBadgePoint, true
	case "爬塔挑战券":
		return GoodsConfigId_TowerTicket, true
	case "领地资源刷新券":
		return GoodsConfigId_ManorResTicket, true
	case "月卡积分":
		return GoodsConfigId_MonthCardPoint, true
	case "竞技场挑战券":
		return GoodsConfigId_ArenaTicket, true
	case "畅玩卡":
		return GoodsConfigId_EnjoymentCard, true
	case "现金充值":
		return GoodsConfigId_CashCharge, true
	case "世界boss扫荡券":
		return GoodsConfigId_WorldBossSweepTicket, true
	case "英雄战魂洗练消耗道具":
		return GoodsConfigId_HeroSoulRefiningExpend, true
	case "付费离线挂机消耗道具":
		return GoodsConfigId_AfkPayExpend, true
	case "VIP经验":
		return GoodsConfigId_VipExp, true
	case "神器突破消耗":
		return GoodsConfigId_ArtifactBreakExpend, true
	case "神器占星货币":
		return GoodsConfigId_AstrologyItem, true
	case "占星升档道具":
		return GoodsConfigId_HolyCrystal, true
	case "公会科技万能碎片":
		return GoodsConfigId_AllianceTechUniversal, true
	case "公会科技能量":
		return GoodsConfigId_AllianceTechTicket, true
	case "收藏宝箱经验":
		return GoodsConfigId_BoxGiftExp, true
	case "英雄招募券":
		return GoodsConfigId_HeroTicket, true
	case "占星之光":
		return GoodsConfigId_ArtifactTicket, true
	case "格斗家收藏盒":
		return GoodsConfigId_Box, true
	case "钓鱼银币":
		return GoodsConfigId_FishSliver, true
	case "英雄万能碎片":
		return GoodsConfigId_RedHeroMagicChip, true
	case "加速券":
		return GoodsConfigId_ReduceTimeTicket, true
	case "赏金值":
		return GoodsConfigId_BountyChallengeScore, true
	case "跨服竞技场门票":
		return GoodsConfigId_CrossArenaTicket, true
	case "贡献":
		return GoodsConfigId_Contribution, true
	case "贡献挑战入场券":
		return GoodsConfigId_ContributeChallengeTicket, true
	case "亚呆娜扫荡券":
		return GoodsConfigId_AthenaSweepTicket, true
	case "乐园券":
		return GoodsConfigId_PlaygroundTicket, true
	case "房卡":
		return GoodsConfigId_RoomCard, true
	case "硬币":
		return GoodsConfigId_Coin, true
	case "幸运石":
		return GoodsConfigId_LuckyStone, true
	case "复活卡":
		return GoodsConfigId_RebornCard, true
	case "清除卡":
		return GoodsConfigId_CleanCard, true
	case "随机卡":
		return GoodsConfigId_RandomCard, true
	case "金币卡":
		return GoodsConfigId_GoldCard, true
	case "英雄卡":
		return GoodsConfigId_RoomHeroCard, true
	case "格子卡":
		return GoodsConfigId_GridCard, true
	case "复制卡":
		return GoodsConfigId_CopyCard, true
	case "基金卡":
		return GoodsConfigId_FundCard, true
	case "龙血卡":
		return GoodsConfigId_LoongCard, true
	case "传奇卡":
		return GoodsConfigId_LegendCard, true
	case "黄金锤子":
		return GoodsConfigId_GoldHammer, true
	case "姜子牙跑鱼":
		return GoodsConfigId_JZYGoneFish, true
	case "姜子牙金鱼":
		return GoodsConfigId_JZYGoldFish, true
	case "自选皮肤体验":
		return GoodsConfigId_ChooseSkinExperience, true
	case "阵容刷新券":
		return GoodsConfigId_FormationRefresh, true
	case "音巢挑战门票":
		return GoodsConfigId_NestTicket, true
	case "公会雕像经验":
		return GoodsConfigId_AllianceBossExp, true
	case "LD玩法代币":
		return GoodsConfigId_LDCoin, true
	case "LD单人匹配次数":
		return GoodsConfigId_LDSingleMatchCnt, true
	case "LD好友匹配次数":
		return GoodsConfigId_LDFriendMatchCnt, true
	case "LD受邀匹配次数":
		return GoodsConfigId_LDInviteeMatchCnt, true
	case "星钻":
		return GoodsConfigId_StarCoin, true
	case "非法道具":
		return GoodsConfigId_InvalidItem, true
	case "神器源石":
		return GoodsConfigId_ArtifactStone, true
	case "公会对决积分":
		return GoodsConfigId_AllianceFightScore, true
	case "公会贡献道具":
		return GoodsConfigId_AllianceContribute, true
	case "南镇金币":
		return GoodsConfigId_SpTDCoin, true
	case "自动钓鱼券":
		return GoodsConfigId_AutoFishTicket, true
	case "祝福经验":
		return GoodsConfigId_BoxGiftBlessExp, true
	case "南镇货币":
		return GoodsConfigId_SpTDCurrency, true
	case "南镇扫荡券":
		return GoodsConfigId_SpTDSweep, true
	case "招募经验":
		return GoodsConfigId_RecruitExp, true
	case "格斗圣殿挑战券":
		return GoodsConfigId_FightingPlaceTicket, true
	case "每日话题投票券":
		return GoodsConfigId_DailyTopicTicket, true
	case "每日话题积分":
		return GoodsConfigId_DailyTopicScore, true
	case "阵营战商店代币":
		return GoodsConfigId_CampBattleShopToken, true
	case "四季活动心愿召唤券":
		return GoodsConfigId_SeasonsTicket, true
	case "庆典兑换币":
		return GoodsConfigId_SeasonsExchangeCoin, true
	case "春芽徽章":
		return GoodsConfigId_SpringFestivalCoin, true
	case "冬日徽章":
		return GoodsConfigId_WinterFestivalCoin, true
	case "夏荷徽章":
		return GoodsConfigId_SummerFestivalCoin, true
	case "翻翻乐-宝匣抽奖券":
		return GoodsConfigId_FlipCardBoxTicket, true
	case "翻翻乐战令积分":
		return GoodsConfigId_FlipCardBattlePassPoint, true
	case "翻翻乐进度积分":
		return GoodsConfigId_FlipFlipCardProgressPoint, true
	case "翻翻乐转盘券":
		return GoodsConfigId_FlipCardTurntableTicket, true
	case "神器万能碎片":
		return GoodsConfigId_ArtifactMagicChip, true
	case "Up招募签到战令积分":
		return GoodsConfigId_UpSignBattlePassScore, true
	case "自选招募签到战令积分":
		return GoodsConfigId_WishSignBattlePassScore, true
	case "神器祈福活动道具":
		return GoodsConfigId_AstrologyActivity, true
	case "格斗家礼物代币":
		return GoodsConfigId_BoxGiftExtraCoin, true
	case "万能本源":
		return GoodsConfigId_UniversalHeroSoul, true
	case "招募冲刺特权卡":
		return GoodsConfigId_RecruitPrivilageCard, true
	case "神器冲刺特权卡":
		return GoodsConfigId_ArtifactPrivilageCard, true
	case "秘术冲刺特权卡":
		return GoodsConfigId_MiShuPrivilageCard, true
	case "团赛积分":
		return GoodsConfigId_TeamLeagueScore, true
	case "升品凭证":
		return GoodsConfigId_PromotionTicket, true
	case "SR凭证":
		return GoodsConfigId_SRTicket, true
	case "SSR凭证":
		return GoodsConfigId_SSRTicket, true
	case "UR凭证":
		return GoodsConfigId_URTicket, true
	case "超刷卡":
		return GoodsConfigId_ManorSuperRefreshCard, true
	case "战力卡":
		return GoodsConfigId_ManorBattlePowerCard, true
	case "插队卡":
		return GoodsConfigId_ManorJumpQueueCard, true
	case "踢飞卡":
		return GoodsConfigId_ManorKickCard, true
	case "便便卡":
		return GoodsConfigId_ManorRubbishCard, true
	case "组队关卡门票":
		return GoodsConfigId_TeamStageTicket, true
	case "低级陨石":
		return GoodsConfigId_ManorMeteor1, true
	case "中级陨石":
		return GoodsConfigId_ManorMeteor2, true
	case "高级陨石":
		return GoodsConfigId_ManorMeteor3, true
	case "终级陨石":
		return GoodsConfigId_ManorMeteor4, true
	case "装扮碎片":
		return GoodsConfigId_FashionChip, true
	case "春节音巢入场券":
		return GoodsConfigId_SpringFestivalNestTicket, true
	case "皇城复活券":
		return GoodsConfigId_ReviveTicket, true
	case "抗扰剂":
		return GoodsConfigId_Purifier, true
	case "灵池甘露":
		return GoodsConfigId_SpiritDew, true

	}
	return -1, false

}

func (e RecruitType) String() string {

	if e == RecruitType_Free {
		return "Free"
	}
	if e == RecruitType_AD {
		return "AD"
	}
	if e == RecruitType_Senior {
		return "Senior"
	}
	if e == RecruitType_LimitTime {
		return "LimitTime"
	}
	if e == RecruitType_Super {
		return "Super"
	}
	if e == RecruitType_Up {
		return "Up"
	}
	if e == RecruitType_Wish {
		return "Wish"
	}
	if e == RecruitType_UpSuper {
		return "UpSuper"
	}
	if e == RecruitType_Seasons {
		return "Seasons"
	}
	if e == RecruitType_FundSuper {
		return "FundSuper"
	}
	if e == RecruitType_WishSuper {
		return "WishSuper"
	}

	return ""

}

func GetRecruitType(tp string) (RecruitType, bool) {
	switch tp {
	case "Free":
		return RecruitType_Free, true
	case "AD":
		return RecruitType_AD, true
	case "Senior":
		return RecruitType_Senior, true
	case "LimitTime":
		return RecruitType_LimitTime, true
	case "Super":
		return RecruitType_Super, true
	case "Up":
		return RecruitType_Up, true
	case "Wish":
		return RecruitType_Wish, true
	case "UpSuper":
		return RecruitType_UpSuper, true
	case "Seasons":
		return RecruitType_Seasons, true
	case "FundSuper":
		return RecruitType_FundSuper, true
	case "WishSuper":
		return RecruitType_WishSuper, true

	}
	return -1, false

}

func GetRecruitType_ByAlias(tp string) (RecruitType, bool) {
	switch tp {
	case "免费招募":
		return RecruitType_Free, true
	case "广告招募":
		return RecruitType_AD, true
	case "高级招募":
		return RecruitType_Senior, true
	case "限时招募":
		return RecruitType_LimitTime, true
	case "大保底":
		return RecruitType_Super, true
	case "Up招募":
		return RecruitType_Up, true
	case "心愿招募":
		return RecruitType_Wish, true
	case "Up超级招募":
		return RecruitType_UpSuper, true
	case "四季活动招募":
		return RecruitType_Seasons, true
	case "基金大保底":
		return RecruitType_FundSuper, true
	case "自选超级招募":
		return RecruitType_WishSuper, true

	}
	return -1, false

}
