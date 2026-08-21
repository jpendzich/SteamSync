package handler

import (
	"net/http"
	"path/filepath"
)

func NewStaticHandler() http.Handler {
	return http.FileServer(http.Dir(filepath.Join(GetExecutableDir(), "static")))
}
