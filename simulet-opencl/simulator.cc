#include <google/protobuf/util/json_util.h>
#include <fstream>
#include "wolong/map/v2//map.pb.h"
#include "proto_use.h"
#include "global.h"
#include "simulet.h"
#include "rpc/routing_client.h"
#include <sys/time.h>
#include <unistd.h>


// 用于从protobuf中载入数据
simulet::PLane *lane =(simulet::PLane *)malloc(180000*sizeof(simulet::PLane));
simulet::Agent *agent = (simulet::Agent *)malloc(820000*sizeof(simulet::Agent));
HostVehicle *hostVehicle = (HostVehicle *)malloc(820000*sizeof(HostVehicle));
drivingAction *drivingActions = (drivingAction *)malloc(820000*sizeof(drivingAction));
localController *localControllers = (localController *)malloc(820000*sizeof(localController));
HostLane *hostLane = (HostLane *)malloc(180000*sizeof(HostLane));
Road *road = (Road *)malloc(18000*sizeof(Road));
// HostAoi *hostAoi = (HostAoi *)malloc(280000*sizeof(HostAoi));
HostAoi *hostAoi = new HostAoi[280000];
DeviceAoi *deviceAoi = (DeviceAoi *)malloc(280000*sizeof(DeviceAoi));
AoiVehicle * aoiVehicle = (AoiVehicle *)malloc(280000*sizeof(AoiVehicle));
// int *InsertNum = (int *)malloc(180000*sizeof(int));
// HostLane hostLane[280000];
std::unordered_map<int, int> laneId2Index;
std::unordered_map<int, int> vehicleId2Index;
std::unordered_map<int, int> AoiID2Index;

RouteInfo routeInfos[1024];  // 表示一次可以接受N个车辆的导航请求
HostSegment hostSegmentsToDevice[131072];  // 平均每个车辆可以分配128段segments
int segmentIndex = 0; // hostSegmentsToDevice 偏移量
int infoIndex = 0; // 指向可用的RouteInfo的位置
int successRouteNum = 0; // 请求成功的导航的个数

//DeviceAoi *AOI1 = (DeviceAoi *)malloc(1*sizeof(DeviceAoi));
int main(int argc,char**argv){
    /* [x] 1. Host端载入数据*/
    InputBin();
    std::cout << "Finished Proto Input." << std::endl;
    /* [x] 2. 构建必须的host端数据结构*/
    InitHostLanes(lane, hostLane, &laneId2Index);
    InitHostVehicles(agent, hostVehicle,aoiVehicle,&vehicleId2Index);
    initHostAoi(hostAoi,deviceAoi, &AoiID2Index);
    initHostRoad(road);
    free(lane);
    free(agent);

    std::cout << "Finished initialize host data structure." << std::endl;
    
    /* [x] 3. 初始化OpenCL环境*/
    std::string programDir = "programs";
    /** 3.1. 构建Context对象
     * 此时内部包含：
     *          1> platform
     *          2> devices
     *          3> context
     **/
    Context *cx = new Context(programDir);

    // 3.2. 注册命令队列
    cx->registerCommandQueue(0);
    
    // 3.3 初始化route_heap以及device_bitmap
    cx->initRouteHeap();

    // 3.4. 构建Memory Object
    cx->createBuffer("points", sizeof(Points), 1);
    cx->createBuffer("vehicles", sizeof(DeviceVehicle), 820000);
    cx->createBuffer("num", sizeof(int), 820000);
    cx->createBuffer("drivingActions", sizeof(drivingAction), 820000);
    cx->createBuffer("num", sizeof(int), 820000);
    cx->createBuffer("localControllers", sizeof(localController), 820000);
    cx->createBuffer("lanes", sizeof(DeviceLane), 180000);
    cx->createBuffer("InsertNum", sizeof(int), 180000);
    cx->createBuffer("roads", sizeof(Road), 18000);
    cx->createBuffer("aois", sizeof(DeviceAoi), 280000);
    cx->createBuffer("AoiRemoveNum", sizeof(int), 280000);
    
    cx->createBuffer("routeMetaInfo", sizeof(RouteMetaInfo), 10000);
    cx->createBuffer("vehicle_start", sizeof(DeviceVehicle), 280000);//供AOI中链表头指针使用
    cx->createBuffer("vehicle_end", sizeof(DeviceVehicle), 280000);//供AOI中链表尾指针使用
    cx->createBuffer("firstVehicles", sizeof(DeviceVehicle), 180000);//供lane中链表尾指针使用
    cx->createBuffer("endVehicles", sizeof(DeviceVehicle), 180000);//供lane中链表尾指针使用
    clFinish(cx->getCommandQueueByDeviceId(0));

    // 全局时间
    cx->createBuffer("globalTime", sizeof(int), 1);

    // 原子结构相关
    cx->createBuffer("removeAto", sizeof(int), 1);
    cx->createBuffer("insertAto", sizeof(int), 1);
    cx->createBuffer("aoiAto", sizeof(int), 1);
    cx->createBuffer("routeMetaInfoAto", sizeof(int), 1);
    cx->createBuffer("aoi2LaneAto", sizeof(int), 1);
    cx->createBuffer("lane2AoiAto", sizeof(int), 1);
    cx->createBuffer("vehicle2LaneAto", sizeof(int), 1);
    cx->createBuffer("vehicle2LaneRemoveAto", sizeof(int), 1);
    int atos[2];
    cx->readBuffer(0, "aoi2LaneAto", atos, true);
    cx->readBuffer(0, "lane2AoiAto", atos+1, true);
    // std::cout << "aoi2laneAto: " << atos[0] << std::endl;
    // std::cout << "lane2aoiAto: " << atos[1] << std::endl;

    // 3.5. 初始化全局指针
    cx->setBufferAsKernelArg("initPoints", 0, POINTER, "points");
    cx->setBufferAsKernelArg("initPoints", 1, POINTER, "vehicles");
    cx->setBufferAsKernelArg("initPoints", 2, POINTER, "aois");
    cx->setBufferAsKernelArg("initPoints", 3, POINTER, "lanes");
    cx->setBufferAsKernelArg("initPoints", 4, POINTER, "routeHeap");
    cx->execKernelNDRangeMode(0, "initPoints");
    clFinish(cx->getCommandQueueByDeviceId(0));
    std::cout << "Finished constructing OpenCL environment." << std::endl;
    //test
    cx->setBufferAsKernelArg("test", 0, POINTER, "num");
    
    // [x] 4. routing服务相关
    routing_client = std::make_shared<RoutingClient>(::grpc::CreateChannel(
        "localhost:52101", ::grpc::InsecureChannelCredentials()));

    /* [x] 5. device端数据初始化*/
    
    initDeviceLane(hostLane, cx);
    initDeviceRoad(road,cx);
    initDeviceVehicle(hostVehicle, cx);
    initDeviceAoi(deviceAoi, aoiVehicle, cx);
    free(deviceAoi);//释放存储设备端数据的空间
   
    // InsertVehicleInAoi(cx);
    // RemoveVehicleInAoi(cx);
    std::cout << "Finished initialize device data structure." << std::endl;

    /* [x] 6. 交通模拟阶段*/
    // 为各kernel设置参数
    cx->setBufferAsKernelArg("RouteInit", 0, POINTER, "points");
    cx->setBufferAsKernelArg("RouteInit", 1, POINTER, "bitIndex");
    cx->setBufferAsKernelArg("RouteInit", 2, POINTER, "bitmap");

    cx->setBufferAsKernelArg("setGlobalTime", 0, POINTER, "points");
    cx->setBufferAsKernelArg("setGlobalTime", 1, POINTER, "globalTime");

    cx->setBufferAsKernelArg("globalTimeInc", 0, POINTER, "points");

    cx->setBufferAsKernelArg("updateVehicleIndex", 0, POINTER, "points");
    cx->setBufferAsKernelArg("updateVehicleIndex", 1, POINTER, "aoi2LaneAto");
    cx->setBufferAsKernelArg("updateVehicleIndex", 2, POINTER, "lane2AoiAto");

    //update参数
    cx->WriteBuffer(0, "drivingActions", drivingActions, true);
    cx->WriteBuffer(0, "localControllers", localControllers, true);
    cx->setBufferAsKernelArg("vehicle_lane_update", 0, POINTER, "points");
    cx->setBufferAsKernelArg("vehicle_lane_update", 1, POINTER, "lane2AoiAto");
    cx->setBufferAsKernelArg("vehicle_lane_update", 2, POINTER, "drivingActions");
    cx->setBufferAsKernelArg("vehicle_lane_update", 3, POINTER, "localControllers");
    cx->setBufferAsKernelArg("vehicle_lane_update", 4, POINTER, "vehicle2LaneAto");
    cx->setBufferAsKernelArg("vehicle_lane_update", 5, POINTER, "vehicle2LaneRemoveAto");

    //模拟相关变量与参数设置
    int totalStep = 30000;
    int timeStep  = 25200;
    //  int totalStep = 20;
    // int timeStep  = 0;
    int globalTime[1];  // 用于设置全局时间
    int routeNum[1];  // 用于从device端读取相应内容并判断当前有多少辆车发出了导航请求
    Points *point =new Points[1];
    RouteMetaInfo routeMetaInfo[10000];  // 用于存储由device端返回的导航元信息
    // [ ]: 测试使用
    // uint bitmap[ROUTE_SEGMENT_NUM/32];  // 用于读取device端的bitmap

    // 设置全局时间
    globalTime[0] = timeStep;
    cx->WriteBuffer(0, "globalTime", globalTime, true);
    cx->execKernelNDRangeMode(0, "setGlobalTime");
    clFinish(cx->getCommandQueueByDeviceId(0));

    std::cout << "Starting simulator---------------------------------------" << std::endl;
    int num =0;
    int timeup1=0;
    int timeup2=0;
    while (timeStep < totalStep)
    { 
        std::cout << "timeStep: " << timeStep << std::endl;
        
        /*5.1 Prepare阶段——负责数据拷贝与迁移*/
        
        /*5.2 主计算流*/
        cx->readBuffer(0, "points", point, true);
        //std::cout << "vehicleInAoi: " << (*point).vehicleInAoi << std::endl;
        VehicleGoLane(cx, (*point).vehicleInAoi);
        RemoveVehicleInAoi(cx);
        unsigned int laneVehicleSize = 820000 - (*point).vehicleInAoi;
        VehicleList(cx);
        cx->execKernelNDRangeMode(0, "test", {1, {820000}, {25}, {0}});
        cx->execKernelNDRangeMode(0, "vehicle_lane_update", {1, {laneVehicleSize}, {1}, {0}});
        clFinish(cx->getCommandQueueByDeviceId(0));
        
       

        // 判断是否有导航请求
        cx->readBuffer(0, "routeMetaInfoAto", routeNum, true);
        if (*routeNum > 0)
        {   // 如果有导航请求
            std::cout << "Route Num: " << *routeNum << std::endl;
            /*[5.3] 请求远程服务并等待返回结果*/
            cx->readBuffer(0, "routeMetaInfo", routeMetaInfo, true);  // 将导航元信息读取出来
            for (size_t i = 0; i < *routeNum; ++i)
            {  
                // TODO: routing服务相关
                //处理一个导航请求
                simulet::PbGetRouteRequest req;
                req.set_type(simulet::PbRouteType::ROUTE_TYPE_DRIVING);
                //导航起点、终点
                req.mutable_start()->mutable_aoi_position()->set_aoi_id(routeMetaInfo[i].startAoiId);
                req.mutable_end()->mutable_aoi_position()->set_aoi_id(routeMetaInfo[i].endAoiId); 
                routing_client->GetRoute(req, routeMetaInfo[i].vehicleIndex, routeMetaInfo[i].startAoiId, *req.mutable_end());
            }
            // 每一个step 等待路径规划处理完成
            routing_client->Wait();
            
            num =successRouteNum+num;
            std::cout<< "导航成功数(successRouteNum):" << successRouteNum << std::endl;
            //std::cout<< "总数" << num << std::endl;
            /*[5.4] route初始化*/
            // 为导航初始化kernel设置额外参数
            // std::cout << "segmentIndex: " << segmentIndex << std::endl;
            // std::cout << "infoIndex: " << infoIndex << std::endl;
            // std::cout << "Total segments: " << segmentIndex + routeInfos[infoIndex].routeLength + 1 << std::endl;
            // for(int i=0;i<successRouteNum;i++)
            //     std::cout<<routeInfos[i].s<<std::endl;
            cx->createBuffer("routeSegments", sizeof(HostSegment), segmentIndex);
            cx->createBuffer("routeInfos", sizeof(RouteInfo), successRouteNum);
            cx->WriteBuffer(0, "routeSegments", hostSegmentsToDevice, true);
            cx->WriteBuffer(0, "routeInfos", routeInfos, true);

            cx->setBufferAsKernelArg("RouteInit", 3, POINTER, "routeSegments");
            cx->setBufferAsKernelArg("RouteInit", 4, POINTER, "routeInfos");
            cx->execKernelNDRangeMode(0, "RouteInit", {1, {(ulong)successRouteNum}, {(ulong)successRouteNum}, {0}});
            clFinish(cx->getCommandQueueByDeviceId(0));

            // // 释放数据对象
            cx->releaseBuffer("routeSegments");
            cx->releaseBuffer("routeInfos");

            // 重置 hostSegmentsToDevice 和 routeInfos 、successRouteNum
            segmentIndex = 0;
            infoIndex = 0;
            successRouteNum = 0;
        }

        routeNum[0] = 0;
        cx->WriteBuffer(0, "routeMetaInfoAto", routeNum, true);

        // [ ]: 测试使用
        // cx->readBuffer(0, "bitmap", bitmap, true);
        // std::cout << "bitmap: " << std::endl;
        // for (size_t i = 0; i < ROUTE_SEGMENT_NUM/32; i++)
        // {
        //     std::cout << std::bitset<sizeof(uint)*8>(bitmap[i]) << std::endl;;
        // }
        // std::cout << std::endl << "**********************************" << std::endl;
         // 更新vehicleIndex数组
        cx->execKernelNDRangeMode(0, "updateVehicleIndex");
        clFinish(cx->getCommandQueueByDeviceId(0));
        // 更新全局时间
        cx->execKernelNDRangeMode(0, "globalTimeInc");
        clFinish(cx->getCommandQueueByDeviceId(0));

        struct timeval time;
 
    // /* 获取时间，理论到us */
    // gettimeofday(&time, NULL);
    // std::cout<<"ms:"<<time.tv_usec<<std::endl<<std::endl;
        ++timeStep;


    }

    /*结束*/
    routing_client->Shutdown();
    return 0;                   
}