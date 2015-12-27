package registry

import (
	"errors"
	"sync"
	"time"

	proto "github.com/micro/go-platform/discovery/proto"
)

type registry struct {
	sync.RWMutex
	services map[string][]*proto.Service

	watchers map[*watcher]string
}

var (
	DefaultRegistry = newRegistry()
	ErrNotFound     = errors.New("not found")
)

func newRegistry() *registry {
	return &registry{
		services: make(map[string][]*proto.Service),
	}
}

func (r *registry) watch(res *proto.Result) {
	r.RLock()
	watchers := r.watchers
	r.RUnlock()

	for watcher, service := range watchers {
		if len(service) == 0 || service == res.Service.Name {
			select {
			case watcher.watch <- res:
			case <-time.After(time.Second):
				// give up
			}
		}
	}
}

func (r *registry) Register(service *proto.Service) error {
	r.Lock()
	defer r.Unlock()

	services, ok := r.services[service.Name]
	if !ok {
		r.services[service.Name] = []*proto.Service{service}
		go r.watch(&proto.Result{Action: "create", Service: service})
		return nil
	}

	go r.watch(&proto.Result{Action: "update", Service: service})

	var srv *proto.Service
	var index int
	for i, s := range services {
		if s.Version == service.Version {
			srv = s
			index = i
			break
		}
	}

	if srv == nil {
		r.services[service.Name] = append(r.services[service.Name], service)
		return nil
	}

	for _, cur := range srv.Nodes {
		var seen bool
		for _, node := range service.Nodes {
			if cur.Id == node.Id {
				seen = true
				break
			}
		}
		if !seen {
			service.Nodes = append(service.Nodes, cur)
		}
	}

	services[index] = service
	r.services[service.Name] = services
	return nil
}

func (r *registry) Deregister(service *proto.Service) error {
	r.Lock()
	defer r.Unlock()

	go r.watch(&proto.Result{Action: "delete", Service: service})

	services, ok := r.services[service.Name]
	if !ok {
		return nil
	}

	if len(service.Nodes) == 0 {
		delete(r.services, service.Name)
		return nil
	}

	var srv *proto.Service
	var index int
	for i, s := range services {
		if s.Version == service.Version {
			srv = s
			index = i
			break
		}
	}

	if srv == nil {
		return nil
	}

	var nodes []*proto.Node

	// filter cur nodes to remove the dead one
	for _, cur := range srv.Nodes {
		var seen bool
		for _, del := range service.Nodes {
			if del.Id == cur.Id {
				seen = true
				break
			}
		}
		if !seen {
			nodes = append(nodes, cur)
		}
	}

	if len(nodes) == 0 {
		if len(services) == 1 {
			delete(r.services, srv.Name)
		} else {
			var srvs []*proto.Service
			for _, s := range services {
				if s.Version != srv.Version {
					srvs = append(srvs, s)
				}
			}
			r.services[srv.Name] = srvs
		}
		return nil
	}

	service.Nodes = nodes
	services[index] = service
	r.services[service.Name] = services

	return nil
}

func (r *registry) GetService(service string) ([]*proto.Service, error) {
	r.RLock()
	defer r.RUnlock()
	services, ok := r.services[service]
	if !ok {
		return nil, ErrNotFound
	}
	return services, nil
}

func (r *registry) ListServices() ([]*proto.Service, error) {
	r.RLock()
	defer r.RUnlock()
	var services []*proto.Service
	for _, service := range r.services {
		services = append(services, service...)
	}
	return services, nil
}

func (r *registry) Watch(service string) (*watcher, error) {
	watch := make(chan *proto.Result, 100)
	stop := make(chan bool)

	watcher := newWatcher(watch, stop)

	r.Lock()
	r.watchers[watcher] = service
	r.Unlock()

	go func() {
		<-stop
		r.Lock()
		delete(r.watchers, watcher)
		r.Unlock()
	}()

	return watcher, nil
}
