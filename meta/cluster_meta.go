package meta


type ClusterMeta struct {
	ClusterConfigTimeout		int
	ClusterNodeTimeout			int
	VoteTimeout          		int
	FailoverRatio        		int
	SafeMode					bool
	Shardings					int
	Replicates					int
	MasterIdc					[]string
	Idc							[]string
}

const (
	DEFAULT_CLUSTER_CONFIG_TIMEOUT		= 100
	DEFAULT_CLUSTER_NODE_TIMEOUT		= 1000
	DEFAULT_VOTE_TIMEOUT           		= 4 * DEFAULT_CLUSTER_NODE_TIMEOUT
	DEFAULT_FAILOVER_RATIO         		= 50
)

func (clusterMeta *ClusterMeta) FetchClusterMeta() error{
	err := FetchMetaDB(".ClusterMeta" , clusterMeta)
	if err != nil {
		return err
	}

	if clusterMeta.ClusterConfigTimeout == CONFIG_NIL {
		clusterMeta.ClusterConfigTimeout = DEFAULT_CLUSTER_CONFIG_TIMEOUT
	}
	if clusterMeta.ClusterNodeTimeout == CONFIG_NIL {
		 clusterMeta.ClusterNodeTimeout = DEFAULT_CLUSTER_NODE_TIMEOUT
	}
	if clusterMeta.VoteTimeout == CONFIG_NIL {
		clusterMeta.VoteTimeout = DEFAULT_VOTE_TIMEOUT
	}
	if clusterMeta.FailoverRatio == CONFIG_NIL {
		clusterMeta.FailoverRatio = DEFAULT_FAILOVER_RATIO
	}

	return nil
}

func (meta *Meta)ClusterIdcs() []string {
	return meta.clusterConfig.Idc
}

func ProbeTimeout() int {
	return meta.clusterConfig.ClusterNodeTimeout
}

func ClusterNodeTimeout() int {
	return meta.clusterConfig.ClusterNodeTimeout
}