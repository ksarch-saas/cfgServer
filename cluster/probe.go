package cluster

import (
	"fmt"
	"time"
	"strconv"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/utils"
	"github.com/ksarch-saas/cfgServer/redis"
	"github.com/ksarch-saas/cfgServer/react/api"
)

const (
	UPDATE_VERSION_INIT 		= 0
	FIRST_HALF_PROBE			= 0
	LAST_HAlE_ROBE 				= 1
)

const (
	SEED_LIVE 					= "connected"
	SEED_PFAIL					= "pfail"
	SEED_FAIL					= "fail"
)
var (
	nodeUpdateVersion 	int64
	seedNodes 			map[string]SeedNode
)

type SeedNode struct {
	node		meta.Node
	status 		string		/* current cfg view, the value is  SEED_PFAIL  SEED_FAIL SEED_LIVE*/
	version		int64
}

func UpdateNodeInfo(node *meta.Node, info map[string]string) {
	if info == nil {
		return 
	}

	node.Status			= meta.NS_ONLINE

	master_host, mh 	:= info["master_host"]
	master_port, mp 	:= info["master_port"]
	if mh && mp {
		node.ParentID 	= fmt.Sprintf("%s:%s", master_host, master_port)
		node.Role		= meta.NR_SLAVE
	} else {
		node.ParentID 	= ""
		node.Role		= meta.NR_MASTER
	}

	node.SlotRange		= []meta.Range{}
	s_num, ok := info["slot_num"]
	if !ok {
		return
	}
	num ,err := strconv.Atoi(s_num)
	if err != nil {
		glog.Info("slot_num atoi error:", err)
		return
	}
	for i := 0; i < num; i++ {
		start := fmt.Sprintf("slot_start_%d", i)
		end   := fmt.Sprintf("slot_end_%d", i)
		slot_left, sl 	:= info[start]
		slot_right, sr 	:= info[end]
		if !sl || !sr {
			continue
		}

		left ,err := strconv.Atoi(slot_left)
		if err != nil {
			glog.Info("slot_left atoi error:", err)
			continue
		}
		right ,err := strconv.Atoi(slot_right)
		if err != nil {
			glog.Info("slot_right atoi error:", err)
			continue
		}

		slotRange := meta.Range{
			Left:		left,
			Right:		right,
		}
		node.SlotRange = append(node.SlotRange, slotRange)
	}

	return
}

func UpdateSeedNodes(seeds map[string]string) {
	nodeUpdateVersion = nodeUpdateVersion + 1
	for host, tag :=range seeds {
		sNode, ok := seedNodes[host]
		if ok {
			continue
		}

		sNode.version  		= nodeUpdateVersion
		sNode.status   		= SEED_LIVE
		node		  	   := &sNode.node
		node.NodeID    		= host
		node.Tag 	   		= tag
		nodeInfo, err 	   := redis.Info(host)
		if err != nil {
			glog.Error("Update node info error:", err, host)
			sNode.status    = SEED_FAIL
			node.SlotRange	= []meta.Range{}
			node.Status		= meta.NS_OFFLINE
			node.ParentID	= ""
			node.Role		= meta.NR_MASTER
			continue
		}
		UpdateNodeInfo(node, nodeInfo)
		seedNodes[host] 	= sNode
	}

	for key, sdn :=range seedNodes{
		if sdn.version == nodeUpdateVersion {
			continue
		}
		glog.Info("Delete node:", sdn)
		delete(seedNodes, key)
	}
	return 
}

func ProbeSeedNodes(){
	for host, sNode :=range seedNodes {
		nodeInfo, err := redis.Info(host)
		if err != nil {
			switch sNode.status {
			case SEED_LIVE:
				sNode.status = SEED_PFAIL
			case SEED_PFAIL:
				sNode.status = SEED_FAIL
			default :
				sNode.status = SEED_FAIL
			}
		}
		UpdateNodeInfo(&sNode.node, nodeInfo)
	}
}

func PostSeedsToMasterCfg(seeds map[string]SeedNode) {
	url := "http://" + meta.MasterCfgAdress() + api.UpdateNodesPath
	req := api.UpdateNodesParams{
			Region: 	meta.Region(),
			CfgID:		meta.CurrID(),
			Seeds:  	seeds,
	}
	glog.Info("Post seeds to master cfg:", req)
	res, err := utils.HttpPost(url, req, 5*time.Second)
	glog.Info(res, err)
}

func ProbeCron(notifyCh chan int) {
	nodeUpdateVersion = UPDATE_VERSION_INIT
	seedNodes 		  = make(map[string]SeedNode)

	// tickChan := time.NewTicker(time.Second * time.Duration(meta.ProbeTimeout())).C
	tickChan := time.NewTicker(time.Second * 1).C
	for {
		select {
		case change := <- notifyCh :
			if change != UPDATE_NEW_SEEDS {
				break
			}
			UpdateSeedNodes(Seeds)
		case <- tickChan:
			ProbeSeedNodes()
			PostSeedsToMasterCfg(seedNodes)
		}
	}
}