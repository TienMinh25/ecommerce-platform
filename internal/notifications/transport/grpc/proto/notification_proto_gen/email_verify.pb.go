// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: email_verify.proto

package notification_proto_gen

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

type TypeVerifyOTP int32

const (
	TypeVerifyOTP_EMAIL TypeVerifyOTP = 0
	TypeVerifyOTP_PHONE TypeVerifyOTP = 1
)

// Enum value maps for TypeVerifyOTP.
var (
	TypeVerifyOTP_name = map[int32]string{
		0: "EMAIL",
		1: "PHONE",
	}
	TypeVerifyOTP_value = map[string]int32{
		"EMAIL": 0,
		"PHONE": 1,
	}
)

func (x TypeVerifyOTP) Enum() *TypeVerifyOTP {
	p := new(TypeVerifyOTP)
	*p = x
	return p
}

func (x TypeVerifyOTP) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TypeVerifyOTP) Descriptor() protoreflect.EnumDescriptor {
	return file_email_verify_proto_enumTypes[0].Descriptor()
}

func (TypeVerifyOTP) Type() protoreflect.EnumType {
	return &file_email_verify_proto_enumTypes[0]
}

func (x TypeVerifyOTP) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TypeVerifyOTP.Descriptor instead.
func (TypeVerifyOTP) EnumDescriptor() ([]byte, []int) {
	return file_email_verify_proto_rawDescGZIP(), []int{0}
}

type PurposeOTP int32

const (
	PurposeOTP_EMAIL_VERIFICATION PurposeOTP = 0
	PurposeOTP_PASSWORD_RESET     PurposeOTP = 1
)

// Enum value maps for PurposeOTP.
var (
	PurposeOTP_name = map[int32]string{
		0: "EMAIL_VERIFICATION",
		1: "PASSWORD_RESET",
	}
	PurposeOTP_value = map[string]int32{
		"EMAIL_VERIFICATION": 0,
		"PASSWORD_RESET":     1,
	}
)

func (x PurposeOTP) Enum() *PurposeOTP {
	p := new(PurposeOTP)
	*p = x
	return p
}

func (x PurposeOTP) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PurposeOTP) Descriptor() protoreflect.EnumDescriptor {
	return file_email_verify_proto_enumTypes[1].Descriptor()
}

func (PurposeOTP) Type() protoreflect.EnumType {
	return &file_email_verify_proto_enumTypes[1]
}

func (x PurposeOTP) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PurposeOTP.Descriptor instead.
func (PurposeOTP) EnumDescriptor() ([]byte, []int) {
	return file_email_verify_proto_rawDescGZIP(), []int{1}
}

type VerifyOTPMessage struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          TypeVerifyOTP          `protobuf:"varint,1,opt,name=type,proto3,enum=TypeVerifyOTP" json:"type,omitempty"`
	Otp           string                 `protobuf:"bytes,2,opt,name=otp,proto3" json:"otp,omitempty"`
	To            string                 `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Fullname      string                 `protobuf:"bytes,4,opt,name=fullname,proto3" json:"fullname,omitempty"`
	Purpose       PurposeOTP             `protobuf:"varint,5,opt,name=purpose,proto3,enum=PurposeOTP" json:"purpose,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VerifyOTPMessage) Reset() {
	*x = VerifyOTPMessage{}
	mi := &file_email_verify_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifyOTPMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyOTPMessage) ProtoMessage() {}

func (x *VerifyOTPMessage) ProtoReflect() protoreflect.Message {
	mi := &file_email_verify_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifyOTPMessage.ProtoReflect.Descriptor instead.
func (*VerifyOTPMessage) Descriptor() ([]byte, []int) {
	return file_email_verify_proto_rawDescGZIP(), []int{0}
}

func (x *VerifyOTPMessage) GetType() TypeVerifyOTP {
	if x != nil {
		return x.Type
	}
	return TypeVerifyOTP_EMAIL
}

func (x *VerifyOTPMessage) GetOtp() string {
	if x != nil {
		return x.Otp
	}
	return ""
}

func (x *VerifyOTPMessage) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *VerifyOTPMessage) GetFullname() string {
	if x != nil {
		return x.Fullname
	}
	return ""
}

func (x *VerifyOTPMessage) GetPurpose() PurposeOTP {
	if x != nil {
		return x.Purpose
	}
	return PurposeOTP_EMAIL_VERIFICATION
}

var File_email_verify_proto protoreflect.FileDescriptor

var file_email_verify_proto_rawDesc = string([]byte{
	0x0a, 0x12, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9b, 0x01, 0x0a, 0x10, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x4f,
	0x54, 0x50, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x4f, 0x54, 0x50, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a,
	0x03, 0x6f, 0x74, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6f, 0x74, 0x70, 0x12,
	0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x07, 0x70,
	0x75, 0x72, 0x70, 0x6f, 0x73, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x50,
	0x75, 0x72, 0x70, 0x6f, 0x73, 0x65, 0x4f, 0x54, 0x50, 0x52, 0x07, 0x70, 0x75, 0x72, 0x70, 0x6f,
	0x73, 0x65, 0x2a, 0x25, 0x0a, 0x0d, 0x54, 0x79, 0x70, 0x65, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x4f, 0x54, 0x50, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x4d, 0x41, 0x49, 0x4c, 0x10, 0x00, 0x12, 0x09,
	0x0a, 0x05, 0x50, 0x48, 0x4f, 0x4e, 0x45, 0x10, 0x01, 0x2a, 0x38, 0x0a, 0x0a, 0x50, 0x75, 0x72,
	0x70, 0x6f, 0x73, 0x65, 0x4f, 0x54, 0x50, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x4d, 0x41, 0x49, 0x4c,
	0x5f, 0x56, 0x45, 0x52, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x00, 0x12,
	0x12, 0x0a, 0x0e, 0x50, 0x41, 0x53, 0x53, 0x57, 0x4f, 0x52, 0x44, 0x5f, 0x52, 0x45, 0x53, 0x45,
	0x54, 0x10, 0x01, 0x42, 0x1a, 0x5a, 0x18, 0x2e, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x67, 0x65, 0x6e, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_email_verify_proto_rawDescOnce sync.Once
	file_email_verify_proto_rawDescData []byte
)

func file_email_verify_proto_rawDescGZIP() []byte {
	file_email_verify_proto_rawDescOnce.Do(func() {
		file_email_verify_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_email_verify_proto_rawDesc), len(file_email_verify_proto_rawDesc)))
	})
	return file_email_verify_proto_rawDescData
}

var file_email_verify_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_email_verify_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_email_verify_proto_goTypes = []any{
	(TypeVerifyOTP)(0),       // 0: TypeVerifyOTP
	(PurposeOTP)(0),          // 1: PurposeOTP
	(*VerifyOTPMessage)(nil), // 2: VerifyOTPMessage
}
var file_email_verify_proto_depIdxs = []int32{
	0, // 0: VerifyOTPMessage.type:type_name -> TypeVerifyOTP
	1, // 1: VerifyOTPMessage.purpose:type_name -> PurposeOTP
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_email_verify_proto_init() }
func file_email_verify_proto_init() {
	if File_email_verify_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_email_verify_proto_rawDesc), len(file_email_verify_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_email_verify_proto_goTypes,
		DependencyIndexes: file_email_verify_proto_depIdxs,
		EnumInfos:         file_email_verify_proto_enumTypes,
		MessageInfos:      file_email_verify_proto_msgTypes,
	}.Build()
	File_email_verify_proto = out.File
	file_email_verify_proto_goTypes = nil
	file_email_verify_proto_depIdxs = nil
}
