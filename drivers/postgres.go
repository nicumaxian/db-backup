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
	spinner, _ := pterm.DefaultSpinner.Start("Killing connections")
	connStr := fmt.Sprintf("host=%v user=%v password=%v sslmode=disable", p.config.Host, p.config.Username, p.config.Password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		spinner.Fail("Failed to open connection to database")
		return err
	}

	_, err = db.Exec("select pg_terminate_backend(pid) from pg_stat_activity where datname= $1", p.config.Database)
	if err != nil {
		spinner.Fail("Failed to reset connections")
		return err
	}

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s", p.config.Database))
	if err != nil {
		spinner.Fail("Failed to recreate database")
		return err
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", p.config.Database))
	if err != nil {
		spinner.Fail("Failed to recreate database")
		return err
	}

	err = db.Close()
	if err != nil {
		spinner.Fail("Failed to close connection")
		return err
	}
	spinner.Success("Database dropped")

	spinner, _ = pterm.DefaultSpinner.Start("Restoring data")
	cmd := exec.Command("psql", "-U", p.config.Username, "-h", p.config.Host, "-f", path, p.config.Database)
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%v", p.config.Password))
	//cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%v", os.Getenv("PATH")))
	err = cmd.Run()
	if err != nil {
		spinner.Fail("Restore failed")
		return err
	}

	spinner.Success("Data restored")
	return nil
}

func createPostgresDbClient(conf configuration.DbConfiguration) (DbClient, error) {

	return PostgresDbClient{
		config: conf,
	}, nil
}
