package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Lyric(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "缺少必要参数 id",
		})
		return
	}

	url := fmt.Sprintf("%s%s?id=%s", APIPath, APILrc, id)
	res, err := fetchAPI(url)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "请求歌词数据失败：" + err.Error(),
		})
		return
	}

	var lyricData map[string]interface{}
	err = json.Unmarshal(res, &lyricData)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "解析歌词数据失败：" + err.Error(),
		})
		return
	}

	lrcInterface, ok := lyricData["lrc"]
	if !ok {
		sendDefaultLyric(c)
		return
	}

	lrcMap, ok := lrcInterface.(map[string]interface{})
	if !ok {
		sendDefaultLyric(c)
		return
	}

	lrcLyric, ok := lrcMap["lyric"].(string)
	if !ok || strings.TrimSpace(lrcLyric) == "" {
		sendDefaultLyric(c)
		return
	}

	tlyricLyric := ""
	tlyricInterface, tlyricExists := lyricData["tlyric"]
	if tlyricExists && tlyricInterface != nil {
		tlyricMap, ok := tlyricInterface.(map[string]interface{})
		if ok {
			tlyricLyric, ok = tlyricMap["lyric"].(string)
			if !ok {
				tlyricLyric = ""
			}
		}
	}

	if tlyricLyric != "" {
		merged := mergeLyrics(lrcLyric, tlyricLyric)
		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(200, merged)
		return
	}

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, lrcLyric)
}

func sendDefaultLyric(c *gin.Context) {
	defaultLyric := "[00:00.000] 当前音乐无歌词"
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(200, defaultLyric)
}

func mergeLyrics(lrc, tlyric string) string {
	lrcLines := strings.Split(lrc, "\n")
	tlyricLines := strings.Split(tlyric, "\n")

	merged := ""
	for i := 0; i < len(lrcLines) && i < len(tlyricLines); i++ {
		if strings.HasPrefix(lrcLines[i], "[") && strings.HasPrefix(tlyricLines[i], "[") {
			timeTag := parseTimeTag(lrcLines[i])
			origin := parseLyricLine(lrcLines[i])
			translate := parseLyricLine(tlyricLines[i])
			if timeTag != "" {
				merged += fmt.Sprintf("%s%s\n", timeTag, origin)
				merged += fmt.Sprintf("%s%s\n", timeTag, translate)
			}
		}
	}
	return merged
}

func parseTimeTag(line string) string {
	idx := strings.Index(line, "]")
	if idx > 0 && line[0] == '[' {
		return line[:idx+1]
	}
	return ""
}

func parseLyricLine(line string) string {
	idx := strings.Index(line, "]")
	if idx != -1 && idx+1 < len(line) {
		return line[idx+1:]
	}
	return ""
}
