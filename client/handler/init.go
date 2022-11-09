package handler

import (
	"flag"
	"sync/conf"
)

var Client *conf.ClientConf
func Init() {
	var config string
	flag.StringVar(&config,"f","./server.conf","指定服务端配置文件")
	flag.Parse()
	Client = conf.NewClient(config)
	return
}