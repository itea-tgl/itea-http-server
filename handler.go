package itea_http_server

import "net/http"

type Handler struct {

}

func (Handler) Handle([]*action) *http.ServeMux {
	mux := http.NewServeMux()
	return mux

}

func DefaultHandler() IHandler {
	return &Handler{}
}