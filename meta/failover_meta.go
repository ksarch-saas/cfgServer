package meta

import (
	"strings"
)

type FailoverEntity struct {
	NodeID				string
	Role				string
	Region				string
	NewID				string
}

type FailoverMeta struct {
	AutoFailover				bool
	FailoverInterval			int
	FailoverConcurrency			int
	FailoverDoing				[]FailoverEntity
	FailoverQueue				[]FailoverEntity
}

const (
	DEFAULT_FAILOVER_INTERVAL 		= 100
	DEFAULT_FAILOVERCONCURRENCY		= 1
)

func (failoverMeta *FailoverMeta) FetchFailoverMeta() error{
	err := FetchMetaDB(".FailoverMeta" , failoverMeta)
	if err != nil {
		return err
	}

	if failoverMeta.FailoverInterval == CONFIG_NIL {
		failoverMeta.FailoverInterval = DEFAULT_FAILOVER_INTERVAL
	}
	if failoverMeta.FailoverConcurrency == CONFIG_NIL {
		failoverMeta.FailoverConcurrency = DEFAULT_FAILOVERCONCURRENCY
	}
	return nil
}

func (failoverMeta *FailoverMeta) FetchFailoverDoing() error{
	err := FetchMetaDB(".FailoverDoing" , &failoverMeta.FailoverDoing)
	if err != nil {
		return err
	}
	return nil
}

func (failoverMeta *FailoverMeta) FetchFailoverQueue() error{
	err := FetchMetaDB(".FailoverQueue" , &failoverMeta.FailoverQueue)
	if err != nil {
		return err
	}
	return nil
}

func GetFailoverInterval() int  {
	return meta.failoverConfig.FailoverInterval
}

func GetFailoverConcurrency() int {
	return meta.failoverConfig.FailoverConcurrency
}

func SetFailoverDoing(list []*FailoverEntity) {
	meta.failoverConfig.FailoverDoing = meta.failoverConfig.FailoverDoing[:0]
	for _, entity := range list {
		meta.failoverConfig.FailoverDoing = append(meta.failoverConfig.FailoverDoing, *entity)
	}
}

func SetFailoverQueue(list []*FailoverEntity) {
	meta.failoverConfig.FailoverQueue = meta.failoverConfig.FailoverQueue[:0]
	for _, entity := range list {
		meta.failoverConfig.FailoverQueue = append(meta.failoverConfig.FailoverQueue, *entity)
	}
}

func (failoverEntity *FailoverEntity) Equal(entity *FailoverEntity) bool{
	if !strings.EqualFold(failoverEntity.NodeID, entity.NodeID){
		return false 
	}
	if !strings.EqualFold(failoverEntity.Role, entity.Role){
		return false 
	}
	if !strings.EqualFold(failoverEntity.Region, entity.Region){
		return false 
	}
	return true
}
