package main

import (
	"encoding/json"
	"fmt"
	"os"
)

/*
## 第 12 课：文件 IO + JSON

### 题 12.1 配置文件读写（⭐⭐）

定义一个 `Config` struct（含 Host, Port, Debug 等字段），实现：

 1. `SaveConfig(path string, cfg Config) error` — 写入 JSON 文件
 2. `LoadConfig(path string) (Config, error)` — 从 JSON 文件读取
    测试：保存一个配置，再读回来，对比是否一致。

---
*/
type Config struct {
	Host  string `json:"host"`
	Port  string `json:"port"`
	Debug bool   `json:"debug"`
}

func main() {
	config := Config{
		Host:  "123.23.23.45",
		Port:  "2206",
		Debug: false,
	}
	path := "exercises/12/test_write.text"
	err := SaveConfig(path, config)
	if err != nil {
		fmt.Println("异常了", err)
	}
	loadConfig, err := LoadConfig(path)
	if err != nil {
		fmt.Println("异常了", err)
	}
	fmt.Println("返回的数据", loadConfig)

}

func SaveConfig(path string, cfg Config) error {
	marshal, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化失败:%w", err)
	}
	err = os.WriteFile(path, []byte(marshal), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败:%w", err)
	}
	return nil
}

func LoadConfig(path string) (Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取文件报错：%w", err)
	}
	var cfg Config
	// 这个反序列化，内部会判断是否是指针的，不是的话会报运行时异常
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("反序列化失败：%w", err)
	}
	return cfg, nil
}
