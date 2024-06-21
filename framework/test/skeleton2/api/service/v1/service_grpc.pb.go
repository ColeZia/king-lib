// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// Skeleton2ServiceClient is the client API for Skeleton2Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Skeleton2ServiceClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type skeleton2ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSkeleton2ServiceClient(cc grpc.ClientConnInterface) Skeleton2ServiceClient {
	return &skeleton2ServiceClient{cc}
}

func (c *skeleton2ServiceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := c.cc.Invoke(ctx, "/api.service.v1.Skeleton2Service/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Skeleton2ServiceServer is the server API for Skeleton2Service service.
// All implementations must embed UnimplementedSkeleton2ServiceServer
// for forward compatibility
type Skeleton2ServiceServer interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
	mustEmbedUnimplementedSkeleton2ServiceServer()
}

// UnimplementedSkeleton2ServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSkeleton2ServiceServer struct {
}

func (UnimplementedSkeleton2ServiceServer) Get(context.Context, *GetRequest) (*GetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSkeleton2ServiceServer) mustEmbedUnimplementedSkeleton2ServiceServer() {}

// UnsafeSkeleton2ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Skeleton2ServiceServer will
// result in compilation errors.
type UnsafeSkeleton2ServiceServer interface {
	mustEmbedUnimplementedSkeleton2ServiceServer()
}

func RegisterSkeleton2ServiceServer(s grpc.ServiceRegistrar, srv Skeleton2ServiceServer) {
	s.RegisterService(&Skeleton2Service_ServiceDesc, srv)
}

func _Skeleton2Service_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Skeleton2ServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.service.v1.Skeleton2Service/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Skeleton2ServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Skeleton2Service_ServiceDesc is the grpc.ServiceDesc for Skeleton2Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Skeleton2Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.service.v1.Skeleton2Service",
	HandlerType: (*Skeleton2ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Skeleton2Service_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/service/v1/service.proto",
}