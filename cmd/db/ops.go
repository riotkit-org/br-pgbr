package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func kickOffConnectedClients(conn *pgx.Conn, dbName string) error {
	logrus.Infof("Logging off all active sessions for selected databases")

	// kick all except self :)
	sql := "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid()"

	// for selected database only
	if dbName != "" {
		sql += fmt.Sprintf(" AND pg_stat_activity.datname = '%s'", dbName)
	}

	sql += fmt.Sprintf(" AND pg_stat_activity.datname != '%s' ;", TechDatabaseName)

	logrus.Debug(sql)
	if _, err := conn.Exec(context.Background(), sql); err != nil {
		return errors.Wrap(err, "Cannot kick off active sessions")
	}
	return nil
}

// setConnectionLimit will set a connection limit to 0, so the applications will not be able to connect while database is not restored
func setMaintenanceMode(client *pgx.Conn, maintenance bool, databaseName string) error {
	var dbs []string
	var err error

	if databaseName != "" {
		dbs = []string{databaseName}
	} else {
		dbs, err = findAllDatabases(client)
		if err != nil {
			return errors.Wrap(err, "Cannot set maintenance mode")
		}
	}

	for _, dbName := range dbs {
		// do not close connection for self
		if dbName == TechDatabaseName {
			continue
		}

		logrus.Infof("Maintenance mode state=%v, for database '%s'", maintenance, dbName)

		sql := fmt.Sprintf("ALTER DATABASE %s WITH ALLOW_CONNECTIONS %v;", dbName, !maintenance)
		logrus.Debugf(sql)

		if _, err = client.Exec(context.Background(), sql); err != nil {
			return err
		}
	}
	return nil
}

func findAllDatabases(client *pgx.Conn) ([]string, error) {
	var dbs []string

	rows, err := client.Query(context.Background(), "SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		return []string{}, err
	}
	for rows.Next() {
		values, _ := rows.Values()
		dbs = append(dbs, values[0].(string))
	}
	return dbs, nil
}
