package database

import (
	"finalProject2/infrastructure/config"

	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func handleDatabaseConnection() {
	appConfig := config.GetAppConfig()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		appConfig.DBHost, appConfig.DBPort, appConfig.DBUser, appConfig.DBPassword, appConfig.DBName,
	)

	db, err = sql.Open(appConfig.DBDialect, psqlInfo)

	if err != nil {
		log.Panic("error occured while trying to validate database arguments:", err)
	}

	err = db.Ping()

	if err != nil {
		log.Panic("error occured while trying to connect to database:", err)
	}

}

func createTables() {
	usersTables := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username text UNIQUE NOT NULL,
			email text UNIQUE NOT NULL,
			password text NOT NULL,
			age integer NOT NULL,
			created_at timestamptz DEFAULT current_timestamp,
			updated_at timestamptz DEFAULT current_timestamp
		)
	`

	photosTable := `
		CREATE TABLE IF NOT EXISTS photos (
			id SERIAL PRIMARY KEY,
			user_id integer NOT NULL,
			title text NOT NULL,
			caption text NOT NULL,
			photo_url text NOT NULL,
			created_at timestamptz DEFAULT current_timestamp,
			updated_at timestamptz DEFAULT current_timestamp,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`

	socialMediasTable := `
		CREATE TABLE IF NOT EXISTS social_medias (
			id SERIAL PRIMARY KEY,
			user_id integer NOT NULL,
			name text NOT NULL,
			social_media_url text NOT NULL,
			created_at timestamptz DEFAULT current_timestamp,
			updated_at timestamptz DEFAULT current_timestamp,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`

	commentsTable := `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			user_id integer NOT NULL,
			photo_id integer NOT NULL,
			message text NOT NULL,
			created_at timestamptz DEFAULT current_timestamp,
			updated_at timestamptz DEFAULT current_timestamp,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE, 
			FOREIGN KEY (photo_id) REFERENCES photos (id) ON DELETE CASCADE
		)
	`

	_, err := db.Exec(usersTables)
	if err != nil {
		log.Panic("error occured while trying to create order table:", err)
	}

	_, err = db.Exec(photosTable)
	if err != nil {
		log.Panic("error occured while trying to create order table:", err)
	}

	_, err = db.Exec(socialMediasTable)
	if err != nil {
		log.Panic("error occured while trying to create order table:", err)
	}

	_, err = db.Exec(commentsTable)
	if err != nil {
		log.Panic("error occured while trying to create order table:", err)
	}

}

func InitiliazeDatabase() {
	handleDatabaseConnection()
	createTables()
}

func GetDatabaseInstance() *sql.DB {
	if db == nil {
		log.Panic("database instance is still nill!!!")
	}

	return db
}
