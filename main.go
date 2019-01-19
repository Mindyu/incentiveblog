/*
	平台主函数
*/
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func pong(c echo.Context) error {
	return c.String(http.StatusOK, "pong\n")
}

func main() {
	fmt.Println("welcome to incentiveblog")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())

	//加一个测试路由
	e.GET("/ping", pong)

	e.Logger.Fatal(e.Start(":8086"))
}
