package server

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	file2 "sync/file"
	"sync/server/tools"
)
var filestatus = make(map[string]int)

func Event(watch *fsnotify.Watcher, ch chan int, opendel bool, excludes []string, addr, basePath string,files map[string]*os.File) {
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
					if tools.Excluddir(ev.Name, excludes) {
						watch.Remove(ev.Name)
						continue
					}
					ok := tools.IsDir(ev.Name)
					if ok {
						watch.Add(ev.Name)
						log.Println("创建目录 : ", ev.Name)
						e := tools.NilDir(ev.Name, watch, excludes, addr, basePath)
						if e != nil {
							ch <- 1
						}
						//path := watch.WatchList()
						//fmt.Println("path: ", path)
					} else {
						log.Println("创建文件 : ", ev.Name)
						f,_ :=os.Stat(ev.Name)
						file := file2.NewFile(basePath)
						if f.Size() > 0 {
							fe,_ := os.Open(ev.Name)
							fe.Read(file2.Bufs)
							fe.Close()
							file.Date = file2.Bufs[:f.Size()]
						}
						file.Name = ev.Name
						file.Operation = "create"
						file.Sendfile(addr)
						filestatus = map[string]int{
							ev.Name:0,
						}
						//path := watch.WatchList()
						//fmt.Println("path: ", path)
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					ok := tools.IsDir(ev.Name)
					if ok {
						file := file2.NewFile(basePath)
						file.Name = ev.Name
						file.Senddir(addr)
						continue
					}

                    if len(files) == 0{
                    	f,_ := os.Open(ev.Name)
						files[ev.Name] = f
					}else if _,ok := files[ev.Name];!ok{
						f,_ := os.Open(ev.Name)
						files[ev.Name] = f
					}
					status := Writefile(files,addr,basePath)
					if status != 0 {
						ch <- 1
						return
					}
					fmt.Println("写入文件")
				}

				if ev.Op&fsnotify.Remove == fsnotify.Remove && opendel {
					var file = &file2.File{
						Name: ev.Name,
					}
					file.Delete(addr)

					log.Println("删除文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					var file = &file2.File{
						Name: ev.Name,
					}
					file.Delete(addr)

					log.Println("重命名文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-watch.Errors:
			{
				log.Println("error : ", err)
				ch <- 1
				return
			}
		}
	}
}
