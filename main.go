/*
	平台主函数
*/
package main

import (
	"fmt"
	_ "incentiveblog/config"
	"incentiveblog/router"
	_ "net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	fmt.Println("welcome to incentiveblog")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	//加一个测试路由
	e.GET("/ping", router.Pong)

	e.Logger.Fatal(e.Start(":8086"))
}
