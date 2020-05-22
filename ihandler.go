package itea_http_server

import "net/http"

type IHandler interface {
	Handle([]*action) *http.ServeMux
}

type HandlerConstruct func() IHandler