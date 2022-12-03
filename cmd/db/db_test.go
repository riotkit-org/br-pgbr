package db_test

import (
	"context"
	"github.com/riotkit-org/br-pg-simple-backup/cmd/db"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"path/filepath"
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

	path, _ := filepath.Abs("../../.build/data")
	cmd := db.NewBackupCommand(path)
	cmd.SetArgs([]string{
		"--host=127.0.0.1",
		"--port=" + port,
		"--user=anarchism",
		"--password=syndicalism",
		"--db-name=riotkit",

		// not using --db-name, in effect dumpall will be used
	})
	assert.Nil(t, cmd.Execute())
}

func createContainer() (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "quay.io/bitnami/postgresql:" + getPostgresVersion(),
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		AutoRemove:   true,
		ReaperImage:  "quay.io/testcontainers/ryuk:0.2.3",
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
