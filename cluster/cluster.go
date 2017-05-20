package cluster

import (
	"time"
	"strings"
	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
)


type NodeState struct {
	node      		*meta.Node
	version 		int64
	time      		time.Time
	voteTotal 		int
	voteFail  		int
	voteLocal 		int
	fail			bool
}

type ReplicateSet struct {
	master *meta.Node
	slaves map[*meta.Node][]*meta.Node
}

type ClusterState struct {
	cluster       	*meta.ClusterMeta
	version			int64
	update			bool
	nodeStates    	map[string]*NodeState
}

var cs *ClusterState

func UpdateNodeStates(seeds map[string]SeedNode ,cfg *meta.CfgNode) {
	cs.version = cs.version + 1
	for host, seed := range seeds {
		ns, ok := cs.nodeStates[host]
		if !ok {
			ns = &NodeState{
				node: 		&seed.node,
				version: 	cs.version,
				time:   	time.Now(),
			}
			cs.update = true
			cs.nodeStates[host] = ns
			continue
		}

		if !ns.node.Equal(&seed.node) {
			cs.update = true
			ns.node = &seed.node
		}
		ns.version = cs.version

		if !ns.fail && strings.EqualFold(seed.status, SEED_FAIL) {
			if ns.voteTotal == 0 {
				ns.time = time.Now()
			}

			endTime := ns.time.Add(time.Duration(meta.ClusterNodeTimeout() * 2))
			if time.Now().After(endTime) {
				cs.update    = true
				ns.voteTotal = 0
				ns.voteFail  = 0
				ns.voteLocal = 0
				ns.node 	 = &seed.node
				continue
			}

			ns.voteTotal = ns.voteTotal + 1
			ns.voteFail  = ns.voteFail +1
			if strings.EqualFold(cfg.Region, ns.node.Region()) {
				ns.voteLocal = ns.voteLocal + 1
			}

			if ns.voteFail > (meta.SlaveCfgNum() + 1)/2 && ns.voteLocal > 0{
				ns.fail = true
				ns.node.Status = meta.NS_FAIL
			}
		}

		if ns.fail && strings.EqualFold(seed.status, SEED_LIVE) {
			cs.update    = true
			ns.voteTotal = 0
			ns.voteFail  = 0
			ns.voteLocal = 0
			ns.node 	 = &seed.node
		}
	}

	for host, ns :=range cs.nodeStates {
		if ns.version == cs.version {
			continue
		}
		cs.update    = true
		delete(cs.nodeStates, host)
	}
}

func UpdateClusterState() error {
	if !cs.update {
		return nil
	}

	cs.update = false
	topo := meta.TopoMeta{
		Nodes: []meta.Node{},
	}
	for _, ns :=range cs.nodeStates {
		topo.Nodes = append(topo.Nodes, *ns.node)
	}

	err := meta.UpdateMetaDB(".TopoMeta", &topo)
	if err != nil {
		glog.Error("Update cluster state error:", err)
		return err
	}
	meta.SetTopo(topo)

	return nil
}

func Init() {
	cs = &ClusterState{
		cluster:       	meta.ClusterConfig(),
		version:		0,
		update:			false,
		nodeStates:    	make(map[string]*NodeState),
	}
	glog.Info(cs)
}
