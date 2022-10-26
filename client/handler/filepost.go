package handler

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	file2 "sync/file"
)



func filepost(w http.ResponseWriter,r *http.Request)  {
	http_body,_ := ioutil.ReadAll(r.Body)
	var file file2.File
	err := json.Unmarshal(http_body,&file)
	if err != nil{
		log.Println("json.Unmarshal err: ",err)
		return
	}
	f,e := os.OpenFile(file.Name,os.O_WRONLY|os.O_CREATE, 0666)
	if e != nil{
		log.Println("file name: ",file.Name)
		log.Println("OpenFile err: ",e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	write := bufio.NewWriter(f)
	_,err = write.Write(file.Date)
	if err != nil{
		log.Println("Write err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	write.Flush()


	w.WriteHeader(http.StatusOK)
	return
}
