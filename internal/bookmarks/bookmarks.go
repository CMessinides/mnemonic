package bookmarks

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

//go:embed schema.sql
var schema string

type BookmarkTags []string

func (bt *BookmarkTags) Scan(src any) error {
	if src == nil {
		*bt = []string{}
		return nil
	}

	var err error
	switch s := src.(type) {
	case []byte:
		err = json.Unmarshal(s, bt)
	case string:
		err = json.Unmarshal([]byte(s), bt)
	default:
		err = fmt.Errorf("cannot handle value of type %T", s)
	}

	if err != nil {
		return fmt.Errorf("unable to parse tags: %w", err)
	}

	return nil
}

func (bt *BookmarkTags) Value() (driver.Value, error) {
	data, err := json.Marshal(bt)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize tags: %w", err)
	}

	return data, nil
}

type Bookmark struct {
	ID        int64
	Title     string
	URL       string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Tags      BookmarkTags
}

type BookmarkStore interface {
	Init() error
	Create(title string, url string, tags []string) (int64, error)
	GetByURL(url string) (*Bookmark, error)
}

func NewBookmarkStore(db *sqlx.DB) BookmarkStore {
	return &PersistentBookmarkStore{db: db}
}

type PersistentBookmarkStore struct {
	db *sqlx.DB
}

func (bs *PersistentBookmarkStore) Init() error {
	_, err := bs.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize bookmarks schema: %w", err)
	}

	return nil
}

func (bs *PersistentBookmarkStore) Create(title string, url string, tags []string) (int64, error) {
	now := time.Now()
	tagsJson, err := json.Marshal(tags)
	if err != nil {
		return 0, fmt.Errorf("could not serialize tags: %w", err)
	}

	result, err := bs.db.Exec("INSERT INTO bookmarks (title, url, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", title, url, tagsJson, now, now)
	if err != nil {
		return 0, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return result.LastInsertId()
}

func (bs *PersistentBookmarkStore) GetByURL(url string) (*Bookmark, error) {
	bookmark := &Bookmark{}
	err := bs.db.Get(bookmark, `
        SELECT * FROM bookmarks WHERE url = ?
    `, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read bookmark from database: %w", err)
	}

	return bookmark, nil
}
