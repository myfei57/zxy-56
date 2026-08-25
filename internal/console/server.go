package console

import (
	"errors"
	"net/http"
)

type Server struct {
	router http.Handler
}

func NewServer(deps Deps) *Server {
	return &Server{router: buildRouter(deps)}
}

func (s *Server) Handler() http.Handler {
	return s.router
}

var errPanic = errors.New("internal server error")
