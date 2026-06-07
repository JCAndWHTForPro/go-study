// 第 9 课：错误处理（error / 自定义错误 / panic / recover / 错误包装）
// 运行方式：go run ./09-errors
//
// Go 的错误处理哲学：「错误是值，不是异常」
//   - 没有 try/catch/throw，用返回值显式处理
//   - error 是一个内置接口，只有一个方法 Error() string
//   - 惯例：error 作为函数的最后一个返回值
package main

import (
	"errors"
	"fmt"
)

func main() {
	// ===== 1. 基础错误处理：if err != nil =====
	// Go 最标志性的写法，你会写无数遍
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("出错了:", err)
	} else {
		fmt.Println("10 / 2 =", result)
	}

	// 故意触发错误
	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("出错了:", err) // 出错了: 除数不能为 0
	}

	// ===== 2. 创建错误的几种方式 =====
	// 方式 A：errors.New —— 最简单，创建一个固定文本的 error
	err1 := errors.New("这是一个简单错误")
	fmt.Println("errors.New:", err1)

	// 方式 B：fmt.Errorf —— 带格式化的错误（可以拼变量进去）
	name := "config.yaml"
	err2 := fmt.Errorf("文件 %s 不存在", name)
	fmt.Println("fmt.Errorf:", err2)

	// ===== 3. 自定义错误类型 =====
	// 实现 error 接口（只需实现 Error() string）就是一个自定义错误
	_, err = withdraw(100, 200)
	if err != nil {
		fmt.Println("取款失败:", err)

		// 类型断言：判断是不是特定的错误类型，拿到更多信息
		//
		// 为什么 balErr 是 *BalanceError（指针），传给 As 还要再加 &？
		//   balErr       类型是 *BalanceError    （一级指针）
		//   &balErr      类型是 **BalanceError   （二级指针：指向指针的指针）
		//
		// 原因：errors.As 的目的是「把匹配到的错误写入 balErr」
		// 和第 7 课的 addByPointer(&a) 完全一样的道理：
		//   函数想改调用者的变量 → 必须传指针
		//   As 想改 balErr 这个指针变量本身（让它指向匹配到的错误）
		//   → 所以要传 &balErr（指针的指针）
		//
		// 如果只传 balErr（不加 &），As 拿到的是副本，
		// 改了副本外面看不到，balErr 仍然是 nil
		var balErr *BalanceError
		if errors.As(err, &balErr) {
			fmt.Printf("  余额不足详情: 余额=%d, 需要=%d, 差额=%d\n",
				balErr.Balance, balErr.Need, balErr.Need-balErr.Balance)
		}
	}

	// ===== 4. 哨兵错误（Sentinel Error）=====
	// 预定义的包级别错误变量，用 errors.Is 判断
	_, err = findUser(0)
	if errors.Is(err, ErrNotFound) {
		fmt.Println("哨兵错误匹配: 用户不存在")
	}

	// ===== 5. 错误包装（Error Wrapping）=====
	// fmt.Errorf + %w 可以「包装」原始错误，保留错误链
	err = loadConfig()
	if err != nil {
		fmt.Println("loadConfig 错误:", err) // 显示完整链
		// errors.Is 可以穿透包装层，找到底层原始错误
		if errors.Is(err, ErrNotFound) {
			fmt.Println("  底层原因: 资源不存在（穿透了包装层）")
		}
	}

	// ===== 5.1 【对比】%w 保留链 vs %v 链断了 =====
	// %w 包装：保留了指向原始错误的链接，errors.Is 能穿透找到
	wrappedW := fmt.Errorf("第1层包装: %w", ErrNotFound)
	wrappedW2 := fmt.Errorf("第2层包装: %w", wrappedW) // 多层包装
	fmt.Println("\n--- %w vs %v 对比 ---")
	fmt.Println("%w 包装后文本:", wrappedW2)
	fmt.Println("%w 穿透2层能否找到 ErrNotFound:", errors.Is(wrappedW2, ErrNotFound)) // true

	// %v 拼接：只是拼了字符串，链断了，errors.Is 找不到
	wrappedV := fmt.Errorf("用 %%v 包装: %v", ErrNotFound)
	fmt.Println("%v 拼接后文本:", wrappedV)
	fmt.Println("%v 能否找到 ErrNotFound:", errors.Is(wrappedV, ErrNotFound)) // false！链断了

	// 结论：
	// %w = 包装（保留链，可穿透）—— 想让调用者判断底层错误，用 %w
	// %v = 只拼文本（链断了）    —— 只想给人看的错误信息，用 %v

	// ===== 6. panic 和 recover =====
	// panic：程序崩溃（类似 Java 的 RuntimeException），用于「不可恢复」的错误
	// recover：在 defer 里捕获 panic，防止程序崩溃
	fmt.Println("\n--- panic / recover 演示 ---")
	safeCall()
	fmt.Println("safeCall 之后程序继续运行（recover 救回来了）")

	// ===== 何时用 error vs panic =====
	// error：业务逻辑中的「预期内」错误（文件不存在、参数非法、网络超时等）
	// panic：程序级的「不可恢复」错误（数组越界、nil 解引用、逻辑 bug）
	// 原则：99% 的情况用 error，panic 极少主动使用
}

// ============================================================
// 基础函数：返回 error
// ============================================================

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("除数不能为 0")
	}
	return a / b, nil
}

// ============================================================
// 自定义错误类型：实现 error 接口
// ============================================================

// BalanceError 自定义错误，携带余额信息
type BalanceError struct {
	Balance int
	Need    int
}

// 实现 error 接口的 Error() string 方法
func (e *BalanceError) Error() string {
	return fmt.Sprintf("余额不足: 当前 %d, 需要 %d", e.Balance, e.Need)
}

func withdraw(balance, amount int) (int, error) {
	if amount > balance {
		return 0, &BalanceError{Balance: balance, Need: amount}
	}
	return balance - amount, nil
}

// ============================================================
// 哨兵错误（Sentinel Error）
// ============================================================

// 包级别预定义的错误变量，用 errors.Is 来判断
var ErrNotFound = errors.New("not found")

func findUser(id int) (string, error) {
	if id == 0 {
		return "", ErrNotFound
	}
	return "Tom", nil
}

// ============================================================
// 错误包装（Error Wrapping）
// ============================================================

func readFile(name string) error {
	return ErrNotFound // 模拟底层错误
}

func loadConfig() error {
	err := readFile("config.yaml")
	if err != nil {
		// %w 包装原始错误，保留错误链（可被 errors.Is / errors.As 穿透）
		return fmt.Errorf("加载配置失败: %w", err)
	}
	return nil
}

// ============================================================
// panic / recover
// ============================================================

func safeCall() {
	// defer + recover 组合：捕获 panic，防止程序崩溃
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover 捕获到 panic:", r)
		}
	}()

	fmt.Println("即将 panic...")
	panic("something went wrong") // 触发 panic
	// 下面这行不会执行
	// fmt.Println("这行不会执行")
}
