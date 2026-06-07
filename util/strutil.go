// 演示「目录即包」+「跨包导入」
//
// 要点：
// 1) 这个文件在 util/ 目录下，所以包名就叫 util（一个目录 = 一个包）
// 2) 包名是「单个简单词」util，不能写成 com.go.util 这种带点的层级
// 3) 想被别的包（如 01-hello）调用的函数，首字母必须「大写」导出
package util

import "strings"

// Reverse 把字符串按字符反转后返回
// 首字母大写 Reverse：导出函数，可被其他包通过 util.Reverse 调用
func Reverse(s string) string {
	runes := []rune(s) // 转成 rune 切片，正确处理中文等多字节字符
	for left, right := 0, len(runes)-1; left < right; left, right = left+1, right-1 {
		runes[left], runes[right] = runes[right], runes[left]
	}
	return string(runes)
}

// ToUpper 把字符串转成大写并返回
func ToUpper(s string) string {
	return strings.ToUpper(s)
}
