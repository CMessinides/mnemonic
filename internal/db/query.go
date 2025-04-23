package db

import "io"

type QueryUpdater interface {
	UpdateQuery(query io.StringWriter, args []any) ([]any, error)
}
