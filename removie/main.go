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
	fmt.Println("Playing: %s", uri)
	cmd := exec.Command("open", "-a", "QuickTime Player", uri)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}


func main() {
	var s service;
	uri := os.Args[1]
	s = findService(uri)
	if s != nil {
		data := s.M3U8(uri)
		macPlay(data.Original)
	}
}
