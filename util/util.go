package util

import (
	"fmt"
	"net/http"
	"time"
)

func JsonPut(w http.ResponseWriter, bs []byte, httpCache bool, cacheTime uint32) {
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept, Accept-Encoding")
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	if httpCache {
		UseHttpCache(w, cacheTime)
	}
	w.Write(bs)
}

func UseHttpCache(w http.ResponseWriter, cacheTime uint32) {
	w.Header().Set("Expires", time.Now().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cacheTime))
}
