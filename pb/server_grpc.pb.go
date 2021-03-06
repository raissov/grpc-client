// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: client.proto

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

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	// client streaming
	MessageManyTimes(ctx context.Context, in *MessageManyTimesRequest, opts ...grpc.CallOption) (MessageService_MessageManyTimesClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) MessageManyTimes(ctx context.Context, in *MessageManyTimesRequest, opts ...grpc.CallOption) (MessageService_MessageManyTimesClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/pb.MessageService/MessageManyTimes", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceMessageManyTimesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_MessageManyTimesClient interface {
	Recv() (*MessageManyTimesResponse, error)
	grpc.ClientStream
}

type messageServiceMessageManyTimesClient struct {
	grpc.ClientStream
}

func (x *messageServiceMessageManyTimesClient) Recv() (*MessageManyTimesResponse, error) {
	m := new(MessageManyTimesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the client API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	// client streaming
	MessageManyTimes(*MessageManyTimesRequest, MessageService_MessageManyTimesServer) error
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) MessageManyTimes(*MessageManyTimesRequest, MessageService_MessageManyTimesServer) error {
	return status.Errorf(codes.Unimplemented, "method MessageManyTimes not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_MessageManyTimes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MessageManyTimesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).MessageManyTimes(m, &messageServiceMessageManyTimesServer{stream})
}

type MessageService_MessageManyTimesServer interface {
	Send(*MessageManyTimesResponse) error
	grpc.ServerStream
}

type messageServiceMessageManyTimesServer struct {
	grpc.ServerStream
}

func (x *messageServiceMessageManyTimesServer) Send(m *MessageManyTimesResponse) error {
	return x.ServerStream.SendMsg(m)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "MessageManyTimes",
			Handler:       _MessageService_MessageManyTimes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "client.proto",
}
