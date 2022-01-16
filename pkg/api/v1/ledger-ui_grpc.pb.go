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

// LedgerUIClient is the client API for LedgerUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LedgerUIClient interface {
	// ListEngineSpecs returns a list of Ledger Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (LedgerUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type ledgerUIClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerUIClient(cc grpc.ClientConnInterface) LedgerUIClient {
	return &ledgerUIClient{cc}
}

func (c *ledgerUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (LedgerUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &LedgerUI_ServiceDesc.Streams[0], "/v1.LedgerUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &ledgerUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type LedgerUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type ledgerUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *ledgerUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *ledgerUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.LedgerUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerUIServer is the server API for LedgerUI service.
// All implementations must embed UnimplementedLedgerUIServer
// for forward compatibility
type LedgerUIServer interface {
	// ListEngineSpecs returns a list of Ledger Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, LedgerUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedLedgerUIServer()
}

// UnimplementedLedgerUIServer must be embedded to have forward compatible implementations.
type UnimplementedLedgerUIServer struct {
}

func (UnimplementedLedgerUIServer) ListEngineSpecs(*ListEngineSpecsRequest, LedgerUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedLedgerUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedLedgerUIServer) mustEmbedUnimplementedLedgerUIServer() {}

// UnsafeLedgerUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LedgerUIServer will
// result in compilation errors.
type UnsafeLedgerUIServer interface {
	mustEmbedUnimplementedLedgerUIServer()
}

func RegisterLedgerUIServer(s grpc.ServiceRegistrar, srv LedgerUIServer) {
	s.RegisterService(&LedgerUI_ServiceDesc, srv)
}

func _LedgerUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LedgerUIServer).ListEngineSpecs(m, &ledgerUIListEngineSpecsServer{stream})
}

type LedgerUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type ledgerUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *ledgerUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _LedgerUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.LedgerUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LedgerUI_ServiceDesc is the grpc.ServiceDesc for LedgerUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LedgerUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.LedgerUI",
	HandlerType: (*LedgerUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _LedgerUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _LedgerUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "ledger-ui.proto",
}