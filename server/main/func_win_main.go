// +build !linux

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync/conf"
	"sync/server/server"
	"sync/server/tools"
	"syscall"
)
func main()  {
	var config string
	flag.StringVar(&config,"f","./server.conf","指定服务端配置文件")
	flag.Parse()
    conf := conf.NewConfing(config)
    if conf == nil{
    	log.Println("服务异常")
		return
	}
	re,_:= regexp.Compile("#")
	conf.DataDIr = tools.RewritePath(conf.DataDIr)
	if conf.SaveFile != "" &&  re.FindString(conf.SaveFile) != ""{
		conf.SaveFile = ""
	}
	basePath := conf.DataDIr
	excludes := strings.Split(conf.Exclude,",")
	// 调用管道初始化集合函数
	Channels := tools.NewChannels()
	// 通知子进程关闭
	signal.Notify(Channels.Sigs, os.Interrupt,os.Kill,  syscall.SIGINT, syscall.SIGTERM)
	// 关闭监听
	defer Channels.Watch.Close()
	//添加要监控的对象，文件或文件夹
	err := Channels.Watch.Add(conf.DataDIr);
	if err != nil {
		log.Fatal(err);
	}
    // 加载同步过的文件
    var savedata *tools.SaveDatas
	savedata = tools.Init(conf.SaveFile)
	//启动监听记录数据的goroutine
	for n:=0;n<conf.SaveThread;n++{
		go tools.SaveData(Channels,savedata,conf)
	}
	//我们另启一个goroutine来处理监控对象的事件
	go server.Event(Channels,excludes,conf)
	//遍历当前监听的目录，全量数据同步一次
	err = tools.NilDir(Channels,basePath,excludes,conf.Clientaddr,basePath)
	if err != nil{
		log.Println("遍历目录异常 err： ",err)
		return
	}
	if savedata.SavePath != "" {
		go tools.CronData(Channels.DataChan, savedata)
	}
	for n:=0;n<conf.SaveThread;n++{
       <-Channels.EndChan
	}
	//
	log.Println("服务异常退出")
}
