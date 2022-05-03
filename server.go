package server

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/valyala/fasthttp"
)

type Server struct {
	listen string
	Routes fasthttp.RequestHandler
	server fasthttp.Server

	errChan    chan error
	exitChan   chan os.Signal
	reloadChan chan os.Signal
}

func NewGCP(routes fasthttp.RequestHandler) (*Server, error) {
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		return nil, errors.New("PORT env var is not set")
	}
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, errors.New("PORT env var is not a valid number")
	}
	return New(port, routes)
}

func New(port int, routes fasthttp.RequestHandler) (*Server, error) {
	if port <= 0 || port >= 65535 {
		return nil, errors.New("invalid port number")
	}

	s := &Server{
		listen: fmt.Sprintf(":%d", port),
		Routes: routes,

		errChan:    make(chan error, 1),
		exitChan:   make(chan os.Signal, 1),
		reloadChan: make(chan os.Signal, 1),
	}

	signal.Notify(s.exitChan, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(s.reloadChan, syscall.SIGHUP)

	return s, nil
}

func (s *Server) Run() error {
	s.start()
	for {
		select {
		case <-s.reloadChan: // reload config
			s.Stop()
			s.start()
		case <-s.exitChan: // shutdown servers
			s.Stop()
			return nil
		case err := <-s.errChan:
			return err
		}
	}
}

func (s *Server) start() {
	s.server = fasthttp.Server{
		Handler: s.Routes,
	}
	go func() {
		err := s.server.ListenAndServe(s.listen)
		if err != nil {
			s.errChan <- err
		}
	}()
}

func (s *Server) Stop() error {
	return s.server.Shutdown()
}
