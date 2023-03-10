// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: wolong/routing/v2/routing_service.proto

package routingv2

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RoutingServiceClient is the client API for RoutingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoutingServiceClient interface {
	// Get route by type/start/end
	GetRoute(ctx context.Context, in *GetRouteRequest, opts ...grpc.CallOption) (*GetRouteResponse, error)
}

type routingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoutingServiceClient(cc grpc.ClientConnInterface) RoutingServiceClient {
	return &routingServiceClient{cc}
}

func (c *routingServiceClient) GetRoute(ctx context.Context, in *GetRouteRequest, opts ...grpc.CallOption) (*GetRouteResponse, error) {
	out := new(GetRouteResponse)
	err := c.cc.Invoke(ctx, "/wolong.routing.v2.RoutingService/GetRoute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoutingServiceServer is the server API for RoutingService service.
// All implementations should embed UnimplementedRoutingServiceServer
// for forward compatibility
type RoutingServiceServer interface {
	// Get route by type/start/end
	GetRoute(context.Context, *GetRouteRequest) (*GetRouteResponse, error)
}

// UnimplementedRoutingServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRoutingServiceServer struct {
}

func (UnimplementedRoutingServiceServer) GetRoute(context.Context, *GetRouteRequest) (*GetRouteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoute not implemented")
}

// UnsafeRoutingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoutingServiceServer will
// result in compilation errors.
type UnsafeRoutingServiceServer interface {
	mustEmbedUnimplementedRoutingServiceServer()
}

func RegisterRoutingServiceServer(s grpc.ServiceRegistrar, srv RoutingServiceServer) {
	s.RegisterService(&RoutingService_ServiceDesc, srv)
}

func _RoutingService_GetRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingServiceServer).GetRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/wolong.routing.v2.RoutingService/GetRoute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingServiceServer).GetRoute(ctx, req.(*GetRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoutingService_ServiceDesc is the grpc.ServiceDesc for RoutingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoutingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wolong.routing.v2.RoutingService",
	HandlerType: (*RoutingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRoute",
			Handler:    _RoutingService_GetRoute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "wolong/routing/v2/routing_service.proto",
}
