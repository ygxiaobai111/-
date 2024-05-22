package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/config"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/web"
)

func main() {
	host := config.File.MustValue("web_server", "host", "127.0.0.1")
	port := config.File.MustValue("web_server", "port", "8088")

	router := gin.Default()
	//路由
	web.Init(router)
	s := &http.Server{
		Addr:           host + ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	log.Println(err)
}
