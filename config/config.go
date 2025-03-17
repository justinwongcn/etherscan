package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体用于解析 YAML 配置
type Config struct {
	Ethereum struct {
		NodeURL string `yaml:"node_url"`
	} `yaml:"ethereum"`
}

// LoadConfig 从 YAML 文件加载配置
func LoadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
} 