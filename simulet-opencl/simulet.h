#include "proto_use.h"
#include "entity/vehicle.h"
#include "entity/aoi.h"
#include "entity/pointer.h"
#include "global.h"
#include "utils/envRela.h"
#include "map"

void InputBin();
void InitHostLanes(simulet::PLane *lane, HostLane *hostLane, std::unordered_map<int, int> *laneId2Index);
void InitHostVehicles(simulet::Agent *agent, HostVehicle *vehicle, AoiVehicle *aoi_vehicle , std::unordered_map<int, int> *vehicleId2Index);
// void InitHostVehicles(simulet::Agent *agent, HostVehicle *vehicle, std::unordered_map<int, int> *vehicleId2Index);
// void initHostAoi(HostAoi *host_aoi);
void initHostAoi(HostAoi *host_aoi, DeviceAoi *device_aoi,std::unordered_map<int, int> *AoiID2Index);
void initHostRoad(Road *hostRoad);
extern std::unordered_map<int, int> laneId2Index;
extern std::unordered_map<int, int> vehicleId2Index;
extern std::unordered_map<int, int> AoiID2Index;
extern HostVehicle *hostVehicle;
extern HostLane *hostLane;
extern HostAoi *hostAoi;
extern RouteInfo routeInfos[1024];  // 表示一次可以接受N个车辆的导航请求
extern HostSegment hostSegmentsToDevice[131072];  // 平均每个车辆可以分配128段segments
extern int segmentIndex; //hostSegmentsToDevice 偏移量
extern int infoIndex; //指向可用的RouteInfo的位置
extern int successRouteNum;//请求成功的导航的个数
