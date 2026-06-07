// 第 13 课：HTTP 服务（net/http 标准库）
// 运行方式：go run ./13-http
// 启动后在浏览器访问：
//   http://localhost:8080/           → 首页
//   http://localhost:8080/hello      → Hello 接口
//   http://localhost:8080/hello?name=Tom → 带参数
//   http://localhost:8080/users      → GET 获取用户列表 / POST 新增用户
//   http://localhost:8080/json       → 返回 JSON
//   按 Ctrl+C 停止服务器
//
// Go 的 HTTP 设计哲学：
//   核心接口只有一个：http.Handler → ServeHTTP(w, r)
//   http.HandleFunc 是便捷注册方式，底层还是 Handler 接口
//   标准库的 net/http 性能极强，不需要框架也能上生产
//
// ============================================================
// 【总结】net/http 核心概念
// ------------------------------------------------------------
//
// ▶ 两个核心参数（每个处理函数都有）：
//   w http.ResponseWriter → 往客户端「写」响应（状态码、Header、Body）
//   r *http.Request       → 从客户端「读」请求（URL、Method、Header、Body）
//
// ▶ 路由注册（HandleFunc vs Handle）：
//   http.HandleFunc("/path", handler)  → 传一个普通函数，简单场景首选
//   http.Handle("/path", handler)      → 传一个实现了 Handler 接口的对象
//
//   区别：
//     HandleFunc：直接传 func(w, r)，不需要定义类型，适合简单逻辑
//     Handle：    传实现了 ServeHTTP(w, r) 的 struct，可以携带状态（DB连接、配置等）
//
//   本质：HandleFunc 是语法糖，底层源码就是：
//     func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
//         Handle(pattern, HandlerFunc(handler))  // 把函数包装成 Handler 接口
//     }
//     type HandlerFunc func(ResponseWriter, *Request)
//     func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) { f(w, r) }
//   → HandlerFunc 是一个「让函数实现接口」的类型，Go 中非常经典的模式
//
// ▶ 请求信息获取：
//   r.Method           → "GET" / "POST" / "PUT" / "DELETE"
//   r.URL.Path         → 请求路径，如 "/users"
//   r.URL.Query()      → URL 参数，如 ?name=Tom → map[name:[Tom]]
//   r.Header.Get("X")  → 获取请求头
//   r.Body             → 请求体（io.ReadCloser），POST/PUT 时用
//
// ▶ 响应写入：
//   w.WriteHeader(code)              → 设置状态码（必须在 Write 之前调用）
//   w.Header().Set("Key", "Value")   → 设置响应头
//   w.Write([]byte("text"))          → 写响应体
//   fmt.Fprintf(w, "格式化 %s", v)    → 格式化写入（w 实现了 io.Writer）
//
// ▶ JSON 响应套路：
//   w.Header().Set("Content-Type", "application/json")
//   json.NewEncoder(w).Encode(data)  → 直接把 struct/slice 序列化写入响应
//
// ▶ JSON 请求体解析套路：
//   var data MyStruct
//   json.NewDecoder(r.Body).Decode(&data)  → 从请求体解析 JSON 到 struct
//
// ▶ 启动服务器：
//   http.ListenAndServe(":8080", nil)  → 监听端口，nil 表示用默认路由
//   这个函数会阻塞，直到出错才返回
// ============================================================
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// User 用户结构体，用于 JSON 交互
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 用内存 slice 模拟数据库，加锁保证并发安全
var (
	users  = []User{
		{ID: 1, Name: "张三", Age: 28},
		{ID: 2, Name: "李四", Age: 33},
	}
	nextID = 3
	mu     sync.Mutex // 保护 users 和 nextID 的并发访问
)

func main() {
	// ===== 1. 最简单的路由：返回纯文本 =====
	http.HandleFunc("/", handleIndex)

	// ===== 2. 带 URL 参数的路由 =====
	http.HandleFunc("/hello", handleHello)

	// ===== 3. 返回 JSON =====
	http.HandleFunc("/json", handleJSON)

	// ===== 4. RESTful 风格：同一路径根据 Method 区分行为 =====
	http.HandleFunc("/users", handleUsers)

	// ===== 启动服务器 =====
	addr := ":9090"
	fmt.Println("🚀 服务器启动: http://localhost" + addr)
	fmt.Println("   试试访问:")
	fmt.Println("   GET  http://localhost:9090/")
	fmt.Println("   GET  http://localhost:9090/hello?name=Tom")
	fmt.Println("   GET  http://localhost:9090/json")
	fmt.Println("   GET  http://localhost:9090/users")
	fmt.Println("   POST http://localhost:9090/users  (Body: {\"name\":\"王五\",\"age\":25})")
	fmt.Println("   按 Ctrl+C 停止")

	// ListenAndServe 会阻塞，直到出错才返回
	// 第二个参数 nil 表示使用默认的 DefaultServeMux（上面 HandleFunc 注册的就在这里）
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}

// ===== 处理函数 =====

// handleIndex 首页，返回纯文本
func handleIndex(w http.ResponseWriter, r *http.Request) {
	// 注意：Go 的 "/" 路由是通配的，匹配所有未被其他路由匹配的路径
	// 如果想严格只匹配 "/"，需要手动判断
	if r.URL.Path != "/" {
		http.NotFound(w, r) // 返回 404
		return
	}
	// fmt.Fprintf 可以直接往 w 写（因为 w 实现了 io.Writer 接口）
	fmt.Fprintf(w, "欢迎来到 Go HTTP 服务！当前路径: %s\n", r.URL.Path)
}

// handleHello 带参数的路由
func handleHello(w http.ResponseWriter, r *http.Request) {
	// r.URL.Query() 返回 url.Values（本质是 map[string][]string）
	name := r.URL.Query().Get("name") // 获取 ?name=xxx 参数
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello, %s! 👋\n", name)
}

// handleJSON 返回 JSON 格式数据
func handleJSON(w http.ResponseWriter, r *http.Request) {
	// 响应 JSON 的标准套路：
	// 1. 设置 Content-Type 头（告诉客户端这是 JSON）
	w.Header().Set("Content-Type", "application/json")

	// 2. 构造数据
	response := map[string]any{
		"message": "这是 JSON 响应",
		"code":    200,
		"data":    []string{"Go", "是", "最好的语言"},
	}

	// 3. 用 json.NewEncoder 直接写入 ResponseWriter
	// 之前第 12 课学过：Encoder 可以直接写入任何 io.Writer
	json.NewEncoder(w).Encode(response)
}

// handleUsers RESTful 风格：根据请求方法分发
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetUsers(w, r)
	case http.MethodPost:
		handleCreateUser(w, r)
	default:
		// 405 Method Not Allowed
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "不支持的方法: "+r.Method, http.StatusMethodNotAllowed)
	}
}

// handleGetUsers GET /users → 返回用户列表
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// handleCreateUser POST /users → 新增用户
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// 1. 从请求体解析 JSON
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "JSON 格式错误: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. 校验必填字段
	if newUser.Name == "" {
		http.Error(w, "name 不能为空", http.StatusBadRequest)
		return
	}

	// 3. 分配 ID 并添加到列表
	mu.Lock()
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)
	mu.Unlock()

	// 4. 返回 201 Created + 新用户 JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201，必须在 Write/Encode 之前调用
	json.NewEncoder(w).Encode(newUser)

	log.Printf("✅ 新增用户: %+v\n", newUser)
}
