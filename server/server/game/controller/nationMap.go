package controller

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
)

var DefaultNationMapController = &nationMapController{}

type nationMapController struct {
}

func (r *nationMapController) Router(router *net.Router) {
	g := router.Group("nationMap")
	g.AddRouter("config", r.config)
}

func (r *nationMapController) config(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//reqObj := &model.ConfigReq{}
	rspObj := &model.ConfigRsp{}

	cfgs := gameConfig.MapBuildConf.Cfg

	rspObj.Confs = make([]model.Conf, len(cfgs))
	for index, v := range cfgs {
		rspObj.Confs[index].Type = v.Type
		rspObj.Confs[index].Name = v.Name
		rspObj.Confs[index].Level = v.Level
		rspObj.Confs[index].Defender = v.Defender
		rspObj.Confs[index].Durable = v.Durable
		rspObj.Confs[index].Grain = v.Grain
		rspObj.Confs[index].Iron = v.Iron
		rspObj.Confs[index].Stone = v.Stone
		rspObj.Confs[index].Wood = v.Wood
	}
	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	rsp.Body.Code = constant.OK
	rsp.Body.Msg = rspObj

}
