package drivers

import "db-backup/configuration"

const MySqlDriver = "mysql"

func createMySqlDbClient(conf configuration.DbConfiguration) (DbClient, error) {
	return nil, nil
}