package main

import (
	"net/http"
	"sync/client/handler"
)

func main() {
	http.HandleFunc("/dir/",handler.Handler)
	http.HandleFunc("/file/",handler.Handler)
	http.HandleFunc("/delete/",handler.Handler)
	http.ListenAndServe("0.0.0.0:8010",nil)
}
