package drivers

import (
	"db-backup/configuration"
	"fmt"
	"log"
)

type DbClient interface {
	Backup(path string) error
	Restore(path string) error
	TestConnection() error
}

func CreateDbClient(conf configuration.DbConfiguration, logger *log.Logger) (DbClient, error) {
	switch conf.Driver {
	case MySqlDriver:
		return createMySqlDbClient(conf, logger)
	case PostgresDriver:
		return createPostgresDbClient(conf, logger)
	default:
		return nil, fmt.Errorf("driver %v not implemented", conf.Driver)
	}
}
