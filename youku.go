package removie

import (
	"fmt"
	"time"
	"regexp"
	"strings"
	"net/url"
	"encoding/json"
	"encoding/base64"
)

type Youku struct {}

func (s *Youku) M3U8(uri string) PlayList {
	vid := s.parseVideoID(uri)
	if vid == "" {
		return PlayList{}
	}
	return s.parseM3u8(vid)
}

var youkuIDRegex = regexp.MustCompile(`^http://v.youku.com/v_show/id_([a-zA-Z0-9]+)`)

func (s *Youku) parseVideoID(uri string) string {
	m := youkuIDRegex.FindStringSubmatch(uri)
	if (m == nil) {
		return ""
	}
	return string(m[1])
}


type youkuMessage struct {
	Data [1]struct {
		Videoid string
		Ep string
		Ip int64
	}
}

func (s *Youku) parseM3u8(vid string) PlayList {
	return parseYoukuM3u8(vid)
}

func parseYoukuM3u8(vid string) PlayList {
	uri := fmt.Sprintf("http://v.youku.com/player/getPlayList/VideoIDS/%s/ctype/12/ev/1", vid)
	body, err := request(uri)
	if (err != nil) {
		return PlayList{}
	}
	var m youkuMessage
	jerr := json.Unmarshal(body, &m)
	if (jerr != nil) {
		return PlayList{}
	}
	data := m.Data[0]
	ep2 := youkuEP2(data.Videoid, data.Ep)
	ts := int(time.Now().Unix())

	u := &url.URL{
		Scheme: "http",
		Host: "pl.youku.com",
		Path: "/playlist/m3u8",
	}
	q := u.Query()
	q.Set("vid", data.Videoid)
	q.Set("ts", fmt.Sprintf("%d", ts))
	q.Set("ep", ep2.Ep)
	q.Set("oip", fmt.Sprintf("%d", data.Ip))
	q.Set("token", ep2.Token)
	q.Set("sid", ep2.Sid)
	q.Set("keyframe", "1")
	q.Set("ctype", "12")
	q.Set("ev", "1")
	u.RawQuery = q.Encode()

	uri = fmt.Sprintf("http://pl.youku.com/playlist/m3u8?vid=%s&ts=%d&keyframe=1&ep=%s&oip=%d&ctype=12&ev=1&token=%s&sid=%s&type=", data.Videoid, ts, ep2.Ep, data.Ip, ep2.Token, ep2.Sid)
	return PlayList{
		Original: fmt.Sprintf("%s&type=%s", u, "hd2"),
		High: fmt.Sprintf("%s&type=%s", u, "hd2"),
		Normal: fmt.Sprintf("%s&type=%s", u, "mp4"),
		Low: fmt.Sprintf("%s&type=%s", u, "flv"),
	}
}

func youkuTransE(a, c string) []byte {
	b := make([]int, 256)
	for i, _ := range b {
		b[i] = i
	}

	f := 0
	h := 0

	for h < 256 {
		f = (f + b[h] + int(a[h % len(a)])) % 256
		b[h], b[f] = b[f], b[h]
		h += 1
	}

	f = 0
	h = 0
	q := 0

	l := len(c)
	var value = make([]byte, l)

	for q < l {
		h = (h + 1) % 256
		f = (f + b[h]) % 256
		b[h], b[f] = b[f], b[h]
		num := int(c[q]) ^ b[(b[h] + b[f]) % 256] 
		value[q] = byte(num)
		q += 1
	}
	return value
}


type youkuToken struct {
	Ep string
	Token string
	Sid string
}

func youkuEP2(vid, ep string) youkuToken {
	data, err := base64.StdEncoding.DecodeString(ep)
	if (err != nil) {
		return youkuToken{}
	}
	ecode := youkuTransE("becaf9be", string(data))

	bits := strings.SplitN(string(ecode), "_", 2)
	sid := bits[0]
	token := bits[1]

	ecode = youkuTransE("bf7e5f01", fmt.Sprintf("%s_%s_%s", sid, vid, token))
	ep = base64.StdEncoding.EncodeToString(ecode)

	return youkuToken{
		Ep: ep,
		Token: token,
		Sid: sid,
	}
}

func Demo(vid, ep string) youkuToken {
	return youkuEP2(vid, ep)
}
