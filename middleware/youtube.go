package middleware

import (
	"dropboxshare/streampipe"
	"fmt"
	"github.com/suconghou/youtubeVideoParser"
	"net/http"
)

var youtube_image_map = map[string]string{
	"large":  "hqdefault",
	"medium": "mqdefault",
	"small":  "default",
}

var youtube_image_host_map = map[string]string{
	"jpg":  "http://i.ytimg.com/vi/",
	"webp": "http://i.ytimg.com/vi_webp/",
}

func ServeYoutubeImage(w http.ResponseWriter, r *http.Request, url string) {
	streampipe.Pipe(w, r, url)
}

func GetYoutubeImageUrl(match []string) string {
	var url string = fmt.Sprintf("%s%s/%s.%s", youtube_image_host_map[match[3]], match[2], youtube_image_map[match[1]], match[3])
	return url
}

func ServeYoutubeVideo(w http.ResponseWriter, r *http.Request, url string) {
	streampipe.Pipe(w, r, url)
}

func GetYoutubeVideoUrl(match []string) (string, []byte, error) {
	var videoId string = match[2]
	var quality string = match[1]
	var types string = match[3]
	info, err := youtubeVideoParser.Parse(videoId)
	if err != nil {
		return "", []byte(""), err
	}
	if types == "json" {
		bs, err := info.ToJson()
		if err != nil {
			return "", []byte(""), err
		}
		return "", bs, nil
	} else {
		url, _, err := info.GetStream(quality, types)
		if err != nil {
			return "", []byte(""), err
		}
		return url, []byte(""), nil
	}

}
