package command

import(	
	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/role"
	"github.com/ksarch-saas/cfgServer/cluster"
)

type UpdateNodesCommand struct {
	Region 			string       	
	CfgID			string			
	Seeds  			map[string]cluster.SeedNode
}

func (self *UpdateNodesCommand) Execute(c *Controller) (Result, error) {
	glog.Info("Receive slave seeds info:", self)
	node := meta.CfgNode{
		Region:			self.Region,
		NodeID:			self.CfgID,
		Status:			meta.STATUS_CONNECTED,
	}
	role.AddSlave(node)

	cluster.UpdateNodeStates(self.Seeds, &node)
	err := cluster.UpdateClusterState()
	return nil, err
}
