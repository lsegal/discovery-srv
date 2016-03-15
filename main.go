package main

import (
	"log"

	"github.com/micro/go-micro"

	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/discovery-srv/handler"

	proto "github.com/micro/discovery-srv/proto/discovery"
	proto2 "github.com/micro/discovery-srv/proto/registry"
)

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.discovery"),
	)

	service.Init(
		micro.BeforeStart(func() error {
			discovery.Init(service)
			discovery.Run()
			return nil
		}),
	)

	service.Server().Subscribe(
		service.Server().NewSubscriber(discovery.HeartbeatTopic, discovery.DefaultDiscovery.ProcessHeartbeat),
	)

	service.Server().Subscribe(
		service.Server().NewSubscriber(discovery.WatchTopic, discovery.DefaultDiscovery.ProcessResult),
	)

	proto.RegisterDiscoveryHandler(service.Server(), new(handler.Discovery))
	proto2.RegisterRegistryHandler(service.Server(), new(handler.Registry))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
