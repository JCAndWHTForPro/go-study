package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
题目：并发价格查询器
模拟一个电商场景：同时从 3 个"供应商"查询同一商品的价格，取最低价返回。要求：

写 3 个函数模拟 3 个供应商的查询（用 time.Sleep 模拟网络延迟，随机返回价格）
用 goroutine 并发查询 3 个供应商
3 秒超时：如果某个供应商超时没返回，忽略它
从已返回的结果中取最低价
打印：最低价、来自哪个供应商、总耗时
Go
// 期望输出类似：
// [供应商A] 返回价格: 299.00 (耗时 1.2s)
// [供应商C] 返回价格: 259.00 (耗时 2.1s)
// [供应商B] 超时未返回
// 最低价: 259.00 来自供应商C，总耗时: 2.1s
提示：会用到 goroutine、channel、select、time.After、sync.WaitGroup
*/
type Supplier struct {
	Name    string
	Price   float32
	Msg     string
	Success bool
}

func main() {
	var wg sync.WaitGroup
	result := make(chan Supplier, 3)
	wg.Add(3)
	go produce1("供应啥A", &wg, result)
	go produce1("供应啥B", &wg, result)
	go produce1("供应啥C", &wg, result)
	go func() {
		wg.Wait()
		close(result)
	}()
	maxSupp := Supplier{}
	for pro := range result {
		fmt.Println(pro.Msg)
		if pro.Success && pro.Price > maxSupp.Price {
			maxSupp = pro
		}
	}
	fmt.Println("", maxSupp.Price, "来自", maxSupp.Name, "总耗时:")
}

func produce1(name string, wg *sync.WaitGroup, result chan<- Supplier) {
	defer wg.Done()
	now := time.Now()
	rstChan := make(chan float32, 1)
	//defer close(rstChan)
	go delayReturn(rstChan)
	select {
	case price := <-rstChan:
		since := time.Since(now)
		result <- Supplier{
			Name:    name,
			Price:   price,
			Msg:     fmt.Sprintf("[%s] 返回价格: %.2f (耗时 %v)", name, price, since),
			Success: true,
		}
	case <-time.After(3 * time.Second):
		result <- Supplier{
			Msg:     fmt.Sprintf("[%s] 差事未返回", name),
			Success: false,
		}
	}
}
func delayReturn(rstChan chan<- float32) {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	rstChan <- float32(rand.Intn(499) + 1)
}
