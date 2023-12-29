package bugmap

import (
	"bug/coordinates"
	"bug/definitions"
)

type BugNode struct {
	Location coordinates.Vector   `json:"Location"`
	Nodetype definitions.NodeType `json:"Nodetype"`
	Nodetile definitions.NodeTile `json:"Nodetile"`
	Previous *BugNode             `json:"-"`
	GCost    int                  `json:"-"`
	HCost    int                  `json:"-"`
}

func (node *BugNode) FCost() int {
	return node.GCost + node.HCost
}

func (n *BugNode) IsInSet(set []*BugNode) bool {
	for _, node := range set {
		if n == node {
			return true
		}
	}
	return false
}

func (n *BugNode) DistanceTo(target *BugNode) int {
	return n.Location.GetManhattanDist(target.Location)
}

func NewBattleNode(loc coordinates.Vector) *BugNode {
	return &BugNode{
		Location: loc,
		Nodetype: definitions.NodeTypeGround,
	}
}
