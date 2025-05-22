/*
Go 代码详细讲解
这个 Go 代码演示了如何用 context.Context 实现“优雅取消 goroutine”的功能。
主线程启动一个工作 goroutine，2 秒后通过 context 通知 goroutine 停止工作。
*/
package main // 声明包名

import (
    "context" // 导入 context 包，用于控制 goroutine 生命周期
    "fmt"     // 格式化输出
    "time"    // 时间相关操作
)

func doWork(ctx context.Context) { // 工作函数，接收 context
    for i := 0; ; i++ {           // 无限循环，i 递增
        select {
        case <-ctx.Done():        // 如果 context 被取消
            fmt.Println("work canceled:", ctx.Err()) // 打印取消原因
            return                // 退出 goroutine
        default:                  // 否则继续工作
            fmt.Println("working:", i) // 打印当前工作进度
            time.Sleep(500 * time.Millisecond) // 模拟工作耗时
        }
    }
}

func main() {
    // 普通取消演示
    go func() {
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()
        fmt.Println("[WithCancel] goroutine start")
        go doWork(ctx)
        time.Sleep(2 * time.Second)
        fmt.Println("[WithCancel] call cancel()")
        cancel()
        time.Sleep(1 * time.Second)
    }()

    // 超时取消演示
    go func() {
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        fmt.Println("[WithTimeout] goroutine start")
        go doWork(ctx)
        time.Sleep(3 * time.Second)
        fmt.Println("[WithTimeout] main end")
    }()

    // 截止时间取消演示
    go func() {
        deadline := time.Now().Add(2 * time.Second)
        ctx, cancel := context.WithDeadline(context.Background(), deadline)
        defer cancel()
        fmt.Println("[WithDeadline] goroutine start")
        go doWork(ctx)
        time.Sleep(3 * time.Second)
        fmt.Println("[WithDeadline] main end")
    }()

    // 传递元数据演示
    go func() {
        type keyType string
        ctx := context.WithValue(context.Background(), keyType("userID"), 12345)
        go func(ctx context.Context) {
            for i := 0; i < 3; i++ {
                userID := ctx.Value(keyType("userID"))
                fmt.Printf("[WithValue] working: %d, userID: %v\n", i, userID)
                time.Sleep(500 * time.Millisecond)
            }
            fmt.Println("[WithValue] goroutine end")
        }(ctx)
    }()

    // 主线程等待所有演示完成
    time.Sleep(5 * time.Second)
}