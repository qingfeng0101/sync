package tools

import (
	"encoding/json"
	"fmt"
	"os"
)

func SaveData(ch chan *ChenData,path string)  {
	for {
		select  {
		case c := <- ch:
			var s = Init(path)
			fmt.Println(s.SaveData)
			info,_ := os.Stat(c.Name)
			s.Record(c.Name,info.ModTime().Unix())
			s.Save()
			s.Empty()
		}
	}
}

type ChenData struct {
	Name string
	Value int64
}

type SaveDatas struct {
	SaveData map[string]int64
	SavePath string `json:"save_path"`
}

func Init(path string)  *SaveDatas {
	var SaveDatas  SaveDatas
	SaveDatas.SaveData = make(map[string]int64)
	SaveDatas.SavePath = path
	if _, err := os.Stat(SaveDatas.SavePath); !os.IsNotExist(err) {
         f,err :=  os.Open(SaveDatas.SavePath)
         if err != nil{
         	fmt.Println("os.Open(SaveDatas.SavePath) err: ",err)
         	return nil
		 }
		 defer f.Close()
		 info,_:= f.Stat()
		 buf := make([]byte,info.Size())
		 f.Read(buf)
		 json.Unmarshal(buf,&SaveDatas.SaveData)
	}
	return &SaveDatas
}
func (s *SaveDatas) Record(name string,d int64) {
	s.SaveData[name] = d
}
func (s *SaveDatas) Exist(name string) bool {
	 _,ok := s.SaveData[name]
	 return ok
}
func (s *SaveDatas) GetDate(name string) int64 {
	d,_ := s.SaveData[name]
	return d
}
func (s *SaveDatas) ContrastDate(name string,new int64) bool {
	old := s.GetDate(name)
	if new > old{
		return false
	}
	return true
}
func (s *SaveDatas) Del(name string)  {
	delete(s.SaveData,name)
}
func (s *SaveDatas) Save() error {
	b,err := json.Marshal(s.SaveData)
	if err != nil{
		return err
	}
	f,err := os.OpenFile(s.SavePath,os.O_CREATE|os.O_TRUNC|os.O_WRONLY,0666)
	if err != nil{
		return err
	}
	defer f.Close()
	f.Write(b)
	f.Sync()
	return nil
}
func (s *SaveDatas) Empty()  {
	s.SaveData = make(map[string]int64)
}


