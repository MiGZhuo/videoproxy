package middleware

import (
	"dropboxshare/streampipe"
	"fmt"
	"net/http"
	"net/url"
	_ "time"
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

var youtube_video_host string = "http://www.youtube.com/get_video_info?video_id="

var youtube_video_map = map[string]string{
	"large":  "high",
	"medium": "medium",
	"small":  "small",
}

type videoInfo struct {
	Id       string
	Title    string
	Duration string
	Keywords string
	Author   string
	Stream   []map[string]string
}

func init() {
	fmt.Println("init youtu")
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

func GetYoutubeVideoUrl(match []string) string {
	var url string = fmt.Sprintf("%s%s%s", youtube_video_host, match[2], "&asv=3&el=detailpage&hl=en_US")
	fmt.Println(url)
	info, err := GetYoutubeVideoMeta(url)
	fmt.Println(info, err)

	return url

}

func GetYoutubeVideoMeta(u string) (videoInfo, error) {
	bytes, err := streampipe.Get(u)
	var info videoInfo
	if err != nil {
		return info, err
	}
	values, err := url.ParseQuery(string(bytes))
	fmt.Println(values)
	if err != nil {
		return info, err
	}

	return info, nil
}
