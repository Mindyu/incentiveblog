package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("hello")
	//1. 连接
	//Dial(url string) (*Session, error)
	sess, err := mgo.Dial("localhost:27017")
	if err != nil {
		fmt.Println("failed to connect ", err)
		return
	}
	fmt.Println("connect ok!")
	//2. 增加一个人
	tablename := sess.DB("yekai").C("person") // 得到集合
	err = tablename.Insert(Person{"yekai", 30})
	if err != nil {
		fmt.Println("failed to insert ", err)
		return
	}

	//3. 修改文档
	tablename.Update(bson.M{"name": "yekai"}, Person{"yekai", 36})
	//4. 查看文档
	var persons []Person
	tablename.Find(nil).All(&persons)
	fmt.Println(persons)
	//5. 删除文档
	tablename.RemoveAll(nil)
}
