package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
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
var pathdir string
func NilDir(path string,watch *fsnotify.Watcher,excludes []string,addr,basePath string) (error) {
		f,e := ioutil.ReadDir(path)

		if e != nil{
			log.Println("ioutil.ReadDir err: ",e)
			return e
		}
		if len(f) == 0 && !Excluddir(path,excludes){
			watch.Add(path)
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
				watch.Add(pathdir)
				file := file2.NewFile(basePath)
				file.Name = pathdir
				file.Senddir(addr)
				NilDir(pathdir,watch,excludes,addr,basePath)
			} else if !dir.IsDir(){

				ok,err := DataSize(pathdir,file2.Buf)
				if err != nil{
					fmt.Println("NilDir DataSize err: ",err)
					return err
				}
				if ok{
					ShardData(pathdir,addr,basePath)
					continue
				}
				read, e := os.Open(pathdir)
				if e != nil{
					fmt.Println("open file err: ",e)
					return e
				}
				s,_ := read.Stat()
				buf := make([]byte,s.Size())
				read.Read(buf)
				read.Close()
				file := file2.NewFile(basePath)
				file.Name = pathdir
				file.Date = buf
				file.Sendfile(addr)

			}
			continue
		}

	return nil
}
// 判断文件大小是否使用分片
// 重试打开文件次数
var restart = 0
func DataSize(path string,size int64) (status bool,err error) {
	status = false
	err = nil
		f,e := os.Open(path)
		defer f.Close()
		if e != nil{
			err = e
			return
		}
		info, e := f.Stat()
		if e != nil{
			fmt.Println("DataSize os.Open err: ",e)
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
		fmt.Println("ShardData f.Stat: ",err)
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
			fmt.Printf("数据同步失败，切片：%d,文件名：%s", file.Shard, file.Name)
			break
		}


	}
	return 0
}
// 分片传入后端服务
func ShardData(path,addr,basePath string) int {
    f,err := os.Open(path)
    if err != nil{
    	fmt.Println("ShardData os.Open err: ",err)
		return 0
	}
	defer f.Close()
    info,err := f.Stat()
    if err != nil{
    	fmt.Println("ShardData f.Stat: ",err)
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
			fmt.Printf("数据同步失败，切片：%d,文件名：%s", file.Shard, file.Name)
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
	var ostype = runtime.GOOS
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

func Gethash(path string) (hash string,err error) {
	file, _ := os.Open(path)
	h_ob := sha256.New()
	_, err = io.Copy(h_ob, file)
	if err == nil {
		hash := h_ob.Sum(nil)
		hashvalue := hex.EncodeToString(hash)
		return hashvalue,nil
	} else {
		return "哈希错误",err
	}
}
