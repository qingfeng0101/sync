package tools

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	file2 "sync/file"
)
var overnum int64
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
func NilDir(path string,watch *fsnotify.Watcher) (error) {
	f,e := ioutil.ReadDir(path)
	if e != nil{
		log.Println("ioutil.ReadDir err: ",e)
		return e
	}
	if len(f) == 0{
		watch.Add(path)
		file := file2.NewFile()
		file.Name = path
		file.Senddir()
		return nil
	}
	for _,dir := range f{
		if dir.IsDir(){
			watch.Add(path +"/"+dir.Name())
			file := file2.NewFile()
			file.Name = path +"/"+dir.Name()
			file.Senddir()
			NilDir(path +"/"+dir.Name(),watch)
		}else {
			fmt.Println("file name: ",path +"/"+dir.Name())
			ok,err := DataSize(path +"/"+dir.Name(),file2.Buf)
			if err != nil{
				fmt.Println("NilDir DataSize err: ",err)
				return err
			}
			if ok{
				ShardData(path +"/"+dir.Name())
				continue
			}
			f, e := os.Open(path +"/"+dir.Name())
			if e != nil{
				fmt.Println("open file err: ",e)
				return e
			}
			s,_ := f.Stat()
			buf := make([]byte,s.Size())
			f.Read(buf)
			file := file2.NewFile()
			file.Name = path +"/"+dir.Name()
			file.Date = buf
			file.Sendfile()
		}
		continue
	}
	return nil
}
// 判断文件大小是否使用分片
func DataSize(path string,size int64) (bool,error) {
	f,e := os.Open(path)
	if e != nil{
		fmt.Println("DataSize os.Open err: ",e)
		return false,e
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil{
		fmt.Println("DataSize f.Stat err: ",e)
		return false,e
	}
	if info.Size() > size{
		fmt.Println("大小1：",info.Size()  )
		return true,nil
	}
	fmt.Println("大小2：",info.Size()  )
	return false,nil
}
// 分片传入后端服务
func ShardData(path string)  {
    f,err := os.Open(path)
    if err != nil{
    	fmt.Println("ShardData os.Open err: ",err)
		return
	}
	defer f.Close()
    info,err := f.Stat()
    if err != nil{
    	fmt.Println("ShardData f.Stat: ",err)
		return
	}
	file := file2.NewFile()
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
				fmt.Println("最后长度：",len(file.Date))
				fmt.Println("最后字符：",string(file.Date))
				fmt.Println("总分片数：",file.Shards)
				fmt.Println("当前分片数：",file.Shards)
				ok := file.Sendfile()
				if !ok {
					fmt.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
					return
				}
				break
			}
			f.Read(file2.Bufs)
			file.Date = file2.Bufs
			ok := file.Sendfile()
			if !ok {
				fmt.Printf("数据同步失败，切片：%d,文件名：%s",file.Shard,file.Name)
				return
			}

			break
		}

		f.Read(file2.Bufs)
		file.Date = file2.Bufs
		ok := file.Sendfile()
		if !ok {
			fmt.Printf("数据同步失败，切片：%d,文件名：%s", file.Shard, file.Name)
			break
		}

		fmt.Println("pppppppp")
	}
}