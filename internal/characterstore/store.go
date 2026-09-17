package characterstore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type database interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
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

func (s *Store) Insert(ctx context.Context, accountId int64, character Character) (uint32, error) {
	var guid uint32
	err := s.db.QueryRow(ctx,
		`INSERT INTO
		 public.characters
   		 (account_id, name, race, class, gender, skin, face, hair_style, hair_color, facial_style, level, map_id, zone_id, position_x, position_y, position_z, orientation)
		 VALUES (@account_id, @name, @race, @class, @gender, @skin, @face, @hair_style, @hair_color, @facial_style, @level, @map_id, @zone_id, @position_x, @position_y, @position_z, @orientation)
		 RETURNING guid`,
		pgx.StrictNamedArgs{
			"account_id":   accountId,
			"name":         character.Name,
			"race":         character.Race,
			"class":        character.Class,
			"gender":       character.Gender,
			"skin":         character.Appearance.Skin,
			"face":         character.Appearance.Face,
			"hair_style":   character.Appearance.HairStyle,
			"hair_color":   character.Appearance.HairColor,
			"facial_style": character.Appearance.FacialStyle,
			"level":        character.State.Level,
			"map_id":       character.State.MapID,
			"zone_id":      character.State.ZoneID,
			"position_x":   character.State.PositionX,
			"position_y":   character.State.PositionY,
			"position_z":   character.State.PositionZ,
			"orientation":  character.State.Orientation}).Scan(&guid)
	if err != nil {
		return 0, fmt.Errorf("insert character: %w", err)
	}
	return guid, nil
}

func (s *Store) ListByAccountId(ctx context.Context, accountId int64) ([]Character, error) {
	rows, err := s.db.Query(
		ctx,
		`SELECT guid, name, race, class, gender, skin, face, hair_style, hair_color, facial_style, level, map_id, zone_id, position_x, position_y, position_z, orientation
		 FROM public.characters
		 WHERE account_id = @account_id`,
		pgx.StrictNamedArgs{"account_id": accountId})
	if err != nil {
		return nil, fmt.Errorf("error querying db: %w", err)
	}
	defer rows.Close()
	var characters []Character
	for rows.Next() {
		var character Character
		err = rows.Scan(
			&character.GUID,
			&character.Name,
			&character.Race,
			&character.Class,
			&character.Gender,
			&character.Appearance.Skin,
			&character.Appearance.Face,
			&character.Appearance.HairStyle,
			&character.Appearance.HairColor,
			&character.Appearance.FacialStyle,
			&character.State.Level,
			&character.State.MapID,
			&character.State.ZoneID,
			&character.State.PositionX,
			&character.State.PositionY,
			&character.State.PositionZ,
			&character.State.Orientation)
		if err != nil {
			return nil, fmt.Errorf("error parsing db row: %w", err)
		}
		characters = append(characters, character)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error parsing db rows: %w", err)
	}
	return characters, nil
}
