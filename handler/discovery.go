package handler

import (
	"github.com/micro/discovery-srv/discovery"
	proto "github.com/micro/discovery-srv/proto/discovery"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Discovery struct{}

func (m *Discovery) Endpoints(ctx context.Context, req *proto.EndpointsRequest, rsp *proto.EndpointsResponse) error {
	if req.Limit == 0 {
		req.Limit = 10
	}
	endpoints, err := discovery.DefaultDiscovery.Endpoints(req.Service, req.Version, int(req.Limit), int(req.Offset))
	if err != nil && err == discovery.ErrNotFound {
		return errors.NotFound("go.micro.srv.discovery.Discovery.Endpoints", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Discovery.Endpoints", err.Error())
	}

	rsp.Endpoints = endpoints
	return nil
}

func (m *Discovery) Heartbeats(ctx context.Context, req *proto.HeartbeatsRequest, rsp *proto.HeartbeatsResponse) error {
	if req.Limit == 0 {
		req.Limit = 10
	}
	hbs, err := discovery.DefaultDiscovery.Heartbeats(req.Id, int64(req.After), int(req.Limit), int(req.Offset))
	if err != nil && err == discovery.ErrNotFound {
		return errors.NotFound("go.micro.srv.discovery.Discovery.Heartbeats", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Discovery.Heartbeats", err.Error())
	}

	rsp.Heartbeats = hbs
	return nil
}

func (m *Discovery) WatchResults(ctx context.Context, req *proto.WatchResultsRequest, rsp *proto.WatchResultsResponse) error {
	if req.Limit == 0 {
		req.Limit = 10
	}
	rs, err := discovery.DefaultDiscovery.WatchResults(req.Service, int64(req.After), int(req.Limit), int(req.Offset))
	if err != nil && err == discovery.ErrNotFound {
		return errors.NotFound("go.micro.srv.discovery.Discovery.WatchResults", err.Error())
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.discovery.Discovery.WatchResults", err.Error())
	}

	rsp.Results = rs
	return nil
}
