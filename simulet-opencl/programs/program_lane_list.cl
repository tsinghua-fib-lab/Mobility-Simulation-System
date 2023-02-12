#include "./programs/header_vehicle.h"

__kernel 
  void vehicle_list(__global struct DeviceLane *lane, __global struct DeviceVehicle *vehicles,__global int *InsertNum){
    size_t tid = get_global_id(0);
    lane[tid].insert_num=InsertNum[tid];
    
    int button= 0;//用于判断lane上是否有车，是否执行链表的构建。
    lane[tid].end->relation[1][1]=0;
    // insert_buffer排序
    if (lane[tid].insert_buffer[0]!=0){
            for(int i=0;i<lane[tid].insert_num-1;i++){
                for(int j=0;j<lane[tid].insert_num-i-1;j++){
                    if(lane[tid].insert_buffer[j+1]->runtime.s < lane[tid].insert_buffer[j]->runtime.s){
                       struct DeviceVehicle *temp = lane[tid].insert_buffer[j];
                        lane[tid].insert_buffer[j] = lane[tid].insert_buffer[j+1];
                        lane[tid].insert_buffer[j+1] = temp;
                    }
                }
            }
    }
    //如果lane上没车
    if (lane[tid].first->relation[1][1] == 0) {
        lane[tid].num_vehicles=0;
        
    }

            
//lane上有车
//else{
    
    //车辆从“主链”上摘除及全lane同步
        struct DeviceVehicle *vehicle_remove=0;
        //remove操作
        if(lane[tid].remove_buffer[0]!=0){
            vehicle_remove = lane[tid].remove_buffer[0];
            //printf("vehicle_remove:%f",vehicle_remove->runtime.s);
            for(int i=0;i<lane[tid].remove_num;i++){
                if(lane[tid].first!=vehicle_remove){
                    vehicle_remove->relation[0][1]->relation[1][1] = vehicle_remove->relation[1][1]; //01代表此车的前一辆车
                    //remove掉的车是否需要此刻更新relation
                    vehicle_remove->relation[1][1]->relation[0][1] = vehicle_remove->relation[0][1];
                }else {
                    lane[tid].first = vehicle_remove->relation[1][1];
                }

            }
        }    
struct DeviceVehicle *veh;
struct DeviceVehicle *veh1;
//insert_buffer测试
//完成了insert_buffer链表构建
int g = 0;

if(lane[tid].insert_num!=0){
    
    veh1 = lane[tid].insert_buffer[0];
    int num =lane[tid].insert_num;
while(num!=0){
    g++;
    if(lane[tid].insert_buffer[g]!=0){
        veh1->runtime.lane = &lane[tid];
        veh1->snapshot.lane = &lane[tid];
        veh = lane[tid].insert_buffer[g];
        veh1->relation[1][1] = veh;
        veh1 = veh;
    }
    num--;
    
}
lane[tid].insert_buffer[lane[tid].insert_num-1]->runtime.lane= &lane[tid];
lane[tid].insert_buffer[lane[tid].insert_num-1]->snapshot.lane= &lane[tid];
// struct DeviceVehicle *u = lane[tid].insert_buffer[0];
// num = lane[tid].insert_num;
// int y=0;
// while(num!=0){
    
//     printf("index%d:u[%d]:%f\n",lane[tid].index,y,u->runtime.s);
//     if(u->relation[1][1]!=0){
//         u = u->relation[1][1];
//     }
//     num--;
//     y++;
// }
// if(veh1==0){
//     printf("1242");
// }
lane[tid].num_vehicles =lane[tid].num_vehicles + lane[tid].insert_num;
lane[tid].insert_buffer[lane[tid].insert_num-1]->relation[1][1]=0;
// if(lane[tid].insert_buffer[lane[tid].insert_num-1]->relation[1][1]=!0){
//     printf("111\n");
// }
}

//主链的更新与排序
if(lane[tid].insert_buffer[0]!=0){
//路上没车，insert_buffer有车
if(lane[tid].first->relation[1][1]==0){
    lane[tid].first->relation[1][1] = lane[tid].insert_buffer[0];
    // struct DeviceVehicle *pre = lane[tid].first->relation[1][1];
    // int a =1;
    // int num = lane[tid].insert_num;
    // while((num-1)!=0){
    //     pre->relation[1][1] = lane[tid].insert_buffer[a];
    //     a++;
    //     pre = pre->relation[1][1];
    // }

}
//路上有车，insert_buffer有车
else{
    //printf("num_vehicles:%d\n",lane[tid].num_vehicles);
    struct DeviceVehicle *p2 = lane[tid].first->relation[1][1];//指向A链的指针
    struct DeviceVehicle *q2= lane[tid].insert_buffer[0];
    struct DeviceVehicle *r2,*s2;
    int i =0;
    r2 = lane[tid].first;//r始终指向C链
    s2 = r2;
    //作为双链表中指向r前一辆车的指针
    while (p2!=0 &&q2!= 0) {//当p、q都不为空时，选取p与q所指结点中的较小者插入C的尾部
		//尾插法建立双向循环链表
        if(p2!=0){
            //printf("111\n");
        }
        if(q2!=0){
             //printf("222\n");
        }
		if (p2->runtime.s < q2->runtime.s) {
			r2->relation[1][1] = p2;
			p2 = p2->relation[1][1];

			r2 = r2->relation[1][1];
            r2->relation[0][1] =s2;
            s2 = r2;

		}
		else
		{
           
            //lane[tid].num_vehicles++;
			r2->relation[1][1] = q2;
            
            q2 = q2->relation[1][1];

			r2 = r2->relation[1][1];
            r2->relation[0][1] = s2;
            s2 = r2;
            //std::cout<<r->relation[0][1]->runtime.s<<std::endl;


                    //lane[tid].num_vehicles++;
        // r2->relation[1][1] = q2;
        // q2->runtime.lane = &lane[tid];
        // q2->snapshot.lane = &lane[tid];
        // i++;
        // q2 = lane[tid].insert_buffer[i];
        // r2 = r2->relation[1][1];
        // //std::cout<<r->runtime.s<<std::endl;
        // r2->relation[0][1] = s2;
        // s2 =r2;
            
		}
    }
    if(p2!=0){
        r2->relation[1][1] = p2;
        p2->relation[0][1] =r2;
    }
    if(q2!=0){
        r2->relation[1][1] = q2;
        q2->relation[0][1] =r2;
    }
    // while (q2!=0) {
        
    //     lane[tid].num_vehicles++;
    //     r2->relation[1][1] = q2;
    //     q2->runtime.lane = &lane[tid];
    //     q2->snapshot.lane = &lane[tid];
    //     i++;
    //     q2 = lane[tid].insert_buffer[i];
    //     r2 = r2->relation[1][1];
    //     //std::cout<<r->runtime.s<<std::endl;
    //     r2->relation[0][1] = s2;
    //     s2 =r2;

    // }

    // while (p2!=0) {
        
    //     //printf("r->attribute.max_acc:%f",r->attribute.max_acc);
    //     p2 = p2->relation[1][1];
    //     r2->relation[1][1] = p2;
    //     r2 = r2->relation[1][1];
    //     r2->relation[0][1] = s2;
    //     s2 = r2;
    // }
    // lane[tid].end->relation[0][1] = r2;
}

}
// if(button==0 &&lane[tid].first->relation[1][1]!=0 ){
//     int i= 0;
//     struct DeviceVehicle *p2 = lane[tid].first->relation[1][1];//指向A链的指针
//     //printf("p->attribute.max_acc:%f",p->attribute.max_acc);
//     struct DeviceVehicle *q2= 0;
//     struct DeviceVehicle *r2,*s2;
//     r2 = lane[tid].first;//r始终指向C链
    
//     s2 = r2;//作为双链表中指向r前一辆车的指针
   
// //当insert_buffer不为空时，对两个链表进行排序
//     if(lane[tid].insert_buffer[i]!=0){
//     q2 = lane[tid].insert_buffer[i];//指向B链的指针  
//     }
   
//     //printf("q->attribute.length:%f",q->attribute.length); 
	
//     //printf("q->attribute.length:%f",q2->attribute.length);

//     while (q2!=0) {
        
//         lane[tid].num_vehicles++;
//         r2->relation[1][1] = q2;
//         q2->runtime.lane = &lane[tid];
//         q2->snapshot.lane = &lane[tid];
//         i++;
//         q2 = lane[tid].insert_buffer[i];
//         r2 = r2->relation[1][1];
//         //std::cout<<r->runtime.s<<std::endl;
//         r2->relation[0][1] = s2;
//         s2 =r2;

//     }
//     lane[tid].end->relation[0][1] = r2; 
//     lane[tid].end->relation[1][1] = 0;
//     //printf("r->runtime.s:%f",r->runtime.s);
//     if(p2!=0){
//         r2->relation[1][1] = p2;
//         p2->relation[0][1] =r2;
//     }
//     while (p2!=0) {
        
//         //printf("r->attribute.max_acc:%f",r->attribute.max_acc);
//         p2 = p2->relation[1][1];
//         r2 = r2->relation[1][1];
//         r2->relation[0][1] = s2;
//         s2 = r2;
//     }

    

    
// }
if(lane[tid].left_lane!=0){
    //printf("lane[tid].left_lane:%f\n",lane[tid].left_lane->length);
    if(lane[tid].left_lane->first->relation[1][1]!=0){
    printf("runtime.s:%f\n",lane[tid].left_lane->first->relation[1][1]->runtime.s);
    if(lane[tid].first->relation[1][1]!=0){
        printf("s:%f\n",lane[tid].first->relation[1][1]->runtime.s);
        
    }
}
}

if(lane[tid].first->relation[1][1]!=0){
    
}

//     //printf("%d:%d\n",tid,lane[tid].num_vehicles);
    //支链“重置”
    //测报错位置
    //std::cout<<index->first->relation[1][1]->runtime.s<<std::endl;
    //printf("num_vehicles:%d",lane[tid].num_vehicles);
        if (lane[tid].left_lane!=0&&lane[tid].left_lane->first->relation[1][1]!=0&&lane[tid].first->relation[1][1]!=0) {//所有的lane
            //测报错位置
            printf("支链“重置”");
            //std::cout<<index->first->relation[1][1]->runtime.s<<std::endl;
            struct DeviceLane *l1 = &lane[tid];//指向当前车道指针
            struct DeviceLane *l2 = lane[tid].left_lane;//指向当前车道左车道指针
            int max_vehicles;
            // //测报错位置
            // //std::cout<<p2->first->relation[1][1]->runtime.s<<std::endl;
            struct DeviceVehicle *v1 = l1->first->relation[1][1];
            struct DeviceVehicle *v2 = l2->first->relation[1][1];
            float v1Tov2S;//v1在v2上对应的位置S
            while(v1&&v2){//               
                v1Tov2S = v1->runtime.s*l2->length/l1->length;
                if (v1Tov2S>v2->runtime.s) {
                    //v1->relation[0][0] = v2;//v2暂时在v1的左后方
                    v2->relation[1][2] = v1;//v2右前方先指向v1
                    v2->relation[0][2] = v1->relation[0][1];
                    v2 = v2->relation[1][1];//v2指向主链下一车
                }
                else{
                    //v2->relation[0][2] = v1;//v1暂时在v2的右后方
                    v1->relation[1][0] = v2;//v1左前方指向v2
                    v1->relation[0][0] = v2->relation[0][1]; //v1左后赋值
                    v1 = v1->relation[1][1];//v1指向主链下一车
                }                
            }
            // //测报错位置
            // //std::cout<<v1->runtime.s<<std::endl;

            // //std::cout<<v2->relation[0][1]->runtime.s<<std::endl;
            if(v1 ==0){
                // p1->end = v1;
                // v1 = v1->relation[0][1];
                v1 = l1->end->relation[0][1];
            }
            if(v2==0){
                //p2->end = v2;
                v2 = l2->end->relation[0][1];
                //v2 = v2->relation[0][1];//两种情况都可以使指针回到最后一辆车
                //std::cout<<v2->runtime.s<<std::endl;

            }

            while (v1!=0) {//v1还有多的车辆
                v1->relation[0][0] = v2;//剩余所有v1的车的左前方都为v2的最后一辆车
                v1 = v1->relation[1][1];
            }
            v1 = 0;
            while (v2!=0) {//v2还有多的车辆
                v2->relation[0][2] = v1;
                v2 = v2->relation[1][1];
            }
            v2 =0;

        }
       
        for(int i=0;i<lane[tid].insert_num;i++){
            lane[tid].insert_buffer[i]= 0;

        }
        
        lane[tid].insert_num = 0;
        InsertNum[tid]=0;
}   


