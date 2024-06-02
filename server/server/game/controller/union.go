package controller

import (
	"github.com/mitchellh/mapstructure"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/logic"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/middleware"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
)

var UnionController = &unionController{}

type unionController struct {
}

func (u *unionController) Router(router *net.Router) {
	g := router.Group("union")
	g.Use(middleware.Log())
	g.AddRouter("list", u.list, middleware.CheckRole())
	g.AddRouter("info", u.info, middleware.CheckRole())
	g.AddRouter("applyList", u.applyList, middleware.CheckRole())
}

func (u *unionController) list(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查询数据库 把所有的联盟 查询出来
	rspObj := &model.ListRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
	uns, err := logic.CoalitionService.List()
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = uns
	rsp.Body.Msg = rspObj

}

func (u *unionController) info(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	reqObj := &model.InfoReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.InfoRsp{}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

	un := logic.CoalitionService.Get(reqObj.Id)
	rspObj.Info = un
	rspObj.Id = un.Id

}

func (u *unionController) applyList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//根据联盟id 去查询申请列表，rid申请人，你角色表 查询详情即可
	// state 0 正在申请 1 拒绝 2 同意
	//什么人能看到申请列表 只有盟主和副盟主能看到申请列表
	reqObj := &model.ApplyReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.ApplyRsp{}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)
	//查询联盟
	un := logic.CoalitionService.GetCoalition(reqObj.Id)
	if un == nil {
		rsp.Body.Code = constant.DBError
		return
	}
	if un.Chairman != role.RId && un.ViceChairman != role.RId {
		rspObj.Id = reqObj.Id
		rspObj.Applys = make([]model.ApplyItem, 0)
		return
	}

	ais, err := logic.CoalitionService.GetListApply(reqObj.Id, 0)
	if err != nil {
		rsp.Body.Code = constant.DBError
		return
	}
	rspObj.Id = reqObj.Id
	rspObj.Applys = ais
}
