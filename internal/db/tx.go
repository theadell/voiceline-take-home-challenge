package db

import (
	"context"
	"database/sql"
	"errors"
)

func ExecTransaction(ctx context.Context, db *sql.DB, fn func(*Queries) error) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {

		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
