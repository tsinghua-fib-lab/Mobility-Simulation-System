// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: wolong/routing/v2/routing.proto

package routingv2

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

type RouteType int32

const (
	RouteType_ROUTE_TYPE_UNSPECIFIED RouteType = 0
	RouteType_ROUTE_TYPE_DRIVING     RouteType = 1
	RouteType_ROUTE_TYPE_WALKING     RouteType = 2
	RouteType_ROUTE_TYPE_BY_BUS      RouteType = 3
)

// Enum value maps for RouteType.
var (
	RouteType_name = map[int32]string{
		0: "ROUTE_TYPE_UNSPECIFIED",
		1: "ROUTE_TYPE_DRIVING",
		2: "ROUTE_TYPE_WALKING",
		3: "ROUTE_TYPE_BY_BUS",
	}
	RouteType_value = map[string]int32{
		"ROUTE_TYPE_UNSPECIFIED": 0,
		"ROUTE_TYPE_DRIVING":     1,
		"ROUTE_TYPE_WALKING":     2,
		"ROUTE_TYPE_BY_BUS":      3,
	}
)

func (x RouteType) Enum() *RouteType {
	p := new(RouteType)
	*p = x
	return p
}

func (x RouteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RouteType) Descriptor() protoreflect.EnumDescriptor {
	return file_wolong_routing_v2_routing_proto_enumTypes[0].Descriptor()
}

func (RouteType) Type() protoreflect.EnumType {
	return &file_wolong_routing_v2_routing_proto_enumTypes[0]
}

func (x RouteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RouteType.Descriptor instead.
func (RouteType) EnumDescriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{0}
}

type JourneyType int32

const (
	JourneyType_JOURNEY_TYPE_UNSPECIFIED JourneyType = 0
	JourneyType_JOURNEY_TYPE_DRIVING     JourneyType = 1
	JourneyType_JOURNEY_TYPE_WALKING     JourneyType = 2
	JourneyType_JOURNEY_TYPE_BY_BUS      JourneyType = 3 // JOURNEY_TYPE_BY_TAXI = 4;
)

// Enum value maps for JourneyType.
var (
	JourneyType_name = map[int32]string{
		0: "JOURNEY_TYPE_UNSPECIFIED",
		1: "JOURNEY_TYPE_DRIVING",
		2: "JOURNEY_TYPE_WALKING",
		3: "JOURNEY_TYPE_BY_BUS",
	}
	JourneyType_value = map[string]int32{
		"JOURNEY_TYPE_UNSPECIFIED": 0,
		"JOURNEY_TYPE_DRIVING":     1,
		"JOURNEY_TYPE_WALKING":     2,
		"JOURNEY_TYPE_BY_BUS":      3,
	}
)

func (x JourneyType) Enum() *JourneyType {
	p := new(JourneyType)
	*p = x
	return p
}

func (x JourneyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (JourneyType) Descriptor() protoreflect.EnumDescriptor {
	return file_wolong_routing_v2_routing_proto_enumTypes[1].Descriptor()
}

func (JourneyType) Type() protoreflect.EnumType {
	return &file_wolong_routing_v2_routing_proto_enumTypes[1]
}

func (x JourneyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use JourneyType.Descriptor instead.
func (JourneyType) EnumDescriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{1}
}

// 描述车辆如何从当前车道进入下一车道
type NextLaneType int32

const (
	NextLaneType_NEXT_LANE_TYPE_UNSPECIFIED NextLaneType = 0
	// 直行进入下一车道（下一车道是当前车道的后继）
	NextLaneType_NEXT_LANE_TYPE_FORWARD NextLaneType = 1
	// 向左变道进入下一车道
	NextLaneType_NEXT_LANE_TYPE_LEFT NextLaneType = 2
	// 向右变道进入下一车道
	NextLaneType_NEXT_LANE_TYPE_RIGHT NextLaneType = 3
	// 当前车道是最后一个车道
	NextLaneType_NEXT_LANE_TYPE_LAST NextLaneType = 4
)

// Enum value maps for NextLaneType.
var (
	NextLaneType_name = map[int32]string{
		0: "NEXT_LANE_TYPE_UNSPECIFIED",
		1: "NEXT_LANE_TYPE_FORWARD",
		2: "NEXT_LANE_TYPE_LEFT",
		3: "NEXT_LANE_TYPE_RIGHT",
		4: "NEXT_LANE_TYPE_LAST",
	}
	NextLaneType_value = map[string]int32{
		"NEXT_LANE_TYPE_UNSPECIFIED": 0,
		"NEXT_LANE_TYPE_FORWARD":     1,
		"NEXT_LANE_TYPE_LEFT":        2,
		"NEXT_LANE_TYPE_RIGHT":       3,
		"NEXT_LANE_TYPE_LAST":        4,
	}
)

func (x NextLaneType) Enum() *NextLaneType {
	p := new(NextLaneType)
	*p = x
	return p
}

func (x NextLaneType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NextLaneType) Descriptor() protoreflect.EnumDescriptor {
	return file_wolong_routing_v2_routing_proto_enumTypes[2].Descriptor()
}

func (NextLaneType) Type() protoreflect.EnumType {
	return &file_wolong_routing_v2_routing_proto_enumTypes[2]
}

func (x NextLaneType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NextLaneType.Descriptor instead.
func (NextLaneType) EnumDescriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{2}
}

// 行人前进的方向与Lane的正方向（s增大的方向）的关系
type MovingDirection int32

const (
	MovingDirection_MOVING_DIRECTION_UNSPECIFIED MovingDirection = 0
	// 与正方向同向
	MovingDirection_MOVING_DIRECTION_FORWARD MovingDirection = 1
	// 与正方向反向
	MovingDirection_MOVING_DIRECTION_BACKWARD MovingDirection = 2
)

// Enum value maps for MovingDirection.
var (
	MovingDirection_name = map[int32]string{
		0: "MOVING_DIRECTION_UNSPECIFIED",
		1: "MOVING_DIRECTION_FORWARD",
		2: "MOVING_DIRECTION_BACKWARD",
	}
	MovingDirection_value = map[string]int32{
		"MOVING_DIRECTION_UNSPECIFIED": 0,
		"MOVING_DIRECTION_FORWARD":     1,
		"MOVING_DIRECTION_BACKWARD":    2,
	}
)

func (x MovingDirection) Enum() *MovingDirection {
	p := new(MovingDirection)
	*p = x
	return p
}

func (x MovingDirection) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MovingDirection) Descriptor() protoreflect.EnumDescriptor {
	return file_wolong_routing_v2_routing_proto_enumTypes[3].Descriptor()
}

func (MovingDirection) Type() protoreflect.EnumType {
	return &file_wolong_routing_v2_routing_proto_enumTypes[3]
}

func (x MovingDirection) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MovingDirection.Descriptor instead.
func (MovingDirection) EnumDescriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{3}
}

type DrivingRouteSegment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LaneId       int32        `protobuf:"varint,1,opt,name=lane_id,json=laneId,proto3" json:"lane_id,omitempty"`
	NextLaneType NextLaneType `protobuf:"varint,2,opt,name=next_lane_type,json=nextLaneType,proto3,enum=wolong.routing.v2.NextLaneType" json:"next_lane_type,omitempty"`
}

func (x *DrivingRouteSegment) Reset() {
	*x = DrivingRouteSegment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DrivingRouteSegment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DrivingRouteSegment) ProtoMessage() {}

func (x *DrivingRouteSegment) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DrivingRouteSegment.ProtoReflect.Descriptor instead.
func (*DrivingRouteSegment) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{0}
}

func (x *DrivingRouteSegment) GetLaneId() int32 {
	if x != nil {
		return x.LaneId
	}
	return 0
}

func (x *DrivingRouteSegment) GetNextLaneType() NextLaneType {
	if x != nil {
		return x.NextLaneType
	}
	return NextLaneType_NEXT_LANE_TYPE_UNSPECIFIED
}

// 车道序列
// 约定：智能体必须通过所有车道的末端（除了最后一个车道以及要求变道的车道）
type DrivingJourneyBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Route []*DrivingRouteSegment `protobuf:"bytes,1,rep,name=route,proto3" json:"route,omitempty"`
}

func (x *DrivingJourneyBody) Reset() {
	*x = DrivingJourneyBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DrivingJourneyBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DrivingJourneyBody) ProtoMessage() {}

func (x *DrivingJourneyBody) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DrivingJourneyBody.ProtoReflect.Descriptor instead.
func (*DrivingJourneyBody) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{1}
}

func (x *DrivingJourneyBody) GetRoute() []*DrivingRouteSegment {
	if x != nil {
		return x.Route
	}
	return nil
}

type WalkingRouteSegment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LaneId          int32           `protobuf:"varint,1,opt,name=lane_id,json=laneId,proto3" json:"lane_id,omitempty"`
	MovingDirection MovingDirection `protobuf:"varint,2,opt,name=moving_direction,json=movingDirection,proto3,enum=wolong.routing.v2.MovingDirection" json:"moving_direction,omitempty"`
}

func (x *WalkingRouteSegment) Reset() {
	*x = WalkingRouteSegment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WalkingRouteSegment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WalkingRouteSegment) ProtoMessage() {}

func (x *WalkingRouteSegment) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WalkingRouteSegment.ProtoReflect.Descriptor instead.
func (*WalkingRouteSegment) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{2}
}

func (x *WalkingRouteSegment) GetLaneId() int32 {
	if x != nil {
		return x.LaneId
	}
	return 0
}

func (x *WalkingRouteSegment) GetMovingDirection() MovingDirection {
	if x != nil {
		return x.MovingDirection
	}
	return MovingDirection_MOVING_DIRECTION_UNSPECIFIED
}

type WalkingJourneyBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Route []*WalkingRouteSegment `protobuf:"bytes,1,rep,name=route,proto3" json:"route,omitempty"`
}

func (x *WalkingJourneyBody) Reset() {
	*x = WalkingJourneyBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WalkingJourneyBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WalkingJourneyBody) ProtoMessage() {}

func (x *WalkingJourneyBody) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WalkingJourneyBody.ProtoReflect.Descriptor instead.
func (*WalkingJourneyBody) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{3}
}

func (x *WalkingJourneyBody) GetRoute() []*WalkingRouteSegment {
	if x != nil {
		return x.Route
	}
	return nil
}

type BusJourneyBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LineId         int32 `protobuf:"varint,1,opt,name=line_id,json=lineId,proto3" json:"line_id,omitempty"`
	StartStationId int32 `protobuf:"varint,2,opt,name=start_station_id,json=startStationId,proto3" json:"start_station_id,omitempty"`
	EndStationId   int32 `protobuf:"varint,3,opt,name=end_station_id,json=endStationId,proto3" json:"end_station_id,omitempty"`
}

func (x *BusJourneyBody) Reset() {
	*x = BusJourneyBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BusJourneyBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BusJourneyBody) ProtoMessage() {}

func (x *BusJourneyBody) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BusJourneyBody.ProtoReflect.Descriptor instead.
func (*BusJourneyBody) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{4}
}

func (x *BusJourneyBody) GetLineId() int32 {
	if x != nil {
		return x.LineId
	}
	return 0
}

func (x *BusJourneyBody) GetStartStationId() int32 {
	if x != nil {
		return x.StartStationId
	}
	return 0
}

func (x *BusJourneyBody) GetEndStationId() int32 {
	if x != nil {
		return x.EndStationId
	}
	return 0
}

type TaxiJourneyBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TaxiJourneyBody) Reset() {
	*x = TaxiJourneyBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaxiJourneyBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaxiJourneyBody) ProtoMessage() {}

func (x *TaxiJourneyBody) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaxiJourneyBody.ProtoReflect.Descriptor instead.
func (*TaxiJourneyBody) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{5}
}

// 路径规划结果的一部分，含且仅含采用一种交通出行方式的完整出行序列
type Journey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    JourneyType         `protobuf:"varint,1,opt,name=type,proto3,enum=wolong.routing.v2.JourneyType" json:"type,omitempty"`
	Driving *DrivingJourneyBody `protobuf:"bytes,2,opt,name=driving,proto3,oneof" json:"driving,omitempty"`
	Walking *WalkingJourneyBody `protobuf:"bytes,3,opt,name=walking,proto3,oneof" json:"walking,omitempty"`
	ByBus   *BusJourneyBody     `protobuf:"bytes,4,opt,name=by_bus,json=byBus,proto3,oneof" json:"by_bus,omitempty"` // optional TaxiJourneyBody by_taxi = 5;
}

func (x *Journey) Reset() {
	*x = Journey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wolong_routing_v2_routing_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Journey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Journey) ProtoMessage() {}

func (x *Journey) ProtoReflect() protoreflect.Message {
	mi := &file_wolong_routing_v2_routing_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Journey.ProtoReflect.Descriptor instead.
func (*Journey) Descriptor() ([]byte, []int) {
	return file_wolong_routing_v2_routing_proto_rawDescGZIP(), []int{6}
}

func (x *Journey) GetType() JourneyType {
	if x != nil {
		return x.Type
	}
	return JourneyType_JOURNEY_TYPE_UNSPECIFIED
}

func (x *Journey) GetDriving() *DrivingJourneyBody {
	if x != nil {
		return x.Driving
	}
	return nil
}

func (x *Journey) GetWalking() *WalkingJourneyBody {
	if x != nil {
		return x.Walking
	}
	return nil
}

func (x *Journey) GetByBus() *BusJourneyBody {
	if x != nil {
		return x.ByBus
	}
	return nil
}

var File_wolong_routing_v2_routing_proto protoreflect.FileDescriptor

var file_wolong_routing_v2_routing_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x2f, 0x76, 0x32, 0x2f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x11, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x76, 0x32, 0x22, 0x75, 0x0a, 0x13, 0x44, 0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6c,
	0x61, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x61,
	0x6e, 0x65, 0x49, 0x64, 0x12, 0x45, 0x0a, 0x0e, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x6c, 0x61, 0x6e,
	0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x77,
	0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32,
	0x2e, 0x4e, 0x65, 0x78, 0x74, 0x4c, 0x61, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x6e,
	0x65, 0x78, 0x74, 0x4c, 0x61, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x52, 0x0a, 0x12, 0x44,
	0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x42, 0x6f, 0x64,
	0x79, 0x12, 0x3c, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x26, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x2e, 0x76, 0x32, 0x2e, 0x44, 0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x22,
	0x7d, 0x0a, 0x13, 0x57, 0x61, 0x6c, 0x6b, 0x69, 0x6e, 0x67, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x53,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x61, 0x6e, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x61, 0x6e, 0x65, 0x49, 0x64, 0x12,
	0x4d, 0x0a, 0x10, 0x6d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x5f, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x77, 0x6f, 0x6c, 0x6f,
	0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32, 0x2e, 0x4d, 0x6f,
	0x76, 0x69, 0x6e, 0x67, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x6d,
	0x6f, 0x76, 0x69, 0x6e, 0x67, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x52,
	0x0a, 0x12, 0x57, 0x61, 0x6c, 0x6b, 0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x3c, 0x0a, 0x05, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32, 0x2e, 0x57, 0x61, 0x6c, 0x6b, 0x69, 0x6e, 0x67, 0x52,
	0x6f, 0x75, 0x74, 0x65, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x72, 0x6f, 0x75,
	0x74, 0x65, 0x22, 0x79, 0x0a, 0x0e, 0x42, 0x75, 0x73, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79,
	0x42, 0x6f, 0x64, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x69, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x69, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x28, 0x0a,
	0x10, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x65, 0x6e, 0x64, 0x5f, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x65, 0x6e, 0x64, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x11, 0x0a,
	0x0f, 0x54, 0x61, 0x78, 0x69, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x42, 0x6f, 0x64, 0x79,
	0x22, 0xab, 0x02, 0x0a, 0x07, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x12, 0x32, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x77, 0x6f, 0x6c,
	0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32, 0x2e, 0x4a,
	0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x44, 0x0a, 0x07, 0x64, 0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x25, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x2e, 0x76, 0x32, 0x2e, 0x44, 0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x75,
	0x72, 0x6e, 0x65, 0x79, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x00, 0x52, 0x07, 0x64, 0x72, 0x69, 0x76,
	0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x44, 0x0a, 0x07, 0x77, 0x61, 0x6c, 0x6b, 0x69, 0x6e,
	0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67,
	0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32, 0x2e, 0x57, 0x61, 0x6c, 0x6b,
	0x69, 0x6e, 0x67, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x42, 0x6f, 0x64, 0x79, 0x48, 0x01,
	0x52, 0x07, 0x77, 0x61, 0x6c, 0x6b, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x3d, 0x0a, 0x06,
	0x62, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x77,
	0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32,
	0x2e, 0x42, 0x75, 0x73, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x42, 0x6f, 0x64, 0x79, 0x48,
	0x02, 0x52, 0x05, 0x62, 0x79, 0x42, 0x75, 0x73, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f,
	0x64, 0x72, 0x69, 0x76, 0x69, 0x6e, 0x67, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x77, 0x61, 0x6c, 0x6b,
	0x69, 0x6e, 0x67, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x62, 0x79, 0x5f, 0x62, 0x75, 0x73, 0x2a, 0x6e,
	0x0a, 0x09, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x16, 0x52,
	0x4f, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43,
	0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x4f, 0x55, 0x54, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x52, 0x49, 0x56, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x16, 0x0a, 0x12, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x57, 0x41,
	0x4c, 0x4b, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x52, 0x4f, 0x55, 0x54, 0x45,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x59, 0x5f, 0x42, 0x55, 0x53, 0x10, 0x03, 0x2a, 0x78,
	0x0a, 0x0b, 0x4a, 0x6f, 0x75, 0x72, 0x6e, 0x65, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x18, 0x4a, 0x4f, 0x55, 0x52, 0x4e, 0x45, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e,
	0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x4a,
	0x4f, 0x55, 0x52, 0x4e, 0x45, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x52, 0x49, 0x56,
	0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x4a, 0x4f, 0x55, 0x52, 0x4e, 0x45, 0x59,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x57, 0x41, 0x4c, 0x4b, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12,
	0x17, 0x0a, 0x13, 0x4a, 0x4f, 0x55, 0x52, 0x4e, 0x45, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x42, 0x59, 0x5f, 0x42, 0x55, 0x53, 0x10, 0x03, 0x2a, 0x96, 0x01, 0x0a, 0x0c, 0x4e, 0x65, 0x78,
	0x74, 0x4c, 0x61, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1e, 0x0a, 0x1a, 0x4e, 0x45, 0x58,
	0x54, 0x5f, 0x4c, 0x41, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x4e, 0x45, 0x58,
	0x54, 0x5f, 0x4c, 0x41, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x4f, 0x52, 0x57,
	0x41, 0x52, 0x44, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x4c, 0x41,
	0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x45, 0x46, 0x54, 0x10, 0x02, 0x12, 0x18,
	0x0a, 0x14, 0x4e, 0x45, 0x58, 0x54, 0x5f, 0x4c, 0x41, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x4e, 0x45, 0x58, 0x54,
	0x5f, 0x4c, 0x41, 0x4e, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x41, 0x53, 0x54, 0x10,
	0x04, 0x2a, 0x70, 0x0a, 0x0f, 0x4d, 0x6f, 0x76, 0x69, 0x6e, 0x67, 0x44, 0x69, 0x72, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x1c, 0x4d, 0x4f, 0x56, 0x49, 0x4e, 0x47, 0x5f, 0x44,
	0x49, 0x52, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x4d, 0x4f, 0x56, 0x49, 0x4e, 0x47,
	0x5f, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x4f, 0x52, 0x57, 0x41,
	0x52, 0x44, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x4d, 0x4f, 0x56, 0x49, 0x4e, 0x47, 0x5f, 0x44,
	0x49, 0x52, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x57, 0x41, 0x52,
	0x44, 0x10, 0x02, 0x42, 0xd3, 0x01, 0x0a, 0x15, 0x63, 0x6f, 0x6d, 0x2e, 0x77, 0x6f, 0x6c, 0x6f,
	0x6e, 0x67, 0x2e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x32, 0x42, 0x0c, 0x52,
	0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x46, 0x67,
	0x69, 0x74, 0x2e, 0x66, 0x69, 0x62, 0x6c, 0x61, 0x62, 0x2e, 0x6e, 0x65, 0x74, 0x2f, 0x73, 0x69,
	0x6d, 0x2f, 0x73, 0x69, 0x6d, 0x75, 0x6c, 0x65, 0x74, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x77, 0x6f, 0x6c, 0x6f, 0x6e, 0x67,
	0x2f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x32, 0x3b, 0x72, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x67, 0x76, 0x32, 0xa2, 0x02, 0x03, 0x57, 0x52, 0x58, 0xaa, 0x02, 0x11, 0x57, 0x6f,
	0x6c, 0x6f, 0x6e, 0x67, 0x2e, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x32, 0xca,
	0x02, 0x11, 0x57, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x5c, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x5c, 0x56, 0x32, 0xe2, 0x02, 0x1d, 0x57, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x5c, 0x52, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x5c, 0x56, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x13, 0x57, 0x6f, 0x6c, 0x6f, 0x6e, 0x67, 0x3a, 0x3a, 0x52, 0x6f,
	0x75, 0x74, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x56, 0x32, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_wolong_routing_v2_routing_proto_rawDescOnce sync.Once
	file_wolong_routing_v2_routing_proto_rawDescData = file_wolong_routing_v2_routing_proto_rawDesc
)

func file_wolong_routing_v2_routing_proto_rawDescGZIP() []byte {
	file_wolong_routing_v2_routing_proto_rawDescOnce.Do(func() {
		file_wolong_routing_v2_routing_proto_rawDescData = protoimpl.X.CompressGZIP(file_wolong_routing_v2_routing_proto_rawDescData)
	})
	return file_wolong_routing_v2_routing_proto_rawDescData
}

var file_wolong_routing_v2_routing_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_wolong_routing_v2_routing_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_wolong_routing_v2_routing_proto_goTypes = []interface{}{
	(RouteType)(0),              // 0: wolong.routing.v2.RouteType
	(JourneyType)(0),            // 1: wolong.routing.v2.JourneyType
	(NextLaneType)(0),           // 2: wolong.routing.v2.NextLaneType
	(MovingDirection)(0),        // 3: wolong.routing.v2.MovingDirection
	(*DrivingRouteSegment)(nil), // 4: wolong.routing.v2.DrivingRouteSegment
	(*DrivingJourneyBody)(nil),  // 5: wolong.routing.v2.DrivingJourneyBody
	(*WalkingRouteSegment)(nil), // 6: wolong.routing.v2.WalkingRouteSegment
	(*WalkingJourneyBody)(nil),  // 7: wolong.routing.v2.WalkingJourneyBody
	(*BusJourneyBody)(nil),      // 8: wolong.routing.v2.BusJourneyBody
	(*TaxiJourneyBody)(nil),     // 9: wolong.routing.v2.TaxiJourneyBody
	(*Journey)(nil),             // 10: wolong.routing.v2.Journey
}
var file_wolong_routing_v2_routing_proto_depIdxs = []int32{
	2, // 0: wolong.routing.v2.DrivingRouteSegment.next_lane_type:type_name -> wolong.routing.v2.NextLaneType
	4, // 1: wolong.routing.v2.DrivingJourneyBody.route:type_name -> wolong.routing.v2.DrivingRouteSegment
	3, // 2: wolong.routing.v2.WalkingRouteSegment.moving_direction:type_name -> wolong.routing.v2.MovingDirection
	6, // 3: wolong.routing.v2.WalkingJourneyBody.route:type_name -> wolong.routing.v2.WalkingRouteSegment
	1, // 4: wolong.routing.v2.Journey.type:type_name -> wolong.routing.v2.JourneyType
	5, // 5: wolong.routing.v2.Journey.driving:type_name -> wolong.routing.v2.DrivingJourneyBody
	7, // 6: wolong.routing.v2.Journey.walking:type_name -> wolong.routing.v2.WalkingJourneyBody
	8, // 7: wolong.routing.v2.Journey.by_bus:type_name -> wolong.routing.v2.BusJourneyBody
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_wolong_routing_v2_routing_proto_init() }
func file_wolong_routing_v2_routing_proto_init() {
	if File_wolong_routing_v2_routing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wolong_routing_v2_routing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DrivingRouteSegment); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DrivingJourneyBody); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WalkingRouteSegment); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WalkingJourneyBody); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BusJourneyBody); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaxiJourneyBody); i {
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
		file_wolong_routing_v2_routing_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Journey); i {
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
	file_wolong_routing_v2_routing_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wolong_routing_v2_routing_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wolong_routing_v2_routing_proto_goTypes,
		DependencyIndexes: file_wolong_routing_v2_routing_proto_depIdxs,
		EnumInfos:         file_wolong_routing_v2_routing_proto_enumTypes,
		MessageInfos:      file_wolong_routing_v2_routing_proto_msgTypes,
	}.Build()
	File_wolong_routing_v2_routing_proto = out.File
	file_wolong_routing_v2_routing_proto_rawDesc = nil
	file_wolong_routing_v2_routing_proto_goTypes = nil
	file_wolong_routing_v2_routing_proto_depIdxs = nil
}