package router

import (
	"crypto/sha256"
	"errors"
	_ "errors"
	"fmt"
	"incentiveblog/db"
	"incentiveblog/utils"

	_ "github.com/gorilla/sessions"
	"github.com/labstack/echo"
	_ "github.com/labstack/echo-contrib/session"
	"gopkg.in/mgo.v2/bson"
)

type Blog struct {
	UserID  string `json:"userid"`
	Title   string `json:"title"`
	BlogUrl string `json:"blogurl"`
}
type PointDetails struct {
	UserID string `json:"userid"`
	Remark string `json:"remark"`
	Points int    `json:"points"`
}

//获取图片
func GetContent(c echo.Context) error {
	contentHash := c.Param("hash")
	return db.Dbconn.DownLoadFile("blogplatform", "blogs", contentHash, c)
}

//发布博客
func PublishBlog(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	title := c.FormValue("title")
	content := c.FormValue("content")
	//cData := []byte(content)
	blogUrl := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	err := db.Dbconn.UploadFile("blogplatform", "blogs", blogUrl, []byte(content))
	if err != nil {
		fmt.Println("failed to uploadfile ", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	user := &db.User{}
	if err := c.Bind(user); err != nil {
		fmt.Println(user)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	err = db.Dbconn.QueryOne("blogplatform", "users", bson.M{"token": user.Token}, &user)

	if err != nil {
		fmt.Println("failed to login", err)
		resp.Errno = utils.RECODE_LOGINERR
		return err
	}

	blog := Blog{user.UserID, title, blogUrl}
	err = db.Dbconn.Insert("blogplatform", "blogs", blog)
	if err != nil {
		fmt.Println("failed to insert into blogs", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	//增加积分
	detail := PointDetails{user.UserID, "publish blog", 10}
	err = db.Dbconn.Insert("blogplatform", "details", detail)
	if err != nil {
		fmt.Println("failed to insert into details", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}

	return nil
}

//上传图片
func UploadContent(c echo.Context) error {
	//1. 响应数据结构初始化
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)
	//2. 保存文件
	url, err := StorageDBFile(c, &resp)
	if err != nil {
		return err
	}

	resp.Data = "/content/" + url

	return nil
}

//上传文件：可以是md或photo
func StorageDBFile(c echo.Context, resp *utils.Resp) (url string, err error) {
	//1. 解析数据
	h, err := c.FormFile("content")
	if err != nil {
		fmt.Println("failed to FormFile ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return "", err
	}
	src, err := h.Open()
	defer src.Close()

	//2. 计算hash
	cData := make([]byte, h.Size)
	n, err := src.Read(cData)
	if err != nil || h.Size != int64(n) {
		resp.Errno = utils.RECODE_UNKNOWERR
		return "", err
	}
	contentHash := fmt.Sprintf("%x", sha256.Sum256(cData))

	//3. 修改数据库
	err = db.Dbconn.UploadFile("blogplatform", "blogs", contentHash, cData)
	if err != nil {
		fmt.Println("failed to uploadfile ", err)
		resp.Errno = utils.RECODE_DBERR
		return "", err
	}

	return contentHash, nil
}

//获取博客列表-其他用户
func ListsLimit(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//从上下文获得被关注对象
	userID := c.Param("userid")
	if userID == "" {
		fmt.Println("failed to get UserID")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("params err")
	}
	var blogs []Blog
	err := db.Dbconn.QueryAll("blogplatform", "blogs", bson.M{"userid": userID}, &blogs, 0, 5)
	if err != nil {
		fmt.Println("failed to query all", err)
		return err
	}
	fmt.Println(blogs)
	resp.Data = blogs
	return nil
}

//获取博客列表-个人
func ListsAll(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//从上下文获得被关注对象
	//	sess, _ := session.Get("session", c)
	//	userID, ok := sess.Values["userid"].(string)
	//	if userID == "" || !ok {
	//		resp.Errno = utils.RECODE_SESSIONERR
	//		return errors.New("no session")
	//	}
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

	var blogs []Blog
	err = db.Dbconn.QueryAll("blogplatform", "blogs", bson.M{"userid": user.UserID}, &blogs, 0, 0)
	if err != nil {
		fmt.Println("failed to query all", err)
		return err
	}
	fmt.Println(blogs)
	resp.Data = blogs
	return nil
}

//获取用户积分明细
func GetDetails(c echo.Context) error {
	var resp utils.Resp
	resp.Errno = utils.RECODE_OK
	defer utils.ResponseData(c, &resp)

	//从上下文获得被关注对象
	//	sess, _ := session.Get("session", c)
	//	userID, ok := sess.Values["userid"].(string)
	//	if userID == "" || !ok {
	//		resp.Errno = utils.RECODE_SESSIONERR
	//		return errors.New("no session")
	//	}

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

	var details []PointDetails
	err = db.Dbconn.QueryAll("blogplatform", "details", bson.M{"userid": user.UserID}, &details, 0, 0)
	if err != nil {
		fmt.Println("failed to query all", err)
		return err
	}
	//fmt.Println(details)
	resp.Data = details
	return nil
}
