package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var WarService = &warService{}

type warService struct {
}

func (w *warService) GetWarReports(rid int) ([]model.WarReport, error) {
	mrs := make([]data.WarReport, 0)
	mr := &data.WarReport{}
	err := db.Engine.Table(mr).
		Where("a_rid=? or d_rid=?", rid, rid).
		Limit(30, 0).
		Desc("ctime").
		Find(&mrs)
	if err != nil {
		log.Println("战报查询出错", err)
		return nil, common.New(constant.DBError, "战报查询出错")
	}
	modelMrs := make([]model.WarReport, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.WarReport))
	}
	return modelMrs, nil
}
