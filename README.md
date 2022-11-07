# sync
基于文件系统实时文件同步服务，功能类似nodfiy+sync。  
main下面是服务端，直接编译出来就是服务端  
client下面的main是客户端，编译出来运行在客户端上  
客户端启动：datadir="/tmp-test1" ./sync-client  
客户端端口默认：8010   
datadir: 客户端存储的数据目录
     
服务端启动命令：datadir=/tmp clientaddr=xxxxxxxx:8010 delete="false"  ./sync-server  
datadir：指定需要同步的目录  
clientaddr： 客户端运行的IP+端口  
delete:  是否同步删除操作，true表示同步删除操作，false表示不同步删除操作  
注意：客户端和服务端指定数据目录的时候要保持一致，比如服务端/tmp 后面没有"/"的时候，客户端也尽量不要有，客户端的格式/tmp-test1，如果服务端/tmp/这样的时候客户端也要有加上"/",比如：/tmp-test1/，否则会造成目录无法识别的问题。
