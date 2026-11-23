// Package webui owns the frontend assets that get embedded in the binary.
// The embed directive must live next to the dist directory, which is why this
// package exists instead of embedding from cmd.
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// FS returns the frontend build rooted at dist, ready for AssetFileServerFS.
func FS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
