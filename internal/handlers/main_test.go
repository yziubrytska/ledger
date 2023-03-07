package handlers

import (
	"fmt"
	"github.com/adrianbrad/psqldocker"
	"github.com/adrianbrad/psqltest"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// psql connection parameters.
	const (
		usr           = "usr"
		password      = "pass"
		dbName        = "tst"
		containerName = "psql_docker_tests"
	)

	// run a new psql docker container.
	c, err := psqldocker.NewContainer(
		usr,
		password,
		dbName,
		psqldocker.WithContainerName(containerName),
		psqldocker.WithSql(`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE accounts(
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email VARCHAR NOT NULL,
			balance bigint
        );
		CREATE TABLE transactions(
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES accounts (id),
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			sum bigint);
		INSERT INTO accounts VALUES ('24f594ab-2207-4b45-897d-7ad948b60896', 'martin@martin.com', '140'), 
('dec1ffba-daef-4f9b-9b2a-8bebb9d1332a', 'michael@michael.com', '200');
		INSERT INTO transactions VALUES ('83435952-4c3d-4e47-bb56-e7390342e7b0', 'dec1ffba-daef-4f9b-9b2a-8bebb9d1332a',
timestamp '2015-01-10 00:51:14', 200);
`),
	)
	if err != nil {
		log.Fatalf("err while creating new psql container: %s", err)
	}

	// exit code
	var ret int

	defer func() {
		// close the psql container
		err = c.Close()
		if err != nil {
			log.Printf("err while tearing down db container: %s", err)
		}

		// exit with the code provided by executing m.Run().
		os.Exit(ret)
	}()

	// compose the psql dsn.
	dsn := fmt.Sprintf(
		"user=%s "+
			"password=%s "+
			"dbname=%s "+
			"host=localhost "+
			"port=%s "+
			"sslmode=disable",
		usr,
		password,
		dbName,
		c.Port(),
	)

	psqltest.Register(dsn)

	// run the package tests.
	ret = m.Run()
}
