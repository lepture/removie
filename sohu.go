package removie

import (
	"os"
	"fmt"
	"regexp"
	"encoding/json"
)


type Sohu struct {}

func (s *Sohu) M3U8(uri string) PlayList {
	vid := s.parseVideoID(uri)
	if vid == "" {
		return PlayList{}
	}
	return s.parseM3u8(vid)
}

var sohuIDRegex = regexp.MustCompile(`.*var\s+vid\s*=\s*"(\d+)"`)
// Find vid in the page
func (s *Sohu) parseVideoID(uri string) string {
	body, _ := request(uri)
	m := sohuIDRegex.FindStringSubmatch(string(body))
	if (m == nil) {
		return ""
	}
	return string(m[1])
}

// SOHU API key. You can set an API key with environment variable: `REMOVIE_SOHU_KEY`
var sohuApiKey = os.Getenv("REMOVIE_SOHU_KEY")
var sohuApiUrl = "http://api.tv.sohu.com/v4/video/info/"

// Struct for JSON unmarshal
type sohuMessage struct {
	Status int
	Data struct {
		Url_original string
		Url_super string
		Url_high string
		Url_nor string
	}
}

// Parse m3u8 play list
func (s *Sohu) parseM3u8(vid string) PlayList {
	if sohuApiKey == "" {
		sohuApiKey = "f351515304020cad28c92f70f002261c"
	}
	uri := fmt.Sprintf("%s%s.json?api_key=%s", sohuApiUrl, vid, sohuApiKey)
	body, err := request(uri)
	if (err != nil) {
		return PlayList{}
	}

	var m sohuMessage
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
