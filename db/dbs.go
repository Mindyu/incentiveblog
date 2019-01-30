package db

import (
	"fmt"
	"incentiveblog/config"
	"io"
	"os"

	"github.com/labstack/echo"

	"gopkg.in/mgo.v2"
)

type User struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	Token    string `json:"token"`
	Points   int    `json:"points"`
}

type PointsDetail struct {
	UserID string `json:"userid"`
	Remark string `json:"remark"`
	Points int    `json:"points"`
}

var DbConn *mgo.Session

func init() {
	sess, err := mgo.Dial(config.ServerConfig.DB.Connstr)
	if err != nil {
		fmt.Println("failed to connect to mongodb", err)
		os.Exit(-1)
	}
	DbConn = sess
}

func Insert(dbname, tabname string, data interface{}) error {
	collect := DbConn.DB(dbname).C(tabname)
	return collect.Insert(data)
}

func QueryOne(dbname, tabname string, cond, result interface{}) error {
	collect := DbConn.DB(dbname).C(tabname)
	return collect.Find(cond).One(result)
}

func Update(dbname, tabname string, cond, data interface{}) error {
	collect := DbConn.DB(dbname).C(tabname)
	return collect.Update(cond, data)
}

//上传
func UploadFile(dbname, tabname, filename string, data []byte) error {
	//GridFS 是mongodb提供的分布式文件存储组件
	f, err := DbConn.DB(dbname).GridFS(tabname).Create(filename)
	if err != nil {
		fmt.Println("failed to open grid file", filename, err)
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

//下载
func DownloadFile(dbname, tabname, filename string, c echo.Context) error {
	f, err := DbConn.DB(dbname).GridFS(tabname).Open(filename)
	if err != nil {
		fmt.Println("failed to open grid file", filename, err)
		return err
	}
	defer f.Close()

	_, err = io.Copy(c.Response(), f)
	if err != nil {
		fmt.Println("failed to copy file", err)
		return err
	}
	return err
}
