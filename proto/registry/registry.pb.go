// Code generated by protoc-gen-go.
// source: github.com/micro/discovery-srv/proto/registry/registry.proto
// DO NOT EDIT!

/*
Package registry is a generated protocol buffer package.

It is generated from these files:
	github.com/micro/discovery-srv/proto/registry/registry.proto

It has these top-level messages:
	RegisterRequest
	RegisterResponse
	DeregisterRequest
	DeregisterResponse
	GetServiceRequest
	GetServiceResponse
	ListServicesRequest
	ListServicesResponse
	WatchRequest
	WatchResponse
*/
package registry

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import discovery "github.com/micro/go-platform/discovery/proto"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RegisterRequest struct {
	Service *discovery.Service `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *RegisterRequest) Reset()                    { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()               {}
func (*RegisterRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RegisterRequest) GetService() *discovery.Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type RegisterResponse struct {
}

func (m *RegisterResponse) Reset()                    { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string            { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()               {}
func (*RegisterResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type DeregisterRequest struct {
	Service *discovery.Service `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *DeregisterRequest) Reset()                    { *m = DeregisterRequest{} }
func (m *DeregisterRequest) String() string            { return proto.CompactTextString(m) }
func (*DeregisterRequest) ProtoMessage()               {}
func (*DeregisterRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DeregisterRequest) GetService() *discovery.Service {
	if m != nil {
		return m.Service
	}
	return nil
}

type DeregisterResponse struct {
}

func (m *DeregisterResponse) Reset()                    { *m = DeregisterResponse{} }
func (m *DeregisterResponse) String() string            { return proto.CompactTextString(m) }
func (*DeregisterResponse) ProtoMessage()               {}
func (*DeregisterResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type GetServiceRequest struct {
	Service string `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *GetServiceRequest) Reset()                    { *m = GetServiceRequest{} }
func (m *GetServiceRequest) String() string            { return proto.CompactTextString(m) }
func (*GetServiceRequest) ProtoMessage()               {}
func (*GetServiceRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type GetServiceResponse struct {
	Services []*discovery.Service `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

func (m *GetServiceResponse) Reset()                    { *m = GetServiceResponse{} }
func (m *GetServiceResponse) String() string            { return proto.CompactTextString(m) }
func (*GetServiceResponse) ProtoMessage()               {}
func (*GetServiceResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *GetServiceResponse) GetServices() []*discovery.Service {
	if m != nil {
		return m.Services
	}
	return nil
}

type ListServicesRequest struct {
}

func (m *ListServicesRequest) Reset()                    { *m = ListServicesRequest{} }
func (m *ListServicesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListServicesRequest) ProtoMessage()               {}
func (*ListServicesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type ListServicesResponse struct {
	Services []*discovery.Service `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

func (m *ListServicesResponse) Reset()                    { *m = ListServicesResponse{} }
func (m *ListServicesResponse) String() string            { return proto.CompactTextString(m) }
func (*ListServicesResponse) ProtoMessage()               {}
func (*ListServicesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *ListServicesResponse) GetServices() []*discovery.Service {
	if m != nil {
		return m.Services
	}
	return nil
}

type WatchRequest struct {
	Service string `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

func (m *WatchRequest) Reset()                    { *m = WatchRequest{} }
func (m *WatchRequest) String() string            { return proto.CompactTextString(m) }
func (*WatchRequest) ProtoMessage()               {}
func (*WatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type WatchResponse struct {
	Action  string             `protobuf:"bytes,1,opt,name=action" json:"action,omitempty"`
	Service *discovery.Service `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
}

func (m *WatchResponse) Reset()                    { *m = WatchResponse{} }
func (m *WatchResponse) String() string            { return proto.CompactTextString(m) }
func (*WatchResponse) ProtoMessage()               {}
func (*WatchResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *WatchResponse) GetService() *discovery.Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func init() {
	proto.RegisterType((*RegisterRequest)(nil), "RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "RegisterResponse")
	proto.RegisterType((*DeregisterRequest)(nil), "DeregisterRequest")
	proto.RegisterType((*DeregisterResponse)(nil), "DeregisterResponse")
	proto.RegisterType((*GetServiceRequest)(nil), "GetServiceRequest")
	proto.RegisterType((*GetServiceResponse)(nil), "GetServiceResponse")
	proto.RegisterType((*ListServicesRequest)(nil), "ListServicesRequest")
	proto.RegisterType((*ListServicesResponse)(nil), "ListServicesResponse")
	proto.RegisterType((*WatchRequest)(nil), "WatchRequest")
	proto.RegisterType((*WatchResponse)(nil), "WatchResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Registry service

type RegistryClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*RegisterResponse, error)
	Deregister(ctx context.Context, in *DeregisterRequest, opts ...client.CallOption) (*RegisterResponse, error)
	GetService(ctx context.Context, in *GetServiceRequest, opts ...client.CallOption) (*GetServiceResponse, error)
	ListServices(ctx context.Context, in *ListServicesRequest, opts ...client.CallOption) (*ListServicesResponse, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...client.CallOption) (Registry_WatchClient, error)
}

type registryClient struct {
	c           client.Client
	serviceName string
}

func NewRegistryClient(serviceName string, c client.Client) RegistryClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "registry"
	}
	return &registryClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *registryClient) Register(ctx context.Context, in *RegisterRequest, opts ...client.CallOption) (*RegisterResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Registry.Register", in)
	out := new(RegisterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Deregister(ctx context.Context, in *DeregisterRequest, opts ...client.CallOption) (*RegisterResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Registry.Deregister", in)
	out := new(RegisterResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) GetService(ctx context.Context, in *GetServiceRequest, opts ...client.CallOption) (*GetServiceResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Registry.GetService", in)
	out := new(GetServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) ListServices(ctx context.Context, in *ListServicesRequest, opts ...client.CallOption) (*ListServicesResponse, error) {
	req := c.c.NewRequest(c.serviceName, "Registry.ListServices", in)
	out := new(ListServicesResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Watch(ctx context.Context, in *WatchRequest, opts ...client.CallOption) (Registry_WatchClient, error) {
	req := c.c.NewRequest(c.serviceName, "Registry.Watch", &WatchRequest{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return &registryWatchClient{stream}, nil
}

type Registry_WatchClient interface {
	RecvR() (*WatchResponse, error)
	client.Streamer
}

type registryWatchClient struct {
	client.Streamer
}

func (x *registryWatchClient) RecvR() (*WatchResponse, error) {
	m := new(WatchResponse)
	err := x.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Registry service

type RegistryHandler interface {
	Register(context.Context, *RegisterRequest, *RegisterResponse) error
	Deregister(context.Context, *DeregisterRequest, *RegisterResponse) error
	GetService(context.Context, *GetServiceRequest, *GetServiceResponse) error
	ListServices(context.Context, *ListServicesRequest, *ListServicesResponse) error
	Watch(context.Context, *WatchRequest, Registry_WatchStream) error
}

func RegisterRegistryHandler(s server.Server, hdlr RegistryHandler) {
	s.Handle(s.NewHandler(&Registry{hdlr}))
}

type Registry struct {
	RegistryHandler
}

func (h *Registry) Register(ctx context.Context, in *RegisterRequest, out *RegisterResponse) error {
	return h.RegistryHandler.Register(ctx, in, out)
}

func (h *Registry) Deregister(ctx context.Context, in *DeregisterRequest, out *RegisterResponse) error {
	return h.RegistryHandler.Deregister(ctx, in, out)
}

func (h *Registry) GetService(ctx context.Context, in *GetServiceRequest, out *GetServiceResponse) error {
	return h.RegistryHandler.GetService(ctx, in, out)
}

func (h *Registry) ListServices(ctx context.Context, in *ListServicesRequest, out *ListServicesResponse) error {
	return h.RegistryHandler.ListServices(ctx, in, out)
}

func (h *Registry) Watch(ctx context.Context, stream server.Streamer) error {
	m := new(WatchRequest)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.RegistryHandler.Watch(ctx, m, &registryWatchStream{stream})
}

type Registry_WatchStream interface {
	SendR(*WatchResponse) error
	server.Streamer
}

type registryWatchStream struct {
	server.Streamer
}

func (x *registryWatchStream) SendR(m *WatchResponse) error {
	return x.Streamer.Send(m)
}

var fileDescriptor0 = []byte{
	// 339 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x8c, 0x92, 0x4f, 0x4f, 0x32, 0x31,
	0x10, 0xc6, 0x81, 0x37, 0x2f, 0xe2, 0xc8, 0xdf, 0x01, 0x12, 0xed, 0x45, 0xd3, 0x78, 0x30, 0x46,
	0x66, 0x11, 0x63, 0x4c, 0x8c, 0xde, 0x4c, 0xbc, 0x78, 0xc2, 0x83, 0x67, 0xa8, 0x15, 0x9a, 0x08,
	0xc5, 0xb6, 0x90, 0xf0, 0x4d, 0xfc, 0xb8, 0x6e, 0x76, 0x0b, 0xbb, 0xeb, 0x62, 0xc2, 0xad, 0xd3,
	0x3c, 0xcf, 0xaf, 0x9b, 0xdf, 0x2c, 0x3c, 0x4c, 0x94, 0x9b, 0x2e, 0xc7, 0x24, 0xf4, 0x2c, 0x98,
	0x29, 0x61, 0x74, 0xf0, 0xae, 0xac, 0xd0, 0x2b, 0x69, 0xd6, 0x3d, 0x6b, 0x56, 0xc1, 0xc2, 0x68,
	0xa7, 0x03, 0x23, 0x27, 0xca, 0x3a, 0xb3, 0xde, 0x1e, 0x28, 0xba, 0x67, 0xf9, 0xf6, 0x44, 0xf7,
	0x16, 0x9f, 0x23, 0xf7, 0xa1, 0xcd, 0x2c, 0x21, 0x79, 0xca, 0x76, 0x8e, 0xdb, 0xfc, 0x0a, 0x1a,
	0xc3, 0x88, 0x27, 0xcd, 0x50, 0x7e, 0x2d, 0xa5, 0x75, 0x78, 0x02, 0x07, 0x56, 0x9a, 0x95, 0x12,
	0xf2, 0xb8, 0x78, 0x56, 0xbc, 0x38, 0x1a, 0x54, 0xe8, 0x35, 0x9e, 0x39, 0x42, 0x33, 0x49, 0xdb,
	0x85, 0x9e, 0x5b, 0xc9, 0x09, 0x5a, 0x4f, 0xd2, 0xec, 0xcf, 0xe8, 0x00, 0xa6, 0xf3, 0x9e, 0x72,
	0x0e, 0xad, 0x67, 0xe9, 0x7c, 0x66, 0x43, 0x69, 0x64, 0x29, 0x87, 0xbc, 0x0f, 0x98, 0x4e, 0xc5,
	0x5d, 0x64, 0x50, 0xf1, 0x31, 0x1b, 0xe6, 0xfe, 0x65, 0x5e, 0xeb, 0x42, 0xfb, 0x25, 0x7c, 0xc9,
	0x8f, 0xd6, 0x93, 0xf9, 0x00, 0x3a, 0xd9, 0xeb, 0x3d, 0x50, 0xa7, 0x50, 0x7d, 0x1b, 0x39, 0x31,
	0xfd, 0xf3, 0xeb, 0xee, 0xa1, 0xe6, 0x03, 0x9e, 0x56, 0x87, 0xf2, 0x48, 0x38, 0xa5, 0xe7, 0x71,
	0x20, 0x6d, 0xa5, 0x94, 0xb5, 0x32, 0xf8, 0x2e, 0x41, 0x65, 0xe8, 0x17, 0x8b, 0xd7, 0x9b, 0xb3,
	0x34, 0xd8, 0xa4, 0x5f, 0xfb, 0x61, 0x2d, 0xca, 0xed, 0xa0, 0x80, 0xb7, 0x00, 0x89, 0x55, 0x44,
	0xca, 0xad, 0x64, 0x77, 0xed, 0x0e, 0x20, 0x11, 0x1a, 0xd6, 0x72, 0x3b, 0x60, 0x6d, 0xca, 0x1b,
	0x0f, 0x8b, 0x8f, 0x50, 0x4d, 0x0b, 0xc4, 0x0e, 0xed, 0xd0, 0xcc, 0xba, 0xb4, 0xcb, 0x72, 0x58,
	0xbf, 0x84, 0xff, 0x91, 0x2a, 0xac, 0x51, 0xda, 0x29, 0xab, 0x53, 0xc6, 0x20, 0x2f, 0xf4, 0x8b,
	0xe3, 0x72, 0xf4, 0xa7, 0xde, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0x13, 0x62, 0x81, 0xce, 0x27,
	0x03, 0x00, 0x00,
}
