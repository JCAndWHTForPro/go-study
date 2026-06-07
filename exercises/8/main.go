package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*
## 第 8 课：并发

### 题 8.1 并发求和（⭐⭐）

把一个 100 个元素的 slice 分成 4 段，每段用一个 goroutine 求和，通过 channel 汇总，打印总和。

### 题 8.2 生产者消费者（⭐⭐⭐）

启动 3 个生产者 goroutine（各生产 5 个商品），1 个消费者 goroutine 消费所有商品。用 channel 通信，所有生产完成后关闭 channel，
消费者用 for-range 读取。

### Channel Close 注意事项总结

1. **先 Wait 再 Close**：必须确保所有往 channel 写入的 goroutine 都完成后，才能 close(channel)。
   否则 goroutine 向已关闭的 channel 发送数据会触发 panic: send on closed channel。
   正确写法：sy.Wait() → close(channel)，而非 close(channel) → sy.Wait()。

2. **只能关闭一次**：对同一个 channel 调用两次 close 也会 panic: close of closed channel。

3. **由发送方关闭**：channel 应由发送方（生产者）负责关闭，接收方（消费者）不应关闭 channel。

4. **关闭后仍可读取**：已关闭的 channel 仍可读取缓冲区中剩余的数据，读完后返回零值。

5. **for-range 自动退出**：用 for v := range ch 读取 channel 时，channel 被关闭且读完后循环自动结束。
   如果不 close，for-range 会永远阻塞导致 deadlock。

6. **判断 channel 是否关闭**：v, ok := <-ch，ok 为 false 表示 channel 已关闭且无数据。
*/
func main() {
	goroutineSum()
}

func goroutineSum() {
	array := []int{}
	for i := 0; i < 100; i++ {
		array = append(array, rand.Intn(100))
	}
	channel := make(chan int, 4)
	length := len(array)
	var sy sync.WaitGroup
	fmt.Println("当前的slice长度是：", length, ",这个整个数据打印一下：", array)
	for start := 0; start < length; start += length / 4 {
		end := start + length/4 - 1
		sy.Add(1)
		go func(s, e int) {
			defer sy.Done()
			sum := 0
			for i := s; i <= e; i++ {
				sum += array[i]
			}
			channel <- sum
		}(start, end)
	}
	go func() {
		sy.Wait()      // 先等所有 goroutine 写入完成
		close(channel) // 再关闭 channel，for-range 才能安全退出
	}()
	result := 0
	for rrr := range channel {
		result += rrr
	}
	//close(channel)
	fmt.Println("最后的和是", result)
	res := 0
	for _, i := range array {
		res += i
	}
	fmt.Println("这里是验证的结果：", res)

}
