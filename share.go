package main

import (
	"dropboxshare/route"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var port string
var doc string

func init() {
	flag.StringVar(&port, "port", "8090", "give me a port number")
	flag.StringVar(&doc, "doc", "./", "document root dir")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !filepath.IsAbs(doc) {
		doc = filepath.Join(pwd, doc)
	}
	if f, err := os.Stat(doc); err == nil {
		if !f.Mode().IsDir() {
			fmt.Println(doc + " is not directory")
			os.Exit(3)
		}
	} else {
		fmt.Println(doc + " not exists")
		os.Exit(2)
	}

}

func main() {
	flag.Parse()
	http.HandleFunc("/status", status)
	http.HandleFunc("/", routeMatch)
	fmt.Println("Starting up on port " + port)
	fmt.Println("Document root " + doc)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func status(w http.ResponseWriter, r *http.Request) {

}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	found := false
	for _, p := range route.RoutePath {
		if p.Reg.MatchString(r.URL.Path) {
			found = true
			p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path))
			break
		}
	}
	if !found {
		fallback(w, r)
	}
}

func fallback(w http.ResponseWriter, r *http.Request) {
	var realpath string = filepath.Join(doc, r.URL.Path)
	if f, err := os.Stat(realpath); err == nil {
		if f.Mode().IsRegular() {
			http.ServeFile(w, r, realpath)
			return
		}
	}
	http.NotFound(w, r)
}
