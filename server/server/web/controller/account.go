package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web/logic"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web/model"

	"net/http"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {
}

func (a *AccountController) Register(ctx *gin.Context) {
	/**
	1. 获取请求参数
	2. 根据用户名 查询数据库是否有 有 用户名已存在 没有 注册
	3. 告诉前端 注册成功即可
	*/
	rq := &model.RegisterReq{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}
	//一般web服务 错误格式 自定义
	err = logic.DefaultAccountLogic.Register(rq)
	if err != nil {
		log.Println("注册业务出错", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil))
}
