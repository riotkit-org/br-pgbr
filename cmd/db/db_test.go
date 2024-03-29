package db_test

import (
	"bytes"
	"context"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strings"
	"testing"
)

func setupContainer() (testcontainers.Container, string) {
	container, err := createContainer()
	if err != nil {
		log.Fatal(err)
	}

	endpoint, _ := container.Endpoint(context.Background(), "")
	port := strings.Split(endpoint, ":")[1]

	return container, port
}

func TestBackupAndRestoreSingleDatabase(t *testing.T) {
	container, _ := setupContainer()
	defer container.Terminate(context.Background())
	ip, _ := container.ContainerIP(context.TODO())

	logrus.SetLevel(logrus.DebugLevel)

	// Backup first
	backupCmd, backupApp := db.NewBackupCommand(true, &bytes.Buffer{})
	backupCmd.SetArgs([]string{
		"--host=" + ip,
		"--port=5432",
		"--user=postgres",
		"--password=postgres",
		"--db-name=riotkit",
	})

	bErr := backupCmd.Execute()
	assert.Nil(t, bErr, "Expected that the command will not return error: "+string(backupApp.Output))
	out := backupApp.Output

	restoreStdinBuff := &bytes.Buffer{}
	restoreStdinBuff.Write(out)

	// Then restore
	restoreCmd, restoreApp := db.NewRestoreCommand(true, restoreStdinBuff)
	restoreCmd.SetArgs([]string{
		"--host=" + ip,
		"--port=5432",
		"--user=postgres",
		"--password=postgres",
		"--db-name=riotkit",
	})
	err := backupCmd.Execute()
	assert.Nil(t, err, "Expected that the restore will succeed. Output: "+string(restoreApp.Output))
}

func TestBackupAndRestoreAllDatabases(t *testing.T) {
	container, _ := setupContainer()
	defer container.Terminate(context.Background())
	ip, _ := container.ContainerIP(context.TODO())

	logrus.SetLevel(logrus.DebugLevel)

	// Backup first
	backupCmd, backupApp := db.NewBackupCommand(true, &bytes.Buffer{})
	backupCmd.SetArgs([]string{
		"--host=" + ip,
		"--port=5432",
		"--user=postgres",
		"--password=postgres",
		// not using --db-name, in effect dumpall will be used
	})

	bErr := backupCmd.Execute()
	assert.Nil(t, bErr, "Expected that the command will not return error: "+string(backupApp.Output))
	out := backupApp.Output

	restoreStdinBuff := &bytes.Buffer{}
	restoreStdinBuff.Write(out)

	// Then restore
	restoreCmd, restoreApp := db.NewRestoreCommand(true, restoreStdinBuff)
	restoreCmd.SetArgs([]string{
		"--host=" + ip,
		"--port=5432",
		"--user=postgres",
		"--password=postgres",
	})
	err := backupCmd.Execute()
	assert.Nil(t, err, "Expected that the restore will succeed. Output: "+string(restoreApp.Output))
}

func createContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "bitnami/postgresql:" + getPostgresVersion(),
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRESQL_DATABASE":          "riotkit",
			"POSTGRESQL_USERNAME":          "anarchism",
			"POSTGRESQL_PASSWORD":          "syndicalism",
			"POSTGRESQL_POSTGRES_PASSWORD": "postgres",
		},
	}
	pg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	return pg, err
}

func getPostgresVersion() string {
	version := os.Getenv("POSTGRES_VERSION")
	return strings.Split(version, ".")[0] // e.g. 11, 12, 13 or 14
}
