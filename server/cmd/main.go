package main

import (
	"fmt"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/config"
)

func main() {
	fmt.Println("无线三国")
	mysqlConfig, err := config.File.GetSection("mysql")
	if err != nil {
		panic(err)
	}
	fmt.Println(mysqlConfig)
}
