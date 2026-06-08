package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

/*
## 第 13 课：HTTP 服务

### 题 13.1 简易 TODO API（⭐⭐⭐）

用 `net/http` 标准库实现一个 TODO 待办 API：

  - `GET /todos` — 返回所有待办（JSON）
  - `POST /todos` — 新增待办（请求体 `{"title":"xxx"}`)
  - `PUT /todos/:id` — 标记完成（`{"done": true}`）
  - `DELETE /todos/:id` — 删除
    内存存储即可，不用数据库。

### 踩坑总结

1. **defer 必须放在 err 检查之后**：
   如果 os.OpenFile 失败，file 是 nil，defer file.Close() 会 panic。
   正确顺序：先 if err != nil { return }，再 defer file.Close()。

2. **文件写入前必须 Seek + Truncate**：
   读完文件后游标在末尾，直接写会追加在后面，导致文件里出现多行 JSON。
   正确做法：file.Seek(0, 0) 回到开头 + file.Truncate(0) 清空内容，再写入。

3. **fmt.Fprintf 中不能用 %w，要用 %v**：
   %w 只在 fmt.Errorf 中有效（用于包装 error）。
   其他 fmt 函数（Printf/Fprintf/Sprintf）中要用 %v 打印 error。

4. **HTTP 状态码要匹配语义**：
   GET 查询返回 200 OK（http.StatusOK）。
   POST 新增返回 201 Created（http.StatusCreated）。
   不要搞反。

5. **POST 处理完要给客户端返回响应**：
   否则客户端收到空内容，不知道是否成功。
   应返回新增的资源 JSON + 合适的状态码。

6. **文件权限用八进制 0644，不能写十进制 644**：
   Go 里数字前面加 0 表示八进制。644（十进制）≠ 0644（八进制）。
   错误的权限会导致文件读写失败。
*/
type Todo struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	lock   sync.Mutex
	nextId = 0
)

func main() {
	http.HandleFunc("/todos", todos)
	addr := ":9090"
	fmt.Println("服务器已经启动，端口是", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

func todos(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		todoList(writer, request)
	case http.MethodPost:
		addTodo(writer, request)
	}
}

func addTodo(writer http.ResponseWriter, request *http.Request) {
	file, err := os.OpenFile("exercises/13/db_file.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintf(writer, "这里有问题，文件db打开失败：%v", err)
		return
	}
	defer file.Close()
	var todos []Todo
	json.NewDecoder(file).Decode(&todos)
	var todo Todo
	json.NewDecoder(request.Body).Decode(&todo)
	lock.Lock()
	nextId++
	todo.Id = int32(nextId)
	lock.Unlock()
	todos = append(todos, todo)
	// 写入前必须：1.回到文件开头 2.清空文件内容，否则会追加导致多行JSON
	file.Seek(0, 0)
	file.Truncate(0)
	json.NewEncoder(file).Encode(todos)
	// 返回新增的 todo 给客户端
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(todo)
}

func todoList(writer http.ResponseWriter, request *http.Request) {
	file, err := os.OpenFile("exercises/13/db_file.json", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Fprintf(writer, "这里有问题，文件db读取失败：%v", err)
		return
	}
	defer file.Close()
	var todos []Todo
	json.NewDecoder(file).Decode(&todos)
	if todos == nil {
		todos = []Todo{} // 文件为空时返回空数组而不是 null
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK) // GET 查询用 200，不是 201
	json.NewEncoder(writer).Encode(todos)
}
