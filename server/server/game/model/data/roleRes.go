package data

import "github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"

type RoleRes struct {
	Id     int `xorm:"id pk autoincr"`
	RId    int `xorm:"rid"`
	Wood   int `xorm:"wood"`
	Iron   int `xorm:"iron"`
	Stone  int `xorm:"stone"`
	Grain  int `xorm:"grain"`
	Gold   int `xorm:"gold"`
	Decree int `xorm:"decree"` //令牌
}

func (r *RoleRes) TableName() string {
	return "role_res"
}

// ToModel 模型转换
func (r *RoleRes) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	p.GoldYield = 100
	p.GrainYield = 100
	p.StoneYield = 100
	p.IronYield = 100
	p.WoodYield = 100
	p.DepotCapacity = 10000
	return p
}
