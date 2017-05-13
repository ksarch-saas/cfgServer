package utils

import (
	"fmt"
	"testing"
)

func TestGetAddrFromBNS(t *testing.T) {
	resp, err := GetAddrFromBNS("group.redis3-video-userbehaviordata.osp.cn")
	fmt.Println(resp, err)
}