package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

const Buf = 1048576
var Bufs = make([]byte, Buf)
type File struct {
	Name string `json:"name"`
	Date []byte `json:"date"`
	Shard int64 `json:"shard"`
	Shards int64 `json:"shards"`
	Size  int64 `json:"size"`
	Operation string `json:"operation"`
    Path string `json:"path"`
	Systype string `json:"systype"`
}

func NewFile(basePath string) (file *File) {
	var ostype = runtime.GOOS
	return &File{
		Path: basePath,
		Systype: ostype,
	}
}
func (f *File) Sendfile(addr string,timeout time.Duration) bool {
	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("encoder.Encode err: ",e)
		return false
	}
	for {
		re,e := http.Post("http://"+addr+"/file/","application/json;utf-8",bytes.NewReader(b))
		if e != nil{
			fmt.Println("post err: ",e)
			time.Sleep(time.Second*timeout)
			continue
		}
		re.Body.Close()
		if re.StatusCode != http.StatusOK {
			return false
		}
		b=[]byte{}
		return true
	}

}
func (f *File) Senddir(addr string)  {
	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("json.Marshal err: ",e)
		return
	}
	re,e := http.Post("http://"+addr+"/dir/","application/octet-stream",bytes.NewReader(b))
	if e != nil{
		fmt.Println("post err: ",e)
		return
	}
	defer re.Body.Close()
	if re.StatusCode != http.StatusOK {
		fmt.Println("code: ",re.StatusCode)
		return
	}
	b=[]byte{}
}
func (f *File) Delete(addr string)  {
	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("json.Marshal err: ",e)
		return
	}
	re,e := http.Post("http://"+addr+"/delete/","application/octet-stream",bytes.NewReader(b))
	defer re.Body.Close()
	if e != nil{
		fmt.Println("post err: ",e)
		return
	}
	if re.StatusCode != http.StatusOK {
		fmt.Println("code: ",re.StatusCode)
		return
	}
	b=[]byte{}
}

