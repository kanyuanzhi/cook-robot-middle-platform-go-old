// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.3
// source: update.proto

package v1

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

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version      string `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	MachineModel string `protobuf:"bytes,2,opt,name=machineModel,proto3" json:"machineModel,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_update_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_update_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_update_proto_rawDescGZIP(), []int{0}
}

func (x *CheckRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *CheckRequest) GetMachineModel() string {
	if x != nil {
		return x.MachineModel
	}
	return ""
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsLatest      bool     `protobuf:"varint,1,opt,name=isLatest,proto3" json:"isLatest,omitempty"`          //是否为最新版
	LatestVersion string   `protobuf:"bytes,2,opt,name=latestVersion,proto3" json:"latestVersion,omitempty"` //最新版的版本号
	HasFile       bool     `protobuf:"varint,3,opt,name=hasFile,proto3" json:"hasFile,omitempty"`            //是否有可供更新的压缩包
	FilePath      []string `protobuf:"bytes,4,rep,name=filePath,proto3" json:"filePath,omitempty"`           //压缩包下载路径
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_update_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_update_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_update_proto_rawDescGZIP(), []int{1}
}

func (x *CheckResponse) GetIsLatest() bool {
	if x != nil {
		return x.IsLatest
	}
	return false
}

func (x *CheckResponse) GetLatestVersion() string {
	if x != nil {
		return x.LatestVersion
	}
	return ""
}

func (x *CheckResponse) GetHasFile() bool {
	if x != nil {
		return x.HasFile
	}
	return false
}

func (x *CheckResponse) GetFilePath() []string {
	if x != nil {
		return x.FilePath
	}
	return nil
}

var File_update_proto protoreflect.FileDescriptor

var file_update_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x76, 0x31, 0x22, 0x4c, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c,
	0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x6d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x6c,
	0x22, 0x87, 0x01, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x24,
	0x0a, 0x0d, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x68, 0x61, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x68, 0x61, 0x73, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x32, 0x38, 0x0a, 0x06, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x05, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x10, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_update_proto_rawDescOnce sync.Once
	file_update_proto_rawDescData = file_update_proto_rawDesc
)

func file_update_proto_rawDescGZIP() []byte {
	file_update_proto_rawDescOnce.Do(func() {
		file_update_proto_rawDescData = protoimpl.X.CompressGZIP(file_update_proto_rawDescData)
	})
	return file_update_proto_rawDescData
}

var file_update_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_update_proto_goTypes = []interface{}{
	(*CheckRequest)(nil),  // 0: v1.CheckRequest
	(*CheckResponse)(nil), // 1: v1.CheckResponse
}
var file_update_proto_depIdxs = []int32{
	0, // 0: v1.Update.Check:input_type -> v1.CheckRequest
	1, // 1: v1.Update.Check:output_type -> v1.CheckResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_update_proto_init() }
func file_update_proto_init() {
	if File_update_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_update_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_update_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
			RawDescriptor: file_update_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_update_proto_goTypes,
		DependencyIndexes: file_update_proto_depIdxs,
		MessageInfos:      file_update_proto_msgTypes,
	}.Build()
	File_update_proto = out.File
	file_update_proto_rawDesc = nil
	file_update_proto_goTypes = nil
	file_update_proto_depIdxs = nil
}
