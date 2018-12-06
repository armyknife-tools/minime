// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc_controller.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_23c2c7e42feab570, []int{0}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Empty)(nil), "proto.Empty")
}

func init() { proto.RegisterFile("grpc_controller.proto", fileDescriptor_23c2c7e42feab570) }

var fileDescriptor_23c2c7e42feab570 = []byte{
	// 97 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0x2f, 0x2a, 0x48,
	0x8e, 0x4f, 0xce, 0xcf, 0x2b, 0x29, 0xca, 0xcf, 0xc9, 0x49, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0xec, 0x5c, 0xac, 0xae, 0xb9, 0x05, 0x25, 0x95, 0x46, 0x16,
	0x5c, 0x7c, 0xee, 0x41, 0x01, 0xce, 0xce, 0x70, 0x75, 0x42, 0x6a, 0x5c, 0x1c, 0xc1, 0x19, 0xa5,
	0x25, 0x29, 0xf9, 0xe5, 0x79, 0x42, 0x3c, 0x10, 0x5d, 0x7a, 0x60, 0xb5, 0x52, 0x28, 0xbc, 0x24,
	0x36, 0x30, 0xc7, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x69, 0xa1, 0xad, 0x79, 0x69, 0x00, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GRPCControllerClient is the client API for GRPCController service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GRPCControllerClient interface {
	Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type gRPCControllerClient struct {
	cc *grpc.ClientConn
}

func NewGRPCControllerClient(cc *grpc.ClientConn) GRPCControllerClient {
	return &gRPCControllerClient{cc}
}

func (c *gRPCControllerClient) Shutdown(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.GRPCController/Shutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GRPCControllerServer is the server API for GRPCController service.
type GRPCControllerServer interface {
	Shutdown(context.Context, *Empty) (*Empty, error)
}

func RegisterGRPCControllerServer(s *grpc.Server, srv GRPCControllerServer) {
	s.RegisterService(&_GRPCController_serviceDesc, srv)
}

func _GRPCController_Shutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCControllerServer).Shutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.GRPCController/Shutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCControllerServer).Shutdown(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _GRPCController_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GRPCController",
	HandlerType: (*GRPCControllerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Shutdown",
			Handler:    _GRPCController_Shutdown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc_controller.proto",
}
