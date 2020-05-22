package itea_http_server

type action struct {
	Uri 		string
	Method 		string
	Controller 	string
	Action 		string
	Middleware  []string
}

type IRouter interface {
	Init(file string)
	Action() []*action
}

type RouterConstruct func() IRouter