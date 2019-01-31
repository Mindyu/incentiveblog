package router

import (
	"errors"
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
	u.Token = utils.MakeToken([]byte(u.UserID), []byte(u.PassWord))
	//注册送积分
	u.Points = 100
	//3. 操作数据库
	err = db.Insert(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, u)
	if err != nil {
		c.Logger().Error("failed to insert into user", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	// -- 维护积分明细表
	detail := db.PointsDetail{u.UserID, "注册赠送", 100}
	err = db.Insert(config.ServerConfig.DB.DBName, config.ServerConfig.DB.DetailTab, &detail)
	if err != nil {
		c.Logger().Error("failed to insert into detail", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = u.Token
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
	u.Token = utils.MakeToken([]byte(u.UserID), []byte(u.PassWord))
	//4. 新token保存到数据库
	err = db.Update(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, bson.M{"userid": u.UserID, "password": u.PassWord}, u)
	if err != nil {
		c.Logger().Error("failed to update  user", err, u.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = u
	return nil
}

//检测登陆
func CheckLogin(c echo.Context) error {
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
	//3. 操作数据库 查询，匹配token
	err = db.QueryOne(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, bson.M{"token": u.Token}, u)
	if err != nil {
		c.Logger().Error("failed to get  user", err, u.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = u
	return nil
}

//积分明细查询
func GetDetail(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	u := new(db.QueryDetail)

	err := c.Bind(u)
	if err != nil || u.UserID == "" {
		c.Logger().Error("failed to get param user", err, u)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//3. 操作数据库 查询，匹配userID,返回一个数组
	var details []db.PointsDetail
	err = db.QueryAll(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.DetailTab,
		bson.M{"userid": u.UserID},
		&details,
		u.Skip,
		u.Limit)
	if err != nil {
		c.Logger().Error("failed to get  details", err, u.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = details
	return nil
}

func UserConcern(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	u := new(db.User)

	err := c.Bind(u)
	if err != nil || u.UserID == "" {
		c.Logger().Error("failed to get param user", err, u)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	//3. 验证token，同时得到本人userid
	toUserID := u.UserID
	err = db.QueryOne(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.UserTab,
		bson.M{"token": u.Token},
		u)
	if err != nil {
		c.Logger().Error("failed to get  user", err, u.Token)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//4. 插入到数据库 user_relations
	err = db.Insert(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.RelationTab,
		bson.M{"fromid": u.UserID, "toid": toUserID},
	)
	if err != nil {
		c.Logger().Error("failed to insert  user_relation", err, u.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	return nil
}

//统计粉丝数量
func CountConcern(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	userID := c.Param("userid")
	if userID == "" {
		c.Logger().Error("failed to get param user")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("not found user")
	}
	//3. 查数据，2次
	fromCount, err := db.DbConn.DB(config.ServerConfig.DB.DBName).C(config.ServerConfig.DB.RelationTab).Find(bson.M{"fromid": userID}).Count()
	if err != nil {
		c.Logger().Error("failed to query  user_relation", err, userID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	toCount, err := db.DbConn.DB(config.ServerConfig.DB.DBName).C(config.ServerConfig.DB.RelationTab).Find(bson.M{"toid": userID}).Count()
	if err != nil {
		c.Logger().Error("failed to query  user_relation", err, userID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//4. 响应消息处理
	cntMap := make(map[string]int)
	cntMap["fromcount"] = fromCount
	cntMap["tocount"] = toCount
	resp.Data = cntMap
	return nil
}
