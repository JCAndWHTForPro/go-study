// 第 18~19 课：GORM（Go 最流行的 ORM 框架）
// 运行方式：go run ./18~19-gorm
//
// 前置条件：
//  1. 本地 MySQL 已启动，端口 3306，用户 root，密码 root123
//  2. 已有数据库 go_study
//  3. 已安装：go get gorm.io/gorm gorm.io/driver/mysql
//
// GORM vs database/sql（第 17 课）：
//
//	database/sql → 手写 SQL，灵活但繁琐（类似 Java JDBC）
//	GORM         → 用 Go 代码操作数据库，自动生成 SQL（类似 Java MyBatis-Plus / JPA）
//
// ============================================================
// 【总结】GORM 核心概念
// ------------------------------------------------------------
//
// ▶ 连接：
//
//	gorm.Open(mysql.Open(dsn), &gorm.Config{})  → 返回 *gorm.DB
//	和 database/sql 一样底层是连接池
//
// ▶ 模型定义（struct → 表）：
//
//	type User struct {
//	    gorm.Model            // 内嵌：自动添加 ID、CreatedAt、UpdatedAt、DeletedAt
//	    Name string           // 字段名 → 列名：Name → name（自动蛇形转换）
//	    Age  int
//	}
//	gorm.Model 包含：
//	  ID        uint           `gorm:"primaryKey"`
//	  CreatedAt time.Time      // 创建时自动填充
//	  UpdatedAt time.Time      // 更新时自动填充
//	  DeletedAt gorm.DeletedAt // 软删除：删除不真删，只标记时间
//
// ▶ 自动建表：
//
//	db.AutoMigrate(&User{})  → 根据 struct 自动创建/更新表结构
//
// ▶ CRUD 操作：
//
//	db.Create(&user)                          → INSERT
//	db.First(&user, id)                       → SELECT ... WHERE id=? LIMIT 1
//	db.Find(&users)                           → SELECT *（查全部）
//	db.Where("name = ?", "张三").Find(&users) → SELECT ... WHERE name='张三'
//	db.Model(&user).Update("age", 30)         → UPDATE ... SET age=30
//	db.Model(&user).Updates(User{Age:30})     → UPDATE（多字段）
//	db.Delete(&user, id)                       → 软删除（UPDATE deleted_at）
//	db.Unscoped().Delete(&user, id)            → 硬删除（真正 DELETE）
//
// ▶ 链式调用（Builder 模式）：
//
//	db.Where(...).Order(...).Limit(...).Find(&results)
//	每个方法返回 *gorm.DB，可以一直链下去
//
// ▶ 自定义表名（三种方式）：
//
//	方式一（推荐）：给 struct 实现 TableName() 方法
//	  func (Product) TableName() string { return "my_products" }
//	  → 隐式接口，实现了方法就自动生效，不需要 implements
//
//	方式二：查询时临时指定
//	  db.Table("my_custom_table").Find(&products)
//
//	方式三：全局配置命名策略
//	  gorm.Open(dsn, &gorm.Config{
//	      NamingStrategy: schema.NamingStrategy{
//	          TablePrefix:   "t_",   // 所有表加前缀 → t_products
//	          SingularTable: true,   // 不用复数 → product（不加 s）
//	      },
//	  })
//
//	默认规则：struct 名 → 复数 + 蛇形（Product → products, UserOrder → user_orders）
//
// ▶ struct tag：
//
//	`gorm:"column:user_name"`     → 指定列名
//	`gorm:"type:varchar(100)"`    → 指定列类型
//	`gorm:"not null"`             → 非空约束
//	`gorm:"uniqueIndex"`          → 唯一索引
//	`gorm:"default:18~19"`           → 默认值
//
// ▶ 软删除：
//
//	GORM 默认使用软删除——Delete 只是设置 deleted_at 字段
//	查询时自动过滤 deleted_at IS NOT NULL 的记录
//	想查被软删的记录：db.Unscoped().Find(&users)
//	想真删：db.Unscoped().Delete(&user, id)
//
// ▶ 多表关联（重要！）：
//
//	定义关联：通过 struct 字段名约定，不用配置文件
//	  type Author struct {
//	      gorm.Model
//	      Name  string
//	      Books []Book     // Has Many：一个作者有多本书
//	  }
//	  type Book struct {
//	      gorm.Model
//	      Title    string
//	      AuthorID uint    // 外键：字段名 = 关联类型名+ID（自动推断）
//	      Author   Author  // Belongs To：一本书属于一个作者
//	  }
//
//	查询关联数据：
//	  db.Preload("Books").Find(&authors)            → 预加载：查作者时自动查书
//	  db.Preload("Author").First(&book)              → 反向：查书时自动查作者
//	  db.Preload("Books", "title LIKE ?", "%x%")     → 条件预加载
//	  db.Joins("Author").Where("authors.name=?", x)  → JOIN 查询
//
//	关联类型：
//	  Has Many    → Author 有 []Book        一对多
//	  Belongs To  → Book 有 Author+AuthorID  属于
//	  Has One     → User 有 Profile          一对一
//	  Many2Many   → User 有 []Role `gorm:"many2many:user_roles"`  多对多
//
// ▶ Preload vs Joins：
//
//	Preload → 额外发一条 SELECT（两条 SQL），适合加载关联数据展示
//	Joins   → 用 SQL JOIN（一条 SQL），适合用关联表的字段做 WHERE 过滤
//
// ▶ 对比 Java：
//
//	Java MyBatis-Plus               GORM
//	@TableName("users")             表名自动从 struct 名推断（Users → users）
//	@TableId(type=AUTO)             gorm.Model 自带 ID
//	baseMapper.selectById(1)        db.First(&user, 1)
//	baseMapper.insert(user)         db.Create(&user)
//	lambdaQuery().eq("name","张三") db.Where("name = ?", "张三")
//	@TableLogic                     gorm.DeletedAt（软删除）
//	逻辑删除自动过滤               同样自动过滤
//	@TableField(exist=false)+嵌套查 Preload("关联名")
//	XML 里写 JOIN                   db.Joins("关联名").Where(...)
//
// ============================================================
package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ===== 1. 模型定义 =====
// Product 商品表（演示自定义字段标签）
type Product struct {
	gorm.Model        // 内嵌 ID + CreatedAt + UpdatedAt + DeletedAt
	Name       string `gorm:"type:varchar(100);not null"` // 商品名
	Price      int    `gorm:"not null"`                   // 价格（分）
	Stock      int    `gorm:"default:0"`                  // 库存，默认 0
}

// ===== 关联模型定义 =====
// GORM 通过 struct 的字段名约定来推断关联关系：
//   外键字段命名规则：关联类型名 + ID → 如 AuthorID 对应 Author
//
// ▶ 一对多（Has Many）：一个作者有多本书
//   Author struct 里有 []Book → 一对多
//   Book struct 里有 AuthorID → 外键（自动推断）
//
// ▶ 属于（Belongs To）：一本书属于一个作者
//   Book struct 里有 Author + AuthorID → 属于关系

// Author 作者表
type Author struct {
	gorm.Model
	Name  string `gorm:"type:varchar(50);not null"`
	Books []Book // 一对多：一个作者有多本书（Has Many）
}

// Book 书籍表
type Book struct {
	gorm.Model
	Title    string `gorm:"type:varchar(100);not null"`
	AuthorID uint   // 外键：GORM 自动推断 → 关联 authors 表的 id
	Author   Author // 属于关系（Belongs To）：通过 AuthorID 关联
}

func main() {
	// ===== 2. 连接数据库 =====
	dsn := "root:root123@tcp(127.0.0.1:3306)/go_study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	fmt.Println("✅ GORM 连接成功")

	// ===== 3. 自动建表（AutoMigrate）=====
	// 根据 struct 自动创建表，已有表则自动新增字段（不会删字段）
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal("自动建表失败:", err)
	}
	fmt.Println("✅ 表 products 自动创建/更新成功")

	// 清空表数据（硬删除所有记录，包括软删除的）
	db.Exec("TRUNCATE TABLE products")

	// ===== 4. 创建（INSERT）=====
	fmt.Println("\n--- CREATE ---")

	// 单条插入
	p1 := Product{Name: "Go 语言圣经", Price: 5900, Stock: 100}
	result := db.Create(&p1)
	// Create 会自动填充 p1.ID、CreatedAt、UpdatedAt
	fmt.Printf("插入成功: ID=%d, 影响行数=%d\n", p1.ID, result.RowsAffected)

	// 批量插入
	products := []Product{
		{Name: "机械键盘", Price: 29900, Stock: 50},
		{Name: "显示器", Price: 199900, Stock: 20},
		{Name: "鼠标垫", Price: 2900, Stock: 200},
	}
	result = db.Create(&products)
	fmt.Printf("批量插入: 影响行数=%d\n", result.RowsAffected)
	for _, p := range products {
		fmt.Printf("  ID=%d, Name=%s\n", p.ID, p.Name)
	}

	// ===== 5. 查询（SELECT）=====
	fmt.Println("\n--- READ ---")

	// 查一条：按主键查
	var product Product
	db.First(&product, 1) // SELECT * FROM products WHERE id=1 LIMIT 1
	fmt.Printf("First(1): %s, ¥%.2f\n", product.Name, float64(product.Price)/100)

	// 查一条：按条件查
	var found Product
	db.Where("name = ?", "机械键盘").First(&found)
	fmt.Printf("Where(机械键盘): ID=%d, Stock=%d\n", found.ID, found.Stock)

	// 查全部
	var allProducts []Product
	db.Find(&allProducts) // SELECT * FROM products
	fmt.Println("Find 全部:")
	for _, p := range allProducts {
		fmt.Printf("  ID=%d, %s, ¥%.2f, 库存=%d\n", p.ID, p.Name, float64(p.Price)/100, p.Stock)
	}

	// 条件查询 + 排序 + 限制
	var expensive []Product
	db.Where("price > ?", 10000).Order("price desc").Limit(2).Find(&expensive)
	fmt.Println("价格>100元，按价格降序，前2条:")
	for _, p := range expensive {
		fmt.Printf("  %s, ¥%.2f\n", p.Name, float64(p.Price)/100)
	}

	// 只查特定字段
	var names []string
	db.Model(&Product{}).Pluck("name", &names) // 只取 name 列
	fmt.Println("所有商品名:", names)

	// 统计
	var count int64
	db.Model(&Product{}).Where("stock > ?", 30).Count(&count)
	fmt.Println("库存>30 的商品数量:", count)

	// ===== 6. 更新（UPDATE）=====
	fmt.Println("\n--- UPDATE ---")

	// 方式一：更新单个字段
	db.Model(&Product{}).Where("id = ?", 1).Update("stock", 88)
	fmt.Println("更新 ID=1 的库存为 88")

	// 方式二：更新多个字段（用 struct）
	// ⚠️ struct 更新会忽略零值字段！（0, "", false 不会被更新）
	db.Model(&Product{}).Where("id = ?", 2).Updates(Product{Price: 25900, Stock: 45})
	fmt.Println("更新 ID=2: Price=25900, Stock=45")

	// 方式三：用 map 更新（不会忽略零值）
	db.Model(&Product{}).Where("id = ?", 3).Updates(map[string]any{
		"stock": 0, // 用 map 可以把 stock 更新为 0
	})
	fmt.Println("更新 ID=3: Stock=0（用 map 更新零值）")

	// 验证更新结果
	var updated Product
	db.First(&updated, 1)
	fmt.Printf("验证 ID=1: Stock=%d\n", updated.Stock)

	// ===== 7. 删除（DELETE）=====
	fmt.Println("\n--- DELETE ---")

	// 软删除：只设置 deleted_at，不真删数据
	db.Delete(&Product{}, 4) // DELETE（软删除）鼠标垫
	fmt.Println("软删除 ID=4")

	// 查询全部（软删除的记录自动被过滤）
	var afterDelete []Product
	db.Find(&afterDelete)
	fmt.Printf("软删除后查到 %d 条（ID=4 被过滤了）\n", len(afterDelete))

	// 查包含软删除的记录
	var withDeleted []Product
	db.Unscoped().Find(&withDeleted)
	fmt.Printf("Unscoped 查到 %d 条（包含软删除的）\n", len(withDeleted))

	// 硬删除：真正从数据库删除
	db.Unscoped().Delete(&Product{}, 4) // 真删
	fmt.Println("硬删除 ID=4")

	var afterHardDelete []Product
	db.Unscoped().Find(&afterHardDelete)
	fmt.Printf("硬删除后 Unscoped 查到 %d 条\n", len(afterHardDelete))

	// ===== 8. 事务 =====
	fmt.Println("\n--- 事务 ---")
	err = db.Transaction(func(tx *gorm.DB) error {
		// 在事务内操作（用 tx 而不是 db）
		if err := tx.Model(&Product{}).Where("id = ?", 1).Update("stock", gorm.Expr("stock - ?", 10)).Error; err != nil {
			return err // 返回错误自动回滚
		}
		if err := tx.Model(&Product{}).Where("id = ?", 2).Update("stock", gorm.Expr("stock + ?", 10)).Error; err != nil {
			return err // 返回错误自动回滚
		}
		return nil // 返回 nil 自动提交
	})
	if err != nil {
		fmt.Println("事务失败:", err)
	} else {
		fmt.Println("事务成功: ID=1 库存-10, ID=2 库存+10")
	}

	// 查看最终结果
	fmt.Println("\n--- 最终结果 ---")
	var finalProducts []Product
	db.Find(&finalProducts)
	for _, p := range finalProducts {
		fmt.Printf("  ID=%d, %s, ¥%.2f, 库存=%d\n", p.ID, p.Name, float64(p.Price)/100, p.Stock)
	}

	// ===== 9. 多表关联 =====
	fmt.Println("\n--- 多表关联 ---")
	demonstrateAssociations(db)

	fmt.Println("\n✅ 第 18~19 课全部完成！")
}

// demonstrateAssociations 演示 GORM 多表关联查询
func demonstrateAssociations(db *gorm.DB) {
	// 自动建表（GORM 会自动创建外键约束）
	db.AutoMigrate(&Author{}, &Book{})
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE books")
	db.Exec("TRUNCATE TABLE authors")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// ===== 9.1 创建关联数据 =====
	// 方式一：先创建作者，再创建书
	author1 := Author{Name: "鲁迅"}
	db.Create(&author1)
	db.Create(&Book{Title: "呐喊", AuthorID: author1.ID})
	db.Create(&Book{Title: "彷徨", AuthorID: author1.ID})

	// 方式二：一次性创建作者和书（嵌套创建）
	author2 := Author{
		Name: "老舍",
		Books: []Book{
			{Title: "骆驼祥子"},
			{Title: "茶馆"},
			{Title: "四世同堂"},
		},
	}
	db.Create(&author2) // GORM 自动填充 Book 的 AuthorID
	fmt.Println("✅ 关联数据创建完成")

	// ===== 9.2 查询：Preload（预加载关联数据）=====
	// 不加 Preload：只查 author，Books 字段是空的
	var authorOnly Author
	db.First(&authorOnly, author1.ID)
	fmt.Printf("不加 Preload: %s, Books=%v（空的！）\n", authorOnly.Name, authorOnly.Books)

	// 加 Preload：自动查关联的 books 表并填充
	var authorWithBooks Author
	db.Preload("Books").First(&authorWithBooks, author1.ID)
	fmt.Printf("加 Preload: %s, 有 %d 本书:\n", authorWithBooks.Name, len(authorWithBooks.Books))
	for _, b := range authorWithBooks.Books {
		fmt.Printf("  - %s\n", b.Title)
	}

	// ===== 9.3 查询：查所有作者 + 各自的书 =====
	var allAuthors []Author
	db.Preload("Books").Find(&allAuthors)
	fmt.Println("\n所有作者及其作品:")
	for _, a := range allAuthors {
		fmt.Printf("  %s（%d 本）:\n", a.Name, len(a.Books))
		for _, b := range a.Books {
			fmt.Printf("    - %s\n", b.Title)
		}
	}

	// ===== 9.4 反向查询：从书查作者（Belongs To）=====
	var book Book
	db.Preload("Author").First(&book, "title = ?", "茶馆")
	fmt.Printf("\n《%s》的作者是: %s\n", book.Title, book.Author.Name)

	// ===== 9.5 条件预加载：只加载标题含"驼"的书 =====
	var authorFiltered Author
	db.Preload("Books", "title LIKE ?", "%驼%").First(&authorFiltered, author2.ID)
	fmt.Printf("\n%s 的书（标题含'驼'）:\n", authorFiltered.Name)
	for _, b := range authorFiltered.Books {
		fmt.Printf("  - %s\n", b.Title)
	}

	// ===== 9.6 Joins 查询（类似 SQL JOIN）=====
	// 当你想用 WHERE 条件过滤关联表时，用 Joins
	var booksOfLuxun []Book
	db.Joins("Author").Where("Author.name = ?", "鲁迅").Find(&booksOfLuxun)
	fmt.Println("\nJoins 查鲁迅的书:")
	for _, b := range booksOfLuxun {
		fmt.Printf("  - %s（作者: %s）\n", b.Title, b.Author.Name)
	}
}
