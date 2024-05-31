package data

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"log"
)

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

func (r *RoleRes) ToModel() interface{} {
	p := model.RoleRes{}
	p.Gold = r.Gold
	p.Grain = r.Grain
	p.Stone = r.Stone
	p.Iron = r.Iron
	p.Wood = r.Wood
	p.Decree = r.Decree

	yield := GetYield(r.RId)
	p.GoldYield = yield.Gold
	p.GrainYield = yield.Grain
	p.StoneYield = yield.Stone
	p.IronYield = yield.Iron
	p.WoodYield = yield.Wood
	p.DepotCapacity = 10000
	return p
}

var RoleResDao = &roleResDao{
	rrChan: make(chan *RoleRes, 100),
}

type roleResDao struct {
	rrChan chan *RoleRes
}

func (r *roleResDao) running() {
	for {
		select {
		case rr := <-r.rrChan:
			_, err := db.Engine.
				Table(new(RoleRes)).
				ID(rr.Id).
				Cols("wood", "iron", "stone", "grain", "gold").
				Update(rr)
			if err != nil {
				log.Println("RoleResDao update error", err)
			}
		}
	}
}

func init() {
	go RoleResDao.running()
}

type Yield struct {
	Wood  int
	Iron  int
	Stone int
	Grain int
	Gold  int
}

// SyncExecute 同步更新数据库
func (r *RoleRes) SyncExecute() {
	RoleResDao.rrChan <- r
}
