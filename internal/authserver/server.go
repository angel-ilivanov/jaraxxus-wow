package authserver

import "github.com/angel-ilivanov/wow-server/internal/accountstore"

type Server struct {
	accountStore *accountstore.Store
}

func New(accountStore *accountstore.Store) *Server {
	return &Server{accountStore: accountStore}
}
