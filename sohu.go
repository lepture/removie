package removie

import (
	"os"
	"fmt"
	"regexp"
	"encoding/json"
)


type Sohu struct {}

func (s *Sohu) M3U8(url string) PlayList {
	vid := parseVideoID(url)
	if vid == "" {
		return PlayList{}
	}
	return parseM3u8(vid)
}


var apiKey = os.Getenv("REMOVIE_SOHU_KEY")
var apiUrl string = "http://api.tv.sohu.com/v4/video/info/"
var vidRegex = regexp.MustCompile(`.*var\s+vid\s*=\s*"(\d+)"`)

type message struct {
	Status int
	Data struct {
		Url_original string
		Url_super string
		Url_high string
		Url_nor string
	}
}


func parseVideoID(url string) string {
	body, _ := request(url)
	m := vidRegex.FindStringSubmatch(string(body))
	if (m == nil) {
		return ""
	}
	return string(m[1])
}


func parseM3u8(vid string) PlayList {
	if apiKey == "" {
		apiKey = "f351515304020cad28c92f70f002261c"
	}
	url := fmt.Sprintf("%s%s.json?api_key=%s", apiUrl, vid, apiKey)
	body, err := request(url)
	if (err != nil) {
		return PlayList{}
	}

	var m message
	jerr := json.Unmarshal(body, &m)
	if (jerr != nil) {
		return PlayList{}
	}

	if m.Status != 200 {
		return PlayList{}
	}

	return PlayList{
		Original: m.Data.Url_original,
		High: m.Data.Url_super,
		Normal: m.Data.Url_high,
		Low: m.Data.Url_nor,
	}
}
