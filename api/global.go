package api

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	APIPath     = "https://api.sayqz.com/tunefree/ncmapi/"
	APIPlaylist = "playlist/track/all"
	APIMusic    = "song/url/v1"
	APILrc      = "lyric"
	DevPath     = "https://dev.moguq.top/"
)

type Song struct {
	MusicID     string `json:"id"`
	MusicTitle  string `json:"title"`
	MusicAuthor string `json:"author"`
	MusicAlbum  string `json:"album"`
	MusicCover  string `json:"cover,omitempty"`
	URL         string `json:"url,omitempty"`
	MD5         string `json:"md5,omitempty"`
	Lrc         string `json:"lrc,omitempty"`
}

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func fetchAPI(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	body := new(bytes.Buffer)
	if err != nil {
		log.Println(err)
	} else {
		defer resp.Body.Close()
		io.Copy(body, resp.Body)
	}
	return body.Bytes(), nil
}
