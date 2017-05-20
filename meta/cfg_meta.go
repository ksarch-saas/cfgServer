package meta

import (
	"strings"
)

type CfgNode struct {
	NodeID 		string
	Region 		string
	Status		string
}

type CfgMeta struct {
	CheckoutMutexTimeout int
	MutexExpire          int
	MasterRegion		 string
	MasterCfgNode        CfgNode
	CfgNodes             []CfgNode
}

const (
	DEFAULT_CHECKOUT_MUTEX_TIMEOUT = 20
	DEFAULT_MUTEX_EXPIRE           = 60
)

const (
	STATUS_CONNECTED = "connected"
	STATUS_NORMAL    = "normal"
	STATUS_PFAIL     = "PFail"
	STATUS_FAIL      = "Fail"
)

func (cfgMeta *CfgMeta) FetchCfgMeta() error {
	err := FetchMetaDB(".CfgMeta", cfgMeta)
	if err != nil {
		return err
	}

	if cfgMeta.CheckoutMutexTimeout == CONFIG_NIL {
		cfgMeta.CheckoutMutexTimeout = DEFAULT_CHECKOUT_MUTEX_TIMEOUT
	}
	if cfgMeta.MutexExpire == CONFIG_NIL {
		cfgMeta.MutexExpire = DEFAULT_MUTEX_EXPIRE
	}
	
	return nil
}

func NewCfgNode(nodeID string, region string, status string) CfgNode {
	node := CfgNode{NodeID:nodeID, Region:region, Status:status}
	return node
}

func UpdateMasterCfgNode(node CfgNode) {
	masterNode := meta.cfgConfig.MasterCfgNode
	masterNode.NodeID = node.NodeID
	masterNode.Region = node.Region
	masterNode.Status = node.Status
}

func UpdateSlaveCfgNode(nodes []CfgNode) {
	meta.cfgConfig.CfgNodes = append(meta.cfgConfig.CfgNodes[:0] ,nodes[:]...)
}

func SlaveCfgNodes() []CfgNode {
	return meta.cfgConfig.CfgNodes
}

func CfgMutexExpire() int {
	return meta.cfgConfig.MutexExpire
}

func CfgCheckMutextTimeout() int {
	return meta.cfgConfig.CheckoutMutexTimeout 
}

func MasterCfgAdress() string{
	return meta.cfgConfig.MasterCfgNode.NodeID
}

func SlaveCfgNum() int {
	return len(meta.cfgConfig.CfgNodes)
}

func MasterCfgRegion() string {
	return meta.cfgConfig.MasterRegion
}

func IsMasterCfgRegion(region string) bool {
	return strings.EqualFold(meta.cfgConfig.MasterRegion, region)
}