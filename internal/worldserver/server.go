package worldserver

import (
	"github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"
)

type Server struct {
	store *accountstore.Store
}

func New(store *accountstore.Store) *Server {
	return &Server{store: store}
}
