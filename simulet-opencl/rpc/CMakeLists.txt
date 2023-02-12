find_package(gRPC CONFIG REQUIRED)

# 生成链接库
add_library(rpc
  routing_client.h
  routing_client.cc
)

# 添加链接库
target_link_libraries(rpc
  # PUBLIC generated_proto
  # PUBLIC CONAN_PKG::grpc
  # PRIVATE CONAN_PKG::spdlog
  wolong
  gRPC::grpc++
)

