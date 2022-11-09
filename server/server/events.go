package server

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	file2 "sync/file"
	"sync/server/tools"
)

func Event(watch *fsnotify.Watcher,ch chan int,opendel bool,excludes []string,addr,basePath string) {
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
					if tools.Excluddir(ev.Name,excludes){
						watch.Remove(ev.Name)
						continue
					}

						ok := tools.IsDir(ev.Name)
						if ok {

							watch.Add(ev.Name)
							log.Println("创建目录 : ", ev.Name);
							e := tools.NilDir(ev.Name, watch,excludes,addr,basePath)
							if e != nil{
								ch <- 1
							}
							path := watch.WatchList()
							fmt.Println("path: ", path)
						} else {
							log.Println("创建文件 : ", ev.Name);
							file := file2.NewFile(basePath)
							file.Name = ev.Name
							file.Operation = "create"
							file.Sendfile(addr)
							path := watch.WatchList()
							fmt.Println("path: ", path)
						}
					}


				if ev.Op&fsnotify.Write == fsnotify.Write {
					ok , err := tools.DataSize(ev.Name,file2.Buf)
					if err != nil{
						ch <- 1
						return
					}
					if ok {
						tools.ShardData(ev.Name,addr,basePath)
						fmt.Println("00000000")
					}else {
						fmt.Println("1111111")
						f, e := os.Open(ev.Name)
						if e != nil {
							fmt.Println("open file err: ", e)
							ch <- 1
							return
						}
						s, _ := f.Stat()

						f.Read(file2.Bufs)
						file := file2.NewFile(basePath)
						file.Name = ev.Name
						file.Date = file2.Bufs[:s.Size()]
						file.Shard = 0
						file.Operation = "append"

						fmt.Println("写入数据：",string(file.Date))
						file.Sendfile(addr)
					}

					fmt.Println("写入文件")
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove && opendel {
					var file = &file2.File{
						Name: ev.Name,
					}
					file.Delete(addr)

					log.Println("删除文件 : ", ev.Name);
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					var file = &file2.File{
						Name: ev.Name,
					}
					file.Delete(addr)

					log.Println("重命名文件 : ", ev.Name);
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", ev.Name);
				}
			}
		case err := <-watch.Errors:
			{
				log.Println("error : ", err);
				ch <- 1
				return;
			}
		}
	}
}
