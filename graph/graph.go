package graph

//import "fmt"

// The constant variables for graph topology
const (
	GRID = "grid"
	RING = "circle"
)

// The constant variables for distance measures
const (
	HOP       = "hop"
	EUCLIDEAN = "euclidean"
)

//var network Network

// The Node struct provides a node abstraction.
type Node struct {
	ID     []int
	Memory int
}

// The Link struct provides a link abstraction.
type Link struct {
	ID          []int
	Age         int
	Reservation int
	IsReserved  bool
	IsActive    bool
	// The IsPruned field is for managing the pruning process.
	IsPruned bool
}

//type Network struct {
//	Nodes []Node
//	Links []Link
//}

// The Topology interface defines the necessary methods for a network abstraction.
type Topology interface {
	Build()
	Clear()
	generateNodes()
	generateLinks()
	//GetLinks() []Link
	GetNodeIDs() [][]int
	GetSize() int
	GetLinks() []*Link
	GetType() string
	//GetNodes() [][]
	GetNeighbors(*Node) ([]*Node, bool)
	GetLinkBetween(n1, n2 *Node) *Link
	Distance(*Node, *Node, string) int
	//Prune([]*Link)
}

// MakeNode creates a node.
//func (node *Node) MakeNode(id []int, memory int) {
//	node.ID = id
//	node.Memory = memory
//}

// IsEqual checks if two nodes are equal, meaning that their fields are the same thing.
func (node *Node) IsEqual(otherNode *Node) bool {
	if node.Memory != otherNode.Memory {
		return false
	}
	if (node.ID == nil) != (otherNode.ID == nil) {
		return false
	}
	if len(node.ID) != len(otherNode.ID) {
		return false
	}
	for i, v := range node.ID {
		if v != otherNode.ID[i] {
			return false
		}
	}
	return true
}

func IsEqual(id1, id2 []int) bool {
	if (id1 == nil) != (id2 == nil) {
		return false
	}
	if len(id1) != len(id2) {
		return false
	}
	for i, v := range id1 {
		if v != id2[i] {
			return false
		}
	}
	return true
}

func (link *Link) IsActivated() bool {
	return link.IsActive
}

// MakeNode creates a node.
func MakeNode(id []int, memory int) *Node {
	var node Node
	node = Node{ID: id, Memory: memory}
	return &node
}

// MakeLink creates a link.
func MakeLink(id []int, age int, isActive bool) *Link {
	var link Link
	link = Link{ID: id, Age: age, IsActive: isActive, IsPruned: false, Reservation: -1, IsReserved: false}
	return &link
}

////////////// TODO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//func removeLink(link Link)

////////////// TODO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func Prune(links []*Link) {
	for _, link := range links {
		link.IsPruned = true
	}
}

func DepruneLinks(links []*Link) {
	for _, link := range links {
		link.IsPruned = false
	}
}

func Deprune(network Topology) {
	links := network.GetLinks()
	for _, link := range links {
		link.IsPruned = false
	}
}

func FindPosition(id []int, nodes []*Node) int {
	for i, _ := range nodes {
		if IsEqual(id, nodes[i].ID) {
			return i
		}
	}
	return -1
}

//func CopyLinks()

// BuildGraph builds the desired graph.
//func BuildGraph(topology string) Topology {
//	if topology == GRID {
//		var grid *Grid = new(Grid)
//		grid.Build()
//		return grid
//	}
//	return nil
//}
