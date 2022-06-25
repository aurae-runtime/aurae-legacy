//===========================================================================*\
//           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
//                                                                           *
//                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
//                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
//                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
//                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
//                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
//                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
//                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
//                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
//                                                                           *
//                       This machine kills fascists.                        *
//                                                                           *
//\*===========================================================================

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: rpc/aurae.proto

package rpc

import (
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

type SetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Val string `protobuf:"bytes,2,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *SetReq) Reset() {
	*x = SetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_aurae_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetReq) ProtoMessage() {}

func (x *SetReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_aurae_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetReq.ProtoReflect.Descriptor instead.
func (*SetReq) Descriptor() ([]byte, []int) {
	return file_rpc_aurae_proto_rawDescGZIP(), []int{0}
}

func (x *SetReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetReq) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

type SetResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *SetResp) Reset() {
	*x = SetResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_aurae_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResp) ProtoMessage() {}

func (x *SetResp) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_aurae_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResp.ProtoReflect.Descriptor instead.
func (*SetResp) Descriptor() ([]byte, []int) {
	return file_rpc_aurae_proto_rawDescGZIP(), []int{1}
}

func (x *SetResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type GetReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetReq) Reset() {
	*x = GetReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_aurae_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReq) ProtoMessage() {}

func (x *GetReq) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_aurae_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReq.ProtoReflect.Descriptor instead.
func (*GetReq) Descriptor() ([]byte, []int) {
	return file_rpc_aurae_proto_rawDescGZIP(), []int{2}
}

func (x *GetReq) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type GetResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Val  string `protobuf:"bytes,1,opt,name=val,proto3" json:"val,omitempty"`
	Code int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *GetResp) Reset() {
	*x = GetResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_aurae_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResp) ProtoMessage() {}

func (x *GetResp) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_aurae_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetResp.ProtoReflect.Descriptor instead.
func (*GetResp) Descriptor() ([]byte, []int) {
	return file_rpc_aurae_proto_rawDescGZIP(), []int{3}
}

func (x *GetResp) GetVal() string {
	if x != nil {
		return x.Val
	}
	return ""
}

func (x *GetResp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

var File_rpc_aurae_proto protoreflect.FileDescriptor

var file_rpc_aurae_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x70, 0x63, 0x2f, 0x61, 0x75, 0x72, 0x61, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x61, 0x75, 0x72, 0x61, 0x65, 0x22, 0x2c, 0x0a, 0x06, 0x53, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22, 0x1d, 0x0a, 0x07, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x1a, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x22, 0x2f, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x10, 0x0a, 0x03,
	0x76, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x32, 0x5d, 0x0a, 0x0b, 0x43, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x26, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x0d, 0x2e, 0x61, 0x75, 0x72, 0x61, 0x65,
	0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x0e, 0x2e, 0x61, 0x75, 0x72, 0x61, 0x65, 0x2e,
	0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x26, 0x0a, 0x03, 0x47, 0x65, 0x74,
	0x12, 0x0d, 0x2e, 0x61, 0x75, 0x72, 0x61, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x1a,
	0x0e, 0x2e, 0x61, 0x75, 0x72, 0x61, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x00, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6b, 0x72, 0x69, 0x73, 0x2d, 0x6e, 0x6f, 0x76, 0x61, 0x2f, 0x61, 0x75, 0x72, 0x61, 0x65, 0x2f,
	0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_aurae_proto_rawDescOnce sync.Once
	file_rpc_aurae_proto_rawDescData = file_rpc_aurae_proto_rawDesc
)

func file_rpc_aurae_proto_rawDescGZIP() []byte {
	file_rpc_aurae_proto_rawDescOnce.Do(func() {
		file_rpc_aurae_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_aurae_proto_rawDescData)
	})
	return file_rpc_aurae_proto_rawDescData
}

var file_rpc_aurae_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_aurae_proto_goTypes = []interface{}{
	(*SetReq)(nil),  // 0: aurae.SetReq
	(*SetResp)(nil), // 1: aurae.SetResp
	(*GetReq)(nil),  // 2: aurae.GetReq
	(*GetResp)(nil), // 3: aurae.GetResp
}
var file_rpc_aurae_proto_depIdxs = []int32{
	0, // 0: aurae.CoreService.Set:input_type -> aurae.SetReq
	2, // 1: aurae.CoreService.Get:input_type -> aurae.GetReq
	1, // 2: aurae.CoreService.Set:output_type -> aurae.SetResp
	3, // 3: aurae.CoreService.Get:output_type -> aurae.GetResp
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rpc_aurae_proto_init() }
func file_rpc_aurae_proto_init() {
	if File_rpc_aurae_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_aurae_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetReq); i {
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
		file_rpc_aurae_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetResp); i {
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
		file_rpc_aurae_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReq); i {
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
		file_rpc_aurae_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResp); i {
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
			RawDescriptor: file_rpc_aurae_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_aurae_proto_goTypes,
		DependencyIndexes: file_rpc_aurae_proto_depIdxs,
		MessageInfos:      file_rpc_aurae_proto_msgTypes,
	}.Build()
	File_rpc_aurae_proto = out.File
	file_rpc_aurae_proto_rawDesc = nil
	file_rpc_aurae_proto_goTypes = nil
	file_rpc_aurae_proto_depIdxs = nil
}
