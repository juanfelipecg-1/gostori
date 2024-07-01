package testutils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

type CloseFunc func() error

type TestingDB struct {
	Conn    *pgxpool.Pool
	closeFn CloseFunc
}

func SetupTestDb(migratePath string) (*TestingDB, CloseFunc, error) {
	var testDB *pgxpool.Pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
		return nil, nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
		return nil, nil, err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
		return nil, nil, err
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	err = resource.Expire(580)
	if err != nil {
		log.Fatalf("Resource expire err: %s", err)
		return nil, nil, err
	}

	pool.MaxWait = 30 * time.Second
	ctx := context.Background()
	if err = pool.Retry(func() error {
		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}
		testDB, err = pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			return err
		}
		return testDB.Ping(ctx)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
		return nil, nil, err
	}

	absMigratePath, err := filepath.Abs(migratePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	log.Printf("Running migrations from path: %s", absMigratePath)

	if err := applyMigrations(databaseUrl, absMigratePath); err != nil {
		return nil, nil, err
	}

	tdb := &TestingDB{
		Conn: testDB,
		closeFn: func() error {
			testDB.Close()
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("Could not purge resource: %s", err)
				return err
			}
			return nil
		},
	}

	return tdb, tdb.closeFn, nil
}

func applyMigrations(databaseUrl, migratePath string) error {
	command := fmt.Sprintf("migrate -database %s -path %s up", databaseUrl, migratePath)
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = migratePath
	return cmd.Run()
}
