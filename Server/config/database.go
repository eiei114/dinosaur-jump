package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

// NewDBConnection returns initialized bun.DB
func NewDBConnection() (*bun.DB, error) {
	user := getEnvWithDefault("MYSQL_USER", "root")
	password := getEnvWithDefault("MYSQL_PASSWORD", "dinosaur")
	host := getEnvWithDefault("MYSQL_HOST", "localhost")
	port := getEnvWithDefault("MYSQL_PORT", "3306")
	database := getEnvWithDefault("MYSQL_DATABASE", "user_database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	sqldb, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error while connecting to database: %+v", err)
		return nil, err
	}

	db := bun.NewDB(sqldb, mysqldialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, "SELECT 1"); err != nil {
		log.Printf("Can't connect to database: %+v", err)
		return nil, err
	}

	fmt.Println("DB接続")

	return db, nil
}

func getEnvWithDefault(name, def string) string {
	env := os.Getenv(name)
	if len(env) != 0 {
		return env
	}
	return def
}
