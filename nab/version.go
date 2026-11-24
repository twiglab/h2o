package nab

import (
	"fmt"
	"net/http"
)

const VERSION = "v0.0.5"

func VersionHandle(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "{%q:%q}", "version", VERSION); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
