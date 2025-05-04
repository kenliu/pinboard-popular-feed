package data

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
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

func CreateDBConfigFromEnv() (DBConfig, error) {
	required := []string{"DB_USERNAME", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	missing := []string{}
	env := map[string]string{}
	for _, key := range required {
		value := os.Getenv(key)
		if value == "" {
			missing = append(missing, key)
		}
		env[key] = value
	}
	if len(missing) > 0 {
		return DBConfig{}, fmt.Errorf("missing required environment variables: %v", missing)
	}
	config := DBConfig{
		Username: env["DB_USERNAME"],
		Password: env["DB_PASSWORD"],
		Host:     env["DB_HOST"],
		Port:     env["DB_PORT"],
		Database: env["DB_NAME"],
	}
	return config, nil
}

func (store *BookmarkStore) InitStore(config DBConfig) error {
	// create connection string in this format: "postgresql://username:password@hostname:port/dbname?sslmode=require"
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require", config.Username, config.Password, config.Host, config.Port, config.Database)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("error connecting to the database", "error", err)
		return err
	}
	// defer db.Close()

	store.conn = db
	slog.Info("successfully connected to the database")
	return err
}

func (store BookmarkStore) StoreBookmark(bookmark Bookmark) error {
	_, err := store.conn.Exec("INSERT INTO bookmarks (bookmark_id, title, url) VALUES ($1, $2, $3)",
		bookmark.BookmarkId, bookmark.Title, bookmark.Url)

	if err != nil {
		slog.Error("error inserting into the database", "error", err)
		return err
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
