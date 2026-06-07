package main

import (
	"fmt"
	"go-study/exercises/mathutil"
)

/*
## 第 11 课：包管理

### 题 11.1 工具包（⭐⭐）

在项目中创建一个 `exercises/mathutil` 包，提供以下导出函数：

  - `GCD(a, b int) int` — 最大公约数
  - `LCM(a, b int) int` — 最小公倍数
  - `IsPrime(n int) bool` — 判断质数
    在 main 中导入并测试。
*/
func main() {
	// GCD 测试
	fmt.Println("GCD(12, 8) =", mathutil.GCD(12, 8)) // 4
	fmt.Println("GCD(15, 5) =", mathutil.GCD(15, 5)) // 5
	fmt.Println("GCD(7, 3) =", mathutil.GCD(7, 3))   // 1

	// LCM 测试
	fmt.Println("LCM(4, 6) =", mathutil.LCM(4, 6)) // 12
	fmt.Println("LCM(3, 7) =", mathutil.LCM(3, 7)) // 21

	// IsPrime 测试
	for _, n := range []int{-1, 0, 1, 2, 3, 4, 7, 9, 11, 15, 17} {
		fmt.Printf("IsPrime(%d) = %t\n", n, mathutil.IsPrime(n))
	}
}
