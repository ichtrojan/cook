package helpers

import (
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
)

type PaginateData struct {
	Page      int `json:"page"`
	Total     int `json:"total"`
	PageCount int `json:"pageCount"`
	PerPage   int `json:"perPage"`
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()

		page, _ := strconv.Atoi(q.Get("page"))

		if page == 0 {
			page = 1
		}

		perPage, _ := strconv.Atoi(q.Get("perPage"))

		switch {
		case perPage > 100:
			perPage = 100
		case perPage <= 0:
			perPage = 10
		}

		offset := (page - 1) * perPage

		return db.Offset(offset).Limit(perPage)
	}
}

func GetPagination(r *http.Request, total int) PaginateData {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if page == 0 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("perPage"))

	switch {
	case perPage > 100:
		perPage = 100
	case perPage <= 0:
		perPage = 10
	}

	pageCount := int(math.Ceil(float64(total) / float64(perPage)))

	if pageCount == 0 {
		pageCount = 1
	}

	return PaginateData{
		Page:      page,
		Total:     total,
		PageCount: pageCount,
		PerPage:   perPage,
	}
}
