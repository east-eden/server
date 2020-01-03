package utils

import (
	"fmt"
	"math/rand"

	"github.com/micro/go-micro/client/selector"
	"github.com/micro/go-micro/registry"
)

// select node by section id: game_id / 10
func PlayerInfoAvgSelector(id int64) selector.SelectOption {
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
