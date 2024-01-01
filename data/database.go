package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Bookmark struct {
	BookmarkId string
	Title      string
	Url        string
}

type BookmarkStore struct {
	conn *sql.DB
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func CreateDBConfigFromEnv() DBConfig {
	// populate an instance of config struct using environment variables
	config := DBConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
	return config
}

func (store *BookmarkStore) InitStore(config DBConfig) error {
	// create connection string in this format: "postgresql://username:password@hostname:port/dbname?sslmode=require"
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	// defer db.Close()

	store.conn = db
	fmt.Println("Successfully connected to the database")
	return err
}

func (store BookmarkStore) StoreBookmark(bookmark Bookmark) error {
	_, err := store.conn.Exec("INSERT INTO bookmarks (bookmark_id, title, url) VALUES ($1, $2, $3)",
		bookmark.BookmarkId, bookmark.Title, bookmark.Url)

	if err != nil {
		log.Fatal("error inserting into the database: ", err)
	}
	return nil
}

func (store BookmarkStore) FindBookmark(bookmarkId string) (bool, error) {
	// query the database to see if the bookmark exists for the given bookmarkId
	rows, err := store.conn.Query("SELECT id FROM bookmarks WHERE bookmark_id = $1", bookmarkId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// iterate over the rows and return true if a row was found
	for rows.Next() {
		return true, nil
	}

	// return false if no row was found
	return false, nil
}
