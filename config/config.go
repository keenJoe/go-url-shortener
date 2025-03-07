package config

type Config struct {
	DB struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

var GlobalConfig Config

func LoadConfig() *Config {
	// 这里可以从配置文件加载配置
	// 示例中使用硬编码配置
	GlobalConfig.DB.Host = "localhost"
	GlobalConfig.DB.Port = "3306"
	GlobalConfig.DB.User = "root"
	GlobalConfig.DB.Password = "your_password"
	GlobalConfig.DB.DBName = "your_db_name"

	return &GlobalConfig
}
