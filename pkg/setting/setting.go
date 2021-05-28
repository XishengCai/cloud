package setting

import (
	"cloud/pkg/common"
	"cloud/pkg/file"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/klog"
)

var (
	configFilePaths = [3]string{"conf/config.yaml", "../conf/config.yaml", "../../conf/config.yaml"}
	EnvConfig       *envConfig
	MongodbSetting  *mongoDB
	MysqlSetting    *mysql
	HttpSetting     *httpSetting
	RedisSetting    *redis
)

type envConfig struct {
	Title      string                  `yaml:"title"`
	ReleaseEnv string                  `yaml:"releaseEnv"`
	Version    string                  `yaml:"version"`
	Server     map[string]ServerConfig `yaml:"server"`
	RunMode    string                  `yaml:"runMode"`
}

type ServerConfig struct {
	MongoDB      mongoDB     `yaml:"mongoDB"`
	Mysql        mysql       `yaml:"mysql"`
	HttpSetting  httpSetting `yaml:"httpSetting"`
	RedisSetting redis       `yaml:"redisSetting"`
}

type mysql struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
	DateBase string `yaml:"database"`
}

type httpSetting struct {
	HttpPort     string        `yaml:"httpPort"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	RunMode      string        `yaml:"runMode"`
}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type mongoDB struct {
	Addresses string `yaml:"addresses"`
	Timeout   int    `yaml:"timeout"`
	Password  string `yaml:"password"`
	User      string `yaml:"user"`
	DateBase  string `yaml:"database"`
}

func SetUp(path string) {
	loadConfig(path)
}

func loadConfig(path string) {
	if path != "" {
		// 根据传参初始化环境配置
		if file.NotExists(path) {
			panic(fmt.Errorf("config file path: %s not exists", path))
		}

		loadConfigWithPath(path)
		return
	}
	for _, configPath := range configFilePaths {
		//从不同的层级目录初始化环境配置，直到有一次初始化成功后退出
		currentDir, _ := os.Getwd()
		klog.Infof("current directory: %s, load config: %s", currentDir, configPath)
		if file.NotExists(configPath) {
			continue
		}
		loadConfigWithPath(configPath)
		break
	}
	if EnvConfig == nil {
		panic("envConfig init fail")
	}
}

func loadConfigWithPath(configPath string) {
	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		klog.Error("read config file err: ", err)
		panic(err)
	}

	err = yaml.Unmarshal(config, &EnvConfig)
	if err != nil {
		panic(err)
	}
	klog.Infof("read Config: %s", configPath)

	releaseEnv := EnvConfig.ReleaseEnv
	serverConfig := EnvConfig.Server[releaseEnv]
	MongodbSetting = &serverConfig.MongoDB
	MysqlSetting = &serverConfig.Mysql
	HttpSetting = &serverConfig.HttpSetting
	RedisSetting = &serverConfig.RedisSetting
	HttpSetting.ReadTimeout = HttpSetting.ReadTimeout * time.Second
	HttpSetting.WriteTimeout = HttpSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	klog.Infof("config: %+v", common.PrettifyJson(EnvConfig.Server[releaseEnv], true))
}
