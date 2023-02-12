#include "./programs/header_vehicle.h"

__kernel 
void init_lane(__global struct HostLane *lane1, __global struct DeviceLane *lane2,__global struct DeviceVehicle *firstVehicles,__global struct DeviceVehicle *endVehicles)
{
    size_t tid = get_global_id(0);
    lane2[tid].first = &firstVehicles[tid];
    lane2[tid].end = &endVehicles[tid];
    lane2[tid].valid = lane1[tid].valid;
    lane2[tid].roadID = lane1[tid].roadID;
    lane2[tid].index = lane1[tid].index;
    lane2[tid].length=lane1[tid].length;
    lane2[tid].max_speed=lane1[tid].max_speed;
    if(lane1[tid].left_lane_index!=-1)
    {
        lane2[tid].left_lane = lane2+lane1[tid].left_lane_index;
    }
    else
    {
        lane2[tid].left_lane = 0;
    }
    if(lane1[tid].right_lane_index!=-1)
    {
        lane2[tid].right_lane = lane2 + lane1[tid].right_lane_index;
    }
    else
    {
        lane2[tid].right_lane=0;
    }
    
    lane2[tid].next_lane_size=lane1[tid].next_lane_size;

    if(lane2[tid].next_lane_size!=0)
    {
        for(int j=0;j<lane2[tid].next_lane_size;j++)
        {
            lane2[tid].next_lanes[j] = lane2 + lane1[tid].next_lanes[j].index;
        }    
    }
    lane2[tid].overlap_size=lane1[tid].overlap_size;
    if(lane2[tid].overlap_size!=0)
    {
        for(int j=0;j<lane2[tid].overlap_size;j++)
        {
            lane2[tid].overlaps[j].self.lane = lane2 + lane1[tid].overlaps[j].self.index;
            lane2[tid].overlaps[j].self.s = lane1[tid].overlaps[j].self.s;
            lane2[tid].overlaps[j].other_s.lane = lane2 + lane1[tid].overlaps[j].other_s.index;
            lane2[tid].overlaps[j].other_s.s = lane1[tid].overlaps[j].other_s.s;
            lane2[tid].overlaps[j].self_first = lane1[tid].overlaps[j].self_first;
        }    
    }    
}
__kernel 
void init_road(__global struct Road *road1, __global struct Road *road2){
    size_t tid = get_global_id(0);
    road2[tid].id = road1[tid].id;
    road2[tid].size = road1[tid].size;
    road2[tid].valid = road1[tid].valid;
    for(int i = 0;i<road2[tid].size;i++){
        road2[tid].lane_ids[i] = road1[tid].lane_ids[i];
    }

}