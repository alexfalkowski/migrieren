// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: migrieren/v1/service.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Migration for a specific database and version with logs.
type Migration struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Database      string                 `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	Logs          []string               `protobuf:"bytes,3,rep,name=logs,proto3" json:"logs,omitempty"`
	unknownFields protoimpl.UnknownFields
	Version       uint64 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	sizeCache     protoimpl.SizeCache
}

func (x *Migration) Reset() {
	*x = Migration{}
	mi := &file_migrieren_v1_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Migration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Migration) ProtoMessage() {}

func (x *Migration) ProtoReflect() protoreflect.Message {
	mi := &file_migrieren_v1_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Migration.ProtoReflect.Descriptor instead.
func (*Migration) Descriptor() ([]byte, []int) {
	return file_migrieren_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *Migration) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

func (x *Migration) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Migration) GetLogs() []string {
	if x != nil {
		return x.Logs
	}
	return nil
}

// MigrateRequest for a specific database and version.
type MigrateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Database      string                 `protobuf:"bytes,1,opt,name=database,proto3" json:"database,omitempty"`
	unknownFields protoimpl.UnknownFields
	Version       uint64 `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	sizeCache     protoimpl.SizeCache
}

func (x *MigrateRequest) Reset() {
	*x = MigrateRequest{}
	mi := &file_migrieren_v1_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MigrateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrateRequest) ProtoMessage() {}

func (x *MigrateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_migrieren_v1_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrateRequest.ProtoReflect.Descriptor instead.
func (*MigrateRequest) Descriptor() ([]byte, []int) {
	return file_migrieren_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *MigrateRequest) GetDatabase() string {
	if x != nil {
		return x.Database
	}
	return ""
}

func (x *MigrateRequest) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

// MigrateResponse for a specific database and version.
type MigrateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Meta          map[string]string      `protobuf:"bytes,1,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Migration     *Migration             `protobuf:"bytes,2,opt,name=migration,proto3" json:"migration,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MigrateResponse) Reset() {
	*x = MigrateResponse{}
	mi := &file_migrieren_v1_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MigrateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MigrateResponse) ProtoMessage() {}

func (x *MigrateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_migrieren_v1_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MigrateResponse.ProtoReflect.Descriptor instead.
func (*MigrateResponse) Descriptor() ([]byte, []int) {
	return file_migrieren_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *MigrateResponse) GetMeta() map[string]string {
	if x != nil {
		return x.Meta
	}
	return nil
}

func (x *MigrateResponse) GetMigration() *Migration {
	if x != nil {
		return x.Migration
	}
	return nil
}

var File_migrieren_v1_service_proto protoreflect.FileDescriptor

var file_migrieren_v1_service_proto_rawDesc = string([]byte{
	0x0a, 0x1a, 0x6d, 0x69, 0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6d, 0x69,
	0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2e, 0x76, 0x31, 0x22, 0x55, 0x0a, 0x09, 0x4d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62,
	0x61, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x6f, 0x67,
	0x73, 0x22, 0x46, 0x0a, 0x0e, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xbe, 0x01, 0x0a, 0x0f, 0x4d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a,
	0x04, 0x6d, 0x65, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6d, 0x69,
	0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69, 0x67, 0x72, 0x61,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x12, 0x35, 0x0a, 0x09, 0x6d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e,
	0x6d, 0x69, 0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69, 0x67,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x1a, 0x37, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x53, 0x0a, 0x07, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x07, 0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65,
	0x12, 0x1c, 0x2e, 0x6d, 0x69, 0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x69, 0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d,
	0x2e, 0x6d, 0x69, 0x67, 0x72, 0x69, 0x65, 0x72, 0x65, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x69,
	0x67, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42,
	0x45, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c,
	0x65, 0x78, 0x66, 0x61, 0x6c, 0x6b, 0x6f, 0x77, 0x73, 0x6b, 0x69, 0x2f, 0x6d, 0x69, 0x67, 0x72,
	0x69, 0x65, 0x72, 0x65, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6d, 0x69, 0x67, 0x72, 0x69, 0x65,
	0x72, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0xea, 0x02, 0x0d, 0x4d, 0x69, 0x67, 0x72, 0x69, 0x65, 0x72,
	0x65, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_migrieren_v1_service_proto_rawDescOnce sync.Once
	file_migrieren_v1_service_proto_rawDescData []byte
)

func file_migrieren_v1_service_proto_rawDescGZIP() []byte {
	file_migrieren_v1_service_proto_rawDescOnce.Do(func() {
		file_migrieren_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_migrieren_v1_service_proto_rawDesc), len(file_migrieren_v1_service_proto_rawDesc)))
	})
	return file_migrieren_v1_service_proto_rawDescData
}

var file_migrieren_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_migrieren_v1_service_proto_goTypes = []any{
	(*Migration)(nil),       // 0: migrieren.v1.Migration
	(*MigrateRequest)(nil),  // 1: migrieren.v1.MigrateRequest
	(*MigrateResponse)(nil), // 2: migrieren.v1.MigrateResponse
	nil,                     // 3: migrieren.v1.MigrateResponse.MetaEntry
}
var file_migrieren_v1_service_proto_depIdxs = []int32{
	3, // 0: migrieren.v1.MigrateResponse.meta:type_name -> migrieren.v1.MigrateResponse.MetaEntry
	0, // 1: migrieren.v1.MigrateResponse.migration:type_name -> migrieren.v1.Migration
	1, // 2: migrieren.v1.Service.Migrate:input_type -> migrieren.v1.MigrateRequest
	2, // 3: migrieren.v1.Service.Migrate:output_type -> migrieren.v1.MigrateResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_migrieren_v1_service_proto_init() }
func file_migrieren_v1_service_proto_init() {
	if File_migrieren_v1_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_migrieren_v1_service_proto_rawDesc), len(file_migrieren_v1_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_migrieren_v1_service_proto_goTypes,
		DependencyIndexes: file_migrieren_v1_service_proto_depIdxs,
		MessageInfos:      file_migrieren_v1_service_proto_msgTypes,
	}.Build()
	File_migrieren_v1_service_proto = out.File
	file_migrieren_v1_service_proto_goTypes = nil
	file_migrieren_v1_service_proto_depIdxs = nil
}
