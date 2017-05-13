package meta


type ClusterMeta struct {
	ClusterConfigTimeout		int
	ClusterNodeTimeout			int
	SafeMode					bool
	Shardings					int
	Replicates					int
	MasterIdc					[]string
	Idc							[]string
}

const (
	DEFAULT_CLUSTER_CONFIG_TIMEOUT		= 100
	DEFAULT_CLUSTER_NODE_TIMEOUT		= 1000
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

	return nil
}