#ifndef DEVICEVEHICLE_H
#define DEVICEVEHICLE_H

#include "./programs/header_aoi.h"
#include "./programs/header_lane.h"

/**一些静态数据 */
#define MAX_TRIP_NUM 20
#define MAX_SCHEDULE_NUM 20
#define TOTAL_VEHICLE 820000
const int MAX_BIT_INDEX = 262144 * 128;
const float stepInterval = 0.1;

typedef struct Trip{
    //Trip终点AoiID 
    int AoiId;
    // 期望的出发时间（单位: 秒）
    float departure_time;
    // 期望的等待时间（单位：秒）
    float wait_time;
}Trip;

typedef struct Schedule{
	Trip Trips[MAX_TRIP_NUM]; 
    //Trip数量
    int TripSize;
	// trips的执行次数，0表示无限循环，大于0表示执行几次
	int LoopCount;
	// 期望的出发时间（单位: 秒）
	float DepartureTime;
	// 期望的等待时间（单位：秒）
	float WaitTime;
}Schedule;

/**导航相关数据结构 */
typedef struct Segment{  // device segment
    int bitIndex;
    DeviceLane *lane;
    int nextLaneType;
    float distanceToEnd;
    struct Segment *nextSegment;
} Segment;

typedef struct hostSegment{  // host segment
    int laneIndex;
    int nextLaneType;
    float distanceToEnd;
} hostSegment;

typedef struct RouteMetaInfo{  // 返回给host端的导航元信息
    int vehicleIndex;
    int startAoiId;
    int endAoiId;
} RouteMetaInfo;

typedef struct RouteInfo{  // 从host端发送给device的导航初始化信息
    int vehicleIndex;
    int startIndex;
    int routeLength;
    float s;//第一个segment车道连接AOI坐标
} RouteInfo;

/**Vehicle相关数据结构 */
typedef struct VehicleRoute {
    int size;  // 总长度——指路段数
    DeviceAoi *startAoi;  // 起始Aoi指针
    DeviceAoi *endAoi;  // 终点Aoi指针
    Segment *segment;  // Segment链表
} VehicleRoute;

typedef struct VehicleAttribute {
  float length;
  float width;
  float max_speed;
  float max_acc;
  float usual_acc;
  float max_braking_acc;
  float usual_braking_acc;
  float min_gap;
  float lane_change_length;
} VehicleAttribute;

typedef struct HostVehicleMotion {
    int lane_index;//车辆所在车道索引
    float speed;//速度
    float s;//坐标
    float lane_change_total_length;//变道所需总长度
    float lane_change_completed_length;//完成变道所需长度
    bool  has_shadow;
    bool  to_left;
    float DistanceToEnd;
} HostVehicleMotion;

typedef struct HostVehicle {
    bool valid;
    int index;
    //起始AoiId
    int AoiId;
    // 常量
    struct VehicleAttribute attribute;
    // snapshot
    struct HostVehicleMotion snapshot;
	// runtime
    struct HostVehicleMotion runtime;
    // 为输出预留
    float x, y, direction;
    //Trip信息
    Schedule schedule[MAX_SCHEDULE_NUM];
    //Schedule数量
    int ScheduleSize;

    // 导航（组织为链表结构，走完一段删一段）
    // 需要能够自描述，使得车辆的计算过程不需要依赖路网拓扑结构
    // 记录终点的s坐标，记录要走的是哪个lane，车辆直接“飞”过去
    VehicleRoute vehicleRoute;;		

    // 关系链表
    int relation[2][3]; // 前后，左中右
    int shadow_relation[2][3]; // 实际开发时可以略去一些
} HostVehicle;

typedef struct DeviceVehicleMotion {
    struct DeviceLane *lane;//车辆所在道路索引
    float speed;//速度
    float s;//坐标
    float lane_change_total_length;//变道所需总长度
    float lane_change_completed_length;//完成变道所需长度
    bool has_shadow;
    bool to_left;
    float DistanceToEnd;
    int LaneChangeStatus;//车变道状态
} DeviceVehicleMotion;

typedef struct DeviceVehicle { 
    bool valid;
    bool canRoute;
    //起始AoiId
    int AoiId;
    // 常量
    struct VehicleAttribute attribute;
    // snapshot
    struct DeviceVehicleMotion snapshot;
	// runtime
    struct DeviceVehicleMotion runtime;
    // 为输出预留
    float x, y, direction;

    // //Trip信息
    Schedule schedule[MAX_SCHEDULE_NUM];
    //Schedule数量
    int ScheduleSize;
    // schedule执行次数，0表示无限循环，大于0表示执行几次
	int LoopCount;
    // 当前schedule下标
    int scheduleIndex;
    // 当前trip下标
	int tripIndex; 
    // 上次trip结束时间  
	float lastTripEndTime;

    // 导航（组织为链表结构，走完一段删一段）
    // 需要能够自描述，使得车辆的计算过程不需要依赖路网拓扑结构
    // 记录终点的s坐标，记录要走的是哪个lane，车辆直接“飞”过去
    VehicleRoute vehicleRoute;	

    // 关系链表 index
    struct DeviceVehicle *relation[2][3]; // 前后，左中右
    struct DeviceVehicle *shadow_relation[2][3]; // 实际开发时可以略去一些

    //供AOI建立AOI内车辆链表使用
    struct DeviceVehicle *pre; 
    struct DeviceVehicle *next; 
} DeviceVehicle;


/**
 * @brief 车辆感知模块
 * @details 主要依照车辆路径规划获取前方车道中会影响驾驶行为的环境信息，
 * 包括其他智能体、车道通行状况等
 */
/// 记录感知结果
typedef struct localController
{
	struct DeviceVehicle *vehicle;		   // 模块所在车辆
	float leftMotivation, rightMotivation; // 变道意愿

} localController;

typedef struct laneAhead
{
	struct DeviceLane *ptr; // 下一车道指针
	float relativeDistance ; // 相对距离
} laneAhead;
typedef struct Agent{
    struct DeviceVehicle *ptr; // 前车/后车指针（可空）
    float relativeDistance; // 车头到车尾的相对距离，恒为正
}Agent;
typedef struct laneChange
{
	bool enable;				 // 本结构体是否启用，判定条件：路径规划是否要求变道
	struct laneAhead targetLane; // 变道目标
	struct laneAhead laneAhead;	 // 变道目标的下一车道
	struct Agent agentAhead;	 // 变道目标对应位置的前车
	struct Agent agentBehind;	 // 变道目标对应位置的后车
} laneChange;

typedef struct drivingAction
{
	float acceleration;
    float laneChangeLengh;
} drivingAction;


#endif