package cnf

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

// Conf 与配置文件内容以一一对应
type Conf struct {
	Env string `yaml:"env"`
	DB  struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Schema   string `yaml:"schema"`
		Username string `yaml:"username"`
		Password string `json:"password"`
	}
}

const (
	confRoot   = "conf"
	confEnv    = "env"
	confSuffix = ".yml"
)

// NewConf 初始化配置
func NewConf() (*Conf, error) {
	// 打开文件 env.yml
	file, err := os.Open(confRoot + string(os.PathSeparator) + confEnv + confSuffix)
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 初始化 Conf
	conf := new(Conf)

	// 解析 yaml
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}

	// 根据 profile 构建配置文件路径
	path := confRoot + string(os.PathSeparator) + conf.Env + confSuffix

	// 打开文件
	file, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// 解析 yaml
	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
