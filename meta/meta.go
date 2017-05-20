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
	"bj":  {"jx", "tc"},
	"nj":  {"nj", "nj03"},
	"gz":  {"gz"},
	"hz":  {"hz"},
	"sh":  {"sh"},
	"sz":  {"sz"},
	"nmg": {"nmg"},
	"yq":  {"yq"},
}

var meta *Meta

type Meta struct {
	appName        string
	currIdc        string
	currId         string
	clusterVersion int
	configVersion  int
	metaDB         string
	metaDBConn     *cluster.Cluster
	clusterConfig  *ClusterMeta
	cfgConfig      *CfgMeta
	failoverConfig *FailoverMeta
	migrateConfig  *MigrateMeta
	topo           *TopoMeta
}

const (
	CONFIG_NIL                  = 0
	CHECKOUT_CFGVERSION_TIMEOUT = 1
	INIT_META_SUCCESS           = 1
)

/*
 * Meta Struct Operations
 */

func AppName() string {
	return meta.appName
}

func Idc() string {
	return meta.currIdc
}

func CurrID() string {
	return meta.currId
}

func MetaDbName() string {
	return meta.metaDB
}

func MetaDBConn() *cluster.Cluster {
	return meta.metaDBConn
}

func ClusterVersion() int {
	return meta.clusterVersion
}

func ConfigVersion() int {
	return meta.configVersion
}

func ClusterConfig() *ClusterMeta {
	return meta.clusterConfig
}

func CfgConfig() *CfgMeta {
	return meta.cfgConfig
}

func FailoverConfig() *FailoverMeta {
	return meta.failoverConfig
}

func MigrateConfig() *MigrateMeta {
	return meta.migrateConfig
}

func Topo() *TopoMeta {
	return meta.topo
}

func Region() string {
	return IdcToRegion[meta.currIdc]
}

func IsMasterRegion(region string) bool {
	for _, midc := range meta.clusterConfig.MasterIdc {
		if strings.EqualFold(IdcToRegion[midc], region) {
			return true
		}
	}

	return false
}

func (meta *Meta) String() string {
	type output struct {
		AppName        string
		MetaDBName     string
		CfgIdc         string
		Address        string
		ClusterVersion int
		ConfigVersion  int
		ClusterConfig  ClusterMeta
		CfgConfig      CfgMeta
		FailoverConfig FailoverMeta
		MigrateConfig  MigrateMeta
		Topo           TopoMeta
	}
	info := &output{
		AppName:        AppName(),
		MetaDBName:     MetaDbName(),
		CfgIdc:         Idc(),
		Address:        CurrID(),
		ClusterVersion: ClusterVersion(),
		ConfigVersion:  ConfigVersion(),
		ClusterConfig:  *ClusterConfig(),
		CfgConfig:      *CfgConfig(),
		FailoverConfig: *FailoverConfig(),
		MigrateConfig:  *MigrateConfig(),
		Topo:           *Topo(),
	}

	infoByte, _ := json.Marshal(info)
	return fmt.Sprintln(string(infoByte))
}

func SetClusterVersion(clusterVersion int) error {
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
	bns := "redis4db-" + MetaDbName() + ".osp." + Idc()
	addrString, err := utils.GetAddrFromBNS(bns)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	hosts := strings.Fields(addrString)

	var addrs []string
	for _, host := range hosts {
		addrs = append(addrs, host)
	}
	if len(addrs) == 0 {
		glog.Error("BNS return empty")
		return "", errors.New("BNS return empty")
	}
	return addrs[rand.Intn(len(addrs))], nil
}

func FetchMetaDB(jsonObj string, data interface{}) error {
	conn := MetaDBConn()
	reply, err := conn.Cmd("json.get", AppName(), jsonObj).Bytes()
	if err != nil {
		glog.Error(err, jsonObj)
		return err
	}

	err = json.Unmarshal(reply, data)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func UpdateMetaDB(jsonObj string, data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	dataString := string(dataByte)

	conn := MetaDBConn()
	reply, err := conn.Cmd("json.set", AppName(), jsonObj, dataString).Str()
	if err != nil {
		glog.Error(err, jsonObj, dataString)
		return err
	}
	if !strings.Contains(reply, "OK") {
		glog.Error(reply, jsonObj, dataString)
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

func FetchConfigMeta() error {
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

func FetchFeedMeta() error {
	err := meta.failoverConfig.FetchFailoverDoing()
	if err != nil{
		return err
	}
	err = meta.failoverConfig.FetchFailoverQueue()
	if err != nil{
		return err
	}
	err = meta.migrateConfig.FetchMigrateDoing()
	if err != nil{
		return err
	}
	err = meta.migrateConfig.FetchMigrateTasks()
	if err != nil{
		return err
	}
	err = meta.topo.FetchTopoMeta()
	if err != nil {
		return err
	}

	return nil
}

func Run(appName, metaDbName, idc, address string, notify chan int, initCh chan error) {
	meta = &Meta{
		appName:        appName,
		currIdc:        idc,
		currId:         address,
		clusterVersion: 0,
		configVersion:  0,
		metaDB:         metaDbName,
		metaDBConn:     &cluster.Cluster{},
		clusterConfig:  &ClusterMeta{},
		cfgConfig:      &CfgMeta{},
		failoverConfig: &FailoverMeta{},
		migrateConfig:  &MigrateMeta{},
		topo:           &TopoMeta{},
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
	meta.metaDBConn = dbConn
	err = FetchConfigMeta()
	if err != nil {
		initCh <- fmt.Errorf("Fetch meta error: %v", err)
		return
	}
	err = FetchFeedMeta()
	if err != nil {
		initCh <- fmt.Errorf("Fetch feed error: %v", err)
		return
	}

	glog.Info("Fetch meta data:", meta)
	notify <- INIT_META_SUCCESS

	var cfgVer = meta.configVersion
	tickChan := time.NewTicker(time.Second * CHECKOUT_CFGVERSION_TIMEOUT).C
	for {
		select {
		case <-tickChan:
			err := FetchConfigVersion(&cfgVer)
			if err != nil {
				initCh <- fmt.Errorf("Fetch CfgVersion error: %v", err)
				break
			}
			if cfgVer > meta.configVersion {
				glog.Info("Old configversion:", meta.configVersion, ", ", "Current configVersion:", cfgVer)
				glog.Info("Updat meta begin:", meta)
				err = FetchConfigMeta()
				if err != nil {
					initCh <- fmt.Errorf("Fetch Meta error: %v", err)
					break
				}
				glog.Info("Updat meta end:", meta)
			}
		}
	}
}
