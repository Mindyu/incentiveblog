package db

import (
	"fmt"
	"incentiveblog/config"
	"os"

	"gopkg.in/mgo.v2"
)

type User struct {
	UserID   string `json:"userid"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
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
