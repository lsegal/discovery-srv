package handler

import (
	"time"

	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/go-micro/errors"

	proto "github.com/micro/discovery-srv/proto/registry"
	proto2 "github.com/micro/go-platform/discovery/proto"
	"golang.org/x/net/context"
)

type Registry struct{}

func (r *Registry) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.RegisterResponse) error {
	if req.Service == nil || len(req.Service.Nodes) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Register", "bad request")
	}

	if err := discovery.DefaultDiscovery.Register(toService(req.Service)); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Register", err.Error())
	}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, req *proto.DeregisterRequest, rsp *proto.DeregisterResponse) error {
	if req.Service == nil || len(req.Service.Nodes) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Deregister", "bad request")
	}

	if err := discovery.DefaultDiscovery.Deregister(toService(req.Service)); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Deregister", err.Error())
	}

	return nil
}

func (r *Registry) GetService(ctx context.Context, req *proto.GetServiceRequest, rsp *proto.GetServiceResponse) error {
	if len(req.Service) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", "bad request")
	}

	services, err := discovery.DefaultDiscovery.GetService(req.Service)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", err.Error())
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func (r *Registry) ListServices(ctx context.Context, req *proto.ListServicesRequest, rsp *proto.ListServicesResponse) error {
	services, err := discovery.DefaultDiscovery.ListServices()
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", err.Error())
	}

	for _, service := range services {
		rsp.Services = append(rsp.Services, toProto(service))
	}

	return nil
}

func (r *Registry) Watch(ctx context.Context, req *proto.WatchRequest, stream proto.Registry_WatchStream) error {
	watcher, err := discovery.DefaultDiscovery.Watch()
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
	}
	for {
		update, err := watcher.Next()
		if err != nil {
			return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
		}
		if len(req.Service) > 0 && update.Service.Name != req.Service {
			continue
		}

		if err := stream.Send(&proto.WatchResponse{
			Result: &proto2.Result{
				Action:    update.Action,
				Service:   toProto(update.Service),
				Timestamp: time.Now().Unix(),
			},
		}); err != nil {
			return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
		}
	}
	return nil
}
