package gameConfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type cfg struct {
	Type     int8   `json:"type"`
	Name     string `json:"name"`
	Level    int8   `json:"level"`
	Grain    int    `json:"grain"`
	Wood     int    `json:"wood"`
	Iron     int    `json:"iron"`
	Stone    int    `json:"stone"`
	Durable  int    `json:"durable"`
	Defender int    `json:"defender"`
}

type mapBuildConf struct {
	Title  string `json:"title"`
	Cfg    []cfg  `json:"cfg"`
	cfgMap map[int8][]cfg
}

var MapBuildConf = &mapBuildConf{}

const mapBuildConfFile = "/config/game/map_build.json"

func (m *mapBuildConf) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + mapBuildConfFile

	//参数  mssgserver.exe  D:/xxx
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + mapBuildConfFile
		}
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, m)
	if err != nil {
		log.Println("json格式不正确，解析出错")
		panic(err)
	}
}
