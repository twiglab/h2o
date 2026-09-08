package base

import (
	"net/http"

	"github.com/twiglab/h2o/nab/box/internal/export/base/restful"
)

type BaseExport interface {
	HttpListen() string
	HandleFunc(method, pattern string, handler restful.Handler)
	HandlerFunc(method, path string, handler http.HandlerFunc)
}
