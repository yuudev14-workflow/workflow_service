package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
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

// func DbExecAndReturnMany[T any](DB *sqlx.DB, query string, args ...interface{}) ([]T, error) {
// 	var dest []T
// 	logging.Logger.Debug(query, args)
// 	err := DB.Select(&dest, query, args...)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return []T{}, nil
// 		}
// 		return []T{}, err
// 	}
// 	return dest, nil
// }

func DbExecAndReturnMany[T any](execer sqlx.ExtContext, query string, args ...interface{}) ([]T, error) {
	var dest []T
	query = execer.Rebind(query)
	logging.Logger.Debug(query, args)
	err := sqlx.SelectContext(context.Background(), execer, &dest, query, args...)
	if err != nil {
		logging.Logger.Warn(query, args, err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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

// transaction function
func Transact[T any](DB *sqlx.DB, fn func(*sqlx.Tx) (*T, error)) (*T, error) {
	tx, err := DB.Beginx()
	if err != nil {
		return nil, err
	}

	result, err := fn(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	commitErr := tx.Commit()

	if commitErr != nil {
		return nil, commitErr
	}

	return result, nil
}

func GenerateKeyValueQuery(payload map[string]types.Nullable[any], index int) ([]any, string) {

	var values []any
	var queries []string
	logging.Logger.Debugf("payload: %v", payload)

	for key, val := range payload {
		logging.Logger.Debugf("index: %v", index)
		logging.Logger.Debugf("key: %v", key)
		logging.Logger.Debugf("set: %v", val.Set)
		if val.Set {

			queries = append(queries, fmt.Sprintf("%v = $%v", key, index))
			index += 1
			values = append(values, val.Value)
		}
	}

	logging.Logger.Debugf("queries: %v", queries)
	logging.Logger.Debugf("values: %v", values)

	return values, strings.Join(queries, ",")
}
