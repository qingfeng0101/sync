package server

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"sync/conf"
	file2 "sync/file"
	"sync/server/tools"
	"syscall"
)
func Event(channels *tools.Channels, excludes []string, conf *conf.Config) {
	for {
		select {
		case ev := <-channels.Watch.Events:
			{
				//判断事件发生的类型，如下5种
				// Create 创建
				// Write 写入
				// Remove 删除
				// Rename 重命名
				// Chmod 修改权限
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if tools.Excluddir(ev.Name, excludes) {
						continue
					}
					ok := tools.IsDir(ev.Name)
					if ok {
						err := channels.Watch.Add(ev.Name)
						if err != nil {
							log.Println("添加目录监听失败：", err)
							channels.EndChan <- 1
							return
						}
						e := tools.NilDir(channels, ev.Name, excludes, conf.Clientaddr, conf.DataDIr)
						if e != nil {
							log.Println("遍历目录异常11 err：", e)
							channels.EndChan <- 1
							return
						}
					} else {
						fmt.Println("IsTemp: ",tools.IsTemp(ev.Name))
						if !tools.IsTemp(ev.Name) {
							log.Println("创建文件 : ", ev.Name)
							var ChenData tools.ChenData
							ChenData.Name = ev.Name
							ChenData.Operation = "create"
							channels.ChanDatas <- &ChenData
						}
					}

					if ev.Op&fsnotify.Write == fsnotify.Write {
						ok := tools.IsDir(ev.Name)
						if ok {
							file := file2.NewFile(conf.DataDIr)
							file.Name = ev.Name
							file.Senddir(conf.Clientaddr)
							continue
						}
						if !tools.IsTemp(ev.Name) {
							var s tools.ChenData
							s.Name = ev.Name
							channels.ChanDatas <- &s
							fmt.Println("写入文件")
						}
					}

					if ev.Op&fsnotify.Remove == fsnotify.Remove && conf.Delete {
						ok := tools.IsDir(ev.Name)
						if ok {
							channels.Watch.Remove(ev.Name)
						}
						file := file2.NewFile(conf.DataDIr)
						file.Name = ev.Name
						file.Delete(conf.Clientaddr)

						log.Println("删除文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						var file = &file2.File{
							Name: ev.Name,
						}
						file.Delete(conf.Clientaddr)

						log.Println("重命名文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("修改权限 : ", ev.Name)
					}
				}
			}
		case Sigs:= <-channels.Sigs:
			for n:=0;n<conf.SaveThread;n++{
				channels.SaveStop <- Sigs
			}
			return

		case err := <-channels.Watch.Errors:
			{
			if err != nil{
				log.Println("error : ", err)
			}
			for n:=0;n<conf.SaveThread;n++{
					channels.SaveStop <- syscall.SIGKILL
				}
			return
			}
		}
	}
}
