package main

import (
	log "github.com/golang/glog"
	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/discovery-srv/handler"
	proto "github.com/micro/discovery-srv/proto/discovery"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
)

func main() {
	cmd.Init()

	server.Init(
		server.Name("go.micro.srv.discovery"),
	)

	proto.RegisterDiscoveryHandler(server.DefaultServer, new(handler.Discovery))

	server.Subscribe(
		server.NewSubscriber(discovery.HeartbeatTopic, discovery.DefaultDiscovery.ProcessHeartbeat),
	)

	server.Subscribe(
		server.NewSubscriber(discovery.WatchTopic, discovery.DefaultDiscovery.ProcessResult),
	)

	discovery.DefaultDiscovery.Run()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
