package middleware

import (
	"fmt"
	"github.com/tj/go-dropbox"
	"net/http"
	"os"
)

var BoxClient *dropbox.Client

func init() {
	fmt.Println("init dropbox")

	token := os.Getenv("DROPBOX_ACCESS_TOKEN")
	BoxClient = dropbox.New(dropbox.NewConfig(token))
}

func ServeBoxFile(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello"))

}
