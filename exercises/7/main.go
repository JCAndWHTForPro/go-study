package main

import "fmt"

/*
## 第 7 课：指针

### 知识点总结

1. 指针基础
   - *T 表示指向 T 类型的指针，&x 取变量 x 的地址
   - 通过指针可以在函数内修改外部变量的值（如 swap 函数）

2. 指针的指针 **T
   - 当需要在函数内修改指针本身的值时，需要传 **T
   - 典型场景：push 函数要修改 head 指针指向新节点，所以参数是 **Node

3. 链表反转（三指针法）
   - 用 prev、current、next 三个指针
   - 每轮循环：保存 next → 反转指向 → prev/current 前进
   - 循环结束后 prev 就是新的头节点

4. 合并有序链表
   - 方法一（插入法）：以一条链为主链，把另一条的节点插入合适位置
   - 方法二（双指针法）：用 dummy 哨兵节点 + current 指针，每次比较两条链当前节点，小的接到 current 后面
   - 双指针法逻辑更对称、更不容易出错，推荐优先掌握

5. 哨兵节点（dummy node）
   - 创建一个虚拟头节点，避免对"头节点为空"等边界情况的特殊处理
   - 最终返回 dummy.Next 即为真正的头节点

### 题 7.1 交换函数（⭐）

写一个函数 `swap(a, b *int)` 通过指针交换两个整数的值。

### 题 7.2 链表（⭐⭐⭐）

定义一个单链表节点 `type Node struct { Val int; Next *Node }`：

1. 写 `push(head **Node, val int)` 在头部插入
2. 写 `printList(head *Node)` 打印整个链表
3. 写 `reverse(head *Node) *Node` 反转链表
*/
type Node struct {
	Val  int
	Next *Node
}

func (node *Node) push(head **Node, val int) {
	newNode := &Node{Val: val, Next: *head}
	*head = newNode
}
func (node *Node) printList(head *Node) {
	for ; head != nil; head = head.Next {
		fmt.Print(head.Val, " ")
	}
	fmt.Println()
}
func (node *Node) reverse(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}
	current := head
	var pre *Node
	for current != nil {
		next := current.Next
		current.Next = pre
		pre = current
		current = next
	}
	return pre
}

/*
### 题 7.3 合并两个有序链表（⭐⭐⭐）

给定两个升序链表，合并为一个升序链表并返回。
示例：

	链表1: 1 → 3 → 5
	链表2: 2 → 4 → 6
	结果:  1 → 2 → 3 → 4 → 5 → 6

提示：用一个哨兵节点（dummy）简化头节点处理
*/
func mergeTwoLists(l1, l2 *Node) *Node {
	dummy := &Node{Val: 0, Next: l1}
	pre := dummy
	for l1 != nil {
		if l2 != nil && l1.Val >= l2.Val {
			next := l2.Next
			l2.Next = l1
			pre.Next = l2
			pre = l2
			l2 = next
		} else {
			pre = l1
			l1 = l1.Next
		}
	}
	if l2 != nil {
		pre.Next = l2
	}
	return dummy.Next
}

func main() {
	// 测试反转
	var head *Node
	head.push(&head, 1)
	head.push(&head, 2)
	head.push(&head, 3)
	head.printList(head)
	head = head.reverse(head)
	head.printList(head)

	// 测试合并有序链表
	var l1 *Node
	l1.push(&l1, 5)
	l1.push(&l1, 3)
	l1.push(&l1, 1)
	fmt.Print("链表1: ")
	l1.printList(l1)

	var l2 *Node
	l2.push(&l2, 6)
	l2.push(&l2, 4)
	l2.push(&l2, 2)
	fmt.Print("链表2: ")
	l2.printList(l2)

	merged := mergeTwoLists(l1, l2)
	fmt.Print("合并后: ")
	merged.printList(merged)
}

func swap(a, b *int) {
	*a, *b = *b, *a
}
