package bugmap

import (
	"bug/definitions"
	"fmt"
	"math"
	"sort"
)

const (
	SAFETY_CHECK = 25000
)

func AStar(start, goal *BugNode, grid *Level) []*BugNode {
	var iteration int = 0
	var closedSet []*BugNode
	var openSet = []*BugNode{start}
	start.GCost = 0
	start.HCost = start.DistanceTo(goal)

	for len(openSet) > 0 {
		iteration++
		if iteration > SAFETY_CHECK {
			return nil
		}

		sort.Slice(openSet, func(i, j int) bool {
			return openSet[i].FCost() < openSet[j].FCost()
		})
		var current = openSet[0]

		if current == goal {
			var path []*BugNode
			for current != start {

				path = append([]*BugNode{current}, path...)
				current = current.Previous

			}
			return path
		}

		openSet = openSet[1:]
		closedSet = append(closedSet, current)

		//grab neighbours
		neighbours := GetNeighboursB(current, goal, grid)

		//update and compute costs for better paths
		for _, n := range neighbours {
			if n.IsInSet(closedSet) {
				continue
			}

			tmpscore := current.GCost + int(current.Nodetype*10) + current.DistanceTo(n)
			if tmpscore < n.GCost || !n.IsInSet(openSet) {
				n.Previous = current
				n.GCost = tmpscore
				n.HCost = goal.DistanceTo(n)

				if !n.IsInSet(openSet) {
					openSet = append(openSet, n)
				}
			}
		}
	}
	return nil
}

func GetNeighbours(node *BugNode, goal *BugNode, grid *Level) []*BugNode {
	var res []*BugNode

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			checkx := node.Location.X + dx
			checky := node.Location.Y + dy

			if 0 <= checkx && checkx < grid.Dimensions.Width && 0 <= checky && checky < grid.Dimensions.Height {
				idx := checky*grid.Dimensions.Width + checkx
				if grid.Nodes[idx].Nodetype < definitions.NodeTypeWall {
					res = append(res, grid.Nodes[idx])
				}
			}
		}
	}

	return res
}

func GetNeighboursB(node *BugNode, goal *BugNode, grid *Level) []*BugNode {
	var res []*BugNode

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			//force orthogonal moves
			c1 := dx == 0 && dy == 0
			c2 := dx == dy
			c3 := dx == -dy

			if c1 || c2 || c3 {
				continue
			}

			checkx := node.Location.X + dx
			checky := node.Location.Y + dy

			if 0 <= checkx && checkx < grid.Dimensions.Width && 0 <= checky && checky < grid.Dimensions.Height {
				idx := checky*grid.Dimensions.Width + checkx

				dz := math.Abs(float64(grid.Nodes[idx].Location.Z - node.Location.Z))
				if grid.Nodes[idx].Nodetype < definitions.NodeTypeWall && dz < 10 {
					res = append(res, grid.Nodes[idx])
				}
			}
		}
	}

	return res
}

func PrintPath(path []*BugNode) {
	fmt.Printf("Path length: %d\n", len(path))
	if len(path) > 0 {
		for i, n := range path {
			fmt.Printf("%d: (%d,%d)\n", i, n.Location.X, n.Location.Y)
		}
	}
}
