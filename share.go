package main

import (
	"dropboxshare/route"
	"dropboxshare/util"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	startTime = time.Now()
	port      string
	doc       string
)

var sysStatus struct {
	Uptime       string
	GoVersion    string
	Hostname     string
	MemAllocated uint64
	MemTotal     uint64
	MemSys       uint64
	NumGoroutine int
	CpuNum       int
	Pid          int
}

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
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	log.Fatal(http.ListenAndServe(bind, nil))
}

func status(w http.ResponseWriter, r *http.Request) {
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	sysStatus.Uptime = time.Since(startTime).String()
	sysStatus.NumGoroutine = runtime.NumGoroutine()
	sysStatus.MemAllocated = memStat.Alloc
	sysStatus.MemTotal = memStat.TotalAlloc
	sysStatus.MemSys = memStat.Sys
	sysStatus.CpuNum = runtime.NumCPU()
	sysStatus.GoVersion = runtime.Version()
	sysStatus.Hostname, _ = os.Hostname()
	sysStatus.Pid = os.Getpid()
	if bs, err := json.Marshal(&sysStatus); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
	} else {
		util.JsonPut(w, bs)
	}
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
