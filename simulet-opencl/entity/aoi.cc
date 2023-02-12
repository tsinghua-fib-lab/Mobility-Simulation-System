#include "aoi.h"
//#include<iostream>
//#include "vehicle.h"

const int AoiSize =280000;
//const int VehicleSize =820000;
// cl_mem initDeviceAoi(HostAoi *aoi, AoiVehicle *aoiVehicle,Context *cx) {
//     cx->createBuffer("hostAoi", sizeof(HostAoi), AoiSize);
//     cx->WriteBuffer(0, "hostAoi", aoi, true);
//     cx->setBufferAsKernelArg("init_aoi", 0, POINTER, "hostAoi");
//     cx->setBufferAsKernelArg("init_aoi", 1, POINTER, "aois");
//     //cx->setBufferAsKernelArg("init_aoi", 2, POINTER, "lanes");
//     cx->execKernelNDRangeMode(0, "init_aoi", {1, {AoiSize}, {28}, {0}});
//     clFinish(cx->getCommandQueueByDeviceId(0));
//     cx->releaseBuffer("hostAoi");
// }
cl_mem initDeviceAoi(DeviceAoi *aoi, AoiVehicle *aoiVehicle,Context *cx) {
    cx->createBuffer("hostAoi", sizeof(DeviceAoi), AoiSize);
    cx->WriteBuffer(0, "hostAoi", aoi, true);
    cx->createBuffer("aoiVehicle", sizeof(AoiVehicle), AoiSize);
    cx->WriteBuffer(0, "aoiVehicle", aoiVehicle, true);
    cx->setBufferAsKernelArg("init_aoi", 0, POINTER, "hostAoi");
    cx->setBufferAsKernelArg("init_aoi", 1, POINTER, "aois");
    cx->setBufferAsKernelArg("init_aoi", 2, POINTER, "aoiVehicle");
    cx->setBufferAsKernelArg("init_aoi", 3, POINTER, "vehicles");
    cx->setBufferAsKernelArg("init_aoi", 4, POINTER, "vehicle_start");
    cx->setBufferAsKernelArg("init_aoi", 5, POINTER, "vehicle_end");
    //cx->setBufferAsKernelArg("init_aoi", 2, POINTER, "lanes");
    cx->execKernelNDRangeMode(0, "init_aoi", {1, {AoiSize}, {20}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
    cx->releaseBuffer("hostAoi");
    cx->releaseBuffer("aoiVehicle");
}
cl_mem InsertVehicleInAoi(Context *cx) {
    cx->setBufferAsKernelArg("insert_vehicle", 0, POINTER, "aois");
    cx->setBufferAsKernelArg("insert_vehicle", 1, POINTER, "vehicles");
    cx->execKernelNDRangeMode(0, "insert_vehicle", {1, {AoiSize}, {20}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
}

cl_mem RemoveVehicleInAoi(Context *cx) {
    cx->setBufferAsKernelArg("remove_buff", 0, POINTER, "aois");
    cx->setBufferAsKernelArg("remove_buff", 1, POINTER, "vehicles");
    cx->setBufferAsKernelArg("remove_buff", 2, POINTER, "AoiRemoveNum");
    cx->execKernelNDRangeMode(0, "remove_buff", {1, {AoiSize}, {20}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
}
//AOI内车辆请求导航上路
cl_mem VehicleGoLane(Context *cx,unsigned long size) {
    //std::cout<<"size:"<<size<<std::endl;
    cx->setBufferAsKernelArg("vehicle_aoi_update", 0, POINTER, "points");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 1, POINTER, "aoi2LaneAto");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 2, POINTER, "routeMetaInfoAto");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 3, POINTER, "routeMetaInfo");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 4, POINTER, "aois");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 5, POINTER, "InsertNum");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 6, POINTER, "AoiRemoveNum");
    cx->setBufferAsKernelArg("vehicle_aoi_update", 7, POINTER, "num");
    cx->execKernelNDRangeMode(0, "vehicle_aoi_update", {1, {size}, {1}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
}

