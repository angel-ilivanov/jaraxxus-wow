package authserver

import "github.com/jackc/pgx/v5/pgxpool"

type Server struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Server {
	return &Server{db: db}
}
