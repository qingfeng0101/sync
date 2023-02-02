package tools

import (
	"github.com/fsnotify/fsnotify"
	"os"
)
// 将需要的管道做一个聚合
type Channels struct {
	DataChan chan int
	EndChan chan int
	 Watch *fsnotify.Watcher
     ChanDatas chan *ChenData
     Sigs chan os.Signal
	 SaveStop chan os.Signal

}
// 初始化一个聚合的管道集合
func NewChannels() *Channels {
	var Channels Channels
	Channels.ChanDatas = make(chan *ChenData,100)
	Channels.Sigs = make(chan os.Signal, 1)
	Channels.SaveStop = make(chan os.Signal,1)
	Channels.DataChan = make(chan int,10)
	Channels.EndChan = make(chan int)
	Channels.Watch,_ = fsnotify.NewWatcher()
	return &Channels
}