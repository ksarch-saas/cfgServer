package redis

import (
	"time"
	"github.com/mediocregopher/radix.v2/cluster"
)

const (
	NUM_RETRY     = 3
	POOL_SIZE     = 1
	CONN_TIMEOUT  = 5 * time.Second
	WRITE_TIMEOUT = 120 * time.Second
)

func DialCluster(addr string) (*cluster.Cluster, error) {
	options := cluster.Opts{
		Addr:     addr,
		Timeout:  WRITE_TIMEOUT,
		PoolSize: POOL_SIZE,
	}

	client, err := cluster.NewWithOpts(options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func RedisCusterCli(addr string, cmd string, args ...interface{}) (interface{}, error) {
	conn, err := DialCluster(addr)
	if err != nil {
		return "redis: connection error", err
	}
	defer conn.Close()
	reply := conn.Cmd(cmd, args...)

	return reply, reply.Err
}