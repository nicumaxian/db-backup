package drivers

import (
	"db-backup/configuration"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
)

const MySqlDriver = "mysql"

type MySqlDbClient struct {
	config configuration.DbConfiguration
}

func (m MySqlDbClient) Backup(path string) error {
	spinner, err := pterm.DefaultSpinner.Start("Backing up...")
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"mysqldump",
		"-u", m.config.Username,
		"-h", m.config.Host,
		"-P", m.config.Port,
		"-h", m.config.Host,
		m.config.Database,
	)
	cmd.Env = append(cmd.Env, fmt.Sprintf("MYSQL_PWD=%v", m.config.Password))

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

func (m MySqlDbClient) Restore(path string) error {
	spinner, _ := pterm.DefaultSpinner.Start("Restoring..")
	cmd := exec.Command(
		"mysql",
		"-u", m.config.Username,
		"-h", m.config.Host,
		"-P", m.config.Port,
		"-h", m.config.Host,
		m.config.Database,
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("MYSQL_PWD=%v", m.config.Password))
	inputFile, err := os.Open(path)
	if err != nil {
		spinner.Fail("couldn't open backup file")
		return err
	}
	cmd.Stdin = inputFile
	defer inputFile.Close()

	err = cmd.Run()
	if err != nil {
		spinner.Fail("Restore failed")
		return err
	}

	spinner.Success("Data restored")
	return nil
}

func (m MySqlDbClient) TestConnection() error {
	spinner, _ := pterm.DefaultSpinner.Start()
	cmd := exec.Command(
		"mysql",
		"-u", m.config.Username,
		"-h", m.config.Host,
		"-P", m.config.Port,
		"-h", m.config.Host,
		"--protocol=TCP",
		m.config.Database,
	)
	cmd.Env = append(cmd.Env, fmt.Sprintf("MYSQL_PWD=%v", m.config.Password))

	err := cmd.Run()
	if err != nil {
		spinner.Fail("connection test failed")
		return err
	}

	spinner.Success("connection test succeeded")

	return nil
}

func createMySqlDbClient(conf configuration.DbConfiguration) (DbClient, error) {
	return MySqlDbClient{
		config: conf,
	}, nil
}
