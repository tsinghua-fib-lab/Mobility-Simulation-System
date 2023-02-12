#ifndef PGANDKER_H
#define PGANDKER_H

#include "CL/opencl.h"
#include <vector>
#include <string>
#include <unordered_map>

#define MAX_KERNELS_IN_ONE_FILE 20

typedef struct {
    size_t unit_size;
    size_t unit_num;
}ArgType;

const ArgType POINTER = { sizeof(cl_mem), 1 };

typedef struct {
    int dimension;
    size_t global_work_items[3];
    size_t local_work_items[3];
    size_t global_offset[3];
} NDRange;

class ProgramController {
public:
    ProgramController();
    ProgramController(std::string, cl_context);
    void setKernelArg(std::string, uint, uint, uint, const void *);
    cl_kernel getKernelByName(std::string);
    void execKernelNDRangeMode(cl_command_queue, std::string, cl_uint, size_t *, size_t *, size_t *);
    void execKernelTaskMode(cl_command_queue, std::string);
private:
    void loadAllProgramByDir(std::string, cl_context);
    void loadProgramFromFile(std::string, cl_context);
    void loadProgramAll(std::vector<std::string>, cl_context);
    cl_program createProgram(cl_context, char *);
    std::unordered_map<std::string, cl_program> allPrograms;
    std::unordered_map<std::string, cl_kernel> allKernels;
};

#endif