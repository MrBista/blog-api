package services

import (
	"fmt"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/mapper"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
)

type PostService interface {
	FindAllPost() ([]dto.PostResponse, error)
	FindDetailPost(slug string) (*dto.PostResponse, error)
	CreatePost(reqBody *dto.CreatePostRequest) error
	UpdatePost(reqBody *dto.UpdatePostRequest) error
	DeletePost(id int)
}

type PostServiceImpl struct {
	PostRepository repository.PostRepository
}

func NewPostService(postRepostiory repository.PostRepository) PostService {
	return &PostServiceImpl{
		PostRepository: postRepostiory,
	}
}

func (p *PostServiceImpl) FindAllPost() ([]dto.PostResponse, error) {
	var postResponse []dto.PostResponse

	posts, err := p.PostRepository.GetAllPost()
	if err != nil {
		return postResponse, err
	}

	postResponse = mapper.MapPostsToReponse(posts)

	return postResponse, nil
}

func (p *PostServiceImpl) FindDetailPost(slug string) (*dto.PostResponse, error) {

	post, err := p.PostRepository.GetDetailPost(slug)
	if err != nil {
		return nil, err
	}

	postResponse := mapper.MapPostToResponse(*post)

	return &postResponse, nil

}

func (p *PostServiceImpl) CreatePost(reqBody *dto.CreatePostRequest) error {
	catId := int64(reqBody.CategoryId)
	modelPost := models.Post{
		Title:      reqBody.Title,
		Slug:       reqBody.Slug,
		CategoryID: &catId,
		Content:    reqBody.Content,
	}
	err := p.PostRepository.CreatePost(&modelPost)

	if err != nil {
		// handle error db
		return fmt.Errorf("failed insert post %w", err)
	}

	reqBody.Id = int(modelPost.ID)

	return nil
}

func (p *PostServiceImpl) UpdatePost(reqBody *dto.UpdatePostRequest) error {
	panic("not implemented") // TODO: Implement

	return nil
}

func (p *PostServiceImpl) DeletePost(id int) {
	panic("not implemented") // TODO: Implement
}
