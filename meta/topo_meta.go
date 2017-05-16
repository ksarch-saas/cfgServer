package meta


type Range struct {
	Left		int
	Right		int
}

type Node struct {
	NodeID					string
	Tag						string
	Role					string
	Capacity				int64
	ReplOffset				int64
	SlotRange				[]Range
	Connected				bool
	Readable				bool
	Writeable				bool
	Status					string
	ParentID				string
	Ping					int64
	Pong					int64
}

type TopoMeta struct {
	Nodes 			[]Node
}

func (topoMeta *TopoMeta) FetchTopoMeta() error{
	err := FetchMetaDB(".TopoMeta" , topoMeta)
	if err != nil {
		return err
	}
	return nil
}