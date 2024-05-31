package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/global"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	utils "github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/util"
	"log"
	"sync"
)

var RoleBuildService = &roleBuildService{
	posRB:  make(map[int]*data.MapRoleBuild),
	roleRB: make(map[int][]*data.MapRoleBuild),
}

func (r *roleBuildService) Load() {
	//加载系统的建筑以及玩家的建筑
	//首先需要判断数据库 是否保存了系统的建筑 没有 进行一个保存
	total, err := db.Engine.
		Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).
		Count(new(data.MapRoleBuild))
	if err != nil {
		panic(err)
	}
	if int64(len(gameConfig.MapRes.SysBuild)) != total {
		//证明数据库存储的系统建筑 有问题
		db.Engine.
			Where("type=? or type=?", gameConfig.MapBuildSysCity, gameConfig.MapBuildSysFortress).
			Delete(new(data.MapRoleBuild))
		for _, v := range gameConfig.MapRes.SysBuild {
			build := &data.MapRoleBuild{
				RId:   0,
				Type:  v.Type,
				Level: v.Level,
				X:     v.X,
				Y:     v.Y,
			}
			build.Init()
			db.Engine.InsertOne(build)
		}
	}
	//查询所有的角色建筑
	dbRB := make(map[int]*data.MapRoleBuild)
	db.Engine.Find(dbRB)

	for _, v := range dbRB {
		posId := global.ToPosition(v.X, v.Y)
		r.posRB[posId] = v
		_, ok := r.roleRB[v.RId]
		if !ok {
			r.roleRB[v.RId] = make([]*data.MapRoleBuild, 0)
		} else {
			r.roleRB[v.RId] = append(r.roleRB[v.RId], v)
		}
	}
}

type roleBuildService struct {
	mutex sync.RWMutex
	//位置 key posId
	posRB map[int]*data.MapRoleBuild
	//key 角色id
	roleRB map[int][]*data.MapRoleBuild
}

func (r *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	mrs := make([]data.MapRoleBuild, 0)
	mr := &data.MapRoleBuild{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("建筑查询出错", err)
		return nil, common.New(constant.DBError, "建筑查询出错")
	}
	modelMrs := make([]model.MapRoleBuild, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.MapRoleBuild))
	}
	return modelMrs, nil
}

func (r *roleBuildService) ScanBlock(req *model.ScanBlockReq) ([]model.MapRoleBuild, interface{}) {
	x := req.X
	y := req.Y
	length := req.Length
	var mrbs = make([]model.MapRoleBuild, 0)
	if x < 0 || x >= global.MapWith || y < 0 || y >= global.MapHeight {
		return mrbs, nil
	}
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	maxX := utils.MinInt(global.MapWith, x+length-1)
	maxY := utils.MinInt(global.MapHeight, y+length-1)

	//范围  x-length  x + length  y-length y+length
	for i := x - length; i <= maxX; i++ {
		for j := y - length; j <= maxY; j++ {
			posId := global.ToPosition(i, j)
			mrb, ok := r.posRB[posId]
			if ok {
				mrbs = append(mrbs, mrb.ToModel().(model.MapRoleBuild))
			}
		}
	}
	log.Println("玩家建筑", mrbs)
	return mrbs, nil
}
func (r *roleBuildService) GetYield(rid int) data.Yield {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	rbs, ok := r.roleRB[rid]
	var yield data.Yield
	if ok {
		for _, v := range rbs {
			yield.Stone += v.Stone
			yield.Wood += v.Wood
			yield.Iron += v.Iron
			yield.Grain += v.Grain
		}
	}
	return yield
}
