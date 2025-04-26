package bookmark

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cmessinides/mnemonic/internal/pagination"
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
	Archived  bool      `json:"archived" db:"archived"`
	Tags      tag.Tags  `json:"tags"`
}

type BookmarkPatch struct {
	ID       int64
	Title    *string
	URL      *string
	Archived *bool
	Tags     tag.Tags
}

type BookmarkStore interface {
	Create(title string, url string, tags []string) (*Bookmark, error)
	Update(patch BookmarkPatch) error
	Get(id int64) (*Bookmark, error)
	GetByURL(url string) (*Bookmark, error)
	GetPage(page uint64, pageSize uint64) (*pagination.Page[*Bookmark], error)
	Delete(id int64) error
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

func (bs *SQLiteBookmarkStore) Update(patch BookmarkPatch) error {
	now := time.Now()

	args := []any{}
	query := &strings.Builder{}
	query.WriteString("UPDATE bookmarks SET ")

	if patch.Title != nil {
		args = append(args, *patch.Title)
		query.WriteString("title = ?, ")
	}

	if patch.URL != nil {
		args = append(args, *patch.URL)
		query.WriteString("url = ?, ")
	}

	if patch.Archived != nil {
		if *patch.Archived {
			args = append(args, now)
		} else {
			args = append(args, nil)
		}
		query.WriteString("archived_at = ?, ")
	}

	if patch.Tags != nil {
		args = append(args, patch.Tags)
		query.WriteString("tags = ?, ")
	}

	if len(args) == 0 {
		// nothing to update
		return nil
	}

	query.WriteString("updated_at = ? WHERE id = ?")
	args = append(args, now, patch.ID)

	result, err := bs.db.Exec(query.String(), args...)
	if err != nil {
		if isDuplicateUrl(err) {
			return &URLExistsError{
				URL: *patch.URL,
				Err: err,
			}
		} else {
			return fmt.Errorf("failed to update bookmark: %w", err)
		}
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not update bookmark: %w", err)
	}

	if n == 0 {
		return &NotFoundError{
			Field: "id",
			Value: patch.ID,
		}
	}

	return nil
}

func (bs *SQLiteBookmarkStore) GetPage(page uint64, pageSize uint64) (*pagination.Page[*Bookmark], error) {
	bookmarks := []*Bookmark{}

	limit := pageSize
	offset := (page - 1) * pageSize
	err := bs.db.Select(&bookmarks, "SELECT * FROM active_bookmarks LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not select bookmarks: %w", err)
	}

	var total uint64
	err = bs.db.Get(&total, "SELECT COUNT(1) FROM active_bookmarks")
	if err != nil {
		return nil, fmt.Errorf("could not select bookmark total: %w", err)
	}

	return &pagination.Page[*Bookmark]{
		Items:      bookmarks,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: max(total/pageSize, 1),
	}, nil
}

func (bs *SQLiteBookmarkStore) Get(id int64) (*Bookmark, error) {
	bookmark := &Bookmark{}
	err := bs.db.Get(bookmark, `
        SELECT * FROM all_bookmarks WHERE id = ?
    `, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &NotFoundError{
				Field: "id",
				Value: id,
				Err:   err,
			}
		}

		return nil, fmt.Errorf("failed to read bookmark from database: %w", err)
	}

	return bookmark, nil
}

func (bs *SQLiteBookmarkStore) GetByURL(url string) (*Bookmark, error) {
	bookmark := &Bookmark{}
	err := bs.db.Get(bookmark, `
        SELECT * FROM all_bookmarks WHERE url = ?
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

func (bs *SQLiteBookmarkStore) Delete(id int64) error {
	result, err := bs.db.Exec(`DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete bookmark from database: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete bookmark from database: %w", err)
	}

	if n == 0 {
		return &NotFoundError{Field: "id", Value: id}
	}

	return nil
}

func isDuplicateUrl(err error) bool {
	var sqliteErr *sqlite.Error

	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE && strings.Contains(sqliteErr.Error(), "bookmarks.url")
	}

	return false
}
