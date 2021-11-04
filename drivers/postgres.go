package drivers

import (
	"database/sql"
	"db-backup/configuration"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

const PostgresDriver = "postgres"

type PostgresDbClient struct {
	config configuration.DbConfiguration
}

func (p PostgresDbClient) Backup(path string) error {
	spinner, err := pterm.DefaultSpinner.Start("Backing up...")
	if err != nil {
		return err
	}

	cmd := exec.Command("pg_dump", "-U", p.config.Username, "-h", p.config.Host, p.config.Database)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", p.config.Password))

	file, err := os.Create(path)
	if err != nil {
		spinner.Fail("Unable to create file for backup")
		return err
	}

	cmd.Stdout = file
	defer file.Close()

	err = cmd.Run()
	if err != nil {
		spinner.Fail("Backup failed")
		return err
	}

	spinner.Success("Done")
	return nil
}

func (p PostgresDbClient) Restore(path string) error {
	connStr := fmt.Sprintf("host=%v user=%v password=%v", p.config.Host, p.config.Username, p.config.Password)
	_, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	panic("implement me")
}

func createPostgresDbClient(conf configuration.DbConfiguration) (DbClient, error) {

	return PostgresDbClient{
		config: conf,
	}, nil
}
