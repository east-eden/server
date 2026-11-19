package micro_service

import (
	"errors"
	"fmt"

	"github.com/spf13/cast"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
)

var (
	ErrNodeNotFound = errors.New("node not found")
)

// select node by full node name
func NodeNameSelector(nodeName string) selector.SelectOption {
	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		var node *registry.Node

	loop:
		for _, service := range srvs {
			for _, nd := range service.Nodes {
				if nd.Metadata["stale"] == "true" {
					continue
				}

				if nd.Id == nodeName {
					node = nd
					break loop
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

// select node by int node id
func NodeIdSelector(nodeId int32) selector.SelectOption {
	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		var node *registry.Node

	loop:
		for _, service := range srvs {
			for _, nd := range service.Nodes {
				if nd.Metadata["stale"] == "true" {
					continue
				}

				if nd.Metadata["nodeId"] == cast.ToString(nodeId) {
					node = nd
					break loop
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
func ConsistentHashSelector(id string) selector.SelectOption {
	return selector.WithStrategy(func(srvs []*registry.Service) selector.Next {
		entry, err := getConsistentRingEntry(srvs)
		if err != nil {
			return func() (*registry.Node, error) {
				return nil, fmt.Errorf("consistent hash select id:%s, err: %w", id, err)
			}
		}

		node, err := entry.selectNode(id)
		return func() (*registry.Node, error) {
			return node, err
		}
	})
}
