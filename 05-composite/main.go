// 第 5 课：复合类型（数组 / 切片 / map / struct）
// 运行方式：go run ./05-composite
package main

import "fmt"

// ============================================================
// 【疑问 1】数组(array) 和 切片(slice) 的区别
// ------------------------------------------------------------
// 维度        | 数组 array            | 切片 slice
// -----------|----------------------|---------------------------
// 长度        | 固定，编译期定死         | 可变，可 append 动态扩容
// 长度是类型吗  | 是！[3]int≠[4]int    | 不是，[]int 就是 []int
// 声明        | var a [3]int         | var s []int (零值是 nil)
// 创建        | [3]int{1,2,3}        | []int{1,2,3} / make([]int,n)
// 传递语义     | 值传递(整体拷贝一份)     | 引用语义(共享底层数组)
// 常用度      | 少，长度固定才用         | 多，日常主力
//
// 关键记忆：
//   1) 数组方括号里「有数字」[3]int；切片「没数字」[]int —— 这是肉眼区分两者最快的方法
//   2) 数组是值：传给函数 / 赋值会「整体拷贝」，改副本不影响原数组
//      切片是引用：多个切片可能共享同一底层数组，改一个影响另一个（见下方 sub 演示）
//   3) 切片底层其实就是「指向某个数组的指针 + 长度 len + 容量 cap」
//
// ============================================================
// 【疑问 2】容器类（slice / map）的几种创建方式汇总
// ------------------------------------------------------------
// 切片 slice 的创建：
//   a) 字面量：       s := []int{1, 2, 3}
//   b) make 指定长度： s := make([]int, 3)        // len=3, cap=3, 元素是零值
//   c) make 指定容量： s := make([]int, 3, 10)    // len=3, cap=10
//   d) 声明为 nil：    var s []int                // nil 切片，可直接 append
//   e) 从数组/切片截取：s := arr[1:3]              // 共享底层数组
//
// map 的创建：
//   a) 字面量：       m := map[string]int{"a": 1}
//   b) make：        m := make(map[string]int)   // 推荐：创建空 map
//   c) 声明为 nil：    var m map[string]int       // nil map！只能读不能写，写会 panic
//
// ⚠️ 重点区别：
//   - nil 切片可以直接 append（安全）
//   - nil map 不能直接赋值写入（会 panic），必须先 make 或用字面量初始化
//
// ============================================================
// 【疑问 3】make vs new 的区别
// ------------------------------------------------------------
//                | make                         | new
// ---------------|-----------------------------|--------------------
// 返回类型       | 值本身 T                      | 指针 *T
// 能用于哪些类型  | 只能 slice / map / channel    | 任意类型
// 作用           | 初始化内部数据结构(底层数组等)   | 分配零值内存
// 举例           | make([]int,3) → []int        | new(int) → *int
//
// 为什么 make 不返回指针？
//   因为 slice / map / channel 本身就是「引用类型」：
//   - slice 底层 = {指向数组的指针, len, cap}
//   - map   底层 = 指向哈希表的指针
//   - chan   底层 = 指向队列的指针
//   它们传递时已经在「共享」，不需要再包一层指针。
//
// 实战选择：
//   slice / map / channel → 用 make
//   结构体想拿指针         → 用 &T{}（比 new 更常见，能带初始值）
//   new 在实战中其实很少用
// ============================================================

// User 结构体：把多个字段组合成一个自定义类型
// 字段首字母大写 = 导出（可被其他包访问）
type User struct {
	Name string
	Age  int
}

func main() {
	// ===== 1. 数组：固定长度，长度是类型的一部分 =====
	// [3]int 和 [4]int 是两种不同的类型
	var arr [3]int        // 零值：[0 0 0]
	arr[0] = 10           // 通过下标赋值
	nums := [3]int{1, 2, 3} // 声明并初始化
	fmt.Println("数组 arr:", arr, " nums:", nums, " 长度:", len(nums))

	// ===== 2. 切片 slice：可变长度，实际开发中远比数组常用 =====
	// 2.1 用字面量创建
	fruits := []string{"apple", "banana"}

	// 2.2 用 append 追加元素（切片会自动扩容）
	fruits = append(fruits, "cherry")
	fmt.Println("切片:", fruits, " 长度:", len(fruits), " 容量:", cap(fruits))

	// 2.3 切片操作 s[low:high]：左闭右开，含 low 不含 high
	sub := fruits[1:3] // 取下标 1、2
	fmt.Println("切片截取 [1:3]:", sub)

	// 2.4 用 make 创建指定长度/容量的切片
	scores := make([]int, 3, 10) // 长度 3，容量 10
	fmt.Println("make 切片:", scores, " len:", len(scores), " cap:", cap(scores))

	// ⚠️ 切片是「引用」语义：截取出的子切片和原切片共享底层数组
	sub[0] = "BANANA"               // 改子切片
	fmt.Println("改 sub 后原切片:", fruits) // 原切片也跟着变了！

	// ===== 3. map：键值对（哈希表），无序 =====
	// 3.1 用字面量创建
	ages := map[string]int{
		"Tom":  18,
		"Jerry": 20,
	}

	// 3.2 增 / 改：直接赋值
	ages["Spike"] = 25

	// 3.3 查：用「值, ok」双返回值判断 key 是否存在
	age, ok := ages["Tom"]
	fmt.Println("Tom 的年龄:", age, " 存在吗:", ok)
	_, exist := ages["Unknown"]
	fmt.Println("Unknown 存在吗:", exist)

	// 3.4 删：delete
	delete(ages, "Spike")

	// 3.5 遍历 map（注意：顺序是随机的）
	for name, a := range ages {
		fmt.Printf("map 遍历 -> %s: %d\n", name, a)
	}

	// ===== 4. struct：结构体，自定义复合类型 =====
	// 4.1 结构体初始化的几种方式
	// 方式 A：指定字段名（最推荐）—— 可只写部分字段，其余取零值，字段顺序随意
	u1 := User{Name: "Alice", Age: 30}

	// 方式 B：按字段顺序（不推荐）—— 必须写全所有字段且顺序要对，加字段时易错
	u2 := User{"Bob", 25}

	// 方式 C：只声明不初始化 —— 得到零值结构体 {"" 0}
	var u3 User

	// 方式 D：部分字段初始化 —— 没写的 Age 自动取零值 0
	u4 := User{Name: "Carol"}

	// 方式 E：取地址得到「结构体指针」—— u5 的类型是 *User
	u5 := &User{Name: "Dave", Age: 40}

	// 方式 F：用 new(T) —— 返回零值结构体的指针 *User，等价于 &User{}
	u6 := new(User)
	u6.Name = "Eve" // 通过指针赋值（. 自动解引用）

	fmt.Println("A 指定字段名:", u1)
	fmt.Println("B 按顺序:", u2)
	fmt.Println("C 零值:", u3)
	fmt.Println("D 部分字段:", u4)
	fmt.Println("E 指针字面量:", u5, " 解引用:", *u5)
	fmt.Println("F new(User):", u6, " 解引用:", *u6)

	// 4.2 访问 / 修改字段
	u1.Age = 31
	fmt.Println("修改后 u1.Age:", u1.Age)

	// 4.3 结构体指针：用 & 取地址，Go 会自动解引用访问字段
	p := &u1
	p.Name = "Alice 改名" // 等价于 (*p).Name，Go 自动帮你解引用
	fmt.Println("通过指针改名后:", u1)

	// ===== 5. 组合：切片 + 结构体（实际开发最常见的数据结构）=====
	users := []User{
		{Name: "张三", Age: 28},
		{Name: "李四", Age: 33},
	}
	for i, user := range users {
		fmt.Printf("用户列表[%d]: 姓名=%s 年龄=%d\n", i, user.Name, user.Age)
	}

	// ===== 6. 【疑问1 演示】数组「值传递」vs 切片「引用语义」=====
	demoArray := [3]int{1, 2, 3}
	demoSlice := []int{1, 2, 3}
	modify(demoArray, demoSlice) // 把数组和切片都传进函数里改
	fmt.Println("传入函数后 -> 数组(没变):", demoArray, " 切片(被改了):", demoSlice)

	// ===== 7. 【疑问2 演示】nil 切片可 append，nil map 不可写 =====
	var nilSlice []int           // nil 切片
	nilSlice = append(nilSlice, 1, 2) // ✅ 安全：可直接 append
	fmt.Println("nil 切片 append 后:", nilSlice)

	// var nilMap map[string]int  // nil map
	// nilMap["x"] = 1            // ❌ 取消注释会 panic: assignment to entry in nil map
	safeMap := make(map[string]int) // ✅ 必须先 make
	safeMap["x"] = 1
	fmt.Println("make 后的 map 可写:", safeMap)

	// ===== 8. 【字面量初始化】右边 {} 必须写，但里面可以为空 =====
	// 字面量格式是「类型{...}」：{} 不能省，但里面给不给元素都行
	emptyLiteralSlice := []int{}          // ✅ 空切片「字面量」，合法
	emptyLiteralMap := map[string]int{}   // ✅ 空 map「字面量」，可直接写入
	emptyStruct := User{}                 // ✅ 空结构体，所有字段取零值
	fmt.Println("空字面量 -> slice:", emptyLiteralSlice, " map:", emptyLiteralMap, " struct:", emptyStruct)

	// 8.1 三种「空切片」的细微差别：nil 切片 vs 空字面量 vs make
	var sNil []int          // nil 切片
	sLiteral := []int{}     // 空字面量
	sMake := make([]int, 0) // make 创建
	// 三者 len 都是 0、都能 append；唯一区别：sNil == nil 为 true，另两个为 false
	fmt.Printf("nil切片:   len=%d, 是nil吗=%t\n", len(sNil), sNil == nil)
	fmt.Printf("空字面量:  len=%d, 是nil吗=%t\n", len(sLiteral), sLiteral == nil)
	fmt.Printf("make切片:  len=%d, 是nil吗=%t\n", len(sMake), sMake == nil)

	// 8.2 核心结论：
	// - 字面量的 {} 必须写，但里面可空 → []int{}、User{} 都合法
	// - 想完全不写 {} → 用 var 声明，拿到「零值/nil」
	// - 大坑：var map 是 nil map 不能写；var slice 是 nil 切片却能 append

	// ===== 9. 切片扩容演示 =====
	// 触发条件：append 时 len 达到 cap → 分配新底层数组，容量翻倍（小于256时）
	// 注意：append 必须写成 s = append(s, x)，扩容后返回的可能是全新的切片头

	// 9.1 观察 cap 的变化过程
	var grow []int
	prevCap := cap(grow)
	fmt.Println("扩容过程:")
	for i := 0; i < 20; i++ {
		grow = append(grow, i)
		if cap(grow) != prevCap {
			fmt.Printf("  len=%2d → cap 从 %d 扩到 %d\n", len(grow), prevCap, cap(grow))
			prevCap = cap(grow)
		}
	}

	// 9.2 扩容后旧引用失效的演示
	original := make([]int, 2, 2) // len=2, cap=2，已满
	original[0], original[1] = 10, 20
	alias := original // alias 和 original 共享底层数组

	alias[0] = 99                         // 没扩容，original 也变了
	fmt.Println("扩容前共享: original:", original, "alias:", alias)

	alias = append(alias, 30)             // len > cap，触发扩容！alias 指向新数组
	alias[0] = 777                        // 改 alias 不再影响 original
	fmt.Println("扩容后脱钩: original:", original, "alias:", alias)
	// 结论：扩容后 alias 和 original 不再共享，各改各的
}

// modify 演示传参语义：数组是「值传递」改不动原值，切片是「引用」能改动原值
func modify(arr [3]int, s []int) {
	arr[0] = 999 // 改的是数组副本，外面看不到
	s[0] = 999   // 改的是共享的底层数组，外面看得到
}
