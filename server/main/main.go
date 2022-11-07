package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"strconv"
	"sync/server/server"
	"sync/server/tools"
)
func main()  {
	basePath := os.Getenv("datadir")
	opendel,err := strconv.ParseBool(os.Getenv("delete"))
	if err != nil{
		fmt.Println("输入类型有误，err: ",err)
		return
	}
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
	tools.NilDir(basePath,watch)
	//我们另启一个goroutine来处理监控对象的事件
	go server.Event(watch,ch,opendel)
	//循环
	<-ch
	log.Println("服务异常退出")
}
