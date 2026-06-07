// 第 14 课：数据结构与算法常用工具（刷题必备）
// 运行方式：go run ./14-datastructures
//
// Go 刷算法和 Java 最大的区别：没有 Collections 框架，
// 大部分数据结构用 slice + 标准库组合实现。
// 本课把刷题最常用的数据结构和操作全部列出来，方便速查。
//
// ============================================================
// 【总结】Go vs Java 数据结构对照表
// ------------------------------------------------------------
//   Java                    Go
//   ArrayList<E>            []T（slice）
//   LinkedList<E>           container/list
//   Stack<E>                []T（slice 模拟）
//   Queue<E>                []T（slice 模拟）或 channel
//   Deque<E>                []T（slice 模拟）或 container/list
//   PriorityQueue<E>        container/heap（需实现接口）
//   HashSet<E>              map[T]bool 或 map[T]struct{}
//   HashMap<K,V>            map[K]V
//   TreeMap<K,V>            无内置，需排序 slice 或第三方库
//   Arrays.sort()           sort.Slice() / slices.Sort()
//   Collections.reverse()   slices.Reverse()
//   Math.max/min            max() / min()（Go 1.21+ 内置）
//   Integer.MAX_VALUE       math.MaxInt
//   StringBuilder           strings.Builder
//
// 【刷题常用 import】
//   "fmt"              → 输入输出
//   "sort"             → 排序
//   "math"             → MaxInt, MinInt, Abs
//   "strings"          → 字符串操作
//   "strconv"          → 字符串 ↔ 数字转换
//   "container/heap"   → 优先级队列
//   "container/list"   → 双向链表
//
// ============================================================
// 【深入理解】Slice 传值原理 & heap 为什么用指针
// ------------------------------------------------------------
//
// ▶ Slice 的本质是一个 struct（描述符），包含 3 个字段：
//     type slice struct {
//         array *底层数组   // 指针 → 传副本时两份共享同一个底层数组
//         len   int        // 值  → 传副本时各自一份，互不影响
//         cap   int        // 值  → 传副本时各自一份，互不影响
//     }
//
// ▶ 所以 slice 方法的接收者选择规则：
//     操作                     改了什么         值接收者行吗？
//     h[i] = x                底层数组元素      ✅ 行（共享指针）
//     h[i], h[j] = h[j], h[i] 底层数组元素      ✅ 行（Swap 不需要指针）
//     h = append(h, x)        len/cap/可能换指针 ❌ 不行（改的是副本的 len）
//     h = h[:n-1]             len              ❌ 不行（改的是副本的 len）
//
//   总结：改元素不用指针，改长度必须用指针。
//
// ▶ heap 为什么声明时要用 &IntHeap{} 而不是 IntHeap{}？
//   1. Push/Pop 要改 slice 的 len → 必须用指针接收者 func (h *IntHeap) Push/Pop
//   2. 指针接收者实现的接口 → 只有指针能赋给接口（第 6 课规则）
//   3. heap.Push(h, val) 的参数 h 是 heap.Interface 接口
//   4. 所以必须传 *IntHeap 指针 → 声明时 minHeap := &IntHeap{}
//
//   而 Len/Less/Swap 只读或只改数组元素，用值接收者就行，
//   但这不影响——指针的方法集包含值接收者的方法，所以 *IntHeap 也有它们。
//
// ▶ 同一个 struct 混用值接收者和指针接收者（heap 就是这样！）：
//   func (h IntHeap) Len() int           // 值接收者：只读
//   func (h IntHeap) Less(i, j int) bool // 值接收者：只读
//   func (h IntHeap) Swap(i, j int)      // 值接收者：改数组元素（共享指针）
//   func (h *IntHeap) Push(x any)        // 指针接收者：要改 len
//   func (h *IntHeap) Pop() any          // 指针接收者：要改 len
//
//   混用完全合法，实际项目中很常见：只读方法用值接收者，修改方法用指针接收者。
//
//   方法集规则：
//     值类型 IntHeap  → 方法集只有：Len、Less、Swap（3 个）
//     指针类型 *IntHeap → 方法集包含全部：Len、Less、Swap + Push、Pop（5 个）
//
//   所以赋给接口时必须用指针（&IntHeap{}），因为接口要求 5 个方法，
//   值类型只有 3 个，不满足；指针类型有全部 5 个，满足。
// ============================================================
package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"math"
	"sort"
	"strings"
)

func main() {
	// ===== 1. Slice 当动态数组（ArrayList）=====
	fmt.Println("===== 1. 动态数组 (slice) =====")
	arr := []int{3, 1, 4, 1, 5, 9}
	arr = append(arr, 2, 6)                      // 尾部追加
	fmt.Println("原始:", arr)

	// 删除索引 i 的元素（不保序，O(1)）
	i := 2 // 删除 arr[2]
	arr[i] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	fmt.Println("删除索引2(不保序):", arr)

	// 删除索引 i 的元素（保序，O(n)）
	arr2 := []int{10, 20, 30, 40, 50}
	i = 2
	arr2 = append(arr2[:i], arr2[i+1:]...)
	fmt.Println("删除索引2(保序):", arr2) // [10 20 40 50]

	// 插入元素到索引 i
	arr3 := []int{1, 2, 4, 5}
	i = 2
	arr3 = append(arr3[:i+1], arr3[i:]...)
	arr3[i] = 3 // 在索引 2 插入 3
	fmt.Println("在索引2插入3:", arr3) // [1 2 3 4 5]

	// 复制 slice（避免共享底层数组）
	copied := make([]int, len(arr3))
	copy(copied, arr3)
	fmt.Println("复制:", copied)

	// ===== 2. 排序 =====
	// sort.Slice 底层用的是 pdqsort（Pattern-Defeating Quicksort，Go 1.19+）
	// 它是一个混合排序算法，会根据数据特征自动切换：
	//   数据量小（≤12）       → 插入排序（常数小，小规模更快）
	//   正常情况              → 快速排序（平均 O(n log n)，缓存友好）
	//   检测到快排退化        → 堆排序（保证最坏 O(n log n)，不退化到 O(n²)）
	//   数据接近有序          → 识别 pattern 加速（pdqsort 的特色）
	//
	// ⚠️ sort.Slice 是「不稳定排序」，相等元素顺序可能变
	//    需要稳定排序用 sort.SliceStable
	//
	// 对比 Java：
	//   Go sort.Slice           → pdqsort（快排+堆排+插入排序）
	//   Java Arrays.sort(int[]) → Dual-Pivot Quicksort（双轴快排）
	//   Java Arrays.sort(Object[]) → TimSort（归并+插入排序，稳定）
	fmt.Println("\n===== 2. 排序 =====")
	nums := []int{5, 3, 8, 1, 9, 2}

	// 升序
	sort.Ints(nums)
	fmt.Println("升序:", nums)

	// 降序
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	fmt.Println("降序:", nums)

	// 自定义排序（sort.Slice 最常用）
	people := []struct {
		Name string
		Age  int
	}{
		{"Tom", 25}, {"Jerry", 20}, {"Alice", 30},
	}
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age // 按年龄升序
	})
	fmt.Println("自定义排序:", people)

	// 二分查找（要求已排序）
	sorted := []int{1, 3, 5, 7, 9, 11}
	idx := sort.SearchInts(sorted, 7) // 找 7 的位置
	fmt.Printf("二分查找 7 的位置: %d (值=%d)\n", idx, sorted[idx])

	// ===== 3. Stack（用 slice 模拟）=====
	fmt.Println("\n===== 3. 栈 (Stack) =====")
	stack := []int{}
	stack = append(stack, 1, 2, 3) // push
	fmt.Println("push 1,2,3:", stack)

	top := stack[len(stack)-1]     // peek（看栈顶）
	stack = stack[:len(stack)-1]   // pop（弹出栈顶）
	fmt.Println("pop:", top, " 剩余:", stack)

	isEmpty := len(stack) == 0 // 判空
	fmt.Println("isEmpty:", isEmpty)

	// ===== 4. Queue（用 slice 模拟）=====
	fmt.Println("\n===== 4. 队列 (Queue) =====")
	queue := []int{}
	queue = append(queue, 1, 2, 3) // enqueue（入队）
	fmt.Println("enqueue 1,2,3:", queue)

	front := queue[0]   // peek（看队首）
	queue = queue[1:]    // dequeue（出队）
	fmt.Println("dequeue:", front, " 剩余:", queue)
	// ⚠️ slice 模拟队列频繁出队会导致底层数组不释放
	// 刷题够用，生产环境建议用 container/list 或 channel

	// ===== 5. HashSet（用 map 模拟）=====
	fmt.Println("\n===== 5. HashSet =====")
	// 方式一：map[T]bool（直观）
	set := map[int]bool{}
	set[1] = true
	set[2] = true
	set[1] = true // 重复添加无效
	fmt.Println("contains 1:", set[1])
	fmt.Println("contains 3:", set[3]) // false（零值）
	delete(set, 1)                     // 删除
	fmt.Println("删除1后 contains 1:", set[1])
	fmt.Println("size:", len(set))

	// 方式二：map[T]struct{}（更省内存，struct{} 占 0 字节）
	set2 := map[string]struct{}{}
	set2["go"] = struct{}{}
	if _, exists := set2["go"]; exists {
		fmt.Println("set2 contains 'go'")
	}

	// ===== 6. HashMap（map）=====
	fmt.Println("\n===== 6. HashMap (map) =====")
	counter := map[string]int{} // 计数器模式（刷题超常用）
	words := []string{"go", "java", "go", "python", "go"}
	for _, w := range words {
		counter[w]++ // 不存在时零值是 0，直接 ++ 就行
	}
	fmt.Println("词频统计:", counter)

	// 遍历 map（顺序不固定！）
	for key, count := range counter {
		fmt.Printf("  %s: %d 次\n", key, count)
	}

	// 检查 key 是否存在（两种写法）
	val, ok := counter["rust"]
	fmt.Printf("rust 存在: %v, 值: %d\n", ok, val)

	// ===== 7. 双向链表（container/list）=====
	fmt.Println("\n===== 7. 双向链表 (list) =====")
	l := list.New()
	l.PushBack(1)             // 尾部插入
	l.PushBack(2)
	l.PushFront(0)            // 头部插入
	fmt.Print("链表: ")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()
	l.Remove(l.Front())       // 删除头部
	fmt.Println("删除头部后, 头部:", l.Front().Value)
	fmt.Println("链表长度:", l.Len())
	// 常见用途：LRU 缓存（map + 双向链表）

	// ===== 8. 优先级队列 / 最小堆（container/heap）=====
	fmt.Println("\n===== 8. 优先级队列 (heap) =====")
	// 需要自定义类型并实现 heap.Interface（5 个方法）
	// 见文件底部的 IntHeap 实现

	// 最小堆（默认）
	minHeap := &IntHeap{}
	heap.Init(minHeap)
	heap.Push(minHeap, 5)
	heap.Push(minHeap, 3)
	heap.Push(minHeap, 8)
	heap.Push(minHeap, 1)
	fmt.Printf("最小堆 peek: %d\n", (*minHeap)[0]) // 看堆顶（不弹出）
	for minHeap.Len() > 0 {
		fmt.Printf("  pop: %d\n", heap.Pop(minHeap))  // 按从小到大弹出
	}

	// 最大堆：把 Less 反过来即可（见底部 MaxHeap）
	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)
	heap.Push(maxHeap, 5)
	heap.Push(maxHeap, 3)
	heap.Push(maxHeap, 8)
	heap.Push(maxHeap, 1)
	fmt.Printf("最大堆 peek: %d\n", (*maxHeap)[0])
	for maxHeap.Len() > 0 {
		fmt.Printf("  pop: %d\n", heap.Pop(maxHeap)) // 按从大到小弹出
	}

	// ===== 9. 常用数学函数 =====
	fmt.Println("\n===== 9. 数学工具 =====")
	fmt.Println("MaxInt:", math.MaxInt)     // 等价于 Java Integer.MAX_VALUE
	fmt.Println("MinInt:", math.MinInt)
	fmt.Println("max(3,7):", max(3, 7))     // Go 1.21+ 内置 max/min
	fmt.Println("min(3,7):", min(3, 7))
	fmt.Println("abs(-5):", abs(-5))         // Go 没有 int 版 abs，自己写
	fmt.Println("Pow(2,10):", int(math.Pow(2, 10))) // 1024

	// ===== 10. 字符串常用操作 =====
	fmt.Println("\n===== 10. 字符串工具 =====")
	s := "hello, world"
	fmt.Println("包含:", strings.Contains(s, "world"))
	fmt.Println("分割:", strings.Split(s, ", "))
	fmt.Println("替换:", strings.ReplaceAll(s, "l", "L"))
	fmt.Println("大写:", strings.ToUpper(s))

	// strings.Builder（高效拼接，等价 Java StringBuilder）
	var builder strings.Builder
	for i := 0; i < 5; i++ {
		fmt.Fprintf(&builder, "%d ", i)
	}
	fmt.Println("Builder:", builder.String())

	// 字符串 ↔ rune 切片（处理中文/Unicode 必须用 rune）
	cnStr := "你好Go"
	runes := []rune(cnStr)
	fmt.Println("rune 长度:", len(runes))    // 4（字符数）
	fmt.Println("byte 长度:", len(cnStr))     // 8（字节数）
	fmt.Println("反转:", string(reverseRunes(runes)))

	// ===== 11. 位运算（刷题常用）=====
	fmt.Println("\n===== 11. 位运算 =====")
	a, b := 6, 3 // 110, 011
	fmt.Printf("AND  %b & %b = %b (%d)\n", a, b, a&b, a&b)   // 010 = 2
	fmt.Printf("OR   %b | %b = %b (%d)\n", a, b, a|b, a|b)   // 111 = 7
	fmt.Printf("XOR  %b ^ %b = %b (%d)\n", a, b, a^b, a^b)   // 101 = 5
	fmt.Printf("左移 %b << 1 = %b (%d)\n", a, a<<1, a<<1)     // 1100 = 12
	fmt.Printf("右移 %b >> 1 = %b (%d)\n", a, a>>1, a>>1)     // 11 = 3

	// 常用技巧
	n := 12
	fmt.Println("是否偶数 (n&1==0):", n&1 == 0)         // true
	fmt.Println("除以2 (n>>1):", n>>1)                    // 6
	fmt.Println("乘以2 (n<<1):", n<<1)                    // 24

	// XOR 交换两个数（不用临时变量）
	x, y := 10, 20
	x, y = x^y, x^y^(x^y) // 但 Go 有多赋值，直接 x, y = y, x 更清晰
	fmt.Println("Go 多赋值交换: 直接 x, y = y, x")
}

// ===== 辅助函数 =====

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func reverseRunes(runes []rune) []rune {
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return runes
}

// ===== 最小堆（heap.Interface 实现）=====
// container/heap 要求实现 5 个方法：Len, Less, Swap, Push, Pop
// 这套代码可以当模板直接用！

// IntHeap 最小堆
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] } // < 是最小堆
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}

// MaxHeap 最大堆：只需要把 Less 反过来
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] } // > 是最大堆
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[:n-1]
	return val
}
