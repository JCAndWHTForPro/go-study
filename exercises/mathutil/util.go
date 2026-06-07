package mathutil

import "math"

// GCD 最大公约数（辗转相除法）
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM 最小公倍数 = a * b / GCD(a, b)
func LCM(a, b int) int {
	return a * b / GCD(a, b)
}

// IsPrime 判断质数
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	limit := int(math.Sqrt(float64(n)))
	for i := 2; i <= limit; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
