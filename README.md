# Discovery Server

Discovery server is a microservice which layers on the registry to provide heartbeating, in memory caching and much more. 
It subscribes to heartbeats and maintains a registry based on liveness.

It's built with the [Eureka 2.0](https://github.com/Netflix/eureka/wiki/Eureka-2.0-Architecture-Overview) design in mind. 
The Discovery service acts as a read layer cache where the usage of a Registry like Consul, Etcd, Zookeeper act as the 
write layer. With the combination of the Registry, Discovery and Platform we can develop a highly available discovery 
system.


## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Download and start the service

	```shell
	go get github.com/micro/discovery-srv
	discovery-srv
	```

	OR as a docker container

	```shell
	docker run microhq/discovery-srv --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Discovery server implements the following RPC Methods

### Discovery.Heartbeats
```shell
micro query go.micro.srv.discovery Discovery.Heartbeats
{
	"heartbeats": [
		{
			"id": "foo-123",
			"interval": 1,
			"service": {
				"endpoints": [
					{
						"metadata": {
							"index": "Handles index requests"
						},
						"name": "/index",
						"request": {
							"name": "request",
							"type": "Request"
						},
						"response": {
							"name": "response",
							"type": "Response"
						}
					}
				],
				"metadata": {
					"foo": "bar"
				},
				"name": "go.micro.srv.foo",
				"nodes": [
					{
						"address": "localhost",
						"id": "foo-123",
						"metadata": {
							"bar": "baz"
						},
						"port": 8080
					}
				],
				"version": "latest"
			},
			"timestamp": 1.451177551e+09,
			"ttl": 5
		}
	]
}
```

### Discovery.Endpoints

```shell
micro query go.micro.srv.discovery Discovery.Endpoints
{
	"endpoints": [
		{
			"endpoint": {
				"metadata": {
					"index": "Handles index requests"
				},
				"name": "/index",
				"request": {
					"name": "request",
					"type": "Request"
				},
				"response": {
					"name": "response",
					"type": "Response"
				}
			},
			"service": "go.micro.srv.foo",
			"version": "latest"
		}
	]
}
```

### Sending Heartbeats

Heartbeats are sent to the discovery service using [go-platform/discovery](https://github.com/micro/go-platform/tree/master/discovery)

