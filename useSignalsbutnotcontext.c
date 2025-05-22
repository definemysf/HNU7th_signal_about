#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>
#include <stdatomic.h>

// 取消信号，主线程设置，工作线程检测
volatile int cancel_flag = 0; 

// 工作线程执行的函数
void* do_work(void* arg) {
    int i = 0;
    while (1) {
        // 检查是否收到取消信号
        if (cancel_flag) { 
            printf("son work canceled\n");
            break; // 收到信号后退出循环
        }
        printf("son is working: %d\n", i++); // 打印工作进度
        usleep(250000); // 休眠500毫秒
    }
    return NULL; // 线程返回
}

int main() {
    pthread_t tid; // 线程ID变量
    printf("[parent] main thread start\n"); // 主线程开始
    // 创建工作线程，执行do_work函数
    pthread_create(&tid, NULL, do_work, NULL); 

    printf("[parent] main thread sleep 2s, let child work...\n"); // 主线程等待2秒，让工作线程运行一段时间
    sleep(2); 

    printf("[parent] main thread set cancel_flag=1, notify child to exit\n"); // 设置取消信号，通知工作线程退出
    cancel_flag = 1; 

    printf("[parent] main thread sleep 1s, wait child to exit...\n"); // 再等待1秒，确保工作线程检测到信号并退出
    sleep(1); 

    printf("[parent] main thread wait for child to join\n"); // 等待工作线程结束，回收资源
    pthread_join(tid, NULL); 
    printf("[parent] main thread end\n"); // 程序结束
    return 0; // 程序结束
}