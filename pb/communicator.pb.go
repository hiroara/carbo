// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: communicator.proto

package pb

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

type FillBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Limit int32  `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *FillBatchRequest) Reset() {
	*x = FillBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_communicator_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FillBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FillBatchRequest) ProtoMessage() {}

func (x *FillBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_communicator_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FillBatchRequest.ProtoReflect.Descriptor instead.
func (*FillBatchRequest) Descriptor() ([]byte, []int) {
	return file_communicator_proto_rawDescGZIP(), []int{0}
}

func (x *FillBatchRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *FillBatchRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type FillBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	End bool `protobuf:"varint,1,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *FillBatchResponse) Reset() {
	*x = FillBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_communicator_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FillBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FillBatchResponse) ProtoMessage() {}

func (x *FillBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_communicator_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FillBatchResponse.ProtoReflect.Descriptor instead.
func (*FillBatchResponse) Descriptor() ([]byte, []int) {
	return file_communicator_proto_rawDescGZIP(), []int{1}
}

func (x *FillBatchResponse) GetEnd() bool {
	if x != nil {
		return x.End
	}
	return false
}

type GetBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetBatchRequest) Reset() {
	*x = GetBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_communicator_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBatchRequest) ProtoMessage() {}

func (x *GetBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_communicator_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBatchRequest.ProtoReflect.Descriptor instead.
func (*GetBatchRequest) Descriptor() ([]byte, []int) {
	return file_communicator_proto_rawDescGZIP(), []int{2}
}

type GetBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token    string     `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Messages []*Message `protobuf:"bytes,2,rep,name=messages,proto3" json:"messages,omitempty"`
}

func (x *GetBatchResponse) Reset() {
	*x = GetBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_communicator_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBatchResponse) ProtoMessage() {}

func (x *GetBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_communicator_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBatchResponse.ProtoReflect.Descriptor instead.
func (*GetBatchResponse) Descriptor() ([]byte, []int) {
	return file_communicator_proto_rawDescGZIP(), []int{3}
}

func (x *GetBatchResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *GetBatchResponse) GetMessages() []*Message {
	if x != nil {
		return x.Messages
	}
	return nil
}

var File_communicator_proto protoreflect.FileDescriptor

var file_communicator_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x22, 0x25, 0x0a, 0x11, 0x46, 0x69, 0x6c, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65,
	0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x4e, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x24, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x32, 0x77, 0x0a,
	0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x34, 0x0a,
	0x09, 0x46, 0x69, 0x6c, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x11, 0x2e, 0x46, 0x69, 0x6c,
	0x6c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e,
	0x46, 0x69, 0x6c, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x31, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12,
	0x10, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1d, 0x5a, 0x1b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x69, 0x72, 0x6f, 0x61, 0x72, 0x61, 0x2f, 0x63, 0x61, 0x72,
	0x62, 0x6f, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_communicator_proto_rawDescOnce sync.Once
	file_communicator_proto_rawDescData = file_communicator_proto_rawDesc
)

func file_communicator_proto_rawDescGZIP() []byte {
	file_communicator_proto_rawDescOnce.Do(func() {
		file_communicator_proto_rawDescData = protoimpl.X.CompressGZIP(file_communicator_proto_rawDescData)
	})
	return file_communicator_proto_rawDescData
}

var file_communicator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_communicator_proto_goTypes = []interface{}{
	(*FillBatchRequest)(nil),  // 0: FillBatchRequest
	(*FillBatchResponse)(nil), // 1: FillBatchResponse
	(*GetBatchRequest)(nil),   // 2: GetBatchRequest
	(*GetBatchResponse)(nil),  // 3: GetBatchResponse
	(*Message)(nil),           // 4: Message
}
var file_communicator_proto_depIdxs = []int32{
	4, // 0: GetBatchResponse.messages:type_name -> Message
	0, // 1: Communicator.FillBatch:input_type -> FillBatchRequest
	2, // 2: Communicator.GetBatch:input_type -> GetBatchRequest
	1, // 3: Communicator.FillBatch:output_type -> FillBatchResponse
	3, // 4: Communicator.GetBatch:output_type -> GetBatchResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_communicator_proto_init() }
func file_communicator_proto_init() {
	if File_communicator_proto != nil {
		return
	}
	file_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_communicator_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FillBatchRequest); i {
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
		file_communicator_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FillBatchResponse); i {
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
		file_communicator_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBatchRequest); i {
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
		file_communicator_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBatchResponse); i {
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
			RawDescriptor: file_communicator_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_communicator_proto_goTypes,
		DependencyIndexes: file_communicator_proto_depIdxs,
		MessageInfos:      file_communicator_proto_msgTypes,
	}.Build()
	File_communicator_proto = out.File
	file_communicator_proto_rawDesc = nil
	file_communicator_proto_goTypes = nil
	file_communicator_proto_depIdxs = nil
}
