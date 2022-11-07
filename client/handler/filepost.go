package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	file2 "sync/file"
)

var filestatus = make(map[string]int)
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func filepost(w http.ResponseWriter,r *http.Request)  {
	var Path = os.Getenv("datadir")
	ok,err := PathExists(Path)
	if err != nil{
		fmt.Println("PathExists err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		os.MkdirAll(Path,os.ModePerm)
	}
	http_body,_ := ioutil.ReadAll(r.Body)
	var file file2.File
	err = json.Unmarshal(http_body,&file)
	if err != nil{
		log.Println("json.Unmarshal err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	file.Name = strings.Replace(file.Name,file.Path,Path,-1)
	if file.Operation == "create" {
		err := sync(file, os.O_CREATE)
		fmt.Println("创建数据")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		filestatus[file.Name] = 1
		w.WriteHeader(http.StatusOK)
		return
	}
	if _,ok := filestatus[file.Name];ok && file.Shard == 0{
		err = sync(file,os.O_APPEND)
		fmt.Println("添加数据")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		delete(filestatus,file.Name)
		w.WriteHeader(http.StatusOK)
		return
	}
	if _,ok := filestatus[file.Name];ok && file.Shard != 0{
		err = sync(file,os.O_APPEND)
		fmt.Println("添加数据")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if file.Shard == file.Shards{
			delete(filestatus,file.Name)
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	if file.Shard == 0 {
		err := create(file)
		fmt.Println("创建数据1 ",string(file.Date))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
	err = create(file)
	fmt.Println("创建数据2")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	filestatus[file.Name] = 1
	w.WriteHeader(http.StatusOK)
	return
}
