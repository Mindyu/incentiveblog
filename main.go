/*
	平台主函数
*/
package main

import (
	"fmt"
	"incentiveblog/router"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
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

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	//加一个测试路由
	e.GET("/ping", pong)
	e.POST("/user", router.Register)
	e.GET("/session", router.GetSession)
	e.POST("/login", router.Login)
	e.POST("/content", router.UploadContent)
	e.POST("/usericon", router.UploadIcon)
	e.GET("/content/:hash", router.GetContent)
	e.GET("/blog/:hash", router.GetContent)
	e.GET("/concern", router.Concern)             // 关注对象
	e.POST("/blog", router.PublishBlog)           // 发表博客
	e.GET("/blogs/:userid", router.ListsLimit)    // 查看博客-指定userid
	e.GET("/blogs", router.ListsAll)              // 查看自己全部博客
	e.GET("/concern/:userid", router.UserConcern) // 查看用户关注人数和粉丝数量
	e.GET("/blogs", router.ListsAll)              // 查看自己全部博客
	e.GET("/fans", router.GetFans)                // 获取粉丝列表
	e.GET("/concerns", router.GetConcerns)        // 获取关注列表
	e.GET("/details", router.GetDetails)          // 获取积分明细

	e.Logger.Fatal(e.Start(":8086"))
}
