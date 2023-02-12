#include "./programs/header_global.h"

// Route初始化——这是一个NDRange并行版本
__kernel
void RouteInit(__global struct Points *points, __global uint *bitIndex, __global uint* bitmap,  
               __global struct hostSegment *routeSegments, __global struct RouteInfo *routeInfo) {
    size_t tid = get_global_id(0);  // 获取当前item的global id  
    // [ ]: 测试使用
    // printf("RouteInit vehicleIndex: %d\n", routeInfo[tid].vehicleIndex);

    /**获取导航元信息 */
    int vehicleIndex = routeInfo[tid].vehicleIndex;  // 获取当前item对应的vehicle id
    //更新车辆坐标
    points->vehicles[vehicleIndex].snapshot.s = routeInfo[tid].s;
    points->vehicles[vehicleIndex].runtime.s = routeInfo[tid].s;
    int startIndex = routeInfo[tid].startIndex;  // 获取当前vehicle的Lane信息
    int routeLength = routeInfo[tid].routeLength;  // 获取当前vehicle对应的route长度
    int totalSize = routeSegments[startIndex].laneIndex;
    int startAoiIndex = routeSegments[startIndex].nextLaneType;
    int endAoiIndex = routeSegments[startIndex].distanceToEnd;
    points->vehicles[vehicleIndex].vehicleRoute.size = totalSize;
    points->vehicles[vehicleIndex].vehicleRoute.startAoi = points->aois+startAoiIndex;
    points->vehicles[vehicleIndex].vehicleRoute.endAoi = points->aois+endAoiIndex;
    //printf("test1\n");
    /**辅助信息 */
    int laneIndex;  // laneIndex: Lane的索引
    int nextLaneType;  // nextLaneType: 下一车道的类型
    float distanceToEnd;  // distanceToEnd：当前lane的distanceToEnd
    bool flag = true;  // 循环标志
    uint oldUint;  // 用于临时存储bitmap中的内容
    uint index = 0;  // 标识当前的循环次数
    uint preIndex;  // 记录上一次循环中得到的bit index
    uint thisIndex;  // 记录当前循环中得到的bit index——根据bit index能够索引到segment
    int mask;  // 作为掩码使用
    while(index != routeLength) {
        laneIndex = routeSegments[startIndex+1+index].laneIndex;  // 获取导航信息中的下一条lane的index
        nextLaneType = routeSegments[startIndex+1+index].nextLaneType;  // 获取导航信息中下一条nextLaneType
        distanceToEnd = routeSegments[startIndex+1+index].distanceToEnd;  // 获取当前lane的distanceToEnd
        // 申请空间——每次循环申请一个segment
        while (flag) 
        {
            // 1.获取当前的bit index值
            // 这里返回当前时刻探测到的值，并对该值加1
            // 其余item不会获得相同bit index——这意味着不可能存在两个item对同一个bit进行操作
            thisIndex = atomic_inc(bitIndex);
            thisIndex = thisIndex%(MAX_BIT_INDEX);

            // [ ]: 测试使用
            // printf("vehicleId: %d, thisIndex: %d\n", vehicleIndex, thisIndex);

            // 2.验证当前bit是否可用
            oldUint = bitmap[(int)(thisIndex/32)];  // 2.1 获取对应uint
            if (((oldUint>>(31-(thisIndex%32)))&0x1) == 0)  // 2.2 判断对应bit是否为0
            {
                // 如果为0，则表示该位可用
                /**需要完成的目标：
                    1. 设置对应bit为1
                    2. 保证同步性质
                */
                // 设置循环跳出标志
                flag = false;  
                // 原子或操作——这里不用oldUint是因为可能中途发生了变化
                mask = 0x1<<(31-thisIndex%32);
                atomic_or(bitmap+(int)(thisIndex/32), mask);  
            }
            // 如果当前探测位不为0，则继续探测
        }
        // 判断当前循环是不是首次循环
        if(index == 0)
        {
            // 首次循环是对vehicle中的指针进行设置
            points->vehicles[vehicleIndex].vehicleRoute.segment = points->heap+thisIndex;
            points->heap[thisIndex].bitIndex = thisIndex;
            points->heap[thisIndex].lane = points->lanes+laneIndex;
            points->heap[thisIndex].nextLaneType = nextLaneType;
            points->heap[thisIndex].distanceToEnd = distanceToEnd;
            points->heap[thisIndex].nextSegment = 0;
            preIndex = thisIndex;
        } else {
            // 其余循环是对Segment中的指针进行设置
            // 设置当前segment的内容
            points->heap[thisIndex].bitIndex = thisIndex;
            points->heap[thisIndex].lane = points->lanes+laneIndex;
            points->heap[thisIndex].nextLaneType = nextLaneType;
            points->heap[thisIndex].distanceToEnd = distanceToEnd;
            points->heap[thisIndex].nextSegment = 0;
            // 设置前一个segment中的指针——使其指向当前segment
            points->heap[preIndex].nextSegment = points->heap+thisIndex;
            preIndex = thisIndex;
        }
        ++index;  // 循环次数加1
        flag = true;
    }
    points->heap[thisIndex].nextSegment = 0;
    //printf("test2\n");
}

// vehicleCalculate: 主计算流
__kernel
void vehicleCalculate(__global Points *points, __global RouteMetaInfo *routeMetaInfo, 
                      __global uint* bitmap, __global uint *bitIndex,
                      __global int *removeAto, __global int *insertAto, 
                      __global int *aoiAto, __global int *routeMetaInfoAto) {
    // 获取tid
    size_t tid = get_global_id(0);

    // [ ]: 测试使用
    if (tid == 0)
    {
        // printf("Device sizeof Vehicle: %d\n", sizeof(DeviceVehicle));
        // printf("Device sizeof Aoi: %d\n", sizeof(DeviceAoi));
        // printf("Device sizeof Lane: %d\n", sizeof(DeviceLane));
        // printf("Device sizeof Segment: %d\n", sizeof(Segment));
        // printf("globalTime: %d\n", points->globalTime);
    }

    if(points->vehicles[tid].valid)  // 判断是否为可用数据
    {
        if(points->vehicles[tid].runtime.s > 0) 
        {   // 说明车辆已经在路上了
            // FIXME: 这里假设每一个timeStep就到达下一个segment
            if(points->vehicles[tid].vehicleRoute.segment != 0)  // 这里可以用于表示是否请求成功
            {
                // 修正segment
                uint vehicleBitIndex = points->vehicles[tid].vehicleRoute.segment->bitIndex;  // 取得当前segment在bitmap中的index
                uint mask;

                // FIXME: 这里模拟可能发生的自主变道
                // TODO: 如何在自主变道过程中考虑distanceToEnd
                if(tid%3 == 1 && points->vehicles[tid].runtime.s==4)   // 1,4,7,10...
                {   // 这里是模拟的需要主动变道的车辆
                    // 这里假设从当前segment开始，插入两个segment并跳过原线路的一个segment
                    // 删除原线路的一个segment
                    uint nextBitIndex = (points->vehicles[tid].vehicleRoute.segment->nextSegment)->bitIndex;  // 这里获取需要删除的segment对应的bitIndex

                    mask = 0xFFFFFFFF-(0x1<<(31-(nextBitIndex%32)));
                    atomic_and(bitmap+(nextBitIndex/32), mask);

                    // 插入两个segment
                    uint preIndex = vehicleBitIndex;
                    uint thisIndex;
                    uint oldUint;
                    bool flag = true;
                    for(int i = 0; i< 2; i++) // 循环两侧——每次循环插入一个segment
                    {
                        while(flag)
                        {
                            thisIndex = atomic_inc(bitIndex);  // 在循环中不断获取当前探测位
                            thisIndex = thisIndex%(MAX_BIT_INDEX);
                            oldUint = bitmap[(int)(thisIndex/32)];
                            if (((oldUint>>(31-(thisIndex%32)))&0x1) == 0)
                            {
                                flag = false;
                                mask = 0x1<<(31-(thisIndex%32));
                                atomic_or(bitmap+(int)(thisIndex/32), mask);  
                            }
                        }
                        points->heap[thisIndex].bitIndex = thisIndex;
                        points->heap[thisIndex].lane = points->lanes;
                        points->heap[thisIndex].nextLaneType = 1;
                        points->heap[thisIndex].nextSegment = 0;
                        points->heap[preIndex].nextSegment = points->heap+thisIndex;
                        preIndex = thisIndex;
                        flag = true;
                    }
                    points->heap[thisIndex].nextSegment = points->heap+((vehicleBitIndex+2)%MAX_BIT_INDEX);  // 链接回原来的可用部分
                }

                mask = 0xFFFFFFFF-(0x1<<(31-(vehicleBitIndex%32)));
                atomic_and(bitmap+(vehicleBitIndex/32), mask);  // 对应位置0
                if(points->vehicles[tid].vehicleRoute.segment->nextSegment == 0)  // 表示到达了当前Trip的终点
                {
                    points->vehicles[tid].runtime.s = 0;
                    points->vehicles[tid].scheduleIndex += 1;
                    points->vehicles[tid].tripIndex = 0;
                }else 
                {
                    points->vehicles[tid].vehicleRoute.segment = (Segment *)(points->vehicles[tid].vehicleRoute.segment->nextSegment);
                    points->vehicles[tid].tripIndex += 1;
                    points->vehicles[tid].runtime.s += 1;
                }
            }
        } else if (points->vehicles[tid].ScheduleSize > 0 && 
                        points->globalTime >= points->vehicles[tid].schedule[points->vehicles[tid].scheduleIndex].DepartureTime)
        {   // 车辆需要从Aoi中出发——请求route
            // 获取导航元信息
            int ret = atomic_inc(routeMetaInfoAto);  // 在该ret位置插入导航元信息

            routeMetaInfo[ret].vehicleIndex = tid;
            routeMetaInfo[ret].startAoiId = points->vehicles[tid].AoiId;
            routeMetaInfo[ret].endAoiId = points->vehicles[tid].schedule[points->vehicles[tid].scheduleIndex].Trips[points->vehicles[tid].tripIndex].AoiId;

            points->vehicles[tid].runtime.s = 1;
        }
    }
}