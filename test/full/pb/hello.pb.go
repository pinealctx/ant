// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: test/full/pb/hello.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Hello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Halo string `protobuf:"bytes,1,opt,name=halo,proto3" json:"halo,omitempty"`
}

func (x *Hello) Reset() {
	*x = Hello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_full_pb_hello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hello) ProtoMessage() {}

func (x *Hello) ProtoReflect() protoreflect.Message {
	mi := &file_test_full_pb_hello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hello.ProtoReflect.Descriptor instead.
func (*Hello) Descriptor() ([]byte, []int) {
	return file_test_full_pb_hello_proto_rawDescGZIP(), []int{0}
}

func (x *Hello) GetHalo() string {
	if x != nil {
		return x.Halo
	}
	return ""
}

var File_test_full_pb_hello_proto protoreflect.FileDescriptor

var file_test_full_pb_hello_proto_rawDesc = []byte{
	0x0a, 0x18, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x66, 0x75, 0x6c, 0x6c, 0x2f, 0x70, 0x62, 0x2f, 0x68,
	0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1b, 0x0a, 0x05, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x6c, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x6c, 0x6f, 0x32, 0xb5, 0x01, 0x0a, 0x07, 0x45, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x12, 0x1c, 0x0a, 0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x09, 0x2e, 0x70,
	0x62, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x1a, 0x09, 0x2e, 0x70, 0x62, 0x2e, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x12, 0x29, 0x0a, 0x04, 0x50, 0x72, 0x6f, 0x63, 0x12, 0x09, 0x2e, 0x70, 0x62, 0x2e,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x28, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x09, 0x2e, 0x70,
	0x62, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x37, 0x0a, 0x05, 0x54, 0x6f, 0x75, 0x63, 0x68,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70,
	0x69, 0x6e, 0x65, 0x61, 0x6c, 0x63, 0x74, 0x78, 0x2f, 0x61, 0x6e, 0x74, 0x2f, 0x74, 0x65, 0x73,
	0x74, 0x2f, 0x66, 0x75, 0x6c, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_test_full_pb_hello_proto_rawDescOnce sync.Once
	file_test_full_pb_hello_proto_rawDescData = file_test_full_pb_hello_proto_rawDesc
)

func file_test_full_pb_hello_proto_rawDescGZIP() []byte {
	file_test_full_pb_hello_proto_rawDescOnce.Do(func() {
		file_test_full_pb_hello_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_full_pb_hello_proto_rawDescData)
	})
	return file_test_full_pb_hello_proto_rawDescData
}

var file_test_full_pb_hello_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_full_pb_hello_proto_goTypes = []interface{}{
	(*Hello)(nil),         // 0: pb.Hello
	(*emptypb.Empty)(nil), // 1: google.protobuf.Empty
}
var file_test_full_pb_hello_proto_depIdxs = []int32{
	0, // 0: pb.Example.Echo:input_type -> pb.Hello
	0, // 1: pb.Example.Proc:input_type -> pb.Hello
	1, // 2: pb.Example.Get:input_type -> google.protobuf.Empty
	1, // 3: pb.Example.Touch:input_type -> google.protobuf.Empty
	0, // 4: pb.Example.Echo:output_type -> pb.Hello
	1, // 5: pb.Example.Proc:output_type -> google.protobuf.Empty
	0, // 6: pb.Example.Get:output_type -> pb.Hello
	1, // 7: pb.Example.Touch:output_type -> google.protobuf.Empty
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_full_pb_hello_proto_init() }
func file_test_full_pb_hello_proto_init() {
	if File_test_full_pb_hello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_full_pb_hello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hello); i {
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
			RawDescriptor: file_test_full_pb_hello_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_test_full_pb_hello_proto_goTypes,
		DependencyIndexes: file_test_full_pb_hello_proto_depIdxs,
		MessageInfos:      file_test_full_pb_hello_proto_msgTypes,
	}.Build()
	File_test_full_pb_hello_proto = out.File
	file_test_full_pb_hello_proto_rawDesc = nil
	file_test_full_pb_hello_proto_goTypes = nil
	file_test_full_pb_hello_proto_depIdxs = nil
}
