package meta

import (
	"fmt"
	"time"
	"errors"
	"strings"
	"math/rand"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/redis"
	"github.com/ksarch-saas/cfgServer/utils"
	"github.com/mediocregopher/radix.v2/cluster"
)


var IdcToRegion = map[string]string{
	"jx":   "bj",
	"tc":   "bj",
	"nj":   "nj",
	"nj03": "nj",
	"gz":   "gz",
	"hz":   "hz",
	"sh":   "sh",
	"sz":   "sz",
	"nmg":  "nmg",
	"yq":   "yq",
}

var RegionToIdcs = map[string][]string{
	"bj":   {"jx","tc"},
	"nj":   {"nj","nj03"},
	"gz":   {"gz"},
	"hz":   {"hz"},
	"sh":   {"sh"},
	"sz":   {"sz"},
	"nmg":  {"nmg"},
	"yq":   {"yq"},
}

var meta *Meta

type Meta struct {
	appName						string
	currIdc						string
	currId						string
	clusterVersion				int
	configVersion				int
	metaDBConn					*cluster.Cluster
	clusterConfig				*ClusterMeta
	cfgConfig					*CfgMeta
	failoverConfig				*FailoverMeta
	migrateConfig				*MigrateMeta
	topo						*TopoMeta
}

const (
	CONFIG_NIL							= 0 
	CHECKOUT_CFGVERSION_TIMEOUT			= 1
)


/*
 * Meta Struct Operations
 */
func (meta *Meta)AppName() string {
	return meta.appName
}

func (meta *Meta)Idc() string {
	return meta.currIdc
}

func (meta *Meta)CurrID() string {
	return meta.currId
}

func (meta *Meta)MetaDBConn() *cluster.Cluster {
	return meta.metaDBConn
}

func (meta *Meta)ClusterVersion() int {
	return meta.clusterVersion
}

func (meta *Meta)ConfigVersion() int {
	return meta.configVersion
}

func (meta *Meta)ClusterConfig() *ClusterMeta {
	return meta.clusterConfig
}

func (meta *Meta)CfgConfig() *CfgMeta {
	return meta.cfgConfig
}

func (meta *Meta)FailoverConfig() *FailoverMeta {
	return meta.failoverConfig
}

func (meta *Meta)MigrateConfig() *MigrateMeta {
	return meta.migrateConfig
}

func (meta *Meta)Topo() *TopoMeta {
	return meta.topo
}

func (meta *Meta)String() string {
	type output struct{
		AppName						string
		CfgIdc						string
		Address						string
		ClusterVersion				int
		ConfigVersion				int
		ClusterConfig				ClusterMeta
		CfgConfig					CfgMeta
		FailoverConfig				FailoverMeta
		MigrateConfig				MigrateMeta
		Topo						TopoMeta
	}
	info := &output{
		AppName				: 		meta.AppName(),
		CfgIdc				: 		meta.Idc(),
		Address				:		meta.CurrID(),
		ClusterVersion		:		meta.ClusterVersion(),
		ConfigVersion		:		meta.ConfigVersion(),
		ClusterConfig		:		*meta.ClusterConfig(),
		CfgConfig			:		*meta.CfgConfig(),
		FailoverConfig		:		*meta.FailoverConfig(),
		MigrateConfig		:		*meta.MigrateConfig(),
		Topo				:		*meta.Topo(),
	}

	infoByte , _ :=json.Marshal(info)
	return fmt.Sprintln(string(infoByte))
}

/* 
 * Set function rule adapt to all structs:
 * 	1. set redis, return success ,go on
 * 	2. set meta struct, return true
 * 	3. 1-2 is sync and ordered
 */
func (meta *Meta)SetClusterVersion(clusterVersion int) error{
	err := UpdateMetaDB(".ClusterVersion", &clusterVersion)
	if err != nil {
		glog.Error(err)
		return err
	}
	meta.clusterVersion = clusterVersion
	return nil
}


/*
 * General functions for all structs
 */
func MetaAddress() (string, error) {
	bns := "redis4db-" + meta.AppName() + ".osp." + meta.Idc()
	addrString, err := utils.GetAddrFromBNS(bns)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	hosts := strings.Fields(addrString)

	var addrs []string
	for _ , host := range hosts {
		addrs = append(addrs, host)
	}
	if len(addrs) == 0 {
		glog.Error("BNS return empty")
		return "", errors.New("BNS return empty")
	}
	return addrs[rand.Intn(len(addrs))], nil
}

func FetchMetaDB(jsonObj string, data interface{}) error{
	conn := meta.MetaDBConn()
	reply, err := conn.Cmd("json.get", meta.AppName(), jsonObj).Bytes()
	if err != nil {
		glog.Error(err)
		return err
	}
	err = json.Unmarshal(reply, data)
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func UpdateMetaDB(jsonObj string , data interface{}) error {
	dataByte , err :=json.Marshal(data)
	if err != nil {
		return err
	}
	dataString := string(dataByte)
	conn := meta.MetaDBConn()
	reply, err := conn.Cmd("json.set", meta.AppName(), jsonObj, dataString).Str()
	if err != nil {
		glog.Error(err)
		return err
	}
	if !strings.Contains(reply, "OK") {
		glog.Error(errors.New(reply))
		return errors.New(reply)
	}
	return nil
}

/*
 * Meta fetch
 */

func FetchClusterVersion(clusterVersion *int) error {
	return FetchMetaDB(".ClusterVersion", clusterVersion)
}

func FetchConfigVersion(configVersion *int) error {
	return FetchMetaDB(".ConfigVersion", configVersion)
}

func FetchClusterConfig(clusterMeta *ClusterMeta) error {
	return clusterMeta.FetchClusterMeta()
}

func FetchCfgConfig(cfgMeta *CfgMeta) error {
	return cfgMeta.FetchCfgMeta()
}

func FetchFailoverConfig(failoverMeta *FailoverMeta) error {
	return failoverMeta.FetchFailoverMeta()
}

func FetchMigrateConfig(migrateMeta *MigrateMeta) error {
	return migrateMeta.FetchMigrateMeta()
}

func FetchTopo(topoMeta *TopoMeta) error {
	return topoMeta.FetchTopoMeta()
}

func FetchMeta() (error) {
	err := FetchClusterVersion(&meta.clusterVersion)
	if err != nil {
		return err
	}
	err = FetchConfigVersion(&meta.configVersion)
	if err != nil {
		return err
	}
	err = FetchClusterConfig(meta.clusterConfig)
	if err != nil {
		return err
	}
	err = FetchCfgConfig(meta.cfgConfig)
	if err != nil {
		return err
	}
	err = FetchFailoverConfig(meta.failoverConfig)
	if err != nil {
		return err
	}
	err = FetchMigrateConfig(meta.migrateConfig)
	if err != nil {
		return err
	}
	err = FetchTopo(meta.topo)
	if err != nil {
		return err
	}
	return nil
}

func Run(appName, idc, address string, initCh chan error){
	meta = &Meta{
		appName			:		appName,
		currIdc			:		idc,
		currId			:		address,
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
		initCh <- fmt.Errorf("BNS error: %v", err)
		return 
	}
	dbConn, err := redis.DialCluster(addr)
	if err != nil {
		initCh <- fmt.Errorf("MetaDB: can't connect: %v", err)
		return 
	}
	meta.metaDBConn	= dbConn
	err = FetchMeta()
	if err != nil {
		initCh <- fmt.Errorf("Fetch Meta error: %v", err)
		return 
	}

	var cfgVer = meta.configVersion
	tickChan := time.NewTicker(time.Second * CHECKOUT_CFGVERSION_TIMEOUT).C
	for {
		select{
		case <- tickChan :
			err := FetchConfigVersion(&cfgVer)
			if err != nil {
				initCh <- fmt.Errorf("Fetch CfgVersion error: %v", err)
				break
			}
			if cfgVer != meta.configVersion {
				glog.Info("Current configversion:", cfgVer, ", ", "MetaDb configVersion:", meta.configVersion)
				glog.Info("Updat meta begin:", meta)
				err = FetchMeta()
				if err != nil {
					initCh <- fmt.Errorf("Fetch Meta error: %v", err)
					continue
				}
				glog.Info("Updat meta end:", meta)
			}
		}
	}
}