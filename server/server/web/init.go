package web

import (
	"github.com/gin-gonic/gin"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/models"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web/controller"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web/middleware"
)

func Init(router *gin.Engine) {
	//测试数据库，并且初始化数据库
	db.TestDB().Sync(new(models.User))
	//还有别的初始化方法
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	router.Use(middleware.Cors())
	router.Any("/account/register", controller.DefaultAccountController.Register)
}
