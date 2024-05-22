package controller

import (
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/login/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/login/proto"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/models"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/util"

	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	/**
		1. 用户名 密码 硬件id
		2. 根据用户名 查询user表 得到数据
		3. 进行密码比对，如果密码正确 登录成功
	    4. 保存用户登录记录
		5. 保存用户的最后一次登录信息
	    6. 客户端 需要 一个session，jwt 生成一个加密字符串的加密算法
	    7. 客户端 在发起需要用户登录的行为时，判断用户是否合法
	*/
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil {
		log.Println("用户表查询出错", err)
		return
	}
	if !ok {
		//有没有查出来数据
		rsp.Body.Code = constant.UserNotExist
		return
	}
	pwd := util.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		rsp.Body.Code = constant.PwdIncorrect
		return
	}
	//jwt A.B.C 三部分 A定义加密算法 B定义放入的数据 C部分 根据秘钥+A和B生成加密字符串
	token, _ := util.Award(user.UId)
	rsp.Body.Code = constant.OK
	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

	//保存用户登录记录
	ul := &model.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, State: model.Login,
	}
	db.Engine.Table(ul).Insert(ul)
	//最后一次登录的状态记录
	ll := &model.LoginLast{}
	ok, _ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		//有数据 更新
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engine.Table(ll).Update(ll)
	} else {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		_, err := db.Engine.Table(ll).Insert(ll)
		if err != nil {
			log.Println(err)
		}
	}
	//缓存一下 此用户和当前的ws连接
	net.Mgr.UserLogin(req.Conn, user.UId, token)
}
