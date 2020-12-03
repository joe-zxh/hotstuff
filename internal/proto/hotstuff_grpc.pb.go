// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// HotstuffClient is the client API for Hotstuff service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HotstuffClient interface {
	Propose(ctx context.Context, in *Block, opts ...grpc.CallOption) (*empty.Empty, error)
	Vote(ctx context.Context, in *PartialCert, opts ...grpc.CallOption) (*empty.Empty, error)
}

type hotstuffClient struct {
	cc grpc.ClientConnInterface
}

func NewHotstuffClient(cc grpc.ClientConnInterface) HotstuffClient {
	return &hotstuffClient{cc}
}

func (c *hotstuffClient) Propose(ctx context.Context, in *Block, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.Hotstuff/Propose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hotstuffClient) Vote(ctx context.Context, in *PartialCert, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/proto.Hotstuff/Vote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HotstuffServer is the server API for Hotstuff service.
// All implementations must embed UnimplementedHotstuffServer
// for forward compatibility
type HotstuffServer interface {
	Propose(context.Context, *Block) (*empty.Empty, error)
	Vote(context.Context, *PartialCert) (*empty.Empty, error)
	mustEmbedUnimplementedHotstuffServer()
}

// UnimplementedHotstuffServer must be embedded to have forward compatible implementations.
type UnimplementedHotstuffServer struct {
}

func (UnimplementedHotstuffServer) Propose(context.Context, *Block) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Propose not implemented")
}
func (UnimplementedHotstuffServer) Vote(context.Context, *PartialCert) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vote not implemented")
}
func (UnimplementedHotstuffServer) mustEmbedUnimplementedHotstuffServer() {}

// UnsafeHotstuffServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HotstuffServer will
// result in compilation errors.
type UnsafeHotstuffServer interface {
	mustEmbedUnimplementedHotstuffServer()
}

func RegisterHotstuffServer(s grpc.ServiceRegistrar, srv HotstuffServer) {
	s.RegisterService(&_Hotstuff_serviceDesc, srv)
}

func _Hotstuff_Propose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Block)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotstuffServer).Propose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Hotstuff/Propose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotstuffServer).Propose(ctx, req.(*Block))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hotstuff_Vote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PartialCert)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HotstuffServer).Vote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Hotstuff/Vote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HotstuffServer).Vote(ctx, req.(*PartialCert))
	}
	return interceptor(ctx, in, info, handler)
}

var _Hotstuff_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Hotstuff",
	HandlerType: (*HotstuffServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Propose",
			Handler:    _Hotstuff_Propose_Handler,
		},
		{
			MethodName: "Vote",
			Handler:    _Hotstuff_Vote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hotstuff.proto",
}