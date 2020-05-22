package itea_http_server

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type routeConf struct {
	Groups 		[]groupConf				`yaml:"groups"`
	Action 		map[string]actionConf	`yaml:"action"`
}

type groupConf struct {
	Name 		string					`yaml:"name"`
	Prefix 		string					`yaml:"prefix"`
	Middleware 	string					`yaml:"middleware"`
}

type actionConf struct {
	Method 		string					`yaml:"method"`
	Uses 		string					`yaml:"uses"`
	Middleware 	string					`yaml:"middleware"`
	Group 		string					`yaml:"group"`
}

type Router struct {
	conf 	routeConf
}

// init router by config file
func (r *Router) Init(file string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile(path + file)
	if err != nil {
		panic(err)
	}

	var routeConf routeConf
	err = yaml.Unmarshal(dat, &routeConf)
	if err != nil {
		panic(err)
	}
	r.conf = routeConf
}

// get actions all by config
func (r *Router) Action() []*action {
	groups := groupMap(r.conf.Groups)
	var actions []*action
	var wg sync.WaitGroup
	for u, a := range r.conf.Action {
		wg.Add(1)
		go func(u string, a actionConf) {
			defer wg.Done()
			if action := extract(u, a, groups); action != nil {
				actions = append(actions, action)
			}
		}(u, a)
	}
	wg.Wait()
	return actions
}

// make group list to map[name]group
func groupMap(g []groupConf) map[string]groupConf {
	gc := make(map[string]groupConf)
	for _, i := range g {
		gc[i.Name] = i
	}
	return gc
}

// extract action info
func extract(u string, a actionConf, g map[string]groupConf) *action {
	if a.Uses == "" {
		return nil
	}

	pa := strings.Split(a.Uses, "@")
	if len(pa) != 2 {
		return nil
	}

	m, c, f := "get", pa[0], pa[1]
	var ml []string

	ua := strings.Split(u, " ")
	switch len(ua) {
	case 1:
		u = ua[0]
		break
	case 2:
		m, u = ua[0], ua[1]
	default:
		return nil
	}

	if a.Method != "" {
		m = a.Method
	}

	if a.Group != "" {
		pre, mid := group(a.Group, g)
		u = pre + u
		ml = append(ml, mid...)
	}

	if a.Middleware != "" {
		ml = append(ml, strings.Split(a.Middleware, "|")...)
	}

	return &action{
		Uri:        u,
		Method:     m,
		Controller: c,
		Action:     f,
		Middleware: ml,
	}
}

// extract group info
// get prefix and middleware of group config
func group(sg string, g map[string]groupConf) (pre string, mid []string) {
	for _, n := range strings.Split(sg, "|") {
		if i, ok := g[n]; ok {
			if i.Prefix != "" {
				pre = pre + i.Prefix
			}
			if i.Middleware != "" {
				mid = append(mid, strings.Split(i.Middleware, "|")...)
			}
		}
	}
	return
}

func DefaultRouter() IRouter {
	return &Router{}
}
