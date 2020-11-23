// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: networking/networking.proto

package networking

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Node_Status int32

const (
	Node_NODE_STATUS_UNSET  Node_Status = 0
	Node_NODE_STATUS_JOINED Node_Status = 1
	Node_NODE_STATUS_FAILED Node_Status = 2
)

// Enum value maps for Node_Status.
var (
	Node_Status_name = map[int32]string{
		0: "NODE_STATUS_UNSET",
		1: "NODE_STATUS_JOINED",
		2: "NODE_STATUS_FAILED",
	}
	Node_Status_value = map[string]int32{
		"NODE_STATUS_UNSET":  0,
		"NODE_STATUS_JOINED": 1,
		"NODE_STATUS_FAILED": 2,
	}
)

func (x Node_Status) Enum() *Node_Status {
	p := new(Node_Status)
	*p = x
	return p
}

func (x Node_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Node_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_networking_networking_proto_enumTypes[0].Descriptor()
}

func (Node_Status) Type() protoreflect.EnumType {
	return &file_networking_networking_proto_enumTypes[0]
}

func (x Node_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Node_Status.Descriptor instead.
func (Node_Status) EnumDescriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{0, 0}
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr   string      `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	Status Node_Status `protobuf:"varint,3,opt,name=status,proto3,enum=networking.Node_Status" json:"status,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{0}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Node) GetStatus() Node_Status {
	if x != nil {
		return x.Status
	}
	return Node_NODE_STATUS_UNSET
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{1}
}

func (x *RegisterRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RegisterRequest) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterResponse) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type EmptyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyRequest) Reset() {
	*x = EmptyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyRequest) ProtoMessage() {}

func (x *EmptyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyRequest.ProtoReflect.Descriptor instead.
func (*EmptyRequest) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{3}
}

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Node *Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{4}
}

func (x *PingRequest) GetNode() *Node {
	if x != nil {
		return x.Node
	}
	return nil
}

type PingReqResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *PingReqResponse) Reset() {
	*x = PingReqResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_networking_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingReqResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingReqResponse) ProtoMessage() {}

func (x *PingReqResponse) ProtoReflect() protoreflect.Message {
	mi := &file_networking_networking_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingReqResponse.ProtoReflect.Descriptor instead.
func (*PingReqResponse) Descriptor() ([]byte, []int) {
	return file_networking_networking_proto_rawDescGZIP(), []int{5}
}

func (x *PingReqResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_networking_networking_proto protoreflect.FileDescriptor

var file_networking_networking_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x22, 0xac, 0x01, 0x0a, 0x04, 0x4e, 0x6f,
	0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x4f, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x15, 0x0a, 0x11, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x55, 0x4e, 0x53, 0x45, 0x54, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x4e, 0x4f, 0x44, 0x45,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4a, 0x4f, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x16, 0x0a, 0x12, 0x4e, 0x4f, 0x44, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x02, 0x22, 0x35, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x61,
	0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x22,
	0x26, 0x0a, 0x10, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72, 0x22, 0x0e, 0x0a, 0x0c, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x33, 0x0a, 0x0b, 0x50, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e,
	0x67, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x22, 0x2b, 0x0a, 0x0f,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0xcd, 0x01, 0x0a, 0x11, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x3d, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x36,
	0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x18, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x10, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x41, 0x0a, 0x07, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x71, 0x12, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6e, 0x65, 0x6d, 0x61, 0x73, 0x2f,
	0x6c, 0x69, 0x62, 0x65, 0x72, 0x65, 0x63, 0x6f, 0x2f, 0x62, 0x75, 0x73, 0x69, 0x6e, 0x65, 0x73,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_networking_networking_proto_rawDescOnce sync.Once
	file_networking_networking_proto_rawDescData = file_networking_networking_proto_rawDesc
)

func file_networking_networking_proto_rawDescGZIP() []byte {
	file_networking_networking_proto_rawDescOnce.Do(func() {
		file_networking_networking_proto_rawDescData = protoimpl.X.CompressGZIP(file_networking_networking_proto_rawDescData)
	})
	return file_networking_networking_proto_rawDescData
}

var file_networking_networking_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_networking_networking_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_networking_networking_proto_goTypes = []interface{}{
	(Node_Status)(0),         // 0: networking.Node.Status
	(*Node)(nil),             // 1: networking.Node
	(*RegisterRequest)(nil),  // 2: networking.RegisterRequest
	(*RegisterResponse)(nil), // 3: networking.RegisterResponse
	(*EmptyRequest)(nil),     // 4: networking.EmptyRequest
	(*PingRequest)(nil),      // 5: networking.PingRequest
	(*PingReqResponse)(nil),  // 6: networking.PingReqResponse
}
var file_networking_networking_proto_depIdxs = []int32{
	0, // 0: networking.Node.status:type_name -> networking.Node.Status
	1, // 1: networking.PingRequest.node:type_name -> networking.Node
	2, // 2: networking.NetworkingService.Register:input_type -> networking.RegisterRequest
	4, // 3: networking.NetworkingService.Ping:input_type -> networking.EmptyRequest
	5, // 4: networking.NetworkingService.PingReq:input_type -> networking.PingRequest
	1, // 5: networking.NetworkingService.Register:output_type -> networking.Node
	1, // 6: networking.NetworkingService.Ping:output_type -> networking.Node
	6, // 7: networking.NetworkingService.PingReq:output_type -> networking.PingReqResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_networking_networking_proto_init() }
func file_networking_networking_proto_init() {
	if File_networking_networking_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_networking_networking_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_networking_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_networking_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_networking_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_networking_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_networking_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingReqResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_networking_networking_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_networking_networking_proto_goTypes,
		DependencyIndexes: file_networking_networking_proto_depIdxs,
		EnumInfos:         file_networking_networking_proto_enumTypes,
		MessageInfos:      file_networking_networking_proto_msgTypes,
	}.Build()
	File_networking_networking_proto = out.File
	file_networking_networking_proto_rawDesc = nil
	file_networking_networking_proto_goTypes = nil
	file_networking_networking_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NetworkingServiceClient is the client API for NetworkingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NetworkingServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (NetworkingService_RegisterClient, error)
	Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (NetworkingService_PingClient, error)
	PingReq(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReqResponse, error)
}

type networkingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNetworkingServiceClient(cc grpc.ClientConnInterface) NetworkingServiceClient {
	return &networkingServiceClient{cc}
}

func (c *networkingServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (NetworkingService_RegisterClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NetworkingService_serviceDesc.Streams[0], "/networking.NetworkingService/Register", opts...)
	if err != nil {
		return nil, err
	}
	x := &networkingServiceRegisterClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NetworkingService_RegisterClient interface {
	Recv() (*Node, error)
	grpc.ClientStream
}

type networkingServiceRegisterClient struct {
	grpc.ClientStream
}

func (x *networkingServiceRegisterClient) Recv() (*Node, error) {
	m := new(Node)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *networkingServiceClient) Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (NetworkingService_PingClient, error) {
	stream, err := c.cc.NewStream(ctx, &_NetworkingService_serviceDesc.Streams[1], "/networking.NetworkingService/Ping", opts...)
	if err != nil {
		return nil, err
	}
	x := &networkingServicePingClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type NetworkingService_PingClient interface {
	Recv() (*Node, error)
	grpc.ClientStream
}

type networkingServicePingClient struct {
	grpc.ClientStream
}

func (x *networkingServicePingClient) Recv() (*Node, error) {
	m := new(Node)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *networkingServiceClient) PingReq(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingReqResponse, error) {
	out := new(PingReqResponse)
	err := c.cc.Invoke(ctx, "/networking.NetworkingService/PingReq", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NetworkingServiceServer is the server API for NetworkingService service.
type NetworkingServiceServer interface {
	Register(*RegisterRequest, NetworkingService_RegisterServer) error
	Ping(*EmptyRequest, NetworkingService_PingServer) error
	PingReq(context.Context, *PingRequest) (*PingReqResponse, error)
}

// UnimplementedNetworkingServiceServer can be embedded to have forward compatible implementations.
type UnimplementedNetworkingServiceServer struct {
}

func (*UnimplementedNetworkingServiceServer) Register(*RegisterRequest, NetworkingService_RegisterServer) error {
	return status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (*UnimplementedNetworkingServiceServer) Ping(*EmptyRequest, NetworkingService_PingServer) error {
	return status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (*UnimplementedNetworkingServiceServer) PingReq(context.Context, *PingRequest) (*PingReqResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingReq not implemented")
}

func RegisterNetworkingServiceServer(s *grpc.Server, srv NetworkingServiceServer) {
	s.RegisterService(&_NetworkingService_serviceDesc, srv)
}

func _NetworkingService_Register_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RegisterRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NetworkingServiceServer).Register(m, &networkingServiceRegisterServer{stream})
}

type NetworkingService_RegisterServer interface {
	Send(*Node) error
	grpc.ServerStream
}

type networkingServiceRegisterServer struct {
	grpc.ServerStream
}

func (x *networkingServiceRegisterServer) Send(m *Node) error {
	return x.ServerStream.SendMsg(m)
}

func _NetworkingService_Ping_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NetworkingServiceServer).Ping(m, &networkingServicePingServer{stream})
}

type NetworkingService_PingServer interface {
	Send(*Node) error
	grpc.ServerStream
}

type networkingServicePingServer struct {
	grpc.ServerStream
}

func (x *networkingServicePingServer) Send(m *Node) error {
	return x.ServerStream.SendMsg(m)
}

func _NetworkingService_PingReq_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NetworkingServiceServer).PingReq(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/networking.NetworkingService/PingReq",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NetworkingServiceServer).PingReq(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _NetworkingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "networking.NetworkingService",
	HandlerType: (*NetworkingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PingReq",
			Handler:    _NetworkingService_PingReq_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Register",
			Handler:       _NetworkingService_Register_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Ping",
			Handler:       _NetworkingService_Ping_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "networking/networking.proto",
}