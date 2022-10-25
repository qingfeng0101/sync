package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	file2 "sync/file"
)

func delete(w http.ResponseWriter,r *http.Request)  {
	size := r.ContentLength
	buf := make([]byte,size)
	r.Body.Read(buf)
	var file file2.File
	json.Unmarshal(buf,&file)
	err := os.Remove(file.Name)
	if err != nil{
		log.Println("os.Remove err: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
