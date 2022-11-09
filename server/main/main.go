package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"strings"
	"sync/conf"
	"sync/server/server"
	"sync/server/tools"
)
func main()  {
	var config string
	flag.StringVar(&config,"f","./server.conf","指定服务端配置文件")
	flag.Parse()
    conf := conf.NewConfing(config)
    if conf == nil{
    	fmt.Println("服务异常")
		return
	}

	basePath := tools.RewritePath(conf.DataDIr)
	//basePath := conf.DataDIr
	excludes := strings.Split(conf.Exclude,",")
	opendel := conf.Delete

	// goroutine状态标识
	ch := make(chan int)
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher();
	if err != nil {
		log.Fatal(err);
	}
	defer watch.Close();
	//patharr := make([]string,0)
	//添加要监控的对象，文件或文件夹
	//patharr = append(patharr,"./tmp")
	//watch.WatchList()
	err = watch.Add(basePath);
	if err != nil {
		log.Fatal(err);
	}

	// 遍历当前监听的目录，全量数据同步一次
	tools.NilDir(basePath,watch,excludes,conf.Clientaddr,basePath)
	fmt.Println("watch1: ",watch.WatchList())
	//我们另启一个goroutine来处理监控对象的事件
	go server.Event(watch,ch,opendel,excludes,conf.Clientaddr,basePath)
	fmt.Println("watch2: ",watch.WatchList())
	//循环
	<-ch
	log.Println("服务异常退出")
}
