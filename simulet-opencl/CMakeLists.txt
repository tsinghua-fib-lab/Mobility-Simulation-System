cmake_minimum_required( VERSION 3.16 )
project( cmakeTest )
find_package( protobuf CONFIG REQUIRED )
find_package( gRPC CONFIG REQUIRED )

set( CMAKE_EXPORT_COMPILE_COMMANDS ON )

include_directories( ./utils )
include_directories( ./programs )

add_subdirectory( wolong )
add_subdirectory( entity )
add_subdirectory( utils )
add_subdirectory( rpc )

add_executable( demo simulator.cc simulet.cc ) 
SET(EXECUTABLE_OUTPUT_PATH ${PROJECT_SOURCE_DIR})

target_link_libraries( demo
wolong
OpenCL
protobuf::libprotobuf
rpc
entity
utils
)