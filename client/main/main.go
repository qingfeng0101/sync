package main

import (
	"fmt"
	"net/http"
	"sync/client/handler"
)

func main() {
	handler.Init()
	fmt.Println(handler.Client.Ipaddr+handler.Client.Port)
	http.HandleFunc("/dir/",handler.Handler)
	http.HandleFunc("/file/",handler.Handler)
	http.HandleFunc("/delete/",handler.Handler)
	http.ListenAndServe(handler.Client.Ipaddr+":"+handler.Client.Port,nil)
}
