package command

import(
	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/role"
)

type MergeSeedsCommand struct {
	Region 			string       	
	CfgID			string			
	Seeds  			[]*meta.Node 	
}

func (self *MergeSeedsCommand) Execute(c *Controller) (Result, error) {
	glog.Info("Receive slave seeds info:", self)
	node := meta.CfgNode{
		Region:			self.Region,
		NodeID:			self.CfgID,
		Status:			meta.STATUS_CONNECTED,
	}

	role.AddSlave(node)

	return nil, nil
}
