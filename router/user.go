package router

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"incentiveblog/db"
	"incentiveblog/utils"
	"strconv"

	_ "github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"gopkg.in/mgo.v2/bson"
)

//注册功能
func Register(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//2. 解析数据
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	user.Token = utils.MakeToken([]byte(user.UserID), []byte(user.UserName))
	user.Points = 100
	//3. 操作mongodb
	err := db.Dbconn.Insert("blogplatform", "users", user)
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	//4. session处理
	//	sess, _ := session.Get("session", c)
	//	sess.Options = &sessions.Options{
	//		Path:     "/",
	//		MaxAge:   86400 * 7,
	//		HttpOnly: true,
	//	}
	//	sess.Values["userid"] = user.UserID
	//	sess.Values["username"] = user.UserName
	//	sess.Save(c.Request(), c.Response())
	//5. 赠送积分
	detail := PointDetails{user.UserID, "register", 100}
	err = db.Dbconn.Insert("blogplatform", "details", detail)
	if err != nil {
		fmt.Println("failed to insert into details", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = user.Token
	return nil
}

//session获取
func GetSession(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//处理session
	sess, err := session.Get("session", c)
	if err != nil {
		fmt.Println("failed to get session")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	userid := sess.Values["userid"]
	if userid == nil {
		fmt.Println("failed to get session,userid is nil")
		resp.Errno = utils.RECODE_SESSIONERR
		return err
	}
	user := &db.User{}
	err = db.Dbconn.QueryOne("blogplatform", "users", bson.M{"userid": userid}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	resp.Data = user
	return nil
}

func CheckLogin(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//处理session
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	err := db.Dbconn.QueryOne("blogplatform", "users", bson.M{"token": user.Token}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	resp.Data = user
	return nil
}

//登陆功能
func Login(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//2. 解析数据
	user := &db.User{}
	fmt.Println(c)
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	fmt.Println("-----------------------", user)
	//3. 查看mongodb
	//db.users.find({"userid":"yekai","password":"123"})
	err := db.Dbconn.QueryOne("blogplatform", "users", bson.M{"userid": user.UserID, "password": user.PassWord}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	user.Token = utils.MakeToken([]byte(user.UserID), []byte(user.UserName))

	//更新 token Update(DBName, collection string, cond, data interface{})
	err = db.Dbconn.Update("blogplatform", "users", bson.M{"userid": user.UserID, "password": user.PassWord}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	resp.Data = user
	//4. session处理
	//	sess, _ := session.Get("session", c)
	//	sess.Options = &sessions.Options{
	//		Path:     "/",
	//		MaxAge:   86400 * 7,
	//		HttpOnly: true,
	//	}
	//	sess.Values["username"] = user.UserName
	//	sess.Values["userid"] = user.UserID
	//	sess.Save(c.Request(), c.Response())
	return nil
}

//上传头像
func UploadIcon(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//2. 解析数据

	h, err := c.FormFile("icon")
	if err != nil {
		fmt.Println("failed to FormFile ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	src, err := h.Open()
	defer src.Close()

	//计算hash
	cData := make([]byte, h.Size)
	n, err := src.Read(cData)
	if err != nil || h.Size != int64(n) {
		resp.Errno = utils.RECODE_UNKNOWERR
		return err
	}
	contentHash := fmt.Sprintf("%x", sha256.Sum256(cData))

	photoUrl := contentHash
	//3. 修改数据库
	err = db.Dbconn.UploadFile("blogplatform", "blogs", photoUrl, cData)
	if err != nil {
		fmt.Println("failed to uploadfile ", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	//保存映射关系
	//	sess, _ := session.Get("session", c)
	//	userID, ok := sess.Values["userid"].(string)
	//	if userID == "" || !ok {
	//		resp.Errno = utils.RECODE_SESSIONERR
	//		return errors.New("no session")
	//	}
	//Update(DBName, collection string, cond, data interface{})
	user := &db.User{}
	fmt.Println(c)
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	err = db.Dbconn.Update("blogplatform", "users", bson.M{"token": user.Token}, bson.M{"$set": bson.M{"iconurl": photoUrl}})
	if err != nil {
		resp.Errno = utils.RECODE_DBERR
		fmt.Println("failed to update user", err, user.Token)
		return err
	}
	return nil
}

//用户关注 - GET
func Concern(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//从session获取userID
	//	sess, err := session.Get("session", c)
	//	if err != nil {
	//		fmt.Println("failed to get session", err)
	//		resp.Errno = utils.RECODE_SESSIONERR
	//		return err
	//	}
	//	userID := sess.Values["userid"]
	//	if userID == nil {
	//		fmt.Println("failed to get session,userID is nil")
	//		resp.Errno = utils.RECODE_SESSIONERR
	//		return err
	//	}
	user := &db.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	toUserID := user.UserID
	err := db.Dbconn.QueryOne("blogplatform", "users", bson.M{"token": user.Token}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}
	//从上下文获得被关注对象

	err = db.Dbconn.Insert("blogplatform", "user_relations", bson.M{"fromid": user.UserID, "toid": toUserID})
	if err != nil {
		fmt.Println("failed to insert into  user_relations", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	return nil
}

//获取用户关注人数和粉丝数量
func UserConcern(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//从session获取userID

	//从上下文获得被关注对象
	userID := c.Param("userid")
	if userID == "" {
		fmt.Println("failed to get toUserID")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("params err")
	}
	fromCount, err := db.Dbconn.Count("blogplatform", "user_relations", bson.M{"toid": userID})
	if err != nil {
		fmt.Println("failed to count", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	toCount, err := db.Dbconn.Count("blogplatform", "user_relations", bson.M{"fromid": userID})
	if err != nil {
		fmt.Println("failed to count", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	cntMap := make(map[string]int)
	cntMap["fromcount"] = fromCount
	cntMap["tocount"] = toCount
	fmt.Println(cntMap)
	resp.Data = cntMap
	return nil
}

//获得粉丝列表
func GetFans(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//从session获取userID

	//从上下文获得被关注对象
	userID := c.QueryParam("userid")
	strskip := c.QueryParam("skip")
	strlimit := c.QueryParam("limit")
	if userID == "" {
		fmt.Println("failed to get toUserID")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("params err")
	}

	skip, _ := strconv.Atoi(strskip)
	limit, _ := strconv.Atoi(strlimit)

	var userrel []db.UserRelation
	err := db.Dbconn.QueryAll("blogplatform", "user_relations", bson.M{"toid": userID}, &userrel, skip, limit)
	if err != nil {
		fmt.Println("failed to GetFans", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	resp.Data = userrel
	return nil
}

//获取关注列表
func GetConcerns(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//从session获取userID

	//从上下文获得被关注对象
	userID := c.QueryParam("userid")
	strskip := c.QueryParam("skip")
	strlimit := c.QueryParam("limit")
	if userID == "" {
		fmt.Println("failed to get toUserID")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("params err")
	}

	skip, _ := strconv.Atoi(strskip)
	limit, _ := strconv.Atoi(strlimit)

	var userrel []db.UserRelation
	err := db.Dbconn.QueryAll("blogplatform", "user_relations", bson.M{"fromid": userID}, &userrel, skip, limit)
	if err != nil {
		fmt.Println("failed to GetConcerns", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	resp.Data = userrel
	return nil
}
