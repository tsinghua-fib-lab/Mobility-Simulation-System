#include "./programs/header_vehicle.h"

__kernel 
void init_vehicle(__global struct HostVehicle *vehicle1, __global struct DeviceVehicle *vehicle2,  
                  __global struct DeviceLane *lane)
{
    size_t tid = get_global_id(0);
    vehicle2[tid].valid = vehicle1[tid].valid;
    vehicle2[tid].canRoute = true;
    vehicle2[tid].AoiId = vehicle1[tid].AoiId;

    //初始化vehicle.attribute
    vehicle2[tid].attribute.length = vehicle1[tid].attribute.length;
    vehicle2[tid].attribute.width = vehicle1[tid].attribute.width;  
    vehicle2[tid].attribute.max_speed = vehicle1[tid].attribute.max_speed;
    vehicle2[tid].attribute.max_acc = vehicle1[tid].attribute.max_acc; 
    vehicle2[tid].attribute.usual_acc = vehicle1[tid].attribute.usual_acc;
    vehicle2[tid].attribute.max_braking_acc = vehicle1[tid].attribute.max_braking_acc;
    vehicle2[tid].attribute.usual_braking_acc = vehicle1[tid].attribute.usual_braking_acc; 
    vehicle2[tid].attribute.min_gap = vehicle1[tid].attribute.min_gap;
    vehicle2[tid].attribute.lane_change_length = vehicle1[tid].attribute.lane_change_length;  

    //初始化vehicle.snapshot
    // vehicle2[tid].snapshot.lane = &lane[vehicle1[tid].snapshot.lane_index];
    // vehicle2[tid].snapshot.speed = vehicle1[tid].snapshot.speed;
    // vehicle2[tid].snapshot.s = vehicle1[tid].snapshot.s;
    // vehicle2[tid].snapshot.lane_change_total_length = vehicle1[tid].snapshot.lane_change_total_length;
    // vehicle2[tid].snapshot.lane_change_completed_length = vehicle1[tid].snapshot.lane_change_completed_length;
    // vehicle2[tid].snapshot.has_shadow = vehicle1[tid].snapshot.has_shadow;
    // vehicle2[tid].snapshot.to_left = vehicle1[tid].snapshot.to_left;
    // vehicle2[tid].snapshot.DistanceToEnd = vehicle1[tid].snapshot.DistanceToEnd;

    // 初始化trip信息
    if(vehicle1[tid].ScheduleSize!=0){
        vehicle2[tid].ScheduleSize = vehicle1[tid].ScheduleSize;
        vehicle2[tid].LoopCount = 1;
        vehicle2[tid].scheduleIndex = 0;
        vehicle2[tid].tripIndex = 0;
        vehicle2[tid].lastTripEndTime = 0;
        for(int j = 0; j < vehicle1[tid].ScheduleSize; j++)
        {
            vehicle2[tid].schedule[j].DepartureTime = vehicle1[tid].schedule[j].DepartureTime;
            vehicle2[tid].schedule[j].WaitTime = vehicle1[tid].schedule[j].WaitTime;
            vehicle2[tid].schedule[j].LoopCount = vehicle1[tid].schedule[j].LoopCount;
            if(vehicle1[tid].schedule[j].TripSize!=0){
                vehicle2[tid].schedule[j].TripSize = vehicle1[tid].schedule[j].TripSize;
                for(int k = 0;k < vehicle1[tid].schedule[j].TripSize;k++)
                {
                    vehicle2[tid].schedule[j].Trips[k].AoiId = vehicle1[tid].schedule[j].Trips[k].AoiId;
                    vehicle2[tid].schedule[j].Trips[k].departure_time = vehicle1[tid].schedule[j].Trips[k].AoiId;
                    vehicle2[tid].schedule[j].Trips[k].wait_time = vehicle1[tid].schedule[j].Trips[k].AoiId;
                }
            }
        }
    }
}