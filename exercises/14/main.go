package main

import (
	"container/heap"
	"fmt"
)

/*
## 第 14 课：数据结构

### 题 14.1 用栈判断括号匹配（⭐⭐）

写一个函数 `isValid(s string) bool`，判断字符串中的括号是否匹配。
支持 `()`, `[]`, `{}`。
示例：`"([{}])"` → true，`"([)]"` → false

### 题 14.2 TopK 问题（⭐⭐⭐）

给定一个整数 slice，用 `container/heap` 实现找出最大的 K 个元素。
要求用最小堆，堆大小保持 K。

### 踩坑总结

 1. **栈为空时不能 break 退出循环**：
    原来写了 `if len(stack) == 0 { break }`，导致第一个字符进来时栈为空直接退出，整个循环等于没执行。
    正确做法：栈为空的判断应该放在遇到"右括号"的分支里，表示右括号多了 → return false。
    左括号入栈不需要关心栈是否为空。
*/
func main() {
	fmt.Println("===== 14.1 括号匹配 =====")
	fmt.Println("([{}]):", isValid("([{}])")) // true
	fmt.Println("([)]:", isValid("([)]"))     // false
	fmt.Println("()[]{}:", isValid("()[]{}")) // true
	fmt.Println("(:", isValid("("))           // false
	fmt.Println("):", isValid(")"))           // false
	fmt.Println(":", isValid(""))             // false

	fmt.Println("\n===== 14.2 TopK =====")
	nums := []int{3, 1, 5, 12, 2, 11, 9, 7, 4}
	fmt.Println("原始数据:", nums)
	fmt.Println("Top 3:", topK(nums, 3)) // [12, 11, 9]
	fmt.Println("Top 5:", topK(nums, 5)) // [12, 11, 9, 7, 5]
	fmt.Println("Top 1:", topK(nums, 1)) // [12]
}

func topK(ss []int, k int) []int {
	queue := &MaxIntPriorityQueue{}
	for _, s := range ss {
		heap.Push(queue, s)
	}
	rst := []int{}
	for i := 0; i < k; i++ {
		pop := heap.Pop(queue)
		rst = append(rst, pop.(int))
	}
	return rst
}

type MaxIntPriorityQueue []int

func (m MaxIntPriorityQueue) Len() int {
	return len(m)
}

func (m MaxIntPriorityQueue) Less(i, j int) bool {
	return m[i] > m[j]
}

func (m MaxIntPriorityQueue) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m *MaxIntPriorityQueue) Push(x any) {
	*m = append(*m, x.(int))
}

func (m *MaxIntPriorityQueue) Pop() any {
	old := *m
	l := len(old)
	val := old[l-1]
	*m = old[:l-1]
	return val
}

func isValid(s string) bool {
	if s == "" {
		return false
	}
	mapping := map[byte]byte{
		']': '[',
		')': '(',
		'}': '{',
	}
	stack := []byte{}
	bytes := []byte(s)
	for _, by := range bytes {
		if by == '{' || by == '[' || by == '(' {
			stack = append(stack, by)
		} else {
			// 栈为空说明右括号多了，不匹配
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if mapping[by] != top {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	if len(stack) != 0 {
		return false
	}
	return true
}
