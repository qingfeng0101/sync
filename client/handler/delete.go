package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	file2 "sync/file"
	"sync/server/tools"
)

func del(w http.ResponseWriter,r *http.Request)  {
	size := r.ContentLength
	buf := make([]byte,size)
	var Path = tools.RewritePath(Client.DataDIr)
	r.Body.Read(buf)
	var file file2.File
	json.Unmarshal(buf,&file)
	file.Name = strings.Replace(file.Name,file.Path,Path,-1)
	if file.Systype == "windows"{
		file.Name = strings.Replace(file.Name,"\\","/",-1)
	}
	err := os.Remove(file.Name)
	if err != nil{
		log.Println("os.Remove err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
