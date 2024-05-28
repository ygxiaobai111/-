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
	utils "github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/util"
	"log"
)

var DefaultRoleController = &RoleController{}

type RoleController struct {
}

func (r *RoleController) Router(router *net.Router) {
	g := router.Group("role")
	g.Use(middleware.Log())
	g.AddRouter("enterServer", r.enterServer)
	g.AddRouter("myProperty", r.myProperty, middleware.CheckRole())
	g.AddRouter("posTagList", r.posTagList, middleware.CheckRole())
}

func (r *RoleController) enterServer(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//进入的游戏的逻辑
	//Session 需要验证是否合法 合法的情况下 可以取出登录的用户id
	//根据用户id 去查询对应的游戏角色，如果有 就继续 没有 提示无角色
	//根据角色id 查询角色拥有的资源 roleRes，如果资源有 返回，没有 初始化资源
	reqObj := &model.EnterServerReq{}
	rspObj := &model.EnterServerRsp{}
	err := mapstructure.Decode(req.Body.Msg, reqObj)
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	if err != nil {
		log.Println("enterServer err:", err)
		rsp.Body.Code = constant.InvalidParam
		return
	}
	session := reqObj.Session
	_, claim, err := utils.ParseToken(session)
	if err != nil {
		log.Println("enterServer err:", err)

		rsp.Body.Code = constant.SessionInvalid
		return
	}
	uid := claim.Uid
	err = logic.RoleService.EnterServer(uid, rspObj, req)
	if err != nil {
		log.Println("enterServer err:", err)
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
	log.Println("enterServer resp:", rsp)
}

func (r *RoleController) myProperty(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//分别根据角色id 去查询 军队 资源 建筑 城池 武将
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rid := role.(*data.Role).RId
	rspObj := &model.MyRolePropertyRsp{}
	//资源
	rspObj.RoleRes, err = logic.RoleService.GetRoleRes(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//城池
	rspObj.Citys, err = logic.RoleCityService.GetRoleCitys(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//建筑
	rspObj.MRBuilds, err = logic.RoleBuildService.GetBuilds(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//军队
	rspObj.Armys, err = logic.ArmyService.GetArmys(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	//武将
	rspObj.Generals, err = logic.GeneralService.GetGenerals(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj
}
func (r *RoleController) posTagList(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	rspObj := &model.PosTagListRsp{}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	//去 角色属性 表去查询
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId
	pts, err := logic.RoleAttrService.GetTagList(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.PosTags = pts
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

}
