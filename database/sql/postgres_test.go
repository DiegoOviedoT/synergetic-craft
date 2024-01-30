//go:build integration
// +build integration

package sql_test

import (
	"testing"

	"synergetic-craft/database/sql"

	"github.com/stretchr/testify/assert"
)

func TestPostgres_Migrate(t *testing.T) {
	conf := sql.PGConfig{
		ConnURL:  "postgres://test:passwordtest@localhost:5432/dbtest?sslmode=disable",
		Database: "dbtest",
		PoolSize: 100,
	}

	conn := sql.NewPostgres(conf)
	defer conn.Close()

	err := conn.Migrate("file://migrations_test")

	assert.NoError(t, err)
}
