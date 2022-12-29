package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync/conf"
	file2 "sync/file"
	"time"
)

func SaveData(channels *Channels,s *SaveDatas,conf *conf.Config)  {
	for {
		select  {
		case c := <- channels.ChanDatas:
			info,err := os.Stat(c.Name)
			if os.IsNotExist(err){
				log.Println("IsNotExist filename：",c.Name)
				continue
			}
			if s.Exist(c.Name) && s.ContrastDate(c.Name,info.ModTime().Unix()){
				continue
			}
			ok,err := DataSize(c.Name,file2.Buf)
			if err != nil{
				log.Println("NilDir DataSize err: ",err)
				channels.EndChan <- 1
				return
			}
			if ok{
				ShardData(c.Name,conf.Clientaddr,conf.DataDIr)
				if conf.SourceDelete {
					os.Remove(c.Name)
				}
			}else {
				c.Value = info.Size()
				SmallData(c,conf)
				if conf.SourceDelete {
					os.Remove(c.Name)
				}
			}
            if s.SavePath != ""{
				s.Record(c.Name,info.ModTime().Unix())
				if len(s.SaveData) % 10 ==0 {
					s.Save()
					if len(channels.DataChan) < 100{
						channels.DataChan <- 1
					}
			}
			}
		case  <- channels.Sigs:
			channels.EndChan <- 1
			break
		}
	}
}

type ChenData struct {
	Name string
	Value int64
	Operation string
}

type SaveDatas struct {
	SaveData map[string]int64
	//SaveData sync.Map
	SavePath string `json:"save_path"`
}

func Init(path string)  *SaveDatas {
	var SaveDatas  SaveDatas
	SaveDatas.SaveData = map[string]int64{}
	SaveDatas.SavePath = path
	if path != ""{
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
func CronData(c chan int, s *SaveDatas)  {
	lendata :=  len(s.SaveData)
	for  {
		time.Sleep(time.Second * 10)
		if len(c) == 0 && len(s.SaveData) != lendata {
			s.Save()
			lendata = len(s.SaveData)
		}else {
			for n := 0;n <len(c);n++{
				go func() {
					<- c
					return
				}()
			}

		}
	}

}

