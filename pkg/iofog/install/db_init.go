package install

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/datasance/potctl/pkg/util"
)

func CreateControllerDatabase(host, user, password, provider, dbName string, port int) {

	// Create MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port) 
	Verbose(dsn)

	// Connect to MySQL server
	db, err := sql.Open(provider, dsn) 
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer db.Close()

	// Create the new database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	Verbose(fmt.Sprintf("Database %s created successfully\n", dbName))

	// Switch to the newly created database
	_, err = db.Exec("USE `" + dbName + "`")
	if err != nil {
		log.Fatalf("Failed to switch to database: %v", err)
	}

	// Read migration SQL file
	migrationSQL, err := util.GetStaticFile("database/db_migration_v1.0.0.sql")
	if err != nil {
		log.Fatalf("Failed to read migration SQL file: %v", err)
	}

	// Split SQL script into individual statements
	migrationStatements := strings.Split(string(migrationSQL), ";")

	// Execute each migration SQL statement individually
	for _, statement := range migrationStatements {
		// Trim any leading/trailing whitespace
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue // Skip empty statements
		}

		// Execute the SQL statement
		_, err := db.Exec(statement)
		if err != nil {
			log.Println("Failed to execute migration SQL statement: %v", err)
		}
	}

	Verbose("Migration SQL executed successfully")

	// Read seeding SQL file
	seedingSQL, err := util.GetStaticFile("database/db_seeder_v1.0.0.sql")
	if err != nil {
		log.Fatalf("Failed to read seeding SQL file: %v", err)
	}

	// Split SQL script into individual statements
	seedStatements := strings.Split(string(seedingSQL), ";")

	// Execute each seeding SQL statement individually
	for _, statement := range seedStatements {
		// Trim any leading/trailing whitespace
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue // Skip empty statements
		}

		// Execute the SQL statement
		_, err := db.Exec(statement)
		if err != nil {
			log.Println("Failed to execute seeding SQL statement: %v", err)
		}
	}

	Verbose("Seeding SQL executed successfully")
}
