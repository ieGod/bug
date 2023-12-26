package definitions

type NodeType int

const (
	NodeTypeGround NodeType = iota
	NodeTypeWall
)

type NodeTile int

const (
	NodeTileGround NodeTile = iota
	NodeTileWallTop
	NodeTileWallTopLeft
	NodeTileWallLeft
	NodeTileWallBottomLeft
	NodeTileWallBottom
	NodeTileWallBottomRight
	NodeTileWallRight
	NodeTileWallTopRight
)
