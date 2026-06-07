// 第 7 课：指针（& 取地址 / * 解引用 / nil / 指针 vs 值 / 为什么用指针）
// 运行方式：go run ./07-pointer
//
// Go 的指针比 C 简单很多：没有指针运算（不能 p++），不能随意转换，
// 但保留了「直接操作内存地址」的能力，核心就两个符号：
//   &x  取变量 x 的「地址」          → 得到一个指针
//   *p  取指针 p 指向的「值」(解引用)  → 拿回原值
//
// ============================================================
// 【总结】Go 中指针与接口的几种引用方式
// ------------------------------------------------------------
//
// 方式 1：裸指针（最基础）
//   p := &x         // 取地址得到 *int 指针
//   *p = 20         // 通过指针改值（需要解引用）
//   fmt.Println(*p) // 读值也要解引用
//
// 方式 2：结构体指针 + 自动解引用（语法糖）
//   up := &User{Name:"Tom"}  // *User 指针
//   up.Name = "Jerry"        // Go 自动解引用，等价于 (*up).Name
//
// 方式 3：指针接收者方法（方法绑在指针上）
//   func (r *Rectangle) Scale(f float64) { r.Width *= f }
//   rect.Scale(2)    // Go 自动取地址：等价于 (&rect).Scale(2)
//
// 方式 4：指针装进接口（最常见、最重要！）
//   var err error = &BalanceError{...}   // *BalanceError → error 接口
//   var s Shape = &Circle{Radius:5}      // *Circle → Shape 接口
//   ⚠️ 调用者拿到的是接口类型，不是裸指针
//   ⚠️ 不需要解引用，直接调接口方法即可：err.Error(), s.Area()
//   ⚠️ fmt.Println(err) 会自动调 Error() 方法，不会打印地址
//
// 方式 5：引用类型（天生含指针，不需要 & 取地址）
//   slice / map / channel 内部已含指针，传递时自动共享
//   s := make([]int, 3)   // s 不是指针，但底层含指针
//   m := make(map[string]int) // m 不是指针，但底层是指针
//
// 关键对比：
//   裸指针 *T      → 需要解引用 *p 才能拿到值
//   接口里的指针    → 不需要解引用，调方法即可（接口帮你封装了）
//   引用类型        → 不需要 & 也不需要 *，天生共享
//
// ============================================================
// 【总结】函数参数 vs 方法接收者：传值/传指针规则完全不同！
// ------------------------------------------------------------
//
// ▶ 普通函数参数：类型必须严格匹配，不会自动转换
//   func greet(u User)   → 只能传 User 值，传 &u 编译报错
//   func update(u *User) → 只能传 *User 指针，传 u 编译报错
//
// ▶ 方法接收者：Go 编译器会自动帮你加 & 或 *
//   u.Hello()   // 值调值方法 ✅
//   u.Update()  // 值调指针方法 ✅ → Go 自动取地址 (&u).Update()
//   p.Hello()   // 指针调值方法 ✅ → Go 自动解引用 (*p).Hello()
//   p.Update()  // 指针调指针方法 ✅
//
// 口诀：普通函数严格匹配，方法调用随便来。
// ============================================================
//
// 接口底层结构（iface）：
//   接口变量 = { 类型信息指针, 数据指针 }
//   所以接口既知道"里面装的是什么类型"（运行时反射/断言），
//   又能通过数据指针找到实际数据——这是 Go 多态的实现基础。
// ============================================================
package main

import "fmt"

type User struct {
	Name string
	Age  int
}

func main() {
	// ===== 1. & 取地址，* 解引用 =====
	x := 10
	p := &x // p 是指向 x 的指针，类型是 *int

	fmt.Println("x 的值:", x)
	fmt.Println("x 的地址 &x:", p)  // 打印出来是一个内存地址，如 0xc0000...
	fmt.Println("*p 解引用拿回值:", *p) // 通过指针拿回 x 的值 = 10

	// 通过指针「改」原变量：改 *p 就等于改 x 本身
	*p = 20
	fmt.Println("通过 *p=20 修改后, x =", x) // x 变成了 20

	// ===== 2. 指针的零值是 nil =====
	var np *int // 只声明不赋值，零值是 nil
	fmt.Println("未初始化的指针:", np, " 是 nil 吗:", np == nil)
	// fmt.Println(*np) // ❌ 对 nil 指针解引用会 panic，使用前要判空

	// ===== 2.1 nil 指针的安全使用 =====
	// nil 指针能做的事：
	fmt.Println("打印 nil 指针:", np)           // ✅ 输出 <nil>
	fmt.Println("判空:", np == nil)             // ✅ 输出 true
	var np3 *int = np                           // ✅ 赋值给另一个指针变量
	fmt.Println("赋值给另一个指针:", np3 == nil) // ✅ true
	safeUsePointer(np)                          // ✅ 传给函数当参数（函数内要判空）

	// nil 指针不能做的事（取消注释会 panic）：
	// *np = 10        // ❌ panic: runtime error: invalid memory address (写)
	// fmt.Println(*np)// ❌ panic: 对 nil 解引用 (读)
	// var u *User; u.Name // ❌ panic: 结构体指针为 nil 也不能访问字段

	// ===== 2.2 标准套路：使用前先判空 =====
	var userPtr *User // nil
	// 不判空直接用 → panic；判空了再用 → 安全
	if userPtr != nil {
		fmt.Println("用户名:", userPtr.Name) // 不会执行
	} else {
		fmt.Println("userPtr 是 nil，不能访问字段")
	}
	// 给它赋值后就能正常使用了
	userPtr = &User{Name: "Tom", Age: 18}
	if userPtr != nil {
		fmt.Println("赋值后可用:", userPtr.Name, userPtr.Age)
	}

	// ===== 3. 指针 vs 值：传参时的本质区别 =====
	a := 5
	addByValue(a)   // 传「值」：函数内改副本，外面不变
	fmt.Println("传值后 a =", a) // 仍是 5

	addByPointer(&a) // 传「指针」：函数内通过指针改原值
	fmt.Println("传指针后 a =", a) // 变成了 105

	// ===== 4. 为什么用指针？两大核心理由 =====

	// 理由一：想在函数里「修改」调用者的变量
	u := User{Name: "Tom", Age: 18}
	growUp(&u) // 传地址进去，函数才能改到原对象
	fmt.Println("长大后:", u) // Age 变成 19

	// 理由二：避免「大对象拷贝」，提升性能
	// 传值会把整个结构体复制一份；传指针只复制一个地址（8 字节），更高效
	fmt.Println("用户名首字母:", firstLetter(&u))

	// ===== 5. new 函数：另一种创建指针的方式 =====
	// new(T) 会分配一个 T 的零值，并返回指向它的指针 *T
	np2 := new(int) // 等价于：var tmp int; np2 := &tmp
	fmt.Println("new(int) 得到指针，解引用值为零值:", *np2) // 0
	*np2 = 100
	fmt.Println("赋值后:", *np2)

	// ===== 6. 结构体指针：. 会自动解引用（语法糖）=====
	up := &User{Name: "Jerry", Age: 20}
	// up.Name 等价于 (*up).Name，Go 自动帮你解引用，不用写成 (*up).Name
	fmt.Println("结构体指针访问字段:", up.Name, up.Age)
	up.Age = 21 // 直接通过指针改字段
	fmt.Println("通过指针改字段后:", *up)

	// ===== 7. 【总结演示】五种引用方式对比 =====
	fmt.Println("\n===== 五种引用方式对比 =====")

	// 方式 1：裸指针 —— 需要解引用
	raw := 42
	rawPtr := &raw
	fmt.Printf("方式1 裸指针: 类型=%T, 打印指针=%v, 解引用=%d\n", rawPtr, rawPtr, *rawPtr)

	// 方式 2：结构体指针 + 自动解引用
	tom := &User{Name: "Tom", Age: 20}
	tom.Age = 21 // 自动解引用，等价于 (*tom).Age = 21
	fmt.Printf("方式2 结构体指针: 类型=%T, 直接访问 tom.Name=%s\n", tom, tom.Name)

	// 方式 3：指针接收者方法
	rect := Rectangle{Width: 3, Height: 4}
	rect.Scale(2) // Go 自动取地址：等价于 (&rect).Scale(2)
	fmt.Printf("方式3 指针接收者: rect.Scale(2) 后 Width=%.0f\n", rect.Width)

	// 方式 4：指针装进接口（最常见！）
	var animal Animal = &Dog{Name: "旺财"}
	// animal 的类型是 Animal 接口，不是 *Dog 裸指针
	// 不需要解引用，直接调接口方法
	fmt.Printf("方式4 指针装进接口: 类型=%T, 调方法=%s\n", animal, animal.Speak())
	// fmt.Println 遇到 error/Stringer 接口会自动调方法，所以不打印地址
	fmt.Println("  直接 Println:", animal) // 如果实现了 String()，打印文本而非地址

	// 方式 5：引用类型（天生含指针）
	s1 := []int{1, 2, 3}
	s2 := s1       // 不用 &，s2 和 s1 已经共享底层数组
	s2[0] = 999    // 改 s2 影响 s1
	fmt.Printf("方式5 引用类型: s1=%v (被 s2 改了), 无需 & 或 *\n", s1)

	// ===== 裸指针 vs 接口里的指针：关键区别 =====
	fmt.Println("\n--- 裸指针 vs 接口里的指针 ---")
	rawDog := &Dog{Name: "小黑"}
	var ifaceDog Animal = rawDog

	fmt.Printf("裸指针打印:   %v (打印地址)\n", rawPtr)      // 打印地址
	fmt.Printf("接口打印:     %v (自动调方法)\n", ifaceDog)    // 自动调 String()
	fmt.Printf("裸指针类型:   %T\n", rawDog)                 // *main.Dog
	fmt.Printf("接口类型:     %T\n", ifaceDog)               // *main.Dog（接口里装的）
	fmt.Println("裸指针要解引用: (*rawDog).Name =", (*rawDog).Name)
	fmt.Println("接口调方法即可:  ifaceDog.Speak() =", ifaceDog.Speak())
}

// addByValue 值传递：n 是 a 的副本，改它不影响外面
func addByValue(n int) {
	n += 100
}

// addByPointer 指针传递：通过 *n 直接修改调用者的变量
func addByPointer(n *int) {
	*n += 100
}

// growUp 用指针接收，才能修改原结构体的字段
func growUp(u *User) {
	u.Age++ // 等价于 (*u).Age++
}

// firstLetter 用指针只为「避免拷贝整个结构体」，这里并不修改它
func firstLetter(u *User) string {
	if len(u.Name) == 0 {
		return ""
	}
	return string([]rune(u.Name)[0])
}

// safeUsePointer 接收一个可能为 nil 的指针，演示函数内如何判空保护
func safeUsePointer(p *int) {
	if p == nil {
		fmt.Println("safeUsePointer: 收到 nil 指针，跳过使用")
		return // 提前返回，避免解引用 nil
	}
	fmt.Println("safeUsePointer: 值 =", *p)
}

// ===== 以下类型用于第 7 节「五种引用方式」演示 =====

// Rectangle 演示指针接收者方法
type Rectangle struct {
	Width  float64
	Height float64
}

// Scale 指针接收者：能修改原对象
func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

// Animal 接口：演示指针装进接口
type Animal interface {
	Speak() string
}

// Dog 实现 Animal 接口
type Dog struct {
	Name string
}

func (d *Dog) Speak() string {
	return d.Name + " 说: 汪汪!"
}

// String 实现 Stringer 接口，让 fmt.Println 自动调用此方法（不打印地址）
func (d *Dog) String() string {
	return "Dog{" + d.Name + "}"
}
