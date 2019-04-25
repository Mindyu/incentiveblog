package db

import (
	"fmt"
	"github.com/mindyu/incentiveblog/config"
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

type QueryDetail struct {
	UserID string `json:"userid"`
	Skip   int    `json:"skip"`
	Limit  int    `json:"limit"`
}

type BlogInfo struct {
	UserID string `json:"userid"`
	Text   string `json:"text"`
	Url    string `json:"url"`
	IsAll  bool   `json:"isall"`
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

//skip代表跳过记录，limit代表取记录数 skip=10,limit=10 ==> 11-20
func QueryAll(dbname, tabname string, cond, result interface{}, skip, limit int) error {
	collect := DbConn.DB(dbname).C(tabname)
	return collect.Find(cond).Skip(skip).Limit(limit).All(result)
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
