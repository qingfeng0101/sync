package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync/conf"
	file2 "sync/file"
)
var ostype = runtime.GOOS
var overnum int64
// 判断创建文件是否为目录
func IsDir(path string) bool {
	f,e := os.Stat(path)
    if os.IsNotExist(e){
		log.Println("IsNotExist: ",path)
		return false
	}
	if e != nil{
		log.Println("os.Stat err: ",e)
		return false
	}
	return f.IsDir()
}
// 遍历目录加入watch
var m = 0
var pathdir string
func NilDir(channels *Channels,path string,excludes []string,addr,basePath string) (error) {
	    fmt.Println("数目： ",m)
		f,e := ioutil.ReadDir(path)
		if e != nil{
			log.Println("ioutil.ReadDir err: ",e)
			return e
		}
		if len(f) == 0 && !Excluddir(path,excludes){
			err := channels.Watch.Add(path)
			if err != nil{
				return err
			}
			file := file2.NewFile(basePath)
			file.Name = path
			file.Senddir(addr)
			return nil
		}
		for _,dir := range f{
			var ostype = runtime.GOOS
			if ostype == "windows"{
				pathdir = path +"\\"+dir.Name()
			}else if ostype == "linux"{
				pathdir = path +"/"+dir.Name()
			}
			if dir.IsDir() &&!Excluddir(pathdir,excludes){
				err := channels.Watch.Add(pathdir)
				if err != nil{
					return err
				}
				file := file2.NewFile(basePath)
				file.Name = pathdir
				file.Senddir(addr)
				NilDir(channels,pathdir,excludes,addr,basePath)
			} else if !dir.IsDir(){
				finfo,err := os.Stat(pathdir)
				if os.IsNotExist(err){
					log.Println("IsNotExist filename：",pathdir)
					continue
				}
				var ch ChenData
				ch.Name = pathdir
				ch.Value = finfo.ModTime().Unix()
				channels.ChanDatas <- &ch
				m++
			}
			continue
		}
	return nil
}
// 发送小文件
func SmallData(ChenData *ChenData,conf *conf.Config)  {
	read, e := os.Open(ChenData.Name)
	if e != nil{
		log.Println("open file err: ",e)
	}
	//s,_ := read.Stat()
	buf := make([]byte,ChenData.Value)
	read.Read(buf)
	read.Close()
	file := file2.NewFile(conf.DataDIr)
	file.Name = ChenData.Name
	file.Date = buf
	if ChenData.Operation != ""{
		file.Operation = ChenData.Operation
	}
	file.Sendfile(conf.Clientaddr)
}
// 判断文件大小是否使用分片
// 重试打开文件次数
var restart = 0
func DataSize(path string,size int64) (status bool,err error) {
	status = false
	err = nil
	info, e := os.Stat(path)
	if e != nil{
		log.Println("DataSize os.Open err: ",e)
		status = false
		err = e
		return
	}
	if info.Size() > size{
		status = true
		err = nil
		return true,nil
	}
	return
}
// 分片
func ShardData1(f *os.File,path,addr,basePath string) int {
	defer f.Close()
	info,err := f.Stat()
	if err != nil{
		log.Println("ShardData f.Stat: ",err)
		return 0
	}
	file := file2.NewFile(basePath)
	file.Operation = "append"
	file.Name = path
	if info.Size() % file2.Buf == 0{
		file.Shards = info.Size() / file2.Buf
	}else {
		overnum = info.Size() % file2.Buf
		file.Shards = info.Size() / file2.Buf + 1
	}
	for {
		file.Shard +=1
		if file.Shard == file.Shards {
			if overnum != 0 {
				f.Read(file2.Bufs[:overnum])
				file.Date = file2.Bufs[:overnum]
				ok := file.Sendfile(addr)
				if !ok {
					fmt.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
					return 0
				}
				break
			}
			f.Read(file2.Bufs)
			file.Date = file2.Bufs
			ok := file.Sendfile(addr)
			if !ok {
				log.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
				return 0
			}

			break
		}

		f.Read(file2.Bufs)
		file.Date = file2.Bufs
		ok := file.Sendfile(addr)
		if !ok {
			log.Printf("数据同步失败，切片：%d,文件名：%s", file.Shard, file.Name)
			break
		}


	}
	return 0
}
// 分片传入后端服务
func ShardData(path,addr,basePath string) int {
    f,err := os.Open(path)
    if err != nil{
    	log.Println("ShardData os.Open err: ",err)
		return 0
	}
	defer f.Close()
    info,err := f.Stat()
    if err != nil{
    	log.Println("ShardData f.Stat: ",err)
		return 0
	}
	file := file2.NewFile(basePath)
    file.Operation = "append"
    file.Name = path
    if info.Size() % file2.Buf == 0{
    	file.Shards = info.Size() / file2.Buf
	}else {
		overnum = info.Size() % file2.Buf
		file.Shards = info.Size() / file2.Buf + 1
	}
	for {
		file.Shard +=1
		if file.Shard == file.Shards {
			if overnum != 0 {
				f.Read(file2.Bufs[:overnum])
				file.Date = file2.Bufs[:overnum]

				ok := file.Sendfile(addr)
				if !ok {
					fmt.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
					return 0
				}
				break
			}
			f.Read(file2.Bufs)
			file.Date = file2.Bufs
			ok := file.Sendfile(addr)
			if !ok {
				fmt.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
				return 0
			}

			break
		}

		f.Read(file2.Bufs)
		file.Date = file2.Bufs
		ok := file.Sendfile(addr)
		if !ok {
			log.Printf("数据同步失败，切片：%d,文件名：%s", file.Shard, file.Name)
			break
		}


	}
	return 0
}
// 判断目录是否排除
func Excluddir(path string,exclude []string) bool  {
	for _,name := range exclude{

		if RewritePath(name) == path{

			return true
		}
	}
	return false
}
// 去除目录结尾/
func RewritePath(path string) string {
	var fn func(rune) bool

	if ostype == "windows"{
		 fn = func(c rune) bool {
			return strings.ContainsRune("\\", c)
		}
	}else if ostype == "linux"{
		 fn = func(c rune) bool {
			return strings.ContainsRune("/", c)
		}
	}

	return strings.TrimRightFunc(path, fn)
}

func IsTemp(path string) (bool) {
	re,_ := regexp.Compile("^\\.")
	if  ostype == "windows"{
		paths := strings.Split(path,"\\")
		str := paths[len(paths)-1]
		if re.FindString(str) != ""{
			return true
		}
	}
	if  ostype == "linux"{
		paths := strings.Split(path,"/")

		str := paths[len(paths)-1]
		if re.FindString(str) != ""{
			return true
		}
	}
	return false
}
