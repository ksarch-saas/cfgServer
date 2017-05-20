package cluster

import(
	"flag"
	"fmt"
	"testing"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/redis"
)

func TestUpdateDiscovery(t *testing.T) {
	initCh := make(chan int)
	UpdateDiscovery("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", initCh)
}

func TestDiscoverCron(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	notifyCh := make(chan int)
	initCh := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh, initCh)
	up := <- notifyCh
	fmt.Println(up)
}

func TestUpdateNodeInfo(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	notifyCh := make(chan int)
	initCh   := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh, initCh)
	up := <- notifyCh


	if up == 1 {
		for host, _ :=range Seeds {
			nodeInfo, _ := redis.Info(host)
			node := meta.Node{}
			UpdateNodeInfo(&node, nodeInfo)
			fmt.Println(node)
		}
	}
}

func TestUpdateSeedNodes(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	notifyCh := make(chan int)
	initCh   := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh, initCh)
	up := <- notifyCh

	if up == 1 {
		UpdateSeedNodes(Seeds)
	}
}

func TestProbeSeedNodes(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	notifyCh := make(chan int)
	initCh   := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh, initCh)
	up := <- notifyCh

	if up == 1 {
		UpdateSeedNodes(Seeds)
		ProbeSeedNodes()
	}
}

func TestProbeCron(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	notifyCh := make(chan int)
	initCh   := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh, initCh)
	
	up := <- notifyCh
	if up == 1 {
		ProbeCron(notifyCh)
	}
}


func TestPostSeedsToMasterCfg(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.94.46.20:2335", initCh, notifyCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}


	notifyCh2 := make(chan int)
	initCh2   := make(chan error)
	go DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh2, initCh2)
	
	up := <- notifyCh2
	if up == 1 {
		ProbeCron(notifyCh2)
	}

}

func TestUpdateNodeStates(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.94.46.20:2335", notifyCh, initCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}


	Init()
	seeds := make(map[string]SeedNode)

	n1 := meta.Node{
		NodeID:		"13.29.08.98:2300",
		Tag:		"bj:yf:tc",
		Role:		"master",
		Status:		"fail",
		ParentID:	"",
	}
	n2 := meta.Node{
		NodeID:		"19.29.08.90:2300",
		Tag:		"bj:yf:tc",
		Role:		"slave",
		Status:		"fail",
		ParentID:	"13.29.08.98:2300",
	}
	sn1 := SeedNode{
		node:		n1,
		status:		"pfail",
		version:	9,
	}
	sn2 := SeedNode{
		node:		n2,
		status:		"fail",
		version:	10,
	}

	seeds[n1.NodeID] = sn1
	seeds[n2.NodeID] = sn2
	glog.Info(seeds)
	cfg   := &meta.CfgNode {
		NodeID:		"10.234.18.90",
		Region:		"bj",
		Status:		"connected",
	}
	UpdateNodeStates(seeds ,cfg)
	UpdateClusterState() 

}