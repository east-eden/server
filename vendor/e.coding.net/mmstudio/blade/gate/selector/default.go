package selector

import (
	"e.coding.net/mmstudio/blade/golib/container/cmap"
	"e.coding.net/mmstudio/blade/kvs"
	"github.com/rs/zerolog/log"
	"sync"
)

type serviceGroup struct {
	sync.RWMutex
	registry     kvs.RegistryV2
	name         string
	entries      kvs.Entries
	mark         cmap.ConcurrentMap
	registryPath func(serviceKey string) string
}

func (g *serviceGroup) init(stop chan struct{}) error {
	g.Lock()
	defer g.Unlock()
	if g.name == "" {
		log.Warn().Str("serviceName", g.name).Msg("service group init ignore")
		return nil
	}
	es, err := g.registry.GetAll(g.registryPath(g.name), kvs.WithReadOptionConsistent(true))
	if err != nil {
		log.Error().Err(err).Str("service", g.name).Msg("get service fail")
	} else {
		g.entries = es
	}
	go func() {
		esCh, errCh := g.registry.Watch(g.registryPath(g.name))
		for {
			select {
			case es := <-esCh:
				if len(es) == 0 {
					log.Warn().Str("service", g.name).Msg("no entry")
				}
				g.Lock()
				g.entries = es
				g.Unlock()
			case err = <-errCh:
				log.Error().Err(err).Str("service", g.name).Msg("watch service recv fail")
			case <-stop:
				return
			}
		}
	}()

	return nil
}

func (g *serviceGroup) getEntries() kvs.Entries {
	g.RLock()
	defer g.RUnlock()
	return g.entries
}

type registrySelector struct {
	services ConcurrentMapStringService
	op       *Options
	stop     chan struct{}
}

func NewSelector(opts ...Option) Selector {
	ret := &registrySelector{
		services: NewConcurrentMapStringService(),
		op:       NewOptions(opts...),
		stop:     make(chan struct{}),
	}
	return ret
}

func (r *registrySelector) Select(service string, opts ...SelectOption) (*kvs.Entry, error) {
	es, err := r.getEntries(service)
	if err != nil {
		return nil, err
	}

	so := NewSelectOptions(opts...)

	for _, filter := range so.Filters {
		es = filter(es)
	}

	return so.Strategy(es)
}

func (r *registrySelector) Mark(service string, node *kvs.Entry, err error) {
	if group, ok := r.services.Get(service); ok {
		group.mark.Set(node.Identifier, err)
	}
}

func (r *registrySelector) Reset(service string) {
	if group, ok := r.services.Get(service); ok {
		group.mark.Clear()
	}
}

func (r *registrySelector) Close() error {
	close(r.stop)
	return nil
}

func (r *registrySelector) String() string {
	return "registry"
}

func (r *registrySelector) getEntries(service string) (kvs.Entries, error) {
	if group, ok := r.services.Get(service); ok {
		return group.getEntries(), nil
	}
	group := &serviceGroup{
		registry:     r.op.Registry,
		registryPath: r.op.RegistryPathResolver,
		name:         service,
		entries:      make(kvs.Entries, 0),
		mark:         cmap.NewConcurrentMap(),
	}
	if r.services.SetNX(service, group) {
		_ = group.init(r.stop)
	} else {
		s, ok := r.services.Get(service)
		if !ok {
			return nil, ErrNotFound
		}
		group = s
	}
	return group.getEntries(), nil
}
