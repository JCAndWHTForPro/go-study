// 第 16 课：Context（上下文控制）
// 运行方式：go run ./16-context
//
// Context 是 Go 并发编程的"遥控器"，用来控制 goroutine 的生命周期。
// 核心能力：超时控制、取消传播、携带请求级别的数据。
//
// 为什么需要 Context？
//   场景：HTTP 请求进来 → 查数据库 → 调远程服务 → 返回结果
//   如果用户中途断开连接，后面的数据库查询和远程调用都应该立即停止，
//   否则白白浪费资源。Context 就是用来通知所有下游"别干了，取消吧"。
//
// ============================================================
// 【总结】Context 核心概念
// ------------------------------------------------------------
//
// ▶ 四种创建方式：
//   context.Background()           → 根 context，一般在 main/初始化时用
//   context.TODO()                 → 占位用，还不确定用哪个时临时用
//   context.WithCancel(parent)     → 手动取消：返回 ctx + cancel 函数
//   context.WithTimeout(parent, d) → 超时自动取消：d 时间后自动 cancel
//   context.WithDeadline(parent, t)→ 截止时间取消：到 t 时刻自动 cancel
//   context.WithValue(parent, k, v)→ 携带数据（谨慎使用，别滥用）
//
// ▶ 核心方法（context.Context 接口）：
//   ctx.Done()     → 返回一个 channel，ctx 被取消时会关闭
//   ctx.Err()      → 取消原因：context.Canceled 或 context.DeadlineExceeded
//   ctx.Deadline() → 返回截止时间（如果有）
//   ctx.Value(key) → 获取携带的数据
//
// ▶ 使用套路：
//   select {
//   case <-ctx.Done():   // ctx 被取消了，停止工作
//       return ctx.Err()
//   case result := <-ch: // 正常拿到结果
//       return result
//   }
//
// ▶ 重要规则：
//   1. Context 是第一个参数：func DoSomething(ctx context.Context, ...)
//   2. 不要把 Context 存到 struct 里，每次调用传进去
//   3. cancel() 必须调用！通常 defer cancel()，防止资源泄漏
//   4. WithValue 只放请求级别的数据（traceID、userID），不放业务参数
//
// ▶ Context 的传播是树状的：
//   Background
//     └── WithTimeout (5s)         ← 父 context
//           ├── WithCancel         ← 子 context（父取消时自动取消）
//           └── WithValue(traceID) ← 子 context
//   父 context 取消 → 所有子 context 自动取消（向下传播）
//   子 context 取消 → 不影响父 context（不向上传播）
//
// ▶ 对比 Java：
//   Java 用 Future.cancel() / ExecutorService.shutdownNow() 取消任务
//   Go 用 context 统一管理，更优雅，贯穿整个调用链
// ============================================================
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// ===== 1. WithCancel：手动取消 =====
	fmt.Println("===== 1. WithCancel =====")
	ctx, cancel := context.WithCancel(context.Background())

	// 启动一个 goroutine 模拟后台任务
	go worker(ctx, "任务A")

	time.Sleep(300 * time.Millisecond) // 让任务跑一会儿
	cancel()                           // 手动取消！worker 会收到通知
	time.Sleep(100 * time.Millisecond) // 等 worker 打印退出消息

	// ===== 2. WithTimeout：超时自动取消 =====
	fmt.Println("\n===== 2. WithTimeout =====")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2() // 即使超时也要调 cancel，释放资源

	// 模拟一个可能很慢的操作
	result, err := slowOperation(ctx2, 200*time.Millisecond) // 200ms < 500ms 超时
	if err != nil {
		fmt.Println("操作失败:", err)
	} else {
		fmt.Println("操作成功:", result)
	}

	// 这次故意让操作超时
	ctx3, cancel3 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel3()

	result, err = slowOperation(ctx3, 500*time.Millisecond) // 500ms > 100ms 超时
	if err != nil {
		fmt.Println("操作超时:", err) // context deadline exceeded
	} else {
		fmt.Println("操作成功:", result)
	}

	// ===== 3. WithValue：携带请求级别的数据 =====
	fmt.Println("\n===== 3. WithValue =====")
	ctx4 := context.WithValue(context.Background(), "traceID", "abc-123-xyz")
	ctx4 = context.WithValue(ctx4, "userID", 42) // 可以链式添加多个值
	processRequest(ctx4)

	// ===== 4. Context 传播：父取消，子自动取消 =====
	fmt.Println("\n===== 4. 取消传播 =====")
	parentCtx, parentCancel := context.WithCancel(context.Background())

	// 从父 context 派生两个子 context
	childCtx1, childCancel1 := context.WithCancel(parentCtx)
	defer childCancel1()
	childCtx2, childCancel2 := context.WithCancel(parentCtx)
	defer childCancel2()

	go worker(childCtx1, "子任务1")
	go worker(childCtx2, "子任务2")

	time.Sleep(300 * time.Millisecond)
	fmt.Println("→ 取消父 context（所有子任务应该自动停止）")
	parentCancel() // 取消父 → 子 1 和子 2 自动取消
	time.Sleep(200 * time.Millisecond)

	// ===== 5. 实际场景：模拟 HTTP 请求处理 =====
	fmt.Println("\n===== 5. 模拟 HTTP 请求处理 =====")
	handleRequest()
}

// worker 模拟一个持续运行的后台任务，收到取消信号后退出
func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done(): // ctx 被取消时，Done() 的 channel 会关闭
			fmt.Printf("  [%s] 收到取消信号，退出。原因: %v\n", name, ctx.Err())
			return
		default:
			fmt.Printf("  [%s] 工作中...\n", name)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// slowOperation 模拟一个可能很慢的操作（如数据库查询、远程调用）
func slowOperation(ctx context.Context, duration time.Duration) (string, error) {
	select {
	case <-ctx.Done(): // 超时或被取消
		return "", ctx.Err()
	case <-time.After(duration): // 操作完成
		return "查询结果: 42", nil
	}
}

// processRequest 演示从 context 中取值
func processRequest(ctx context.Context) {
	traceID := ctx.Value("traceID")
	userID := ctx.Value("userID")
	fmt.Printf("  处理请求: traceID=%v, userID=%v\n", traceID, userID)

	// 取不存在的 key 返回 nil
	missing := ctx.Value("notExist")
	fmt.Println("  不存在的 key:", missing) // <nil>
}

// handleRequest 模拟一个完整的 HTTP 请求处理流程
func handleRequest() {
	// 模拟：请求最多允许 1 秒
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 添加请求级别的数据
	ctx = context.WithValue(ctx, "traceID", "req-001")

	fmt.Println("  开始处理请求...")

	// 步骤 1：查数据库（300ms）
	dbResult, err := queryDatabase(ctx)
	if err != nil {
		fmt.Println("  数据库查询失败:", err)
		return
	}
	fmt.Println("  数据库结果:", dbResult)

	// 步骤 2：调远程服务（400ms）
	apiResult, err := callRemoteAPI(ctx)
	if err != nil {
		fmt.Println("  远程调用失败:", err)
		return
	}
	fmt.Println("  远程服务结果:", apiResult)

	fmt.Println("  ✅ 请求处理完成！")
}

func queryDatabase(ctx context.Context) (string, error) {
	traceID := ctx.Value("traceID")
	fmt.Printf("  [DB] traceID=%v, 查询中...\n", traceID)
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("数据库查询被取消: %w", ctx.Err())
	case <-time.After(300 * time.Millisecond):
		return "用户列表: [张三, 李四]", nil
	}
}

func callRemoteAPI(ctx context.Context) (string, error) {
	traceID := ctx.Value("traceID")
	fmt.Printf("  [API] traceID=%v, 调用中...\n", traceID)
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("远程调用被取消: %w", ctx.Err())
	case <-time.After(400 * time.Millisecond):
		return "天气: 晴, 25°C", nil
	}
}
