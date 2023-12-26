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
