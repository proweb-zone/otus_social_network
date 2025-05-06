package enum

type PageFilter int

const (
	Limit PageFilter = iota
	Offset
)
