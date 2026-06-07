package main

import (
	"errors"
	"fmt"
)

/*
*
## 第 9 课：错误处理

### 题 9.1 安全除法（⭐）

写一个函数 `divide(a, b float64) (float64, error)`，当 b 为 0 时返回自定义错误 `DivideByZeroError`。用 `errors.Is` 判断错误类型。

### 题 9.2 多层错误包装（⭐⭐）

模拟三层调用链：`handler() → service() → dao()`。
dao 返回原始错误，service 用 `%w` 包装，handler 用 `errors.As` 提取原始错误。打印整个错误链。
*/
// 自定义错误类型：实现 error 接口（有 Error() string 方法就行）
type DivideError struct {
	Dividend float64
	Divisor  float64
}

func (e *DivideError) Error() string {
	return fmt.Sprintf("不能用 %.2f 除以 %.2f", e.Dividend, e.Divisor)
}

// 哨兵错误（用于 errors.Is 判断）
var DivideByZeroError = errors.New("DivideByZeroError")

func main() {
	err := handler()
	fmt.Println("完整错误链:", err)

	// errors.As 用法：从包装链中提取出特定类型的错误
	// 第二个参数传目标类型的指针，As 会沿着错误链逐层解包，找到类型匹配的就赋值给它
	var divErr *DivideError
	if errors.As(err, &divErr) {
		fmt.Printf("errors.As 提取成功！原始错误: %s\n", divErr)
		fmt.Printf("  Dividend=%.2f, Divisor=%.2f\n", divErr.Dividend, divErr.Divisor)
	}

	// 对比：errors.Is 判断是否包含某个哨兵错误值
	if errors.Is(err, DivideByZeroError) {
		fmt.Println("errors.Is 也能匹配到哨兵错误")
	}
}

func handler() error {
	err := service()
	if err != nil {
		return fmt.Errorf("handler123:%w\n", err)
	}
	return nil
}

func service() error {
	err := dao()
	if err != nil {
		return fmt.Errorf("service456:%w\n", err)
	}
	return nil
}

func dao() error {
	// 返回自定义错误类型（errors.As 能按类型提取到它）
	return &DivideError{Dividend: 10, Divisor: 0}
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, DivideByZeroError
	}
	return (a / b), nil
}
