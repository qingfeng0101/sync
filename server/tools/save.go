package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"sync/conf"
	file2 "sync/file"
	"time"
)
var ok bool
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
            if info == nil{
            	continue
			}
			if info.Size() > file2.Buf {
				ok = ShardData(c.Name,conf)
				if conf.SourceDelete {
					os.Remove(c.Name)
				}

			}else {
				c.Value = info.Size()
				ok = SmallData(c,conf)
				if conf.SourceDelete {
					os.Remove(c.Name)
				}
			}
			if s.SavePath != "" && ok {
				s.Record(c.Name,info.ModTime().Unix())
				if len(*s.SaveData) % 10 ==0 {
					s.Save()
					if len(channels.DataChan) < 10{
						channels.DataChan <- 1
					}
				}

			}


		case  <- channels.SaveStop:
            channels.EndChan <- 1
			return
		}
	}
}

type ChenData struct {
	Name string
	Value int64
	Operation string
}

type SaveDatas struct {
	SaveData *map[string]int64
	Mutex *sync.Mutex
	//SaveData sync.Map
	SavePath string `json:"save_path"`
}

func Init(path string)  *SaveDatas {
	var SaveDatas  SaveDatas
	SaveDatas.SaveData = &map[string]int64{}
	SaveDatas.SavePath = path
	SaveDatas.Mutex = &sync.Mutex{}
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
	s.Mutex.Lock()
	(*s.SaveData)[name] = d
    s.Mutex.Unlock()

}
func (s *SaveDatas) Exist(name string) bool {
	s.Mutex.Lock()
	 _,ok := (*s.SaveData)[name]
	s.Mutex.Unlock()
	 return ok
}
func (s *SaveDatas) GetDate(name string) int64 {

	d,_ := (*s.SaveData)[name]

	return d
}
func (s *SaveDatas) ContrastDate(name string,new int64) bool {
	s.Mutex.Lock()
	old := s.GetDate(name)
	s.Mutex.Unlock()
	if new > old{
		return false
	}
	return true
}
func (s *SaveDatas) Del(name string)  {
	s.Mutex.Lock()
	delete((*s.SaveData),name)
	s.Mutex.Unlock()
}
func (s *SaveDatas) Load() (*map[string]int64,error) {
	f,err := os.Open(s.SavePath)
	if err != nil{
		return nil,err
	}
	b,err := ioutil.ReadAll(f)
	if err != nil{
		return nil,err
	}
	temp := make(map[string]int64)
	json.Unmarshal(b,&temp)
	return &temp,nil
}
func (s *SaveDatas) Merge(temp,old *map[string]int64)  {
	for tname,v:= range *temp{
		for name,_:= range *old{
            if tname == name{
				(*old)[name] = v
			}
		}
	}

}
func (s *SaveDatas) Save() error {
	s.Mutex.Lock()
	b,err := json.Marshal(*s.SaveData)
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
	s.Mutex.Unlock()
	return nil
}
func (s *SaveDatas) Empty()  {
	s.SaveData = &map[string]int64{}

}
func CronData(c chan int, s *SaveDatas)  {
	lendata :=  len(*s.SaveData)
	for  {
		time.Sleep(time.Second * 10)
		if len(c) == 0 && len(*s.SaveData) != lendata {
			s.Save()
			lendata = len(*s.SaveData)
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

