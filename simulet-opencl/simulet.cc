#include <iostream>
#include <fstream>
#include "simulet.h"
#include "CL/opencl.h"


simulet::PMap node;
simulet::Pagents Pvehicle;
cl_context context = NULL;
cl_device_id device = NULL;

//将bin数据读入proto结构中
void InputBin(){
    
    const char * filename="map.bin";
    
    std::fstream input("./data/map.bin", std::ios::in | std::ios::binary);
    if (!input)
    {
        std::cout<<"intput file : "<< filename << " is not found."<<std::endl;
    }

    if(!node.ParseFromIstream(&input)){
        std::cerr<<"Failed to parse maps."<<std::endl;

    }
    const char * filename2="vehicle_data.bin";
    std::fstream input2("./data/vehicle_data.bin", std::ios::in | std::ios::binary);
    if (!input2)
    {
        std::cout<<"intput file : "<< filename2 << " is not found."<<std::endl;
    }
    if(!Pvehicle.ParseFromIstream(&input2)){
        std::cerr<<"Failed to parse vehicles."<<std::endl;
    }
}

//将数据存入车辆数据结构中
void InitHostVehicles(simulet::Agent *agent, HostVehicle *vehicle, AoiVehicle *aoi_vehicle , std::unordered_map<int, int> *vehicleId2Index){
    int size = Pvehicle.agents_size();
    for(int i=0;i<size;i++){
        //读入数据
        agent[i] = Pvehicle.agents(i);

        vehicle[i].valid = true;

        (*vehicleId2Index).insert(std::make_pair(agent[i].id(), i));       

        vehicle[i].index = agent[i].id();
        vehicle[i].attribute.length = agent[i].attribute().length();
        vehicle[i].attribute.width = agent[i].attribute().width();
        vehicle[i].attribute.max_speed = agent[i].attribute().max_speed();
        vehicle[i].attribute.max_acc = agent[i].attribute().max_acceleration();
        vehicle[i].attribute.usual_acc = agent[i].attribute().usual_acceleration();
        vehicle[i].attribute.max_braking_acc = agent[i].attribute().max_braking_acceleration();
        vehicle[i].attribute.usual_braking_acc = agent[i].attribute().usual_braking_acceleration();
        vehicle[i].attribute.min_gap = agent[i].vehicle_attribute().min_gap();
        vehicle[i].attribute.lane_change_length = agent[i].vehicle_attribute().lane_change_length();
        vehicle[i].snapshot.speed = agent[i].motion().speed();
        vehicle[i].AoiId = agent[i].home().aoi_position().aoi_id();
        //初始化aoi中的车辆信息
        aoi_vehicle[vehicle[i].AoiId-500000000].vehicle_index[aoi_vehicle[vehicle[i].AoiId-500000000].number]=i;
        aoi_vehicle[vehicle[i].AoiId-500000000].number++;
        if(agent[i].schedules_size()!=0){
            vehicle[i].ScheduleSize = agent[i].schedules_size();
            for(int j = 0;j < agent[i].schedules_size();j++)
            {
                vehicle[i].schedule[j].DepartureTime = agent[i].schedules(j).departure_time();
                vehicle[i].schedule[j].WaitTime = agent[i].schedules(j).wait_time();
                vehicle[i].schedule[j].LoopCount = agent[i].schedules(j).loop_count();
                if(agent[i].schedules(j).trips_size()!=0){
                    vehicle[i].schedule[j].TripSize = agent[i].schedules(j).trips_size();       
                    for(int k = 0;k < agent[i].schedules(j).trips_size();k++)
                    {
                        vehicle[i].schedule[j].Trips[k].AoiId = agent[i].schedules(j).trips(k).end().aoi_position().aoi_id();
                        vehicle[i].schedule[j].Trips[k].departure_time = agent[i].schedules(j).trips(k).departure_time();
                        vehicle[i].schedule[j].Trips[k].wait_time = agent[i].schedules(j).trips(k).wait_time();
                    }
                }
            }
        }
    }
    // for(int i=0;i<size;i++){
    //     //初始化aoi中的车辆信息
    //     aoi_vehicle[vehicle[i].AoiId-500000000].vehicle_index[aoi_vehicle[vehicle[i].AoiId-500000000].number]=i;
    //     aoi_vehicle[vehicle[i].AoiId-500000000].number++;
    //       if(aoi_vehicle[vehicle[i].AoiId-500000000].number>1600)
    //         std::cout<<"number"<<vehicle[i].AoiId-500000000<<":"<<aoi_vehicle[vehicle[i].AoiId-500000000].number<<std::endl;
    // }
    for (size_t i = size; i < 820000; i++)
    {
        vehicle[i].valid = false;
    }
}

//将数据存入lane数据结构中
void InitHostLanes(simulet::PLane *lane, HostLane *hostlane, std::unordered_map<int, int> *laneId2Index){
    int size = node.lanes_size();
        //读入地图数据并转出符合数据结构的数据
    for(int i=0; i<size; i++){
        //读入数据
        lane[i] = node.lanes(i);
        hostlane[i].valid = true;
    
        (*laneId2Index).insert(std::make_pair(lane[i].id(), i));

        //输出数据
        hostlane[i].index = lane[i].id();
        hostlane[i].roadID = lane[i].parent_id();
        hostlane[i].length = lane[i].length();
        hostlane[i].max_speed = lane[i].max_speed();
        hostlane[i].next_lane_size = lane[i].successors_size();
        hostlane[i].overlap_size = lane[i].overlaps_size();
        if (lane[i].left_lane_ids().empty()) {
            hostlane[i].left_lane_index = -1;
        }
        else {
            hostlane[i].left_lane_index = lane[i].left_lane_ids(0);
        }

        if(lane[i].right_lane_ids().empty()){
            hostlane[i].right_lane_index = -1;
        }
        else {
            hostlane[i].right_lane_index = lane[i].right_lane_ids(0);
        }

        if(!lane[i].overlaps().empty()){
            for(int j=0;j<lane[i].overlaps_size();j++)
            {
                hostlane[i].overlaps[j].self.index=lane[i].overlaps(j).self().lane_id();
                hostlane[i].overlaps[j].self.s=lane[i].overlaps(j).self().s();
                hostlane[i].overlaps[j].other_s.index=lane[i].overlaps(j).other().lane_id();
                hostlane[i].overlaps[j].other_s.s=lane[i].overlaps(j).other().s();
                hostlane[i].overlaps[j].self_first=lane[i].overlaps(j).self_first();
            }
        }
        //结构体内只有index
        if(lane[i].successors_size()){
            for(int j=0;j<lane[i].successors_size();j++){
                hostlane[i].next_lanes[j].index = lane[i].successors(j).id();
            }
        } 
    }   

    for (size_t i = size; i < 180000; i++)
    {
        hostlane[i].valid = false;
    }
}

void initHostAoi(HostAoi *host_aoi, DeviceAoi *device_aoi, std::unordered_map<int, int> *AoiID2Index)
{
    int size = node.aois_size();
    for(int i=0;i<size;i++)
    {
        host_aoi[i].AoiID = node.aois(i).id();
        (*AoiID2Index).insert(std::make_pair(host_aoi[i].AoiID, i)); 
        host_aoi[i].DrivingPositionSize = node.aois(i).driving_positions_size();
        host_aoi[i].PositionSize = node.aois(i).positions_size();
        host_aoi[i].DrivingGateSize = node.aois(i).driving_gates_size();
        host_aoi[i].area = node.aois(i).area();
        //获取设备端初始化数据
        device_aoi[i].AoiID = node.aois(i).id();
        device_aoi[i].DrivingPositionSize = node.aois(i).driving_positions_size();
        device_aoi[i].PositionSize = node.aois(i).positions_size();
        device_aoi[i].DrivingGateSize = node.aois(i).driving_gates_size();
        device_aoi[i].area = node.aois(i).area();
        for(int j = 0;j<host_aoi[i].DrivingPositionSize;j++)
        {

            host_aoi[i].DrivingPositions[j].index = node.aois(i).driving_positions(j).lane_id();
            host_aoi[i].DrivingPositions[j].s = node.aois(i).driving_positions(j).s();
            // if(host_aoi[i].DrivingPositions[j].s!=0)
                // std::cout<<"AOIID:"<<i<<"laneID"<<host_aoi[i].DrivingPositions[j].index<<"s:"<<host_aoi[i].DrivingPositions[j].s<<std::endl;
            host_aoi[i].lanS.insert(std::make_pair(host_aoi[i].DrivingPositions[j].index,host_aoi[i].DrivingPositions[j].s));
            device_aoi[i].DrivingPositions[j].index = node.aois(i).driving_positions(j).lane_id();
            device_aoi[i].DrivingPositions[j].s = node.aois(i).driving_positions(j).s();
        }
        for(int j = 0;j<host_aoi[i].PositionSize;j++)
        {
            host_aoi[i].Positions[j].x = node.aois(i).positions(j).x();
            host_aoi[i].Positions[j].y = node.aois(i).positions(j).y();
            device_aoi[i].Positions[j].x = node.aois(i).positions(j).x();
            device_aoi[i].Positions[j].y = node.aois(i).positions(j).y();
        }
        for(int j = 0;j<host_aoi[i].DrivingGateSize;j++)
        {
            host_aoi[i].DrivingGates[j].x = node.aois(i).driving_gates(j).x();
            host_aoi[i].DrivingGates[j].y = node.aois(i).driving_gates(j).y();
            device_aoi[i].DrivingGates[j].x = node.aois(i).driving_gates(j).x();
            device_aoi[i].DrivingGates[j].y = node.aois(i).driving_gates(j).y();
        }
    }
    // for(int i=0;i<100;i++){
    //     for(int j = 0;j<host_aoi[i].DrivingPositionSize;j++)
    //     {
    //         std::cout<<"index:"<<host_aoi[i].DrivingPositions[j].index<<std::endl;
    //         std::cout<<"s:"<<host_aoi[i].DrivingPositions[j].s<<std::endl;
    //         std::cout<<"map_s:"<<host_aoi[i].lanS[host_aoi[i].DrivingPositions[j].index]<<std::endl;
    //     }
    // } 

}
void initHostRoad(Road *hostRoad)
{
    int size =node.roads_size();
    for (int i = 0; i < size; i++)
    {
        hostRoad[i].valid = true;
        hostRoad[i].size = node.roads(i).lane_ids_size();
        hostRoad[i].id = node.roads(i).id();
        for (int j = 0; j < node.roads(i).lane_ids_size(); j++)
        {
            hostRoad[i].lane_ids[j] = node.roads(i).lane_ids(j);
        }
    }  
    for (size_t i = size; i < 18000; i++)
    {
        hostRoad[i].valid = false;
    }
}