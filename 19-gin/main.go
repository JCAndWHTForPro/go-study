// 第 19 课：Gin Web 框架
// 运行方式：go run ./19-gin
// 启动后测试：
//   GET  http://localhost:9090/
//   GET  http://localhost:9090/hello?name=Tom
//   GET  http://localhost:9090/users/3
//   GET  http://localhost:9090/api/v1/products
//   POST http://localhost:9090/api/v1/products  (Body: {"name":"键盘","price":299})
//   按 Ctrl+C 停止
//
// Gin vs net/http（第 13 课）：
//   net/http → 标准库，功能够用但写起来繁琐（路由简陋、参数解析手动）
//   Gin      → 最流行的 Go Web 框架，性能极高，开发体验好
//
// ============================================================
// 【总结】Gin 核心概念
// ------------------------------------------------------------
//
// ▶ 对比 net/http：
//   net/http                           Gin
//   http.HandleFunc("/path", handler)  r.GET("/path", handler)
//   r.URL.Query().Get("name")          c.Query("name")
//   手动解析路径参数                    c.Param("id")（路由 /users/:id）
//   手动 json.NewDecoder(r.Body)       c.ShouldBindJSON(&data)
//   手动 w.Header().Set + json.Encode  c.JSON(200, data)
//   没有路由分组                        r.Group("/api/v1")
//   没有中间件机制                      r.Use(middleware)
//
// ▶ 核心对象：
//   gin.Default()     → 创建引擎（自带 Logger + Recovery 中间件）
//   gin.New()         → 创建引擎（不带任何中间件）
//   *gin.Context (c)  → 每个请求的上下文（读请求 + 写响应，全靠它）
//
// ▶ 路由注册：
//   r.GET("/path", handler)            → GET 请求
//   r.POST("/path", handler)           → POST 请求
//   r.PUT("/path", handler)            → PUT 请求
//   r.DELETE("/path", handler)         → DELETE 请求
//   r.Any("/path", handler)            → 任意方法
//
// ▶ 参数获取（都通过 c *gin.Context）：
//   c.Query("name")                    → URL 参数 ?name=Tom
//   c.DefaultQuery("name", "World")    → 有默认值的 URL 参数
//   c.Param("id")                      → 路径参数 /users/:id
//   c.ShouldBindJSON(&data)            → 解析 JSON 请求体到 struct
//   c.PostForm("field")                → 表单参数
//   c.GetHeader("Authorization")       → 请求头
//
// ▶ 响应（都通过 c *gin.Context）：
//   c.JSON(200, gin.H{"key": "value"}) → JSON 响应（gin.H 是 map 的别名）
//   c.JSON(200, structData)            → JSON 响应（直接传 struct）
//   c.String(200, "text")              → 纯文本响应
//   c.Status(204)                      → 只返回状态码，无 Body
//
// ▶ 路由分组（Group）：
//   v1 := r.Group("/api/v1")           → 创建分组，公共前缀
//   v1.GET("/users", handler)          → 实际路径 /api/v1/users
//   v1.Use(authMiddleware)             → 给分组加中间件
//
// ▶ 中间件（Middleware）：
//   中间件 = 在处理请求前后插入的公共逻辑（日志、鉴权、CORS 等）
//   func MyMiddleware() gin.HandlerFunc {
//       return func(c *gin.Context) {
//           // 请求前的逻辑
//           c.Next()   // 调用下一个处理函数
//           // 请求后的逻辑
//       }
//   }
//   r.Use(MyMiddleware())              → 全局中间件
//   group.Use(MyMiddleware())          → 分组中间件
//   r.GET("/path", MyMiddleware(), handler) → 单个路由中间件
//
// ▶ 对比 Java Spring Boot：
//   Spring Boot                        Gin
//   @RestController                    r.GET / r.POST（函数式注册）
//   @GetMapping("/users")              r.GET("/users", handler)
//   @PathVariable Long id              c.Param("id")
//   @RequestParam String name          c.Query("name")
//   @RequestBody User user             c.ShouldBindJSON(&user)
//   ResponseEntity.ok(data)            c.JSON(200, data)
//   @RequestMapping("/api/v1")         r.Group("/api/v1")
//   @Component Filter / Interceptor    r.Use(middleware)
// ============================================================
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Product 商品
type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required"`  // binding:"required" → ShouldBind 时校验非空
	Price int    `json:"price" binding:"required"`
}

// 内存数据模拟
var (
	products = []Product{
		{ID: 1, Name: "Go 语言圣经", Price: 59},
		{ID: 2, Name: "机械键盘", Price: 299},
		{ID: 3, Name: "显示器", Price: 1999},
	}
	nextID = 4
	mu     sync.Mutex
)

func main() {
	// ===== 1. 创建 Gin 引擎 =====
	// Default 自带两个中间件：Logger（打印请求日志）+ Recovery（panic 不崩服务）
	r := gin.Default()

	// ===== 2. 最简单的路由 =====
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "欢迎来到 Gin Web 服务！")
	})

	// ===== 3. Query 参数 =====
	r.GET("/hello", func(c *gin.Context) {
		name := c.DefaultQuery("name", "World") // 有默认值
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Hello, %s! 👋", name),
		})
	})

	// ===== 4. 路径参数（:id）=====
	r.GET("/users/:id", func(c *gin.Context) {
		idStr := c.Param("id") // 取路径里的 :id
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id 必须是数字"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("查询用户 ID=%d", id),
		})
	})

	// ===== 5. 路由分组（Group）=====
	v1 := r.Group("/api/v1")
	v1.Use(logMiddleware()) // 给 v1 分组加中间件
	{
		v1.GET("/products", getProducts)       // GET /api/v1/products
		v1.GET("/products/:id", getProductByID) // GET /api/v1/products/1
		v1.POST("/products", createProduct)    // POST /api/v1/products
		v1.PUT("/products/:id", updateProduct)  // PUT /api/v1/products/1
		v1.DELETE("/products/:id", deleteProduct) // DELETE /api/v1/products/1
	}

	// ===== 6. 全局中间件 =====
	// gin.Default() 已自带 Logger 和 Recovery
	// 这里再加一个自定义的 CORS 中间件
	r.Use(corsMiddleware())

	// ===== 启动 =====
	addr := ":9090"
	fmt.Println("🚀 Gin 服务器启动: http://localhost" + addr)
	fmt.Println("   测试命令:")
	fmt.Println("   curl http://localhost:9090/")
	fmt.Println("   curl 'http://localhost:9090/hello?name=Tom'")
	fmt.Println("   curl http://localhost:9090/users/3")
	fmt.Println("   curl http://localhost:9090/api/v1/products")
	fmt.Println("   curl http://localhost:9090/api/v1/products/1")
	fmt.Println("   curl -X POST http://localhost:9090/api/v1/products -H 'Content-Type: application/json' -d '{\"name\":\"鼠标\",\"price\":99}'")
	fmt.Println("   curl -X PUT http://localhost:9090/api/v1/products/1 -H 'Content-Type: application/json' -d '{\"name\":\"Go圣经(第2版)\",\"price\":69}'")
	fmt.Println("   curl -X DELETE http://localhost:9090/api/v1/products/3")

	err := r.Run(addr) // 等价于 http.ListenAndServe
	if err != nil {
		log.Fatal("启动失败:", err)
	}
}

// ===== RESTful CRUD 处理函数 =====

// getProducts GET /api/v1/products → 获取所有商品
func getProducts(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, products)
}

// getProductByID GET /api/v1/products/:id → 获取单个商品
func getProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id 必须是数字"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
}

// createProduct POST /api/v1/products → 新增商品
func createProduct(c *gin.Context) {
	var input Product
	// ShouldBindJSON：解析 JSON + 校验 binding 标签
	// 比第 13 课的 json.NewDecoder(r.Body).Decode(&input) 更强大
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	input.ID = nextID
	nextID++
	products = append(products, input)
	mu.Unlock()

	c.JSON(http.StatusCreated, input) // 201 Created
}

// updateProduct PUT /api/v1/products/:id → 更新商品
func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id 必须是数字"})
		return
	}

	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, p := range products {
		if p.ID == id {
			input.ID = id
			products[i] = input
			c.JSON(http.StatusOK, input)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
}

// deleteProduct DELETE /api/v1/products/:id → 删除商品
func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id 必须是数字"})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("商品 %s 已删除", p.Name)})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
}

// ===== 中间件 =====

// logMiddleware 自定义日志中间件：记录每个请求的耗时
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// c.Next() 调用后面的处理函数（和其他中间件）
		c.Next()

		// 处理完之后，计算耗时
		duration := time.Since(start)
		log.Printf("[自定义日志] %s %s → %d (%v)\n",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}

// corsMiddleware 跨域中间件（前后端分离时必须）
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 预检请求直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
