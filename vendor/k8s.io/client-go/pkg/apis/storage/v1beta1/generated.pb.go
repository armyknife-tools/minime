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
// source: k8s.io/kubernetes/pkg/apis/storage/v1beta1/generated.proto
// DO NOT EDIT!

/*
	Package v1beta1 is a generated protocol buffer package.

	It is generated from these files:
		k8s.io/kubernetes/pkg/apis/storage/v1beta1/generated.proto

	It has these top-level messages:
		StorageClass
		StorageClassList
*/
package v1beta1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import strings "strings"
import reflect "reflect"
import github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func (m *StorageClass) Reset()                    { *m = StorageClass{} }
func (*StorageClass) ProtoMessage()               {}
func (*StorageClass) Descriptor() ([]byte, []int) { return fileDescriptorGenerated, []int{0} }

func (m *StorageClassList) Reset()                    { *m = StorageClassList{} }
func (*StorageClassList) ProtoMessage()               {}
func (*StorageClassList) Descriptor() ([]byte, []int) { return fileDescriptorGenerated, []int{1} }

func init() {
	proto.RegisterType((*StorageClass)(nil), "k8s.io.client-go.pkg.apis.storage.v1beta1.StorageClass")
	proto.RegisterType((*StorageClassList)(nil), "k8s.io.client-go.pkg.apis.storage.v1beta1.StorageClassList")
}
func (m *StorageClass) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StorageClass) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.ObjectMeta.Size()))
	n1, err := m.ObjectMeta.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Provisioner)))
	i += copy(dAtA[i:], m.Provisioner)
	if len(m.Parameters) > 0 {
		for k := range m.Parameters {
			dAtA[i] = 0x1a
			i++
			v := m.Parameters[k]
			mapSize := 1 + len(k) + sovGenerated(uint64(len(k))) + 1 + len(v) + sovGenerated(uint64(len(v)))
			i = encodeVarintGenerated(dAtA, i, uint64(mapSize))
			dAtA[i] = 0xa
			i++
			i = encodeVarintGenerated(dAtA, i, uint64(len(k)))
			i += copy(dAtA[i:], k)
			dAtA[i] = 0x12
			i++
			i = encodeVarintGenerated(dAtA, i, uint64(len(v)))
			i += copy(dAtA[i:], v)
		}
	}
	return i, nil
}

func (m *StorageClassList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StorageClassList) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.ListMeta.Size()))
	n2, err := m.ListMeta.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0x12
			i++
			i = encodeVarintGenerated(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeFixed64Generated(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Generated(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *StorageClass) Size() (n int) {
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Provisioner)
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Parameters) > 0 {
		for k, v := range m.Parameters {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovGenerated(uint64(len(k))) + 1 + len(v) + sovGenerated(uint64(len(v)))
			n += mapEntrySize + 1 + sovGenerated(uint64(mapEntrySize))
		}
	}
	return n
}

func (m *StorageClassList) Size() (n int) {
	var l int
	_ = l
	l = m.ListMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *StorageClass) String() string {
	if this == nil {
		return "nil"
	}
	keysForParameters := make([]string, 0, len(this.Parameters))
	for k := range this.Parameters {
		keysForParameters = append(keysForParameters, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForParameters)
	mapStringForParameters := "map[string]string{"
	for _, k := range keysForParameters {
		mapStringForParameters += fmt.Sprintf("%v: %v,", k, this.Parameters[k])
	}
	mapStringForParameters += "}"
	s := strings.Join([]string{`&StorageClass{`,
		`ObjectMeta:` + strings.Replace(strings.Replace(this.ObjectMeta.String(), "ObjectMeta", "k8s_io_apimachinery_pkg_apis_meta_v1.ObjectMeta", 1), `&`, ``, 1) + `,`,
		`Provisioner:` + fmt.Sprintf("%v", this.Provisioner) + `,`,
		`Parameters:` + mapStringForParameters + `,`,
		`}`,
	}, "")
	return s
}
func (this *StorageClassList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&StorageClassList{`,
		`ListMeta:` + strings.Replace(strings.Replace(this.ListMeta.String(), "ListMeta", "k8s_io_apimachinery_pkg_apis_meta_v1.ListMeta", 1), `&`, ``, 1) + `,`,
		`Items:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Items), "StorageClass", "StorageClass", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *StorageClass) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StorageClass: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StorageClass: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Provisioner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Provisioner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Parameters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var keykey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				keykey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			var stringLenmapkey uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLenmapkey |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLenmapkey := int(stringLenmapkey)
			if intStringLenmapkey < 0 {
				return ErrInvalidLengthGenerated
			}
			postStringIndexmapkey := iNdEx + intStringLenmapkey
			if postStringIndexmapkey > l {
				return io.ErrUnexpectedEOF
			}
			mapkey := string(dAtA[iNdEx:postStringIndexmapkey])
			iNdEx = postStringIndexmapkey
			if m.Parameters == nil {
				m.Parameters = make(map[string]string)
			}
			if iNdEx < postIndex {
				var valuekey uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					valuekey |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				var stringLenmapvalue uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGenerated
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLenmapvalue |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLenmapvalue := int(stringLenmapvalue)
				if intStringLenmapvalue < 0 {
					return ErrInvalidLengthGenerated
				}
				postStringIndexmapvalue := iNdEx + intStringLenmapvalue
				if postStringIndexmapvalue > l {
					return io.ErrUnexpectedEOF
				}
				mapvalue := string(dAtA[iNdEx:postStringIndexmapvalue])
				iNdEx = postStringIndexmapvalue
				m.Parameters[mapkey] = mapvalue
			} else {
				var mapvalue string
				m.Parameters[mapkey] = mapvalue
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *StorageClassList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: StorageClassList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StorageClassList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ListMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, StorageClass{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGenerated
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGenerated(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGenerated = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("k8s.io/client-go/pkg/apis/storage/v1beta1/generated.proto", fileDescriptorGenerated)
}

var fileDescriptorGenerated = []byte{
	// 479 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x4d, 0x8b, 0xd3, 0x40,
	0x1c, 0xc6, 0x33, 0x2d, 0xc5, 0xdd, 0xa9, 0x62, 0x89, 0x1e, 0x4a, 0x0f, 0xd9, 0xb2, 0xa7, 0x2a,
	0x38, 0x63, 0xd7, 0x17, 0xca, 0x82, 0x97, 0x8a, 0xa0, 0xa0, 0xb8, 0xc4, 0x9b, 0x28, 0x38, 0xe9,
	0xfe, 0x4d, 0xc7, 0x34, 0x99, 0x30, 0xf3, 0x4f, 0xb0, 0xe0, 0xc1, 0x8f, 0xe0, 0xc7, 0xea, 0xcd,
	0x3d, 0x7a, 0x5a, 0x6c, 0xf4, 0x83, 0x48, 0x5e, 0xdc, 0x04, 0xe3, 0x62, 0xf1, 0x96, 0x79, 0xf9,
	0x3d, 0xff, 0xe7, 0x79, 0x26, 0xf4, 0x38, 0x98, 0x19, 0x26, 0x15, 0x0f, 0x12, 0x0f, 0x74, 0x04,
	0x08, 0x86, 0xc7, 0x81, 0xcf, 0x45, 0x2c, 0x0d, 0x37, 0xa8, 0xb4, 0xf0, 0x81, 0xa7, 0x53, 0x0f,
	0x50, 0x4c, 0xb9, 0x0f, 0x11, 0x68, 0x81, 0x70, 0xca, 0x62, 0xad, 0x50, 0xd9, 0xb7, 0x4b, 0x96,
	0xd5, 0x2c, 0x8b, 0x03, 0x9f, 0xe5, 0x2c, 0xab, 0x58, 0x56, 0xb1, 0xa3, 0x3b, 0xbe, 0xc4, 0x65,
	0xe2, 0xb1, 0x85, 0x0a, 0xb9, 0xaf, 0x7c, 0xc5, 0x0b, 0x09, 0x2f, 0x79, 0x5f, 0xac, 0x8a, 0x45,
	0xf1, 0x55, 0x4a, 0x8f, 0xee, 0x57, 0xb6, 0x44, 0x2c, 0x43, 0xb1, 0x58, 0xca, 0x08, 0xf4, 0xba,
	0x36, 0x16, 0x02, 0x0a, 0x9e, 0xb6, 0x0c, 0x8d, 0xf8, 0x65, 0x94, 0x4e, 0x22, 0x94, 0x21, 0xb4,
	0x80, 0x87, 0xff, 0x02, 0xcc, 0x62, 0x09, 0xa1, 0x68, 0x71, 0xf7, 0x2e, 0xe3, 0x12, 0x94, 0x2b,
	0x2e, 0x23, 0x34, 0xa8, 0x5b, 0x50, 0x23, 0x93, 0x01, 0x9d, 0x82, 0xae, 0x03, 0xc1, 0x47, 0x11,
	0xc6, 0x2b, 0xf8, 0x4b, 0xa6, 0xc3, 0x9f, 0x1d, 0x7a, 0xf5, 0x55, 0x59, 0xe6, 0xe3, 0x95, 0x30,
	0xc6, 0x7e, 0x47, 0xf7, 0xf2, 0xfc, 0xa7, 0x02, 0xc5, 0x90, 0x8c, 0xc9, 0xa4, 0x7f, 0x74, 0x97,
	0x55, 0x0f, 0xd1, 0xb4, 0x53, 0x3f, 0x45, 0x7e, 0x9b, 0xa5, 0x53, 0xf6, 0xd2, 0xfb, 0x00, 0x0b,
	0x7c, 0x01, 0x28, 0xe6, 0xf6, 0xe6, 0xfc, 0xc0, 0xca, 0xce, 0x0f, 0x68, 0xbd, 0xe7, 0x5e, 0xa8,
	0xda, 0x0f, 0x68, 0x3f, 0xd6, 0x2a, 0x95, 0x46, 0xaa, 0x08, 0xf4, 0xb0, 0x33, 0x26, 0x93, 0xfd,
	0xf9, 0x8d, 0x0a, 0xe9, 0x9f, 0xd4, 0x47, 0x6e, 0xf3, 0x9e, 0xfd, 0x89, 0xd2, 0x58, 0x68, 0x11,
	0x02, 0x82, 0x36, 0xc3, 0xee, 0xb8, 0x3b, 0xe9, 0x1f, 0x3d, 0x65, 0xbb, 0xff, 0x23, 0xac, 0x19,
	0x93, 0x9d, 0x5c, 0x48, 0x3d, 0x89, 0x50, 0xaf, 0x6b, 0xcb, 0xf5, 0x81, 0xdb, 0x98, 0x37, 0x7a,
	0x44, 0xaf, 0xff, 0x81, 0xd8, 0x03, 0xda, 0x0d, 0x60, 0x5d, 0x94, 0xb4, 0xef, 0xe6, 0x9f, 0xf6,
	0x4d, 0xda, 0x4b, 0xc5, 0x2a, 0x81, 0x32, 0x93, 0x5b, 0x2e, 0x8e, 0x3b, 0x33, 0x72, 0xf8, 0x95,
	0xd0, 0x41, 0x73, 0xfe, 0x73, 0x69, 0xd0, 0x7e, 0xd3, 0xaa, 0x9a, 0xed, 0x56, 0x75, 0x4e, 0x17,
	0x45, 0x0f, 0x2a, 0xd7, 0x7b, 0xbf, 0x77, 0x1a, 0x35, 0xbf, 0xa5, 0x3d, 0x89, 0x10, 0x9a, 0x61,
	0xa7, 0xa8, 0x6a, 0xf6, 0xbf, 0x55, 0xcd, 0xaf, 0x55, 0x43, 0x7a, 0xcf, 0x72, 0x39, 0xb7, 0x54,
	0x9d, 0xdf, 0xda, 0x6c, 0x1d, 0xeb, 0x6c, 0xeb, 0x58, 0xdf, 0xb6, 0x8e, 0xf5, 0x39, 0x73, 0xc8,
	0x26, 0x73, 0xc8, 0x59, 0xe6, 0x90, 0xef, 0x99, 0x43, 0xbe, 0xfc, 0x70, 0xac, 0xd7, 0x57, 0x2a,
	0xb5, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7e, 0xb2, 0x39, 0x2a, 0x05, 0x04, 0x00, 0x00,
}
