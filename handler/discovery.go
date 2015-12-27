package handler

import (
	"github.com/micro/discovery-srv/discovery"
	proto "github.com/micro/discovery-srv/proto/discovery"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Discovery struct{}

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
