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

	tlyricMap := make(map[string]string)
	for _, line := range tlyricLines {
		tag := parseTimeTag(line)
		text := parseLyricLine(line)
		if tag != "" && text != "" {
			tlyricMap[tag] = text
		}
	}

	mergedTags := make(map[string]bool)
	var builder strings.Builder

	for _, line := range lrcLines {
		tag := parseTimeTag(line)
		if tag == "" {
			continue
		}
		content := parseLyricLine(line)
		if content == "" {
			continue
		}

		builder.WriteString(fmt.Sprintf("%s %s\n", tag, content))

		if trans, ok := tlyricMap[tag]; ok && trans != "" && !mergedTags[tag] {
			builder.WriteString(fmt.Sprintf("%s %s\n", tag, trans))
			mergedTags[tag] = true
		}
	}

	return builder.String()
}

func parseTimeTag(line string) string {
	start := strings.Index(line, "[")
	end := strings.Index(line, "]")
	if start == -1 || end == -1 || end <= start {
		return ""
	}

	raw := line[start+1 : end]

	if strings.Count(raw, ":") >= 2 {
		parts := strings.Split(raw, ":")
		if len(parts) < 3 {
			return ""
		}
		mm := parts[0]
		ss := parts[1]
		ff := parts[2]

		mm = fmt.Sprintf("%02s", mm)
		ss = fmt.Sprintf("%02s", ss)
		ff = fmt.Sprintf("%02s", ff)

		return fmt.Sprintf("[%s:%s.%s]", mm, ss, ff)
	}

	parts := strings.Split(raw, ":")
	if len(parts) < 2 {
		return ""
	}

	mm := parts[0]
	ssAndRest := strings.Join(parts[1:], ":")

	secParts := strings.Split(ssAndRest, ".")
	ss := secParts[0]
	ff := ""
	if len(secParts) > 1 {
		ff = secParts[1]
	}

	ff = fmt.Sprintf("%-2s", ff)
	ff = strings.ReplaceAll(ff, " ", "0")
	if len(ff) > 2 {
		ff = ff[:2]
	}

	mm = fmt.Sprintf("%02s", mm)
	ss = fmt.Sprintf("%02s", ss)

	return fmt.Sprintf("[%s:%s.%s]", mm, ss, ff)
}

func parseLyricLine(line string) string {
	idx := strings.Index(line, "]")
	if idx != -1 && idx+1 < len(line) {
		return line[idx+1:]
	}
	return ""
}
