package discovery

import (
	"errors"
	"sync"
	"time"

	log "github.com/golang/glog"
	disco "github.com/micro/go-platform/discovery"
	proto "github.com/micro/go-platform/discovery/proto"
	"golang.org/x/net/context"
)

type discovery struct {
	disco.Discovery

	wtx          sync.RWMutex
	watchResults map[string][]*proto.Result

	htx        sync.RWMutex
	heartbeats map[string][]*proto.Heartbeat
}

var (
	DefaultDiscovery = newDiscovery()
	ErrNotFound      = errors.New("not found")
	HeartbeatTopic   = "micro.discovery.heartbeat"
	WatchTopic       = "micro.discovery.watch"

	TickInterval = time.Duration(time.Minute)
	History      = int64(3600)
)

func newDiscovery() *discovery {
	return &discovery{
		watchResults: make(map[string][]*proto.Result),
		heartbeats:   make(map[string][]*proto.Heartbeat),
	}
}

func filterRes(r []*proto.Result, after int64, limit, offset int) []*proto.Result {
	if len(r) < offset {
		return []*proto.Result{}
	}

	if (limit + offset) > len(r) {
		limit = len(r) - offset
	}

	var rs []*proto.Result
	for i := 0; i < limit; i++ {
		if r[offset].Timestamp > after {
			rs = append(rs, r[offset])
		}
		offset++
	}
	return rs
}

func filter(hb []*proto.Heartbeat, after int64, limit, offset int) []*proto.Heartbeat {
	if len(hb) < offset {
		return []*proto.Heartbeat{}
	}

	if (limit + offset) > len(hb) {
		limit = len(hb) - offset
	}

	var hbs []*proto.Heartbeat
	for i := 0; i < limit; i++ {
		if hb[offset].Timestamp > after {
			hbs = append(hbs, hb[offset])
		}
		offset++
	}
	return hbs
}

func (d *discovery) reapHB() {
	d.htx.Lock()
	defer d.htx.Unlock()

	t := time.Now().Unix()
	for id, hb := range d.heartbeats {
		var beats []*proto.Heartbeat
		for _, beat := range hb {
			if t > (beat.Timestamp+beat.Interval) && t > (beat.Timestamp+beat.Ttl) {
				continue
			}
			beats = append(beats, beat)
		}
		d.heartbeats[id] = beats
	}
}

func (d *discovery) reapRes() {
	d.wtx.Lock()
	defer d.wtx.Unlock()

	t := time.Now().Unix()
	for service, results := range d.watchResults {
		var rs []*proto.Result
		for _, res := range results {
			if res.Timestamp < (t - History) {
				continue
			}
			rs = append(rs, res)
		}
		d.watchResults[service] = rs
	}
}

func (d *discovery) run() {
	t := time.NewTicker(TickInterval)

	for _ = range t.C {
		d.reapHB()
		d.reapRes()
	}
}

func (d *discovery) Init() {
	d.Discovery = disco.NewDiscovery(
		disco.EnableDiscovery(false),
	)
}

func (d *discovery) Heartbeats(id string, after int64, limit, offset int) ([]*proto.Heartbeat, error) {
	d.htx.RLock()
	defer d.htx.RUnlock()

	if len(id) == 0 {
		var hbs []*proto.Heartbeat
		for _, hb := range d.heartbeats {
			hbs = append(hbs, hb...)
		}
		return filter(hbs, after, limit, offset), nil
	}

	hbs, ok := d.heartbeats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return filter(hbs, after, limit, offset), nil
}

func (d *discovery) WatchResults(service string, after int64, limit, offset int) ([]*proto.Result, error) {
	d.wtx.RLock()
	defer d.wtx.RUnlock()

	if len(service) == 0 {
		var rs []*proto.Result
		for _, res := range d.watchResults {
			rs = append(rs, res...)
		}
		return filterRes(rs, after, limit, offset), nil
	}

	rs, ok := d.watchResults[service]
	if !ok {
		return nil, ErrNotFound
	}
	return filterRes(rs, after, limit, offset), nil
}

func (d *discovery) ProcessHeartbeat(ctx context.Context, hb *proto.Heartbeat) error {
	d.htx.Lock()
	defer d.htx.Unlock()

	hbs, ok := d.heartbeats[hb.Id]
	if !ok {
		d.heartbeats[hb.Id] = append(d.heartbeats[hb.Id], hb)
		return nil
	}

	for i, h := range hbs {
		if h.Service.Nodes[0].Id == hb.Service.Nodes[0].Id {
			hbs[i] = hb
			d.heartbeats[hb.Id] = hbs
			return nil
		}
	}

	hbs = append(hbs, hb)
	d.heartbeats[hb.Id] = hbs
	return nil
}

func (d *discovery) ProcessResult(ctx context.Context, r *proto.Result) error {
	d.wtx.Lock()
	defer d.wtx.Unlock()
	d.watchResults[r.Service.Name] = append(d.watchResults[r.Service.Name], r)
	return nil
}

func (d *discovery) Run() {
	if d.Discovery == nil {
		log.Fatal("Discovery not initialised")
	}

	d.Discovery.Start()
	go d.run()
}
