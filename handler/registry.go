package handler

import (
	"github.com/micro/discovery-srv/discovery"
	proto "github.com/micro/discovery-srv/proto/registry"
	"github.com/micro/discovery-srv/registry"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Registry struct{}

func (r *Registry) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.RegisterResponse) error {
	if req.Service == nil || len(req.Service.Nodes) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Register", "bad request")
	}

	if err := registry.DefaultRegistry.Register(req.Service); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Register", err.Error())
	}

	if err := discovery.Register(ctx, req.Service); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Register", err.Error())
	}

	return nil
}

func (r *Registry) Deregister(ctx context.Context, req *proto.DeregisterRequest, rsp *proto.DeregisterResponse) error {
	if req.Service == nil || len(req.Service.Nodes) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Deregister", "bad request")
	}

	if err := registry.DefaultRegistry.Deregister(req.Service); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Deregister", err.Error())
	}

	if err := discovery.Deregister(ctx, req.Service); err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Deregister", err.Error())
	}

	return nil
}

func (r *Registry) GetService(ctx context.Context, req *proto.GetServiceRequest, rsp *proto.GetServiceResponse) error {
	if len(req.Service) == 0 {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", "bad request")
	}

	services, err := registry.DefaultRegistry.GetService(req.Service)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", err.Error())
	}

	rsp.Services = services

	return nil
}

func (r *Registry) ListServices(ctx context.Context, req *proto.ListServicesRequest, rsp *proto.ListServicesResponse) error {
	services, err := registry.DefaultRegistry.ListServices()
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.GetService", err.Error())
	}

	rsp.Services = services

	return nil
}

func (r *Registry) Watch(ctx context.Context, req *proto.WatchRequest, stream proto.Registry_WatchStream) error {
	watcher, err := registry.DefaultRegistry.Watch(req.Service)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
	}
	for {
		update, err := watcher.Next()
		if err != nil {
			return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
		}
		if err := stream.SendR(&proto.WatchResponse{
			Action:  update.Action,
			Service: update.Service,
		}); err != nil {
			return errors.InternalServerError("go.micro.srv.discovery.Registry.Watch", err.Error())
		}
	}
	return nil
}
