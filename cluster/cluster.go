package cluster

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
)

type NodeState struct {
	Node      *meta.Node
	Version   int64
	Time      time.Time
	VoteTotal int
	VoteFail  int
	VoteLocal int
	Fail      bool
}

type ReplicateSet struct {
	Master *meta.Node
	Slaves []*meta.Node
}

type ClusterState struct {
	Cluster    *meta.ClusterMeta
	Version    int64
	Update     bool
	NodeStates map[string]*NodeState
}

var cs *ClusterState

func GetClusterState() *ClusterState {
	return cs
}

func UpdateNodeStates(seeds map[string]SeedNode, cfg *meta.CfgNode) {
	cs.Version = cs.Version + 1
	for host, seed := range seeds {
		ns, ok := cs.NodeStates[host]
		if !ok {
			ns = &NodeState{
				Node:    &seed.Node,
				Version: cs.Version,
				Time:    time.Now(),
			}
			cs.Update = true
			cs.NodeStates[host] = ns
			continue
		}

		if !ns.Node.Equal(&seed.Node) {
			cs.Update = true
			ns.Node = &seed.Node
		}
		ns.Version = cs.Version

		if !ns.Fail && strings.EqualFold(seed.Status, SEED_FAIL) {
			if ns.VoteTotal == 0 {
				ns.Time = time.Now()
			}

			endTime := ns.Time.Add(time.Duration(meta.ClusterNodeTimeout() * 2))
			if time.Now().After(endTime) {
				cs.Update = true
				ns.VoteTotal = 0
				ns.VoteFail = 0
				ns.VoteLocal = 0
				ns.Node = &seed.Node
				continue
			}

			ns.VoteTotal = ns.VoteTotal + 1
			ns.VoteFail = ns.VoteFail + 1
			if strings.EqualFold(cfg.Region, ns.Node.Region()) {
				ns.VoteLocal = ns.VoteLocal + 1
			}

			if ns.VoteFail > (meta.SlaveCfgNum()+1)/2 && ns.VoteLocal > 0 {
				ns.Fail = true
				ns.Node.Status = meta.NS_FAIL

				if strings.EqualFold(ns.Node.Role, meta.NR_SLAVE) {
					continue
				}
				entity := &meta.FailoverEntity{
					NodeID: ns.Node.NodeID,
					Role:   ns.Node.Role,
					Region: ns.Node.Region(),
					NewID:  "",
				}
				AddFailoverTasks(entity)
			}

		}

		if ns.Fail && strings.EqualFold(seed.Status, SEED_LIVE) {
			cs.Update = true
			ns.VoteTotal = 0
			ns.VoteFail = 0
			ns.VoteLocal = 0
			ns.Node = &seed.Node
		}
	}

	for host, ns := range cs.NodeStates {
		if ns.Version == cs.Version {
			continue
		}
		cs.Update = true
		delete(cs.NodeStates, host)
	}
}

func UpdateClusterState() error {
	if !cs.Update {
		return nil
	}

	cs.Update = false
	topo := meta.TopoMeta{
		Nodes: []meta.Node{},
	}
	for _, ns := range cs.NodeStates {
		topo.Nodes = append(topo.Nodes, *ns.Node)
	}

	err := meta.UpdateMetaDB(".TopoMeta", &topo)
	if err != nil {
		glog.Error("Update cluster state error:", err)
		return err
	}
	meta.SetTopo(topo)

	// update clusterVersion to notify proxy update topo
	clusterVersion := meta.ClusterVersion()
	clusterVersion = clusterVersion + 1
	err = meta.SetClusterVersion(clusterVersion)
	if err != nil {
		glog.Error("Update cluster version error:", err)
		return err
	}

	return nil
}

func Init() {
	cs = &ClusterState{
		Cluster:    meta.ClusterConfig(),
		Version:    0,
		Update:     false,
		NodeStates: make(map[string]*NodeState),
	}
}
