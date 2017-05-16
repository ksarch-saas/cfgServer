package role

import (
	"fmt"
	"flag"
	"testing"

	"github.com/ksarch-saas/cfgServer/meta"
)

func TestLock(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.67.17.43:3700", initCh, notifyCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}
	reply, err := AcquireLock(CFG_MUTEX, meta.CurrID(), meta.CfgMutexExpire())
	fmt.Println(reply, err)
	reply, err = ExtendLock(CFG_MUTEX, meta.CurrID(), meta.CfgMutexExpire())
	fmt.Println(reply, err)
}

func TestUpdateCfg(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.67.17.43:3700", initCh, notifyCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}

	node := meta.CfgNode{
		NodeID: "10.120.39.44:2400",
		Region: "nj",
		Status: "failed",
	}
	err := UpdateMasterCfg(node)

	nodes := []meta.CfgNode{
		meta.CfgNode{
			NodeID: "10.120.39.44:2400",
			Region: "nj",
			Status: "failed",
		},
		meta.CfgNode{
			NodeID: "10.120.38.44:2400",
			Region: "nj",
			Status: "pfaile",
		},
	}
	err = UpdateSlaveCfg(nodes)
	fmt.Println(err)
}

func TestSlaveMangaer(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.67.17.43:3700", initCh, notifyCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}

	roleManager = &RoleManager {
		own:		meta.CfgNode{},
		slaves:		[]meta.CfgNode{},
		role: 		0,
	}
	n1 := meta.CfgNode{
		NodeID: "109.119.39.44:2400",
		Region: "nj",
		Status: "failed",
	}
	n2 := meta.CfgNode{
		NodeID: "109.129.39.44:2400",
		Region: "nj",
		Status: "failed",
	}
	roleManager.AddSlave(n1)
	roleManager.AddSlave(n2)
	fmt.Println(roleManager.slaves)
	roleManager.RemoveSlave(n1)
	fmt.Println(roleManager.slaves)
	for i := 0; i < 5; i++ {
		roleManager.UpdateSlaveStatus()
		fmt.Println(roleManager.slaves)
	}
}

func TestRun(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.67.17.43:3700", initCh, notifyCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}

	Run(initCh)
}
