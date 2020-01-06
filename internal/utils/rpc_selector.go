package utils

import (
	"fmt"
	"math/rand"

	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
)

// select node by section id: game_id / 10
func SectionIDRandSelector(id int64) selector.SelectOption {
	gameId := MachineIDHigh(id)
	section := gameId / 10

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

// select node by specific game_id
func SpecificIDSelector(id int64) selector.SelectOption {
	gameId := MachineIDHigh(id)

	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		var node *registry.Node
		for _, service := range srvs {
			for _, nd := range service.Nodes {
				if nd.Id == fmt.Sprintf("%s-%d", service.Name, gameId) {
					node = nd
				}
			}
		}

		return func() (*registry.Node, error) {
			if node == nil {
				return nil, fmt.Errorf("error selector")
			}

			return node, nil
		}
	})
}
