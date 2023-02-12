#include "lane.h"
#include "vehicle.h"
#include "aoi.h"

typedef struct {
    int globalTime;
    int vehicleIndex[820000];
    int vehicleInAoi;
    int aoi2Lane[1024];
    int lane2Aoi[1024];
    DeviceVehicle *vehicles;
    DeviceAoi *aois;
    DeviceLane *lanes;
    void *heap;
}Points;