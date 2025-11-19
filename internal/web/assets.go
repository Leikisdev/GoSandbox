package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/* static/assets/*
var assetsFS embed.FS

// Assets returns an http.FileSystem rooted at the embedded "static" directory.
func Assets() http.FileSystem {
	sub, _ := fs.Sub(assetsFS, "static")
	return http.FS(sub)
}
