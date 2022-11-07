# sync
基于文件系统实时文件同步服务，功能类似nodfiy+sync。  
main下面是服务端，直接编译出来就是服务端  
client下面的main是客户端，编译出来运行在客户端上  
客户端启动：./sync-client  
客户端端口默认：8010  
服务端启动命令：datadir=/tmp clientaddr=xxxxxxxx:8010 delete="false"  ./sync-server  
datadir：指定需要同步的目录  
clientaddr： 客户端运行的IP+端口  
delete:  是否同步删除操作，true表示同步删除操作，false表示不同步删除操作  
