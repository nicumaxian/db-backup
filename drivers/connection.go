package drivers

import (
	"db-backup/configuration"
	"fmt"
)

type DbClient interface {
	Backup(path string) error
	Restore(path string) error
}

func CreateDbClient(conf configuration.DbConfiguration) (DbClient, error) {
	switch conf.Driver {
	case MySqlDriver:
		return createMySqlDbClient(conf)
	case PostgresDriver:
		return createPostgresDbClient(conf)
	default:
		return nil, fmt.Errorf("driver %v not implemented", conf.Driver)
	}
}
