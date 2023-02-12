#include "./programs/header_global.h"

// initPoints: 用于初始化指针
__kernel 
void initPoints(__global struct Points *points, __global struct DeviceVehicle *vehicles,
                __global struct DeviceAoi *aois, __global struct DeviceLane *lanes, 
                __global struct Segment *segments) {
                    points->vehicles = vehicles;
                    points->aois = aois;
                    points->lanes = lanes;
                    points->heap = segments;
                    for(int i = 0; i < TOTAL_VEHICLE; ++i) 
                    {
                        points->vehicleIndex[i] = i;
                    }
                    points->vehicleInAoi = TOTAL_VEHICLE;
                }

// setGlobalTime：用于设置全局时间
__kernel
void setGlobalTime(__global struct Points *points, __global int *globalTime) {
    points->globalTime = *globalTime;
}

// globalTimeInc：全局时间加1
__kernel
void globalTimeInc(__global struct Points *points) {
    points->globalTime += 1;
}

// 更新vehicleIndex数组
__kernel
void updateVehicleIndex(__global struct Points *points, __global int *aoi2LaneAto, __global int *lane2AoiAto) {
    int aoi2Lane = *aoi2LaneAto;
    int lane2Aoi = *lane2AoiAto;
    // TODO: 这里要加排序
    int temp;
    int temp2;
    // 对aoi2Lane而言，需要大的先换
     for(int i = 0; i < aoi2Lane; ++i)  // 执行n次——每次增加一个有序元素
    {
        for(int j = aoi2Lane-1; j > i; --j)
        {
            temp = points->aoi2Lane[j];
            temp2 = points->aoi2Lane[j-1];
            if(temp > temp2)
            {
                points->aoi2Lane[j-1] = temp;
                points->aoi2Lane[j] = temp2;
            }
        }
    }

    // 对lane2Aoi而言，需要小的先换
    for(int i = 0; i < lane2Aoi; ++i)  // 执行n次——每次增加一个有序元素
    {
        for(int j = lane2Aoi-1; j > i; --j)
        {
            temp = points->lane2Aoi[j];
            temp2 = points->lane2Aoi[j-1];
            if(temp < temp2)
            {
                points->lane2Aoi[j-1] = temp;
                points->lane2Aoi[j] = temp2;
            }
        }
    }
    // printf("aoi2LaneAto: %d\n", *aoi2LaneAto);
    // printf("lane2AoiAto: %d\n", *lane2AoiAto);
    // printf("vehicle in aoi: %d\n", points->vehicleInAoi);
    if (aoi2Lane >= lane2Aoi && aoi2Lane != 0) // 进入Lane的车辆较返回Aoi的车更多(或相同)
    {
        int tmp;
        for(int i = 0; i < lane2Aoi; ++i)
        {
            tmp = points->vehicleIndex[points->aoi2Lane[i]];
            points->vehicleIndex[points->aoi2Lane[i]] = points->vehicleIndex[points->lane2Aoi[i]];
            points->vehicleIndex[points->lane2Aoi[i]] = tmp;
        }
        for(int i = lane2Aoi; i < aoi2Lane; ++i)
        {
            tmp = points->vehicleIndex[points->aoi2Lane[i]];
            points->vehicleIndex[points->aoi2Lane[i]] = points->vehicleIndex[points->vehicleInAoi - 1];
            points->vehicleIndex[points->vehicleInAoi - 1] = tmp;
            points->vehicleInAoi = points->vehicleInAoi - 1;
        }
    } else if (lane2Aoi != 0)  // 进入Aoi的车辆多于进入Lane的车辆
    {
        int tmp;
        for(int i = 0; i < aoi2Lane; ++i)
        {
            tmp = points->vehicleIndex[points->lane2Aoi[i]];
            points->vehicleIndex[points->lane2Aoi[i]] = points->vehicleIndex[points->aoi2Lane[i]];
            points->vehicleIndex[points->aoi2Lane[i]] = tmp;
        }
        for(int i = aoi2Lane; i < lane2Aoi; ++i)
        {
            tmp = points->vehicleIndex[points->lane2Aoi[i]];
            points->vehicleIndex[points->lane2Aoi[i]] = points->vehicleIndex[points->vehicleInAoi];
            points->vehicleIndex[points->vehicleInAoi] = tmp;
            points->vehicleInAoi = points->vehicleInAoi + 1;
        }
    }
    // 重置原子变量
    *lane2AoiAto = 0;
    *aoi2LaneAto = 0;

    // printf("After aoi2LaneAto: %d\n", *aoi2LaneAto);
    // printf("After lane2AoiAto: %d\n", *lane2AoiAto);
    // printf("After vehicle in aoi: %d\n", points->vehicleInAoi);
}