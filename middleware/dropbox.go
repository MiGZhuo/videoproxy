package middleware

import (
	"dropboxshare/util"
	"encoding/json"
	_ "fmt"
	"github.com/tj/go-dropbox"
	"github.com/tj/go-dropy"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var publicDir string = "/Public"

var client *dropy.Client

type fileInfo struct {
	Name    string
	Size    int64
	IsDir   bool
	Path    string
	ModTime time.Time
}

type fileInfoList struct {
	List  []fileInfo
	Total int
}

func init() {
	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	if token == "" {
		token = ""
	}
	client = dropy.New(dropbox.New(dropbox.NewConfig(token)))
}

func ServeBoxFile(w http.ResponseWriter, r *http.Request, match []string) error {
	var filePath string = path.Join(publicDir, match[0])
	file, err := client.Stat(filePath)
	if err != nil {
		return err
	}
	if file.IsDir() {
		fileList, err := client.List(filePath)
		if err != nil {
			return err
		}
		infoList := make([]fileInfo, 0)
		for _, item := range fileList {
			info := fileInfo{item.Name(), item.Size(), item.IsDir(), path.Join(match[0], item.Name()), item.ModTime()}
			infoList = append(infoList, info)
		}
		res := fileInfoList{Total: len(fileList), List: infoList}
		if bs, err := json.Marshal(&res); err != nil {
			return err
		} else {
			util.JsonPut(w, bs, true, 3600)
		}
	} else {
		util.CrossShare(w)
		util.UseHttpCache(w, 86400)
		w.Header().Set("Content-Length", strconv.Itoa(int(file.Size())))
		_, err := io.Copy(w, client.Open(filePath))
		if err != nil {
			return err
		}
	}
	return nil
}
