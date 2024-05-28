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

var SkillController = &skillController{}

type skillController struct {
}

func (s *skillController) Router(router *net.Router) {
	g := router.Group("skill")
	g.AddRouter("list", s.list, middleware.CheckRole())
}

func (s *skillController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查找战报表 得出数据
	rspObj := &model.SkillListRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId

	skills, err := logic.SkillService.GetSkills(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = skills
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
}
