package router

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"incentiveblog/config"
	"incentiveblog/db"
	"incentiveblog/utils"
	_ "net/http"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type RespOneBlog struct {
	Text  string `json:"text"`
	Url   string `json:"url"`
	IsAll bool   `json:"isall"`
}

type RespBlogs struct {
	UserID   string        `json:"userid"`
	UserName string        `json:"username"`
	Blogs    []RespOneBlog `json:"blogs"`
}

//文件上传，返回一个url
func UploadContent(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 解析请求消息
	file, err := c.FormFile("content") // content=@xxx.jpg
	if err != nil {
		c.Logger().Error("failed to formfile ", err)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	src, err := file.Open()
	if err != nil {
		c.Logger().Error("failed to openfile ", err, file.Filename)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}
	defer src.Close()
	cData := make([]byte, file.Size)
	n, err := src.Read(cData)
	if err != nil || int64(n) != file.Size {
		c.Logger().Error("failed to read file ", err, file.Filename)
		resp.Errno = utils.RECODE_PARAMERR
		return err
	}

	contentHash := fmt.Sprintf("%x", sha256.Sum256(cData))

	//3. 操作数据库
	err = db.UploadFile(config.ServerConfig.DB.DBName, config.ServerConfig.DB.BlogTab, contentHash, cData)
	if err != nil {
		c.Logger().Error("failed to UploadFile file ", err, file.Filename)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//4. 返回一个url
	resp.Data = "/content/" + contentHash
	return nil
}

//文件下载
func DownloadContent(c echo.Context) error {
	//1. 解析要请求的文件
	contentHash := c.Param("hash")
	//2. 从mongodb下载写到流
	return db.DownloadFile(config.ServerConfig.DB.DBName, config.ServerConfig.DB.BlogTab, contentHash, c)
}

//发表博客
func PublishBlog(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)

	//2. 解析请求消息
	//token + content
	token := c.FormValue("token")
	content := c.FormValue("content")
	//3. 验证token
	user := &db.User{}
	user.Token = token
	err := db.QueryOne(config.ServerConfig.DB.DBName, config.ServerConfig.DB.UserTab, bson.M{"token": user.Token}, user)
	if err != nil {
		c.Logger().Error("failed to get  user", err, user.UserID)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//4. 判断正文大小
	cData := []byte(content)
	clen := len(cData)
	contentHash := fmt.Sprintf("%x", sha256.Sum256(cData))
	bloginfo := &db.BlogInfo{}
	bloginfo.UserID = user.UserID
	if clen > 300 {
		//存文件,数据库
		err = db.UploadFile(config.ServerConfig.DB.DBName, config.ServerConfig.DB.BlogTab, contentHash, cData)
		if err != nil {
			c.Logger().Error("failed to UploadFile file ", err, contentHash)
			resp.Errno = utils.RECODE_DBERR
			return err
		}
		bloginfo.Text = string(cData[:300])
		bloginfo.IsAll = false
		bloginfo.Url = "/content/" + contentHash

	} else {
		//数据库操作
		bloginfo.Text = content
		bloginfo.IsAll = true
		bloginfo.Url = ""
	}
	//5. 插入数据库
	err = db.Insert(config.ServerConfig.DB.DBName, config.ServerConfig.DB.BlogTab, bloginfo)
	if err != nil {
		c.Logger().Error("failed to insert into  blogs ", err, bloginfo)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//6. 赠送积分
	err = db.Update(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.UserTab,
		bson.M{"userid": bloginfo.UserID},
		bson.M{"$inc": bson.M{"points": 10}},
	)
	if err != nil {
		c.Logger().Error("failed to insert update  users ", err, bloginfo)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//7. detail 新增
	detail := db.PointsDetail{bloginfo.UserID, "博客发表", 10}
	err = db.Insert(config.ServerConfig.DB.DBName, config.ServerConfig.DB.DetailTab, &detail)
	if err != nil {
		c.Logger().Error("failed to insert into detail", err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	return nil
}

//查看博客
func GetBlogs(c echo.Context) error {
	//1. 组织响应消息
	var resp utils.Resp
	resp.Errno = "0"
	defer ResponseData(c, &resp)
	//2. 获得userid
	userID := c.Param("userid")
	if userID == "" {
		c.Logger().Error("failed to get userID")
		resp.Errno = utils.RECODE_PARAMERR
		return errors.New("not found user")
	}
	//3. 查询数据库 - username, []blogs
	var respBlogs RespBlogs
	err := db.QueryOne(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.UserTab,
		bson.M{"userid": userID},
		&respBlogs,
	)
	if err != nil {
		c.Logger().Error("failed to query user", userID, err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	//4. 查看所有的博客
	err = db.QueryAll(config.ServerConfig.DB.DBName,
		config.ServerConfig.DB.BlogTab,
		bson.M{"userid": userID},
		&respBlogs.Blogs,
		0,
		0,
	)
	if err != nil {
		c.Logger().Error("failed to query blogs", userID, err)
		resp.Errno = utils.RECODE_DBERR
		return err
	}
	resp.Data = respBlogs
	return nil

}
