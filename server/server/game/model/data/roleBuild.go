package data

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"time"
)

const (
	MapBuildSysFortress = 50 //系统要塞
	MapBuildSysCity     = 51 //系统城市
	MapBuildFortress    = 56 //玩家要塞
)

type MapRoleBuild struct {
	Id         int       `xorm:"id pk autoincr"`
	RId        int       `xorm:"rid"`
	Type       int8      `xorm:"type"`
	Level      int8      `xorm:"level"`
	OPLevel    int8      `xorm:"op_level"` //操作level
	X          int       `xorm:"x"`
	Y          int       `xorm:"y"`
	Name       string    `xorm:"name"`
	Wood       int       `xorm:"-"`
	Iron       int       `xorm:"-"`
	Stone      int       `xorm:"-"`
	Grain      int       `xorm:"-"`
	Defender   int       `xorm:"-"`
	CurDurable int       `xorm:"cur_durable"`
	MaxDurable int       `xorm:"max_durable"`
	OccupyTime time.Time `xorm:"occupy_time"`
	EndTime    time.Time `xorm:"end_time"` //建造或升级完的时间
	GiveUpTime int64     `xorm:"giveUp_time"`
}

func (m *MapRoleBuild) TableName() string {
	return "map_role_build"
}

func (m *MapRoleBuild) ToModel() interface{} {
	p := model.MapRoleBuild{}
	p.RNick = "111"
	p.UnionId = 0
	p.UnionName = ""
	p.ParentId = 0
	p.X = m.X
	p.Y = m.Y
	p.Type = m.Type
	p.RId = m.RId
	p.Name = m.Name

	p.OccupyTime = m.OccupyTime.UnixNano() / 1e6
	p.GiveUpTime = m.GiveUpTime * 1000
	p.EndTime = m.EndTime.UnixNano() / 1e6

	p.CurDurable = m.CurDurable
	p.MaxDurable = m.MaxDurable
	p.Defender = m.Defender
	p.Level = m.Level
	p.OPLevel = m.OPLevel
	return p
}
