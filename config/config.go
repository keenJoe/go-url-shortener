package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"name"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

var globalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 获取环境变量
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "dev" // 默认开发环境
	}

	// 构建配置文件路径
	configPath := filepath.Join(workDir, "config", fmt.Sprintf("config.%s.yaml", env))

	// 如果环境特定的配置文件不存在，使用默认配置文件
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = filepath.Join(workDir, "config", "config.yaml")
	}

	// 如果通过环境变量指定了配置文件路径，则使用指定的路径
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	log.Printf("正在加载配置文件: %s", configPath)

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证必要的配置项
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	// 保存到全局变量
	globalConfig = config

	// 打印关键配置信息（注意隐藏敏感信息）
	log.Printf("配置加载成功: Server.Port=%d, Database.Host=%s, Redis.Addr=%s",
		config.Server.Port,
		config.Database.Host,
		config.Redis.Addr)

	return config, nil
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server.port 未配置")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database.host 未配置")
	}
	if config.Database.DBName == "" {
		return fmt.Errorf("database.name 未配置")
	}
	return nil
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return globalConfig
}
