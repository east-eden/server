package utils

import (
	"errors"
	"fmt"

	"github.com/asim/go-micro/v3/client/selector"
	"github.com/asim/go-micro/v3/registry"
	"stathat.com/c/consistent"
)

var (
	ErrNodeNotFound = errors.New("node not found")
)

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
