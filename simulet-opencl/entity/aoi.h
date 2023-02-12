#ifndef AOI_H
#define AOI_H
#include"lane.h"
#include <unordered_map>
// #include <map>
const int MAX_DRIVING_POSITIONS = 150;
const int MAX_POSITIONS = 350;

typedef struct XYPosition {
  double x;
  double y;
}XYPosition;
typedef struct AoiVehicle{
    int vehicle_index[2000];
    int number;
} AoiVehicle;
//aoi中车的链表
typedef struct HostAoi{
    int AoiID;
    //Aoi与行车路网的连接点数量
    int DrivingPositionSize;
    // Aoi与行车路网的连接点
    LanePosition DrivingPositions[MAX_DRIVING_POSITIONS];
    //存储Aoi与行车路网的连接点laneID到坐标s的映射 
    std::unordered_map<int, double> lanS;
    // std::map<int, double> lanS;
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
}HostAoi;
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
} DeviceAoi;
cl_mem initDeviceAoi(DeviceAoi *aoi, AoiVehicle *aoiVehicle, Context *cx);
cl_mem InsertVehicleInAoi(Context *cx);
cl_mem RemoveVehicleInAoi(Context *cx);
cl_mem VehicleGoLane(Context *cx,unsigned long size);
#endif
