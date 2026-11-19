package worldboss

import (
	"github.com/rs/zerolog/log"
	"time"

	"github.com/east-eden/server/excel/config"
	"github.com/east-eden/server/utils"
)

func GetWorldBossId() (int32, time.Time) {
	bossCfg := config.GetBossConfigEntry()

	firstBossTime, err := time.ParseInLocation(utils.Layout, bossCfg.FirstBossTime, utils.NowInLocation().Location())
	if err != nil {
		log.Error().Msgf("parse boss time error: %v %s", err, bossCfg.FirstBossTime)
		return 1, utils.LongLongAfter()
	}
	nowTime := utils.NowInLocation()
	diff := nowTime.Sub(firstBossTime)
	period := int32(diff.Seconds() / float64(bossCfg.BossDuration)) // 转换为期数
	// 获取下一期boss起始时间
	nextBossTime := firstBossTime.Add(time.Second * time.Duration(period+1) * time.Duration(bossCfg.BossDuration))
	groupId := bossCfg.GroupRotation[period%int32(len(bossCfg.GroupRotation))]
	return groupId, nextBossTime
}
