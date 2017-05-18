package role

import (
	// "fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/meta"
)

var roleManager *RoleManager
var CFG_MUTEX string

const (
	ROLE_MASTER = iota
	ROLE_SLAVE
)

type RoleManager struct {
	own    meta.CfgNode
	role   int
	slaves []meta.CfgNode
}

func SelfRole() int {
	return roleManager.role
}

func IsMasterCfg() bool {
	if roleManager.role == ROLE_MASTER {
		return true
	}

	return false
}

/*
 * slave manger
 */
func AddSlave(node meta.CfgNode) {
	for _, slave := range roleManager.slaves {
		if strings.EqualFold(slave.NodeID, node.NodeID) {
			slave.Status = meta.STATUS_CONNECTED
			return
		}
	}

	node.Status = meta.STATUS_CONNECTED
	roleManager.slaves = append(roleManager.slaves, node)
	glog.Info("Add slave cfg:", node)

	return
}

func RemoveSlave(node meta.CfgNode) {
	for index, slave := range roleManager.slaves {
		if strings.EqualFold(slave.NodeID, node.NodeID) {
			roleManager.slaves = append(roleManager.slaves[:index], roleManager.slaves[index+1:]...)
			return
		}
	}
	glog.Info("Remove slave cfg:", node)

	return
}

// period update
func (roleManager *RoleManager) UpdateSlaveStatus() {
	for index, slave := range roleManager.slaves {
		switch slave.Status {
		case meta.STATUS_CONNECTED:
			roleManager.slaves[index].Status = meta.STATUS_NORMAL
		case meta.STATUS_NORMAL:
			roleManager.slaves[index].Status = meta.STATUS_PFAIL
		case meta.STATUS_PFAIL:
			roleManager.slaves[index].Status = meta.STATUS_FAIL
		case meta.STATUS_FAIL:
			RemoveSlave(roleManager.slaves[index])
		}
		glog.Info("Slave cfg status change:", roleManager.slaves)
	}

	return
}

// period check
func (roleManager *RoleManager) CheckCfgConfig() bool {
	nodes := roleManager.slaves
	if len(nodes) == 0 {
		return false
	}

	oldNodes := meta.SlaveCfgNodes()
	if len(oldNodes) != len(nodes) {
		return true
	}

	for _, n := range oldNodes {
		for _, m := range nodes {
			if !strings.EqualFold(n.NodeID, m.NodeID) ||
				!strings.EqualFold(n.Region, m.Region) {
				return true
			}
		}
	}

	return false
}

/*
 * interact with metadb
 */
func AcquireLock(key, value string, expire int) (bool, error) {
	conn := meta.MetaDBConn()
	reply, err := conn.Cmd("conlock.setnx", key, value, expire).Str()
	if !strings.Contains(reply, "OK") && !strings.Contains(reply, "wrong type") {
		return false, err
	}

	return true, nil
}

func ExtendLock(key, value string, expire int) (bool, error) {
	conn := meta.MetaDBConn()
	reply, err := conn.Cmd("conlock.extend", key, value, expire).Str()
	if !strings.Contains(reply, "OK") && !strings.Contains(reply, "wrong type") {
		return false, err
	}

	return true, nil
}

func UpdateMasterCfg(node meta.CfgNode) error {
	err := meta.UpdateMetaDB(".CfgMeta.MasterCfgNode", &node)
	if err != nil {
		return err
	}
	meta.UpdateMasterCfgNode(node)
	glog.Info("Update master cfgs:", node)

	return nil
}

func UpdateSlaveCfg(nodes []meta.CfgNode) error {
	err := meta.UpdateMetaDB(".CfgMeta.CfgNodes", &nodes)
	if err != nil {
		return err
	}
	meta.UpdateSlaveCfgNode(nodes)
	glog.Info("Update slave cfgs:", nodes)

	return nil
}

func Run(initCh chan error) {
	roleManager = &RoleManager{
		own:    meta.NewCfgNode(meta.CurrID(), meta.Region(), meta.STATUS_CONNECTED),
		slaves: []meta.CfgNode{},
		role:   ROLE_SLAVE,
	}
	CFG_MUTEX = meta.AppName() + "_cfgserver_mutex"

	ok, err := AcquireLock(CFG_MUTEX, meta.CurrID(), meta.CfgMutexExpire())
	if ok {
		glog.Info("Acquire lock, change master")
		roleManager.role = ROLE_MASTER
		err = UpdateMasterCfg(roleManager.own)
		if err != nil {
			glog.Error("Update master failed:", err)
			initCh <- err
		}
	} else if err != nil {
		glog.Error("Acquire lock error:", err)
		initCh <- err
	}

	roleTickChan := time.NewTicker(time.Second * time.Duration(meta.CfgCheckMutextTimeout())).C
	for {
		select {
		case <- roleTickChan:
			if !meta.IsMasterRegion(roleManager.own.Region) {
				glog.Info("Cross region slave do not participate in compete")
				break
			}
			if SelfRole() == ROLE_MASTER {
				ok, err = ExtendLock(CFG_MUTEX, meta.CurrID(), meta.CfgMutexExpire())
				if !ok {
					roleManager.role = ROLE_SLAVE
					glog.Info("Master change slave")
					break
				} else if err != nil {
					glog.Error("Extend lock error:", err)
					initCh <- err
				}

				roleManager.UpdateSlaveStatus()
				change := roleManager.CheckCfgConfig()
				if change {
					err = UpdateSlaveCfg(roleManager.slaves)
					if err != nil {
						glog.Error("Update slave cfgs failed:", roleManager.slaves)
						initCh <- err
					}
				}
			} else if SelfRole() == ROLE_SLAVE {
				glog.Info("Try to acquire lock")
				ok, err = AcquireLock(CFG_MUTEX, meta.CurrID(), meta.CfgMutexExpire())
				if ok {
					roleManager.role = ROLE_MASTER
					glog.Info("Slave change master")
					err = UpdateMasterCfg(roleManager.own)
					if err != nil {
						glog.Error("Update master failed:", err)
						initCh <- err
					}
					break
				} else if err != nil {
					glog.Error("Acquire lock error:", err)
					initCh <- err
				}
			}
		}
	}

}
