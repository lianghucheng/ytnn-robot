package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/name5566/leaf/log"
)

// Config 配置类型
type PublicConfig struct {
	CfgGameInfo CfgGameInfo
}

type CfgGameInfo struct {
	HallAddress string
	RobotIndex  int
	LogPath     string
}

// Opts Config 默认配置
var opts *PublicConfig

// ParseToml 解析配置文件
func init() {
	ReloadConfig()
}

func ReloadConfig() {
	opts = &PublicConfig{}
	_, err := toml.DecodeFile("conf/nn-robot.toml", opts)
	fmt.Printf("%+v\n", opts.CfgGameInfo)
	if err != nil {
		log.Fatal("配置文件解析错误:%s", err)
	}
}

func GetCfgGameInfo() *CfgGameInfo {
	return &opts.CfgGameInfo
}
