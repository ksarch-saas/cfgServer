package redis

import (
	"fmt"
	"time"
	"strings"

	"github.com/mediocregopher/radix.v2/cluster"
	"github.com/mediocregopher/radix.v2/redis"
)

const (
	NUM_RETRY     = 3
	POOL_SIZE     = 1
	CONN_TIMEOUT  = 5 * time.Second
	WRITE_TIMEOUT = 120 * time.Second
)

/*
 * DialCluster and RedisCusterCli is used to metadb
 */
func DialCluster(addr string) (*cluster.Cluster, error) {
	options := cluster.Opts{
		Addr:     addr,
		Timeout:  WRITE_TIMEOUT,
		PoolSize: POOL_SIZE,
	}

	client, err := cluster.NewWithOpts(options)
	if err == nil {
		return client, nil
	}

	for retry := 0; retry <  NUM_RETRY; retry++ {
		client, err = cluster.NewWithOpts(options)
		if err == nil {
			break
		}
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

/*
 * above funcs is used to ssdb
 */

func dail(address string) (*redis.Client, error){
	client, err := redis.DialTimeout("tcp", address, CONN_TIMEOUT)
	if err == nil {
		return client, err
	}

	for retry := 0; retry <  NUM_RETRY; retry++ {
		client, err = redis.DialTimeout("tcp", address, CONN_TIMEOUT)
		if err == nil {
			break
		}
	}

	return client, err
}

func RedisCli(addr string, cmd string, args ...interface{}) (interface{}, error) {
	client, err := dail(addr)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	reply := client.Cmd(cmd, args ...)
	if reply.Err != nil {
		return nil, reply.Err
	}
	return reply, nil
}


func Info(addr string) (map[string]string, error){
	reply, err := RedisCli(addr, "info")
	if err != nil {
		return nil, err
	}

	reply_array, err := reply.(*redis.Resp).List()
	if err != nil {
		return nil, err
	}

	var seg []string
	for _, ele := range reply_array {
		seg = append(seg, strings.Fields(ele)...)
	}

	info := make(map[string]string)
	for idx, rep := range seg {
		if strings.Contains(rep, "dbsize") {
			info["dbsize"] = rep[strings.Index(rep, ":")+1 : len(rep)]
		} else if strings.Contains(rep, "master_host"){
			info["master_host"] = rep[strings.Index(rep, ":")+1 : len(rep)]
		} else if strings.Contains(rep, "master_port"){
			info["master_port"] = rep[strings.Index(rep, ":")+1 : len(rep)]
		} else if strings.Contains(rep, "last_seq"){
			info["last_seq"] = rep[strings.Index(rep, ":")+1 : len(rep)]
		} else if strings.Contains(rep, "readonly"){
			info["readonly"] = rep[strings.Index(rep, ":")+1 : len(rep)]
		} else if strings.Contains(rep, "slot"){
			slot_num := 0
			idx = idx + 1
			for {
				slot := seg[idx]
				if !strings.Contains(slot, "[") || 
				!strings.Contains(slot, "-") || 
				!strings.Contains(slot, "]") {
					break
				}

				start := fmt.Sprintf("slot_start_%d", slot_num)
				end   := fmt.Sprintf("slot_end_%d", slot_num)
				info[start] = slot[strings.Index(slot, "[")+1 : strings.Index(slot, "-")]
				info[end] = slot[strings.Index(slot, "-")+1 : strings.Index(slot, "]")]
				slot_num = slot_num + 1
				idx = idx +1

			}
			info["slot_num"] = fmt.Sprintf("%d", slot_num)
		}
	}

	return info, nil
}






