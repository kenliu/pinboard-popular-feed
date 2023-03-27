package data

import (
	"encoding/json"
	"io"
	"os"
)

type Bookmark struct {
	Id    string
	Title string
	Url   string
}

type BookmarkStore struct {
	bookmarks map[string]Bookmark
}

type DBConfig struct {
}

func Init() BookmarkStore {
	return BookmarkStore{
		bookmarks: make(map[string]Bookmark),
	}
}

func (store BookmarkStore) InitStore(config DBConfig) error {
	err := loadDB(&store.bookmarks)
	if err != nil {
		panic(err)
	}
	return nil
}

func (store BookmarkStore) StoreBookmark(bookmark Bookmark) error {
	store.bookmarks[bookmark.Id] = bookmark

	j, err := json.Marshal(store.bookmarks)
	if err != nil {
		panic(err)
	}

	_ = os.WriteFile("database.json", j, 0644)
	return nil
}

func (store BookmarkStore) FindBookmark(id string) (bool, error) {
	_, found := store.bookmarks[id]
	if found {
		return true, nil
	} else {
		return false, nil
	}
}

func (store BookmarkStore) UpdatePostedBookmark() error {
	return nil
}

func loadDB(db *map[string]Bookmark) error {
	jsonFile, err := os.Open("database.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &db)
	return nil
}
