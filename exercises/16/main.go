package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
## 第 16 课：Context

### 题 16.1 超时搜索（⭐⭐⭐）

模拟一个搜索场景：同时向 3 个"搜索引擎"发请求（用 goroutine + 随机延迟模拟），设置 500ms 超时。

  - 如果某个引擎在超时内返回 → 打印结果
  - 如果全部超时 → 打印"搜索超时"
    用 `context.WithTimeout` + `select` 实现。
*/
//func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
//	defer cancel()
//
//	// 用 channel 收集结果
//	resultCh := make(chan string, 3)
//
//	// 启动 3 个 goroutine 模拟搜索引擎
//	for i := 0; i < 3; i++ {
//		id := i + 1
//		delay := time.Duration((rand.Intn(5)+1)*100) * time.Millisecond
//		go searchFunc(ctx, id, delay, resultCh)
//	}
//
//	// 等待结果：谁先返回用谁的，超时就放弃
//	select {
//	case result := <-resultCh:
//		fmt.Println("搜索成功:", result)
//	case <-ctx.Done():
//		fmt.Println("搜索超时:", ctx.Err())
//	}
//}

func searchFunc(ctx context.Context, id int, delay time.Duration, ch chan<- string) {
	fmt.Printf("引擎%d 开始搜索，预计耗时 %v\n", id, delay)
	select {
	case <-ctx.Done():
		// 超时或被取消，直接退出
		fmt.Printf("引擎%d 被取消\n", id)
		return
	case <-time.After(delay):
		// 模拟搜索完成
		ch <- fmt.Sprintf("引擎%d 的结果（耗时 %v）", id, delay)
	}
}

/*
### 题 16.2 批量下载器（⭐⭐⭐）

模拟批量下载 5 个文件，要求：

1. 每个下载任务用一个 goroutine，随机耗时 200~800ms
2. 设置总超时 1 秒（context.WithTimeout）
3. 用 channel 收集每个下载的结果（成功 or 被取消）
4. 在 main 中打印：哪些下载成功了，哪些因超时被取消了
5. 所有 goroutine 结束后（用 sync.WaitGroup），程序才退出

提示：
  - 每个 goroutine 里用 select 监听 ctx.Done() 和 time.After()
  - WaitGroup 确保所有 goroutine 都退出后再收尾
  - channel 用 buffered channel（容量 5）

这道题综合了：context + goroutine + channel + select + WaitGroup
*/
type DownInfo struct {
	msg     string
	success bool
}

func main() {
	channal := make(chan DownInfo, 5)
	var wg sync.WaitGroup
	ctx, cancelFunc := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelFunc()
	for i := 0; i < 5; i++ {
		duration := time.Duration((rand.Intn(601) + 200)) * time.Millisecond
		wg.Add(1)
		go download(ctx, duration, channal, &wg, i+1)
	}
	go func() {
		wg.Wait()
		close(channal)
	}()
	for info := range channal {
		fmt.Println(info.msg)
	}

}
func download(ctx context.Context, duration time.Duration, channal chan<- DownInfo, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		channal <- DownInfo{
			msg:     fmt.Sprintf("执行超时,当前任务编号是:%d", i),
			success: false,
		}
	case <-time.After(duration):
		channal <- DownInfo{
			msg:     fmt.Sprintf("下载文件成功,当前任务编号是:%d", i),
			success: true,
		}
	}
}
