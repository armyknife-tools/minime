/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-gogo.
// source: k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto
// DO NOT EDIT!

/*
	Package schema is a generated protocol buffer package.

	It is generated from these files:
		k8s.io/kubernetes/vendor/k8s.io/apimachinery/pkg/runtime/schema/generated.proto

	It has these top-level messages:
*/
package schema

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.GoGoProtoPackageIsVersion1

var fileDescriptorGenerated = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0xce, 0x2f, 0x4e, 0x05, 0x31,
	0x10, 0xc7, 0xf1, 0xd6, 0x20, 0x90, 0xc8, 0x27, 0x46, 0x12, 0x0c, 0x1d, 0x81, 0x41, 0x73, 0x01,
	0x3c, 0xae, 0xbb, 0x6f, 0xe8, 0x36, 0xa5, 0x7f, 0xd2, 0x4e, 0x49, 0x70, 0x1c, 0x81, 0x63, 0xad,
	0x5c, 0x89, 0x64, 0xcb, 0x45, 0x48, 0xda, 0x15, 0x84, 0x04, 0xd7, 0x5f, 0x9a, 0xcf, 0xe4, 0x7b,
	0xf9, 0xe8, 0xee, 0x8b, 0xb2, 0x11, 0x5d, 0x9d, 0x28, 0x07, 0x62, 0x2a, 0xf8, 0x4a, 0xe1, 0x1c,
	0x33, 0x1e, 0x1f, 0x3a, 0x59, 0xaf, 0xe7, 0xc5, 0x06, 0xca, 0x6f, 0x98, 0x9c, 0xc1, 0x5c, 0x03,
	0x5b, 0x4f, 0x58, 0xe6, 0x85, 0xbc, 0x46, 0x43, 0x81, 0xb2, 0x66, 0x3a, 0xab, 0x94, 0x23, 0xc7,
	0xab, 0xeb, 0xe1, 0xd4, 0x6f, 0xa7, 0x92, 0x33, 0xea, 0x70, 0x6a, 0xb8, 0xd3, 0xad, 0xb1, 0xbc,
	0xd4, 0x49, 0xcd, 0xd1, 0xa3, 0x89, 0x26, 0x62, 0xe7, 0x53, 0x7d, 0xee, 0xab, 0x8f, 0xfe, 0x1a,
	0x67, 0x4f, 0x77, 0xff, 0xe5, 0x54, 0xb6, 0x2f, 0x68, 0x03, 0x17, 0xce, 0x7f, 0x5b, 0x1e, 0x6e,
	0xd6, 0x1d, 0xc4, 0xb6, 0x83, 0xf8, 0xdc, 0x41, 0xbc, 0x37, 0x90, 0x6b, 0x03, 0xb9, 0x35, 0x90,
	0x5f, 0x0d, 0xe4, 0xc7, 0x37, 0x88, 0xa7, 0x8b, 0x51, 0xf3, 0x13, 0x00, 0x00, 0xff, 0xff, 0xd9,
	0x82, 0x09, 0xbe, 0x07, 0x01, 0x00, 0x00,
}
