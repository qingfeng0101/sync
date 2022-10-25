package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	file2 "sync/file"
)
// 判断创建文件是否为目录
func IsDir(path string) bool {
	f,e := os.Stat(path)
	if e != nil{
		log.Println("os.Stat err: ",e)
		return false
	}
	return f.IsDir()
}
// 遍历目录加入watch
func NilDir(path string,watch *fsnotify.Watcher)  {
	f,e := ioutil.ReadDir(path)
	if e != nil{
		log.Println("ioutil.ReadDir err: ",e)
		return
	}

	if len(f) == 0{
		watch.Add(path)
		var file =  &file2.File{
			Name: path,
		}
		file.Senddir()
		return
	}
	for _,dir := range f{
		if dir.IsDir(){
			watch.Add(path +"/"+dir.Name())
			var file =  &file2.File{
				Name: path +"/"+dir.Name(),
			}
			file.Senddir()
			NilDir(path +"/"+dir.Name(),watch)
		}else {
			f, e := os.Open(path +"/"+dir.Name())
			if e != nil{
				fmt.Println("open file err: ",e)
				return
			}
			s,_ := f.Stat()
			buf := make([]byte,s.Size())
			f.Read(buf)
			var file =  &file2.File{
				Name: path +"/"+dir.Name(),
				Date: buf,
			}
			file.Sendfile()
		}
		continue
	}
	return
}


func main()  {
	basePath := os.Getenv("datadir")
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
	NilDir(basePath,watch)
	//我们另启一个goroutine来处理监控对象的事件
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					//判断事件发生的类型，如下5种
					// Create 创建
					// Write 写入
					// Remove 删除
					// Rename 重命名
					// Chmod 修改权限
					if ev.Op&fsnotify.Create == fsnotify.Create {
						ok := IsDir(ev.Name)
						if ok{
							watch.Add(ev.Name)
							log.Println("创建目录 : ", ev.Name);
							NilDir(ev.Name,watch)
							path := watch.WatchList()
							fmt.Println("path: ",path)
						}else {
							log.Println("创建文件 : ", ev.Name);
							f, e := os.Open(ev.Name)
							if e != nil{
								fmt.Println("open file err: ",err)
								return
							}
							s,_ := f.Stat()
							buf := make([]byte,s.Size())
							f.Read(buf)
							var file =  &file2.File{
								Name: ev.Name,
								Date: buf,
							}
							file.Sendfile()

							path := watch.WatchList()
							fmt.Println("path: ",path)
						}

					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						f, e := os.Open(ev.Name)
						if e != nil{
							fmt.Println("open file err: ",err)
							return
						}
						s,_ := f.Stat()
						buf := make([]byte,s.Size())
						f.Read(buf)
						var file =  &file2.File{
							Name: ev.Name,
							Date: buf,
						}
						file.Sendfile()
						fmt.Println("写入文件")
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						var file =  &file2.File{
							Name: ev.Name,
						}
						file.Delete()

						log.Println("删除文件 : ", ev.Name);
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						var file =  &file2.File{
							Name: ev.Name,
						}
						file.Delete()

						log.Println("重命名文件 : ", ev.Name);
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("修改权限 : ", ev.Name);
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err);
					return;
				}
			}
		}
	}();

	//循环
	select {};
}
