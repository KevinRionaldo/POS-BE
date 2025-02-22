package paging

import (
	"math"
	"strconv"

	"POS-BE/libraries/models"
)

// SetPageLimit return current page and limit in int
func SetPageLimit(page string, limit string) (int, int) {
	newLimit, err := strconv.Atoi(limit)
	if err != nil {
		newLimit = 10
	}
	newPage, err := strconv.Atoi(page)
	if err != nil {
		newPage = 1
	}
	return newPage, newLimit
}

// GetPaginationDetail return pagination detail result that count from count data, page size, and current page of data
func GetPaginationDetail(totalRecords int64, PageSize int, CurrentPage int) models.Pagination {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(PageSize)))
	return models.Pagination{
		TotalRecords:    totalRecords,
		PageSize:        PageSize,
		CurrentPage:     CurrentPage,
		TotalPages:      totalPages,
		HasNextPage:     CurrentPage < totalPages,
		HasPreviousPage: CurrentPage > 1,
	}
}
