package server

//func Writefile(files map[string]*os.File,addr,basePath string) int {
//	// 性能优化可以改为携程方式
//	for filename,f := range files{
//		if f == nil{
//			delete(files,filename)
//			continue
//		}
//         info,err :=  f.Stat()
//         if err != nil{
//         	fmt.Println("f.Stat err: ",err)
//			 return 1
//		 }
//         if info.Size() > file.Buf{
//         	//分片操作
//			 tools.ShardData1(f,filename,addr,basePath)
//			 delete(files,filename)
//			 return 0
//		 }
//		 f.Read(file.Bufs)
//		 filestu := file.NewFile(basePath)
//		 filestu.Name = filename
//		 filestu.Date = file.Bufs[:info.Size()]
//		 filestu.Shard = 0
//		 filestu.Operation = "append"
//		 filestu.Sendfile(addr)
//		 f.Close()
//		 delete(files,filename)
//	}
//	return 0
//}