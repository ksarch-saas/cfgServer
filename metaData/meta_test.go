package meta

import (
	"fmt"
	"flag"
	"testing"

	"github.com/ksarch-saas/cfgServer/redis"
	"github.com/mediocregopher/radix.v2/cluster"
)


func TestMetaAddress(t *testing.T) {
	meta = &Meta{
		appName		:		"ssdb-test",
		currIdc		:		"tc",
		currId		:		"10.67.17.43:3700",
	}
	res, err := MetaAddress()
	fmt.Println(res, err)
}

func TestFetchMetaDB(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	var c int
	err := FetchMetaDB(".ClusterVersion", &c)
	fmt.Println(c, err)
}

func TestUpdateMetaDB(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	var c = 4
	err := UpdateMetaDB(".ClusterVersion", &c)
	fmt.Println(err)
}

func TestFetchCfgMeta(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := meta.CfgConfig().FetchCfgMeta()
	fmt.Println(meta.CfgConfig(), err)
}

func TestFetchClusterMeta(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := meta.ClusterConfig().FetchClusterMeta()
	fmt.Println(meta.ClusterConfig(), err)
}

func TestFetchFailoverMeta(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := meta.FailoverConfig().FetchFailoverMeta()
	fmt.Println(meta.FailoverConfig(), err)
}

func TestFetchMigrateMeta(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := meta.MigrateConfig().FetchMigrateMeta()
	fmt.Println(meta.MigrateConfig(), err)
}

func TestFetchTopoMeta(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := meta.Topo().FetchTopoMeta()
	fmt.Println(meta.Topo(), err)
}

func TestFetchMeta(t *testing.T) {
	// initCh := make(chan error)
	// Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	// fmt.Println(meta)
	meta = &Meta{
		appName			:		"ssdb-test",
		currIdc			:		"tc",
		currId			:		"",
		clusterVersion	:		0,
		configVersion	:		0,
		metaDBConn		:		&cluster.Cluster{},
		clusterConfig	: 		&ClusterMeta{},
		cfgConfig		:		&CfgMeta{},
		failoverConfig	:		&FailoverMeta{},
		migrateConfig 	:		&MigrateMeta{},
		topo			:		&TopoMeta{},
	}
	addr, err := MetaAddress()
	if err != nil {
		fmt.Println(err)
		return
	}
	dbConn, err := redis.DialCluster(addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	meta.metaDBConn	= dbConn

	err = FetchMeta()
	fmt.Println(meta)
}

func TestRun(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	go Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	err := <-initCh
	if err != nil {
		fmt.Println(err)
	}
}

func TestSetClusterVersion(t *testing.T) {
	initCh := make(chan error)
	Run("ssdb-test", "tc", "10.67.17.43:3700", initCh)
	meta.SetClusterVersion(3)
}