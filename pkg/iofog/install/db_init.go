package install

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

func CreateControllerDatabase(host, user, password, provider, dbName, port) {


	// Create MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, password, host, port) // Fixed port type
	fmt.Printf(dsn)

	// Connect to MySQL server
	db, err := sql.Open(provider, dsn) // Fixed reference to planedb.Provider
	if err != nil {
		log.Fatalf("Failed to connect to MySQL server: %v", err)
	}
	defer db.Close()

	// Create the new database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	fmt.Printf("Database %s created successfully\n", dbName)

	// Switch to the newly created database
	_, err = db.Exec("USE `" + dbName + "`")
	if err != nil {
		log.Fatalf("Failed to switch to database: %v", err)
	}

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Assuming pkg is defined somewhere in your code
	// Read migration SQL file
	migrationSQL, err := ioutil.ReadFile(filepath.Join(dir, pkg.scriptDatabaseMigration))
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
			log.Fatalf("Failed to execute migration SQL statement: %v", err)
		}
	}

	fmt.Println("Migration SQL executed successfully")

	// Read seeding SQL file
	seedingSQL, err := ioutil.ReadFile(filepath.Join(dir, pkg.scriptDatabaseSeeder))
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
			log.Fatalf("Failed to execute seeding SQL statement: %v", err)
		}
	}

	fmt.Println("Seeding SQL executed successfully")
}
