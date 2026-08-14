package accountstore

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("account not found")
	ErrUsernameTaken = errors.New("username already taken")
)

type Credentials struct {
	Salt     []byte
	Verifier []byte
}

type Authentication struct {
	ID int64
	Credentials
}

type database interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type Store struct {
	db database
}

func New(pool *pgxpool.Pool) *Store {
	return newStore(pool)
}

func newStore(db database) *Store {
	return &Store{db: db}
}

func (s *Store) Insert(
	ctx context.Context,
	username string,
	credentials Credentials,
) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		ctx,
		`INSERT INTO public.account (username, salt, verifier)
		 VALUES (@username, @salt, @verifier)
		 ON CONFLICT ON CONSTRAINT account_username_unique DO NOTHING
		 RETURNING id`,
		pgx.StrictNamedArgs{
			"username": strings.ToUpper(username),
			"salt":     credentials.Salt,
			"verifier": credentials.Verifier,
		},
	).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("insert account: %w", ErrUsernameTaken)
		}

		return 0, fmt.Errorf("insert account: %w", err)
	}

	return id, nil
}

func (s *Store) FindForAuthentication(
	ctx context.Context,
	username string,
) (Authentication, error) {
	var authentication Authentication
	err := s.db.QueryRow(
		ctx,
		`SELECT id, salt, verifier
		 FROM public.account
		 WHERE username = @username`,
		pgx.StrictNamedArgs{"username": strings.ToUpper(username)},
	).Scan(
		&authentication.ID,
		&authentication.Salt,
		&authentication.Verifier,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Authentication{}, fmt.Errorf("find account for authentication: %w", ErrNotFound)
		}

		return Authentication{}, fmt.Errorf("find account for authentication: %w", err)
	}

	return authentication, nil
}

// SetSessionKey stores sessionKey, or clears it when sessionKey is nil.
func (s *Store) SetSessionKey(
	ctx context.Context,
	id int64,
	sessionKey []byte,
) error {
	commandTag, err := s.db.Exec(
		ctx,
		`UPDATE public.account
		 SET session_key = @session_key,
		     updated_at = CURRENT_TIMESTAMP
		 WHERE id = @id`,
		pgx.StrictNamedArgs{
			"id":          id,
			"session_key": sessionKey,
		},
	)
	if err != nil {
		return fmt.Errorf("set account session key: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("set account session key: %w", ErrNotFound)
	}

	return nil
}
