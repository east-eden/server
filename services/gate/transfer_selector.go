package gate

import (
	"strings"

	"bitbucket.org/funplus/gate/msg"
	"bitbucket.org/funplus/gate/selector"
	"bitbucket.org/funplus/kvs"
	"bitbucket.org/funplus/server/utils"
	"github.com/micro/go-micro/v2/registry"
	log "github.com/rs/zerolog/log"
	"stathat.com/c/consistent"
)

var ()

type TransferSelector struct {
	selector.Selector
	g          *Gate
	consistent *consistent.Consistent
}

func NewTransferSelector(g *Gate) *TransferSelector {
	ts := &TransferSelector{
		g:          g,
		consistent: consistent.New(),
	}

	ts.consistent.NumberOfReplicas = maxGameNode

	return ts
}

func (s *TransferSelector) Select(service string, opts ...selector.SelectOption) (*kvs.Entry, error) {
	so := selector.NewSelectOptions(opts...)

	es := kvs.Entries{}
	for _, filter := range so.Filters {
		es = filter(es)
	}

	return so.Strategy(es)
}

func (s *TransferSelector) Mark(service string, node *kvs.Entry, err error) {

}

func (s *TransferSelector) Reset(service string) {

}

func (s *TransferSelector) Close() error {
	return nil
}

func (s *TransferSelector) String() string {
	return ""
}

func (s *TransferSelector) ConsistentHashFilter(hs *msg.Handshake) selector.Filter {
	srvs, err := s.g.mi.srv.Server().Options().Registry.GetService("game")
	if !utils.ErrCheck(err, "GetService failed when TransferSelector.ConsistentHashFilter", hs) {
		srvs = make([]*registry.Service, 0)
	}

	nodes := make(map[string]*registry.Node)
	var names []string
	for _, service := range srvs {
		for _, node := range service.Nodes {
			names = append(names, node.Id)
			nodes[node.Id] = node
		}
	}

	s.consistent.Set(names)
	nodeName, err := s.consistent.Get(hs.UserID)
	_ = utils.ErrCheck(err, "consistent Get failed when TransferSelector.ConsistentHashFilter", hs)

	node, ok := nodes[nodeName]
	if !ok {
		log.Error().Str("node_name", nodeName).Str("user_id", hs.UserID).Msg("consistent get node failed")
	}

	return func(kvs.Entries) kvs.Entries {
		if !ok {
			return kvs.Entries{}
		}

		addr := strings.Split(node.Metadata["publicTcpAddr"], ":")
		return kvs.Entries{
			kvs.NewEntry(
				addr[0],
				addr[1],
				kvs.WithEntryIdentifier(nodeName),
			),
		}
	}
}
