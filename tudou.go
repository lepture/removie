package removie

import "regexp"


type Tudou struct {}

func (s *Tudou) M3U8(uri string) PlayList {
	vid := s.parseVideoID(uri)
	if vid == "" {
		return PlayList{}
	}
	return s.parseM3u8(vid)
}

var tudouIDRegex = regexp.MustCompile(`.*vcode\:\s*'([a-zA-Z0-9]+)'`)

func (s *Tudou) parseVideoID(uri string) string {
	body, _ := request(uri)
	m := tudouIDRegex.FindStringSubmatch(string(body))
	if (m == nil) {
		return ""
	}
	return string(m[1])
}

func (s *Tudou) parseM3u8(vid string) PlayList {
	return parseYoukuM3u8(vid)
}
