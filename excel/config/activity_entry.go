package config

import (
	"sort"

	"github.com/east-eden/server/excel"
	"github.com/east-eden/server/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

var activityEntries *ActivityEntries //Activity.json全局变量

// 属性表 活动表-拳皇/Activity
type ActivityEntry struct {
	Id                     string   `json:"id,omitempty" bson:"id"`                                         //##活动id
	Group                  string   `json:"group,omitempty" bson:"group"`                                   //活动组
	Type                   string   `json:"type,omitempty" bson:"type"`                                     //活动类型
	BattePassRound         int32    `json:"battePassRound,omitempty" bson:"battePassRound"`                 //battpass轮数
	AbClass                string   `json:"abClass,omitempty" bson:"abClass"`                               //AB活动
	AbSettingId            int32    `json:"abSettingId,omitempty" bson:"abSettingId"`                       //AB分组id
	AbTag                  string   `json:"abTag,omitempty" bson:"abTag"`                                   //AB配置tag
	IsStandalone           bool     `json:"isStandalone,omitempty" bson:"isStandalone"`                     //是否独立活动
	IsMainInitTask         bool     `json:"isMainInitTask,omitempty" bson:"isMainInitTask"`                 //是否在主任务初始化全部任务
	PrevActivityId         string   `json:"prevActivityId,omitempty" bson:"prevActivityId"`                 //前置活动id
	SubActivities          []string `json:"subActivities,omitempty" bson:"subActivities"`                   //子活动
	SubArrangement         string   `json:"subArrangement,omitempty" bson:"subArrangement"`                 //子活动编排类型
	ShopId                 int32    `json:"shopId,omitempty" bson:"shopId"`                                 //挂载商店id
	ShopIds                []int32  `json:"shopIds,omitempty" bson:"shopIds"`                               //挂载多个商店ids
	CounterTask            []int32  `json:"counterTask,omitempty" bson:"counterTask"`                       //计数任务索引ActivityCounterTask表
	Tasks                  []int32  `json:"tasks,omitempty" bson:"tasks"`                                   //任务列表索引ActivityTask表
	TaskDayRefresh         bool     `json:"taskDayRefresh,omitempty" bson:"taskDayRefresh"`                 //是否任务每日刷新
	PointTasks             []int32  `json:"pointTasks,omitempty" bson:"pointTasks"`                         //积分任务索引ActivityPointTask
	SignInDays             int32    `json:"signInDays,omitempty" bson:"signInDays"`                         //签到持续天数
	ShopRefresh            bool     `json:"shopRefresh,omitempty" bson:"shopRefresh"`                       //活动重置时商店是否一并重置
	ServerOpenDays         int32    `json:"serverOpenDays,omitempty" bson:"serverOpenDays"`                 //服务器开启第x天检测活动（开服当天算第一天）
	CreateDaysAtStart      int32    `json:"createDaysAtStart,omitempty" bson:"createDaysAtStart"`           //活动开启时服务器需满足的开启天数(开服当天算第一天)
	StartMode              string   `json:"startMode,omitempty" bson:"startMode"`                           //开启时间模式
	CreatetimeSpecialParam string   `json:"createtimeSpecialParam,omitempty" bson:"createtimeSpecialParam"` //createtime或者serveropen模式参数(创建时间+开启duration和配置时间取更大的)
	Start                  string   `json:"start,omitempty" bson:"start"`                                   //开启时间
	EndMode                string   `json:"endMode,omitempty" bson:"endMode"`                               //结束时间模式
	End                    string   `json:"end,omitempty" bson:"end"`                                       //结束时间Xm/h
	IsAllianceActivity     bool     `json:"isAllianceActivity,omitempty" bson:"isAllianceActivity"`         //是否是公会活动
}

// Activity.json属性表集合
type ActivityEntries struct {
	Rows   map[string]*ActivityEntry `json:"rows,omitempty" bson:"rows"`     //
	Arrays []*ActivityEntry          `json:"arrays,omitempty" bson:"arrays"` //
}

func init() {
	excel.AddEntryLoader("Activity.json", (*ActivityEntries)(nil))
}

func (e *ActivityEntries) Load(excelFileRaw *excel.ExcelFileRaw) (any, error) {

	activityEntries = &ActivityEntries{
		Rows:   make(map[string]*ActivityEntry),
		Arrays: make([]*ActivityEntry, 0, 8),
	}

	for _, v := range excelFileRaw.RowData {
		entry := &ActivityEntry{}
		if len(v) == 0 {
			continue
		}

		err := mapstructure.Decode(v, entry)
		if !utils.ErrCheck(err, "decode excel data to struct failed", v) {
			return activityEntries, err
		}

		activityEntries.Rows[entry.Id] = entry
		activityEntries.Arrays = append(activityEntries.Arrays, entry)
	}

	sort.Slice(activityEntries.Arrays, func(i, j int) bool {
		return activityEntries.Arrays[i].Id < activityEntries.Arrays[j].Id
	})

	excel.AddArrayGetter("activity.json", activityEntries)
	log.Info().Str("excel_file", excelFileRaw.Filename).Msg("excel load success")
	return activityEntries, nil

}

func (e *ActivityEntries) GetArrays() any {
	return activityEntries.Arrays
}

func GetActivityEntry(id string) (*ActivityEntry, bool) {
	entry, ok := activityEntries.Rows[id]
	return entry, ok
}

func GetActivitySize() int {
	return int(len(activityEntries.Rows))
}

func GetActivityRows() map[string]*ActivityEntry {
	return activityEntries.Rows
}

func GetActivityArrays() []*ActivityEntry {
	return activityEntries.Arrays
}
