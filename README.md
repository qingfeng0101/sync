# sync
基于文件系统实时文件同步服务，功能类似nodfiy+sync。  
main下面是服务端，直接编译出来就是服务端  
client下面的main是客户端，编译出来运行在客户端上  
客户端启动： ./sync-client  -f 配置文件路径  
默认在client/main下有配置文件的yaml示例。  
  
服务端启动命令：  ./sync-server  -f 配置文件路径   
默认在server/main下有配置文件的yaml示例。  
