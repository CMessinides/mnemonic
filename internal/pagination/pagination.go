package pagination

type Page[T any] struct {
	Items      []T    `json:"items"`
	Page       uint64 `json:"page"`
	PageSize   uint64 `json:"pageSize"`
	TotalPages uint64 `json:"totalPages"`
}
