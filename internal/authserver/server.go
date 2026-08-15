package authserver

import "github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"

type Server struct {
	requestHandler *RequestHandler
}

func New(accountStore *accountstore.Store) *Server {
	return &Server{requestHandler: NewRequestHandler(accountStore)}
}
