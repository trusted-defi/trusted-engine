// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: trusted/v1/service.proto

package trustedv1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/emptypb"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("trusted/v1/service.proto", fileDescriptor_ceb9ab9b8e4b32ac) }

var fileDescriptor_ceb9ab9b8e4b32ac = []byte{
	// 687 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0xb7, 0x4e, 0xa0, 0xe1, 0x55, 0x43, 0xe4, 0x30, 0xd6, 0x6e, 0x30, 0xb6, 0x3b, 0x29,
	0x1d, 0xb7, 0x21, 0x21, 0xad, 0x65, 0x6d, 0x85, 0xc6, 0x14, 0xad, 0xd1, 0x34, 0x60, 0x12, 0x72,
	0x93, 0x77, 0x5d, 0x44, 0x1a, 0x17, 0xdb, 0x29, 0xdd, 0xd7, 0xe1, 0x88, 0xc4, 0x17, 0xe1, 0x63,
	0x20, 0x4e, 0x1c, 0xf9, 0x04, 0xc8, 0xf5, 0x9f, 0xda, 0x5d, 0x82, 0x90, 0xd8, 0xad, 0xf6, 0xf3,
	0xf8, 0xe7, 0xc7, 0x7d, 0xed, 0x37, 0x68, 0x93, 0xd3, 0x9c, 0x71, 0x88, 0x1b, 0x93, 0x66, 0x83,
	0x01, 0x9d, 0x24, 0x11, 0xf8, 0x63, 0x4a, 0x38, 0xf1, 0x90, 0x52, 0xfc, 0x49, 0xb3, 0xbe, 0x35,
	0x24, 0x64, 0x98, 0x42, 0x63, 0xa6, 0x0c, 0xf2, 0xcb, 0x06, 0x8c, 0xc6, 0xfc, 0x5a, 0x1a, 0xeb,
	0xbb, 0x16, 0x82, 0xc2, 0xa7, 0x1c, 0x18, 0xff, 0x40, 0x81, 0x8d, 0x49, 0xc6, 0x14, 0x6b, 0xff,
	0x37, 0x42, 0xeb, 0xa1, 0x74, 0xf5, 0xe5, 0x26, 0xde, 0x11, 0xaa, 0x06, 0x84, 0xa4, 0x7d, 0xe0,
	0x01, 0x15, 0xe3, 0x2d, 0x7f, 0xbe, 0x9f, 0xaf, 0x67, 0x4f, 0x25, 0xae, 0xbe, 0xe1, 0xcb, 0x00,
	0xbe, 0x0e, 0xe0, 0x1f, 0x89, 0x00, 0x7b, 0x4b, 0x5e, 0x47, 0x62, 0xba, 0x98, 0x49, 0x4c, 0x89,
	0xb3, 0xbe, 0x6d, 0xe3, 0xb5, 0xfb, 0x54, 0xa5, 0xdc, 0x5b, 0xf2, 0xfa, 0xa8, 0x1a, 0x40, 0x16,
	0x27, 0xd9, 0xf0, 0x84, 0x64, 0x11, 0x78, 0x3b, 0xb6, 0xdf, 0x56, 0x74, 0xa4, 0x27, 0xe5, 0x06,
	0x03, 0x6d, 0xa1, 0xd5, 0xd9, 0x19, 0x39, 0xe6, 0xff, 0x16, 0x4c, 0xbb, 0x2d, 0x46, 0x80, 0xd6,
	0xc4, 0x6c, 0x9b, 0x64, 0x1c, 0x32, 0xee, 0x3d, 0x5e, 0xb4, 0x2b, 0x41, 0xc7, 0xda, 0x29, 0xd5,
	0x0d, 0x31, 0x44, 0xf7, 0x2d, 0xa1, 0x43, 0xc9, 0xe8, 0x36, 0xa8, 0x3d, 0x99, 0x53, 0xfd, 0x13,
	0xa5, 0xc7, 0xbd, 0x41, 0x52, 0x0b, 0x2c, 0x52, 0x07, 0x21, 0x21, 0x1c, 0x93, 0x08, 0xa7, 0xac,
	0x14, 0x74, 0x23, 0xb2, 0xf4, 0x3b, 0x9c, 0xb5, 0xc3, 0x38, 0x96, 0xd3, 0xe1, 0xd4, 0xab, 0xd9,
	0x0b, 0x0e, 0xe3, 0x38, 0x9c, 0x32, 0x7d, 0xbc, 0x7a, 0x91, 0xb4, 0xc0, 0x39, 0x85, 0x11, 0xe1,
	0xf0, 0x3f, 0x9c, 0x2e, 0x5a, 0x0d, 0xa7, 0xa2, 0xba, 0x39, 0x73, 0x6f, 0xbb, 0x9e, 0xd5, 0x98,
	0xed, 0x62, 0xd1, 0x80, 0x5e, 0xa2, 0x3b, 0xe1, 0xb4, 0x0b, 0xdc, 0xdb, 0x74, 0x8d, 0x5d, 0x30,
	0x05, 0xab, 0x15, 0x28, 0xee, 0xfa, 0x1e, 0x66, 0x8b, 0xeb, 0x7b, 0x98, 0x95, 0xac, 0x9f, 0x29,
	0x66, 0x7d, 0x8c, 0x1e, 0xf6, 0xf3, 0x01, 0x8b, 0x68, 0x32, 0x80, 0x13, 0xf8, 0x1c, 0x52, 0x9c,
	0x31, 0x1c, 0xf1, 0x84, 0x64, 0xde, 0xae, 0xf3, 0x8a, 0x6d, 0xd3, 0x54, 0xa3, 0xf7, 0xfe, 0x66,
	0xd1, 0x7b, 0x3c, 0x5b, 0x16, 0x29, 0xdb, 0xf4, 0x7a, 0xbc, 0x70, 0xca, 0xd9, 0x54, 0x61, 0x4a,
	0xa5, 0x98, 0x94, 0xe7, 0xe8, 0x81, 0x2e, 0xbf, 0x6a, 0x3d, 0xe1, 0xd4, 0x7d, 0xd6, 0xa2, 0x42,
	0x5a, 0x29, 0x7c, 0xd6, 0xae, 0xc1, 0x90, 0xdf, 0x22, 0x6f, 0x7e, 0x21, 0x6e, 0x17, 0xfd, 0x1a,
	0xdd, 0xeb, 0x24, 0x69, 0xda, 0x4a, 0x49, 0xf4, 0xd1, 0x73, 0xee, 0x81, 0x99, 0xd6, 0xb8, 0x47,
	0x25, 0xaa, 0x66, 0xed, 0xff, 0x5c, 0x41, 0xd5, 0xf6, 0x15, 0x4e, 0x32, 0xdd, 0x72, 0x0f, 0xd1,
	0x6a, 0x17, 0xb8, 0x64, 0x3b, 0x7f, 0xaa, 0xc3, 0xad, 0x15, 0x28, 0xd6, 0x1d, 0x46, 0x02, 0x81,
	0x53, 0x2c, 0x9a, 0xa4, 0x73, 0xdf, 0xd5, 0xa4, 0xc6, 0x6c, 0x15, 0x6a, 0x06, 0x24, 0xb3, 0xc8,
	0x5e, 0xeb, 0x64, 0x71, 0x9a, 0x6c, 0xad, 0x40, 0xb1, 0x5b, 0x76, 0x3b, 0xa7, 0x14, 0x32, 0x75,
	0x24, 0xa7, 0x00, 0xb6, 0x52, 0x58, 0x00, 0xd7, 0x60, 0x43, 0x8f, 0x31, 0x07, 0xc6, 0x7b, 0x80,
	0x63, 0xa0, 0x2e, 0xd4, 0x56, 0x0a, 0xa1, 0xae, 0xc1, 0x40, 0xdf, 0xa3, 0xf5, 0x59, 0x21, 0x84,
	0x70, 0x34, 0x11, 0x6d, 0xdc, 0x79, 0x27, 0xae, 0x56, 0xf8, 0x4e, 0x16, 0x2d, 0xf3, 0x77, 0xd2,
	0xfa, 0xb6, 0x8c, 0xd6, 0x23, 0x32, 0xb2, 0xdc, 0xad, 0xaa, 0xaa, 0x78, 0x20, 0x5a, 0x66, 0xb0,
	0xfc, 0xee, 0xd5, 0x30, 0xe1, 0x57, 0xf9, 0xc0, 0x8f, 0xc8, 0xa8, 0xa1, 0x6c, 0x4f, 0x63, 0xb8,
	0x4c, 0xcc, 0x00, 0xb2, 0x61, 0x92, 0xa9, 0xcf, 0x7b, 0x44, 0xd2, 0xc6, 0xfc, 0x8b, 0xfe, 0x42,
	0xfd, 0x9c, 0x34, 0xbf, 0x54, 0x56, 0xc2, 0xf3, 0xf3, 0xaf, 0x15, 0xa4, 0xee, 0xad, 0x7f, 0xd6,
	0xfc, 0x6e, 0x06, 0x17, 0x67, 0xcd, 0x1f, 0x95, 0x8d, 0xf9, 0xe0, 0xa2, 0x1b, 0xb4, 0xde, 0x00,
	0xc7, 0x31, 0xe6, 0xf8, 0x57, 0x65, 0x4d, 0x09, 0x07, 0x07, 0x67, 0xcd, 0xc1, 0xdd, 0xd9, 0x2e,
	0xcf, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x7d, 0x76, 0x30, 0x9b, 0x7a, 0x08, 0x00, 0x00,
}
