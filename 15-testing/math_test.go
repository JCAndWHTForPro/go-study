package main

import "testing"

// ===== 1. 最简单的测试：一个函数一个用例 =====
func TestAddSimple(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d, want %d", got, want)
	}
}

// ===== 2. 表驱动测试（Table-Driven Test）—— Go 的标准风格 =====
// 把所有测试用例放进一个 slice，用 for range + t.Run 逐个跑
// 新增用例只需加一行数据，不用写新函数
func TestAdd(t *testing.T) {
	tests := []struct {
		name     string // 子测试名（会显示在输出里）
		a, b     int
		expected int
	}{
		{"正数相加", 2, 3, 5},
		{"负数相加", -1, -2, -3},
		{"正负相加", 10, -3, 7},
		{"零值", 0, 0, 0},
		{"大数", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ===== 3. Abs 的表驱动测试 =====
func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"正数", 5, 5},
		{"负数", -5, 5},
		{"零", 0, 0},
		{"最小负数+1", -2147483647, 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Abs(tt.input)
			if got != tt.expected {
				t.Errorf("Abs(%d) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

// ===== 4. Max 测试 =====
func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"a大", 10, 5, 10},
		{"b大", 3, 8, 8},
		{"相等", 7, 7, 7},
		{"负数", -1, -5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ===== 5. Fibonacci 测试 =====
func TestFibonacci(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		expected int
	}{
		{"Fib(0)", 0, 0},
		{"Fib(1)", 1, 1},
		{"Fib(2)", 2, 1},
		{"Fib(3)", 3, 2},
		{"Fib(5)", 5, 5},
		{"Fib(10)", 10, 55},
		{"Fib(20)", 20, 6765},
		{"负数", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fibonacci(tt.n)
			if got != tt.expected {
				t.Errorf("Fibonacci(%d) = %d, want %d", tt.n, got, tt.expected)
			}
		})
	}
}

// ===== 6. IsPalindrome 测试（含中文）=====
func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"英文回文", "racecar", true},
		{"英文非回文", "hello", false},
		{"中文回文", "上海自来水来自海上", true},
		{"中文非回文", "你好世界", false},
		{"单字符", "a", true},
		{"空字符串", "", true},
		{"两字符回文", "aa", true},
		{"两字符非回文", "ab", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPalindrome(tt.input)
			if got != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// ===== 7. 基准测试（Benchmark）=====
// 用于测量函数性能，Go 会自动调整 b.N 直到结果稳定
// 运行：go test ./15-testing -bench .
// 输出示例：BenchmarkFibonacci20-8    50000000    25.3 ns/op
//          → 8 核，跑了 5000 万次，每次 25.3 纳秒

func BenchmarkFibonacci10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

func BenchmarkFibonacci20(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(20)
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPalindrome("上海自来水来自海上")
	}
}
