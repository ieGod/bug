package definitions

type NodeType int

const (
	NodeTypeBlank NodeType = iota
	NodeTypeGround
	NodeTypeWall
)

type NodeTile int

const (
	NodeTileBlank NodeTile = iota
	NodeTileGround
	NodeTileWallTop
	NodeTileWallTopLeft
	NodeTileWallLeft
	NodeTileWallBottomLeft
	NodeTileWallBottom
	NodeTileWallBottomRight
	NodeTileWallRight
	NodeTileWallTopRight
	NodeTileInsideTopLeft
	NodeTileInsideTopRight
)
