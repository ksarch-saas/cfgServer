package redis

import (
	"fmt"
	"testing"
)

func TestRedisCusterCli(t *testing.T) {
	reply, _ := RedisCusterCli("10.26.188.20:2900", "json.get", "ssdb-test", ".ClusterVersion")
	fmt.Println(reply)
}

func TestInfo(t *testing.T) {
	// info, err := Info("10.50.90.21:3900")
	// fmt.Println(info, err)
	// info, err = Info("10.36.17.38:2500")
	// fmt.Println(info, err)
	info, err := Info("10.36.17.38:2500")
	fmt.Println(info, err)
}