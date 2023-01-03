// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: grpc/interface.proto

package exam

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

// DictionaryClient is the client API for Dictionary service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DictionaryClient interface {
	//Bid = Add
	Add(ctx context.Context, in *RequestAdd, opts ...grpc.CallOption) (*Ack, error)
	//Result = Read
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadOutcome, error)
}

type dictionaryClient struct {
	cc grpc.ClientConnInterface
}

func NewDictionaryClient(cc grpc.ClientConnInterface) DictionaryClient {
	return &dictionaryClient{cc}
}

func (c *dictionaryClient) Add(ctx context.Context, in *RequestAdd, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/exam.Dictionary/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dictionaryClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadOutcome, error) {
	out := new(ReadOutcome)
	err := c.cc.Invoke(ctx, "/exam.Dictionary/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DictionaryServer is the server API for Dictionary service.
// All implementations must embed UnimplementedDictionaryServer
// for forward compatibility
type DictionaryServer interface {
	//Bid = Add
	Add(context.Context, *RequestAdd) (*Ack, error)
	//Result = Read
	Read(context.Context, *ReadRequest) (*ReadOutcome, error)
	mustEmbedUnimplementedDictionaryServer()
}

// UnimplementedDictionaryServer must be embedded to have forward compatible implementations.
type UnimplementedDictionaryServer struct {
}

func (UnimplementedDictionaryServer) Add(context.Context, *RequestAdd) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedDictionaryServer) Read(context.Context, *ReadRequest) (*ReadOutcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedDictionaryServer) mustEmbedUnimplementedDictionaryServer() {}

// UnsafeDictionaryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DictionaryServer will
// result in compilation errors.
type UnsafeDictionaryServer interface {
	mustEmbedUnimplementedDictionaryServer()
}

func RegisterDictionaryServer(s grpc.ServiceRegistrar, srv DictionaryServer) {
	s.RegisterService(&Dictionary_ServiceDesc, srv)
}

func _Dictionary_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestAdd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DictionaryServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exam.Dictionary/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DictionaryServer).Add(ctx, req.(*RequestAdd))
	}
	return interceptor(ctx, in, info, handler)
}

func _Dictionary_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DictionaryServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/exam.Dictionary/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DictionaryServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Dictionary_ServiceDesc is the grpc.ServiceDesc for Dictionary service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Dictionary_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "exam.Dictionary",
	HandlerType: (*DictionaryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Dictionary_Add_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _Dictionary_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/interface.proto",
}
