package wp

import (
	"embed"
	"net/http"
)

//go:embed _admin
var adminFs embed.FS

func AdminPage() (string, http.Handler) {
	return "_admin", http.FileServer(http.FS(adminFs))
}
