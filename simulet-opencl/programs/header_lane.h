#ifndef DEVICELANE_H
#define DEVICELANE_H

// Lane结构体中冲突点上限
#define kMaxNextLanes 64
// Road中前驱后继lane数量上限
#define kMaxOverlaps 100
//road包含的最大车道数量
#define MaxLanes 16
typedef struct Road {
  bool valid;
  int id;
  // 属于该道路Road的所有车道/人行道等lane
  // lane_id是按从最左侧车道到最右侧车道(从前进方向来看)的顺序给出的
  int size;
  int lane_ids[MaxLanes];
}Road;

/**Lane相关数据结构 */
typedef struct InsertInfo{
    int vehicleIndex;
    int s;
} InsertInfo;

typedef struct HostNextLane {
  int index;
}HostNextLane;

typedef struct LanePosition{
  int index;
  double s;
} LanePosition;

typedef struct Overlap {
    //冲突点在本车道坐标
    struct LanePosition self;
     //冲突点在目标车道坐标
    struct LanePosition other_s;
    /// 是否本车道优先
    bool self_first;
} Overlap;

typedef struct HostLane {
    bool valid;
    int index;
    int roadID;
    float length , max_speed;
    int left_lane_index;
    int right_lane_index;
    int next_lane_size;
    struct HostNextLane next_lanes[kMaxNextLanes]; //kMaxRoadNextLanes  Road中前驱后继lane数量上限

    // for junction lane
    int  overlap_size;
    struct Overlap  overlaps[kMaxOverlaps];
} HostLane;

typedef struct DeviceLanePosition{
  struct DeviceLane *lane;
  double s;
} DeviceLanePosition;

typedef struct DeviceOverlap {
   //冲突点在本车道坐标
    struct DeviceLanePosition self;
     //冲突点在目标车道坐标
    struct DeviceLanePosition other_s;
    /// 是否本车道优先
    bool self_first;
} DeviceOverlap;
typedef struct DeviceLane {  
    bool valid;
    int index;
    int roadID;
    float length , max_speed;
    struct DeviceLane *left_lane;
    struct DeviceLane *right_lane;
    int next_lane_size;//后继车道数量
    struct DeviceLane *next_lanes[kMaxNextLanes]; //kMaxRoadNextLanes  Road中前驱后继lane数量上限
    // for junction lane
    int overlap_size;
    struct DeviceOverlap  overlaps[kMaxOverlaps];
    // 链表表头
	  struct DeviceVehicle* first;
    struct DeviceVehicle* end;
    // 一些基本统计值
    int num_vehicles;
    // 增量更新的维护
    int insert_num;
    int remove_num;
    struct DeviceVehicle* insert_buffer[1024];  // 新插入的车辆，需要归并排序后接入主链
    struct DeviceVehicle* remove_buffer[1024];  // 删除的车辆，要从主链中移除。
    bool isInJunction;
} DeviceLane;
typedef struct InsertNum {
    int * insertNum;
}InsertNum;

#endif