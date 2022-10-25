package handler

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	file2 "sync/file"
)



func filepost(w http.ResponseWriter,r *http.Request)  {
	size := r.ContentLength
	buf := make([]byte,size)
	r.Body.Read(buf)
	var file file2.File
	json.Unmarshal(buf,&file)
	f,e := os.OpenFile(file.Name,os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil{
		log.Println("OpenFile err: ",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	_,err := write.Write(file.Date)
	if err != nil{
		log.Println("Write err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	write.Flush()


	w.WriteHeader(http.StatusOK)
	return
}
