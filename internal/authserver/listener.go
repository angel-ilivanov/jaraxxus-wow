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
	session := &AuthSession{}
	logger := slog.With(
		slog.String("component", "authserver"),
		//slog.String("remote_addr", conn.RemoteAddr().String()),
		//slog.String("local_addr", conn.LocalAddr().String())
	)

	slog.InfoContext(ctx, "connection received")

	for {
		err := s.processRequest(ctx, logger, session, conn)
		if err != nil {
			return
		}
	}
}

func (s *Server) processRequest(ctx context.Context, logger *slog.Logger, session *AuthSession, conn net.Conn) error {
	request, err := protocol.DecodeRequest(conn)
	if err != nil {
		logDecodeError(ctx, logger, err)
		return fmt.Errorf("failed to decode request: %w", err)
	}
	logRequestInfo(ctx, logger, request)

	response, err := s.requestHandler.HandleRequest(ctx, session, request)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"failed to handle client request",
			slog.Any("err", err))
		return fmt.Errorf("failed to handle client request: %w", err)
	}
	logResponseInfo(ctx, logger, response)

	packet, err := protocol.EncodeResponse(response)
	if err != nil {
		logger.WarnContext(ctx, "failed to encode unknown response", "err", err)
		return fmt.Errorf("failed to encode unknown response: %w", err)
	}
	logger.DebugContext(ctx, "encoded response packet",
		slog.Int("packet_length", len(packet)))

	bytesWritten, err := conn.Write(packet)
	if err == nil {
		logger.DebugContext(ctx, "wrote bytes to connection",
			slog.Int("bytes_written", bytesWritten))
	}
	if bytesWritten != len(packet) {
		logger.ErrorContext(ctx, "failed to write full packet to connection")
		return fmt.Errorf("write response: wrote %d of %d bytes: %w",
			bytesWritten, len(packet), io.ErrShortWrite)
	}
	return nil
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

func logRequestInfo(ctx context.Context, logger *slog.Logger, request protocol.Request) {
	logger = logger.With(slog.String("msg_type", "request"))
	switch request := request.(type) {
	case protocol.LogonChallengeRequest:
		logger.InfoContext(ctx, "received logon challenge request",
			slog.String("username", request.AccountName))
	case protocol.LogonProofRequest:
		logger.InfoContext(ctx, "received logon proof request")
	case protocol.RealmListRequest:
		logger.InfoContext(ctx, "received realmlist request")
	}
}

func logResponseInfo(ctx context.Context, logger *slog.Logger, response protocol.Response) {
	logger = logger.With(slog.String("msg_type", "response"))
	switch response := response.(type) {
	case protocol.LogonChallengeResponse:
		logger.InfoContext(ctx, "created logon challenge response",
			slog.Int("result", int(response.Result)))
	case protocol.LogonProofResponse:
		logger.InfoContext(ctx, "created logon proof response",
			slog.Int("result", int(response.Result)))
	case protocol.RealmListResponse:
		logger.InfoContext(ctx, "created realmlist response",
			slog.Int("num_chars", int(response.NumChars)))
	}
}
