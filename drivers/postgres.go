package drivers

import (
	"database/sql"
	"db-backup/configuration"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"os/exec"
)

const PostgresDriver = "postgres"

type PostgresDbClient struct {
	db     *sql.DB
	config configuration.DbConfiguration
}

func (p PostgresDbClient) Backup(path string) error {
	cmd := exec.Command( "pg_dump", "-U", p.config.Username, "-h", p.config.Host, p.config.Database)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", p.config.Password))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	cmd.Stdout = file
	defer file.Close()

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresDbClient) Restore(path string) error {
	panic("implement me")
}

func createPostgresDbClient(conf configuration.DbConfiguration) (DbClient, error) {
	connStr := fmt.Sprintf("host=%v user=%v password=%v", conf.Host, conf.Username, conf.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return PostgresDbClient{
		db: db,
		config: conf,
	}, nil
}
