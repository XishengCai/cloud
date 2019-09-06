package common

import (
	"cloud/util"
	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
	"os"
)

var tomlFilePath = "conf/cloud_config.toml"
var tomlConf *TomlConfig

type TomlConfig struct {
	Title      string                  `toml:"title"`
	Env        string                  `toml:"env"`
	Port       string                  `toml:"port"`
	Version    string                  `toml:"version"`
	Server     map[string]ServerConfig `toml:"server"`
	ConfigPath string                  `toml:"string"`
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

//GetConf get toml conf
func GetConf(config string) *TomlConfig {
	if config != "" {
		tomlFilePath = config
	}
	if !util.IsFile(config){
		log.Errorf("toml config  file:%s not exits", config)
		os.Exit(1)
	}
	loadConfig()
	return tomlConf
}

func loadConfig() {

	_, err := toml.DecodeFile(tomlFilePath, &tomlConf)
	if err != nil {
		panic(err)
	}
	log.Infof("read tomlConfig: %s", tomlFilePath)

	serverConfig := tomlConf.Server[tomlConf.Env]
	password, err := Base64Decode(tomlConf.Server[tomlConf.Env].Mysql.Password)
	if err != nil {
		panic(err)
	}
	serverConfig.Mysql.Password = string(password)
	tomlConf.Server[tomlConf.Env] = serverConfig

}
