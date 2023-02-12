#include <iostream>
#include "lane.h"

using namespace std;
const int LaneSize =180000;
const int RoadSize =18000;

void initDeviceLane(HostLane *lane, Context *cx) {
    // 创建hostVehicle mem object
    cx->createBuffer("hostLanes", sizeof(HostLane), LaneSize);
    // 填充
    cx->WriteBuffer(0, "hostLanes", lane, true);
    // 为kernel设置参数
    cx->setBufferAsKernelArg("init_lane", 0, POINTER, "hostLanes");
    cx->setBufferAsKernelArg("init_lane", 1, POINTER, "lanes");
    cx->setBufferAsKernelArg("init_lane", 2, POINTER, "firstVehicles");
    cx->setBufferAsKernelArg("init_lane", 3, POINTER, "endVehicles");
    // 执行kernel
    cx->execKernelNDRangeMode(0, "init_lane", {1, {LaneSize}, {30}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
    // 释放不必要的mem object
    cx->releaseBuffer("hostLanes");
}

void initDeviceRoad(Road *road, Context *cx) {
    // 创建hostVehicle mem object
    cx->createBuffer("hostRoads", sizeof(road), RoadSize);
    // 填充
    cx->WriteBuffer(0, "hostRoads", road, true);
    // 为kernel设置参数
    cx->setBufferAsKernelArg("init_road", 0, POINTER, "hostRoads");
    cx->setBufferAsKernelArg("init_road", 1, POINTER, "roads");
    // 执行kernel
    cx->execKernelNDRangeMode(0, "init_road", {1, {RoadSize}, {30}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
    // 释放不必要的mem object
    cx->releaseBuffer("hostRoads");
}