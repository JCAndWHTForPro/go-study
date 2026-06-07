// 第 12 课：文件 IO + JSON 序列化
// 运行方式：go run ./12-fileio
//
// Go 的 IO 设计哲学：一切皆 Reader/Writer 接口
//   io.Reader: Read(p []byte) (n int, err error)
//   io.Writer: Write(p []byte) (n int, err error)
// 文件、网络连接、标准输入输出都实现了这两个接口，操作方式统一。
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// User 用于 JSON 序列化/反序列化演示
// struct tag（反引号里的 `json:"xxx"`）控制 JSON 字段名
type User struct {
	Name  string `json:"name"`            // JSON 字段名是小写 name
	Age   int    `json:"age"`             // JSON 字段名是小写 age
	Email string `json:"email,omitempty"` // omitempty：值为空时不输出该字段
}

func main() {
	// ===== 1. 写文件：os.WriteFile（最简单的方式）=====
	content := "Hello, Go 文件操作!\n这是第二行\n"
	err := os.WriteFile("12-fileio/test_output.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("写文件失败:", err)
		return
	}
	fmt.Println("✅ 文件写入成功: 12-fileio/test_output.txt")

	// ===== 2. 读文件：os.ReadFile（一次性读完）=====
	data, err := os.ReadFile("12-fileio/test_output.txt")
	if err != nil {
		fmt.Println("读文件失败:", err)
		return
	}
	fmt.Println("读到的内容:")
	fmt.Println(string(data))

	// ===== 3. 逐行读取：os.Open + bufio.Scanner =====
	// 适合大文件，不会一次性把整个文件加载到内存
	file, err := os.Open("12-fileio/test_output.txt")
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer file.Close() // defer 确保函数结束时关闭文件（不管是否出错）

	fmt.Println("逐行读取:")
	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() { // 每次读一行
		fmt.Printf("  第%d行: %s\n", lineNum, scanner.Text())
		lineNum++
	}

	// ===== 4. 追加写入：os.OpenFile =====
	// os.O_APPEND: 追加模式  os.O_CREATE: 不存在则创建  os.O_WRONLY: 只写
	appendFile, err := os.OpenFile("12-fileio/test_output.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	defer appendFile.Close()
	appendFile.WriteString("这是追加的第三行\n")
	fmt.Println("✅ 追加写入成功")

	// ===== 5. JSON 序列化：struct → JSON 字符串 =====
	fmt.Println("\n--- JSON 操作 ---")
	user := User{Name: "Tom", Age: 25, Email: "tom@example.com"}

	// Marshal：struct → []byte（JSON）
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化失败:", err)
		return
	}
	fmt.Println("Marshal:", string(jsonBytes))

	// MarshalIndent：带缩进的美化输出
	prettyJSON, _ := json.MarshalIndent(user, "", "  ")
	fmt.Println("美化输出:\n" + string(prettyJSON))

	// omitempty 演示：Email 为空时不输出
	user2 := User{Name: "Jerry", Age: 20}
	jsonBytes2, _ := json.Marshal(user2)
	fmt.Println("omitempty 效果:", string(jsonBytes2)) // 没有 email 字段

	// ===== 6. JSON 反序列化：JSON 字符串 → struct =====
	jsonStr := `{"name":"Alice","age":30,"email":"alice@go.dev"}`
	var parsed User
	err = json.Unmarshal([]byte(jsonStr), &parsed)
	if err != nil {
		fmt.Println("反序列化失败:", err)
		return
	}
	fmt.Printf("Unmarshal: %+v\n", parsed)

	// ===== 7. JSON 数组 =====
	jsonArr := `[{"name":"A","age":1},{"name":"B","age":2}]`
	var users []User
	json.Unmarshal([]byte(jsonArr), &users)
	fmt.Println("JSON 数组反序列化:", users)

	// ===== 8. JSON 写入文件 / 从文件读取 =====
	// 写入
	jsonFile, _ := os.Create("12-fileio/users.json")
	defer jsonFile.Close()
	encoder := json.NewEncoder(jsonFile) // 直接写入 Writer（文件）
	encoder.SetIndent("", "  ")
	encoder.Encode([]User{
		{Name: "张三", Age: 28, Email: "zhangsan@go.dev"},
		{Name: "李四", Age: 33},
	})
	fmt.Println("✅ JSON 写入文件: 12-fileio/users.json")

	// 读取
	readFile, _ := os.Open("12-fileio/users.json")
	defer readFile.Close()
	var readUsers []User
	decoder := json.NewDecoder(readFile) // 直接从 Reader（文件）读取
	decoder.Decode(&readUsers)
	fmt.Println("从文件读取 JSON:", readUsers)

	fmt.Println("\n✅ 生成的文件保留在 12-fileio/ 目录下，可直接查看")
}
