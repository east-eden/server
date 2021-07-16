package liveness

import (
	"sync"
	"time"

	"e.coding.net/mmstudio/blade/golib/sync2"
	"e.coding.net/mmstudio/blade/golib/time2"
)

var (
	defaultCheckInterval = time.Duration(5) * time.Second
	Default              = New(defaultCheckInterval)
)

type ServiceToken struct {
	id uint64
}

type ServiceInfo struct {
	Token    ServiceToken
	TsNotify sync2.AtomicInt64
	TTL      int64
	Desc     string
}

type Notify interface {
	Notify()
}

func (g *ServiceInfo) Notify() {
	g.TsNotify.Set(time2.Unix())
}

type Liveness struct {
	watchMap          sync.Map
	stop              chan struct{}
	serviceCounter    sync2.AtomicUint64
	checkInterval     sync2.AtomicDuration
	onLivenessFailedL sync.RWMutex
	onLivenessFailed  func([]ServiceInfo)
	running           sync2.AtomicBool
}

func New(duration time.Duration) *Liveness {
	d := &Liveness{}
	d.checkInterval.Set(duration)
	d.stop = make(chan struct{})
	return d
}

func (d *Liveness) ServeBackground() {
	go d.Serve()
}

func (d *Liveness) Close() {
	close(d.stop)
}

func (d *Liveness) Remove(token ServiceToken) {
	d.watchMap.Delete(token)
}

func (d *Liveness) SetOnLivenessFailed(f func([]ServiceInfo)) {
	d.onLivenessFailedL.Lock()
	d.onLivenessFailed = f
	d.onLivenessFailedL.Unlock()
}

func (d *Liveness) SetCheckInterval(duration time.Duration) {
	d.checkInterval.Set(duration)
}

func (d *Liveness) Add(desc string, ttl int64) (ServiceToken, Notify) {
	info := &ServiceInfo{
		Token: ServiceToken{id: d.serviceCounter.Add(1)},
		TTL:   ttl,
		Desc:  desc,
	}
	info.TsNotify.Set(time2.Unix())
	d.watchMap.Store(info.Token, info)
	return info.Token, info
}

func (d *Liveness) Serve() {
	if d.running.Get() {
		return
	}
	d.running.Set(true)
	defer d.running.Set(false)
	ticker := time.NewTicker(d.checkInterval.Get())
	for {
		select {
		case <-ticker.C:

			d.onLivenessFailedL.RLock()
			if d.onLivenessFailed == nil {
				d.onLivenessFailedL.RUnlock()
				continue
			}
			d.onLivenessFailedL.RUnlock()

			expired := make([]ServiceInfo, 0)
			d.watchMap.Range(func(key, value interface{}) bool {
				info := value.(*ServiceInfo)
				if (time2.Unix() - info.TsNotify.Get()) > info.TTL {
					expired = append(expired, ServiceInfo{
						Token:    info.Token,
						TsNotify: info.TsNotify,
						TTL:      info.TTL,
						Desc:     info.Desc,
					})
				}
				return true
			})
			if d.checkInterval != 0 {
				duration := d.checkInterval.Get()
				d.checkInterval.Set(0)
				ticker.Stop()
				ticker = time.NewTicker(duration)
			}
			d.onLivenessFailedL.RLock()
			if d.onLivenessFailed != nil && len(expired) != 0 {
				d.onLivenessFailed(expired)
			}
			d.onLivenessFailedL.RUnlock()
		case <-d.stop:
			return
		}
	}
}
