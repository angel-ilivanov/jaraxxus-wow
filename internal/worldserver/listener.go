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
		conn, err := listener.Accept()
		if err != nil {
			slog.ErrorContext(
				ctx,
				"failed to accept TCP connection",
				slog.String("address", listener.Addr().String()),
				slog.Any("err", err))
			continue
		}
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	logger := slog.With(
		slog.String("component", "worldserver"))
	logger.InfoContext(ctx, "received world server connection")
	userSession := &session{}

	// SMS_AUTH_CHALLENGE kicks off client-worldServer communication
	authChallengeMessage := handleAuthChallenge(userSession)
	bytesWritten, err := conn.Write(authChallengeMessage.Encode())
	if err != nil {
		logger.ErrorContext(ctx, "failed to write response")
		return
	}
	logger.InfoContext(ctx, "wrote bytes to connection",
		slog.Int("bytes_written", bytesWritten))
	for {
		_, err = DecodeRequest(conn)
		if err != nil {
			slog.InfoContext(ctx, "failed to decode request", slog.Any("err", err))
			return
		}
		//response...
	}
}
