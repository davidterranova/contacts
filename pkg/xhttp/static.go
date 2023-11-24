package xhttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MountStatic(router *mux.Router, path string, dir string) {
	router.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(dir))))
}
