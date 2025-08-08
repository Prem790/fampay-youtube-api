package utils

import "fmt"

type PaginatedResponse struct {
	Results  interface{} `json:"results"`
	Count    int64       `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
}

func NewPaginatedResponse(data interface{}, total int64, page, pageSize int) *PaginatedResponse {
	response := &PaginatedResponse{
		Results: data,
		Count:   total,
	}

	// Calculate next and previous URLs
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	if page < int(totalPages) {
		nextPage := page + 1
		nextURL := generatePageURL(nextPage, pageSize)
		response.Next = &nextURL
	}

	if page > 1 {
		prevPage := page - 1
		prevURL := generatePageURL(prevPage, pageSize)
		response.Previous = &prevURL
	}

	return response
}

func generatePageURL(page, pageSize int) string {
	return fmt.Sprintf("/api/videos?page=%d&page_size=%d", page, pageSize)
}
