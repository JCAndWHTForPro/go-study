package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // 只注册驱动，不直接使用
)

/*
## 第 17 课：数据库

### 题 17.1 学生管理系统（⭐⭐）

用 `database/sql` 实现一个学生表的 CRUD：

- 建表（id, name, grade, score）
- 插入 5 个学生
- 查询分数大于 80 的学生
- 更新某个学生的分数
- 删除一个学生
- 事务：给所有学生加 5 分（原子操作）
*/

type student struct {
	name  string
	grade string
	score int
}

func main() {
	dsn := "root:root123@tcp(127.0.0.1:3306)/go_study?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("打开数据库失败", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("打开数据库链接失败", err)
	}
	db.SetMaxOpenConns(10) // 最大打开连接数
	db.SetMaxIdleConns(5)  // 最大空闲连接数

	createSql := `
		CREATE TABLE IF NOT EXISTS students (
			id    INT AUTO_INCREMENT PRIMARY KEY,
			name  VARCHAR(50)  NOT NULL,
			grade VARCHAR(20)  NOT NULL,
			score INT          NOT NULL DEFAULT 0
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`
	_, err = db.Exec(createSql)
	if err != nil {
		log.Fatal("创建表失败", err)
	}
	db.Exec("TRUNCATE TABLE students")
	prepare, _ := db.Prepare("INSERT INTO students (name, grade, score) VALUES (?, ?, ?)")
	defer prepare.Close()
	prepare.Exec("学生1", "一年级", 20)
	prepare.Exec("学生2", "二年级", 85)
	prepare.Exec("学生3", "一年级", 30)
	prepare.Exec("学生4", "三年级", 90)
	prepare.Exec("学生5", "六年级", 10)

	querySql := "select name,grade,score from students where score>80"
	query, err := db.Query(querySql)
	if err != nil {
		log.Fatal("查询报错", err)
	}
	defer query.Close()
	sts := []student{}
	for query.Next() {
		var st student
		query.Scan(&st.name, &st.grade, &st.score)
		sts = append(sts, st)
	}
	fmt.Println("查到的结果是", sts)

	deleteSql := "delete from students where name = '学生5'"
	exec, err := db.Exec(deleteSql)
	if err != nil {
		log.Fatal("删除失败", err)
	}
	affected, _ := exec.RowsAffected()
	fmt.Println("删除影响的行数是：", affected)

	// 更新某个学生的分数
	updateOneSql := "UPDATE students SET score = ? WHERE name = ?"
	result, err := db.Exec(updateOneSql, 100, "学生2")
	if err != nil {
		log.Fatal("更新学生分数失败", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Println("更新学生2的分数，影响行数：", rowsAffected)

	// 事务操作：给所有学生加5分
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("开启事务失败", err)
	}

	txUpdateSql := "UPDATE students SET score = score + 5 WHERE 1=1"
	_, err = tx.Exec(txUpdateSql)
	if err != nil {
		fmt.Println("更新失败，回滚事务", err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("提交事务失败", err)
		tx.Rollback()
		return
	}
	fmt.Println("事务提交成功，所有学生加5分")

}
