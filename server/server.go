package server

import (
	"net/http"

	"github.com/foomo/dumpster/config"
	"github.com/foomo/dumpster/server/handler"
	"github.com/julienschmidt/httprouter"
)

// Server dumpster server
type Server struct {
	Config *config.Config
	Router *httprouter.Router
}

func newServer(c *config.Config) (s *Server, err error) {
	s = &Server{
		Config: c,
		Router: httprouter.New(),
	}
	handler.SetupHandlers(s.Router)
	return s, nil
}

// Run run as a server
func Run(c *config.Config) error {
	s, err := newServer(c)
	if err != nil {
		return err
	}
	return http.ListenAndServe(c.Address, s.Router)
}
