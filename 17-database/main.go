// 第 17 课：数据库 CRUD（database/sql 标准库 + MySQL）
// 运行方式：go run ./17-database
//
// 前置条件：
//   1. 本地安装 MySQL，端口 3306，用户 root，密码 root123
//   2. 已创建数据库：CREATE DATABASE go_study DEFAULT CHARSET utf8mb4;
//   3. 已安装驱动：go get github.com/go-sql-driver/mysql
//
// Go 操作数据库的设计哲学：
//   标准库 database/sql 提供统一接口，驱动负责具体实现（类似 Java JDBC）。
//   import _ "github.com/go-sql-driver/mysql" 只注册驱动，不直接用它的 API。
//
// ============================================================
// 【总结】database/sql 核心概念
// ------------------------------------------------------------
//
// ▶ 连接：
//   sql.Open("mysql", dsn) → 返回 *sql.DB（连接池，不是单个连接！）
//   db.Ping()              → 验证连接是否真的通了
//   defer db.Close()       → 程序结束时关闭连接池
//
// ▶ DSN（数据源名称）格式：
//   "用户名:密码@tcp(地址:端口)/数据库名?参数"
//   "root:root123@tcp(127.0.0.1:3306)/go_study?charset=utf8mb4&parseTime=True"
//
// ▶ 增删改查四大操作：
//   db.Exec(sql, args...)       → INSERT / UPDATE / DELETE（不返回行）
//   db.QueryRow(sql, args...)   → 查一行（SELECT ... LIMIT 1）
//   db.Query(sql, args...)      → 查多行（SELECT ...）
//
// ▶ QueryRow 使用套路：
//   var name string
//   err := db.QueryRow("SELECT name FROM users WHERE id=?", 1).Scan(&name)
//   // Scan 把结果填进变量（又是 any 参数传 & 的套路！第 6 课讨论过）
//   // 如果没查到行，err == sql.ErrNoRows
//
// ▶ Query 使用套路：
//   rows, err := db.Query("SELECT id, name FROM users")
//   defer rows.Close()       // 必须关！否则连接泄漏
//   for rows.Next() {        // 逐行读取
//       rows.Scan(&id, &name)
//   }
//
// ▶ 预编译语句（Prepared Statement）：
//   stmt, err := db.Prepare("INSERT INTO users(name, age) VALUES(?, ?)")
//   defer stmt.Close()
//   stmt.Exec("张三", 28)    // 复用语句，只换参数
//   stmt.Exec("李四", 33)    // 防 SQL 注入 + 性能更好
//
// ▶ 事务：
//   tx, err := db.Begin()       // 开启事务
//   tx.Exec(...)                // 在事务内执行
//   tx.Commit()                 // 提交（成功）
//   tx.Rollback()               // 回滚（失败）
//
// ▶ ? 占位符：
//   Go 的 MySQL 驱动用 ? 做占位符（和 Java JDBC 一样）
//   db.Query("SELECT * FROM users WHERE id = ?", 1)
//   永远用 ? 传参，不要拼字符串，防止 SQL 注入！
//
// ▶ 对比 Java JDBC：
//   Java                              Go
//   DriverManager.getConnection()     sql.Open()
//   Connection                        *sql.DB（注意：是连接池！）
//   PreparedStatement                 *sql.Stmt
//   ResultSet.next() + getString()    rows.Next() + rows.Scan()
//   connection.setAutoCommit(false)   db.Begin()
//   connection.commit()               tx.Commit()
// ============================================================
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // 只注册驱动，不直接使用
)

// User 对应数据库表的一行
type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// ===== 1. 连接数据库 =====
	// sql.Open 不会立即连接，只是创建连接池对象
	dsn := "root:root123@tcp(127.0.0.1:3306)/go_study?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("打开数据库失败:", err)
	}
	defer db.Close() // 程序结束时关闭连接池

	// Ping 才是真正尝试连接
	err = db.Ping()
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	fmt.Println("✅ 数据库连接成功")

	// 配置连接池（生产环境建议设置）
	db.SetMaxOpenConns(10)  // 最大打开连接数
	db.SetMaxIdleConns(5)   // 最大空闲连接数

	// ===== 2. 建表（Exec 执行 DDL）=====
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id   INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		age  INT NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("建表失败:", err)
	}
	fmt.Println("✅ 表 users 创建成功（或已存在）")

	// 清空表，保证每次运行从干净状态开始
	db.Exec("TRUNCATE TABLE users")

	// ===== 3. 插入数据（INSERT）=====
	fmt.Println("\n--- INSERT ---")

	// 方式一：直接 Exec
	result, err := db.Exec("INSERT INTO users(name, age) VALUES(?, ?)", "张三", 28)
	if err != nil {
		log.Fatal("插入失败:", err)
	}
	lastID, _ := result.LastInsertId() // 获取自增 ID
	rowsAff, _ := result.RowsAffected() // 影响行数
	fmt.Printf("插入成功: lastID=%d, 影响行数=%d\n", lastID, rowsAff)

	// 方式二：Prepared Statement（预编译，推荐批量插入时用）
	stmt, err := db.Prepare("INSERT INTO users(name, age) VALUES(?, ?)")
	if err != nil {
		log.Fatal("预编译失败:", err)
	}
	defer stmt.Close()

	// 复用同一个 stmt，只换参数
	stmt.Exec("李四", 33)
	stmt.Exec("王五", 25)
	stmt.Exec("赵六", 40)
	fmt.Println("批量插入完成（预编译语句）")

	// ===== 4. 查询单行（QueryRow）=====
	fmt.Println("\n--- QueryRow（查一行）---")
	var user User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id = ?", 1).Scan(&user.ID, &user.Name, &user.Age)
	// Scan 把查询结果填进变量 → 必须传 &（第 6 课讨论的 any 参数传指针）
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("没有找到该用户")
		} else {
			log.Fatal("查询失败:", err)
		}
	} else {
		fmt.Printf("查到用户: %+v\n", user)
	}

	// ===== 5. 查询多行（Query）=====
	fmt.Println("\n--- Query（查多行）---")
	rows, err := db.Query("SELECT id, name, age FROM users ORDER BY id")
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close() // ⚠️ 必须关闭！否则连接泄漏

	var allUsers []User
	for rows.Next() { // 逐行读取
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Age)
		if err != nil {
			log.Fatal("Scan 失败:", err)
		}
		allUsers = append(allUsers, u)
	}
	// 检查遍历过程中是否有错误
	if err = rows.Err(); err != nil {
		log.Fatal("遍历行出错:", err)
	}
	fmt.Println("所有用户:")
	for _, u := range allUsers {
		fmt.Printf("  ID=%d, Name=%s, Age=%d\n", u.ID, u.Name, u.Age)
	}

	// ===== 6. 更新数据（UPDATE）=====
	fmt.Println("\n--- UPDATE ---")
	result, err = db.Exec("UPDATE users SET age = ? WHERE name = ?", 29, "张三")
	if err != nil {
		log.Fatal("更新失败:", err)
	}
	rowsAff, _ = result.RowsAffected()
	fmt.Printf("更新成功: 影响行数=%d\n", rowsAff)

	// 验证更新结果
	var newAge int
	db.QueryRow("SELECT age FROM users WHERE name = ?", "张三").Scan(&newAge)
	fmt.Println("张三的新年龄:", newAge)

	// ===== 7. 删除数据（DELETE）=====
	fmt.Println("\n--- DELETE ---")
	result, err = db.Exec("DELETE FROM users WHERE name = ?", "赵六")
	if err != nil {
		log.Fatal("删除失败:", err)
	}
	rowsAff, _ = result.RowsAffected()
	fmt.Printf("删除成功: 影响行数=%d\n", rowsAff)

	// ===== 8. 事务（Transaction）=====
	fmt.Println("\n--- 事务 ---")
	err = transferAge(db)
	if err != nil {
		fmt.Println("事务失败:", err)
	} else {
		fmt.Println("事务成功")
	}

	// 查看事务后的结果
	rows2, _ := db.Query("SELECT id, name, age FROM users ORDER BY id")
	defer rows2.Close()
	fmt.Println("事务后的用户列表:")
	for rows2.Next() {
		var u User
		rows2.Scan(&u.ID, &u.Name, &u.Age)
		fmt.Printf("  ID=%d, Name=%s, Age=%d\n", u.ID, u.Name, u.Age)
	}

	fmt.Println("\n✅ 第 17 课全部完成！")
}

// transferAge 演示事务：张三年龄 -5，李四年龄 +5（模拟"转账"）
func transferAge(db *sql.DB) error {
	// 1. 开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}

	// 2. 在事务内执行多个操作
	_, err = tx.Exec("UPDATE users SET age = age - 5 WHERE name = ?", "张三")
	if err != nil {
		tx.Rollback() // 出错就回滚
		return fmt.Errorf("扣减失败: %w", err)
	}

	_, err = tx.Exec("UPDATE users SET age = age + 5 WHERE name = ?", "李四")
	if err != nil {
		tx.Rollback() // 出错就回滚
		return fmt.Errorf("增加失败: %w", err)
	}

	// 3. 全部成功，提交事务
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	fmt.Println("  张三 age-5, 李四 age+5 → 事务提交成功")
	return nil
}
