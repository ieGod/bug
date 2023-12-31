package shaders

import (
	_ "embed"
)

type ShaderType int

const (
	ShaderGlitch ShaderType = iota
)

var (
	//go:embed glitch.kage
	Glitch []byte
)

func init() {

}
