package data

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"time"
)

type RoleAttribute struct {
	Id              int            `xorm:"id pk autoincr"`
	RId             int            `xorm:"rid"`
	UnionId         int            `xorm:"-"`                 //联盟id
	ParentId        int            `xorm:"parent_id"`         //上级id（被沦陷）
	CollectTimes    int8           `xorm:"collect_times"`     //征收次数
	LastCollectTime time.Time      `xorm:"last_collect_time"` //最后征收的时间
	PosTags         string         `xorm:"pos_tags"`          //位置标记
	PosTagArray     []model.PosTag `xorm:"-"`
}

func (r *RoleAttribute) TableName() string {
	return "role_attribute"
}
