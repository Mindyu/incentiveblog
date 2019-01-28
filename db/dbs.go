package db

import (
	"fmt"
	"io"
	"log"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"
)

type DBConnect struct {
	Sess *mgo.Session
}

var Dbconn *DBConnect

type User struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	IconUrl  string `json:"iconurl"`
	Token    string `json:"token"`
	Points   int    `json:"points"`
}

type UserRelation struct {
	FromID string `json:"froid"`
	ToID   string `json:"toid"`
}

func init() {
	Dbconn = &DBConnect{}
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal("failed to connect to mongodb", err)
	}
	Dbconn.Sess = sess

	fmt.Println("connect to mongodb ok")
}

func (c *DBConnect) Insert(DBName, collection string, data interface{}) error {
	return c.Sess.DB(DBName).C(collection).Insert(data)
}

func (c *DBConnect) Update(DBName, collection string, cond, data interface{}) error {
	return c.Sess.DB(DBName).C(collection).Update(cond, data)
}

func (c *DBConnect) QueryOne(DBName, collection string, cond, result interface{}) error {
	fmt.Println(cond)
	err := c.Sess.DB(DBName).C(collection).Find(cond).One(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *DBConnect) QueryAll(DBName, collection string, cond, result interface{}, skip, limit int) error {

	return c.Sess.DB(DBName).C(collection).Find(cond).Skip(skip).Limit(limit).All(result)
}

func (c *DBConnect) UploadFile(DBName, collection, fileName string, data []byte) error {

	fmt.Println(DBName, collection, fileName)
	f, err := c.Sess.DB(DBName).GridFS(collection).Create(fileName)
	if err != nil {
		fmt.Println("failed to create gridfs file", err)
		return err
	}
	_, err = f.Write(data)
	f.Close()
	return err
}

func (c *DBConnect) DownLoadFile(DBName, collection, fileName string, ec echo.Context) error {
	fmt.Println("DownLoadFile", fileName, collection, DBName)
	f, err := c.Sess.DB(DBName).GridFS(collection).Open(fileName)
	if err != nil {
		fmt.Println("failed to open gridfs file", err, fileName)
		return err
	}
	defer f.Close()
	_, err = io.Copy(ec.Response(), f)
	if err != nil {
		fmt.Println("failed to copy file", err)
		return err
	}
	return err
}

func (c *DBConnect) GetFileText(DBName, collection, fileName string) ([]byte, error) {
	fmt.Println("GetFileText", fileName, collection, DBName)
	f, err := c.Sess.DB(DBName).GridFS(collection).Open(fileName)
	if err != nil {
		fmt.Println("failed to open gridfs file", err, fileName)
		return nil, err
	}
	defer f.Close()
	data := make([]byte, f.Size())
	n, err := f.Read(data)
	if err != nil || int64(n) != f.Size() {
		fmt.Println("failed to read data", err)
		return data, err
	}
	return data, nil
}

func (c *DBConnect) Count(DBName, collection string, cond interface{}) (int, error) {
	return c.Sess.DB(DBName).C(collection).Find(cond).Count()
}
