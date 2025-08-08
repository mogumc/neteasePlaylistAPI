package api

import (
	"io"
	"net/http"
)

const (
	APIPath     = "https://api.sayqz.com/tunefree/ncmapi/"
	APIPlaylist = "playlist/track/all"
	APIMusic    = "song/url/v1"
	APILrc      = "lyric"
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

func fetchAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
