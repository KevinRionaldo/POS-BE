package apiResponse

import (
	"errors"
	"math"

	"gorm.io/gorm"
)

type pagination struct {
	TotalRecords    int64 `json:"total_records"`
	PageSize        int   `json:"page_size"`
	CurrentPage     int   `json:"current_page"`
	TotalPages      int   `json:"total_pages"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
}

type getPluralsAPIResponse[T interface{}] struct {
	Status        string     `json:"status"`
	StatusMessage string     `json:"status_message"`
	Data          []T        `json:"data"`
	Pagination    pagination `json:"pagination"`
}

type getSingularAPIResponse struct {
	Status        string      `json:"status"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
}

type getErrorResponse struct {
	Status        string `json:"status"`         // Status of the API response
	StatusMessage string `json:"status_message"` // Short message explaining the error
	Code          string `json:"code"`           // Error code that represents the specific error
	Details       error  `json:"details"`        // More detailed explanation of the error
}

func gormErrorTracking(dbError error) string {
	if errors.Is(dbError, gorm.ErrDuplicatedKey) {
		return "Data already exists"
	}
	if errors.Is(dbError, gorm.ErrForeignKeyViolated) {
		return "Invalid data request: Foreign key constraint violated"
	}
	if errors.Is(dbError, gorm.ErrRecordNotFound) {
		return "Data not found"
	}
	return dbError.Error()
}

func GetPaginationDetail(totalRecords int64, PageSize int, CurrentPage int) pagination {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(PageSize)))
	return pagination{
		TotalRecords:    totalRecords,
		PageSize:        PageSize,
		CurrentPage:     CurrentPage,
		TotalPages:      totalPages,
		HasNextPage:     CurrentPage < totalPages,
		HasPreviousPage: CurrentPage > 1,
	}
}

func SuccessSingularResponse(data interface{}) getSingularAPIResponse {
	return getSingularAPIResponse{
		Status:        "success",
		StatusMessage: "Data retrieved successfully",
		Data:          data,
		// Pagination:    paging.GetPaginationDetail(totalData, limitData, currentPage),
	}
}

func SuccessPluralResponse[T interface{}](data []T, totalData int64, limitData int, currentPage int) getPluralsAPIResponse[T] {
	return getPluralsAPIResponse[T]{
		Status:        "success",
		StatusMessage: "Data retrieved successfully",
		Data:          data,
		Pagination:    GetPaginationDetail(totalData, limitData, currentPage),
	}
}

func GeneralErrorResponse(errorMessage error) getErrorResponse {
	return getErrorResponse{
		Status:        "error",
		StatusMessage: errorMessage.Error(),
		Details:       errorMessage,
	}
}

func DBErrorResponse(errorMessage error) getErrorResponse {
	return getErrorResponse{
		Status:        "error",
		StatusMessage: gormErrorTracking(errorMessage),
		Details:       errorMessage,
	}
}
