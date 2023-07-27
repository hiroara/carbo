// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: communicator.proto

package pb

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

const (
	Communicator_FillBatch_FullMethodName = "/Communicator/FillBatch"
	Communicator_GetBatch_FullMethodName  = "/Communicator/GetBatch"
)

// CommunicatorClient is the client API for Communicator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommunicatorClient interface {
	FillBatch(ctx context.Context, in *FillBatchRequest, opts ...grpc.CallOption) (*FillBatchResponse, error)
	GetBatch(ctx context.Context, in *GetBatchRequest, opts ...grpc.CallOption) (*GetBatchResponse, error)
}

type communicatorClient struct {
	cc grpc.ClientConnInterface
}

func NewCommunicatorClient(cc grpc.ClientConnInterface) CommunicatorClient {
	return &communicatorClient{cc}
}

func (c *communicatorClient) FillBatch(ctx context.Context, in *FillBatchRequest, opts ...grpc.CallOption) (*FillBatchResponse, error) {
	out := new(FillBatchResponse)
	err := c.cc.Invoke(ctx, Communicator_FillBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *communicatorClient) GetBatch(ctx context.Context, in *GetBatchRequest, opts ...grpc.CallOption) (*GetBatchResponse, error) {
	out := new(GetBatchResponse)
	err := c.cc.Invoke(ctx, Communicator_GetBatch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommunicatorServer is the server API for Communicator service.
// All implementations must embed UnimplementedCommunicatorServer
// for forward compatibility
type CommunicatorServer interface {
	FillBatch(context.Context, *FillBatchRequest) (*FillBatchResponse, error)
	GetBatch(context.Context, *GetBatchRequest) (*GetBatchResponse, error)
	mustEmbedUnimplementedCommunicatorServer()
}

// UnimplementedCommunicatorServer must be embedded to have forward compatible implementations.
type UnimplementedCommunicatorServer struct {
}

func (UnimplementedCommunicatorServer) FillBatch(context.Context, *FillBatchRequest) (*FillBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FillBatch not implemented")
}
func (UnimplementedCommunicatorServer) GetBatch(context.Context, *GetBatchRequest) (*GetBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBatch not implemented")
}
func (UnimplementedCommunicatorServer) mustEmbedUnimplementedCommunicatorServer() {}

// UnsafeCommunicatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommunicatorServer will
// result in compilation errors.
type UnsafeCommunicatorServer interface {
	mustEmbedUnimplementedCommunicatorServer()
}

func RegisterCommunicatorServer(s grpc.ServiceRegistrar, srv CommunicatorServer) {
	s.RegisterService(&Communicator_ServiceDesc, srv)
}

func _Communicator_FillBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FillBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicatorServer).FillBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communicator_FillBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicatorServer).FillBatch(ctx, req.(*FillBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Communicator_GetBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommunicatorServer).GetBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Communicator_GetBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommunicatorServer).GetBatch(ctx, req.(*GetBatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Communicator_ServiceDesc is the grpc.ServiceDesc for Communicator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Communicator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Communicator",
	HandlerType: (*CommunicatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FillBatch",
			Handler:    _Communicator_FillBatch_Handler,
		},
		{
			MethodName: "GetBatch",
			Handler:    _Communicator_GetBatch_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "communicator.proto",
}