package admin

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/index.html
var adminPage string

//go:embed assets/assets
var embeddedAdminAssets embed.FS

func adminAssetHandler() http.Handler {
	root, _ := fs.Sub(embeddedAdminAssets, "assets/assets")
	return http.StripPrefix("/assets/", http.FileServer(http.FS(root)))
}
