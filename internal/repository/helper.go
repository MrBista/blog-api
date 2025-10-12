package repository

import (
	"github.com/MrBista/blog-api/internal/dto"
	"gorm.io/gorm"
)

func applyPagination(db *gorm.DB, params dto.PaginationParams) *gorm.DB {
	query := db.Offset(params.GetOffset()).Limit(params.PageSize)

	if params.Sort != "" {
		query = query.Order(params.Sort)
	}

	return query
}
