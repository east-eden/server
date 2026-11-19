package activity

import (
	"time"

	"github.com/east-eden/server/utils"
	"github.com/spf13/cast"
)

// ActivityConf 活动时间配置。
//
// kof 里这些字段来自 ts 协议结构 NetActivityConf，这里改为等价的本地结构，
// 由业务侧从配置表或协议填充后传入，不依赖具体序列化格式。
type ActivityConf struct {
	StartMode              string `json:"startMode"`              // 开始模式：cron/datetime/createtime/servercreate/special/forever/playerLevel
	Start                  string `json:"start"`                  // 开始参数
	EndMode                string `json:"endMode"`                // 结束模式：cron/datetime/duration/special/forever
	End                    string `json:"end"`                    // 结束参数
	CreatetimeSpecialParam string `json:"createtimeSpecialParam"` // 创建时间特殊参数
}

// 通过时间格式获取活动时间
func GetActivityTime(avyConf *ActivityConf, fromWhen time.Time, createTime time.Time, playerLevel int32) (startTm, endTm time.Time) {
	// cron模式下需要获取配置的上一个时间和下一个时间
	var cronNextStartTm, cronNextEndTm, cronPrevStartTm, cronPrevEndTm time.Time

	// 获取活动开始时间
	switch avyConf.StartMode {
	case "cron":
		cron, err := utils.ParseCron(avyConf.Start)
		if !utils.ErrCheck(err, "parse cron failed", avyConf) {
			cronNextStartTm = utils.EnsureLocation(time.Time{})
			cronPrevStartTm = utils.EnsureLocation(time.Time{})
			break
		}
		cronNextStartTm = utils.EnsureLocation(cron.Next(fromWhen))
		cronPrevStartTm = utils.EnsureLocation(cron.Prev(fromWhen))

	case "datetime":
		tm, err := utils.ParseDateTime(avyConf.Start)
		if !utils.ErrCheck(err, "parse datetime failed", avyConf) {
			startTm = utils.EnsureLocation(time.Time{})
			break
		}
		startTm = tm

	case "createtime":
		startTm = time.Date(createTime.Year(), createTime.Month(), createTime.Day(), 0, 0, 0, 0, createTime.Location())
		// 创建时间 + 创建之后开启时间
		if avyConf.Start != "" {
			dur, err := time.ParseDuration(avyConf.Start)
			if !utils.ErrCheck(err, "parse duration failed", avyConf) {
				endTm = utils.EnsureLocation(time.Time{})
				break
			}
			startTm = startTm.Add(dur)
		}
		// 特殊参数
		if avyConf.CreatetimeSpecialParam != "" {
			specialStartTm, err := utils.ParseDateTime(avyConf.CreatetimeSpecialParam)
			if !utils.ErrCheck(err, "parse datetime failed", avyConf) {
				break
			}
			// 如果配置时间在创建时间之后则使用配置时间
			if specialStartTm.After(startTm) {
				startTm = specialStartTm
			}
		}
	case "servercreate":
		startTm = utils.EnsureLocation(time.Time{})

	case "special":
		startTm = utils.EnsureLocation(time.Time{})

	case "forever":
		startTm = utils.FromTimestamp(0)
	case "playerLevel":
		if playerLevel >= cast.ToInt32(avyConf.Start) {
			startTm = time.Date(createTime.Year(), createTime.Month(), createTime.Day(), 0, 0, 0, 0, createTime.Location())
		} else {
			startTm, _ = time.ParseInLocation(time.DateTime, "2038-01-18 23:59:59", utils.NowInLocation().Location())
		}
	default:
		startTm = utils.NowInLocation()
	}

	// 获取活动结束时间
	switch avyConf.EndMode {
	case "cron":
		cron, err := utils.ParseCron(avyConf.End)
		if !utils.ErrCheck(err, "parse cron failed", avyConf) {
			cronNextEndTm = utils.EnsureLocation(time.Time{})
			cronPrevEndTm = utils.EnsureLocation(time.Time{})
			break
		}
		cronNextEndTm = utils.EnsureLocation(cron.Next(cronNextStartTm))
		cronPrevEndTm = utils.EnsureLocation(cron.Next(cronPrevStartTm))

	case "datetime":
		tm, err := utils.ParseDateTime(avyConf.End)
		if !utils.ErrCheck(err, "parse datetime failed", avyConf) {
			endTm = utils.EnsureLocation(time.Time{})
			break
		}
		endTm = tm

	case "duration":
		dur, err := time.ParseDuration(avyConf.End)
		if !utils.ErrCheck(err, "parse duration failed", avyConf) {
			endTm = utils.EnsureLocation(time.Time{})
			break
		}
		endTm = startTm.Add(dur)
		cronNextEndTm = cronNextStartTm.Add(dur)
		cronPrevEndTm = cronPrevStartTm.Add(dur)

	case "special":
		endTm, _ = time.ParseInLocation(time.DateTime, "2038-01-18 23:59:59", utils.NowInLocation().Location())

	case "forever":
		endTm, _ = time.ParseInLocation(time.DateTime, "2038-01-18 23:59:59", utils.NowInLocation().Location())
	default:
		endTm = utils.EnsureLocation(time.Time{})
	}

	// cron模式单独处理 需判断cron上一次匹配的时间和下一次匹配的时间
	if avyConf.StartMode == "cron" {
		if fromWhen.After(cronPrevEndTm) || fromWhen.Equal(cronPrevEndTm) {
			startTm = cronNextStartTm
			endTm = cronNextEndTm
		} else {
			startTm = cronPrevStartTm
			endTm = cronPrevEndTm
		}
	}

	return
}
