package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type File struct {
	Name string `json:"name"`
	Date []byte `json:"date"`
}

func (f *File) Sendfile()  {

	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("encoder.Encode err: ",e)
		return
	}
	fmt.Println("file Sendfile file name: ",f.Name)
	re,e := http.Post("http://"+os.Getenv("clientaddr")+"/file/","application/json;utf-8",bytes.NewReader(b))
	defer re.Body.Close()
	if e != nil{
		fmt.Println("post err: ",e)
		return
	}
	if re.StatusCode != http.StatusOK {
		fmt.Println("code: ",re.StatusCode)
		return
	}
}
func (f *File) Senddir()  {
	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("json.Marshal err: ",e)
		return
	}
	re,e := http.Post("http://"+os.Getenv("clientaddr")+"/dir/","application/octet-stream",bytes.NewReader(b))
	//defer re.Body.Close()
	if e != nil{
		fmt.Println("post err: ",e)
		return
	}
	if re.StatusCode != http.StatusOK {
		fmt.Println("code: ",re.StatusCode)
		return
	}
}
func (f *File) Delete()  {
	b,e := json.Marshal(f)
	if e != nil{
		fmt.Println("json.Marshal err: ",e)
		return
	}
	re,e := http.Post("http://"+os.Getenv("clientaddr")+"/delete/","application/octet-stream",bytes.NewReader(b))
	//defer re.Body.Close()
	if e != nil{
		fmt.Println("post err: ",e)
		return
	}
	if re.StatusCode != http.StatusOK {
		fmt.Println("code: ",re.StatusCode)
		return
	}
}

