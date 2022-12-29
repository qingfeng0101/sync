package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	DataDIr string `json:"datadir"`
	Clientaddr string `json:"clientaddr"`
	Delete bool `json:"delete"`
	Exclude string `json:"exclude"`
	SaveFile string `json:"savefile"`
	SourceDelete bool `json:"sourcedelete"`
}
type ClientConf struct {
	Ipaddr string `json:"ipaddr"`
	Port string `json:"port"`
	DataDIr string `json:"datadir"`
}

func NewClient(path string)  *ClientConf {
	f,err := os.Open(path)
	if err != nil{
		fmt.Println("config open err: ",err)
		return nil
	}
	defer f.Close()
	b,err := ioutil.ReadAll(f)
	if err != nil{
		fmt.Println("config ioutil.ReadAll err: ",err)
		return nil
	}
	var Config ClientConf
	err = yaml.Unmarshal(b,&Config)
	if err != nil{
		fmt.Println("config json.Unmarshal err: ",err)
		return nil
	}
	return &Config
}
func NewConfing(path string) *Config {
	f,err := os.Open(path)
	if err != nil{
		fmt.Println("config open err: ",err)
		return nil
	}
	defer f.Close()
	b,err := ioutil.ReadAll(f)
	if err != nil{
		fmt.Println("config ioutil.ReadAll err: ",err)
		return nil
	}
	var Config Config
	err = yaml.Unmarshal(b,&Config)
	if err != nil{
		fmt.Println("config json.Unmarshal err: ",err)
		return nil
	}
	return &Config
}