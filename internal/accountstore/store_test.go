package accountstore

import (
	"context"
	"reflect"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestStoreInsertReturnsAccountID(t *testing.T) {
	t.Parallel()

	credentials := Credentials{
		Salt:     make([]byte, 32),
		Verifier: make([]byte, 32),
	}
	db := &stubDatabase{
		row: rowFunc(func(dest ...any) error {
			*(dest[0].(*int64)) = 42
			return nil
		}),
	}

	store := newStore(db)
	id, err := store.Insert(context.Background(), "ARTHAS", credentials)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if id != 42 {
		t.Fatalf("Insert() ID = %d, want 42", id)
	}
	assertArgs(t, db,
		pgx.StrictNamedArgs{
			"username": "ARTHAS",
			"salt":     credentials.Salt,
			"verifier": credentials.Verifier,
		},
	)
}

func TestStoreFindForAuthenticationReturnsCredentials(t *testing.T) {
	t.Parallel()

	want := Authentication{
		ID: 42,
		Credentials: Credentials{
			Salt:     make([]byte, 32),
			Verifier: make([]byte, 32),
		},
	}
	db := &stubDatabase{
		row: rowFunc(func(dest ...any) error {
			*(dest[0].(*int64)) = want.ID
			*(dest[1].(*[]byte)) = want.Salt
			*(dest[2].(*[]byte)) = want.Verifier
			return nil
		}),
	}

	store := newStore(db)
	authentication, err := store.FindForAuthentication(context.Background(), "ARTHAS")
	if err != nil {
		t.Fatalf("FindForAuthentication() error = %v", err)
	}
	if !reflect.DeepEqual(authentication, want) {
		t.Fatalf("FindForAuthentication() = %#v, want %#v", authentication, want)
	}
	assertArgs(t, db,
		pgx.StrictNamedArgs{"username": "ARTHAS"},
	)
}

func TestStoreSetSessionKeyUpdatesAccount(t *testing.T) {
	t.Parallel()

	sessionKey := make([]byte, 40)
	db := &stubDatabase{
		commandTag: pgconn.NewCommandTag("UPDATE 1"),
	}

	store := newStore(db)
	if err := store.SetSessionKey(context.Background(), 42, sessionKey); err != nil {
		t.Fatalf("SetSessionKey() error = %v", err)
	}
	assertArgs(t, db,
		pgx.StrictNamedArgs{
			"id":          int64(42),
			"session_key": sessionKey,
		},
	)
}

// TODO: Replace this stub with a PostgreSQL Testcontainer for integration tests.
type stubDatabase struct {
	row        pgx.Row
	commandTag pgconn.CommandTag
	args       []any
}

func (db *stubDatabase) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	db.args = args
	return db.row
}

func (db *stubDatabase) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (pgconn.CommandTag, error) {
	db.args = args
	return db.commandTag, nil
}

type rowFunc func(...any) error

func (f rowFunc) Scan(dest ...any) error {
	return f(dest...)
}

func assertArgs(t *testing.T, db *stubDatabase, want pgx.StrictNamedArgs) {
	t.Helper()

	if !reflect.DeepEqual(db.args, []any{want}) {
		t.Fatalf("arguments = %#v, want %#v", db.args, []any{want})
	}
}
