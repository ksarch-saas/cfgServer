package utils

import (
	"fmt"
	"testing"
)

func TestGetAddrFromBNS(t *testing.T) {
	resp, err := GetAddrFromBNS("group.redis3-video-userbehaviordata.osp.cn")
	fmt.Println(resp, err)
}

func TestLocalIP(t *testing.T) {
	resp, err := LocalIP()
	fmt.Println(resp, err)
}

func TestHostname(t *testing.T) {
	resp, err := Hostname()
	fmt.Println(resp, err)
}