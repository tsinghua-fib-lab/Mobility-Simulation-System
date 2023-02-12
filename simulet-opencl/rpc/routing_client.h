#ifndef ROUTING_CLIENT
#define ROUTING_CLIENT
#include <grpcpp/grpcpp.h>
#include <condition_variable>
#include <memory>
#include <mutex>
#include <string>
#include <thread>
#include "proto_use.h"

/// 路径规划gRPC Client
class RoutingClient {
 public:
  struct AsyncClientCall {
    // Vehicle被销毁或请求被放弃
    void GiveUp();
    /* contain message from server */
    simulet::PbGetRouteResponse response;
    grpc::ClientContext context;
    grpc::Status status;
    std::unique_ptr<grpc::ClientAsyncResponseReader<simulet::PbGetRouteResponse>>
      response_reader;
    // 指向发出此次请求的车辆
    int vehicleIndex = -1;
    int startAoiId;
    simulet::PbPosition end_position;
    // 路径规划有效性互斥锁
    std::mutex mtx;
    // 标记此次请求是否有效
    bool valid = true;
  };
  explicit RoutingClient(std::shared_ptr<grpc::Channel> channel);
  AsyncClientCall* GetRoute(const simulet::PbGetRouteRequest& request,
                            int vehicleIndex,
                            int startAoiId,
                            simulet::PbPosition end_position);
  /// RoutingClient处理响应的例程
  void AsyncCompleteRpc();
  /// 关闭RoutingClient
  void Shutdown();
  /// 等待上一个step发出的所有请求均处理完成
  void Wait();

 private:
  std::unique_ptr<::wolong::routing::v2::RoutingService::Stub> stub_;
  grpc::CompletionQueue cq_;
  // 记录正在等待的请求个数
  int count_waiting_req_ = 0;
  // 路径规划完成等待互斥锁
  std::mutex mtx_;
  // 路径规划完成等待条件变量
  std::condition_variable cv_;
  // 路径规划接收与处理线程
  std::thread worker_;
  //操作 segmentIndex 和 infoIndex 的互斥锁
  std::mutex m_mtx;
};

/// 全局导航组件
extern std::shared_ptr<RoutingClient> routing_client;

#endif