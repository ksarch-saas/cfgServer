package command

import(	
	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/role"
	// "github.com/ksarch-saas/cfgServer/cluster"
)

type UpdateNodesCommand struct {
	Region 			string       	
	CfgID			string			
	Seeds  			interface{}
}

func (self *UpdateNodesCommand) Execute(c *Controller) (Result, error) {
	glog.Info("Receive slave seeds info:", self)
	node := meta.CfgNode{
		Region:			self.Region,
		NodeID:			self.CfgID,
		Status:			meta.STATUS_CONNECTED,
	}
	role.AddSlave(node)

	// seeds := self.Seeds.([]cluster.SeedNode)
	// cluster.UpdateNodeStates(seeds, &node)
	// err := cluster.UpdateClusterState()
	return nil, nil
}
