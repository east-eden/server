package gate

import (
	"e.coding.net/mmstudio/blade/gate/msg"
	"e.coding.net/mmstudio/blade/gate/selector"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

type FilterPlugin func(*msg.Handshake) selector.Filter

type pluginManager struct {
	sync.RWMutex
	plugin map[string]FilterPlugin
	path   string
}

func newPluginManager(path string) *pluginManager {
	ret := &pluginManager{
		plugin: make(map[string]FilterPlugin, 0),
		path:   path,
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("[pluginManager] create patch path fail")
		}
	}
	return ret
}

func (p *pluginManager) getFilter(handshake *msg.Handshake, fs []selector.Filter) []selector.Filter {
	p.RLock()
	defer p.RUnlock()
	for _, v := range p.plugin {
		fs = append(fs, v(handshake))
	}
	return fs
}

func getPluginName(name string) string {
	patchName := filepath.Base(name)
	return patchName[:len(patchName)-3]
}

func (p *pluginManager) removePlugin(name string) {
	p.Lock()
	defer p.Unlock()
	pluginName := getPluginName(name)
	delete(p.plugin, pluginName)
	log.Warn().Str("name", pluginName).Msg("[pluginManager] remove plugin success")
}

func (p *pluginManager) addPlugin(name string) {
	pluginName := getPluginName(name)
	pluginPath := filepath.Join(p.path, pluginName+".so")
	pg, err := plugin.Open(pluginPath)
	if err != nil {
		log.Error().Err(err).Str("name", pluginName).Str("path", pluginPath).Msg("[pluginManager] load plugin fail")
		return
	}
	f, err := pg.Lookup("Filter")
	if err != nil {
		log.Error().Err(err).Str("name", pluginName).Str("path", pluginPath).Msg("[pluginManager] lookup Filter fail")
		return
	}
	if fn, ok := f.(func(*msg.Handshake) selector.Filter); ok {
		p.Lock()
		p.plugin[pluginName] = fn
		p.Unlock()
		log.Warn().Str("name", pluginName).Str("path", pluginPath).Msg("[pluginManager] add plugin success")
	} else {
		log.Error().Str("name", pluginName).Str("path", pluginPath).Msg("[pluginManager] assert FilterPlugin fail")
	}
}

func (p *pluginManager) watchPlugin(stop chan struct{}) {
	_ = filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".so") {
			return nil
		}
		log.Debug().Str("path", p.path).Str("plugin", path).Msg("[pluginManager] load plugin")
		p.addPlugin(path)
		return nil
	})

	var err error
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Str("path", p.path).Msg("[pluginManager] new fswatch fail")
	}
	defer watcher.Close()
	err = watcher.Add(p.path)
	if err != nil {
		log.Fatal().Err(err).Str("path", p.path).Msg("[pluginManager] fswatch add path fail")
	}
	log.Debug().Str("path", p.path).Msg("[pluginManager] start watch plugin")
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			log.Debug().Str("path", p.path).Str("event", event.String()).Msg("[pluginManager] watch event")

			if !strings.HasSuffix(event.Name, ".so") {
				continue
			}
			if (event.Op&fsnotify.Write) == fsnotify.Write ||
				(event.Op&fsnotify.Create) == fsnotify.Create {
				p.addPlugin(event.Name)
			} else if (event.Op & fsnotify.Remove) == fsnotify.Remove {
				p.removePlugin(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("[pluginManager] watch recv error")
		case <-stop:
			return
		}
	}
}

// MakeFilter
func (g *Gate) MakeFilter(handshake *msg.Handshake) []selector.Filter {
	filters := make([]selector.Filter, 0)
	// custom input filter
	if g.spec.Filter != nil {
		for _, f := range g.spec.Filter {
			filters = append(filters, f(handshake))
		}
	}

	// plugin filter
	filters = g.pluginMgr.getFilter(handshake, filters)

	return filters
}
