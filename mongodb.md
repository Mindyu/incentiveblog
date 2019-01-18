## MongoDB的简介
MongoDB是一个基于分布式文件存储的数据库，使用C++编写，旨在为WEB应用提供可扩展的高性能数据存储解决方案。MongoDB本身属于非关系型数据库，但是它又是非关系型中最像关系型的。
特点：
- 高性能
- 易使用
- 易部署
- 模式自由
- 动态查询（支持js）
- BSON数据存储
- 复制和数据恢复
- 索引，分片都支持

> MongoDB公司原名是10gen公司，后来更名为MongoDB公司，典型的产品比公司名气大。名字来源于Humongous，目标在于处理海量的数据。


## MongoDB的安装

[官网下载](https://www.mongodb.com/download-center/community)

![官网示意图](https://note.youdao.com/yws/public/resource/6ea92b8dadc261228e50cf0af238a0ce/xmlnote/0A15B01768CA462F834C118DE4A4033D/22737)
下载后，进行解压缩就可以直接获得可执行文件进行使用！重点关注2个文件：
- mongod  服务器启动程序
- mongo   客户端启动程序

对于类unix平台来说，可以使用命令直接安装，相对很方便！

- for mac os 
 
```
brew install mongodb
```

- for ubuntu


```
sudo apt-get install mongodb
```

如果是类unix平台，推荐使用命令行安装，执行文件默认在系统可执行路径下！

## MongoDB的增删改查

### MongoDB的组织结构

MongoDB默认没有用户验证机制，使用mongo可以直接登陆，我们可以认为MongoDB的组织结构是三层！

库 - 集合(表) - 文档(记录)

### 启动MongoDB

- 服务器启动

```
/usr/local/bin/mongod -f /usr/local/etc/mongod.conf &
```

配置文件内容如下：
```
localhost:incentiveblog yk$ cat /usr/local/etc/mongod.conf
systemLog:
  destination: file
  path: /usr/local/var/log/mongodb/mongo.log
  logAppend: true
storage:
  dbPath: /usr/local/var/mongodb
net:
  bindIp: 127.0.0.1

```


- 客户端连接

输入mongo就连接了
```
localhost:incentiveblog yk$ mongo
MongoDB shell version v3.6.4
connecting to: mongodb://127.0.0.1:27017
MongoDB server version: 3.6.4
Server has startup warnings: 
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] 
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] 
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] 
2019-01-18T11:29:05.899+0800 I CONTROL  [initandlisten] ** WARNING: soft rlimits too low. Number of files is 256, should be at least 1000
> 

```


### 文档的增删改查

- 新增文档

语法：

```
db.COLLECTION.insert(JSON) 
```

举例
```
> db.person.insert({name:"yekai",age:24,sex:"man"})
WriteResult({ "nInserted" : 1 })
> db.person.insert({name:"fuhongxue",age:"23",sex:"man"})
WriteResult({ "nInserted" : 1 })
> db.person.insert({name:"luxiaojia",age:22,sex:"man",info:{like:"huasheng",wuqi:"jian"}})
WriteResult({ "nInserted" : 1 })
```

- 查询文档

语法：

```
db.COLLECTTION.find(COND_JSON,SHOW_JSON)
```
示例：

```
> db.person.find()
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196c"), "name" : "yekai", "age" : 24, "sex" : "man" }
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196d"), "name" : "fuhongxue", "age" : "23", "sex" : "man" }
{ "_id" : ObjectId("5c414aab466f6a72a6e5196e"), "name" : "luxiaojia", "age" : 22, "sex" : "man", "info" : { "like" : "huasheng", "wuqi" : "jian" } }
```

```
> db.person.find({name:"yekai"})
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196c"), "name" : "yekai", "age" : 24, "sex" : "man" }
> db.person.find({name:"yekai"},{_id:0})
{ "name" : "yekai", "age" : 24, "sex" : "man" }
```




- 修改文档

类比关系型语法：update tablename set col= val where cond;
语法：

```
db.COLLECTION.update(COND_JSON,UPDATE_JSON)
```

示例1，想把yekai的年龄修改为30

```
> db.person.update({name:"yekai"},{age:30})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
```

修改后yekai消失了，结论后面的json要写全部信息。
```
> db.person.update({name:"yekai"},{age:30})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.person.find()
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196c"), "age" : 30 }
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196d"), "name" : "fuhongxue", "age" : "23", "sex" : "man" }
{ "_id" : ObjectId("5c414aab466f6a72a6e5196e"), "name" : "luxiaojia", "age" : 22, "sex" : "man", "info" : { "like" : "huasheng", "wuqi" : "jian" } }

```

示例2，也可以使用$set操作符，进行精确修改


```
> db.person.update({name:"fuhongxue"},{$set:{age:30}})
WriteResult({ "nMatched" : 1, "nUpserted" : 0, "nModified" : 1 })
> db.person.find()
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196c"), "age" : 30 }
{ "_id" : ObjectId("5c414aa9466f6a72a6e5196d"), "name" : "fuhongxue", "age" : 30, "sex" : "man" }
{ "_id" : ObjectId("5c414aab466f6a72a6e5196e"), "name" : "luxiaojia", "age" : 22, "sex" : "man", "info" : { "like" : "huasheng", "wuqi" : "jian" } }

```
实验证明，这样ok！




- 删除文档

语法：

```
db.COLLECTION.remove(COND_JSON)
```

示例：删除年龄是30的人


```
> db.person.remove({age:30})
WriteResult({ "nRemoved" : 2 })
> db.person.find()
{ "_id" : ObjectId("5c414aab466f6a72a6e5196e"), "name" : "luxiaojia", "age" : 22, "sex" : "man", "info" : { "like" : "huasheng", "wuqi" : "jian" } }
```

总结：MongoDB对json非常依赖，我们只是简单的介绍了几个例子，有些更复杂的写法留给读者自学。

