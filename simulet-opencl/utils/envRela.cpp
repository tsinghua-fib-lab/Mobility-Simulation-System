#include "envRela.h"
#include <iostream>
#include <cstdlib>
#include <string.h>

Context::Context(std::string dir, int device_type, int device_num) {
    initPlatform();
    initDevices(this->platform, this->devices, device_type, device_num);
    createContext(this->devices, device_num);
    for (size_t i = 0; i < device_num; i++)
    {
        DeviceQueue dq = {this->devices[i], nullptr};
        this->commandQueueByDeviceId[i] = dq;
    }
    initMemoryController(this->context);
    initProgramController(dir);
}

void Context::initPlatform() {
    char versionBuffer[1024];
    char nameBuffer[1024];
    cl_uint num;
    clGetPlatformIDs(0, NULL, &num);
    std::cout << "There are " << num << " available platform" << std::endl;
    cl_platform_id platforms[num];
    clGetPlatformIDs(num, platforms, NULL);
    uint numPlatCompatible = 0;
    for (size_t i = 0; i < num; i++)
    {
        memset(versionBuffer, '\0', sizeof(versionBuffer));
        memset(nameBuffer, '\0', sizeof(nameBuffer));
        clGetPlatformInfo(platforms[i], CL_PLATFORM_VERSION, sizeof(versionBuffer), (void *)versionBuffer, NULL);
        clGetPlatformInfo(platforms[i], CL_PLATFORM_NAME, sizeof(nameBuffer), (void *)nameBuffer, NULL);
        if (atof(versionBuffer+6) >= PLATFORM_VERSION)
        {
            ++numPlatCompatible;
            std::cout << i+1 << ". " << nameBuffer << ", Version: " << versionBuffer << std::endl;
        }
    }
    if (numPlatCompatible > 0)
    {
        std::cout << "There are totally " << numPlatCompatible << " platforms compatible with OpenCL 1.2" << std::endl;
        std::cout << "Please select a platform to use: (waiting your choice) ";
        int index;
        std::cin >> index;
        if (index > numPlatCompatible)
        {
            std::cerr << "Invalid platform id" << std::endl;
            exit(-1);
        }
        this->platform = platforms[index-1];
    } else 
    {
        std::cerr << "There are no compatible platform" << std::endl;
        exit(-1);
    }
}

void Context::initDevices(cl_platform_id platform, cl_device_id *devices, int deviceType, int num_device) {
    cl_uint num;
    uint types[5] = {CL_DEVICE_TYPE_CPU, 
                     CL_DEVICE_TYPE_ACCELERATOR, CL_DEVICE_TYPE_CUSTOM, 
                     CL_DEVICE_TYPE_DEFAULT,  CL_DEVICE_TYPE_ALL};
    std::cout << "Init platform related device......" << std::endl;
    cl_int ret = clGetDeviceIDs(platform, deviceType, 0, NULL, &num);
    if (ret == CL_INVALID_PLATFORM)
    {
        std::cerr << "Using invalid platform when getting device." << std::endl;
        std::cout << "Error code: " << ret << std::endl;
        exit(-1);
    } else if (ret == CL_OUT_OF_RESOURCES)
    {
        std::cerr << "There is a failure to allocate resources required by the OpenCL implementation on the device" << std::endl;
        std::cout << "Error code: " << ret << std::endl;
        exit(-1);
    } else if (ret == CL_OUT_OF_HOST_MEMORY)
    {
        std::cerr << "There is a failure to allocate resources required by the OpenCL implementation on the host" << std::endl;
        std::cout << "Error code: " << ret << std::endl;
        exit(-1);
    }
    int i = 0;
    while (ret == CL_DEVICE_NOT_FOUND && i < 5)
    {
        ret = clGetDeviceIDs(platform, types[i], 0, NULL, &num);
        ++i;
    }
    if (i == 5 && ret == CL_DEVICE_NOT_FOUND)
    {
        std::cerr << "Error when finding device." << std::endl;
        exit(-1);
    } else {
        if (i == 0)  // GPU
        {
            std::cout << "There totally " << num << " device of specific type (GPU)" << std::endl;
            if (num < num_device)
            {
                std::cerr << "There are not enough devices to use" << std::endl;
                exit(-1);
            }
            cl_int err;
            err = clGetDeviceIDs(platform, deviceType, num_device, devices, NULL);
            if (err != CL_SUCCESS)
            {
                std::cerr << "Error when getting devices" << std::endl;
                std::cerr << "Error code: " << err << std::endl;
                exit(-1);
            }
        } else  // other device type
        {
            --i;
            std::cout << "There totally " << num << " device of specific type " << i << std::endl;
            if (num < num_device)
            {
                std::cerr << "There are not enough devices to use" << std::endl;
                exit(-1);
            }
            cl_int err;
            err = clGetDeviceIDs(platform, types[i], num_device, devices, NULL);
            if (err != CL_SUCCESS)
            {
                std::cerr << "Error when getting devices" << std::endl;
                std::cerr << "Error code: " << err << std::endl;
                exit(-1);
            }
        }
    }
}

void pfn_notify(const char *errinfo, const void *private_info, size_t cb, void *user_data) {
    std::cout << "***********************" << std::endl;
    std::cout << errinfo << std::endl;
    std::cout << "***********************" << std::endl;
}

void Context::createContext(cl_device_id *devices, int device_num) {
    cl_int err;
    cl_context context = clCreateContext(NULL, device_num, devices, pfn_notify, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Error when creating context" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    this->context = context;
}

void Context::initMemoryController(cl_context Context) {
    this->mController = new MemoryController(context);
}

cl_mem Context::getRouteHeap() {
    return mController->doGetRouteHeap();
}

cl_mem Context::getBitmap() {
    return mController->doGetBitmap();
}

void Context::registerCommandQueue(int device_id, int flags) {
    cl_int err;
    cl_device_id device = this->commandQueueByDeviceId[device_id].device;
    cl_context context = this->context;
    cl_command_queue cq = clCreateCommandQueue(context, device, flags, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Error when creating command queue" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    this->commandQueueByDeviceId[device_id].commandQueue = cq;
}

void Context::initRouteHeap() {
    cl_command_queue cq = getCommandQueueByDeviceId(0);
    mController->doInitRouteHeap(cq);
}

void Context::createBuffer(std::string buffer_name, size_t unit_size, uint unit_num) {
    this->mController->createBuffer(buffer_name, unit_size, unit_num);
}

void Context::createReadOnlyBuffer(std::string buffer_name, size_t unit_size, uint unit_num) {
    this->mController->createReadOnlyBuffer(buffer_name, unit_size, unit_num);
}

void Context::createWriteOnlyBuffer(std::string buffer_name, size_t unit_size, uint unit_num) {
    this->mController->createWriteOnlyBuffer(buffer_name, unit_size, unit_num);
}

void Context::createBufferAndInit(std::string buffer_name, size_t unit_size, uint unit_num, void *ptr) {
    this->mController->createBufferAndInit(buffer_name, unit_size, unit_num, ptr);
}

void Context::releaseBuffer(std::string bName) {
    this->mController->releaseBuffer(bName);
}

void Context::deviceMalloc(std::string bName, size_t unit_size, uint unit_num) {
    this->mController->doDeviceMalloc(bName, unit_size, unit_num);
}

void Context::deviceFree(std::string bName) {
    this->mController->doDeviceFree(bName);
}

void Context::readBuffer(int deviceId, std::string buffer_name, void *ptr, bool block, size_t size, size_t offset) {
    cl_command_queue cq = getCommandQueueByDeviceId(deviceId);
    this->mController->readOrWriteBuffer(true, cq, buffer_name, ptr, block, offset, size);
}

void Context::WriteBuffer(int deviceId, std::string buffer_name, void *ptr, bool block, size_t size, size_t offset) {
    cl_command_queue cq = getCommandQueueByDeviceId(deviceId);
    this->mController->readOrWriteBuffer(false, cq, buffer_name, ptr, block, offset, size);
}

void Context::copyBetweenBuffer(int deviceId, std::string src_buffer, std::string des_buffer, size_t size, size_t src_offset, size_t des_offset){
    cl_command_queue cq = getCommandQueueByDeviceId(deviceId);
    this->mController->copyBetweenBuffer(cq, src_buffer, des_buffer, size, src_offset, des_offset);
}

void Context::initProgramController(std::string dir) {
    this->pController = new ProgramController(dir, this->context);
}

cl_kernel Context::getKernelByName(std::string kernel_name) {
    return this->pController->getKernelByName(kernel_name);
}

void Context::setBufferAsKernelArg(std::string kernel_name, uint argIndex, ArgType type, std::string buffer_name) {
    BufferObj buffer = this->mController->getBufferObjByName(buffer_name);
    cl_mem bufferObject = buffer.memObj;
    this->pController->setKernelArg(kernel_name, argIndex, type.unit_size, type.unit_num, &bufferObject);
}

void Context::execKernelNDRangeMode(int deviceId, std::string kernel_name, NDRange ndMessage) {
    cl_command_queue cq = getCommandQueueByDeviceId(deviceId);
    this->pController->execKernelNDRangeMode(cq, kernel_name, ndMessage.dimension, 
                ndMessage.global_work_items, ndMessage.global_offset, ndMessage.local_work_items);
}

cl_command_queue Context::getCommandQueueByDeviceId(int deviceId) {
    return this->commandQueueByDeviceId[deviceId].commandQueue;
}