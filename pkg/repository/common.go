package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

func DbExecAndReturnOne[T any](DB *sqlx.DB, query string, args ...interface{}) (*T, error) {
	var dest T
	logging.Logger.Debug(query, args)
	err := DB.Get(&dest, query, args...)
	if err != nil {
		logging.Logger.Warn(err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dest, nil
}

func DbExecAndReturnMany[T any](DB *sqlx.DB, query string, args ...interface{}) ([]T, error) {
	var dest []T
	logging.Logger.Debug(query, args)
	err := DB.Select(&dest, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []T{}, nil
		}
		return []T{}, err
	}
	return dest, nil
}

func DbSelectOne[T any](DB *sqlx.DB, query string, args ...interface{}) (*T, error) {
	var dest T
	logging.Logger.Debug(query, args)
	err := DB.Get(&dest, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dest, nil
}

func Transact(DB *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := DB.Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
