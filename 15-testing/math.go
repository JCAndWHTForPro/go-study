// 第 15 课：测试（go test）
// 运行方式：
//   go test ./15-testing            → 运行所有测试
//   go test ./15-testing -v         → 显示详细输出（每个用例的结果）
//   go test ./15-testing -run TestAdd → 只运行名字匹配的测试
//   go test ./15-testing -cover     → 查看测试覆盖率
//   go test ./15-testing -bench .   → 运行基准测试
//
// Go 测试的设计哲学：
//   1. 测试是一等公民，go test 内置在工具链里，不需要 JUnit 之类的框架
//   2. 测试文件必须以 _test.go 结尾，和被测代码放同一个包里
//   3. 测试函数必须以 Test 开头，参数是 *testing.T
//   4. 基准测试函数以 Benchmark 开头，参数是 *testing.B
//
// ============================================================
// 【总结】Go 测试 vs Java 测试
// ------------------------------------------------------------
//   Java (JUnit)                  Go (testing)
//   @Test                         func TestXxx(t *testing.T)
//   @BeforeEach / @AfterEach      TestMain(m *testing.M) 或 t.Cleanup()
//   assertEquals(expected, actual) if got != want { t.Errorf(...) }
//   @ParameterizedTest            表驱动测试（Table-Driven Test）
//   JMH 基准测试                  func BenchmarkXxx(b *testing.B)
//   Mockito                       接口 + 手写 mock（或 gomock 库）
//   assertThrows                  没有异常机制，直接判 err != nil
//
// ▶ Go 没有 assert 函数！用 if + t.Errorf 手动判断
//   这是故意的设计：Go 认为 assert 隐藏了测试逻辑，不够清晰
//
// ▶ 表驱动测试（Table-Driven Test）是 Go 的标准风格：
//   把所有测试用例放进一个 slice，用 for range + t.Run 逐个跑
//   优点：新增用例只需加一行数据，不用写新函数
//
// ▶ t.Run("子测试名", func(t *testing.T) { ... })
//   子测试：可以用 -run TestAdd/正数 只跑某个子用例
//
// ▶ t.Errorf vs t.Fatalf：
//   t.Errorf → 报错但继续跑后面的用例（推荐，能看到所有失败）
//   t.Fatalf → 报错并立即停止当前测试函数
//
// ▶ 基准测试（Benchmark）：
//   func BenchmarkXxx(b *testing.B) { for i := 0; i < b.N; i++ { ... } }
//   Go 会自动调整 b.N 的值，直到结果稳定
// ============================================================
package main

// Add 两数相加
func Add(a, b int) int {
	return a + b
}

// Abs 取绝对值
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Max 返回两个数中较大的
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Fibonacci 返回第 n 个斐波那契数（从 0 开始）
// Fib(0)=0, Fib(1)=1, Fib(2)=1, Fib(3)=2, Fib(4)=3, Fib(5)=5
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	prev, curr := 0, 1
	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

// IsPalindrome 判断字符串是否是回文（支持中文）
func IsPalindrome(s string) bool {
	runes := []rune(s)
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		if runes[left] != runes[right] {
			return false
		}
	}
	return true
}
