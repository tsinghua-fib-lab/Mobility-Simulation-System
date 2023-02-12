#include "mm.h"
#include <iostream>

MemoryController::MemoryController(cl_context cx): context(cx) {
    cl_int err;
    /*构建设备端堆空间*/
    this->heap = clCreateBuffer(context, CL_MEM_READ_WRITE, DEVICE_HEAP_SIZE, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when creating device heap (with size " << DEVICE_HEAP_SIZE << " bytes)" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    // 初始化表头
    this->free_head->next = new FreePeriod{0, DEVICE_HEAP_SIZE, nullptr};
    this->totalFreeSize = DEVICE_HEAP_SIZE;

    /*构建设备端route专用堆空间*/
    this->route_heap = clCreateBuffer(context, CL_MEM_READ_WRITE, ROUTE_HEAP_SIZE, nullptr, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when creating route specific heap (with size " << ROUTE_HEAP_SIZE << " bytes)" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    BufferObj obj = {this->route_heap, ROUTE_HEAP_SIZE};
    this->memObjs.insert(std::make_pair("routeHeap", obj));

    // 构建bitmap数据对象
    this->deviceBitmap = clCreateBuffer(context, CL_MEM_READ_WRITE, ROUTE_SEGMENT_NUM/8, nullptr, &err);  // 每个bit对应了一个segments对象
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when creating device bitmap" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    obj = {this->deviceBitmap, ROUTE_SEGMENT_NUM/8};
    this->memObjs.insert(std::make_pair("bitmap", obj));

    // 构建bitIndex数据对象
    this->bitIndex = clCreateBuffer(context, CL_MEM_READ_WRITE, sizeof(int), nullptr, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when creating device bitIndex" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    obj = {this->bitIndex, sizeof(int)};
    this->memObjs.insert(std::make_pair("bitIndex", obj));
}

void MemoryController::createBuffer(std::string bName, size_t basic_size, uint num) {
    cl_int err;
    cl_mem bufferObj = clCreateBuffer(this->context, CL_MEM_READ_WRITE, basic_size*num, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't create buffer: " << bName << std::endl;
        std::cerr << "error code: " << err << std::endl;
        exit(-1);
    }
    BufferObj obj = {bufferObj, basic_size*num};
    this->memObjs.insert(std::make_pair(bName, obj));
}

void MemoryController::createReadOnlyBuffer(std::string bName, size_t basic_size, uint num) {
    cl_int err;
    cl_mem bufferObj = clCreateBuffer(this->context, CL_MEM_READ_ONLY, basic_size*num, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't create read only buffer: " << bName << std::endl;
        std::cerr << "error code: " << err << std::endl;
        exit(-1);
    }
    BufferObj obj = {bufferObj, basic_size*num};
    this->memObjs.insert(std::make_pair(bName, obj));
}

void MemoryController::createWriteOnlyBuffer(std::string bName, size_t basic_size, uint num) {
    cl_int err;
    cl_mem bufferObj = clCreateBuffer(this->context, CL_MEM_WRITE_ONLY, basic_size*num, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't create write only buffer: " << bName << std::endl;
        std::cerr << "error code: " << err << std::endl;
        exit(-1);
    }
    BufferObj obj = {bufferObj, basic_size*num};
    this->memObjs.insert(std::make_pair(bName, obj));
}

/** 在创建Buffer对象的同时进行初始化
 * 初始化方法为指定host端数据地址ptr，以该地址内数据为初始化数据
 **/
void MemoryController::createBufferAndInit(std::string bName, size_t basic_size, uint num, void *ptr) {
    cl_int err;
    cl_mem bufferObj = clCreateBuffer(this->context, CL_MEM_COPY_HOST_PTR | CL_MEM_ALLOC_HOST_PTR, 
                                basic_size*num, ptr, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't create a buffer and init it directly." << std::endl;
        std::cerr << "error code: " << err << std::endl;
        exit(-1);
    }
    BufferObj obj = {bufferObj, basic_size*num};
    this->memObjs.insert(std::make_pair(bName, obj));
}

void MemoryController::doDeviceMalloc(std::string bName, size_t basic_size, uint num) {
    size_t size = basic_size * num;
    if (this->totalFreeSize < size)
    {
        std::cerr << "There has no enough heap space (size " << size << " )" << std::endl;
        exit(-1);
    }
    FreePeriod *now = this->free_head->next;
    while (now != nullptr)
    {
        if (now->size >= size)  // 如果当前FreePeriod的空间足够使用
        {
            // 子空间创建
            cl_int err;
            _cl_buffer_region region = _cl_buffer_region{now->start, size};
            cl_mem sub = clCreateSubBuffer(this->heap, CL_MEM_READ_WRITE, CL_BUFFER_CREATE_TYPE_REGION, &region, &err);
            BufferObj obj = {sub, size, now->start, now->start + size};
            this->memObjs.insert(std::make_pair(bName, obj));

            // 链表处理
            size_t left = now->size - size;
            now->size = left;
            now->start = now->start + size;
            break;
        }
        now = now->next;
    }
}

void MemoryController::doDeviceFree(std::string bName) {
    BufferObj obj = this->getBufferObjByName(bName);
    cl_int err = clReleaseMemObject(obj.memObj);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when release heap object" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    // 链表处理
    FreePeriod *prior = this->free_head;
    FreePeriod *now = this->free_head->next;
    bool flag = true;
    while (flag)
    {
        if (now->start > obj.start + obj.size) // 独立
        {
            prior->next = new FreePeriod{obj.start, obj.size, now};
            flag = false;
        }else if (obj.end == now->start)  // 尾头相接
        {
            now->start = obj.start;
            now->size += obj.size;
            flag = false;
        } else if (obj.start == now->start + now->size)  // 头尾相接
        {
            now->size += obj.size;
            FreePeriod *next = now->next;
            if (next->start == now->start + now->size)  // 合并
            {
                now->size += next->size;
                now->next = next->next;
                free(next);
            }
            flag = false;
        }
        // 切换
        prior = now;
        now = now->next;
    }
    this->memObjs.erase(bName);
}

void MemoryController::readOrWriteBuffer(bool read, cl_command_queue cq, std::string bName, void *ptr, bool block, size_t offset, size_t size) {
    cl_int err;
    BufferObj obj = this->memObjs.at(bName);
    if (size == 0)
    {
        // 如果没有指定要读取的长度，则全部读取
        size = this->memObjs[bName].size - offset;
    }
    if (read)  // 读操作
    {
        err = clEnqueueReadBuffer(cq, this->memObjs[bName].memObj, block, offset, size, ptr, 0, NULL, NULL);
    } else
    {
        err = clEnqueueWriteBuffer(cq, this->memObjs[bName].memObj, block, offset, size, ptr, 0, NULL, NULL);
    }
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't read buffer( " << bName << " ) into ( " << ptr << " )" << std::endl;
        std::cerr << "error code: " << err << std::endl;
        exit(-1);
    }
}

void MemoryController::copyBetweenBuffer(cl_command_queue cq, std::string srcName, std::string desName, size_t size, size_t src_offset, size_t des_offset) {
    cl_int err;
    BufferObj srcObj = this->memObjs[srcName];
    BufferObj desObj = this->memObjs[desName];
    if (src_offset >= srcObj.size || des_offset >= desObj.size)
    {
        std::cerr << "Wrong offset" << std::endl;
        exit(-1);
    } else if (size == 0)
    {
        // 二者取小
        size_t size_src = srcObj.size - src_offset;
        size_t size_des = desObj.size - des_offset;
        size = size_src <= size_des ? size_src : size_des;
    }
    err = clEnqueueCopyBuffer(cq, srcObj.memObj, desObj.memObj, src_offset, des_offset, size, 0, NULL, NULL);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't copy data from " << srcName << " to " << desName << "." << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
}

BufferObj MemoryController::getBufferObjByName(std::string bName) {
    return this->memObjs[bName];
}

void MemoryController::setArgForKernel(cl_kernel kernel, uint argIndex, std::string bName) {
    BufferObj bufferObj = getBufferObjByName(bName);
    cl_int err;
    err = clSetKernelArg(kernel, argIndex, sizeof(cl_mem), (void *)&bufferObj.memObj);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Error when set arg" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
}

void MemoryController::doInitRouteHeap(cl_command_queue cq) {
    cl_int pattern = 0;
    cl_int err = clEnqueueFillBuffer(cq, this->deviceBitmap, &pattern, sizeof(cl_int), 0, ROUTE_SEGMENT_NUM/8, 0, nullptr, nullptr);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when init bitmap" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    err = clEnqueueFillBuffer(cq, this->route_heap, &pattern, sizeof(cl_int), 0, ROUTE_HEAP_SIZE, 0, nullptr, nullptr);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when init route heap" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
}

cl_mem MemoryController::doGetRouteHeap() {
    return this->route_heap;
}

cl_mem MemoryController::doGetBitmap() {
    return this->deviceBitmap;
}

void MemoryController::releaseBuffer(std::string bName) {
    cl_mem buffer = getBufferObjByName(bName).memObj;
    cl_int err = clReleaseMemObject(buffer);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when release buffer object" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    this->memObjs.erase(bName);
}