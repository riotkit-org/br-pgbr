package db_test

import (
	"context"
	"github.com/riotkit-org/br-pg-simple-backup/assets"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"strings"
	"testing"
)

func TestBackup(t *testing.T) {
	container, err := createContainer()
	if err != nil {
		log.Fatal(err)
	}
	defer container.Terminate(context.Background())

	endpoint, _ := container.Endpoint(context.Background(), "")
	port := strings.Split(endpoint, ":")[1]

	path := assets.UnpackOrExit()
	cmd := db.NewBackupCommand(path)
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
