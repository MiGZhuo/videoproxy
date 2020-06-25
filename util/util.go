package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	// Log print to stdout
	Log = log.New(os.Stdout, "", 0)
)

// JSONPut resp json
func JSONPut(w http.ResponseWriter, v interface{}, status int, age int) (int, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	h := w.Header()
	h.Set("Content-Type", "text/json; charset=utf-8")
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Cache-Control", fmt.Sprintf("public,max-age=%d", age))
	w.WriteHeader(status)
	return w.Write(bs)
}
