package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ichtrojan/olympian"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/ichtrojan/cook/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbDriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	if dbDriver == "" {
		log.Fatal("DB_DRIVER not set in .env")
	}

	var dsn string
	var dialect olympian.Dialect

	switch dbDriver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
		dialect = olympian.MySQL()
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
		dialect = olympian.Postgres()
	case "sqlite3":
		dsn = os.Getenv("DB_DSN")
		if dsn == "" {
			dsn = "./database.db"
		}
		dialect = olympian.SQLite()
	default:
		log.Fatalf("Unsupported database driver: %s", dbDriver)
	}

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrator := olympian.NewMigrator(db, dialect)
	if err := migrator.Init(); err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	migrations := olympian.GetMigrations()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "status":
			if err := migrator.Status(migrations); err != nil {
				log.Fatalf("Failed to get status: %v", err)
			}
		case "rollback":
			if err := migrator.Rollback(migrations, 1); err != nil {
				log.Fatalf("Failed to rollback: %v", err)
			}
			fmt.Println("Rollback completed successfully")
		case "reset":
			if err := migrator.Reset(migrations); err != nil {
				log.Fatalf("Failed to reset: %v", err)
			}
			fmt.Println("Reset completed successfully")
		case "fresh":
			if err := migrator.Fresh(migrations); err != nil {
				log.Fatalf("Failed to fresh: %v", err)
			}
			fmt.Println("Fresh migration completed successfully")
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: migrate (default), status, rollback, reset, fresh")
			os.Exit(1)
		}
	} else {
		if err := migrator.Migrate(migrations); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	}
}
