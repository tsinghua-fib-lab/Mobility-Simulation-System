#include "./programs/header_global.h"

/**
    该内核用于实现vehicle在Lane上的主计算流
 */
__kernel
void vehicle_lane_update(__global Points *points, __global int *lane2AoiAto, __global struct drivingAction *action,
__global struct localController *l,__global int *vehicle2LaneAto,__global int *vehicle2LaneremoveAto) {
          size_t tid = get_global_id(0); // 确定tid
    int vehicleInAoi = points->vehicleInAoi;
    int rtid = points->vehicleIndex[vehicleInAoi + tid];
	//struct drivingAction *action;
    //points->vehicles[rtid];
    //printf("111");
    // TODO: 这里后期需要进行更新
    // if (false) {  // 如果车辆确定需要返回到Aoi中
    //     int ret = atomic_inc(lane2AoiAto);
    //     points->lane2Aoi[ret] = vehicleInAoi + tid;
    // }
    struct laneChange laneChange;
	laneChange.laneAhead.ptr = 0;
	laneChange.targetLane.ptr = 0;
	l->vehicle = &points->vehicles[rtid];
	l->vehicle->snapshot.speed =3;
	//l->vehicle->snapshot.
	l->vehicle->runtime.speed =3;
//这部分是获取环境getDrivingEnv

    struct DeviceLane *lanePtr = 0;//下一车道
	struct DeviceLane*lanePtrShadow =0;
    struct DeviceVehicle *vehiclePtr=0;//下一辆车
	struct DeviceVehicle *VehiclePtrShadow=0;

	struct DeviceVehicle *beforeV = 0;
	struct DeviceVehicle *afterV = 0;
	float beforeVRelativeDistance=0;
	float afterVRelativeDistance=0;

    float relativeDistance=0;//俩车间距
	float relativeDistanceShadow=0;
	float laneRelativeDistance=0;//路相对距离
	float laneRelativeDistanceShadow=0;
	//currentLaneMaxSpeed不需要赋值
//getAgentAndLaneAhead
	struct VehicleRoute *route = &points->vehicles[rtid].vehicleRoute;
	struct DeviceLane *lane = points->vehicles[rtid].snapshot.lane;
	float s = points->vehicles[rtid].snapshot.s;
		//获取给定s后的第一辆车GetFirstVehicleAfterS
	bool hasTargetLane =1;
	bool isChangeLane = 1;

		struct DeviceVehicle *aheadV = points->vehicles[rtid].relation[1][1];
		float aheadS;
		if (aheadV == 0){
			aheadS =0;
		}
		else{
			aheadS = aheadV->runtime.s;
		}
		//printf("test1\n");
		//检查前方车道
		if(points->vehicles[rtid].vehicleRoute.segment!=0){
			if(points->vehicles[rtid].vehicleRoute.segment->nextLaneType == 1){//NextLaneType_NEXT_LANE_TYPE_FORWARD NextLaneType = 1
				//printf("test2\n");
        		lanePtr = points->vehicles[rtid].vehicleRoute.segment->nextSegment->lane;
				//printf("test8\n");
				float laneRelativeDistance = points->vehicles[rtid].runtime.lane->length -points->vehicles[rtid].runtime.s;
				//printf("test9\n");
				//printf("lanePtrID:%d\n",lanePtr->index);
				//printf("laneRelativeDistance:%f\n",laneRelativeDistance);
        		if(!aheadV){
					//printf("test3\n");
            		float minDis = -1;//无穷大改为-1
            	// 遍历后继车道找出后车
           	 for(int i = 0; i < points->vehicles[rtid].runtime.lane->next_lane_size; ++i){ 
					//printf("test4\n");
                	struct DeviceVehicle *v = points->vehicles[rtid].runtime.lane->next_lanes[i]->first->relation[1][1];//lane上的头车
                	float vs; 
					if(v==0){
						vs = -1;
					}
					else{
						vs = v->runtime.s;
					}

                	if(vs > minDis){
                    	aheadV = v;
                    	aheadS = vs;
                	}
            		}
            		aheadS += points->vehicles[rtid].runtime.lane->length;
					//printf("test5\n");
        		}	
    		}
		}
        if(aheadV !=0){
			//printf("test6\n");
            vehiclePtr = aheadV;
            relativeDistance = aheadS - s -aheadV->attribute.length;
			//printf("vehiclePtr->runtime.s:%f\n",vehiclePtr->runtime.s);
			//printf("relativeDistance:%f\n",relativeDistance);
        }
		//getAgentAndLaneAhead到此结束，获取了Vehicleptr和relativeDistance
		// if (points->vehicles[rtid].snapshot.LaneChangeStatus!=0){//是否变道
		// 	if (points->vehicles[rtid].runtime.lane->isInJunction){//是否在路口
		// 		printf("lane change in junction");
        //         return;
		// 	}
		// }
	
		if (points->vehicles[rtid].runtime.LaneChangeStatus == 2){
			//GetAgentAndLaneAheadShadow还没变道之前，需要看两个路上的车来判断自己速度
			//获取上一时刻车辆（非影子）所在的车道（需要预先检查车辆是否在车道内）
			struct DeviceLane* Shadowlane=0;
			struct DeviceVehicle *aheadV=0;
            struct VehicleRoute route = points->vehicles[rtid].vehicleRoute;
            float s = points->vehicles[rtid].snapshot.s;
			//printf("test\n");
			if (route.segment->nextSegment->lane!=0) {
				Shadowlane = points->vehicles[rtid].snapshot.lane;//是snapshot的lane
				//printf("test2\n"); 
				
				if(points->vehicles[rtid].snapshot.lane->first->relation[1][1] != NULL){
				//printf("test3\n"); 
				aheadV =Shadowlane->first->relation[1][1]; // lane 的头车
				float s1 = aheadV->runtime.s*Shadowlane->length/(points->vehicles[rtid].runtime.lane->length);
				//printf("test33\n"); 
				while(s > s1){
					if( aheadV->relation[1][1]==0){
						break;
					}else{
						aheadV = aheadV->relation[1][1];
					}	
					s1 = aheadV->runtime.s*Shadowlane->length/(points->vehicles[rtid].runtime.lane->length);
				}
				//在投影车道上的坐标	
			}
			}
			// printf("test1\n"); 
			// //获取上一时刻车辆（非影子）所在的车道（需要预先检查车辆是否在车道内）
			// if (Shadowlane) {
			// 	s = points->vehicles[rtid].snapshot.s;
			// }
			
			//float v1Tov2S = v1->runtime.s*p2->length/Shadowlane->length;
			// 检查当前车道
			// 获取给定s坐标后的第一辆车
			
			//printf("test4\n"); 
			// 检查前面的车道（shadow在index+1，向前是index+2）
			if(aheadV==0 && route.segment->nextLaneType == 1){
				if( route.segment->nextSegment->nextSegment!=0){
					struct DeviceLane* nextLane = route.segment->nextSegment->nextSegment->lane; 
					aheadV = nextLane->first->relation[1][1];
					lanePtrShadow = nextLane;//env->laneAheadShadow.
					//printf("lanePtrShadow0:%f\n",lanePtrShadow->max_speed);
					laneRelativeDistanceShadow = lane->length - s;//env->laneAheadShadow.
				}
				
					
					//printf("lanePtrShadowLaneID:%d\n",lanePtrShadow->index);
					//printf("laneRelativeDistanceShadow:%f",laneRelativeDistanceShadow);
			}
			if(aheadV){
				VehiclePtrShadow = aheadV;//env->agentAheadShadow.
				relativeDistanceShadow = aheadV->snapshot.s - s -aheadV->attribute.length;//env->agentAheadShadow.
				//printf("VehiclePtrShadow:%f",VehiclePtrShadow->runtime.s);
				//printf("relativeDistanceShadow:%f",relativeDistanceShadow);
			}

		}
	//到此GetAgentAndLaneAheadShadow结束，传出了lanePtrShadow,laneRelativeDistanceShadow,VehiclePtrShadow,relativeDistanceShadow
	//getLaneChangeEnv开始
	
	else if(!points->vehicles[rtid].runtime.lane->isInJunction){
		*route = points->vehicles[rtid].vehicleRoute;
		struct Segment *segment = route->segment;
		if(!(segment->nextLaneType == 2//NextLaneType_NEXT_LANE_TYPE_LEFT NextLaneType = 2
				|| segment->nextLaneType == 3)){//NextLaneType_NEXT_LANE_TYPE_RIGHT NextLaneType = 3
			isChangeLane = 0; //无需变道
		}
		if(isChangeLane){

			// // //env->laneChange.enable = true;
			//points->vehicles[rtid].runtime.to_left = true;
			// //获取上一时刻车辆（非影子）所在的车道（需要预先检查车辆是否在车道内）
			struct DeviceLane* lane;
			lane = points->vehicles[rtid].snapshot.lane;//上一时刻的lane
			struct DeviceLane* targetLane;
			// //TODO:这里不确定是否是route.Route[route.Index+1].Lane
			if (points->vehicles[rtid].vehicleRoute.segment->nextSegment!=0)
				targetLane = points->vehicles[rtid].vehicleRoute.segment->nextSegment->lane;
			float s = points->vehicles[rtid].snapshot.s;//上一时刻的s
			float neighborS=0;//执行ProjectFromLane
			//对同一道路内的车道按比例“投影”
			if(targetLane==0){
				hasTargetLane =0;
			}
			else{//
				if(targetLane->length < (s/(lane->length)*(targetLane->length)))
				{
					neighborS = targetLane->length;
				}
				//math.Min(l.length, math.Max(.0, otherS/other.length*l.length))
				else
				{	
					neighborS = s/(lane->length)*(targetLane->length);
				}
			}
	// 		//printf("hasTargetLane:%d\n",hasTargetLane);
			if(hasTargetLane==1){
				// 变道目标
				//laneChange->targetLane.ptr = targetLane;
				laneChange.targetLane.ptr = targetLane;
				//printf("laneChange.targetLane.ptrIndex:%d\n",laneChange.targetLane.ptr->index);
				if(points->vehicles[rtid].vehicleRoute.segment->nextSegment!=0){
					struct Segment * nextSegment=  points->vehicles[rtid].vehicleRoute.segment->nextSegment;	
				// 变道目标的后继
				if (nextSegment->nextLaneType == 1) {//NEXT_LANE_TYPE_FORWARD=1
					if(nextSegment->nextSegment!=0){
						//printf("test222\n");
						// struct DeviceLane *nextLane = nextSegment->nextSegment->lane;//nextLane := route.Route[route.Index+2].Lane
						laneChange.laneAhead.ptr = nextSegment->nextSegment->lane;
						//printf("test222\n");
						laneChange.laneAhead.relativeDistance = laneChange.targetLane.ptr->length - neighborS;
						// printf("test333\n");
						//printf("laneChange.laneAhead.ptrIndex:%d\n",laneChange.laneAhead.ptr->index);
						//printf("test444\n");
						//printf("laneChange.laneAhead.relativeDistance:%f\n",laneChange.laneAhead.relativeDistance);	
						//printf("test555\n");
					}
				}
				if (nextSegment->nextLaneType == 2){//向左变道
					//printf("test1\n");
					beforeV = points->vehicles[rtid].relation[0][0];//当前车尾后面的车
					afterV = points->vehicles[rtid].relation[1][0];//当前车头前面的车
					//printf("test1.1\n");
					//printf("beforeV:%f\n",beforeV->attribute.length);
					if(beforeV!=0){
						//printf("test111");
						if(neighborS-beforeV->snapshot.s>0){
							beforeVRelativeDistance = neighborS-beforeV->snapshot.s;
						}else{
							beforeVRelativeDistance = 0;
						}	
						//printf("beforeVLength:%f\n",beforeV->attribute.length);	
						//printf("beforeVRelativeDistance:%f\n",beforeVRelativeDistance);		
					}
					if(afterV!=0){
						//printf("test112");
						if(afterV->snapshot.s-neighborS-afterV->attribute.length>0){
							afterVRelativeDistance = afterV->snapshot.s-neighborS-afterV->attribute.length;
						}else{
							afterVRelativeDistance = 0;
						}
						//printf("afterVLength:%f\n",afterV->attribute.length);	
						//printf("afterVRelativeDistance:%f\n",afterVRelativeDistance);		
					}
					
				}
				if (nextSegment->nextLaneType == 3){//向右变道
					//printf("test2\n");
					if(points->vehicles[rtid].relation[0][2]!=0)
						beforeV =points->vehicles[rtid].relation[0][2];
					if(points->vehicles[rtid].relation[1][2]!=0)	
						afterV = points->vehicles[rtid].relation[1][2];
					//printf("test2.2\n");
					if(beforeV!=0){
						//printf("test221");
						if(neighborS-beforeV->snapshot.s>0){
							beforeVRelativeDistance = neighborS-beforeV->snapshot.s;
						}else{
							beforeVRelativeDistance = 0;
						}
						//printf("beforeVLength:%f\n",beforeV->attribute.length);	
						//printf("beforeVRelativeDistance:%f\n",beforeVRelativeDistance);	
					}
					if(afterV!=0){
						//printf("test222");
						if(afterV->snapshot.s-neighborS-afterV->attribute.length>0){
							afterVRelativeDistance = afterV->snapshot.s-neighborS-afterV->attribute.length;
						}else{
							afterVRelativeDistance = 0;
						}
						//printf("afterVLength:%f\n",afterV->attribute.length);	
						//printf("afterVRelativeDistance:%f\n",afterVRelativeDistance);		
					}
				}
				// if(beforeV!=0)
				// 	printf("beforeVLength:%f\n",beforeV->attribute.length);	
				// if(afterV!=0)
				// 	printf("afterVLength:%f\n",afterV->attribute.length);	
				// printf("beforeVRelativeDistance:%f\n",beforeVRelativeDistance);	
				// printf("afterVRelativeDistance:%f\n",afterVRelativeDistance);	
			
			}
		}
	}
}
//printf("111");
	// // 自主变道
	points->vehicles[rtid].runtime.DistanceToEnd = points->vehicles[rtid].vehicleRoute.segment->distanceToEnd;
	float distanceToEnd = points->vehicles[rtid].runtime.DistanceToEnd;
	//getDrivingEnv结束
			

	//这部分是planLaneChange
	
	// 	//都是return的按钮开关
	int button = 1;
	


    // struct laneChange lc = env->laneChange;
	// struct laneChange *lc;
    float uba;
    
	

	// 	//if l.vehicle.InJunction()路口不变道
	if(button){
		if (l->vehicle->snapshot.lane->isInJunction){
		//return;
		button =0;
		}
	}
 	//对应if l.vehicle.IsLaneChanging()
	if(button){

		if(l->vehicle->snapshot.LaneChangeStatus!=0){
			// printf("test20\n");
			if (l->vehicle->runtime.LaneChangeStatus== 1){ //车的状态如果为准备变道状态
			//printf("test21\n");
				
				 //将环境进行赋值
				//开始setLaneChange
				//*lc = laneChange;
				//l->vehicle->runtime.LaneChangeStatus = 1;
				uba = l->vehicle->attribute.usual_braking_acc;
				//计算在distance距离以加速度a减速为v所需的初速度ComputeSpeed
				//判断现在直接变道是否会撞前车或者让后车追尾
				if ((afterV != 0 && (afterVRelativeDistance <= 0 || 
					sqrt(afterV->snapshot.speed * afterV->snapshot.speed -2 * uba * afterVRelativeDistance) <
					l->vehicle->snapshot.speed)) ||
					(beforeV != 0 && (beforeVRelativeDistance <= 0 ||
					sqrt(l->vehicle->snapshot.speed * l->vehicle->snapshot.speed -2 * beforeV->attribute.usual_braking_acc *beforeVRelativeDistance ) < 
					beforeV->snapshot.speed)) ){
					// 设置is_intentional并减速(直到停车)等待路况合适
					l->vehicle->runtime.LaneChangeStatus = 1;
					//printf("test22\n");
					if(action[rtid].acceleration>uba){
						action[rtid].acceleration = uba;
					}
				}
				else{
					l->vehicle->runtime.LaneChangeStatus = 2;
					// 设置shadow
					if(afterV!=0){
					VehiclePtrShadow = afterV;//env->agentAheadShadow.
					//对应367，无变化
					}
					relativeDistanceShadow = afterVRelativeDistance;
					if(laneChange.laneAhead.ptr!=0){

					
					lanePtrShadow = laneChange.laneAhead.ptr;
					//printf("lanePtrShadow1:%f\n",lanePtrShadow->max_speed);
					laneRelativeDistanceShadow = laneChange.laneAhead.relativeDistance;
					}
					
					//env.agentAheadShadow = env.laneChange.agentAhead
					//env.laneAheadShadow = env.laneChange.laneAhead
				}
	
				float value1;
				if(l->vehicle->runtime.speed*3<l->vehicle->attribute.length){//LANE_CHANGE_LENGTH_FACTOR=3
					value1 = l->vehicle->attribute.length;
				}else if(l->vehicle->runtime.speed*3>l->vehicle->runtime.lane_change_total_length){
					value1 = l->vehicle->runtime.lane_change_total_length;
				}else{
					value1 = l->vehicle->runtime.speed*3;
				}
				action[rtid].laneChangeLengh = value1;
				
			}
			button = 0;
		}
	}
	int cnt = 0;
	float distance = l->vehicle->snapshot.speed * l->vehicle->snapshot.speed / -uba / 2 +
					 l->vehicle->snapshot.speed * stepInterval + (float)cnt * l->vehicle->attribute.lane_change_length +
					 l->vehicle->attribute.min_gap;
	//printf("distance:%f",distance);
	float remainingDistance = l->vehicle->snapshot.lane->length - l->vehicle->snapshot.s;
	// lc = &laneChange;
	// laneChange.agentAhead.ptr = afterV;
	// laneChange.agentAhead.relativeDistance = afterVRelativeDistance;
	// laneChange.agentBehind.ptr=beforeV;
	// laneChange.agentBehind.relativeDistance = beforeVRelativeDistance;
	uba =l->vehicle->attribute.usual_braking_acc;
	if(button){
		//printf("test111\n");
		// 下面的循环不会越界，因为最后一个lane的类型必不为变道
		struct Segment *IsLaneChange;
		IsLaneChange =l->vehicle->vehicleRoute.segment;
		//printf("Segment越界1\n");
		while(IsLaneChange->nextLaneType==2||IsLaneChange->nextLaneType==3){
			IsLaneChange = IsLaneChange->nextSegment;
			//printf("Segment越界2\n");
			cnt++;
		}
	}
	if(button){
	
 		if (cnt >0){
			//printf("test666\n");
			// // 车道前后5m不变道
			if(l->vehicle->snapshot.s<5||remainingDistance<5){
			//return;
				//printf("test777\n");
				button=0;//设置开关
			}
		}
		if(button){
			//根据车道尽头的距离调整触发变道的概率
			//printf("remainingDistance:%f,distance%f\n",remainingDistance,distance);
			if(remainingDistance <=distance// TODO:少一个变道概率
				//||l.vehicle.PTrue(remainingDistance/l.vehicle.GetLane().Length())
				){//执行一次setLaneChange
				//计算在distance距离以加速度a减速为v所需的初速度ComputeSpeed
				// 判断现在直接变道是否会撞前车或者让后车追尾
				if ((afterV != 0 && (afterVRelativeDistance <= 0 || 
					sqrt(afterV->snapshot.speed * afterV->snapshot.speed -2 * uba * afterVRelativeDistance) <
					l->vehicle->snapshot.speed)) ||
					(beforeV != 0 && (beforeVRelativeDistance <= 0 ||
					sqrt(l->vehicle->snapshot.speed * l->vehicle->snapshot.speed -2 * beforeV->attribute.usual_braking_acc *beforeVRelativeDistance ) < 
					beforeV->snapshot.speed)) ){
		
					// 设置is_intentional并减速(直到停车)等待路况合适
					l->vehicle->runtime.LaneChangeStatus = 1;
					if(action[rtid].acceleration>uba){
						//printf("test7\n");
						action[rtid].acceleration = uba;
						
					}
				}
				else{//删除掉正在变道状态，直接给runtime赋值
					//if(laneChange.laneAhead.ptr!=0)
						//printf("test123:\n");
					l->vehicle->runtime.LaneChangeStatus = 2;
					// // 设置shadow
					if(afterV!=0){
						VehiclePtrShadow = afterV;//env->agentAheadShadow.
						
					}
					relativeDistanceShadow = afterVRelativeDistance;//对应367，无变化
					if(laneChange.laneAhead.ptr!=0){
						lanePtrShadow = laneChange.laneAhead.ptr;
						//printf("test1111111:\n");
						//printf("lanePtrShadow2:%d\n",lanePtrShadow->index);
						laneRelativeDistanceShadow = laneChange.laneAhead.relativeDistance;
					}
					
					//l->vehicle->LaneChangeStatus = LaneChangeStatus_CHANGING;
					// 设置shadow
					// env.agentAheadShadow = env.laneChange.agentAhead
					// env.laneAheadShadow = env.laneChange.laneAhead
				}
				float value1;
				if(l->vehicle->snapshot.speed*3<l->vehicle->attribute.length){//LANE_CHANGE_LENGTH_FACTOR=3
					value1 = l->vehicle->attribute.length;
				}else if(l->vehicle->snapshot.speed*3>l->vehicle->snapshot.lane_change_total_length){
					value1 = l->vehicle->snapshot.lane_change_total_length;
				}else{
					value1 = l->vehicle->snapshot.speed*3;
				}
				action[rtid].laneChangeLengh = value1;
				//button = false;
			}//执行完毕setLaneChange
			button = 0;
		}	
	}   

	l->leftMotivation=0;
	l->rightMotivation = 0;
	//planLaneChange到此结束







// //planAcceleration开始



	float currentLaneMaxSpeed  = l->vehicle->snapshot.lane->max_speed;
	float currentSpeed = l->vehicle->snapshot.speed;
	float usualBrakingAcc = l->vehicle->attribute.usual_braking_acc;
	float maxBrakingAcc = l->vehicle->attribute.max_braking_acc;
	float usualAcc = l->vehicle->attribute.usual_acc;
	float maxAcc = l->vehicle->attribute.max_acc;
	float maxSpeed = min(l->vehicle->attribute.max_speed,currentLaneMaxSpeed);
//Clamp
	float acc = (maxSpeed-currentSpeed)/stepInterval;
	if((maxSpeed-currentSpeed)/stepInterval < usualBrakingAcc){
		acc = usualBrakingAcc;
	}else if((maxSpeed-currentSpeed)/stepInterval > usualAcc){
		acc = usualAcc;
	}
//printf("acc:%f\n",acc);
//这时acc为2
//前方有其他车辆
	if(vehiclePtr !=0||VehiclePtrShadow != 0 ){//|| laneChange.agentAheadShadow
	// Krauss模型：
		struct Agent agent1;
		struct Agent agent2;
		agent1.ptr = vehiclePtr;
		agent1.relativeDistance = relativeDistance;
		agent2.ptr = VehiclePtrShadow;
		agent2.relativeDistance = relativeDistanceShadow;
		struct Agent agent[2] = {agent1,agent2};
		struct Agent *ahead;
		for(int i=0;i<2;i++ ){
			ahead =&agent[i];
			if(ahead->ptr == 0) {
				continue;
			}
		// 记本车为B，前车为A，本时刻速度为v_B和v_A，那么本车下一时刻的速度u需要满足
		// u^2/(2a_B)+(v_B+u)t/2 ≤ (v_A^2)/(2a_A)+d
		
			float a = 0.5/-usualBrakingAcc;
			float b = 0.5*stepInterval;
			float c = 0.5*currentSpeed*stepInterval - ahead->ptr->snapshot.speed *ahead->ptr->snapshot.speed*0.5/
					-ahead->ptr->attribute.max_braking_acc +
					l->vehicle->attribute.min_gap - ahead->relativeDistance;
			float det = b*b - 4*a*c;
			//printf("%f\n",det);
			if(det<0){
			// 紧急刹车
				if(acc>maxBrakingAcc){
					acc = maxBrakingAcc;
				}
			
				break;
			}else{
				float targetSpeed =0.0;
				if((sqrt(det)-b)/a/2>0){
					targetSpeed =(sqrt(det)-b)/a/2;
					//printf("%f\n",targetSpeed);
				}
				float a =0.0;
				if(usualAcc<(targetSpeed-currentSpeed)/stepInterval){
					a = (targetSpeed-currentSpeed)/stepInterval;
				}else if(usualAcc>=(targetSpeed-currentSpeed)/stepInterval&&usualAcc<=usualBrakingAcc){
					a = usualAcc;
				}else{
					a = usualBrakingAcc;
				}
				if(acc >a){
					acc = a;
				}
			
			}
		
		}
	}
//
	currentSpeed = l->vehicle->snapshot.speed;
//printf("max1:%f",currentSpeed);
// 自己和影子前方的lane
//planLaneAheadAcceleration
	if(lanePtr!=0){
		if (currentSpeed > lanePtr->max_speed){
	// 超速需要减速
			float max1 =0.0;
			float max2 =0.0;
			if(laneRelativeDistance-currentSpeed*stepInterval>1e-4){
				max1 = laneRelativeDistance-currentSpeed*stepInterval;
			}else{
				max1 =1e-4;
			}
	
			if(max1>1e-6){
				max2 = max1;
			}else{
				max2 = 1e-6;
			}
	//printf("a111:%f\n",lanePtr->max_speed);
	
			float a = (lanePtr->max_speed * lanePtr->max_speed -currentSpeed * currentSpeed)/ 2 /max2;
	//printf("lanePtr->max_speed:%f,currentSpeed%f\n",lanePtr->max_speed,currentSpeed);
			if(a < l->vehicle->attribute.usual_braking_acc){
				if(acc>a){
					acc =a;
				}
		
			}
		}
	}
	// if(lanePtrShadow!=0){
	// 	printf("lanePtrShadow->max_speed:%f",lanePtrShadow->max_speed);
	// 	printf("test9\n");
	// }
	// acc =min(acc, l->planLaneAheadAcceleration(&env.laneAhead, stepInterval));
	// acc = math.Min(acc, l.planLaneAheadAcceleration(&env.laneAheadShadow, stepInterval))
	//printf("acc:%f",acc);
// 	printf("test10\n");
	if(lanePtrShadow!=0){
		printf("test11,%f\n");
		if (currentSpeed > lanePtrShadow->max_speed){
		// 超速需要减速
			float max1;
			float max2;
			printf("test12\n");
			if(laneRelativeDistanceShadow-currentSpeed*stepInterval>1e-4){
				max1 = laneRelativeDistanceShadow-currentSpeed*stepInterval;
			}else{
				max1 =1e-4;
			}
			printf("max1:%f",max1);
			if(max1>1e-6){
				max2 = max1;
			}else{
				max2 = 1e-6;
			}
			//printf("max2:%f\n",max2);
			if(lanePtrShadow!=0){
				float a = (lanePtrShadow->max_speed * lanePtrShadow->max_speed -currentSpeed * currentSpeed)/ 2 /max2;
				if(a < l->vehicle->attribute.usual_braking_acc){
					if(acc>a){
						acc =a;
					}
		
				}
			}
		//printf("acc:%f\n",acc);
		}
	}
//printf("acc:%f",acc);
//终点
	if(currentSpeed *stepInterval+currentSpeed*currentSpeed /usualBrakingAcc/2 >= distanceToEnd){
		acc = min(acc, -currentSpeed*currentSpeed/distanceToEnd/2);
	}
//printf("acc:%f",acc);
// //lo.Clamp(acc, maxBrakingAcc, maxAcc)


	float acc1=0.0;
	if(acc<=maxBrakingAcc){
	acc1 = maxBrakingAcc;
	//printf("acc1:%f\n",acc1);
	}else if(acc>=maxAcc){
		acc1 = maxAcc;
	}else{
		acc1 = acc;
	}
	// if(action[rtid].acceleration>acc1){
	// 	action[rtid].acceleration = acc1;
	// }else{
	// 	action[rtid].acceleration = action[rtid].acceleration;
	// }
	
	//printf("acc1:%f\n",acc1);
// //和drivingEnv比较,会报错
//printf("test1111111111\n");

	// if(action[rtid].acceleration>acc1){
	// 		printf("test222\n");
	// 	action[rtid].acceleration =acc1;
	// 		printf("test333\n");
		
	// }
	//printf("test444\n");
//printf("action[rtid].acceleration:%f\n",action[rtid].acceleration);



//refreshMotionByDrivingAction
float speed =0.0;
float ds =0.0;
//计算本时刻的速度与移动距离 v(t)=v(t-1)+acc*dt, ds=v(t-1)*dt+acc*dt*dt/2
//computeSpeedAndDistance
int button2 =1;
if(button2){
	float dv = acc1 *stepInterval;
	if (l->vehicle->snapshot.speed+dv < 0) {
		// 刹车到停止
		speed =  0,ds= l->vehicle->snapshot.speed * l->vehicle->snapshot.speed / 2 / -action->acceleration;
		button2 =0;
	}
	if(button2){	
	if(l->vehicle->snapshot.speed+dv > 0){
		speed =l->vehicle->snapshot.speed + dv,ds = (l->vehicle->snapshot.speed + dv/2) * stepInterval;
		l->vehicle->runtime.speed = speed;
		// printf("ds:%f",ds);
		// printf("speed:%f",speed);
	}
}
}
//printf("speed:%f\n",speed);
//printf("ds:%f\n",ds);
// if (l->vehicle->runtime.LaneChangeStatus == 0) {//不变道情况
// 		// 直行
// 		//driveStraightAndRefreshLocation
// 		l->vehicle->runtime.s = l->vehicle->runtime.s+ds;
		
// 		if (l->vehicle->runtime.s >l->vehicle->runtime.lane->length){
			
// 			l->vehicle->runtime.LaneChangeStatus = 0;//更新为不变道
// 			while(l->vehicle->runtime.s>l->vehicle->runtime.lane->length){
// 				struct Segment s = *l->vehicle->vehicleRoute.segment;
// 				l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
// 				l->vehicle->vehicleRoute.size--;
// 				bool teleport = false;
// 				int nextLaneType1;
// 				if(l->vehicle->vehicleRoute.size>=0){
// 					nextLaneType1 =  s.nextLaneType;
// 				}
// 				while(nextLaneType1 ==2||nextLaneType1==3){
// 					l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
// 					l->vehicle->vehicleRoute.size--;
// 					if(l->vehicle->vehicleRoute.size==0){

// 					}
// 				}
// 			}
		
// 		}

//  	}else if(l->vehicle->runtime.LaneChangeStatus == 1 ){
// 	// 	struct DeviceLane *targetLane = l->vehicle->vehicleRoute.segment->nextSegment->lane;
// 	// 	int myIndex = atomic_inc(vehicle2LaneAto);
// 	// 	printf("111\n");
// 	// 	targetLane->insert_buffer[myIndex];
// 	// 	int move = atomic_inc(vehicle2LaneremoveAto);
// 	// 	l->vehicle->runtime.lane->insert_buffer[move];

// 	// 	float neiborS;
// 	// 	float targetS;
// 	// 	neiborS = l->vehicle->runtime.s * targetLane->length/l->vehicle->runtime.lane->length;
// 	// 	if(neiborS+action->laneChangeLengh<targetLane->length){
// 	// 		targetS = neiborS+action->laneChangeLengh;
// 	// 	}else{
// 	// 		targetS = targetLane->length;
// 	// 	}
// 	// 	if (neiborS+ds >= targetS){
// 	// 		l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
// 	//  		l->vehicle->vehicleRoute.size--;
// 	// 		//driveStraightAndRefreshLocation
// 	// 		neiborS = neiborS+ds;
// 	// 		if (neiborS >targetLane->length){
				
// 	// 			l->vehicle->runtime.LaneChangeStatus = 0;//更新为不变道
// 	// 			lanePtrShadow=0;
// 	// 			while (neiborS > targetLane->length) {
// 	// 				struct Segment s = *l->vehicle->vehicleRoute.segment;
// 	// 				l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
// 	//  				l->vehicle->vehicleRoute.size--;
// 	// 				bool teleport = false;
// 	// 				int nextLaneType1;
// 	// 				if(l->vehicle->vehicleRoute.segment->bitIndex < l->vehicle->vehicleRoute.size){
// 	// 					nextLaneType1 = s.nextLaneType;//有问题
// 	// 					while(nextLaneType1) {
// 	// 						s = *l->vehicle->vehicleRoute.segment;
// 	// 						l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
							
// 	//  						l->vehicle->vehicleRoute.size--;
// 	// 					if (l->vehicle->vehicleRoute.segment->bitIndex >= l->vehicle->vehicleRoute.size ){
// 	// 						break;
// 	// 					}
// 	// 					nextLaneType1 = s.nextLaneType;
// 	// 					teleport = true;
// 	// 					}
// 	// 				}
// 	// 				// if(l->vehicle->vehicleRoute.segment->bitIndex > l->vehicle->vehicleRoute.size){
// 	// 				// 	nextLaneType1 = s.nextLaneType;//有问题

// 	// 				// 	targetLane = nextLaneType1
// 	// 				// }
// 	// 				struct DeviceLane *nextLane = l->vehicle->vehicleRoute.segment->lane;
// 	// 				ds = nextLane->length;
// 	// 				targetLane = nextLane;
// 	// 		}
			
			
// 	// 	}
// 	// 	l->vehicle->runtime.lane = targetLane;
// 	// 	l->vehicle->runtime.s = ds;
// 	// 	lanePtrShadow =0;
// 	// 	l->vehicle->runtime.LaneChangeStatus = 0;
// 	// }else{

// 	// }


struct DeviceLane *targetLane = l->vehicle->vehicleRoute.segment->nextSegment->lane;
float neiborS = l->vehicle->runtime.s * targetLane->length/l->vehicle->runtime.lane->length;
//printf("neiborS:%f\n",neiborS);
//targetLane := v.route.Route[v.route.Index+1].Lane
float targetS = min(neiborS+action->laneChangeLengh, targetLane->length);
//printf("targetS:%f\n",targetS);

	if (neiborS+ds >= targetS) {
			// 跳过变道
			// if(l->vehicle->vehicleRoute.segment->nextSegment!=0){
			// 	l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
			// }
			//v.driveStraightAndRefreshLocation(targetLane, neiborS, ds)
			neiborS +=ds;
			if(neiborS>targetLane->length){
				// int myIndex = atomic_inc(vehicle2LaneAto);
				// printf("111\n");
				// targetLane->insert_buffer[myIndex];
				// int move = atomic_inc(vehicle2LaneremoveAto);
				// l->vehicle->runtime.lane->insert_buffer[move];
				l->vehicle->runtime.LaneChangeStatus =0;
				while(neiborS>targetLane->length){
					//
					//if v.route.Index < v.route.Size
					if(l->vehicle->vehicleRoute.segment->bitIndex<l->vehicle->vehicleRoute.size){
						//laneNextType := v.route.Route[v.route.Index-1].NextLaneType
						int laneNextType = l->vehicle->vehicleRoute.segment->nextLaneType;
						l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
						while(laneNextType==2||laneNextType==3){
							l->vehicle->vehicleRoute.segment = l->vehicle->vehicleRoute.segment->nextSegment;
							if(l->vehicle->vehicleRoute.segment->nextSegment==0){
								break;
							}
							laneNextType = l->vehicle->vehicleRoute.segment->nextLaneType;
						}
					}
					// if(l->vehicle->vehicleRoute.segment->bitIndex>l->vehicle->vehicleRoute.size){

					// }
					struct DeviceLane *nextLane = l->vehicle->vehicleRoute.segment->lane;
					neiborS -= targetLane->length;
					targetLane = nextLane;
				}
				

			}
			l->vehicle->runtime.s = neiborS;
			l->vehicle->runtime.lane = targetLane;
		}else{
			// motion.ShadowLane = targetLane
			// motion.ShadowS = neiborS + ds
			// motion.S = motion.Lane.ProjectFromLane(targetLane, neiborS+ds)
			// motion.LaneChangeTotalLength = targetS - neiborS
			// motion.LaneChangeCompletedLength = ds

		}


// }
l->vehicle->snapshot = l->vehicle->runtime;
//更新到终点的距离
if(l->vehicle->vehicleRoute.segment->distanceToEnd-l->vehicle->runtime.s>0){
	l->vehicle->runtime.DistanceToEnd =l->vehicle->vehicleRoute.segment->distanceToEnd-l->vehicle->runtime.s;
}else{
	l->vehicle->runtime.DistanceToEnd = 0;
}



//checkCloseToEndAndRefreshRuntime
// bool reachTarget;
// if(l->vehicle->runtime.DistanceToEnd > 1){
// 	reachTarget = false;
// 	button = 0;
// }
if(button){
if(l->vehicle->vehicleRoute.endAoi!=0){
	l->vehicle->runtime.lane = l->vehicle->vehicleRoute.segment->lane;
	l->vehicle->runtime.s = distanceToEnd;
	// v.runtime.Motion.S = v.route.EndS;
	// v.runtime.Motion.ShadowLane = nil
	// v.runtime.Motion.ShadowS = 0
	// v.runtime.Motion.LaneChangeStatus = LaneChangeStatus_NONE
}
}
// if(l->vehicle->runtime.lane!=0){
// 	printf("rtid:%d\n",rtid);//13715,32108,778506,45637
// }
// if(rtid==13715){
// 	printf("speed:%f\n",l->vehicle->runtime.speed);
// 	printf("s:%f\n",l->vehicle->runtime.s);
// 	printf("laneindex:%d\n",l->vehicle->runtime.lane->index);
// 	if(l->vehicle->relation[1][1]!=0){
// 		printf("afterS:%f\n",l->vehicle->relation[1][1]->runtime.s);
// 	}
// }
// 车辆直接根据更新后的位置信息插入新的lane
*vehicle2LaneAto = 0;
*vehicle2LaneremoveAto =0;
}
