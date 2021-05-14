package define

const (
	Condition_Type_Begin int32 = iota
	Condition_Type_And   int32 = iota - 1 // 需满足所有子条件
	Condition_Type_Or                     // 满足其中一个子条件

	Condition_Type_End
)

const (
	Condition_SubType_Begin                    int32 = iota
	Condition_SubType_TeamLevel_Achieve        int32 = iota - 1 // 0 队伍等级达到**级
	Condition_SubType_KillAllEnemy                              // 1 击杀所有敌方单位
	Condition_SubType_OurUnitAllDead                            // 2 己方单位全部死亡
	Condition_SubType_OurUnitDeadLessThan                       // 3 己方单位死亡人数小于*
	Condition_SubType_InterruptEnemySkill                       // 4 成功打断*次敌方技能
	Condition_SubType_OurUnitCastUltimateSkill                  // 5 己方单位成功使用*次奥义技能
	Condition_SubType_CombatPassTimeLessThan                    // 6 通关时间小于*秒
	Condition_SubType_KillEnemyTypeIdFirst                      // 7 优先击杀id为*的敌方单位

	Condition_SubType_End
)
