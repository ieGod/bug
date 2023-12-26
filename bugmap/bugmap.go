package bugmap

import "bug/coordinates"

type Level struct {
	Dimensions coordinates.Dimension `json:"Dimensions"`
	Nodes      []*BugNode            `json:"Nodes"`
}
