package command

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/ksarch-saas/cfgServer/role"
)

type Result interface{}

type CommandType int

const (
	REGION_COMMAND CommandType = iota
	CLUSTER_COMMAND
	MUTEX_COMMAND
	NOMUTEX_COMMAND
)

type Command interface {
	Type() CommandType
	Mutex() CommandType
	Execute(*Controller) (Result, error)
}

var (
	ErrProcessCommandTimedout = errors.New("cfgserver: process command timeout")
	ErrNotMasterCfg	    	  = errors.New("cfgserver: not master cfgserver")
)

type Controller struct {
	mutex          sync.Mutex
}

func NewController() *Controller {
	c := &Controller{
		mutex:          sync.Mutex{},
	}
	return c
}

func (c *Controller) ProcessCommand(command Command, timeout time.Duration) (result Result, err error) {
	switch command.Type() {
	case CLUSTER_COMMAND:
		if !role.IsMasterCfg() {
			return nil, ErrNotMasterCfg
		}
	}

	// 一次处理一条命令，也即同一时间只能在做一个状态变换
	commandType := strings.Split(reflect.TypeOf(command).String(), ".")
	commandName := ""
	if len(commandType) == 2 && commandType[1] != "UpdateRegionCommand" {
		commandName = commandType[1]
	}
	if commandName != "" {
		glog.Infof("OP", "Command: %s, Event:Start", commandName)
	}
	// command must be excuted serial
	if command.Mutex() == MUTEX_COMMAND {
		c.mutex.Lock()
		defer c.mutex.Unlock()
	}

	resultCh := make(chan Result)
	errorCh := make(chan error)

	go func() {
		result, err := command.Execute(c)
		if err != nil {
			errorCh <- err
		} else {
			resultCh <- result
		}
	}()

	select {
	case result = <-resultCh:
	case err = <-errorCh:
	case <-time.After(timeout):
		err = ErrProcessCommandTimedout
	}
	if commandName != "" {
		glog.Infof("OP", "Command: %s, Event:End", commandName)
	}
	return
}

