// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/vision/v1p3beta1/geometry.proto

package vision // import "google.golang.org/genproto/googleapis/cloud/vision/v1p3beta1"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "google.golang.org/genproto/googleapis/api/annotations"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A vertex represents a 2D point in the image.
// NOTE: the vertex coordinates are in the same scale as the original image.
type Vertex struct {
	// X coordinate.
	X int32 `protobuf:"varint,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y                    int32    `protobuf:"varint,2,opt,name=y,proto3" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vertex) Reset()         { *m = Vertex{} }
func (m *Vertex) String() string { return proto.CompactTextString(m) }
func (*Vertex) ProtoMessage()    {}
func (*Vertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_369491ddc24d6412, []int{0}
}
func (m *Vertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vertex.Unmarshal(m, b)
}
func (m *Vertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vertex.Marshal(b, m, deterministic)
}
func (dst *Vertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vertex.Merge(dst, src)
}
func (m *Vertex) XXX_Size() int {
	return xxx_messageInfo_Vertex.Size(m)
}
func (m *Vertex) XXX_DiscardUnknown() {
	xxx_messageInfo_Vertex.DiscardUnknown(m)
}

var xxx_messageInfo_Vertex proto.InternalMessageInfo

func (m *Vertex) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Vertex) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// A vertex represents a 2D point in the image.
// NOTE: the normalized vertex coordinates are relative to the original image
// and range from 0 to 1.
type NormalizedVertex struct {
	// X coordinate.
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y                    float32  `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NormalizedVertex) Reset()         { *m = NormalizedVertex{} }
func (m *NormalizedVertex) String() string { return proto.CompactTextString(m) }
func (*NormalizedVertex) ProtoMessage()    {}
func (*NormalizedVertex) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_369491ddc24d6412, []int{1}
}
func (m *NormalizedVertex) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NormalizedVertex.Unmarshal(m, b)
}
func (m *NormalizedVertex) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NormalizedVertex.Marshal(b, m, deterministic)
}
func (dst *NormalizedVertex) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NormalizedVertex.Merge(dst, src)
}
func (m *NormalizedVertex) XXX_Size() int {
	return xxx_messageInfo_NormalizedVertex.Size(m)
}
func (m *NormalizedVertex) XXX_DiscardUnknown() {
	xxx_messageInfo_NormalizedVertex.DiscardUnknown(m)
}

var xxx_messageInfo_NormalizedVertex proto.InternalMessageInfo

func (m *NormalizedVertex) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *NormalizedVertex) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

// A bounding polygon for the detected image annotation.
type BoundingPoly struct {
	// The bounding polygon vertices.
	Vertices []*Vertex `protobuf:"bytes,1,rep,name=vertices,proto3" json:"vertices,omitempty"`
	// The bounding polygon normalized vertices.
	NormalizedVertices   []*NormalizedVertex `protobuf:"bytes,2,rep,name=normalized_vertices,json=normalizedVertices,proto3" json:"normalized_vertices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *BoundingPoly) Reset()         { *m = BoundingPoly{} }
func (m *BoundingPoly) String() string { return proto.CompactTextString(m) }
func (*BoundingPoly) ProtoMessage()    {}
func (*BoundingPoly) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_369491ddc24d6412, []int{2}
}
func (m *BoundingPoly) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BoundingPoly.Unmarshal(m, b)
}
func (m *BoundingPoly) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BoundingPoly.Marshal(b, m, deterministic)
}
func (dst *BoundingPoly) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoundingPoly.Merge(dst, src)
}
func (m *BoundingPoly) XXX_Size() int {
	return xxx_messageInfo_BoundingPoly.Size(m)
}
func (m *BoundingPoly) XXX_DiscardUnknown() {
	xxx_messageInfo_BoundingPoly.DiscardUnknown(m)
}

var xxx_messageInfo_BoundingPoly proto.InternalMessageInfo

func (m *BoundingPoly) GetVertices() []*Vertex {
	if m != nil {
		return m.Vertices
	}
	return nil
}

func (m *BoundingPoly) GetNormalizedVertices() []*NormalizedVertex {
	if m != nil {
		return m.NormalizedVertices
	}
	return nil
}

// A normalized bounding polygon around a portion of an image.
type NormalizedBoundingPoly struct {
	// Normalized vertices of the bounding polygon.
	Vertices             []*NormalizedVertex `protobuf:"bytes,1,rep,name=vertices,proto3" json:"vertices,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *NormalizedBoundingPoly) Reset()         { *m = NormalizedBoundingPoly{} }
func (m *NormalizedBoundingPoly) String() string { return proto.CompactTextString(m) }
func (*NormalizedBoundingPoly) ProtoMessage()    {}
func (*NormalizedBoundingPoly) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_369491ddc24d6412, []int{3}
}
func (m *NormalizedBoundingPoly) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NormalizedBoundingPoly.Unmarshal(m, b)
}
func (m *NormalizedBoundingPoly) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NormalizedBoundingPoly.Marshal(b, m, deterministic)
}
func (dst *NormalizedBoundingPoly) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NormalizedBoundingPoly.Merge(dst, src)
}
func (m *NormalizedBoundingPoly) XXX_Size() int {
	return xxx_messageInfo_NormalizedBoundingPoly.Size(m)
}
func (m *NormalizedBoundingPoly) XXX_DiscardUnknown() {
	xxx_messageInfo_NormalizedBoundingPoly.DiscardUnknown(m)
}

var xxx_messageInfo_NormalizedBoundingPoly proto.InternalMessageInfo

func (m *NormalizedBoundingPoly) GetVertices() []*NormalizedVertex {
	if m != nil {
		return m.Vertices
	}
	return nil
}

// A 3D position in the image, used primarily for Face detection landmarks.
// A valid Position must have both x and y coordinates.
// The position coordinates are in the same scale as the original image.
type Position struct {
	// X coordinate.
	X float32 `protobuf:"fixed32,1,opt,name=x,proto3" json:"x,omitempty"`
	// Y coordinate.
	Y float32 `protobuf:"fixed32,2,opt,name=y,proto3" json:"y,omitempty"`
	// Z coordinate (or depth).
	Z                    float32  `protobuf:"fixed32,3,opt,name=z,proto3" json:"z,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_geometry_369491ddc24d6412, []int{4}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (dst *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(dst, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetX() float32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Position) GetY() float32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Position) GetZ() float32 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Vertex)(nil), "google.cloud.vision.v1p3beta1.Vertex")
	proto.RegisterType((*NormalizedVertex)(nil), "google.cloud.vision.v1p3beta1.NormalizedVertex")
	proto.RegisterType((*BoundingPoly)(nil), "google.cloud.vision.v1p3beta1.BoundingPoly")
	proto.RegisterType((*NormalizedBoundingPoly)(nil), "google.cloud.vision.v1p3beta1.NormalizedBoundingPoly")
	proto.RegisterType((*Position)(nil), "google.cloud.vision.v1p3beta1.Position")
}

func init() {
	proto.RegisterFile("google/cloud/vision/v1p3beta1/geometry.proto", fileDescriptor_geometry_369491ddc24d6412)
}

var fileDescriptor_geometry_369491ddc24d6412 = []byte{
	// 316 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0x49, 0x87, 0x63, 0xc4, 0x09, 0x52, 0x41, 0x8a, 0x28, 0xcc, 0xa2, 0xb0, 0x83, 0x24,
	0xcc, 0x79, 0xf3, 0xe4, 0x3c, 0x78, 0x10, 0xa4, 0xf4, 0xe0, 0xc1, 0x8b, 0x66, 0x6d, 0x08, 0x81,
	0x36, 0xaf, 0xa4, 0x59, 0x59, 0x8b, 0xff, 0x95, 0xff, 0x9c, 0x47, 0x69, 0x53, 0x2a, 0x9d, 0x58,
	0xf1, 0xf8, 0xbd, 0xfc, 0xde, 0xf7, 0x3e, 0xf2, 0x1e, 0xbe, 0x12, 0x00, 0x22, 0xe1, 0x34, 0x4a,
	0x60, 0x13, 0xd3, 0x42, 0xe6, 0x12, 0x14, 0x2d, 0x16, 0xd9, 0x72, 0xcd, 0x0d, 0x5b, 0x50, 0xc1,
	0x21, 0xe5, 0x46, 0x97, 0x24, 0xd3, 0x60, 0xc0, 0x3d, 0xb3, 0x34, 0x69, 0x68, 0x62, 0x69, 0xd2,
	0xd1, 0x27, 0xa7, 0xad, 0x19, 0xcb, 0x24, 0x65, 0x4a, 0x81, 0x61, 0x46, 0x82, 0xca, 0x6d, 0xb3,
	0x7f, 0x81, 0xc7, 0xcf, 0x5c, 0x1b, 0xbe, 0x75, 0xa7, 0x18, 0x6d, 0x3d, 0x34, 0x43, 0xf3, 0xbd,
	0x10, 0x35, 0xaa, 0xf4, 0x1c, 0xab, 0x4a, 0x9f, 0xe0, 0xc3, 0x27, 0xd0, 0x29, 0x4b, 0x64, 0xc5,
	0xe3, 0x5d, 0xde, 0xe9, 0xf1, 0x4e, 0xcd, 0x7f, 0x20, 0x3c, 0x5d, 0xc1, 0x46, 0xc5, 0x52, 0x89,
	0x00, 0x92, 0xd2, 0xbd, 0xc3, 0x93, 0x82, 0x6b, 0x23, 0x23, 0x9e, 0x7b, 0x68, 0x36, 0x9a, 0xef,
	0x5f, 0x5f, 0x92, 0xc1, 0xd8, 0xc4, 0x4e, 0x09, 0xbb, 0x36, 0xf7, 0x0d, 0x1f, 0xa9, 0x2e, 0xc3,
	0x6b, 0xe7, 0xe6, 0x34, 0x6e, 0xf4, 0x0f, 0xb7, 0xdd, 0xf4, 0xa1, 0xab, 0x7a, 0x95, 0xda, 0xca,
	0xe7, 0xf8, 0xf8, 0x9b, 0xeb, 0xc5, 0x7f, 0xfc, 0x11, 0xff, 0xdf, 0x03, 0x3b, 0x03, 0xff, 0x06,
	0x4f, 0x02, 0xc8, 0x65, 0xbd, 0x85, 0xa1, 0x4f, 0xac, 0x55, 0xe5, 0x8d, 0xac, 0xaa, 0x56, 0xef,
	0xf8, 0x3c, 0x82, 0x74, 0x78, 0xea, 0xea, 0xe0, 0xa1, 0x3d, 0x8d, 0xa0, 0x5e, 0x6e, 0x80, 0x5e,
	0xee, 0x5b, 0x5e, 0x40, 0xc2, 0x94, 0x20, 0xa0, 0x05, 0x15, 0x5c, 0x35, 0xab, 0xa7, 0xf6, 0x89,
	0x65, 0x32, 0xff, 0xe5, 0xd0, 0x6e, 0x6d, 0xe1, 0x13, 0xa1, 0xf5, 0xb8, 0x69, 0x59, 0x7e, 0x05,
	0x00, 0x00, 0xff, 0xff, 0x79, 0x59, 0xbf, 0x39, 0x9a, 0x02, 0x00, 0x00,
}
