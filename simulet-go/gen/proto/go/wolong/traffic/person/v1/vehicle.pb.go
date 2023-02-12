// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: wolong/traffic/person/v1/vehicle.proto

package personv1

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

type Reason int32

const (
	Reason_REASON_UNSPECIFIED     Reason = 0
	Reason_REASON_TO_SELF_LIMIT   Reason = 1
	Reason_REASON_CAR_FOLLOW      Reason = 2
	Reason_REASON_TO_LANE_LIMIT   Reason = 3
	Reason_REASON_RESTRICTION     Reason = 4
	Reason_REASON_RED_LIGHT       Reason = 5
	Reason_REASON_YELLOW_LIGHT    Reason = 6
	Reason_REASON_END             Reason = 7
	Reason_REASON_JO_WALKING      Reason = 8
	Reason_REASON_JO_STOPCAR      Reason = 9
	Reason_REASON_JO_OCCUPANCY    Reason = 10
	Reason_REASON_JO_INTERSECTION Reason = 11
	Reason_REASON_JO_CLOSESTART   Reason = 12
	Reason_REASON_LC_AHEAD        Reason = 13
	Reason_REASON_LC_BEHIND       Reason = 14
	// others (O)
	Reason_REASON_O_RING Reason = 15
)

// Enum value maps for Reason.
var (
	Reason_name = map[int32]string{
		0:  "REASON_UNSPECIFIED",
		1:  "REASON_TO_SELF_LIMIT",
		2:  "REASON_CAR_FOLLOW",
		3:  "REASON_TO_LANE_LIMIT",
		4:  "REASON_RESTRICTION",
		5:  "REASON_RED_LIGHT",
		6:  "REASON_YELLOW_LIGHT",
		7:  "REASON_END",
		8:  "REASON_JO_WALKING",
		9:  "REASON_JO_STOPCAR",
		10: "REASON_JO_OCCUPANCY",
		11: "REASON_JO_INTERSECTION",
		12: "REASON_JO_CLOSESTART",
		13: "REASON_LC_AHEAD",
		14: "REASON_LC_BEHIND",
		15: "REASON_O_RING",
	}
	Reason_value = map[string]int32{
		"REASON_UNSPECIFIED":     0,
		"REASON_TO_SELF_LIMIT":   1,
		"REASON_CAR_FOLLOW":      2,
		"REASON_TO_LANE_LIMIT":   3,
		"REASON_RESTRICTION":     4,
		"REASON_RED_LIGHT":       5,
		"REASON_YELLOW_LIGHT":    6,
		"REASON_END":             7,
		"REASON_JO_WALKING":      8,
		"REASON_JO_STOPCAR":      9,
		"REASON_JO_OCCUPANCY":    10,
		"REASON_JO_INTERSECTION": 11,
		"REASON_JO_CLOSESTART":   12,
		"REASON_LC_AHEAD":        13,
		"REASON_LC_BEHIND":       14,
		"REASON_O_RING":          15,
	}
)

func (x Reason) Enum() *Reason {
	p := new(Reason)
	*p = x
	return p
}

func (x Reason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Reason) Descriptor() protoreflect.EnumDescriptor {
	return file_wolong_traffic_person_v1_vehicle_proto_enumTypes[0].Descriptor()
}

func (Reason) Type() protoreflect.EnumType {
	return &file_wolong_traffic_person_v1_vehicle_proto_enumTypes[0]
}

func (x Reason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Reason.Descriptor instead.
func (Reason) EnumDescriptor() ([]byte, []int) {
	return file_wolong_traffic_person_v1_vehicle_proto_rawDescGZIP(), []int{0}
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Acc              float64 `protobuf:"fixed64,1,opt,name=acc,proto3" json:"acc,omitempty"`
	EnableLaneChange bool    `protobuf:"varint,2,opt,name=enable_lane_change,json=enableLaneChange,proto3" json:"enable_lane_change,omitempty"`
	LaneChangeLength float64 `protobuf:"fixed64,3,opt,name=lane_change_length,json=laneChangeLength,proto3" json:"lane_change_length,omitempty"`
	Reason           Reason  `protobuf:"varint,4,opt,name=reason,proto3,enum=wolong.traffic.person.v1.Reason" json:"reason,omitempty"`
	RelatedId        int32   `protobuf:"varint,5,opt,name=related_id,json=relatedId,proto3" json:"related_id,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_traffic_person_v1_vehicle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_traffic_person_v1_vehicle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_wolong_traffic_person_v1_vehicle_proto_rawDescGZIP(), []int{0}
}

func (x *Action) GetAcc() float64 {
	if x != nil {
		return x.Acc
	}
	return 0
}

func (x *Action) GetEnableLaneChange() bool {
	if x != nil {
		return x.EnableLaneChange
	}
	return false
}

func (x *Action) GetLaneChangeLength() float64 {
	if x != nil {
		return x.LaneChangeLength
	}
	return 0
}

func (x *Action) GetReason() Reason {
	if x != nil {
		return x.Reason
	}
	return Reason_REASON_UNSPECIFIED
}

func (x *Action) GetRelatedId() int32 {
	if x != nil {
		return x.RelatedId
	}
	return 0
}

type Vehicle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                        int32              `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Base                      *BaseRuntime       `protobuf:"bytes,2,opt,name=base,proto3" json:"base,omitempty"`
	BaseOnRoad                *BaseRuntimeOnRoad `protobuf:"bytes,3,opt,name=base_on_road,json=baseOnRoad,proto3" json:"base_on_road,omitempty"`
	Action                    *Action            `protobuf:"bytes,4,opt,name=action,proto3" json:"action,omitempty"`
	DistanceToEnd             float64            `protobuf:"fixed64,5,opt,name=distance_to_end,json=distanceToEnd,proto3" json:"distance_to_end,omitempty"`
	ShadowLaneId              int32              `protobuf:"varint,6,opt,name=shadow_lane_id,json=shadowLaneId,proto3" json:"shadow_lane_id,omitempty"`
	ShadowS                   float64            `protobuf:"fixed64,7,opt,name=shadow_s,json=shadowS,proto3" json:"shadow_s,omitempty"`
	LaneChangeTotalLength     float64            `protobuf:"fixed64,8,opt,name=lane_change_total_length,json=laneChangeTotalLength,proto3" json:"lane_change_total_length,omitempty"`
	LaneChangeCompletedLength float64            `protobuf:"fixed64,9,opt,name=lane_change_completed_length,json=laneChangeCompletedLength,proto3" json:"lane_change_completed_length,omitempty"`
	IsLaneChanging            bool               `protobuf:"varint,10,opt,name=is_lane_changing,json=isLaneChanging,proto3" json:"is_lane_changing,omitempty"`
}

func (x *Vehicle) Reset() {
	*x = Vehicle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_traffic_person_v1_vehicle_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vehicle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vehicle) ProtoMessage() {}

func (x *Vehicle) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_traffic_person_v1_vehicle_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vehicle.ProtoReflect.Descriptor instead.
func (*Vehicle) Descriptor() ([]byte, []int) {
	return file_wolong_traffic_person_v1_vehicle_proto_rawDescGZIP(), []int{1}
}

func (x *Vehicle) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Vehicle) GetBase() *BaseRuntime {
	if x != nil {
		return x.Base
	}
	return nil
}

func (x *Vehicle) GetBaseOnRoad() *BaseRuntimeOnRoad {
	if x != nil {
		return x.BaseOnRoad
	}
	return nil
}

func (x *Vehicle) GetAction() *Action {
	if x != nil {
		return x.Action
	}
	return nil
}

func (x *Vehicle) GetDistanceToEnd() float64 {
	if x != nil {
		return x.DistanceToEnd
	}
	return 0
}

func (x *Vehicle) GetShadowLaneId() int32 {
	if x != nil {
		return x.ShadowLaneId
	}
	return 0
}

func (x *Vehicle) GetShadowS() float64 {
	if x != nil {
		return x.ShadowS
	}
	return 0
}

func (x *Vehicle) GetLaneChangeTotalLength() float64 {
	if x != nil {
		return x.LaneChangeTotalLength
	}
	return 0
}

func (x *Vehicle) GetLaneChangeCompletedLength() float64 {
	if x != nil {
		return x.LaneChangeCompletedLength
	}
	return 0
}

func (x *Vehicle) GetIsLaneChanging() bool {
	if x != nil {
		return x.IsLaneChanging
	}
	return false
}

var File_wolong_traffic_person_v1_vehicle_proto protoreflect.FileDescriptor

var file_wolong_traffic_person_v1_vehicle_proto_rawDesc = []byte{
	0x0a, 0x26, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63,
	0x2f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x68, 0x69, 0x63,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67,
	0x2e, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e,
	0x76, 0x31, 0x1a, 0x23, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x61, 0x66, 0x66,
	0x69, 0x63, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x61, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x01, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x63, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x03, 0x61, 0x63, 0x63, 0x12, 0x2c, 0x0a, 0x12, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6c,
	0x61, 0x6e, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x10, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x4c, 0x61, 0x6e, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x12, 0x2c, 0x0a, 0x12, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10,
	0x6c, 0x61, 0x6e, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68,
	0x12, 0x38, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x20, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69,
	0x63, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x73,
	0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65,
	0x6c, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x65, 0x64, 0x49, 0x64, 0x22, 0xea, 0x03, 0x0a, 0x07, 0x56, 0x65,
	0x68, 0x69, 0x63, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x61,
	0x66, 0x66, 0x69, 0x63, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x42,
	0x61, 0x73, 0x65, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65,
	0x12, 0x4d, 0x0a, 0x0c, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x6f, 0x6e, 0x5f, 0x72, 0x6f, 0x61, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e,
	0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x4f, 0x6e, 0x52,
	0x6f, 0x61, 0x64, 0x52, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x4f, 0x6e, 0x52, 0x6f, 0x61, 0x64, 0x12,
	0x38, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63,
	0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x0f, 0x64, 0x69, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x74, 0x6f, 0x5f, 0x65, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x0d, 0x64, 0x69, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x54, 0x6f, 0x45, 0x6e,
	0x64, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x68, 0x61, 0x64, 0x6f, 0x77, 0x5f, 0x6c, 0x61, 0x6e, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x73, 0x68, 0x61, 0x64, 0x6f,
	0x77, 0x4c, 0x61, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x64, 0x6f,
	0x77, 0x5f, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x73, 0x68, 0x61, 0x64, 0x6f,
	0x77, 0x53, 0x12, 0x37, 0x0a, 0x18, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x5f, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x15, 0x6c, 0x61, 0x6e, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x3f, 0x0a, 0x1c, 0x6c,
	0x61, 0x6e, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c,
	0x65, 0x74, 0x65, 0x64, 0x5f, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x19, 0x6c, 0x61, 0x6e, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x28, 0x0a, 0x10,
	0x69, 0x73, 0x5f, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x69, 0x73, 0x4c, 0x61, 0x6e, 0x65, 0x43, 0x68,
	0x61, 0x6e, 0x67, 0x69, 0x6e, 0x67, 0x2a, 0xfd, 0x02, 0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45, 0x41,
	0x53, 0x4f, 0x4e, 0x5f, 0x54, 0x4f, 0x5f, 0x53, 0x45, 0x4c, 0x46, 0x5f, 0x4c, 0x49, 0x4d, 0x49,
	0x54, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x43, 0x41,
	0x52, 0x5f, 0x46, 0x4f, 0x4c, 0x4c, 0x4f, 0x57, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x45,
	0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x54, 0x4f, 0x5f, 0x4c, 0x41, 0x4e, 0x45, 0x5f, 0x4c, 0x49, 0x4d,
	0x49, 0x54, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x52,
	0x45, 0x53, 0x54, 0x52, 0x49, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x04, 0x12, 0x14, 0x0a, 0x10,
	0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x52, 0x45, 0x44, 0x5f, 0x4c, 0x49, 0x47, 0x48, 0x54,
	0x10, 0x05, 0x12, 0x17, 0x0a, 0x13, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x59, 0x45, 0x4c,
	0x4c, 0x4f, 0x57, 0x5f, 0x4c, 0x49, 0x47, 0x48, 0x54, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x52,
	0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x45, 0x4e, 0x44, 0x10, 0x07, 0x12, 0x15, 0x0a, 0x11, 0x52,
	0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x5f, 0x57, 0x41, 0x4c, 0x4b, 0x49, 0x4e, 0x47,
	0x10, 0x08, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x5f,
	0x53, 0x54, 0x4f, 0x50, 0x43, 0x41, 0x52, 0x10, 0x09, 0x12, 0x17, 0x0a, 0x13, 0x52, 0x45, 0x41,
	0x53, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x5f, 0x4f, 0x43, 0x43, 0x55, 0x50, 0x41, 0x4e, 0x43, 0x59,
	0x10, 0x0a, 0x12, 0x1a, 0x0a, 0x16, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x5f,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x53, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0b, 0x12, 0x18,
	0x0a, 0x14, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4a, 0x4f, 0x5f, 0x43, 0x4c, 0x4f, 0x53,
	0x45, 0x53, 0x54, 0x41, 0x52, 0x54, 0x10, 0x0c, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x41, 0x53,
	0x4f, 0x4e, 0x5f, 0x4c, 0x43, 0x5f, 0x41, 0x48, 0x45, 0x41, 0x44, 0x10, 0x0d, 0x12, 0x14, 0x0a,
	0x10, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4c, 0x43, 0x5f, 0x42, 0x45, 0x48, 0x49, 0x4e,
	0x44, 0x10, 0x0e, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x4f, 0x5f,
	0x52, 0x49, 0x4e, 0x47, 0x10, 0x0f, 0x42, 0xfd, 0x01, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x77,
	0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2e, 0x70, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x56, 0x65, 0x68, 0x69, 0x63, 0x6c, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x2e, 0x66, 0x69, 0x62,
	0x6c, 0x61, 0x62, 0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x69, 0x6d, 0x2f, 0x73, 0x69, 0x6d, 0x75,
	0x6c, 0x65, 0x74, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x61, 0x66, 0x66,
	0x69, 0x63, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x65, 0x72,
	0x73, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x57, 0x54, 0x50, 0xaa, 0x02, 0x18, 0x57, 0x6f,
	0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x2e, 0x50, 0x65, 0x72,
	0x73, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x18, 0x57, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x5c,
	0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x5c, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x5c, 0x56,
	0x31, 0xe2, 0x02, 0x24, 0x57, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x5c, 0x54, 0x72, 0x61, 0x66, 0x66,
	0x69, 0x63, 0x5c, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x1b, 0x57, 0x6f, 0x6c, 0x6f, 0x6e,
	0x67, 0x3a, 0x3a, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x3a, 0x3a, 0x50, 0x65, 0x72, 0x73,
	0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wolong_traffic_person_v1_vehicle_proto_rawDescOnce sync.Once
	file_wolong_traffic_person_v1_vehicle_proto_rawDescData = file_wolong_traffic_person_v1_vehicle_proto_rawDesc
)

func file_wolong_traffic_person_v1_vehicle_proto_rawDescGZIP() []byte {
	file_wolong_traffic_person_v1_vehicle_proto_rawDescOnce.Do(func() {
		file_wolong_traffic_person_v1_vehicle_proto_rawDescData = protoimpl.X.CompressGZIP(file_wolong_traffic_person_v1_vehicle_proto_rawDescData)
	})
	return file_wolong_traffic_person_v1_vehicle_proto_rawDescData
}

var file_wolong_traffic_person_v1_vehicle_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_wolong_traffic_person_v1_vehicle_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_wolong_traffic_person_v1_vehicle_proto_goTypes = []interface{}{
	(Reason)(0),               // 0: wolong.traffic.person.v1.Reason
	(*Action)(nil),            // 1: wolong.traffic.person.v1.Action
	(*Vehicle)(nil),           // 2: wolong.traffic.person.v1.Vehicle
	(*BaseRuntime)(nil),       // 3: wolong.traffic.person.v1.BaseRuntime
	(*BaseRuntimeOnRoad)(nil), // 4: wolong.traffic.person.v1.BaseRuntimeOnRoad
}
var file_wolong_traffic_person_v1_vehicle_proto_depIdxs = []int32{
	0, // 0: wolong.traffic.person.v1.Action.reason:type_name -> wolong.traffic.person.v1.Reason
	3, // 1: wolong.traffic.person.v1.Vehicle.base:type_name -> wolong.traffic.person.v1.BaseRuntime
	4, // 2: wolong.traffic.person.v1.Vehicle.base_on_road:type_name -> wolong.traffic.person.v1.BaseRuntimeOnRoad
	1, // 3: wolong.traffic.person.v1.Vehicle.action:type_name -> wolong.traffic.person.v1.Action
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_wolong_traffic_person_v1_vehicle_proto_init() }
func file_wolong_traffic_person_v1_vehicle_proto_init() {
	if File_wolong_traffic_person_v1_vehicle_proto != nil {
		return
	}
	file_wolong_traffic_person_v1_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_wolong_traffic_person_v1_vehicle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
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
		file_wolong_traffic_person_v1_vehicle_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vehicle); i {
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
			RawDescriptor: file_wolong_traffic_person_v1_vehicle_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wolong_traffic_person_v1_vehicle_proto_goTypes,
		DependencyIndexes: file_wolong_traffic_person_v1_vehicle_proto_depIdxs,
		EnumInfos:         file_wolong_traffic_person_v1_vehicle_proto_enumTypes,
		MessageInfos:      file_wolong_traffic_person_v1_vehicle_proto_msgTypes,
	}.Build()
	File_wolong_traffic_person_v1_vehicle_proto = out.File
	file_wolong_traffic_person_v1_vehicle_proto_rawDesc = nil
	file_wolong_traffic_person_v1_vehicle_proto_goTypes = nil
	file_wolong_traffic_person_v1_vehicle_proto_depIdxs = nil
}
