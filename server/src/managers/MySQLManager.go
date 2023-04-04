package managers

import(
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"fmt"
	"os"
	"time"
)


func MySqlConnectDatabase() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionName := fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName)
	db, err := sql.Open("mysql", connectionName)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(60)
	db.SetMaxIdleConns(60)
	return db
}
