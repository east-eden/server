package utils

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"
	"stathat.com/c/consistent"
)

var (
	ErrNodeNotFound = errors.New("node not found")
)

// select node by section id: game_id / 10
func SectionIDRandSelector(id int64) selector.SelectOption {
	gameID := MachineIDHigh(id)
	section := gameID / 10

	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		nodes := make([]*registry.Node, 0, len(srvs))

		for _, service := range srvs {
			for _, node := range service.Nodes {
				if node.Metadata["section"] == fmt.Sprintf("%d", section) {
					nodes = append(nodes, node)
				}
			}
		}

		return func() (*registry.Node, error) {
			if len(nodes) == 0 {
				return nil, fmt.Errorf("error selector")
			}

			i := rand.Intn(len(nodes))
			return nodes[i], nil
		}
	})
}

// select node by specific node_id
func SpecificIDSelector(nodeId string) selector.SelectOption {
	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		var node *registry.Node
		for _, service := range srvs {
			for _, nd := range service.Nodes {
				if nd.Id == nodeId {
					node = nd
				}
			}
		}

		return func() (*registry.Node, error) {
			if node == nil {
				return nil, ErrNodeNotFound
			}

			return node, nil
		}
	})
}

// select node by consistent hash
func ConsistentHashSelector(cons *consistent.Consistent, id string) selector.SelectOption {

	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		nodes := make(map[string]*registry.Node)
		var names []string
		for _, service := range srvs {
			for _, node := range service.Nodes {
				names = append(names, node.Id)
				nodes[node.Id] = node
			}
		}

		cons.Set(names)
		nodeName, err := cons.Get(id)
		return func() (*registry.Node, error) {
			if err != nil {
				return nil, fmt.Errorf("error selector with id:%s, err: %w", id, err)
			}

			node, ok := nodes[nodeName]
			if !ok {
				return nil, fmt.Errorf("error selector with id:%s, err: %w", id, err)
			}

			return node, nil
		}
	})
}
