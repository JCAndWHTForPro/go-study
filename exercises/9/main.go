package main

/*
*
## 第 9 课：错误处理

### 题 9.1 安全除法（⭐）

写一个函数 `divide(a, b float64) (float64, error)`，当 b 为 0 时返回自定义错误 `DivideByZeroError`。用 `errors.Is` 判断错误类型。

### 题 9.2 多层错误包装（⭐⭐）

模拟三层调用链：`handler() → service() → dao()`。
dao 返回原始错误，service 用 `%w` 包装，handler 用 `errors.As` 提取原始错误。打印整个错误链。
*/
func main() {}
