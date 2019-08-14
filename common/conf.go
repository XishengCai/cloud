package common

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"runtime"
	"strings"
)

var tomlFilePath = "conf/cloud_config.toml"
var tomlConf *TomlConfig

type TomlConfig struct {
	Title   string                  `toml:"title"`
	Env     string                  `toml:"env"`
	Port    int                     `toml:"port"`
	Version string                  `toml:"version"`
	Server  map[string]ServerConfig `toml:"server"`
}

type ServerConfig struct {
	Mysql MysqlConfig `toml:"mysql"`
}

type MysqlConfig struct {
	IP       string `toml:"ip"`
	Password string `toml:"password"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	DateBase string `toml:"database"`
}

func init() {
	loadConfig()
	//s := make(chan os.Signal, 1)
	//signal.Notify(s, syscall.SIGUSR1)
	//go func() {
	//	for {
	//		<-s
	//		log.Info("Reloaded config")
	//		loadConfig()
	//	}
	//}()
}

//GetConf get toml conf
func GetConf() *TomlConfig {
	return tomlConf
}

func loadConfig() {
	_, filename, _, _ := runtime.Caller(1)
	fmt.Println("filename: ", filename)
	paths := strings.Split(filename,"/")
	paths = paths[0:len(paths)-2]
	paths = append(paths, tomlFilePath)
	confFileFullPath := strings.Join(paths,"/")
	fmt.Println("confFileFullPath: ", confFileFullPath)
	_, err := toml.DecodeFile(confFileFullPath, &tomlConf)
	if err != nil {
		panic(err)
	}
	fmt.Println("read tomlConfig: ", tomlConf)

	serverConfig := tomlConf.Server[tomlConf.Env]
	password, err := Base64Decode(tomlConf.Server[tomlConf.Env].Mysql.Password)
	if err != nil {
		panic(err)
	}
	serverConfig.Mysql.Password = string(password)
	tomlConf.Server[tomlConf.Env] = serverConfig
}
