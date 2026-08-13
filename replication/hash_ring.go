package replication

type Node struct {
	ID   string
	Addr string
	Hash string
}

type Ring struct {
	Nodes []Node
}
