# Go 学习练习题集

每章 2-3 道题，从易到难。建议在 `exercises/` 目录下按章节创建文件来写答案。

---

## 第 1-2 课：Hello World & 变量

### 题 1.1 温度转换器（⭐）
写一个程序，定义一个华氏温度变量 `fahrenheit = 98.6`，将它转为摄氏温度并打印。
公式：`celsius = (fahrenheit - 32) * 5 / 9`
要求：用 `:=` 短声明，打印时保留 2 位小数。

### 题 1.2 多返回值交换（⭐）
声明两个变量 `a = 10, b = 20`，不使用第三个变量，利用 Go 的多赋值特性一行完成交换，然后打印交换后的值。

### 题 1.3 iota 星期枚举（⭐⭐）
用 `const + iota` 定义一周七天的枚举（Sunday=0, Monday=1, ..., Saturday=6），然后写一个函数 `dayName(day int) string`，输入数字返回中文星期名。测试打印 Monday 和 Friday 对应的中文。

---

## 第 3 课：流程控制

### 题 3.1 FizzBuzz（⭐）
打印 1 到 30，如果是 3 的倍数打印 "Fizz"，5 的倍数打印 "Buzz"，同时是 3 和 5 的倍数打印 "FizzBuzz"，否则打印数字本身。

### 题 3.2 九九乘法表（⭐⭐）
用嵌套 for 循环打印九九乘法表，格式对齐。

### 题 3.3 猜数字（⭐⭐）
用 `math/rand` 生成一个 1-100 的随机数，模拟猜 10 次（用 for 循环，每次猜中间值，类似二分查找），打印每次猜的过程和最终结果。

---

## 第 4 课：函数

### 题 4.1 递归阶乘（⭐）
写一个递归函数 `factorial(n int) int`，计算 n 的阶乘。测试 `factorial(5)` 应该返回 120。

### 题 4.2 闭包计数器（⭐⭐）
写一个函数 `makeCounter() func() int`，每次调用返回的函数时，返回值递增（1, 2, 3, ...）。用闭包实现。

### 题 4.3 defer 执行顺序（⭐⭐）
写一个函数，在 for 循环中 defer 打印 0-4，观察并解释为什么输出是 4, 3, 2, 1, 0（栈顺序）。

---

## 第 5 课：复合类型

### 题 5.1 切片去重（⭐⭐）
写一个函数 `unique(nums []int) []int`，返回去重后的切片（保持原顺序）。
示例：`unique([]int{1,3,2,3,1,4})` → `[1,3,2,4]`

### 题 5.2 单词频率统计（⭐⭐）
给定一个字符串 `"go is great go is fast go go go"`，用 map 统计每个单词出现的次数，按频率从高到低排序输出。

### 题 5.3 矩阵转置（⭐⭐⭐）
写一个函数 `transpose(matrix [][]int) [][]int`，将 M×N 矩阵转置为 N×M。
示例：`[[1,2,3],[4,5,6]]` → `[[1,4],[2,5],[3,6]]`

---

## 第 6 课：接口

### 题 6.1 面积计算器（⭐⭐）
定义一个 `Shape` 接口，包含 `Area() float64` 方法。
实现三个类型：`Circle`、`Rectangle`、`Triangle`，各自实现 `Area()`。
写一个函数 `totalArea(shapes []Shape) float64` 计算总面积。

### 题 6.2 排序接口（⭐⭐⭐）
自定义一个 `StringSlice` 类型（底层是 `[]string`），实现 `sort.Interface` 的三个方法（`Len`、`Less`、`Swap`），让字符串按长度排序（短的在前）。

---

## 第 7 课：指针

### 题 7.1 交换函数（⭐）
写一个函数 `swap(a, b *int)` 通过指针交换两个整数的值。

### 题 7.2 链表（⭐⭐⭐）
定义一个单链表节点 `type Node struct { Val int; Next *Node }`：
1. 写 `push(head **Node, val int)` 在头部插入
2. 写 `printList(head *Node)` 打印整个链表
3. 写 `reverse(head *Node) *Node` 反转链表

---

## 第 8 课：并发

### 题 8.1 并发求和（⭐⭐）
把一个 100 个元素的 slice 分成 4 段，每段用一个 goroutine 求和，通过 channel 汇总，打印总和。

### 题 8.2 生产者消费者（⭐⭐⭐）
启动 3 个生产者 goroutine（各生产 5 个商品），1 个消费者 goroutine 消费所有商品。用 channel 通信，所有生产完成后关闭 channel，消费者用 for-range 读取。

---

## 第 9 课：错误处理

### 题 9.1 安全除法（⭐）
写一个函数 `divide(a, b float64) (float64, error)`，当 b 为 0 时返回自定义错误 `DivideByZeroError`。用 `errors.Is` 判断错误类型。

### 题 9.2 多层错误包装（⭐⭐）
模拟三层调用链：`handler() → service() → dao()`。
dao 返回原始错误，service 用 `%w` 包装，handler 用 `errors.As` 提取原始错误。打印整个错误链。

---

## 第 10 课：泛型

### 题 10.1 泛型 Map/Filter（⭐⭐）
写两个泛型函数：
- `Map[T, U](slice []T, fn func(T) U) []U` — 映射
- `Filter[T](slice []T, fn func(T) bool) []T` — 过滤
用例：把 `[]int{1,2,3,4,5}` 过滤出偶数，然后每个乘以 10。

### 题 10.2 泛型缓存（⭐⭐⭐）
实现一个泛型的简易缓存 `type Cache[K comparable, V any] struct`：
- `Set(key K, value V)`
- `Get(key K) (V, bool)`
- `Delete(key K)`
用 map 实现。

---

## 第 11 课：包管理

### 题 11.1 工具包（⭐⭐）
在项目中创建一个 `exercises/mathutil` 包，提供以下导出函数：
- `GCD(a, b int) int` — 最大公约数
- `LCM(a, b int) int` — 最小公倍数
- `IsPrime(n int) bool` — 判断质数
在 main 中导入并测试。

---

## 第 12 课：文件 IO + JSON

### 题 12.1 配置文件读写（⭐⭐）
定义一个 `Config` struct（含 Host, Port, Debug 等字段），实现：
1. `SaveConfig(path string, cfg Config) error` — 写入 JSON 文件
2. `LoadConfig(path string) (Config, error)` — 从 JSON 文件读取
测试：保存一个配置，再读回来，对比是否一致。

---

## 第 13 课：HTTP 服务

### 题 13.1 简易 TODO API（⭐⭐⭐）
用 `net/http` 标准库实现一个 TODO 待办 API：
- `GET /todos` — 返回所有待办（JSON）
- `POST /todos` — 新增待办（请求体 `{"title":"xxx"}`)
- `PUT /todos/:id` — 标记完成（`{"done": true}`）
- `DELETE /todos/:id` — 删除
内存存储即可，不用数据库。

---

## 第 14 课：数据结构

### 题 14.1 用栈判断括号匹配（⭐⭐）
写一个函数 `isValid(s string) bool`，判断字符串中的括号是否匹配。
支持 `()`, `[]`, `{}`。
示例：`"([{}])"` → true，`"([)]"` → false

### 题 14.2 TopK 问题（⭐⭐⭐）
给定一个整数 slice，用 `container/heap` 实现找出最大的 K 个元素。
要求用最小堆，堆大小保持 K。

---

## 第 15 课：测试

### 题 15.1 给练习题加测试（⭐⭐）
选择前面任意 3 道题的函数，用表驱动测试写完整的测试用例（至少 5 个用例）。
要求覆盖正常值、边界值、异常值。

---

## 第 16 课：Context

### 题 16.1 超时搜索（⭐⭐⭐）
模拟一个搜索场景：同时向 3 个"搜索引擎"发请求（用 goroutine + 随机延迟模拟），设置 500ms 超时。
- 如果某个引擎在超时内返回 → 打印结果
- 如果全部超时 → 打印"搜索超时"
用 `context.WithTimeout` + `select` 实现。

---

## 第 17 课：数据库

### 题 17.1 学生管理系统（⭐⭐）
用 `database/sql` 实现一个学生表的 CRUD：
- 建表（id, name, grade, score）
- 插入 5 个学生
- 查询分数大于 80 的学生
- 更新某个学生的分数
- 删除一个学生
- 事务：给所有学生加 5 分（原子操作）

---

## 第 18 课：GORM

### 题 18.1 博客系统模型（⭐⭐⭐）
用 GORM 实现一个简易博客的数据模型：
- `User`（id, username, email）
- `Post`（id, title, content, user_id）→ 属于 User
- `Comment`（id, content, post_id, user_id）→ 属于 Post 和 User
实现：
1. AutoMigrate 建表
2. 创建 2 个用户，每人写 2 篇文章，每篇文章 2 条评论
3. Preload 查询：查某个用户的所有文章及每篇文章的评论
4. Joins 查询：查某篇文章的所有评论及评论者名字

---

## 第 19 课：Gin

### 题 19.1 Gin + GORM 整合（⭐⭐⭐）
把第 18 题的博客系统改成 Gin API 服务：
- `GET /api/users/:id/posts` — 查某用户的所有文章（Preload 评论）
- `POST /api/posts` — 发文章（JSON: title, content, user_id）
- `POST /api/posts/:id/comments` — 给文章评论
- 加一个日志中间件记录请求耗时
- 路由分组 `/api`

---

## 难度说明

- ⭐ 基础题：直接用课上学的知识就能写
- ⭐⭐ 进阶题：需要组合多个知识点
- ⭐⭐⭐ 综合题：需要独立设计 + 综合运用

## 建议

1. **先不看代码**，自己尝试写，写不出来再回去翻对应章节
2. 每道题写完后用 `go run` 验证，确保能运行
3. 关键题目（带 ⭐⭐⭐ 的）写完后加上测试用例
4. 题 14.1（括号匹配）和 14.2（TopK）是经典算法题，推荐认真做
