package api

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Single(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "缺少必要参数 id",
		})
		return
	}
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
	musicURL := fmt.Sprintf("%s%s?id=%s&level=%s", APIPath, APIMusic, id, level)
	res, err := fetchAPI(musicURL)
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

	if err := json.Unmarshal(res, &musicRes); err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "解析音乐链接JSON失败：" + err.Error(),
		})
		return
	}

	json := c.Query("json")
	if json == "1" {
		c.JSON(200, gin.H{
			"code": 200,
			"data": musicRes.Data,
		})
		return
	} else {
		c.Redirect(302, musicRes.Data[0].URL)
	}
}
