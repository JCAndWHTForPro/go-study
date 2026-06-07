// b_handler.go：演示同包多文件 init 执行顺序
// 文件名以 b_ 开头，字母序第二 → 在 a_config.go 之后执行
package main

import "fmt"

func init() {
	fmt.Println("[init] b_handler.go 的 init（字母序第二）")
}
