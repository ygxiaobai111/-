package gameConfig

import (
	"encoding/json"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/global"
	"io/ioutil"
	"log"
	"os"
)

type mapData struct {
	Width  int     `json:"w"`
	Height int     `json:"h"`
	List   [][]int `json:"list"`
}

type NationalMap struct {
	MId   int  `xorm:"mid"`
	X     int  `xorm:"x"`
	Y     int  `xorm:"y"`
	Type  int8 `xorm:"type"`
	Level int8 `xorm:"level"`
}

const (
	MapBuildSysFortress = 50 //系统要塞
	MapBuildSysCity     = 51 //系统城市
	MapBuildFortress    = 56 //玩家要塞
)

var MapRes = &mapRes{
	Confs:    make(map[int]NationalMap),
	SysBuild: make(map[int]NationalMap),
}

type mapRes struct {
	Confs    map[int]NationalMap
	SysBuild map[int]NationalMap
}

const mapFile = "/config/game/map.json"

func (m *mapRes) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + mapFile

	//参数  mssgserver.exe  D:/xxx
	length := len(os.Args)
	if length > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + mapFile
		}
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	mapData := &mapData{}
	err = json.Unmarshal(data, mapData)
	if err != nil {
		log.Println("json格式不正确，解析出错")
		panic(err)
	}
	global.MapWith = mapData.Width
	global.MapHeight = mapData.Height
	log.Println("list len", len(mapData.List))
	for index, v := range mapData.List {
		t := int8(v[0]) //type
		l := int8(v[1]) //level
		nm := NationalMap{
			X:     index % global.MapWith,
			Y:     index / global.MapHeight,
			Type:  t,
			Level: l,
			MId:   index,
		}
		m.Confs[index] = nm
		if t == MapBuildSysCity || t == MapBuildSysFortress {
			m.SysBuild[index] = nm
		}
	}
}
