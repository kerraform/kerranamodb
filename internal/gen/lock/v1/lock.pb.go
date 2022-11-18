// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: lock/v1/lock.proto

package lockv1

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

type LockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid    string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Table  string `protobuf:"bytes,3,opt,name=table,proto3" json:"table,omitempty"`
	Key    string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *LockRequest) Reset() {
	*x = LockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockRequest) ProtoMessage() {}

func (x *LockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockRequest.ProtoReflect.Descriptor instead.
func (*LockRequest) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{0}
}

func (x *LockRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *LockRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *LockRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *LockRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type UnlockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid    string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Table  string `protobuf:"bytes,3,opt,name=table,proto3" json:"table,omitempty"`
	Key    string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *UnlockRequest) Reset() {
	*x = UnlockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockRequest) ProtoMessage() {}

func (x *UnlockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockRequest.ProtoReflect.Descriptor instead.
func (*UnlockRequest) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{1}
}

func (x *UnlockRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *UnlockRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *UnlockRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *UnlockRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type RLockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid    string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Table  string `protobuf:"bytes,3,opt,name=table,proto3" json:"table,omitempty"`
	Key    string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *RLockRequest) Reset() {
	*x = RLockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RLockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RLockRequest) ProtoMessage() {}

func (x *RLockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RLockRequest.ProtoReflect.Descriptor instead.
func (*RLockRequest) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{2}
}

func (x *RLockRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RLockRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *RLockRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *RLockRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type RUnlockRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uid    string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Source string `protobuf:"bytes,2,opt,name=source,proto3" json:"source,omitempty"`
	Table  string `protobuf:"bytes,3,opt,name=table,proto3" json:"table,omitempty"`
	Key    string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *RUnlockRequest) Reset() {
	*x = RUnlockRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RUnlockRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RUnlockRequest) ProtoMessage() {}

func (x *RUnlockRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RUnlockRequest.ProtoReflect.Descriptor instead.
func (*RUnlockRequest) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{3}
}

func (x *RUnlockRequest) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *RUnlockRequest) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *RUnlockRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *RUnlockRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type LockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Available bool `protobuf:"varint,1,opt,name=available,proto3" json:"available,omitempty"`
}

func (x *LockResponse) Reset() {
	*x = LockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockResponse) ProtoMessage() {}

func (x *LockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockResponse.ProtoReflect.Descriptor instead.
func (*LockResponse) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{4}
}

func (x *LockResponse) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

type UnlockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Available bool `protobuf:"varint,1,opt,name=available,proto3" json:"available,omitempty"`
}

func (x *UnlockResponse) Reset() {
	*x = UnlockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockResponse) ProtoMessage() {}

func (x *UnlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockResponse.ProtoReflect.Descriptor instead.
func (*UnlockResponse) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{5}
}

func (x *UnlockResponse) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

type RLockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Available bool `protobuf:"varint,1,opt,name=available,proto3" json:"available,omitempty"`
}

func (x *RLockResponse) Reset() {
	*x = RLockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RLockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RLockResponse) ProtoMessage() {}

func (x *RLockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RLockResponse.ProtoReflect.Descriptor instead.
func (*RLockResponse) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{6}
}

func (x *RLockResponse) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

type RUnlockResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Available bool `protobuf:"varint,1,opt,name=available,proto3" json:"available,omitempty"`
}

func (x *RUnlockResponse) Reset() {
	*x = RUnlockResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lock_v1_lock_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RUnlockResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RUnlockResponse) ProtoMessage() {}

func (x *RUnlockResponse) ProtoReflect() protoreflect.Message {
	mi := &file_lock_v1_lock_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RUnlockResponse.ProtoReflect.Descriptor instead.
func (*RUnlockResponse) Descriptor() ([]byte, []int) {
	return file_lock_v1_lock_proto_rawDescGZIP(), []int{7}
}

func (x *RUnlockResponse) GetAvailable() bool {
	if x != nil {
		return x.Available
	}
	return false
}

var File_lock_v1_lock_proto protoreflect.FileDescriptor

var file_lock_v1_lock_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x22, 0x5f, 0x0a,
	0x0b, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x61,
	0x0a, 0x0d, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x22, 0x60, 0x0a, 0x0c, 0x52, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x62, 0x0a, 0x0e, 0x52, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x2c, 0x0a, 0x0c, 0x4c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2e, 0x0a, 0x0e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2d, 0x0a, 0x0d, 0x52, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x22, 0x2f, 0x0a, 0x0f, 0x52, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x61, 0x76, 0x61, 0x69,
	0x6c, 0x61, 0x62, 0x6c, 0x65, 0x32, 0xfb, 0x01, 0x0a, 0x0b, 0x4c, 0x6f, 0x63, 0x6b, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x04, 0x4c, 0x6f, 0x63, 0x6b, 0x12, 0x14, 0x2e,
	0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x06,
	0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31,
	0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x05, 0x52, 0x4c, 0x6f,
	0x63, 0x6b, 0x12, 0x15, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x4c, 0x6f,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x6c, 0x6f, 0x63, 0x6b,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x4c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x07, 0x52, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x17,
	0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6b, 0x65, 0x72, 0x72, 0x61, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x6b, 0x65, 0x72, 0x72,
	0x61, 0x6e, 0x61, 0x6d, 0x6f, 0x64, 0x62, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x6c, 0x6f, 0x63, 0x6b, 0x2f, 0x76, 0x31, 0x3b, 0x6c, 0x6f, 0x63,
	0x6b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lock_v1_lock_proto_rawDescOnce sync.Once
	file_lock_v1_lock_proto_rawDescData = file_lock_v1_lock_proto_rawDesc
)

func file_lock_v1_lock_proto_rawDescGZIP() []byte {
	file_lock_v1_lock_proto_rawDescOnce.Do(func() {
		file_lock_v1_lock_proto_rawDescData = protoimpl.X.CompressGZIP(file_lock_v1_lock_proto_rawDescData)
	})
	return file_lock_v1_lock_proto_rawDescData
}

var file_lock_v1_lock_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_lock_v1_lock_proto_goTypes = []interface{}{
	(*LockRequest)(nil),     // 0: lock.v1.LockRequest
	(*UnlockRequest)(nil),   // 1: lock.v1.UnlockRequest
	(*RLockRequest)(nil),    // 2: lock.v1.RLockRequest
	(*RUnlockRequest)(nil),  // 3: lock.v1.RUnlockRequest
	(*LockResponse)(nil),    // 4: lock.v1.LockResponse
	(*UnlockResponse)(nil),  // 5: lock.v1.UnlockResponse
	(*RLockResponse)(nil),   // 6: lock.v1.RLockResponse
	(*RUnlockResponse)(nil), // 7: lock.v1.RUnlockResponse
}
var file_lock_v1_lock_proto_depIdxs = []int32{
	0, // 0: lock.v1.LockService.Lock:input_type -> lock.v1.LockRequest
	1, // 1: lock.v1.LockService.Unlock:input_type -> lock.v1.UnlockRequest
	2, // 2: lock.v1.LockService.RLock:input_type -> lock.v1.RLockRequest
	3, // 3: lock.v1.LockService.RUnlock:input_type -> lock.v1.RUnlockRequest
	4, // 4: lock.v1.LockService.Lock:output_type -> lock.v1.LockResponse
	5, // 5: lock.v1.LockService.Unlock:output_type -> lock.v1.UnlockResponse
	6, // 6: lock.v1.LockService.RLock:output_type -> lock.v1.RLockResponse
	7, // 7: lock.v1.LockService.RUnlock:output_type -> lock.v1.RUnlockResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_lock_v1_lock_proto_init() }
func file_lock_v1_lock_proto_init() {
	if File_lock_v1_lock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lock_v1_lock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LockRequest); i {
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
		file_lock_v1_lock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockRequest); i {
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
		file_lock_v1_lock_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RLockRequest); i {
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
		file_lock_v1_lock_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RUnlockRequest); i {
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
		file_lock_v1_lock_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LockResponse); i {
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
		file_lock_v1_lock_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockResponse); i {
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
		file_lock_v1_lock_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RLockResponse); i {
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
		file_lock_v1_lock_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RUnlockResponse); i {
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
			RawDescriptor: file_lock_v1_lock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lock_v1_lock_proto_goTypes,
		DependencyIndexes: file_lock_v1_lock_proto_depIdxs,
		MessageInfos:      file_lock_v1_lock_proto_msgTypes,
	}.Build()
	File_lock_v1_lock_proto = out.File
	file_lock_v1_lock_proto_rawDesc = nil
	file_lock_v1_lock_proto_goTypes = nil
	file_lock_v1_lock_proto_depIdxs = nil
}