package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/models"
	"github.com/mvader/pinlist/api/services"
	"github.com/mvader/pinlist/api/workers"
	"gopkg.in/gorp.v1"
)

var (
	conn = os.Getenv("DB_CONN")
	addr = os.Getenv("RUN_ADDR")
	cert = os.Getenv("SSL_CERT")
	key  = os.Getenv("SSL_KEY")

	migrationsDir = os.Getenv("MIGRATIONS_DIR")
)

func main() {
	if err := runMigrations(); err != nil {
		log.Err(err)
		return
	}

	db, err := databaseConnection()
	if err != nil {
		log.Err(err)
		return
	}

	r := gin.Default()
	v1 := r.Group("/v1")

	services.Services{
		services.NewAccount(db),
		services.NewList(db),
		services.NewPin(db),
		services.NewTag(db),
	}.Register(v1)

	workers.Workers{
		workers.NewSession(db),
	}.Run()

	if cert != "" && key != "" {
		log.Err(r.RunTLS(addr, cert, key))
	} else {
		log.Err(r.Run(addr))
	}
}

func databaseConnection() (*gorp.DbMap, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return models.NewDBMap(db, &gorp.PostgresDialect{}), nil
}

func runMigrations() error {
	errs, ok := migrate.UpSync(conn, migrationsDir)
	if !ok {
		for _, e := range errs {
			log.Err(e)
		}
		return errors.New("error during database migration")
	}

	return nil
}
