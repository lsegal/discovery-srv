package main

import (
	log "github.com/golang/glog"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"

	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/discovery-srv/handler"

	proto "github.com/micro/discovery-srv/proto/discovery"
	proto2 "github.com/micro/discovery-srv/proto/registry"
)

func main() {
	cmd.Init()

	server.Init(
		server.Name("go.micro.srv.discovery"),
	)

	proto.RegisterDiscoveryHandler(server.DefaultServer, new(handler.Discovery))
	proto2.RegisterRegistryHandler(server.DefaultServer, new(handler.Registry))

	server.Subscribe(
		server.NewSubscriber(discovery.HeartbeatTopic, discovery.DefaultDiscovery.ProcessHeartbeat),
	)

	server.Subscribe(
		server.NewSubscriber(discovery.WatchTopic, discovery.DefaultDiscovery.ProcessResult),
	)

	discovery.DefaultDiscovery.Init()
	discovery.DefaultDiscovery.Run()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
