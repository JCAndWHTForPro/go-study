// 第 8 课：并发（goroutine / channel / select / sync）
// 运行方式：go run ./08-concurrency
//
// Go 的并发模型核心理念：「不要通过共享内存来通信，而是通过通信来共享内存」
// 两大武器：
//   goroutine —— 轻量级协程，用 go 关键字启动，一个程序可以跑几十万个
//   channel   —— 协程间的「管道」，用来安全地传递数据
package main

import (
	"fmt"
	"sync"
	"time"
)

// Task 用于演示 chan *T（channel 传指针）
type Task struct {
	ID     int
	Result string
}

func main() {
	// ===== 1. goroutine：用 go 关键字启动 =====
	// go f() 就启动了一个新协程，和 main 并发执行
	go sayHello("goroutine-1")
	go sayHello("goroutine-2")

	// 注意：main 退出后所有 goroutine 都会被杀掉
	// 这里先用 Sleep 等一下，后面会用更优雅的方式（channel / WaitGroup）
	time.Sleep(500 * time.Millisecond)
	fmt.Println("--- 1. goroutine 演示完毕 ---\n")

	// ===== 2. channel：协程间传数据的管道 =====
	// make(chan 类型) 创建 channel，<- 发送/接收
	ch := make(chan string) // 无缓冲 channel

	go func() {
		ch <- "hello from goroutine" // 发送数据到 channel
	}()

	msg := <-ch // 接收数据（会阻塞直到有数据）
	fmt.Println("收到:", msg)

	// ===== 2.1 用 channel 等待 goroutine 完成（替代 Sleep）=====
	done := make(chan bool)
	go func() {
		fmt.Println("后台任务执行中...")
		time.Sleep(200 * time.Millisecond)
		fmt.Println("后台任务完成")
		done <- true // 通知 main 已完成
	}()
	<-done // 阻塞等待完成信号
	fmt.Println("--- 2. channel 演示完毕 ---\n")

	// ===== 3. 带缓冲的 channel =====
	// make(chan T, 容量)：发送不会阻塞，直到缓冲满
	bufCh := make(chan int, 3) // 缓冲容量 3
	bufCh <- 1
	bufCh <- 2
	bufCh <- 3
	// bufCh <- 4  // ❌ 缓冲满了会阻塞（如果没人接收就死锁）
	fmt.Println("缓冲 channel 取出:", <-bufCh, <-bufCh, <-bufCh)

	// ===== 4. range 遍历 channel + close 关闭 =====
	numCh := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			numCh <- i
		}
		close(numCh) // 发完了就关闭，告诉接收方「没有更多数据了」
	}()
	fmt.Print("range 遍历 channel: ")
	for n := range numCh { // 自动在 close 后退出循环
		fmt.Print(n, " ")
	}
	fmt.Println("\n--- 4. range + close 演示完毕 ---\n")

	// ===== 5. select：同时监听多个 channel =====
	// select 类似 switch，但每个 case 是 channel 操作
	ch1 := make(chan string)
	ch2 := make(chan string)
	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自 ch1"
	}()
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "来自 ch2"
	}()

	// 等待两个 channel，谁先到就先处理谁
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("select 收到:", msg1)
		case msg2 := <-ch2:
			fmt.Println("select 收到:", msg2)
		}
	}
	fmt.Println("--- 5. select 演示完毕 ---\n")

	// ===== 6. sync.WaitGroup：等待一组 goroutine 全部完成 =====
	// 比用 channel 计数更方便，实战中非常常用
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1) // 每启动一个 goroutine 就 +1
		go func(id int) {
			defer wg.Done() // 完成后 -1（defer 确保一定执行）
			fmt.Printf("  worker %d 开始\n", id)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("  worker %d 完成\n", id)
		}(i)
	}
	wg.Wait() // 阻塞直到计数归零（所有 worker 完成）
	fmt.Println("--- 6. WaitGroup 演示完毕 ---\n")

	// ===== 7. sync.Mutex：互斥锁，保护共享变量 =====
	// 多个 goroutine 同时改一个变量会「数据竞争」，用 Mutex 加锁保护
	counter := 0
	var mu sync.Mutex
	var wg2 sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mu.Lock()   // 加锁：同一时刻只有一个 goroutine 能进入
			counter++
			mu.Unlock() // 解锁
		}()
	}
	wg2.Wait()
	fmt.Println("1000 个 goroutine 安全累加后 counter =", counter) // 一定是 1000
	fmt.Println("--- 7. Mutex 演示完毕 ---\n")

	// ===== 8. chan *T：channel 传指针 =====
	// make 返回 channel 本身（值），但 channel 里传什么类型由你定义
	// 这里传 *Task 指针：避免大结构体拷贝，接收方还能修改原对象
	taskCh := make(chan *Task, 3) // chan *Task：传递 Task 指针的 channel

	// 生产者：创建 Task 指针发送到 channel
	go func() {
		for i := 1; i <= 3; i++ {
			taskCh <- &Task{ID: i, Result: "待处理"} // 发送指针
		}
		close(taskCh)
	}()

	// 消费者：接收指针，直接修改原对象
	for task := range taskCh {
		// task 的类型是 *Task（指针），可以直接修改
		task.Result = fmt.Sprintf("任务%d已完成", task.ID)
		fmt.Printf("  处理 -> ID=%d, Result=%s\n", task.ID, task.Result)
	}

	// 对比：chan Task（传值）vs chan *Task（传指针）
	// chan Task:   每次发送都拷贝整个结构体，接收方改的是副本
	// chan *Task:  只传 8 字节地址，接收方能改原对象，更高效
	fmt.Println("--- 8. chan *T 演示完毕 ---")
}

func sayHello(name string) {
	for i := 0; i < 3; i++ {
		fmt.Printf("  %s: 第 %d 次打招呼\n", name, i+1)
		time.Sleep(50 * time.Millisecond)
	}
}
