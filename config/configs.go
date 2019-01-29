package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	HttpConf *HttpConfig
	DBConf   *DBConfig
	VerConf  *VersionConfig
}

type DBConfig struct {
	Connstr string
	DBName  string
	UserTab string
	BlogTab string
}

type HttpConfig struct {
	IP   string
	Port string
}

type VersionConfig struct {
	Auth       string
	Version    string
	CreateTime string
}

var ServConf *ServerConfig

func init() {
	//fmt.Println("init=====", flag.Args())
	if getConfig() != nil {
		fmt.Println("failed to config")
		os.Exit(-1)
	}
}

func printUsage() {
	fmt.Println("Usage: %s -c config_file [-v] [-h]\n", os.Args[0])
}
func printVersion() {
	fmt.Println("author  :", ServConf.VerConf.Auth)
	fmt.Println("version :", ServConf.VerConf.Version)
	fmt.Println("time    :", ServConf.VerConf.CreateTime)
}

func getConfig() error {
	fmt.Println(os.Args, flag.Args())

	configFile := flag.String("c", "", "config_confile")
	fmt.Println(*configFile)
	var ver = flag.Bool("v", false, "version")
	var help = flag.Bool("h", false, "help")
	flag.Parse()
	fmt.Println(*help, "get help")
	ServConf = &ServerConfig{}
	if *ver {
		ServConf.VerConf = &VersionConfig{"yekai", "1.0.0", "2019-2-4"}

		printVersion()
		os.Exit(0)
	}
	if *help {
		printUsage()
		os.Exit(0)
	}

	_, err := toml.DecodeFile(*configFile, &ServConf)
	if err != nil {
		fmt.Println("failed to decode config file", *configFile)
		return err
	}
	return nil
}
