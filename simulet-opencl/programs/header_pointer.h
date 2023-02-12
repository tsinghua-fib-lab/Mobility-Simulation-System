#ifndef DEVICEPOINTER_H
#define DEVICEPOINTER_H

#include "./programs/header_lane.h"
#include "./programs/header_vehicle.h"
#include "./programs/header_aoi.h"

/**初始化数据结构 */
typedef struct Points{
    int globalTime;
    int vehicleIndex[820000];
    int vehicleInAoi;
    int aoi2Lane[1024];
    int lane2Aoi[1024];
    DeviceVehicle *vehicles;
    DeviceAoi *aois;
    DeviceLane *lanes;
    Segment *heap;
} Points;

#endif