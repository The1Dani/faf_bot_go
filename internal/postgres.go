package internal

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/The1Dani/faf_bot_go/smig"
)

func migratePostgres(db *sql.DB) bool {
	log.Println("Checking if there are any migrations to apply")
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting the migration transaction:", err)
		return false
	}

	err = smig.MigratePG(tx, "./sql/")
	if err != nil {
		log.Println("Error during the migration:", err)

		err = tx.Rollback()
		if err != nil {
			log.Println("Error rolling back the migration transaction:", err)
		}

		return false
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error during committing the transaction:", err)
		return false
	}

	log.Println("All the migrations are applied")
	return true
}

func StartPostgreSQL() *sql.DB {
	pgsqlConnection, found := os.LookupEnv("FAFBOT_PGSQL_CONNECTION")
	if !found {
		log.Println("Could not find FAFBOT_PGSQL_CONNECTION variable")
		return nil
	}

	db, err := sql.Open("postgres", pgsqlConnection)
	if err != nil {
		log.Println("Could not open PostgreSQL connection:", err)
		return nil
	}

	ok := migratePostgres(db)
	if !ok {
		err := db.Close()
		if err != nil {
			log.Println("Error while closing PostgreSQL connection due to failed migration:", err)
		}
		return nil
	}

	
	
	return db
}
