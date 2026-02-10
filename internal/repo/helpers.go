package repo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func uuidToPgtype(u uuid.UUID) pgtype.UUID {
	var pg pgtype.UUID
	copy(pg.Bytes[:], u[:])
	pg.Valid = true
	return pg
}

func pgtypeToUUID(pg pgtype.UUID) (uuid.UUID, error) {
	if !pg.Valid {
		return uuid.Nil, fmt.Errorf("invalid UUID")
	}
	return uuid.FromBytes(pg.Bytes[:])
}

func timeToPGDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  t,
		Valid: true,
	}
}

func pgTimestamptzToTime(ts pgtype.Timestamptz) (time.Time, error) {
	if !ts.Valid {
		return time.Time{}, fmt.Errorf("invalid timestamptz")
	}
	return ts.Time, nil
}

func pgDateToTime(d pgtype.Date) (*time.Time, error) {
	if !d.Valid {
		return nil, nil
	}
	return &d.Time, nil
}

func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func derefInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
