#include "pgAndKernel.h"
#include <unordered_map>
#include <fstream>
#include <sys/types.h>
#include <dirent.h>
#include <sys/stat.h>
#include <iostream>
#include <string.h>

ProgramController::ProgramController(std::string dir, cl_context context) {
    loadAllProgramByDir(dir, context);
}

void ProgramController::loadAllProgramByDir(std::string dir, cl_context context) {
    const char *path = dir.c_str();
    DIR *d = NULL;
    struct dirent *dp = NULL;
    struct stat st;    
    char p[1024] = {0};
    if(stat(path, &st) < 0 || !S_ISDIR(st.st_mode)) {
        std::cerr << "Wrong program directory" << std::endl;
        exit(-1);
    }
    if(!(d = opendir(path))) {  // 返回文件夹句柄
        std::cerr << "Can't Open program dir" << std::endl;
        exit(-1);
    }
    while((dp = readdir(d)) != NULL) {  // readdir基于文件夹句柄返回各个文件夹条目内容
        if((!strncmp(dp->d_name, ".", 1)) || (!strncmp(dp->d_name, "..", 2)) || (strncmp(dp->d_name, "program", 7)))
            continue;
        snprintf(p, sizeof(p) - 1, "%s/%s", path, dp->d_name);
        stat(p, &st);
        if(!S_ISDIR(st.st_mode)) {  // 如果不是文件夹
            loadProgramFromFile(p, context);
        }
    }
}

// void ProgramController::loadAllProgramByDir(std::string dir, cl_context context) {
//     const char *path = dir.c_str();
//     DIR *d = NULL;
//     struct dirent *dp = NULL;
//     struct stat st;   
//     std::vector<std::string> all_files; 
//     char p[1024] = {0};
//     if(stat(path, &st) < 0 || !S_ISDIR(st.st_mode)) {
//         std::cerr << "Wrong program directory" << std::endl;
//         exit(-1);
//     }
//     if(!(d = opendir(path))) {  // 返回文件夹句柄
//         std::cerr << "Can't Open program dir" << std::endl;
//         exit(-1);
//     }
//     while((dp = readdir(d)) != NULL) {  // readdir基于文件夹句柄返回各个文件夹条目内容
//         if((!strncmp(dp->d_name, ".", 1)) || (!strncmp(dp->d_name, "..", 2)))
//             continue;
//         snprintf(p, sizeof(p) - 1, "%s/%s", path, dp->d_name);
//         stat(p, &st);
//         if(!S_ISDIR(st.st_mode)) {  // 如果不是文件夹
//             all_files.push_back(p);
//         }
//         memset(p, 0, 1024);
//     }
//     std::vector<std::string>::iterator it = all_files.begin();
//     while (it != all_files.end())
//     {
//         std::cout << *it << std::endl;
//         it++; 
//     }
    
//     loadProgramAll(all_files, context);
// }

// void ProgramController::loadProgramAll(std::vector<std::string> all_files, cl_context context) {
//     // 读取文件内容
//     FILE *program_handles[all_files.size()];
//     size_t program_size = 0;
//     char *program_buffer;
//     std::string file_name;
//     for (size_t i = 0; i < all_files.size(); i++)
//     {
//         file_name = all_files[i];
//         program_handles[i] = fopen(file_name.c_str(), "r");
//         if(program_handles[i] == NULL) {
//             std::cerr << "Couldn't find the program file " << file_name << std::endl;
//             exit(-1);   
//         }
//         fseek(program_handles[i], 0, SEEK_END);
//         program_size += ftell(program_handles[i]);
//         rewind(program_handles[i]);
//     }
//     program_buffer = (char*)malloc(program_size + 1);
//     program_buffer[program_size] = '\0';
//     size_t index = 0;
//     for (size_t i = 0; i < all_files.size(); i++)
//     {
//         index = fread(program_buffer+index, sizeof(char), program_size, program_handles[i]);
//         fclose(program_handles[i]);
//     }
//     // 创建并激活program对象
//     cl_program pg = createProgram(context, program_buffer);
    
//     // 将program插入对应的管理容器中
//     this->allPrograms.insert(std::make_pair(file_name, pg));
//     // 创建该program中的所有kernel
//     char name_buffer[1024] = {'\0'};
//     cl_uint num_kernel = 0;
//     clGetProgramInfo(pg, CL_PROGRAM_KERNEL_NAMES, sizeof(name_buffer), (void *)name_buffer, NULL);
//     if (name_buffer[0] != '\0')
//     {
//         for (size_t i = 0; i < 1024; i++)
//         {
//             if (name_buffer[i] == ';')
//             {
//                 ++num_kernel;
//             }
//         }
//         ++num_kernel;
//     }
//     cl_kernel kernels[num_kernel];
//     clCreateKernelsInProgram(pg, num_kernel, kernels, NULL);
//     if (num_kernel == 0)
//     {
//         std::cerr << "No kernel in file: " << file_name << std::endl;
//         std::cerr << "This can caused by the way of writing kernel file" << std::endl;
//         exit(-1);
//     }
//     // 遍历所有kernel并将其插入到管理容器中
//     for (size_t i = 0; i < num_kernel; i++)
//     {
//         memset(name_buffer, '\0', sizeof(name_buffer));
//         clGetKernelInfo(kernels[i], CL_KERNEL_FUNCTION_NAME, sizeof(name_buffer), (void *)name_buffer, NULL);
//         this->allKernels.insert(std::make_pair(name_buffer, kernels[i]));
//     }
// }

void ProgramController::loadProgramFromFile(std::string file_name, cl_context context) {
    // 读取文件内容
    std::cout << "Creating program for " << file_name << std::endl;
    FILE *program_handle;
    size_t program_size;
    char *program_buffer;
    program_handle = fopen(file_name.c_str(), "r");
    if(program_handle == NULL) {
       std::cerr << "Couldn't find the program file" << std::endl;
       exit(-1);   
    }
    fseek(program_handle, 0, SEEK_END);
    program_size = ftell(program_handle);
    rewind(program_handle);
    program_buffer = (char*)malloc(program_size + 1);
    program_buffer[program_size] = '\0';
    fread(program_buffer, sizeof(char), program_size, program_handle);
    fclose(program_handle);
    // 创建并激活program对象
    cl_program pg = createProgram(context, program_buffer);
    
    // 将program插入对应的管理容器中
    this->allPrograms.insert(std::make_pair(file_name, pg));
    // 创建该program中的所有kernel
    char name_buffer[1024] = {'\0'};
    cl_uint num_kernel = 0;
    clGetProgramInfo(pg, CL_PROGRAM_KERNEL_NAMES, sizeof(name_buffer), (void *)name_buffer, NULL);
    if (name_buffer[0] != '\0')
    {
        for (size_t i = 0; i < 1024; i++)
        {
            if (name_buffer[i] == ';')
            {
                ++num_kernel;
            }
        }
        ++num_kernel;
    }
    cl_kernel kernels[num_kernel];
    clCreateKernelsInProgram(pg, num_kernel, kernels, NULL);
    if (num_kernel == 0)
    {
        std::cerr << "No kernel in file: " << file_name << std::endl;
        std::cerr << "This can caused by the way of writing kernel file" << std::endl;
        exit(-1);
    }
    // 遍历所有kernel并将其插入到管理容器中
    for (size_t i = 0; i < num_kernel; i++)
    {
        memset(name_buffer, '\0', sizeof(name_buffer));
        clGetKernelInfo(kernels[i], CL_KERNEL_FUNCTION_NAME, sizeof(name_buffer), (void *)name_buffer, NULL);
        this->allKernels.insert(std::make_pair(name_buffer, kernels[i]));
    }
}

cl_program ProgramController::createProgram(cl_context context, char *program_source) {
    cl_int err;
    cl_program program = clCreateProgramWithSource(context, 1, (const char **)&program_source, NULL, &err);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Wrong when creating program." << std::endl;
        exit(-1);
    }
    err = clBuildProgram(program, 0, NULL, "-I ../programs/", NULL, NULL);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Error when building program " << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
    return program;
}

void ProgramController::setKernelArg(std::string kName, uint argIndex, uint basic_size, uint num, const void * ptr) {
    cl_kernel kernel = this->allKernels[kName];
    if (kernel == nullptr)
    {
        std::cerr << "No such kernel" << std::endl;
        exit(-1);
    }
    cl_int err;
    err = clSetKernelArg(kernel, argIndex, basic_size*num, ptr);
    if (err != CL_SUCCESS)
    {
        std::cerr << "Can't set arg" << std::endl;
        std::cerr << "Error code: " << err << std::endl;
        exit(-1);
    }
}

cl_kernel ProgramController::getKernelByName(std::string kernel_name) {
    return this->allKernels[kernel_name];
}

void ProgramController::execKernelNDRangeMode(cl_command_queue cq, std::string kernel_name, cl_uint dimension, 
                                size_t *global_work_size, size_t *global_work_offset, 
                                size_t *local_work_size) {
    cl_kernel kernel = this->allKernels[kernel_name];
    cl_int err = clEnqueueNDRangeKernel(cq, kernel, dimension, global_work_offset, global_work_size, local_work_size, 0, NULL, NULL);
    if (err != CL_SUCCESS)
    {
        std::cout << "Error when enqueuing kernel: " << kernel_name << std::endl;
        std::cout << "Error code: " << err << std::endl;
        exit(-1);
    }
}

void ProgramController::execKernelTaskMode(cl_command_queue cq, std::string kernel_name) {
    cl_kernel kernel = this->allKernels[kernel_name];
    cl_int err = clEnqueueTask(cq, kernel, 0, NULL, NULL);
    if (err != CL_SUCCESS)
    {
        std::cout << "Error when enqueuing task: " << kernel_name << std::endl;
        std::cout << "Error code: " << err << std::endl;
        exit(-1);
    }
}