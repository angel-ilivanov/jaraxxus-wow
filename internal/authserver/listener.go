package authserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/authserver/protocol"
)

func (s *Server) Start(ctx context.Context, port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error establishing listener %w", err)
	}
	defer listener.Close()
	slog.Info("listening for TCP connections")
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
	slog.InfoContext(ctx, "connection received")
	session := &AuthSession{}
	logger := slog.With(
		slog.String("component", "authserver"),
		slog.String("remote_addr", conn.RemoteAddr().String()),
		slog.String("local_addr", conn.LocalAddr().String()))

	for {
		request, err := protocol.DecodeRequest(conn)
		if err != nil {
			logDecodeError(ctx, logger, err)
			return
		}
		response, err := s.requestHandler.HandleRequest(ctx, session, request)
		if err != nil {
			logger.ErrorContext(
				ctx,
				"failed to handle client request",
				slog.Any("err", err))
			return
		}
		packet, err := protocol.EncodeResponse(response)
		if err != nil {
			logger.WarnContext(ctx, "failed to encode unknown response", "err", err)
		}
		logger.Info("encoded response packet",
			slog.Int("packet_length", len(packet)))
		bytesWritten, err := conn.Write(packet)
		if err == nil {
			logger.Info("wrote bytes to connection",
				slog.Int("bytes_written", bytesWritten))
		}
	}
}

func logDecodeError(ctx context.Context, logger *slog.Logger, err error) {
	switch {
	case errors.Is(err, io.EOF):
		logger.DebugContext(ctx, "client disconnected", "err", err)
	case errors.Is(err, protocol.ErrUnknownOpcode):
		logger.WarnContext(ctx, "received unknown opcode", "err", err)
	default:
		logger.ErrorContext(ctx, "failed to decode client request", "err", err)
	}
}
