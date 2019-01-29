package router

import (
	"incentiveblog/config"
	"incentiveblog/db"
	"incentiveblog/utils"
	"net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

//统一响应消息
func ResponseData(c echo.Context, resp *utils.Resp) {
	resp.ErrMsg = utils.GetRecodeText(resp.Errno)
	c.JSON(http.StatusOK, resp)
}

func Pong(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	return nil
}

func Register(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	u := new(db.User)

	err := c.Bind(u)
	if err != nil {
		c.Logger().Error("failed to get param user", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//3. 操作数据库
	err = db.Insert(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, u)
	if err != nil {
		c.Logger().Error("failed to insert into user", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	return nil
}

//登陆
func Login(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	u := new(db.User)

	err := c.Bind(u)
	if err != nil {
		c.Logger().Error("failed to get param user", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//3. 操作数据库 查询，匹配用户名和密码ok
	err = db.QueryOne(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, bson.M{"userid": u.UserID, "password": u.PassWord}, u)
	if err != nil {
		c.Logger().Error("failed to get  user", err, u.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = u
	return nil
}