package logic

import (
	"log"

	"time"

	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/models"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/util"
)

var DefaultAccountLogic = &AccountLogic{}

type AccountLogic struct {
}

func (l AccountLogic) Register(rq *model.RegisterReq) error {
	username := rq.Username
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", username).Get(user)
	if err != nil {
		log.Println("注册查询失败", err)
		return common.New(constant.DBError, "数据库异常")
	}
	if ok {
		//有数据 提示用户已存在
		return common.New(constant.UserExist, "用户已存在")
	} else {
		user.Mtime = time.Now()
		user.Ctime = time.Now()
		user.Username = rq.Username
		user.Passcode = util.RandSeq(6)
		user.Passwd = util.Password(rq.Password, user.Passcode)
		user.Hardware = rq.Hardware
		_, err := db.Engine.Table(user).Insert(user)
		if err != nil {
			log.Println("注册插入失败", err)
			return common.New(constant.DBError, "数据库异常")
		}
		return nil
	}
}
