#include <iostream>
#include "lane.h"
#include <fstream>
#include <cstring>
#include "vehicle.h"

const int VehicleSize =820000;
const int LaneSize =180000;
cl_mem initDeviceVehicle(HostVehicle *vehicle, Context *cx) {
    cx->createBuffer("hostVehicles", sizeof(HostVehicle), VehicleSize);
    cx->WriteBuffer(0, "hostVehicles", vehicle, true);
    cx->setBufferAsKernelArg("init_vehicle", 0, POINTER, "hostVehicles");
    cx->setBufferAsKernelArg("init_vehicle", 1, POINTER, "vehicles");
    cx->setBufferAsKernelArg("init_vehicle", 2, POINTER, "lanes");
    cx->execKernelNDRangeMode(0, "init_vehicle", {1, {VehicleSize}, {25}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
    cx->releaseBuffer("hostVehicles");
}

cl_mem VehicleList(Context *cx) {
    //std::cout<<"first"<<std::endl;
    cx->setBufferAsKernelArg("vehicle_list", 0, POINTER, "lanes");
    cx->setBufferAsKernelArg("vehicle_list", 1, POINTER, "vehicles");
    cx->setBufferAsKernelArg("vehicle_list", 2, POINTER, "InsertNum");
    //std::cout<<"end1"<<std::endl;
    cx->execKernelNDRangeMode(0, "vehicle_list", {1, {LaneSize}, {30}, {0}});
    clFinish(cx->getCommandQueueByDeviceId(0));
    //std::cout<<"end"<<std::endl;
}
// cl_mem InitDeviceVehicle(HostVehicle *vehicle, cl_context context, cl_device_id device, cl_mem *memObject)
// {   

//     cl_int ret;
//     cl_int status;
//     /** step 4: create command queue */
//     cl_command_queue commandQueue = NULL;
//     commandQueue = clCreateCommandQueueWithProperties(context, device, 0, &ret);
//     if ((CL_SUCCESS != ret) || (NULL == commandQueue))
//     {
//         std::cout << "Error creating command queue: " << ret << std::endl;
//         return 0;
//     }

//     cl_mem memObjects[2]={0,0};
//     memObjects[0] = clCreateBuffer(context, CL_MEM_READ_WRITE | CL_MEM_USE_HOST_PTR, sizeof(HostVehicle)*VehicleSize, vehicle, NULL);
//     clFinish(commandQueue);
//     memObjects[1] = clCreateBuffer(context, CL_MEM_READ_WRITE, sizeof(DeviceVehicle)*VehicleSize, NULL, NULL);
//     if ( memObjects == NULL ) 
//         perror("Error in clCreateBuffer.\n");

//     // 6. program
//     const char * filename = "program_init_vehicle.cl";
//     std::string sourceStr;
//     status = convertToString(filename, sourceStr);
//     if (status)
//         std::cout << status << "  !!!!!!!!" << std::endl;
//     const char * source = sourceStr.c_str();
//     size_t sourceSize[] = { strlen(source) };
//     //创建程序对象
//     cl_program program = clCreateProgramWithSource(
//         context,
//         1,
//         &source,
//         sourceSize,
//         NULL);
//     //编译程序对象
//     status = clBuildProgram(program, 1, &device, NULL, NULL, NULL);
//     if (status)
//         std::cout << status << "  !!!!!!!!" <<std::endl;
//     if (status != 0)
//     {
//         printf("clBuild failed:%d\n", status);
//         char tbuf[0x10000];
//         clGetProgramBuildInfo(program, device, CL_PROGRAM_BUILD_LOG, 0x10000, tbuf,
//             NULL);
//         printf("\n%s\n", tbuf);
//         //return −1;
//     }
    
//     // 7. kernel
//     //创建 Kernel 对象
//     cl_kernel kernel = clCreateKernel(program, "init_vehicle", NULL);
//     status = clSetKernelArg(kernel, 0, sizeof(cl_mem), (void *)&memObjects[0]);
//     status = clSetKernelArg(kernel, 1, sizeof(cl_mem), (void *)&memObjects[1]);
//     status = clSetKernelArg(kernel, 2, sizeof(cl_mem), (void *)memObject);
//     //std::cout<<"status:"<<status<<std::endl;
//     if (status)
//         std::cout << "参数设置错误" << std::endl;
//     size_t global[1];
//     size_t local[1];
//     cl_event prof_event;
//     global[0] = VehicleSize;
//     local[0] = 25;
//     status = clEnqueueNDRangeKernel(commandQueue, kernel, 1, NULL,
//              global, local, 0, NULL, &prof_event);
//     clFinish(commandQueue);
//     if (status)
//         std::cout << "执行内核时错误" << std::endl;
//     clReleaseMemObject(memObjects[0]);
//     return memObjects[1];
// }