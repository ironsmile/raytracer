package shaders

import "embed"

//go:generate glslc triangle.frag -o frag.spv
//go:generate glslc triangle.vert -o vert.spv

// FS embed the vertex and fragnent shaders. Run `go generate` in order to compile
// them again.
//
//go:embed frag.spv
//go:embed vert.spv
var FS embed.FS
