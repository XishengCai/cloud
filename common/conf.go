package common

import (
	"cloud/util"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

var configFilePaths = [3]string{"conf/config.yaml", "../conf/config.yaml", "../../conf/config.yaml"}
var envConfig *EnvConfig
var once sync.Once

type EnvConfig struct {
	Title      string                  `yaml:"title"`
	Env        string                  `yaml:"env"`
	Port       string                  `yaml:"port"`
	Version    string                  `yaml:"version"`
	Server     map[string]ServerConfig `yaml:"server"`
	ConfigPath string                  `yaml:"string"`
}

type ServerConfig struct {
	Mysql
	MongoDB
}
type Mysql struct {
	IP       string `yaml:"ip"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
	DateBase string `yaml:"database"`
}

type MongoDB struct {
	Addresses string `yaml:"addresses"`
	Timeout   int    `yaml:"timeout"`
	Password  string `yaml:"password"`
	User      string `yaml:"user"`
	DateBase  string `yaml:"database"`
}

//GetConf get toml conf
func GetConf(path string) *EnvConfig {
	if envConfig != nil {
		return envConfig
	}
	once.Do(func() {
		LoadConfig(path)
	})
	return envConfig
}

func LoadConfig(path string) {
	if path != "" {
		// 根据传参初始化环境配置
		if !util.FileExists(path) {
			panic(fmt.Errorf("config file path: %s not exists", path))
		}

		loadConfigWithPath(path)
		return
	}
	for _, configPath := range configFilePaths {
		//从不同的层级目录初始化环境配置，直到有一次初始化成功后退出
		currentDir, _ := os.Getwd()
		klog.Infof("current directory: %s, load config: %s", currentDir, configPath)
		if !util.FileExists(configPath) {
			continue
		}
		loadConfigWithPath(configPath)
		break
	}
	if envConfig == nil {
		panic("envConfig init fail")
	}
}

func loadConfigWithPath(configPath string) {
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		klog.Error("read config file err: ", err)
		panic(err)
	}
	err = yaml.Unmarshal(config, &envConfig)
	if err != nil {
		panic(err)
	}
	klog.Infof("read Config: %s", configPath)

	releaseEnv := envConfig.Env
	serverConfig := envConfig.Server[releaseEnv]
	// password, err := Base64Decode(serverConfig.Mysql.Password)
	// if err != nil {
	// 	panic(err)
	// }
	// serverConfig.Mysql.Password = string(password)
	envConfig.Server[releaseEnv] = serverConfig
}
