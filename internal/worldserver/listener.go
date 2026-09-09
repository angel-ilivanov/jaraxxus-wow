package worldserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/angel-ilivanov/jaraxxus-wow/internal/worldserver/protocol"
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
	worldConnection := NewWorldConnection(conn)

	err := s.authenticateConnection(ctx, worldConnection, logger, userSession)
	if err != nil {
		logger.ErrorContext(ctx, "failed to authenticate connection",
			slog.Any("err", err))
		return
	}

	for {
		request, err := worldConnection.ReadMessage()
		if err != nil {
			logger.ErrorContext(ctx, "failed to decode request packet",
				slog.Any("err", err))
			if errors.Is(err, protocol.ErrUnknownOpcode) {
				continue
			}
			return
		}
		logRequestInfo(ctx, logger, request)

		response, err := s.requestHandler.HandleRequest(ctx, userSession, request)
		if err != nil {
			logger.ErrorContext(ctx, "failed to handle request",
				slog.Any("err", err))
			return
		}
		logResponseInfo(ctx, logger, response)

		err = worldConnection.WriteMessage(response)
		if err != nil {
			logger.ErrorContext(ctx, "failed to write response packet",
				slog.Any("err", err))
			return
		}
	}
}

func (s *Server) authenticateConnection(ctx context.Context, worldConnection *WorldConnection, logger *slog.Logger, userSession *session) error {
	// SMSG_AUTH_CHALLENGE
	challengeMessage := handleAuthChallenge(userSession)
	err := worldConnection.WriteMessage(challengeMessage)
	if err != nil {
		logger.ErrorContext(ctx, "failed to write SMSG_AUTH_CHALLENGE",
			slog.Any("err", err))
		return err
	}
	logResponseInfo(ctx, logger, challengeMessage)

	// CMSG_AUTH_SESSION
	request, err := worldConnection.ReadMessage()
	if err != nil {
		logger.ErrorContext(ctx, "failed to decode CMSG_AUTH_SESSION request",
			slog.Any("err", err))
		return err
	}
	logRequestInfo(ctx, logger, request)

	// SMSG_AUTH_PROOF
	response, err := s.requestHandler.HandleRequest(ctx, userSession, request)
	if err != nil {
		logger.ErrorContext(ctx, "failed to handle CMSG_AUTH_SESSION",
			slog.Any("err", err))
		return err
	}
	logResponseInfo(ctx, logger, response)

	err = worldConnection.EnableEncryption(userSession.SessionKey)
	if err != nil {
		logger.ErrorContext(ctx, "failed to enable encryption",
			slog.Any("err", err))
		return err
	}

	err = worldConnection.WriteMessage(response)
	if err != nil {
		logger.ErrorContext(ctx, "failed to write response packet",
			slog.Any("err", err))
		return err
	}

	return nil
}

func logRequestInfo(ctx context.Context, logger *slog.Logger, request protocol.ClientMessage) {
	logger = logger.With(
		slog.String("msg_type", "request"),
		slog.Int("opcode", int(request.Opcode())))
	switch request := request.(type) {
	case protocol.AuthSessionRequest:
		logger.InfoContext(ctx, "received CMSG_AUTH_SESSION packet",
			slog.String("username", request.Username))
	}
}
func logResponseInfo(ctx context.Context, logger *slog.Logger, response protocol.ServerMessage) {
	logger = logger.With(
		slog.String("msg_type", "response"),
		slog.Int("opcode", int(response.Opcode())))
	switch response := response.(type) {
	case protocol.AuthChallengeServerMessage:
		logger.InfoContext(ctx, "sent SMSG_AUTH_CHALLENGE packet, containing server proof",
			slog.Bool("encryption_enabled", false))
	case protocol.AuthResponse:
		logger.InfoContext(ctx, "sent SMSG_AUTH_RESPONSE packet",
			slog.Any("result", response.ResultCode),
			slog.Bool("encryption_enabled", true))
	}
}
