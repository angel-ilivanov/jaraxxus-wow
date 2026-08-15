package authserver

import "github.com/angel-ilivanov/jaraxxus-wow/internal/accountstore"

type Server struct {
	messageHandler *MessageHandler
}

func New(accountStore *accountstore.Store) *Server {
	return &Server{messageHandler: NewMessageHandler(accountStore)}
}
