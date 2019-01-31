package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

/*
[httpconf]
ip = "localhost"
port = ":8086"

[db]
connstr = "localhost"
dbname = "blogs"
usertab = "users"
blogtab = "blogs"
*/

type ServerConf struct {
	Httpconf HttpServer
	DB       Database
}

type HttpServer struct {
	IP   string
	Port string
}

type Database struct {
	Connstr     string
	DBName      string
	UserTab     string
	BlogTab     string
	DetailTab   string
	RelationTab string
}

var ServerConfig *ServerConf

func printHelp() {
	fmt.Printf("%s -c config_file [-h] [-v]\n", os.Args[0])
}

func printVersion() {
	fmt.Printf("auth   :%s\n", "yekai")
	fmt.Printf("version:%s\n", "1.0.0")
	fmt.Printf("time   :%s\n", "2019-2-4")
}

//当config被包含时自动调用，而且多次包含只会执行一次
func init() {
	if getConfig() != nil {
		os.Exit(-1)
	}
	fmt.Println(*ServerConfig)
}

//读取配置文件，解析配置文件
func getConfig() error {
	//1. 解析命令行参数 flag
	configFile := flag.String("c", "", "config_file")
	help := flag.Bool("h", false, "help")
	ver := flag.Bool("v", false, "version")
	flag.Parse() //执行解析
	if *help {
		printHelp()
		os.Exit(0)
	}
	if *ver {
		printVersion()
		os.Exit(0)
	}
	//2. 解析配置文件
	if *configFile == "" {
		printHelp()
		return errors.New("no config file")
	} else {
		ServerConfig = &ServerConf{}
		_, err := toml.DecodeFile(*configFile, ServerConfig)
		if err != nil {
			fmt.Println("failed to decode config file ", configFile, err)
			return err
		}
	}
	return nil
}
