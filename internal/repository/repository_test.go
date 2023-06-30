package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
)

var (
	dbPool         *pgxpool.Pool //nolint:gochecknoglobals  // Explanation: This global variable is needed for tests
	userRepository *User         //nolint:gochecknoglobals  // Explanation: This global variable is needed for tests
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("unix:///home/entetry/.docker/desktop/docker.sock")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "14.1-alpine", []string{"POSTGRES_PASSWORD=password123"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	var dbHostAndPort string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pool.Retry(func() error {
		dbHostAndPort = resource.GetHostPort("5432/tcp")

		dbPool, err = pgxpool.Connect(ctx, fmt.Sprintf("postgresql://postgres:password123@%v/postgres", dbHostAndPort))
		if err != nil {
			return err
		}

		return dbPool.Ping(ctx)
	})
	if err != nil {
		cancel()
		log.Errorf("Could not connect to database: %s", err)
		return
	}

	userRepository = NewUserRepository(dbPool)
	if err != nil {
		log.Error(err)
		return
	}

	cmd := exec.Command("flyway",
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", dbHostAndPort),
		"migrate")

	err = cmd.Run()
	if err != nil {
		log.Errorf("Could not connect to database: %s", err)
		return
	}

	code := m.Run()
	defer os.Exit(code)

	if err = pool.Purge(resource); err != nil {
		log.Panicf("Could not purge resource: %s", err)
	}

	err = resource.Expire(1)
	if err != nil {
		log.Panic(err)
	}
}
