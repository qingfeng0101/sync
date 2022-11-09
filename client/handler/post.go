package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	file2 "sync/file"
	"sync/server/tools"
)


func post(w http.ResponseWriter,r *http.Request)  {
    size := r.ContentLength
    buf := make([]byte,size)
	var Path = tools.RewritePath(Client.DataDIr)
    var file file2.File
	r.Body.Read(buf)
	json.Unmarshal(buf,&file)
	file.Name = strings.Replace(file.Name,file.Path,Path,-1)
    os.MkdirAll(file.Name,os.ModePerm)
    w.WriteHeader(http.StatusOK)
}
