// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: hotstuff.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
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

type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ParentHash []byte      `protobuf:"bytes,1,opt,name=ParentHash,proto3" json:"ParentHash,omitempty"`
	QC         *QuorumCert `protobuf:"bytes,2,opt,name=QC,proto3" json:"QC,omitempty"`
	Height     int64       `protobuf:"varint,3,opt,name=Height,proto3" json:"Height,omitempty"`
	Commands   []*Command  `protobuf:"bytes,4,rep,name=Commands,proto3" json:"Commands,omitempty"`
	ProposerID uint32      `protobuf:"varint,5,opt,name=ProposerID,proto3" json:"ProposerID,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotstuff_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_hotstuff_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_hotstuff_proto_rawDescGZIP(), []int{0}
}

func (x *Block) GetParentHash() []byte {
	if x != nil {
		return x.ParentHash
	}
	return nil
}

func (x *Block) GetQC() *QuorumCert {
	if x != nil {
		return x.QC
	}
	return nil
}

func (x *Block) GetHeight() int64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Block) GetCommands() []*Command {
	if x != nil {
		return x.Commands
	}
	return nil
}

func (x *Block) GetProposerID() uint32 {
	if x != nil {
		return x.ProposerID
	}
	return 0
}

type PartialSig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicaID int32  `protobuf:"varint,1,opt,name=ReplicaID,proto3" json:"ReplicaID,omitempty"`
	R         []byte `protobuf:"bytes,2,opt,name=R,proto3" json:"R,omitempty"`
	S         []byte `protobuf:"bytes,3,opt,name=S,proto3" json:"S,omitempty"`
}

func (x *PartialSig) Reset() {
	*x = PartialSig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotstuff_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialSig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialSig) ProtoMessage() {}

func (x *PartialSig) ProtoReflect() protoreflect.Message {
	mi := &file_hotstuff_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialSig.ProtoReflect.Descriptor instead.
func (*PartialSig) Descriptor() ([]byte, []int) {
	return file_hotstuff_proto_rawDescGZIP(), []int{1}
}

func (x *PartialSig) GetReplicaID() int32 {
	if x != nil {
		return x.ReplicaID
	}
	return 0
}

func (x *PartialSig) GetR() []byte {
	if x != nil {
		return x.R
	}
	return nil
}

func (x *PartialSig) GetS() []byte {
	if x != nil {
		return x.S
	}
	return nil
}

type PartialCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sig  *PartialSig `protobuf:"bytes,1,opt,name=Sig,proto3" json:"Sig,omitempty"`
	Hash []byte      `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *PartialCert) Reset() {
	*x = PartialCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotstuff_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialCert) ProtoMessage() {}

func (x *PartialCert) ProtoReflect() protoreflect.Message {
	mi := &file_hotstuff_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialCert.ProtoReflect.Descriptor instead.
func (*PartialCert) Descriptor() ([]byte, []int) {
	return file_hotstuff_proto_rawDescGZIP(), []int{2}
}

func (x *PartialCert) GetSig() *PartialSig {
	if x != nil {
		return x.Sig
	}
	return nil
}

func (x *PartialCert) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type QuorumCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sigs []*PartialSig `protobuf:"bytes,1,rep,name=Sigs,proto3" json:"Sigs,omitempty"`
	Hash []byte        `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
}

func (x *QuorumCert) Reset() {
	*x = QuorumCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotstuff_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuorumCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuorumCert) ProtoMessage() {}

func (x *QuorumCert) ProtoReflect() protoreflect.Message {
	mi := &file_hotstuff_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuorumCert.ProtoReflect.Descriptor instead.
func (*QuorumCert) Descriptor() ([]byte, []int) {
	return file_hotstuff_proto_rawDescGZIP(), []int{3}
}

func (x *QuorumCert) GetSigs() []*PartialSig {
	if x != nil {
		return x.Sigs
	}
	return nil
}

func (x *QuorumCert) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_hotstuff_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_hotstuff_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_hotstuff_proto_rawDescGZIP(), []int{4}
}

func (x *Command) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_hotstuff_proto protoreflect.FileDescriptor

var file_hotstuff_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x68, 0x6f, 0x74, 0x73, 0x74, 0x75, 0x66, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xae, 0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x1e,
	0x0a, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x0a, 0x50, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x21,
	0x0a, 0x02, 0x51, 0x43, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52, 0x02, 0x51,
	0x43, 0x12, 0x16, 0x0a, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x2a, 0x0a, 0x08, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x08, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x70, 0x6f, 0x73, 0x65,
	0x72, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x50, 0x72, 0x6f, 0x70, 0x6f,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x46, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c,
	0x53, 0x69, 0x67, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49,
	0x44, 0x12, 0x0c, 0x0a, 0x01, 0x52, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x52, 0x12,
	0x0c, 0x0a, 0x01, 0x53, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x53, 0x22, 0x46, 0x0a,
	0x0b, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x43, 0x65, 0x72, 0x74, 0x12, 0x23, 0x0a, 0x03,
	0x53, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x52, 0x03, 0x53, 0x69,
	0x67, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x22, 0x47, 0x0a, 0x0a, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43,
	0x65, 0x72, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x53, 0x69, 0x67, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61,
	0x6c, 0x53, 0x69, 0x67, 0x52, 0x04, 0x53, 0x69, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61,
	0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x22, 0x1d,
	0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x32, 0x73, 0x0a,
	0x08, 0x48, 0x6f, 0x74, 0x73, 0x74, 0x75, 0x66, 0x66, 0x12, 0x31, 0x0a, 0x07, 0x50, 0x72, 0x6f,
	0x70, 0x6f, 0x73, 0x65, 0x12, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x04,
	0x56, 0x6f, 0x74, 0x65, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x61, 0x6c, 0x43, 0x65, 0x72, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6a, 0x6f, 0x65, 0x2d, 0x7a, 0x78, 0x68, 0x2f, 0x68, 0x6f, 0x74, 0x73, 0x74, 0x75, 0x66,
	0x66, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hotstuff_proto_rawDescOnce sync.Once
	file_hotstuff_proto_rawDescData = file_hotstuff_proto_rawDesc
)

func file_hotstuff_proto_rawDescGZIP() []byte {
	file_hotstuff_proto_rawDescOnce.Do(func() {
		file_hotstuff_proto_rawDescData = protoimpl.X.CompressGZIP(file_hotstuff_proto_rawDescData)
	})
	return file_hotstuff_proto_rawDescData
}

var file_hotstuff_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_hotstuff_proto_goTypes = []interface{}{
	(*Block)(nil),       // 0: proto.Block
	(*PartialSig)(nil),  // 1: proto.PartialSig
	(*PartialCert)(nil), // 2: proto.PartialCert
	(*QuorumCert)(nil),  // 3: proto.QuorumCert
	(*Command)(nil),     // 4: proto.Command
	(*empty.Empty)(nil), // 5: google.protobuf.Empty
}
var file_hotstuff_proto_depIdxs = []int32{
	3, // 0: proto.Block.QC:type_name -> proto.QuorumCert
	4, // 1: proto.Block.Commands:type_name -> proto.Command
	1, // 2: proto.PartialCert.Sig:type_name -> proto.PartialSig
	1, // 3: proto.QuorumCert.Sigs:type_name -> proto.PartialSig
	0, // 4: proto.Hotstuff.Propose:input_type -> proto.Block
	2, // 5: proto.Hotstuff.Vote:input_type -> proto.PartialCert
	5, // 6: proto.Hotstuff.Propose:output_type -> google.protobuf.Empty
	5, // 7: proto.Hotstuff.Vote:output_type -> google.protobuf.Empty
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_hotstuff_proto_init() }
func file_hotstuff_proto_init() {
	if File_hotstuff_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_hotstuff_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
		file_hotstuff_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialSig); i {
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
		file_hotstuff_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialCert); i {
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
		file_hotstuff_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuorumCert); i {
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
		file_hotstuff_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
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
			RawDescriptor: file_hotstuff_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hotstuff_proto_goTypes,
		DependencyIndexes: file_hotstuff_proto_depIdxs,
		MessageInfos:      file_hotstuff_proto_msgTypes,
	}.Build()
	File_hotstuff_proto = out.File
	file_hotstuff_proto_rawDesc = nil
	file_hotstuff_proto_goTypes = nil
	file_hotstuff_proto_depIdxs = nil
}
