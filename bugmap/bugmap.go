package bugmap

import (
	"bug/coordinates"
	"bug/definitions"
)

type Level struct {
	Dimensions    coordinates.Dimension `json:"Dimensions"`
	Nodes         []*BugNode            `json:"Nodes"`
	walkablenodes []*BugNode            `json:"-"`
}

func (level *Level) GetWalkableNodes() []*BugNode {

	if len(level.walkablenodes) == 0 {
		for _, node := range level.Nodes {
			if definitions.NodeTypeBlank < node.Nodetype && node.Nodetype < definitions.NodeTypeWall {
				level.walkablenodes = append(level.walkablenodes, node)
			}
		}
	}

	return level.walkablenodes
}
