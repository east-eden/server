package config

import "sort"

// 活动相关的按 activityId 分组索引。
// 生成代码只提供主键查询，活动任务/签到奖励都需要按活动聚合。
var (
	activityTasksByActivity   map[string][]*ActivityTaskEntry
	activitySignInsByActivity map[string][]*ActivitySignInEntry
)

func buildActivityIndex() {
	if activityTasksByActivity != nil {
		return
	}

	tasks := make(map[string][]*ActivityTaskEntry, 64)
	for _, e := range activityTaskEntries.Arrays {
		tasks[e.ActivityId] = append(tasks[e.ActivityId], e)
	}

	for _, list := range tasks {
		sort.Slice(list, func(i, j int) bool {
			if list[i].Order != list[j].Order {
				return list[i].Order < list[j].Order
			}

			return list[i].Id < list[j].Id
		})
	}

	signIns := make(map[string][]*ActivitySignInEntry, 16)
	for _, e := range activitySignInEntries.Arrays {
		signIns[e.ActivityId] = append(signIns[e.ActivityId], e)
	}

	for _, list := range signIns {
		sort.Slice(list, func(i, j int) bool { return list[i].Day < list[j].Day })
	}

	activityTasksByActivity = tasks
	activitySignInsByActivity = signIns
}

// GetActivityTasksByActivityId 返回某活动下的任务，按 Order 升序
func GetActivityTasksByActivityId(activityId string) []*ActivityTaskEntry {
	buildActivityIndex()
	return activityTasksByActivity[activityId]
}

// GetActivitySignInByActivityId 返回某活动的签到奖励，按 Day 升序
func GetActivitySignInByActivityId(activityId string) []*ActivitySignInEntry {
	buildActivityIndex()
	return activitySignInsByActivity[activityId]
}

// ResetActivityIndex 配置重载后重建索引
func ResetActivityIndex() {
	activityTasksByActivity = nil
	activitySignInsByActivity = nil
}
