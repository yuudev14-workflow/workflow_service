package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

var DB *sqlx.DB

// Setup DB
func SetupDB(dataSourceName string) error {
	logging.Logger.Infof("Connecting to DB... %v", dataSourceName)
	var err error
	DB, err = sqlx.Open("postgres", dataSourceName)

	if err != nil {
		logging.Logger.Errorf("error opening database: %v", err.Error())
		return fmt.Errorf("error opening database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		logging.Logger.Errorf("error connecting to database: %v %v", dataSourceName, err.Error())
		return fmt.Errorf("error connecting to database: %w", err)
	}
	return err
}
