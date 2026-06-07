package main

import (
	"fmt"
	"math/rand"
)

/*
### 题 3.1 FizzBuzz（⭐）

打印 1 到 30，如果是 3 的倍数打印 "Fizz"，5 的倍数打印 "Buzz"，同时是 3 和 5 的倍数打印 "FizzBuzz"，否则打印数字本身。

### 题 3.2 九九乘法表（⭐⭐）

用嵌套 for 循环打印九九乘法表，格式对齐。

### 题 3.3 猜数字（⭐⭐）

用 `math/rand` 生成一个 1-100 的随机数，模拟猜 10 次（用 for 循环，每次猜中间值，类似二分查找），打印每次猜的过程和最终结果。

### 二分查找三种写法对比

1. 左闭右闭 [start, end]（标准写法）：
   循环条件：start <= end
   缩边界：start = mid+1 / end = mid-1

2. 左闭右开 [start, end)：
   循环条件：start < end
   缩边界：start = mid+1 / end = mid

3. 九章算法"相邻退出"模板（最保险，不易死循环/越界）：
   循环条件：start+1 < end         ← 退出时 start 和 end 相邻
   缩边界：start = mid / end = mid  ← 不加不减！
   退出后：单独判断 start 和 end 哪个是答案
   优点：不用纠结 <= 还是 <，不用纠结 mid±1，不会死循环
*/
func main() {
	//FizzBuzz()
	//PrintMultiplicationTable()
	randGuess()
}
func randGuess() {
	start, end := 1, 100
	number := rand.Intn(100) + 1
	for i := 0; i < 10 && start <= end; i++ {
		mid := start + (end-start)/2
		fmt.Printf("当前猜测的值是：%d\n", mid)
		if mid == number {
			fmt.Printf("猜中:%d\n", number)
			return
		}
		if mid < number {
			start = mid + 1
		} else {
			end = mid - 1
		}
	}
}
func PrintMultiplicationTable() {
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d ✖️ %d = %-2d    ", j, i, i*j)
		}
		fmt.Printf("\n")
	}
}
func FizzBuzz() {
	for i := 1; i <= 30; i++ {
		if i%5 == 0 && i%3 == 0 {
			fmt.Print("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Print("Fizz")
		} else if i%5 == 0 {
			fmt.Print("Buzz")
		} else {
			fmt.Print(i)
		}
		fmt.Println("")
	}
}
