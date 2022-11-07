package handler

import (
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter,r *http.Request)  {
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if strings.Split(r.URL.EscapedPath(),"/")[1] == "dir"{
		post(w,r)
		return
	}
	if strings.Split(r.URL.EscapedPath(),"/")[1] == "file"{
		filepost(w,r)
		return
	}
	if strings.Split(r.URL.EscapedPath(),"/")[1] == "delete"{
		del(w,r)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}
