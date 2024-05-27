package gameConfig

import (
	"encoding/json"
	"log"
	"os"
)

const (
	TypeDurable        = 1 //耐久
	TypeCost           = 2
	TypeArmyTeams      = 3 //队伍数量
	TypeSpeed          = 4 //速度
	TypeDefense        = 5 //防御
	TypeStrategy       = 6 //谋略
	TypeForce          = 7 //攻击武力
	TypeConscriptTime  = 8 //征兵时间
	TypeReserveLimit   = 9 //预备役上限
	TypeUnkonw         = 10
	TypeHanAddition    = 11
	TypeQunAddition    = 12
	TypeWeiAddition    = 13
	TypeShuAddition    = 14
	TypeWuAddition     = 15
	TypeDealTaxRate    = 16 //交易税率
	TypeWood           = 17
	TypeIron           = 18
	TypeGrain          = 19
	TypeStone          = 20
	TypeTax            = 21 //税收
	TypeExtendTimes    = 22 //扩建次数
	TypeWarehouseLimit = 23 //仓库容量
	TypeSoldierLimit   = 24 //带兵数量
	TypeVanguardLimit  = 25 //前锋数量
)

type conditions struct {
	Type  int `json:"type"`
	Level int `json:"level"`
}

type facility struct {
	Title      string       `json:"title"`
	Des        string       `json:"des"`
	Name       string       `json:"name"`
	Type       int8         `json:"type"`
	Additions  []int8       `json:"additions"`
	Conditions []conditions `json:"conditions"`
	Levels     []fLevel     `json:"levels"`
}

type NeedRes struct {
	Decree int `json:"decree"`
	Grain  int `json:"grain"`
	Wood   int `json:"wood"`
	Iron   int `json:"iron"`
	Stone  int `json:"stone"`
	Gold   int `json:"gold"`
}

type fLevel struct {
	Level  int     `json:"level"`
	Values []int   `json:"values"`
	Need   NeedRes `json:"need"`
	Time   int     `json:"time"` //升级需要的时间
}

type conf struct {
	Name string
	Type int8
}

type facilityConf struct {
	Title     string `json:"title"`
	List      []conf `json:"list"`
	facilitys map[int8]*facility
}

var FacilityConf = &facilityConf{}

const facilityFile = "/config/game/facility/facility.json"
const facilityPath = "/config/game/facility/"

// Load 读取文件中的城池配置
func (f *facilityConf) Load() {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configFile := currentDir + facilityFile
	configPath := currentDir + facilityPath
	//参数  mssgserver.exe  D:/xxx
	length := len(os.Args)
	if length > 1 {
		dir := os.Args[1]
		if dir != "" {
			configFile = dir + facilityFile
			configPath = dir + facilityPath
		}
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, f)
	if err != nil {
		log.Println("json格式不正确，解析出错")
		panic(err)
	}
	f.facilitys = make(map[int8]*facility, len(f.List))

	files, err := os.ReadDir(configPath)
	if err != nil {
		log.Println("读取设施目录出错")
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name() == "facility.json" {
			continue
		}
		data, err := os.ReadFile(configPath + file.Name())
		if err != nil {
			log.Println("读取设施文件出错")
			panic(err)
		}
		fac := &facility{}
		err = json.Unmarshal(data, fac)
		if err != nil {
			log.Println("转json设施数据出错")
			panic(err)
		}
		f.facilitys[fac.Type] = fac
	}
}
