package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	_ "github.com/mattes/migrate/driver/postgres"
	"github.com/mattes/migrate/migrate"
	"github.com/mvader/pinlist/api/models"
	"github.com/mvader/pinlist/api/services"
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
		log.Fatal(err)
	}

	db, err := databaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	v1 := r.Group("/v1")

	services.Services{
		services.NewAccount(db),
		services.NewList(db),
		services.NewPin(db),
		services.NewTag(db),
	}.Register(v1)

	if cert != "" && key != "" {
		log.Fatal(r.RunTLS(addr, cert, key))
	} else {
		log.Fatal(r.Run(addr))
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
			log.Fatal(e)
		}
		return errors.New("error during database migration")
	}

	return nil
}
