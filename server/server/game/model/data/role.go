package data

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"time"
)

type Role struct {
	RId        int       `xorm:"rid pk autoincr"`
	UId        int       `xorm:"uid"`
	NickName   string    `xorm:"nick_name" validate:"min=4,max=20,regexp=^[a-zA-Z0-9_]*$"`
	Balance    int       `xorm:"balance"`
	HeadId     int16     `xorm:"headId"`
	Sex        int8      `xorm:"sex"`
	Profile    string    `xorm:"profile"`
	LoginTime  time.Time `xorm:"login_time"`
	LogoutTime time.Time `xorm:"logout_time"`
	CreatedAt  time.Time `xorm:"created_at"`
}

func (r *Role) TableName() string {
	return "role"
}

func (r *Role) ToModel() interface{} {
	m := model.Role{}
	m.UId = r.UId
	m.RId = r.RId
	m.Sex = r.Sex
	m.NickName = r.NickName
	m.HeadId = r.HeadId
	m.Balance = r.Balance
	m.Profile = r.Profile
	return m
}
