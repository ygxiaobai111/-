package controller

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/logic"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/middleware"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
)

var DefaultGeneralController = &GeneralController{}

type GeneralController struct {
}

func (r *GeneralController) Router(router *net.Router) {
	g := router.Group("general")
	g.AddRouter("myGenerals", r.myGenerals, middleware.CheckRole())
}

func (r *GeneralController) myGenerals(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查询武将的时候 角色拥有的武将 查询出来即可
	// 如果初始化 进入游戏 武将没有 需要随机三个武将 很多游戏 初始化武将是一样的

	rspObj := &model.MyGeneralRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId

	gs, err := logic.GeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Generals = gs
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

}
