package main

import (
	"os"
	"fmt"
	"strings"
)
import "github.com/lepture/removie"

type service interface {
	M3U8(url string) removie.PlayList
}

func findService(url string) service {
	if strings.Contains(url, "tv.sohu.com/") {
		return &removie.Sohu{}
	} else if strings.Contains(url, "youku.com/") {
		return &removie.Youku{}
	} else if strings.Contains(url, "tudou.com/") {
		return &removie.Tudou{}
	} else {
		return nil
	}
}

func main() {
	var s service;
	url := os.Args[1]
	s = findService(url)
	if s != nil {
		data := s.M3U8(url)
		fmt.Println(data.Original)
	}
}
