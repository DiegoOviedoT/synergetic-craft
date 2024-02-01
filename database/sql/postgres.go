package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Postgres interface {
	Migrate(file string) error
	Close()
}

type PGConfig struct {
	ConnURL  string
	Database string
	PoolSize int
}

type postgres struct {
	connection *sql.DB
	config     PGConfig
}

func NewPostgres(conf PGConfig) Postgres {
	return &postgres{
		connection: newConnection(conf),
		config:     conf,
	}
}

func newConnection(conf PGConfig) *sql.DB {
	conn, err := sql.Open("postgres", conf.ConnURL)
	if err != nil {
		panic(err)
	}

	conn.SetMaxOpenConns(conf.PoolSize)

	return conn
}

func (p *postgres) Migrate(file string) error {
	driver, err := pg.WithInstance(p.connection, &pg.Config{})
	if err != nil {
		return fmt.Errorf("error creating driver instance: %+v", err)
	}

	migrations, err := migrate.NewWithDatabaseInstance(file, p.config.Database, driver)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %+v", err)
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error applying migrations: %+v", err)
	}

	return nil
}

func (p *postgres) Close() {
	if p.connection != nil {
		p.connection.Close()
	}
}
