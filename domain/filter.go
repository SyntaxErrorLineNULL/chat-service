package domain

// Filter chat message history filter
type Filter struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Count int64 `json:"count"`
}
