package webdelivery

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/index.html
var dataPage string

//go:embed assets/assets
var embeddedDataAssets embed.FS

func dataAssetHandler() http.Handler {
	root, _ := fs.Sub(embeddedDataAssets, "assets/assets")
	return http.StripPrefix("/assets/", http.FileServer(http.FS(root)))
}
