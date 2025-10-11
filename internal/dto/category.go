package dto

type CategoryResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Desc     string `json:"desc"`
	ParentId int    `json:"parentId"`
}

type CategoryRequst struct {
	Name     string `json:"name" validate:"required"`
	Desc     string `json:"desc"`
	ParentId int    `json:"parentId"`
}
