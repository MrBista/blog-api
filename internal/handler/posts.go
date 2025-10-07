package handler

type Post interface {
	GetAllPost() error
	GetDetailPost() error
	CreatePost() error
	UpdatePost() error
	DeletePost() error
}
