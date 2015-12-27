package discovery

import (
	"errors"
	"sync"
	"time"

	"github.com/micro/discovery-srv/registry"
	"github.com/micro/go-micro/client"
	proto "github.com/micro/go-platform/discovery/proto"
	"golang.org/x/net/context"
)

type discovery struct {
	sync.RWMutex
	heartbeats map[string][]*proto.Heartbeat
}

var (
	DefaultDiscovery = newDiscovery()
	ErrNotFound      = errors.New("not found")
	HeartbeatTopic   = "micro.discovery.heartbeat"
	WatchTopic       = "micro.discovery.watch"
	TickInterval     = time.Duration(time.Minute)
)

func newDiscovery() *discovery {
	return &discovery{
		heartbeats: make(map[string][]*proto.Heartbeat),
	}
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

func (d *discovery) reap() {
	d.Lock()
	defer d.Unlock()

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

func (d *discovery) run() {
	t := time.NewTicker(TickInterval)

	for _ = range t.C {
		d.reap()
	}
}

func (d *discovery) Heartbeats(id string, after int64, limit, offset int) ([]*proto.Heartbeat, error) {
	d.RLock()
	defer d.RUnlock()

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

func (d *discovery) ProcessHeartbeat(ctx context.Context, hb *proto.Heartbeat) error {
	d.Lock()
	defer d.Unlock()

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
	registry.DefaultRegistry.Register(hb.Service)
	return nil
}

func (d *discovery) ProcessResult(ctx context.Context, r *proto.Result) error {
	switch r.Action {
	case "create", "update":
		return registry.DefaultRegistry.Register(r.Service)
	case "delete":
		return registry.DefaultRegistry.Deregister(r.Service)
	}
	return nil
}

func (d *discovery) Run() {
	go d.run()
}

func Register(ctx context.Context, service *proto.Service) error {
	p := client.NewPublication(WatchTopic, &proto.Result{
		Action:  "update",
		Service: service,
	})
	return client.Publish(ctx, p)
}

func Deregister(ctx context.Context, service *proto.Service) error {
	p := client.NewPublication(WatchTopic, &proto.Result{
		Action:  "delete",
		Service: service,
	})
	return client.Publish(ctx, p)
}
