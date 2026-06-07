# Go 学习工程

用于循序渐进地学习 Go，每个知识点放在独立子目录，互不干扰，都能单独运行。

## 目录结构

每个 `NN-主题/` 目录是一个独立的可运行程序，各自有自己的 `func main()`。

```
go-study/
├── go.mod
├── 01-hello/        # 第 1 课：Hello World
│   └── main.go
└── ...
```

## 如何运行

在工程根目录 `~/Project/go-study` 下执行：

```bash
# 运行某一课（注意是目录，不是文件）
go run ./01-hello

# 编译成二进制（可选）
go build -o bin/hello ./01-hello
```

## 新增一课

1. 新建目录，如 `02-variables/`
2. 在里面写 `main.go`，文件开头写 `package main`
3. 用 `go run ./02-variables` 运行

## 小贴士

- Go 强制要求：导入了不用的包、声明了不用的变量都会编译报错
- 格式化代码：`go fmt ./...`
- 同一个目录里只能有一个 `func main()`，所以每课用独立目录隔开
