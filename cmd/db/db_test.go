package db_test

import (
	"bytes"
	"context"
	"github.com/riotkit-org/br-pg-simple-backup/assets"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
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
	container, port := setupContainer()
	defer container.Terminate(context.Background())

	path := assets.UnpackOrExit()
	cmd := db.NewBackupCommand(path, true, bytes.Buffer{})
	cmd.SetArgs([]string{
		"--host=127.0.0.1",
		"--port=" + port,
		"--user=anarchism",
		"--password=syndicalism",
		"--db-name=riotkit",

		// not using --db-name, in effect dumpall will be used
	})
	assert.Nil(t, cmd.Execute(), "Expected that the command will not return error")
}

func TestBackupAndRestoreAllDatabases(t *testing.T) {
	container, port := setupContainer()
	defer container.Terminate(context.Background())

	path := assets.UnpackOrExit()
	output := bytes.Buffer{}

	// Backup first
	cmd := db.NewBackupCommand(path, true, output)
	cmd.SetArgs([]string{
		"--host=127.0.0.1",
		"--port=" + port,
		"--user=anarchism",
		"--password=syndicalism",
		// not using --db-name, in effect dumpall will be used
	})

	out := new(strings.Builder)
	io.Copy(out, &output)

	println(out.String())
	assert.Nil(t, cmd.Execute(), "Expected that the command will not return error")
}

func createContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "bitnami/postgresql:" + getPostgresVersion(),
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRESQL_DATABASE": "riotkit",
			"POSTGRESQL_USERNAME": "anarchism",
			"POSTGRESQL_PASSWORD": "syndicalism",
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
