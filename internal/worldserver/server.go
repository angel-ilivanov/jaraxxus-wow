package worldserver

import (
	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
)

type Server struct {
	requestHandler *RequestHandler
}

func New(store *accountstore.Store) *Server {
	return &Server{requestHandler: NewRequestHandler(store)}
}
