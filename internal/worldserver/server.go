package worldserver

import (
	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
	"github.com/angel-ilivanov/jaraxxus-wow/internal/characterstore"
)

type Server struct {
	requestHandler *RequestHandler
}

func New(accountStore *accountstore.Store, characterStore *characterstore.Store) *Server {
	return &Server{requestHandler: NewRequestHandler(accountStore, characterStore)}
}
