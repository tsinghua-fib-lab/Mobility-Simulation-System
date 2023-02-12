#ifndef ENVRELA_H
#define ENVRELA_H
#include "CL/opencl.h"
#include "pgAndKernel.h"
#include "mm.h"
#include <vector>

#define MAX_DEVICES 10
#define PLATFORM_VERSION 1.2
#define DEFAULT_DEVICE_TYPE CL_DEVICE_TYPE_GPU

typedef struct {
    cl_device_id device;
    cl_command_queue commandQueue;
} DeviceQueue;

class Context{
public:
    Context(std::string dir, int=CL_DEVICE_TYPE_GPU, int=1);
    /*根据传入的设备Id，获取对应的设备对象*/
    cl_device_id getDeviceById(int);
    /*根据传入的kernel name，获取对应的kernel对象*/
    cl_kernel getKernelByName(std::string);
    /*根据传入的设备Id，获取对应的command queue对象（因为command queue总是与一个device想绑定的）*/
    cl_command_queue getCommandQueueByDeviceId(int);
    /*将command queue注册到对应的设备上
        arg1: 设备Id
        arg2: flag标志，默认为0
    */
    void registerCommandQueue(int, int=0);
    /*创建buffer object
        arg1: 为该对象命名
        arg2: 该buffer存储对象的单位大小
        arg3: 单位个数
        例如: 需要创建一个存储100个vehicle对象的buffer, 则arg2=sizeof(vehicle), arg3=100
    */
    void createBuffer(std::string, size_t, uint=1);
    /*创建只读buffer object*/
    void createReadOnlyBuffer(std::string, size_t, uint=1);
    /*创建只写buffer object*/
    void createWriteOnlyBuffer(std::string, size_t, uint=1);
    /*在创建buffer object的同时进行初始化
        前三个参数与createBuffer相同
        arg4: 进行数据填充的host端的结构起点
    */
    void createBufferAndInit(std::string, size_t, uint, void *);
    /*释放指定buffer*/
    void releaseBuffer(std::string);
    /*用于在设备端通用heap分配数据对象*/
    void deviceMalloc(std::string, size_t, uint);
    /*用于在设备端通用heap释放数据对象*/
    void deviceFree(std::string);
    /*从指定的buffer object中读取数据*/
    void readBuffer(int, std::string, void *, bool=false, size_t=0, size_t=0);
    /*向指定的buffer object写入数据*/
    void WriteBuffer(int, std::string, void *, bool=false, size_t=0, size_t=0);
    /*在两个指定的buffer object之间拷贝数据*/
    void copyBetweenBuffer(int, std::string, std::string, size_t=0, size_t=0, size_t=0);
    /*初始化OpenCL programs以及kernels控制器*/
    void initProgramController(std::string);
    /*将指定的buffer object设定为指定kernel的参数*/
    void setBufferAsKernelArg(std::string, uint, ArgType, std::string);
    /*以NDRange的方式执行指定的kernel*/
    void execKernelNDRangeMode(int, std::string, NDRange={1, {1}, {1}, {0}});
    /*以Task的方式执行指定的kernel*/
    void execKernelTaskMode(int, std::string);
    /*将route heap注册到设备上*/
    void initRouteHeap();
    /*获取route heap数据对象*/
    cl_mem getRouteHeap();
    /*获取route bitmap数据对象*/
    cl_mem getBitmap();

private:
    void initPlatform();
    void initDevices(cl_platform_id, cl_device_id *, int, int);
    void createContext(cl_device_id *, int);
    void initMemoryController(cl_context);
    cl_platform_id platform;
    cl_device_id devices[MAX_DEVICES];
    cl_context context;
    std::unordered_map<int, DeviceQueue> commandQueueByDeviceId;
    ProgramController *pController;
    MemoryController *mController;
};

#endif