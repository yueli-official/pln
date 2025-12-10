package conf

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const ApiKeyFile = "./data/apikey.txt"

var ApiKeys = make(map[string]bool)

func AddAPIKey(key string) {
	ApiKeys[key] = true
}

func IsValidAPIKey(key string) bool {
	return ApiKeys[key]
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 4) // 4字节 = 8位hex字符串
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func InitAPIKey() string {
	// 读取已存在的 api key 文件
	data, err := os.ReadFile(ApiKeyFile)
	if err == nil {
		key := string(data)
		AddAPIKey(key)
		return key
	}

	// 文件不存在则创建目录
	if errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(filepath.Dir(ApiKeyFile), 0755)
	} else {
		// 其他错误仍需要处理，但不 panic
		fmt.Println("读取 API key 文件时出错:", err)
	}

	// 生成新的 key
	key, err := GenerateAPIKey()
	if err != nil {
		fmt.Println("生成 API key 失败:", err)
		return ""
	}

	// 写文件
	err = os.WriteFile(ApiKeyFile, []byte(key), 0644)
	if err != nil {
		fmt.Println("写入 API key 文件失败:", err)
		return ""
	}

	AddAPIKey(key)
	return key
}
