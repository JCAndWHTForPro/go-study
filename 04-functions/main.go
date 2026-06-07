// 第 4 课：函数（多返回值 / 命名返回值 / defer / 变长参数 / 闭包）
// 运行方式：go run ./04-functions
package main

import (
	"errors"
	"fmt"
)

func main() {
	// ===== 1. 基础函数调用 =====
	fmt.Println("3 + 5 =", add(3, 5))

	// ===== 2. 多返回值 + 错误处理（Go 最经典的写法）=====
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("出错了:", err)
	} else {
		fmt.Println("10 / 2 =", result)
	}

	// 故意触发错误
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("出错了:", err)
	}

	// ===== 3. 命名返回值 =====
	// 返回值可以提前命名，return 时可以不写变量名（裸 return）
	// 多返回值不能直接塞进单值位置，要先用变量接收
	q, r := quotientAndRemainder(17, 5)
	fmt.Println("17 / 5 的商:", q, "余数:", r)

	// ===== 4. 变长参数（可变参数）=====
	fmt.Println("求和(1,2,3) =", sum(1, 2, 3))
	fmt.Println("求和(无参) =", sum())
	nums := []int{4, 5, 6}
	fmt.Println("切片展开求和 =", sum(nums...)) // 用 ... 把切片展开传入

	// ===== 5. defer：延迟执行，常用于资源释放 =====
	demoDefer()

	// ===== 6. 函数是「一等公民」：可赋值给变量、当参数传 =====
	var op func(int, int) int = add // 把函数赋给变量
	fmt.Println("通过变量调用 add:", op(7, 8))

	// 把函数当参数传进去（高阶函数）
	fmt.Println("apply(add):", apply(add, 2, 3))

	// ===== 7. 匿名函数与闭包 =====
	// 匿名函数：没有名字，定义后可立即调用或赋值
	square := func(x int) int {
		return x * x
	}
	fmt.Println("5 的平方 =", square(5))

	// 闭包：函数「记住」了它外部的变量
	counter := makeCounter()
	fmt.Println("计数器:", counter(), counter(), counter()) // 1 2 3
}

// add 普通函数：两个 int 参数同类型时可合并写 (a, b int)
func add(a, b int) int {
	return a + b
}

// divide 多返回值：返回结果 + error
// Go 的惯例是把 error 作为最后一个返回值
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为 0")
	}
	return a / b, nil
}

// quotientAndRemainder 命名返回值：q 和 r 已声明，可直接赋值后裸 return
func quotientAndRemainder(a, b int) (q int, r int) {
	q = a / b
	r = a % b
	return // 裸 return：自动返回命名的 q 和 r
}

// sum 变长参数：nums 在函数内是一个 []int 切片
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// demoDefer 演示 defer 的执行时机与顺序
func demoDefer() {
	// defer 注册的语句会在函数「返回前」执行
	// 多个 defer 按「后进先出」(栈) 的顺序执行
	defer fmt.Println("defer 1（最后执行）")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3（最先执行）")
	fmt.Println("函数体正常执行")
}

// apply 高阶函数：接收一个函数作为参数
func apply(fn func(int, int) int, a, b int) int {
	return fn(a, b)
}

// makeCounter 返回一个闭包：每次调用 +1
// count 被内部的匿名函数「捕获」，在多次调用间保持状态
func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}
