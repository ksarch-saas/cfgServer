package meta

import (
	"strings"
)

const (	
	NS_ONLINE				= "online"
	NS_FAIL					= "fail"
	NS_OFFLINE				= "offline"

	NR_MASTER				= "master"
	NR_SLAVE				= "slave"
)



type Range struct {
	Left  int
	Right int
}

func (r *Range)Equal(s *Range) bool{
	if r.Left != s.Left {
		return false 
	}
	if r.Right != s.Right {
		return false 
	}
	return true
}

type Node struct {
	NodeID    string
	Tag       string
	Role      string
	SlotRange []Range
	Status    string
	ParentID  string
}

type TopoMeta struct {
	Nodes []Node
}

func (topoMeta *TopoMeta) FetchTopoMeta() error {
	err := FetchMetaDB(".TopoMeta", topoMeta)
	if err != nil {
		return err
	}
	return nil
}

func SetTopo(tp TopoMeta){
	meta.topo = &tp
}

func (node *Node) IsMaster() bool {
	return strings.EqualFold(node.Role, NR_MASTER) 
}

func (node *Node) Region() string {
	index := strings.Index(node.Tag, ":")
	if index == -1 || index == 0 {
		return ""
	} else if index == 1 {
		return node.Tag[0:1]
	} 
	
	return node.Tag[:index-1]
}


func (node *Node) Equal(seed *Node) bool {
	if !strings.EqualFold(node.NodeID, seed.NodeID){
		return false 
	} 
	if !strings.EqualFold(node.Tag, seed.Tag){
		return false 
	} 
	if !strings.EqualFold(node.Role, seed.Role){
		return false 
	} 
	if !strings.EqualFold(node.Status, seed.Status){
		return false 
	} 
	if !strings.EqualFold(node.ParentID, seed.ParentID){
		return false 
	}
	for _, n := range node.SlotRange {
		for _, s := range seed.SlotRange{
			if !s.Equal(&n) {
				return false
			}
		}
	}

	return true
}