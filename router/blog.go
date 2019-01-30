package router

import (
	"crypto/sha256"
	"fmt"
	"incentiveblog/config"
	"incentiveblog/db"
	"incentiveblog/utils"
	_ "net/http"

	"github.com/labstack/echo"
	_ "gopkg.in/mgo.v2/bson"
)

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
