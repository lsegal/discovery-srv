package discovery

import (
	"errors"
	"log"
	"sync"
	"time"

	proto2 "github.com/micro/discovery-srv/proto/discovery"
	"github.com/micro/go-micro"
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

	etx       sync.RWMutex
	endpoints map[string][]*proto2.ServiceEndpoint
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
		endpoints:    make(map[string][]*proto2.ServiceEndpoint),
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

func filterEps(eps []*proto2.ServiceEndpoint, limit, offset int) []*proto2.ServiceEndpoint {
	if len(eps) < offset {
		return []*proto2.ServiceEndpoint{}
	}

	if (limit + offset) > len(eps) {
		limit = len(eps) - offset
	}

	var newEps []*proto2.ServiceEndpoint
	for i := 0; i < limit; i++ {
		newEps = append(newEps, eps[offset])
		offset++
	}
	return newEps
}

func (d *discovery) reapEps() {
	versions := make(map[string][]string)

	// Create a service version map
	d.htx.RLock()
	for _, hb := range d.heartbeats {
		for _, beat := range hb {
			service, ok := versions[beat.Service.Name]
			if !ok {
				versions[beat.Service.Name] = []string{beat.Service.Version}
				continue
			}
			var seen bool
			for _, version := range service {
				if version == beat.Service.Version {
					seen = true
					break
				}
			}
			if !seen {
				service = append(service, beat.Service.Version)
				versions[beat.Service.Name] = service
			}
		}
	}
	d.htx.RUnlock()

	d.etx.Lock()
	defer d.etx.Unlock()

	endpoints := make(map[string][]*proto2.ServiceEndpoint)

	// Create new endpoint map

	for service, epoints := range d.endpoints {
		vers, ok := versions[service]
		if !ok || len(vers) == 0 {
			continue
		}

		var eps []*proto2.ServiceEndpoint

		for _, ep := range epoints {
			for _, version := range vers {
				if version == ep.Version {
					eps = append(eps, ep)
					break
				}
			}
		}

		endpoints[service] = eps
	}

	d.endpoints = endpoints

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
		d.reapEps()
	}
}

func (d *discovery) Init(s micro.Service) {
	d.Discovery = disco.NewDiscovery(
		disco.EnableDiscovery(false),
		disco.Registry(s.Options().Registry),
		disco.Client(s.Client()),
	)
}

func (d *discovery) Endpoints(service, version string, limit, offset int) ([]*proto2.ServiceEndpoint, error) {
	d.etx.RLock()
	defer d.etx.RUnlock()

	if len(service) == 0 {
		var eps []*proto2.ServiceEndpoint
		for _, ep := range d.endpoints {
			eps = append(eps, ep...)
		}
		return filterEps(eps, limit, offset), nil
	}

	eps, ok := d.endpoints[service]
	if !ok {
		return nil, ErrNotFound
	}

	if len(version) > 0 {
		var newEps []*proto2.ServiceEndpoint
		for _, ep := range eps {
			if ep.Version != version {
				continue
			}
			newEps = append(newEps, ep)
		}
		eps = newEps
	}

	return filterEps(eps, limit, offset), nil
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

	heartbeats := hbs

	var seen bool
	for i, h := range heartbeats {
		if h.Service.Nodes[0].Id == hb.Service.Nodes[0].Id {
			heartbeats[i] = hb
			seen = true
			break
		}
	}
	if !seen {
		heartbeats = append(heartbeats, hb)
	}
	d.heartbeats[hb.Id] = heartbeats
	return nil
}

func (d *discovery) ProcessResult(ctx context.Context, r *proto.Result) error {
	d.wtx.Lock()
	d.watchResults[r.Service.Name] = append(d.watchResults[r.Service.Name], r)
	d.wtx.Unlock()

	d.etx.Lock()
	defer d.etx.Unlock()

	// add endpoints
	service, ok := d.endpoints[r.Service.Name]
	if !ok {
		if r.Action == "delete" {
			return nil
		}

		var endpoints []*proto2.ServiceEndpoint

		for _, ep := range r.Service.Endpoints {
			endpoints = append(endpoints, &proto2.ServiceEndpoint{
				Service:  r.Service.Name,
				Version:  r.Service.Version,
				Endpoint: ep,
			})
		}

		d.endpoints[r.Service.Name] = endpoints
	}

	// TODO: handle deletes
	if r.Action == "delete" {
		return nil
	}

	var endpoints []*proto2.ServiceEndpoint

	// keep what we don't touch
	for _, endpoint := range service {
		// skip this service
		if r.Service.Version == endpoint.Version {
			continue
		}
		endpoints = append(endpoints, endpoint)
	}

	// add new endpoints
	for _, endpoint := range r.Service.Endpoints {
		endpoints = append(endpoints, &proto2.ServiceEndpoint{
			Service:  r.Service.Name,
			Version:  r.Service.Version,
			Endpoint: endpoint,
		})
	}

	// update the endpoints
	d.endpoints[r.Service.Name] = endpoints

	return nil
}

func (d *discovery) Run() {
	if d.Discovery == nil {
		log.Fatal("Discovery not initialised")
	}

	d.Discovery.Start()
	go d.run()
}

func Init(s micro.Service) {
	DefaultDiscovery.Init(s)
}

func Run() {
	DefaultDiscovery.Run()
}
