package handler

import (
	"encoding/json"
	"net/http"
	"os"
	file2 "sync/file"
)


func post(w http.ResponseWriter,r *http.Request)  {
    size := r.ContentLength
    buf := make([]byte,size)
    var file file2.File

	r.Body.Read(buf)
	json.Unmarshal(buf,&file)

    os.MkdirAll(file.Name,os.ModePerm)
    w.WriteHeader(http.StatusOK)
}
