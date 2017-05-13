package meta


type CfgNode struct{
	NodeID				string
}

type CfgMeta struct {
	CheckoutMutexTimeout		int
	MutexExpire					int
	VoteTimeout					int
	FailoverRatio				int
	MasterCfgNode				CfgNode
	CfgNodes					[]CfgNode
}

const (
	DEFAULT_CHECKOUT_MUTEX_TIMEOUT		= 1000
	DEFAULT_MUTEX_EXPIRE 				= 10000
	DEFAULT_VOTE_TIMEOUT 				= 4 * DEFAULT_CLUSTER_NODE_TIMEOUT
	DEFAULT_FAILOVER_RATIO			    = 50
)

func (cfgMeta *CfgMeta) FetchCfgMeta() error{
	err := FetchMetaDB(".CfgMeta" , cfgMeta)
	if err != nil {
		return err
	}

	if cfgMeta.CheckoutMutexTimeout == CONFIG_NIL {
		cfgMeta.CheckoutMutexTimeout = DEFAULT_CHECKOUT_MUTEX_TIMEOUT
	}
	if cfgMeta.MutexExpire == CONFIG_NIL {
		 cfgMeta.MutexExpire  = DEFAULT_MUTEX_EXPIRE
	}
	if cfgMeta.VoteTimeout == CONFIG_NIL {
		cfgMeta.VoteTimeout = DEFAULT_VOTE_TIMEOUT
	}
	if cfgMeta.FailoverRatio == CONFIG_NIL {
		cfgMeta.FailoverRatio = DEFAULT_FAILOVER_RATIO
	}

	return nil
}