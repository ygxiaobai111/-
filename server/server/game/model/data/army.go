package data

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"time"
)

const (
	ArmyCmdIdle        = 0 //空闲
	ArmyCmdAttack      = 1 //攻击
	ArmyCmdDefend      = 2 //驻守
	ArmyCmdReclamation = 3 //屯垦
	ArmyCmdBack        = 4 //撤退
	ArmyCmdConscript   = 5 //征兵
	ArmyCmdTransfer    = 6 //调动
)

const (
	ArmyStop    = 0
	ArmyRunning = 1
)

// 军队
type Army struct {
	Id                 int        `xorm:"id pk autoincr"`
	RId                int        `xorm:"rid"`
	CityId             int        `xorm:"cityId"`
	Order              int8       `xorm:"order"`
	Generals           string     `xorm:"generals"`
	Soldiers           string     `xorm:"soldiers"`
	ConscriptTimes     string     `xorm:"conscript_times"` //征兵结束时间，json数组
	ConscriptCnts      string     `xorm:"conscript_cnts"`  //征兵数量，json数组
	Cmd                int8       `xorm:"cmd"`
	FromX              int        `xorm:"from_x"`
	FromY              int        `xorm:"from_y"`
	ToX                int        `xorm:"to_x"`
	ToY                int        `xorm:"to_y"`
	Start              time.Time  `json:"-"xorm:"start"`
	End                time.Time  `json:"-"xorm:"end"`
	State              int8       `xorm:"-"` //状态:0:running,1:stop
	GeneralArray       []int      `json:"-" xorm:"-"`
	SoldierArray       []int      `json:"-" xorm:"-"`
	ConscriptTimeArray []int64    `json:"-" xorm:"-"`
	ConscriptCntArray  []int      `json:"-" xorm:"-"`
	Gens               []*General `json:"-" xorm:"-"`
	CellX              int        `json:"-" xorm:"-"`
	CellY              int        `json:"-" xorm:"-"`
}

func (a *Army) TableName() string {
	return "army"
}

func (a *Army) ToModel() interface{} {
	p := model.Army{}
	p.CityId = a.CityId
	p.Id = a.Id
	p.UnionId = 0
	p.Order = a.Order
	p.Generals = a.GeneralArray
	p.Soldiers = a.SoldierArray
	p.ConTimes = a.ConscriptTimeArray
	p.ConCnts = a.ConscriptCntArray
	p.Cmd = a.Cmd
	p.State = a.State
	p.FromX = a.FromX
	p.FromY = a.FromY
	p.ToX = a.ToX
	p.ToY = a.ToY
	p.Start = a.Start.Unix()
	p.End = a.End.Unix()
	return p
}
