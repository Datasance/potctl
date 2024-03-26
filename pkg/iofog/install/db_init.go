package install

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/datasance/potctl/pkg/iofog"
	"github.com/datasance/potctl/pkg/util"
	_ "github.com/go-sql-driver/mysql"
)


type database struct {
	databaseName string
	provider     string
	host         string
	user         string
	password     string
	port         int
}

var ctrl struct {
	db database
}


func CreateControllerDatabase() {
	// MySQL connection parameters
	username := ctrl.db.user
	password := ctrl.db.password
	host := ctrl.db.host
	port := ctrl.db.port
	dbName := ctrl.db.databaseName

	// Create MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	// Connect to MySQL server
	db, err := sql.Open(ctrl.db.provider , dsn) //db.config
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