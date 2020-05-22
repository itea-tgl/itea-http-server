package itea_http_server

import (
	"fmt"
	"testing"
)

var router *Router

func Test_Init(t *testing.T) {
	router = &Router{}
	router.Init("/test_router.yml")
}

func Test_Action(t *testing.T) {
	actions := router.Action()
	fmt.Println(actions[0])
	fmt.Println(actions[1])
}