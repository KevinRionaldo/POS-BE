package models

type Pagination struct {
	TotalRecords    int64 `json:"total_records"`
	PageSize        int   `json:"page_size"`
	CurrentPage     int   `json:"current_page"`
	TotalPages      int   `json:"total_pages"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
}

type GetPluralsAPIResponse[T interface{}] struct {
	Status        string     `json:"status"`
	StatusMessage string     `json:"status_message"`
	Data          []T        `json:"data"`
	Pagination    Pagination `json:"pagination"`
}

type GetSingularAPIResponse struct {
	Status        string      `json:"status"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
}

type ErrorAPIResponse struct {
	Status        string `json:"status"`         // Status of the API response
	StatusMessage string `json:"status_message"` // Short message explaining the error
	Code          string `json:"code"`           // Error code that represents the specific error
	Details       string `json:"details"`        // More detailed explanation of the error
}
