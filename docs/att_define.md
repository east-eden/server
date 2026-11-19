# 属性定义

属性分三组，分别对应 `define/att_define.go` 里的三组常量与三个名字表
（`AttNames`、`AttBaseNames`、`AttPercentNames`）：

- **最终属性** `Att_Begin..Att_End`（0-45，共 `AttFinalNum` 个）：战斗结算实际使用的属性
- **基础属性** `Att_Base_Begin..Att_Base_End`（0-5，共 `AttBaseNum` 个）：攻击/护甲等基础值
- **百分比属性** `Att_Percent_Begin..Att_Percent_End`（0-5，共 `AttPercentNum` 个）：按基础属性做百分比加成

`internal/att` 的 `AttManager` 会按「基础属性 → 百分比属性 → 最终属性」的顺序聚合，再供战斗使用。

## 最终属性（AttFinalNum）

| 枚举 | 值 | 说明 |
| --- | --- | --- |
| `Att_Atk` | 0 | 攻击力最终值 |
| `Att_Armor` | 1 | 护甲最终值 |
| `Att_SelfDmgInc` | 2 | 自身伤害加成 |
| `Att_EnemyWoundInc` | 3 | 敌方受伤加成 |
| `Att_SelfDmgDec` | 4 | 自身伤害消减 |
| `Att_EnemyWoundDec` | 5 | 敌方受伤消减 |
| `Att_Crit` | 6 | 暴击值 |
| `Att_CritInc` | 7 | 暴击倍数加层 |
| `Att_Tenacity` | 8 | 韧性 |
| `Att_Heal` | 9 | 治疗强度最终值 |
| `Att_RealDmg` | 10 | 真实伤害 |
| `Att_MoveSpeed` | 11 | 战场移动速度 |
| `Att_AtbSpeed` | 12 | 时间槽速度 |
| `Att_EffectHit` | 13 | 技能效果命中 |
| `Att_EffectResist` | 14 | 技能效果抵抗 |
| `Att_MaxHP` | 15 | 生命值上限 |
| `Att_CurHP` | 16 | 当前生命值 |
| `Att_MaxMP` | 17 | 蓝量上限 |
| `Att_CurMP` | 18 | 当前蓝量 |
| `Att_GenMP` | 19 | MP 恢复值 |
| `Att_MaxRage` | 20 | 怒气值上限 |
| `Att_GenRagePercent` | 21 | 怒气增长提高百分比 |
| `Att_InitRage` | 22 | 初始怒气 |
| `Att_Hit` | 23 | 命中 |
| `Att_Dodge` | 24 | 闪避 |
| `Att_MoveScope` | 25 | 移动范围 |
| `Att_DmgPhysics` ~ `Att_DmgDeath` | 26-35 | 各系伤害加成（物理/地/水/火/风/雷/时/空/钢/灭） |
| `Att_ResPhysics` ~ `Att_ResDeath` | 36-45 | 各系伤害抗性（物理/地/水/火/风/雷/时/空/钢/灭） |

## 基础属性（AttBaseNum）

| 枚举 | 值 | 说明 |
| --- | --- | --- |
| `Att_AtkBase` | 0 | 攻击力基础值 |
| `Att_ArmorBase` | 1 | 护甲基础值 |
| `Att_HealBase` | 2 | 治疗强度基础值 |
| `Att_MoveSpeedBase` | 3 | 战场移动速度基础值 |
| `Att_AtbSpeedBase` | 4 | 时间槽速度基础值 |
| `Att_MaxHPBase` | 5 | 生命值上限基础值 |

## 百分比属性（AttPercentNum）

| 枚举 | 值 | 说明 |
| --- | --- | --- |
| `Att_AtkPercent` | 0 | 攻击力百分比 |
| `Att_ArmorPercent` | 1 | 护甲百分比 |
| `Att_HealPercent` | 2 | 治疗强度百分比 |
| `Att_MoveSpeedPercent` | 3 | 战场移动速度百分比 |
| `Att_AtbSpeedPercent` | 4 | 时间槽速度百分比 |
| `Att_MaxHPPercent` | 5 | 生命值上限百分比 |

> 英雄卡牌的属性不直接使用这里的枚举，而是来自 `BattleProperty` 配置表的公式，
> 详见 `services/game/hero/card.go`。
