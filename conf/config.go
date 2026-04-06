package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Server          ServerConfig     `mapstructure:"server"`
	Database        DatabaseConfig   `mapstructure:"database"`
	FileServer      FileServerConfig `mapstructure:"file_server"`
	ThumbnailConfig ThumbnailOption  `mapstructure:"thumbnail"`
	PreviewConfig   ThumbnailOption  `mapstructure:"preview"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Path   string `mapstructure:"path"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

type FileOperationConfig struct {
	Enabled      bool     `mapstructure:"enabled"`
	MaxSize      int64    `mapstructure:"max_size"`      // 字节为单位，0表示无限制
	AllowedTypes []string `mapstructure:"allowed_types"` // 允许的文件类型，如 [".jpg", ".png", ".pdf"]
}

type FileServerConfig struct {
	BaseURL     string `mapstructure:"base_url"`     // 基础URL，本地存储时留空
	APIKey      string `mapstructure:"api_key"`      // API密钥
	AppID       string `mapstructure:"app_id"`       //
	SpaceID     string `mapstructure:"space_id"`     //
	StoragePath string `mapstructure:"storage_path"` // 本地存储路径，如 ./data/uploads
}

type ThumbnailOption struct {
	Enabled bool   `mapstructure:"enabled"`
	Width   int    `mapstructure:"width"`
	Height  int    `mapstructure:"height"`
	Mode    string `mapstructure:"mode"` // fit, fill, stretch 等
	Quality int    `mapstructure:"quality"`
}

var Config *AppConfig

func LoadConfig(configPath string) error {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 设置默认值
	v.SetDefault("server.port", 9000)
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.mode", "debug")
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.path", "./data/artwork.db")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	Config = &AppConfig{}
	if err := v.Unmarshal(Config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

func GetDSN() string {
	return Config.Database.Path
}
