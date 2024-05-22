package config

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/Unknwon/goconfig"
)

const configFile = "/config/conf.ini"

var File *goconfig.ConfigFile

// 加载此文件的时候 会先走初始化方法
func init() {
	//拿到当前的程序的目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}
	//参数  main.exe  D:/xxx/conf.ini
	path := flag.String("conf_path", "", "product config  path")
	flag.Parse()
	if *path != "" {
		configPath = *path
	}

	//文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}
