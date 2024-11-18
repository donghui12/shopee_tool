package shopee

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

// TestConfig 测试配置结构
type TestConfig struct {
    Phone     string `json:"phone"`
    Password  string `json:"password"`
    BaseURL   string `json:"base_url"`
    TestData  struct {
        ProductID int64 `json:"product_id"`
    } `json:"test_data"`
}

// setupTestClient 创建测试客户端
func setupTestClient(t *testing.T) (*Client, *TestConfig) {
    config := loadTestConfig(t)
    client := NewClient(
        WithBaseURL(config.BaseURL),
        WithTimeout(30*time.Second),
        WithRetry(3, 5*time.Second),
    )
    return client, config
}

// loadTestConfig 加载测试配置
func loadTestConfig(t *testing.T) *TestConfig {
    // 获取项目根目录
    rootDir := findProjectRoot(t)
    configPath := filepath.Join(rootDir, "configs", "test.json")
    
    data, err := os.ReadFile(configPath)
    if err != nil {
        t.Fatalf("读取测试配置失败: %v", err)
    }

    var config TestConfig
    if err := json.Unmarshal(data, &config); err != nil {
        t.Fatalf("解析测试配置失败: %v", err)
    }

    // 设置默认值
    if config.BaseURL == "" {
        config.BaseURL = BaseSellerURL
    }

    return &config
}

// findProjectRoot 查找项目根目录
func findProjectRoot(t *testing.T) string {
    dir, err := os.Getwd()
    if err != nil {
        t.Fatalf("获取当前目录失败: %v", err)
    }

    for {
        if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
            return dir
        }
        parent := filepath.Dir(dir)
        if parent == dir {
            t.Fatal("未找到项目根目录")
        }
        dir = parent
    }
}

// loginTestClient 登录测试客户端
func loginTestClient(t *testing.T, client *Client, config *TestConfig) {
    err := client.Login(config.Phone, config.Password, "")
    if err != nil {
        t.Fatalf("登录失败: %v", err)
    }
} 