package util

import (
	"net/http"
)

func JsonPut(w http.ResponseWriter, bs []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Write(bs)
}
