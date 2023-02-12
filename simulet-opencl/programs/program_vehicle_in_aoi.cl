#include "./programs/header_global.h"
/**
    该kernel用于实现vehicle在aoi中的更新流
 */
__kernel
void vehicle_aoi_update(__global Points *points, __global int *aoi2LaneAto, __global int *routeMetaInfoAto,
                        __global RouteMetaInfo *routeMetaInfo, __global DeviceAoi *deviceAoi,__global int *InsertNum,__global int *AoiRemoveNum,__global int *num) {
    size_t tid = get_global_id(0);
    int rtid = points->vehicleIndex[tid];
    //请求导航
    if(points->vehicles[rtid].valid != 0 && points->vehicles[rtid].canRoute && 
        points->globalTime >= points->vehicles[rtid].schedule[points->vehicles[rtid].scheduleIndex].DepartureTime) 
    {
        //printf("tid:%d\n",tid);
        int ret = atomic_inc(routeMetaInfoAto);  // 在该ret位置插入导航元信息
        routeMetaInfo[ret].vehicleIndex = rtid;
        routeMetaInfo[ret].startAoiId = points->vehicles[rtid].AoiId;
        routeMetaInfo[ret].endAoiId = points->vehicles[rtid].schedule[points->vehicles[rtid].scheduleIndex].Trips[points->vehicles[rtid].tripIndex].AoiId;
        points->vehicles[rtid].canRoute = false;
    }
    //车插入lane中
        // if(points->vehicles[rtid].vehicleRoute.segment != 0){
        // int ret3 = atomic_inc(num+rtid);
        // if(ret3==0){
        // int laneIndex =points->vehicles[rtid].vehicleRoute.segment->lane->index;
        // int ret1 =atomic_inc(InsertNum+laneIndex);
        // points->vehicles[rtid].vehicleRoute.segment->lane->insert_buffer[ret1] = &(points->vehicles[rtid]);
        
        // //printf("insert_num:%d\n",points->vehicles[rtid].vehicleRoute.segment->lane->insert_num);
        // //将车插入AOI的removebuff中
        // int index = points->vehicles[rtid].AoiId - 500000000;
        // int ret2 = atomic_inc(AoiRemoveNum+index);
        // deviceAoi[index].removeBuff[ret2] = &(points->vehicles[rtid]);
        // }
        // //将AOI中的车辆加入到aoi2Lane中
        // int ret = atomic_inc(aoi2LaneAto);
        // points->aoi2Lane[ret] = tid;
        // }
    //车插入lane中
    if(points->vehicles[rtid].vehicleRoute.segment != 0){
        int laneIndex =points->vehicles[rtid].vehicleRoute.segment->lane->index;
        int ret1 =atomic_inc(InsertNum+laneIndex);
        points->vehicles[rtid].vehicleRoute.segment->lane->insert_buffer[ret1] = &(points->vehicles[rtid]);
        int ret3 = atomic_inc(num+rtid);
        //printf("insert_num:%d\n",points->vehicles[rtid].vehicleRoute.segment->lane->insert_num);
        //将车插入AOI的removebuff中
        int index = points->vehicles[rtid].AoiId - 500000000;
        int ret2 = atomic_inc(AoiRemoveNum+index);
        deviceAoi[index].removeBuff[ret2] = &(points->vehicles[rtid]);

        //将AOI中的车辆加入到aoi2Lane中
        int ret = atomic_inc(aoi2LaneAto);
        points->aoi2Lane[ret] = tid;
    }
}