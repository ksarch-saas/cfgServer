package meta


type FailoveEntity struct {
	OldMaster			string
	NewMaster			string
}

type FailoverMeta struct {
	AutoFailover				bool
	FailoverInterval			int
	FailoverConcurrency			int
	FailoverDoing				[]FailoveEntity
	FailoverQueue				[]FailoveEntity
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