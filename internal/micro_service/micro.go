package micro_service

import (
	"github.com/east-eden/server/utils"
	"github.com/spf13/cast"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
)

// 通过服务名找到所有服务列表，例如：serviceName = gate、game、lobby
func GetServicesByName(c client.Client, serviceName string) ([]*registry.Service, error) {
	services, err := c.Options().Registry.GetService(serviceName)
	if !utils.ErrCheck(err, "get service failed", serviceName) {
		return nil, err
	}

	return services, nil
}

// 通过唯一节点id找到服务节点，例如：serviceName = game, nodeId = 101
func GetNodeById(c client.Client, serviceName string, nodeId int32) (*registry.Node, error) {
	next, err := c.Options().Selector.Select(serviceName, NodeIdSelector(nodeId))
	if !utils.ErrCheck(err, "select node failed", serviceName, nodeId) {
		return nil, err
	}

	node, err := next()
	if !utils.ErrCheck(err, "next node failed", serviceName, nodeId) {
		return nil, err
	}

	return node, nil
}

// 通过唯一节点名找到服务节点，例如：serviceName = game, nodeName = game-101
func GetNodeByName(c client.Client, serviceName string, nodeName string) (*registry.Node, error) {
	next, err := c.Options().Selector.Select(serviceName, NodeNameSelector(nodeName))
	if !utils.ErrCheck(err, "select node failed", serviceName, nodeName) {
		return nil, err
	}

	node, err := next()
	if !utils.ErrCheck(err, "next node failed", serviceName, nodeName) {
		return nil, err
	}

	return node, nil
}

// 通过一致性哈希key找到服务节点，例如：serviceName = game, hashKey = user11010003
func GetNodeByHashkey(c client.Client, serviceName string, hashKey string) (*registry.Node, error) {
	next, err := c.Options().Selector.Select(serviceName, ConsistentHashSelector(hashKey))
	if !utils.ErrCheck(err, "select node failed", serviceName, hashKey) {
		return nil, err
	}

	node, err := next()
	if !utils.ErrCheck(err, "next node failed", serviceName, hashKey) {
		return nil, err
	}

	return node, nil
}

// 通过哈希key列表返回节点列表
func GetNodesByHashkeyList(c client.Client, serviceName string, hashKeyList []int64) map[string][]int64 {
	result := make(map[string][]int64)
	if len(hashKeyList) == 0 {
		return result
	}

	services, err := c.Options().Registry.GetService(serviceName)
	if !utils.ErrCheck(err, "get service failed", serviceName) {
		return result
	}

	entry, err := getConsistentRingEntry(services)
	if !utils.ErrCheck(err, "get consistent ring entry failed", serviceName) {
		return result
	}

	hashKeyMap := make(map[int64]struct{}, len(hashKeyList))
	for _, hashKey := range hashKeyList {
		hashKeyMap[hashKey] = struct{}{}
	}

	for hashKey := range hashKeyMap {
		node, err := entry.selectNode(cast.ToString(hashKey))
		if !utils.ErrCheck(err, "select node failed", serviceName, hashKey) {
			continue
		}

		if _, found := result[node.Id]; !found {
			result[node.Id] = make([]int64, 0, 4)
		}

		result[node.Id] = append(result[node.Id], hashKey)
	}

	return result
}
