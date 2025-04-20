package bookmark

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cmessinides/mnemonic/internal/tag"
	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Bookmark struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Tags      tag.Tags  `json:"tags"`
}

type BookmarkStore interface {
	Create(title string, url string, tags []string) (*Bookmark, error)
	GetByURL(url string) (*Bookmark, error)
	GetMany(limit uint, offset uint) ([]*Bookmark, error)
}

func NewSQLiteBookmarkStore(db *sql.DB) *SQLiteBookmarkStore {
	return &SQLiteBookmarkStore{db: sqlx.NewDb(db, "sqlite")}
}

type SQLiteBookmarkStore struct {
	db *sqlx.DB
}

//go:embed schema.sql
var schema string

func (bs *SQLiteBookmarkStore) Init() error {
	_, err := bs.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize bookmarks schema: %w", err)
	}

	return nil
}

func (bs *SQLiteBookmarkStore) Create(title string, url string, tags []string) (*Bookmark, error) {
	now := time.Now()
	b := new(Bookmark)

	err := bs.db.Get(b, "INSERT INTO bookmarks (title, url, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id, title, url, tags, created_at, updated_at", title, url, tag.Tags(tags), now, now)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE && strings.Contains(sqliteErr.Error(), "bookmarks.url") {
				return nil, &URLExistsError{
					URL: url,
					Err: err,
				}
			}
		}

		return nil, fmt.Errorf("failed to create bookmark: %w", err)
	}

	return b, nil
}

func (bs *SQLiteBookmarkStore) GetMany(limit uint, offset uint) ([]*Bookmark, error) {
	bookmarks := []*Bookmark{}
	err := bs.db.Select(&bookmarks, "SELECT * FROM bookmarks LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return bookmarks, fmt.Errorf("could not select bookmarks: %w", err)
	}

	return bookmarks, nil
}

func (bs *SQLiteBookmarkStore) GetByURL(url string) (*Bookmark, error) {
	bookmark := &Bookmark{}
	err := bs.db.Get(bookmark, `
        SELECT * FROM bookmarks WHERE url = ?
    `, url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{
				Field: "url",
				Value: url,
				Err:   err,
			}
		}

		return nil, fmt.Errorf("failed to read bookmark from database: %w", err)
	}

	return bookmark, nil
}
