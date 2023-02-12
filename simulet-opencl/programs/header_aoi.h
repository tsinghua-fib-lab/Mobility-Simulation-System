#ifndef AOI_H
#define AOI_H
// #include <unordered_map>
#include"./programs/header_lane.h"
#define MAX_DRIVING_POSITIONS 150
#define MAX_POSITIONS 350

// typedef struct DeviceVehicle;

/**Aoi相关数据结构 */
typedef struct XYPosition {
  double x;
  double y;
}XYPosition;

typedef struct AoiVehicle{
    int vehicle_index[2000];
    int number;
}AoiVehicle;

typedef struct DeviceAoi{
    int AoiID;
    //Aoi与行车路网的连接点数量
    int DrivingPositionSize;
    // Aoi与行车路网的连接点
    LanePosition DrivingPositions[MAX_DRIVING_POSITIONS];
    // Aoi原始位置数量
    int PositionSize;
    // Aoi原始位置（如果只有一个值，则为Aoi所在的点，否则为Aoi多边形的边界）
    XYPosition Positions[MAX_POSITIONS];
    //Aoi与行车路网连接时在自身边界上的连接点数量
    int DrivingGateSize;
    // Aoi与行车路网连接时在自身边界上的连接点, 与driving_positions按索引一一对应
    XYPosition DrivingGates[MAX_DRIVING_POSITIONS];
    // Aoi面积, 若是Poi则无此字段
    double area;
    //AOI中车的数量
    int vehicleNum;
    //Aoi中车辆链表头指针
    struct DeviceVehicle *vehicleListStart;
    //Aoi中车辆链表尾指针
    struct DeviceVehicle *vehicleListEnd;
    //进AOI车的数量
    int insertNum;
    //存储进AOI的车的指针
    struct DeviceVehicle* insertBuff[1024];
    //出AOI车的数量
    int removeNum;
    //存储出AOI车的指针
    struct DeviceVehicle* removeBuff[1024];
}DeviceAoi;

#endif