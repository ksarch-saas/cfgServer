package meta


type FailoverEntity struct {
	NodeID				string
	Role				string
	Region				string
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
