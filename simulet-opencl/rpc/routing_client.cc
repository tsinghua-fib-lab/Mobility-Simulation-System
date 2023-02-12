#include "rpc/routing_client.h"
#include <grpcpp/grpcpp.h>
// #include <spdlog/spdlog.h>
#include <condition_variable>
#include <memory>
#include <mutex>
#include <string>
#include "proto_use.h"
#include "simulet.h"

void RoutingClient::AsyncClientCall::GiveUp() {
  std::unique_lock<std::mutex> lock(mtx);
  valid = false;
}

RoutingClient::RoutingClient(std::shared_ptr<::grpc::Channel> channel)
    : stub_(::wolong::routing::v2::RoutingService::NewStub(channel)),
      worker_(std::thread(&RoutingClient::AsyncCompleteRpc, this)) {}

RoutingClient::AsyncClientCall* RoutingClient:: GetRoute(
    const simulet::PbGetRouteRequest& request, int vehicleIndex,int startAoiId,
    simulet::PbPosition end_position) {
      AsyncClientCall* call = new AsyncClientCall;
      call->vehicleIndex = vehicleIndex;
      call->startAoiId = startAoiId;
      call->end_position = end_position;
      call->response_reader =
          stub_->PrepareAsyncGetRoute(&call->context, request, &cq_);
      call->response_reader->StartCall();
      call->response_reader->Finish(&call->response, &call->status,
                                    reinterpret_cast<void*>(call));
      std::unique_lock<std::mutex> lock(mtx_);
      ++count_waiting_req_;
      // std::cout<<"++count_waiting_req_:\t"<<count_waiting_req_<<std::endl;
      return call;
}

void RoutingClient::AsyncCompleteRpc() {
  AsyncClientCall* call;
  bool ok = false;
  while (cq_.Next(reinterpret_cast<void**>(&call), &ok)) {
    assert(ok);
    if (call->status.ok()) {
      std::unique_lock<std::mutex> lock(call->mtx);
      // 检查消息是否有效
      if (call->valid) {
        const simulet::PbGetRouteResponse& res = call->response;
        if(!res.journeys().empty()){
          assert(res.journeys().size() == 1);
          assert(res.journeys().at(0).type() == simulet::PbJourneyType::JOURNEY_TYPE_DRIVING);
          assert(res.journeys().at(0).has_driving());
          assert(!res.journeys().at(0).driving().route().empty());
          const auto& pb = res.journeys().at(0).driving().route();
          int size = pb.size(); //导航routing返回的路径的大小
          
          HostSegment localSegments[128];  //在局部变量 localSegments 准备数据
          //复用一个 segment
          localSegments[0].laneIndex = size;
          for (int i = 0, size = pb.size(); i < size; ++i) {
            localSegments[i+1].laneIndex = pb.at(i).lane_id();
            localSegments[i+1].nextLaneType = (cl_int)pb.at(i).next_lane_type();
          }
          //计算到终点的距离
          const auto& end_position = call->end_position;
          if (end_position.has_lane_position()) {
            localSegments[size].distanceToEnd = (float)end_position.lane_position().s();
          }
          else if (end_position.has_aoi_position()) {
            int32_t end_aoi_id = end_position.aoi_position().aoi_id();
            HostAoi aoi = hostAoi[AoiID2Index[(int)end_aoi_id]];//通过aoi_id，获取aoi
            int last_lane_id = (int)localSegments[size].laneIndex;
            localSegments[size].distanceToEnd = (float)(aoi).lanS[last_lane_id]; //aoi连接的车道id → 对应道路上位置s
          } else {
            assert(false);
          }
          for (int i = size - 1; i > 0; --i) {
            switch ((simulet::PbNextLaneType)localSegments[i].nextLaneType) {
              case simulet::PbNextLaneType::NEXT_LANE_TYPE_FORWARD:
                // 直行的，累加本车道长度
                localSegments[i].distanceToEnd = 
                    localSegments[i+1].distanceToEnd + hostLane[laneId2Index[localSegments[i].laneIndex]].length;
                break;
              case simulet::PbNextLaneType::NEXT_LANE_TYPE_LEFT:
              case simulet::PbNextLaneType::NEXT_LANE_TYPE_RIGHT:
                // 如果变道，则不累加
                localSegments[i].distanceToEnd = localSegments[i+1].distanceToEnd;
                break;
              default:
                // throw std::runtime_error("vehicle: wrong PbNextLaneType");
                goto outLOOP;
            }
          }

          // 在 hostSegmentsToDevice 分配空间
          int info_index;
          int segment_index;
          {
            std::unique_lock<std::mutex> lock(m_mtx);
            successRouteNum++;  // 每成功一个导航+1
            info_index = infoIndex++;  // 获取当前车辆存放routeInfo的位置，routeInfo加1
            segment_index = segmentIndex;
            segmentIndex += size + 1;
            // std::cout<<std::endl<<"成功的导航数："<<successRouteNum<<"\t分配空间--> infoIndex: "<<info_index<<"\tsegmentIndex: "<<segment_index<<std::endl;
          }
          const auto& startAoiID = call->startAoiId;
          HostAoi startAoi = hostAoi[AoiID2Index[(int)startAoiID]];
          //设置 routeInfos 和 hostSegmentsToDevice
          routeInfos[info_index].vehicleIndex = call->vehicleIndex;
          routeInfos[info_index].startIndex = segment_index;
          routeInfos[info_index].routeLength = size; 
          routeInfos[info_index].s =hostAoi[AoiID2Index[(int)startAoiID]].lanS[localSegments[1].laneIndex];
          memcpy(hostSegmentsToDevice+segment_index, localSegments, (size+1)*sizeof(HostSegment));

          // // 打印 routeInfos 和 hostSegmentsToDevice信息
          // std::cout<<"routeInfos["<<info_index<<"]-->\tvehicleIndex: "<<routeInfos[info_index].vehicleIndex <<" \tstartIndex: "<<routeInfos[info_index].startIndex <<" \trouteLength: "<<routeInfos[info_index].routeLength << std::endl;  
          // int endIndex = segment_index + size;
          // for(int i = segment_index ; i <= endIndex ; ++i){
          //   if(i == segment_index){
          //     std::cout<<"hostSegmentsToDevice["<<i<<"]-->\tsize: "<<hostSegmentsToDevice[i].laneIndex << "(复用一个hostSegmentsToDevice[i])" <<std::endl;  
          //   }
          //   else{
          //       std::cout<<"hostSegmentsToDevice["<<i<<"]-->\tlaneIndex: "<<hostSegmentsToDevice[i].laneIndex <<" \tnextLaneType: "<<hostSegmentsToDevice[i].nextLaneType  <<" \tdistanceToEnd: "<<hostSegmentsToDevice[i].distanceToEnd << std::endl ;  
          //   }  
          // }
        }
      }
      else{
        assert(false);
      } 
    } else {
      // ::spdlog::error("routing_client: rpc failed: {}",
      //                 call->status.error_message());
      throw std::runtime_error("routing_client: rpc failed");
    }
    outLOOP:
    std::unique_lock<std::mutex> lock(mtx_);
    --count_waiting_req_;
    // std::cout<<"--count_waiting_req_:\t"<<count_waiting_req_<<std::endl;
    // 唤醒等待程序，检测请求是否处理完成
    cv_.notify_one();
    delete call;
  }
}

void RoutingClient::Shutdown() {
  // 等待已有的导航请求处理结束
  Wait();
  cq_.Shutdown();
  // 等待线程结束
  worker_.join();
}

void RoutingClient::Wait() {
  std::unique_lock<std::mutex> lock(mtx_);
  // 等待路径规划处理完成
  cv_.wait(lock, [this]() { return count_waiting_req_ == 0; });
}

std::shared_ptr<RoutingClient> routing_client;

