package redis

import (
	"fmt"
	"testing"
)

func TestRedisCusterCli(t *testing.T) {
	reply, _ := RedisCusterCli("10.26.188.20:2900", "json.get", "ssdb-test", ".ClusterVersion")
	fmt.Println(reply)
}