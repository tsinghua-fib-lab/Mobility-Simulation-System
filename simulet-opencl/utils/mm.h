#ifndef MM_H
#define MM_H

#include "CL/opencl.h"
#include <unordered_map>
#include <vector>

#define MEM_CONTROLLER_READ 1
#define MEM_CONTROLLER_WRITE 0
#define DEVICE_HEAP_SIZE 1024
#define ROUTE_SEGMENT_SIZE 32  // 一个segment的大小
#define ROUTE_SEGMENT_NUM 262144 * 128  // device端可以存储的segments的数量
#define ROUTE_HEAP_SIZE 262144 * 128 * 32  // vehicle * segments * size

typedef unsigned char byte;

typedef struct {
    cl_mem memObj;
    ulong size;
    ulong start;
    ulong end;
}BufferObj;

// 用于管理堆区的空闲空间
typedef struct FreePeriod{
    ulong start;
    ulong size;
    FreePeriod *next;
}FreePeriod;

typedef struct {
    int bitIndex;
    int ByteBefor;
}BitHelper;

class MemoryController {
public:
    MemoryController();
    MemoryController(cl_context);
    // 创建Memory对象(Buffer)
    void createBuffer(std::string, size_t, uint);  // 创建可读写buffer
    void createReadOnlyBuffer(std::string, size_t, uint);  // 创建只读buffer
    void createWriteOnlyBuffer(std::string, size_t, uint);  // 创建只写buffer
    // 从buffer中读写数据
    void readOrWriteBuffer(bool, cl_command_queue, std::string, void *, bool, size_t, size_t);
    // 在Buffer间拷贝内容
    void copyBetweenBuffer(cl_command_queue, std::string, std::string, size_t, size_t, size_t);
    void releaseBuffer(std::string);

    // 根据传入的buffer对象名获取对应的BufferObj对象
    BufferObj getBufferObjByName(std::string);

    // 为kernel设置memory对象相关参数
    void setArgForKernel(cl_kernel, uint, std::string);

    // 设备端堆相关
    void createBufferAndInit(std::string, size_t, uint, void *);  // 创建并直接初始化buffer
    void doDeviceMalloc(std::string, size_t, uint);  // 在设备“堆”开辟空间
    void doDeviceFree(std::string);  // 释放设备“堆”空间
    /*在route heap上开辟指定个segments大小的空间，返回首个segment对应的索引*/
    int mallocSegments(int num);
    void doInitRouteHeap(cl_command_queue);
    cl_mem doGetRouteHeap();
    cl_mem doGetBitmap();

private:
    cl_mem heap;  // 通用heap
    FreePeriod *free_head = new FreePeriod{0, 0, nullptr};  // 通用heap的空闲块链表表头
    size_t totalFreeSize;  // 通用heap的总大小

    cl_mem route_heap;  // route专用heap
    // std::vector<cl_mem> segments;  // 用于存储route中的每一个buffer object
    // byte hostBitmap[ROUTE_SEGMENT_NUM/8];  // host端的bitmap
    cl_mem deviceBitmap;  // 用于表示对应segment空间是否空闲的bitmap
    cl_mem bitIndex;
    
    cl_context context;  // 关联的context的对象
    std::unordered_map<std::string, BufferObj> memObjs; // 用于管理所有创建的buffer object
};

#endif