// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gservice.proto

/*
Package gservice is a generated protocol buffer package.

It is generated from these files:
	gservice.proto

It has these top-level messages:
	ViewChangeRequest
	ViewChangeResponse
	View
	ServerNode
*/
package gservice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type ViewChangeRequest_ViewChangeType int32

const (
	ViewChangeRequest_ADD_NODE    ViewChangeRequest_ViewChangeType = 0
	ViewChangeRequest_REMOVE_NODE ViewChangeRequest_ViewChangeType = 1
)

var ViewChangeRequest_ViewChangeType_name = map[int32]string{
	0: "ADD_NODE",
	1: "REMOVE_NODE",
}
var ViewChangeRequest_ViewChangeType_value = map[string]int32{
	"ADD_NODE":    0,
	"REMOVE_NODE": 1,
}

func (x ViewChangeRequest_ViewChangeType) String() string {
	return proto.EnumName(ViewChangeRequest_ViewChangeType_name, int32(x))
}
func (ViewChangeRequest_ViewChangeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor0, []int{0, 0}
}

type ViewChangeRequest struct {
	RequestID   int64                            `protobuf:"varint,1,opt,name=RequestID" json:"RequestID,omitempty"`
	ServerNode  *ServerNode                      `protobuf:"bytes,2,opt,name=ServerNode" json:"ServerNode,omitempty"`
	Type        ViewChangeRequest_ViewChangeType `protobuf:"varint,3,opt,name=Type,enum=ViewChangeRequest_ViewChangeType" json:"Type,omitempty"`
	CurrentView []*View                          `protobuf:"bytes,4,rep,name=CurrentView" json:"CurrentView,omitempty"`
}

func (m *ViewChangeRequest) Reset()                    { *m = ViewChangeRequest{} }
func (m *ViewChangeRequest) String() string            { return proto.CompactTextString(m) }
func (*ViewChangeRequest) ProtoMessage()               {}
func (*ViewChangeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ViewChangeRequest) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ViewChangeRequest) GetServerNode() *ServerNode {
	if m != nil {
		return m.ServerNode
	}
	return nil
}

func (m *ViewChangeRequest) GetType() ViewChangeRequest_ViewChangeType {
	if m != nil {
		return m.Type
	}
	return ViewChangeRequest_ADD_NODE
}

func (m *ViewChangeRequest) GetCurrentView() []*View {
	if m != nil {
		return m.CurrentView
	}
	return nil
}

type ViewChangeResponse struct {
	RequestID   int64   `protobuf:"varint,1,opt,name=RequestID" json:"RequestID,omitempty"`
	CurrentView []*View `protobuf:"bytes,2,rep,name=currentView" json:"currentView,omitempty"`
	Status      string  `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
}

func (m *ViewChangeResponse) Reset()                    { *m = ViewChangeResponse{} }
func (m *ViewChangeResponse) String() string            { return proto.CompactTextString(m) }
func (*ViewChangeResponse) ProtoMessage()               {}
func (*ViewChangeResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ViewChangeResponse) GetRequestID() int64 {
	if m != nil {
		return m.RequestID
	}
	return 0
}

func (m *ViewChangeResponse) GetCurrentView() []*View {
	if m != nil {
		return m.CurrentView
	}
	return nil
}

func (m *ViewChangeResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type View struct {
	CurrentPartition []*ServerNode `protobuf:"bytes,1,rep,name=currentPartition" json:"currentPartition,omitempty"`
}

func (m *View) Reset()                    { *m = View{} }
func (m *View) String() string            { return proto.CompactTextString(m) }
func (*View) ProtoMessage()               {}
func (*View) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *View) GetCurrentPartition() []*ServerNode {
	if m != nil {
		return m.CurrentPartition
	}
	return nil
}

type ServerNode struct {
	IP   string `protobuf:"bytes,1,opt,name=IP" json:"IP,omitempty"`
	Port string `protobuf:"bytes,2,opt,name=Port" json:"Port,omitempty"`
}

func (m *ServerNode) Reset()                    { *m = ServerNode{} }
func (m *ServerNode) String() string            { return proto.CompactTextString(m) }
func (*ServerNode) ProtoMessage()               {}
func (*ServerNode) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ServerNode) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *ServerNode) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

func init() {
	proto.RegisterType((*ViewChangeRequest)(nil), "ViewChangeRequest")
	proto.RegisterType((*ViewChangeResponse)(nil), "ViewChangeResponse")
	proto.RegisterType((*View)(nil), "View")
	proto.RegisterType((*ServerNode)(nil), "ServerNode")
	proto.RegisterEnum("ViewChangeRequest_ViewChangeType", ViewChangeRequest_ViewChangeType_name, ViewChangeRequest_ViewChangeType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Store service

type StoreClient interface {
	// Sends a greeting
	AddServerNode(ctx context.Context, in *ViewChangeRequest, opts ...grpc.CallOption) (*ViewChangeResponse, error)
	RemoveServerNode(ctx context.Context, in *ViewChangeRequest, opts ...grpc.CallOption) (*ViewChangeResponse, error)
}

type storeClient struct {
	cc *grpc.ClientConn
}

func NewStoreClient(cc *grpc.ClientConn) StoreClient {
	return &storeClient{cc}
}

func (c *storeClient) AddServerNode(ctx context.Context, in *ViewChangeRequest, opts ...grpc.CallOption) (*ViewChangeResponse, error) {
	out := new(ViewChangeResponse)
	err := grpc.Invoke(ctx, "/Store/AddServerNode", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storeClient) RemoveServerNode(ctx context.Context, in *ViewChangeRequest, opts ...grpc.CallOption) (*ViewChangeResponse, error) {
	out := new(ViewChangeResponse)
	err := grpc.Invoke(ctx, "/Store/RemoveServerNode", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Store service

type StoreServer interface {
	// Sends a greeting
	AddServerNode(context.Context, *ViewChangeRequest) (*ViewChangeResponse, error)
	RemoveServerNode(context.Context, *ViewChangeRequest) (*ViewChangeResponse, error)
}

func RegisterStoreServer(s *grpc.Server, srv StoreServer) {
	s.RegisterService(&_Store_serviceDesc, srv)
}

func _Store_AddServerNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ViewChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).AddServerNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Store/AddServerNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).AddServerNode(ctx, req.(*ViewChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Store_RemoveServerNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ViewChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StoreServer).RemoveServerNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Store/RemoveServerNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StoreServer).RemoveServerNode(ctx, req.(*ViewChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Store_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Store",
	HandlerType: (*StoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddServerNode",
			Handler:    _Store_AddServerNode_Handler,
		},
		{
			MethodName: "RemoveServerNode",
			Handler:    _Store_RemoveServerNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gservice.proto",
}

func init() { proto.RegisterFile("gservice.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 353 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0xf5, 0xca, 0xb2, 0xa9, 0x46, 0xad, 0xea, 0x6e, 0xa1, 0x08, 0xe3, 0x83, 0xaa, 0x8b, 0x05,
	0x05, 0xb9, 0xa8, 0x94, 0x42, 0x21, 0x04, 0x63, 0xf9, 0xe0, 0x43, 0x6c, 0xb1, 0x0e, 0xbe, 0x06,
	0xc5, 0x1e, 0x14, 0x41, 0xa2, 0x55, 0x76, 0xd7, 0x0e, 0xb9, 0xe6, 0x87, 0xf3, 0x0b, 0xc1, 0x6b,
	0x83, 0xd7, 0xe8, 0x10, 0xc8, 0x6d, 0xf4, 0xde, 0x9b, 0xd1, 0x9b, 0x79, 0x0b, 0x5e, 0x21, 0x51,
	0xec, 0xca, 0x35, 0xc6, 0xb5, 0xe0, 0x8a, 0xf7, 0x07, 0x05, 0xe7, 0xc5, 0x3d, 0x8e, 0xf2, 0xba,
	0x1c, 0xe5, 0x55, 0xc5, 0x55, 0xae, 0x4a, 0x5e, 0xc9, 0x03, 0x1b, 0xbe, 0x12, 0xf8, 0xb6, 0x2a,
	0xf1, 0x69, 0x72, 0x97, 0x57, 0x05, 0x32, 0x7c, 0xdc, 0xa2, 0x54, 0x74, 0x00, 0xce, 0xb1, 0x9c,
	0xa5, 0x3e, 0x09, 0x48, 0xd4, 0x66, 0x27, 0x80, 0xfe, 0x02, 0x58, 0xa2, 0xd8, 0xa1, 0x98, 0xf3,
	0x0d, 0xfa, 0x56, 0x40, 0x22, 0x37, 0x71, 0xe3, 0x13, 0xc4, 0x0c, 0x9a, 0xfe, 0x05, 0xfb, 0xfa,
	0xb9, 0x46, 0xbf, 0x1d, 0x90, 0xc8, 0x4b, 0x7e, 0xc6, 0x8d, 0x9f, 0x19, 0xc8, 0x5e, 0xc8, 0xb4,
	0x9c, 0x0e, 0xc1, 0x9d, 0x6c, 0x85, 0xc0, 0x4a, 0xed, 0x69, 0xdf, 0x0e, 0xda, 0x91, 0x9b, 0x74,
	0xb4, 0x96, 0x99, 0x4c, 0x38, 0x02, 0xef, 0x7c, 0x00, 0xfd, 0x0c, 0x9f, 0xc6, 0x69, 0x7a, 0x33,
	0x5f, 0xa4, 0xd3, 0x5e, 0x8b, 0x7e, 0x05, 0x97, 0x4d, 0xaf, 0x16, 0xab, 0xe9, 0x01, 0x20, 0xa1,
	0x04, 0x6a, 0x7a, 0x90, 0x35, 0xaf, 0x24, 0xbe, 0xb3, 0xf1, 0x10, 0xdc, 0xb5, 0xe1, 0xc6, 0x3a,
	0x73, 0x63, 0x30, 0xf4, 0x07, 0x74, 0xa5, 0xca, 0xd5, 0x56, 0xea, 0x7d, 0x1d, 0x76, 0xfc, 0x0a,
	0x2f, 0xc1, 0xd6, 0xfc, 0x3f, 0xe8, 0x1d, 0xe5, 0x59, 0x2e, 0x54, 0xb9, 0x4f, 0xc2, 0x27, 0x7a,
	0xda, 0xd9, 0x01, 0x1b, 0xa2, 0xf0, 0xb7, 0x79, 0x73, 0xea, 0x81, 0x35, 0xcb, 0xb4, 0x4d, 0x87,
	0x59, 0xb3, 0x8c, 0x52, 0xb0, 0x33, 0x2e, 0x94, 0xce, 0xc2, 0x61, 0xba, 0x4e, 0x5e, 0x08, 0x74,
	0x96, 0x8a, 0x0b, 0xa4, 0xff, 0xe1, 0xcb, 0x78, 0xb3, 0x31, 0xda, 0x69, 0x33, 0x85, 0xfe, 0xf7,
	0xb8, 0x79, 0x95, 0xb0, 0x45, 0x2f, 0xa0, 0xc7, 0xf0, 0x81, 0xef, 0xf0, 0x43, 0xed, 0xb7, 0x5d,
	0xfd, 0xca, 0xfe, 0xbc, 0x05, 0x00, 0x00, 0xff, 0xff, 0xc8, 0xd0, 0xc1, 0x77, 0x95, 0x02, 0x00,
	0x00,
}
