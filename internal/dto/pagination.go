package dto

import "encoding/json"

type PaginationParams struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Sort     string `json:"sort" query:"sort"`
}

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type PaginationResult struct {
	FieldName string
	Data      interface{}    `json:"data"`
	Meta      PaginationMeta `json:"meta"`
}

func NewPaginationResult(data interface{}, total int64, page, pageSize int, field string) *PaginationResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PaginationResult{
		FieldName: field,
		Data:      data,
		Meta: PaginationMeta{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}

func (p *PaginationResult) MarshalJSON() ([]byte, error) {
	field := p.FieldName
	if p.FieldName == "" {
		field = "data"
	}

	return json.Marshal(map[string]interface{}{
		field:  p.Data,
		"meta": p.Meta,
	})
}

// SetDefaults untuk set nilai default pagination
func (p *PaginationParams) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100 // Max limit
	}

}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}
