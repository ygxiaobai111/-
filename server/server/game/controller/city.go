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

var CityController = &cityController{}

type cityController struct {
}

func (c *cityController) Router(router *net.Router) {
	g := router.Group("city")
	g.Use(middleware.Log())
	g.AddRouter("facilities", c.facilities, middleware.CheckRole())
	g.AddRouter("upFacility", c.upFacility, middleware.CheckRole())
}

func (c *cityController) facilities(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//参数 城池id 需要查询城池信息
	//还需要查询城池里面的设施列表
	reqObj := &model.FacilitiesReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.FacilitiesRsp{}

	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

	//角色
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	//查询城池
	rc, ok := logic.RoleCityService.Get(reqObj.CityId)
	if !ok {
		rsp.Body.Code = constant.CityNotExist
		return
	}
	if role.RId != rc.RId {
		rsp.Body.Code = constant.CityNotMe
		return
	}
	//查询城池的设施
	fac := logic.CityFacilityService.GetFacility(role.RId, reqObj.CityId)
	if fac == nil {
		rsp.Body.Code = constant.CityNotExist
		return
	}
	rspObj.CityId = reqObj.CityId
	rspObj.Facilities = make([]model.Facility, len(fac))
	for index, v := range fac {
		rspObj.Facilities[index].Type = v.Type
		rspObj.Facilities[index].Name = v.Name
		rspObj.Facilities[index].Level = v.GetLevel()
		rspObj.Facilities[index].UpTime = v.UpTime
	}

}

func (c *cityController) upFacility(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//1. 需要根据城池id 查询城池 确保城池存在
	//2. 查询城池的设施 确保设施存在
	//3. 升级设施，需要更新升级时间 upTime 升级完 upTime=0
	//4. 升级时候 判断资源是否符合条件 如果符合才能升级，升级完成 进行数据库的更新 设施更新的内容固话到数据库
	//5. 消耗资源，资源减少 固化到数据库
	//6. 资源查询出来 返回前端
	reqObj := &model.UpFacilityReq{}
	mapstructure.Decode(req.Body.Msg, reqObj)
	rspObj := &model.UpFacilityRsp{}
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK

	//角色
	r, _ := req.Conn.GetProperty("role")
	role := r.(*data.Role)

	rc, ok := logic.RoleCityService.Get(reqObj.CityId)
	if !ok {
		rsp.Body.Code = constant.CityNotExist
		return
	}
	if rc.RId != role.RId {
		rsp.Body.Code = constant.CityNotMe
		return
	}
	facs := logic.CityFacilityService.GetFacility(role.RId, reqObj.CityId)
	if facs == nil {
		rsp.Body.Code = constant.CityNotExist
		return
	}
	fac, err := logic.CityFacilityService.UpFacility(role.RId, reqObj.CityId, reqObj.FType)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.Facility.Name = fac.Name
	rspObj.Facility.Level = fac.GetLevel()
	rspObj.Facility.UpTime = fac.UpTime
	rspObj.Facility.Type = fac.Type

	res := logic.RoleResService.GetRoleRes(role.RId)
	if res != nil {
		rspObj.RoleRes = res.ToModel().(model.RoleRes)
	}
	rspObj.CityId = reqObj.CityId
}
