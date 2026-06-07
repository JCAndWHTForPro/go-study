package main

/*
*
## 第 4 课：函数

### 题 4.1 递归阶乘（⭐）

写一个递归函数 `factorial(n int) int`，计算 n 的阶乘。测试 `factorial(5)` 应该返回 120。

### 题 4.2 闭包计数器（⭐⭐）

写一个函数 `makeCounter() func() int`，每次调用返回的函数时，返回值递增（1, 2, 3, ...）。用闭包实现。

### 题 4.3 defer 执行顺序（⭐⭐）

写一个函数，在 for 循环中 defer 打印 0-4，观察并解释为什么输出是 4, 3, 2, 1, 0（栈顺序）。
*/
/*
⚠️ 易错点：string() 不是数字转字符串！

Go 中 string(120) 不会得到 "120"，而是得到 "x"（ASCII 码 120 = 'x'）。
string() 的作用是将整数当作 Unicode 码点，转成对应的字符。

❌ 错误写法：
  fmt.Println("结果：" + string(factorial(5)))  // 输出：结果：x

✅ 正确的数字转字符串方式：
  1. strconv.Itoa(n)             — 推荐，专门用于 int → string
  2. fmt.Sprintf("%d", n)        — 格式化方式，更灵活
  3. fmt.Println("结果：", n)     — 直接用逗号，Println 自动处理
*/
func main() {
	//fmt.Println("factorial(5)返回的结果：", factorial(5))
	/*counter := makeCounter()
	println("第一次调用", counter())
	println("第2次调用", counter())
	println("第3次调用", counter())
	println("第4次调用", counter())*/
	deferPrint()
}

func factorial(n int) int {
	if n == 1 {
		return n
	}
	return factorial(n-1) * n
}

func makeCounter() func() int {
	var count int = 0
	return func() int {
		count = count + 1
		return count
	}
}
func deferPrint() {
	for i := 0; i < 5; i++ {
		defer println("当前输出的是", i)
	}
}
