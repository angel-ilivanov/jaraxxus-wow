package worldserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
)

func (s *Server) Start(ctx context.Context, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error establishing listener %w", err)
	}
	defer listener.Close()
	slog.InfoContext(ctx, "listening for TCP connections")
	for {
		_, err := listener.Accept()
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to accept TCP connection",
				slog.String("address", listener.Addr().String()),
				slog.Any("err", err))
			continue
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	logger := slog.With(
		slog.String("component", "worldserver"))
	logger.InfoContext(ctx, "received world server connection")
	conn.Close()
}
