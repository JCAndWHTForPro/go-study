// 第 10 课：泛型（Go 1.18+）
// 运行方式：go run ./10-generics
//
// 泛型解决什么问题？
//   没有泛型时，想写一个「对 int 和 float64 都适用的求最大值函数」，
//   要么写两个函数（maxInt / maxFloat），要么用 any + 类型断言（丑且不安全）。
//   泛型让你只写一次，编译器自动适配多种类型，且保持类型安全。
//
// ============================================================
// 【知识点】Go 的比较规则 & 与泛型约束的关系
// ------------------------------------------------------------
//
// 一、相等比较 == !=（哪些类型能比）
//   ✅ 基本类型：int / float / bool / string / byte / rune
//   ✅ 指针 *T：比较的是地址是否相同
//   ✅ channel：是否是同一个 channel
//   ✅ 接口 interface：类型和值都相同才 true
//   ✅ 数组 [3]int：逐元素比较（前提：元素类型可比较）
//   ✅ 结构体 struct：逐字段比较（前提：所有字段可比较）
//   ❌ 切片 []T：只能和 nil 比，不能互相比
//   ❌ map：只能和 nil 比
//   ❌ 函数 func：只能和 nil 比
//
// 二、排序比较 > < >= <=（只有极少类型能排序）
//   ✅ 整数（int/int8/int16/int32/int64/uint...）
//   ✅ 浮点数（float32/float64）
//   ✅ 字符串 string（按字典序）
//   ❌ 其他所有类型（bool/struct/数组/指针/接口...）
//
// 三、与泛型约束的对应关系
//   约束             能做什么             对应哪些类型
//   any             啥都不能比            所有类型
//   comparable      能 == !=            基本类型+指针+channel+数组+结构体
//   Ordered(自定义)  能 > < >= <=        整数+浮点+字符串
//   Number(自定义)   能 + - * /          整数+浮点
//
// 四、易踩的坑
//   1) 切片/map 不能 ==：想比较内容要用 slices.Equal (Go 1.21+) 或逐个比
//   2) 含不可比较字段的结构体也不能 ==：
//      type A struct { Data []int } → A{} == A{} 编译报错（Data 是切片）
//      type B struct { Name string } → B{} == B{} ✅（所有字段可比较）
//   3) 能 == 的类型远多于能 > 的类型：
//      结构体/数组可以 == 但不能 >（所以 comparable ≠ Ordered）
// ============================================================
package main

import "fmt"

func main() {
	// ===== 1. 泛型函数：类型参数 [T] =====
	// Max[T] 可以比较任意「可排序」类型（int / float64 / string 等）
	fmt.Println("Max int:", Max(3, 7))
	fmt.Println("Max float:", Max(3.14, 2.71))
	fmt.Println("Max string:", Max("apple", "banana")) // 字符串按字典序比

	// 大多数时候 Go 能自动推断类型参数，不需要显式写 [int]
	// 但你也可以显式指定：
	fmt.Println("显式指定类型:", Max[int](10, 20))

	// ===== 2. 类型约束（constraint）=====
	// comparable：内置约束，表示类型支持 == 和 !=
	fmt.Println("Contains int:", Contains([]int{1, 2, 3}, 2))
	fmt.Println("Contains string:", Contains([]string{"go", "rust"}, "java"))

	// ===== 3. 自定义约束 =====
	// Number 约束：只允许数字类型
	fmt.Println("Sum int:", Sum([]int{1, 2, 3, 4}))
	fmt.Println("Sum float:", Sum([]float64{1.1, 2.2, 3.3}))

	// ===== 4. 泛型结构体 =====
	// 一个通用的「栈」数据结构，支持任意类型
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	fmt.Println("栈 Pop:", intStack.Pop(), intStack.Pop())
	fmt.Println("栈剩余大小:", intStack.Size())

	strStack := &Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	fmt.Println("字符串栈 Pop:", strStack.Pop())

	// ===== 5. 泛型 Map / Filter 函数 =====
	nums := []int{1, 2, 3, 4, 5}

	// Map：对每个元素应用函数，返回新切片
	doubled := Map(nums, func(n int) int { return n * 2 })
	fmt.Println("Map 翻倍:", doubled)

	// Filter：筛选满足条件的元素
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("Filter 偶数:", evens)

	// Map 还可以转换类型：int → string
	strs := Map(nums, func(n int) string { return fmt.Sprintf("第%d个", n) })
	fmt.Println("Map 转字符串:", strs)
}

// ============================================================
// 1. 泛型函数
// ============================================================

// Max 返回两个值中的较大值
// [T Ordered] 表示 T 必须满足 Ordered 约束（支持 < > 比较）
func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ============================================================
// 2. 内置约束 comparable
// ============================================================

// Contains 判断切片中是否包含某个元素
// [T comparable] 表示 T 支持 == 比较
func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

// ============================================================
// 3. 自定义约束
// ============================================================

// Ordered 约束：支持 < > 比较的类型（Go 1.21 后标准库有 cmp.Ordered，这里自定义演示）
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

// Number 约束：只允许数字类型
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum 对数字切片求和
func Sum[T Number](nums []T) T {
	var total T // 零值：0
	for _, n := range nums {
		total += n
	}
	return total
}

// ============================================================
// 4. 泛型结构体
// ============================================================

// Stack 通用栈，支持任意类型
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		var zero T // 返回零值
		return zero
	}
	top := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return top
}

func (s *Stack[T]) Size() int {
	return len(s.items)
}

// ============================================================
// 5. 泛型 Map / Filter 函数
// ============================================================

// Map 对切片每个元素应用函数 fn，返回新切片（支持类型转换）
func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// Filter 筛选满足条件的元素
func Filter[T any](slice []T, fn func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
