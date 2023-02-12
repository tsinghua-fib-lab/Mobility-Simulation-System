// #include "./programs/header_aoi.h"
#include "./programs/header_vehicle.h"
__kernel 
void init_aoi(__global struct DeviceAoi *host_aoi, __global struct DeviceAoi *device_aoi, __global struct AoiVehicle *aoiVehicle,
__global struct DeviceVehicle *device_vehicle, __global struct DeviceVehicle *vehicle_start, __global struct DeviceVehicle *vehicle_end){
    
    size_t tid = get_global_id(0);

    device_aoi[tid].AoiID = host_aoi[tid].AoiID;
    device_aoi[tid].DrivingPositionSize = host_aoi[tid].DrivingPositionSize;
    device_aoi[tid].PositionSize =host_aoi[tid].PositionSize;
    device_aoi[tid].DrivingGateSize = host_aoi[tid].DrivingGateSize;
    device_aoi[tid].area = host_aoi[tid].area;
    for(int i=0;i<host_aoi[tid].DrivingPositionSize;i++)
    {
        device_aoi[tid].DrivingPositions[i].index = host_aoi[tid].DrivingPositions[i].index;
        device_aoi[tid].DrivingPositions[i].s = host_aoi[tid].DrivingPositions[i].s;
    }

    for(int i=0;i<host_aoi[tid].PositionSize;i++)
    {
        device_aoi[tid].Positions[i].x = host_aoi[tid].Positions[i].x;
        device_aoi[tid].Positions[i].y = host_aoi[tid].Positions[i].y;
    }
      
    for(int i=0;i<host_aoi[tid].DrivingGateSize;i++)
    {
        device_aoi[tid].DrivingGates[i].x = host_aoi[tid].DrivingGates[i].x;
        device_aoi[tid].DrivingGates[i].y = host_aoi[tid].DrivingGates[i].y;
    }

    device_aoi[tid].vehicleNum = aoiVehicle[tid].number;
    device_aoi[tid].vehicleListStart = &vehicle_start[tid];
    device_aoi[tid].vehicleListEnd = &vehicle_end[tid];
    
    for(int i=0;i<aoiVehicle[tid].number;i++){
          if(i==0){
            device_aoi[tid].vehicleListStart->next = &device_vehicle[aoiVehicle[tid].vehicle_index[i]];
            device_vehicle[aoiVehicle[tid].vehicle_index[i]].next =&device_vehicle[aoiVehicle[tid].vehicle_index[i+1]];
            device_vehicle[aoiVehicle[tid].vehicle_index[i]].pre = device_aoi[tid].vehicleListStart;
            continue;
        }
        if(i==aoiVehicle[tid].number-1){
            device_vehicle[aoiVehicle[tid].vehicle_index[i]].next = device_aoi[tid].vehicleListEnd;
            device_aoi[tid].vehicleListEnd->pre = &device_vehicle[aoiVehicle[tid].vehicle_index[i]];
            device_vehicle[aoiVehicle[tid].vehicle_index[i]].pre = &device_vehicle[aoiVehicle[tid].vehicle_index[i-1]]; 
            device_aoi[tid].vehicleListEnd->next = 0;
            continue;
        }
        device_vehicle[aoiVehicle[tid].vehicle_index[i]].next = &device_vehicle[aoiVehicle[tid].vehicle_index[i+1]];
        device_vehicle[aoiVehicle[tid].vehicle_index[i]].pre = &device_vehicle[aoiVehicle[tid].vehicle_index[i-1]];
    }
    //To test
    // if(tid==3){
    //           device_aoi[tid].insertNum=5;
    //           device_aoi[tid].removeNum=5;
    //           printf("未排序前insertBuff中的车辆出发时间:\n");
    //           for(int i= 0;i<5;i++){
    //             device_aoi[tid].insertBuff[i] = &device_vehicle[i];
    //             printf("departure_time:%f\n",device_aoi[tid].insertBuff[i]->schedule[device_aoi[tid].insertBuff[i]->scheduleIndex].DepartureTime);
    //             device_aoi[tid].removeBuff[i] = &device_vehicle[i];
    //             printf("remove_time:%f\n",device_aoi[tid].removeBuff[i]->schedule[device_aoi[tid].insertBuff[i]->scheduleIndex].DepartureTime);
    //           } 
    //     //    printf("number:%d\n",aoiVehicle[tid].number);
    //     //    for(int i= 0;i<aoiVehicle[tid].number;i++){
    //     //         printf("index%d:%d\n",i,aoiVehicle[tid].vehicle_index[i]);
    //     //         printf("scheduleInde:%d\n",device_vehicle[aoiVehicle[tid].vehicle_index[i]].scheduleIndex);
                 
    //     //    }
    //     printf("Aoi链表中车辆出发时间:\n");
    //     printf("departure_time:%f\n",device_aoi[tid].vehicleListStart->next->schedule[0].DepartureTime);
    //     printf("departure_time:%f\n",device_aoi[tid].vehicleListStart->next->next->schedule[0].DepartureTime);
    //     printf("departure_time:%f\n",device_aoi[tid].vehicleListStart->next->next->next->schedule[0].DepartureTime);    

    // }

} 

__kernel 
void insert_vehicle(__global struct DeviceAoi *device_aoi, __global struct DeviceVehicle *device_vehicle)
{
   
    size_t tid = get_global_id(0);
    if(device_aoi[tid].insertNum!=0){
        struct DeviceVehicle * list = device_aoi[tid].vehicleListEnd->pre;

        for(int i=0;i<device_aoi[tid].insertNum;i++){
            
            list->next=device_aoi[tid].insertBuff[i];
            device_aoi[tid].insertBuff[i]->pre = list;
            list = device_aoi[tid].insertBuff[i];   
    }
        list->next=device_aoi[tid].vehicleListEnd;
        device_aoi[tid].vehicleListEnd->pre = list;
        device_aoi[tid].vehicleNum = device_aoi[tid].vehicleNum +device_aoi[tid].insertNum;
        for(int i=0;i<device_aoi[tid].insertNum;i++){
            device_aoi[tid].insertBuff[i] = 0;
        }
        device_aoi[tid].insertNum = 0;
    }
    //To test
    // printf("vehicleNum:%d\n",device_aoi[tid].vehicleNum);
    // struct DeviceVehicle *test = device_aoi[tid].vehicleListStart->next;
    // printf("插入后链表中车辆出发时间:\n");
    // for(int i=0;i<device_aoi[tid].vehicleNum;i++){
    //     printf("DepartureTime:%f\n",test->schedule[test->scheduleIndex].DepartureTime);
    //     test =test->next;
    // }

}


__kernel
void remove_buff(__global struct DeviceAoi *device_aoi, __global struct DeviceVehicle *device_vehicle,__global int *AoiRemoveNum)
{
    size_t tid = get_global_id(0);
    device_aoi[tid].removeNum = AoiRemoveNum[tid];
    if(device_aoi[tid].removeNum!=0){
        //printf("removeNum:%d\n",device_aoi[tid].removeNum);
        for(int i =0;i<device_aoi[tid].removeNum;i++){
            device_aoi[tid].removeBuff[i]->pre->next=device_aoi[tid].removeBuff[i]->next;
            device_aoi[tid].removeBuff[i]->next->pre = device_aoi[tid].removeBuff[i]->pre;
            //device_aoi[tid].removeBuff[i]->pre =0;
            //device_aoi[tid].removeBuff[i]->next = 0;
            device_aoi[tid].removeBuff[i] = 0;
            device_aoi[tid].vehicleNum--;
        }
        device_aoi[tid].removeNum=0;
        AoiRemoveNum[tid]=0;
    }

    
}