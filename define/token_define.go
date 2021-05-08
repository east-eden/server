package define

const (
	Token_Begin              int32 = iota
	Token_Gold               int32 = iota - 1 // 0 金币
	Token_Diamond                             // 1 钻石
	Token_Music                               // 2 音石
	Token_Friendship                          // 3 友情点
	Token_Maze                                // 4 迷宫代币
	Token_Arena                               // 5 竞技场代币
	Token_Expedition                          // 6 远征代币
	Token_Home                                // 7 家园代币
	Token_ActivityQuest                       // 8 活跃任务代币
	Token_GuildContrubution                   // 9 帮会贡献
	Token_CrystalExp                          // 10 晶石经验代币
	Token_ExploreReputation1                  // 11 探索声望代币1
	Token_ExploreReputation2                  // 12 探索声望代币2
	Token_ExploreReputation3                  // 13 探索声望代币3
	Token_ExploreReputation4                  // 14 探索声望代币4
	Token_ExploreReputation5                  // 15 探索声望代币5
	Token_Strength                            // 16 体力
	Token_StrengthStore                       // 17 体力存储

	Token_CollectionBegin
	Token_CollectionGreen  int32 = iota - 2 // 18 绿色收集品通用碎片
	Token_CollectionBlue                    // 19 蓝色收集品通用碎片
	Token_CollectionPurple                  // 20 紫色收集品通用碎片
	Token_CollectionYellow                  // 21 黄色收集品通用碎片
	Token_CollectionEnd

	Token_End int32 = iota - 3
)
