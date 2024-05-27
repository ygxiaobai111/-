package controller

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/logic"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
)

var WarController = &warController{}

type warController struct {
}

func (w *warController) Router(router *net.Router) {
	g := router.Group("war")
	g.AddRouter("report", w.report)
}

func (w *warController) report(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//查找战报表 得出数据
	rspObj := &model.WarReportRsp{}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	role, err := req.Conn.GetProperty("role")
	if err != nil {
		rsp.Body.Code = constant.SessionInvalid
		return
	}
	rid := role.(*data.Role).RId

	reports, err := logic.WarService.GetWarReports(rid)
	if err != nil {
		rsp.Body.Code = err.(*common.MyError).Code()
		return
	}
	rspObj.List = reports
	rsp.Body.Msg = rspObj
	rsp.Body.Code = constant.OK
}
