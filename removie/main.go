package main

import (
	"os"
	"fmt"
	"strings"
	"os/exec"
)
import "github.com/lepture/removie"

type service interface {
	M3U8(uri string) removie.PlayList
}

func findService(uri string) service {
	if strings.Contains(uri, "tv.sohu.com/") {
		return &removie.Sohu{}
	} else if strings.Contains(uri, "youku.com/") {
		return &removie.Youku{}
	} else if strings.Contains(uri, "tudou.com/") {
		return &removie.Tudou{}
	} else {
		return nil
	}
}


func macPlay(uri string) {
	cmd := exec.Command("open", "-a", "QuickTime Player", uri)
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
}


func main() {
	var s service;
	uri := os.Args[1]
	s = findService(uri)
	if s == nil {
		return
	}
	data := s.M3U8(uri)
	m3uURL := data.High
	fmt.Println("Playing: ", m3uURL)
	if m3uURL != "" {
		macPlay(m3uURL)
	}
}
