package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Netease(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "缺少必要参数 id",
		})
		return
	}

	limit := c.DefaultQuery("limit", "50")
	lim, _ := strconv.Atoi(limit)
	if lim > 100 || lim <= 0 {
		lim = 100
	}

	url := fmt.Sprintf("%s%s?id=%s&limit=%d", APIPath, APIPlaylist, id, lim)
	res, err := fetchAPI(url)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "请求歌单数据失败：" + err.Error(),
		})
		return
	}

	var songList map[string]interface{}
	if err := json.Unmarshal(res, &songList); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "解析歌单JSON失败：" + err.Error(),
		})
		return
	}

	rawSongs, ok := songList["songs"].([]interface{})
	if !ok {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "未找到歌曲列表数据",
		})
		return
	}

	var formatted []Song
	var ids []string

	for _, raw := range rawSongs {
		song, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		songID, ok := song["id"].(float64)
		if !ok {
			continue
		}
		strID := strconv.FormatFloat(songID, 'f', 0, 64)

		artistName := "未知歌手"
		if arList, ok := song["ar"].([]interface{}); ok && len(arList) > 0 {
			if ar, ok := arList[0].(map[string]interface{}); ok {
				if name, ok := ar["name"].(string); ok {
					artistName = name
				}
			}
		}

		albumName := "未知歌曲"
		if al, ok := song["al"].(map[string]interface{}); ok {
			if name, ok := al["name"].(string); ok {
				albumName = name
			}
		}

		formatted = append(formatted, Song{
			MusicID:     strID,
			MusicTitle:  albumName,
			MusicAuthor: artistName,
		})
		ids = append(ids, strID)
	}

	if len(ids) == 0 {
		c.JSON(200, gin.H{
			"code": 200,
			"data": formatted,
		})
		return
	}
	fmt.Println("Fetching music URLs for IDs:", ids)
	var level string
	level_id := c.Query("level")
	switch level_id {
	case "1":
		level = "higher"
	case "2":
		level = "lossless"
	case "3":
		level = "exhigh"
	case "4":
		level = "hires"
	default:
		level = "standard"
	}
	musicURL := fmt.Sprintf("%s%s?id=%s&level=%s", APIPath, APIMusic, strings.Join(ids, ","), level)
	res2, err := fetchAPI(musicURL)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "获取歌曲链接失败：" + err.Error(),
		})
		return
	}

	var musicRes struct {
		Data []struct {
			ID  float64 `json:"id"`
			URL string  `json:"url,omitempty"`
			MD5 string  `json:"md5,omitempty"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res2, &musicRes); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "解析音乐链接JSON失败：" + err.Error(),
		})
		return
	}

	urlMap := make(map[string]struct {
		URL string
		MD5 string
	})

	for _, item := range musicRes.Data {
		key := strconv.Itoa(int(item.ID))
		urlMap[key] = struct {
			URL string
			MD5 string
		}{
			URL: item.URL,
			MD5: item.MD5,
		}
	}

	for i := range formatted {
		if info, ok := urlMap[formatted[i].MusicID]; ok {
			formatted[i].URL = info.URL
			formatted[i].MD5 = info.MD5
			formatted[i].Lrc = fmt.Sprintf("https://dev.moguq.top/lyric?id=%s", formatted[i].MusicID)
		}
	}

	c.JSON(200, gin.H{
		"code": 200,
		"data": formatted,
	})
}
