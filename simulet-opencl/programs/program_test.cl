__kernel
void test(__global int * num){
    //printf("test1\n");
size_t tid = get_global_id(0);
if(num[tid]>1){ 
    printf("num:%d\n",num[tid]);
 }
}