package itea_http_server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	DefaultIp			= "127.0.0.1"
	DefaultPort			= 80
	DefaultReadTimeout  = 1
	DefaultWriteTimeout = 30
)

type Server struct {
	Ctx				context.Context
	Name 			string
	Ip 				string
	Port 			int
	ReadTimeout 	int
	WriteTimeout	int
	Route 			string
	s				*http.Server
	r 				RouterConstruct
	h 				HandlerConstruct
	wg 				sync.WaitGroup
}

func (s *Server) Init() {
	s.make()
	s.handle()
	go watch(s)
	s.start()
}

func (s *Server) SetRouter(r RouterConstruct) {
	s.r = r
}

func (s *Server) SetHandler(h HandlerConstruct) {
	s.h = h
}

func (s *Server) make() {
	if s.Ip == "" {
		s.Ip = DefaultIp
	}
	if s.Port == 0 {
		s.Port = DefaultPort
	}
	s.s = &http.Server{
		Addr: fmt.Sprintf("%s:%d", s.Ip, s.Port),
		ReadTimeout: DefaultReadTimeout * time.Second,
		WriteTimeout: DefaultWriteTimeout * time.Second,
	}

	if s.ReadTimeout != 0 {
		s.s.ReadTimeout = time.Duration(s.ReadTimeout) * time.Second
	}
	if s.WriteTimeout != 0 {
		s.s.WriteTimeout = time.Duration(s.WriteTimeout) * time.Second
	}
}

func (s *Server) handle() {
	if s.r == nil {
		s.r = DefaultRouter
	}
	if s.h == nil {
		s.h = DefaultHandler
	}
	route := s.r()
	route.Init(s.Route)
	handler := s.h()
	s.s.Handler = handler.Handle(route.Action())
}

func (s *Server) start() {
	log.Println(fmt.Sprintf("=== 【Http】Server [%s] start [%s] ===", s.Name,  s.s.Addr))
	if err := s.s.ListenAndServe(); err != nil {
		log.Println(fmt.Sprintf("http server [%s] stop [%s]", s.Name, err))
	}
}

func watch(s *Server) {
	for {
		select {
		case <-	s.Ctx.Done():
			log.Println("http server stop ...")
			log.Println("wait for all http requests return ...")
			s.wg.Wait()
			err := s.s.Shutdown(s.Ctx)
			if err != nil {
				log.Println("http server shutdown error : ", err)
			}
			log.Println("http server stop success")
			return
		}
	}
}