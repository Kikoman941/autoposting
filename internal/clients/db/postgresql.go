package db

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
)

type Client interface {
	ModelContext(c context.Context, model ...interface{}) *pg.Query
}

func NewClient(ctx context.Context, dsn string) (*pg.DB, error) {
	opt, err := pg.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot parse postgres dsn:\n%s", err)
	}

	db := pg.Connect(opt)

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("cannot ping postgres db:\n%s", err)
	}

	return db, nil
}
