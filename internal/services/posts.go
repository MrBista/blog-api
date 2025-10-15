package services

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/mapper"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/utils"
)

type PostService interface {
	FindAllPost() ([]dto.PostResponse, error)
	FindAllPostWithPaging(filter dto.PostFilterRequest) (*dto.PaginationResult, error)
	FindDetailPost(slug string) (*dto.PostResponse, error)
	FindDetailPostWitInclude(slug string, filter dto.PostFilterRequest) (*dto.PostResponse, error)
	CreatePost(reqBody *dto.CreatePostRequest, user *utils.Claims) error
	UpdatePost(reqBody *dto.UpdatePostRequest, user utils.Claims) error
	DeletePost(id string, user utils.Claims) error
	SaveFileTemp(file *multipart.FileHeader, dst string) (*dto.PostUploadResponse, error)
}

type PostServiceImpl struct {
	PostRepository     repository.PostRepository
	CategoryRepository repository.CategoryRepository
	StorageService     StorageService
}

func NewPostService(postRepostiory repository.PostRepository,
	categoryRepository repository.CategoryRepository,
	storageService StorageService,
) PostService {
	return &PostServiceImpl{
		PostRepository:     postRepostiory,
		CategoryRepository: categoryRepository,
		StorageService:     storageService,
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

func (p *PostServiceImpl) FindAllPostWithPaging(filter dto.PostFilterRequest) (*dto.PaginationResult, error) {
	posts, err := p.PostRepository.FindAllPostWithPaging(filter)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostServiceImpl) FindDetailPost(slug string) (*dto.PostResponse, error) {

	post, err := p.PostRepository.GetDetailPost(slug)
	if err != nil {
		return nil, err
	}

	postResponse := mapper.MapPostToResponse(*post)

	return &postResponse, nil

}

func (p *PostServiceImpl) FindDetailPostWitInclude(slug string, filter dto.PostFilterRequest) (*dto.PostResponse, error) {

	post, err := p.PostRepository.GetDetailPostWithFilter(slug, filter)
	if err != nil {
		return nil, err
	}

	// postResponse := mapper.MapPostToResponse(*post)

	return post, nil

}

func (p *PostServiceImpl) CreatePost(reqBody *dto.CreatePostRequest, user *utils.Claims) error {
	catId := int64(reqBody.CategoryId)

	// validasi dulu apakah ada atau ga categorynya

	if _, err := p.CategoryRepository.FindById(reqBody.CategoryId); err != nil {
		return exception.NewNotFoundErr("category not found")
	}

	// Ubah title jadi slug-friendly (huruf kecil, spasi jadi '-')
	slugBase := strings.ToLower(strings.ReplaceAll(reqBody.Title, " ", "_"))

	// Tambahkan timestamp biar unik
	slugTitle := fmt.Sprintf("%s-%d", slugBase, time.Now().UnixMilli())
	modelPost := models.Post{
		Title:        reqBody.Title,
		Slug:         slugTitle,
		CategoryID:   &catId,
		Content:      reqBody.Content,
		AuthorID:     int64(user.UserId),
		MainImageURI: &reqBody.ImgUrl,
	}
	err := p.PostRepository.CreatePost(&modelPost)

	if err != nil {
		// handle error db
		return err
	}

	reqBody.Id = int(modelPost.ID)

	return nil
}

func (p *PostServiceImpl) UpdatePost(reqBody *dto.UpdatePostRequest, user utils.Claims) error {

	postDetail, err := p.FindDetailPost(reqBody.Slug)
	if err != nil {
		return err
	}

	if postDetail.AuthorId != user.UserId {
		return exception.NewForbiddenErr("posts is not yours")
	}

	dataToUpdate := make(map[string]interface{})
	if reqBody.Title != nil {
		dataToUpdate["title"] = *reqBody.Title
	}
	if reqBody.Content != nil {
		dataToUpdate["content"] = *reqBody.Content
	}
	dataToUpdate["status"] = reqBody.Status

	if err := p.PostRepository.UpdatePost(reqBody.Slug, dataToUpdate); err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) DeletePost(slug string, user utils.Claims) error {
	postDetail, err := p.FindDetailPost(slug)
	if err != nil {
		return err
	}

	if postDetail.AuthorId != user.UserId {
		return exception.NewForbiddenErr("posts is not yours")
	}

	err = p.PostRepository.DeletePost(slug)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) SaveFileTemp(file *multipart.FileHeader, dst string) (*dto.PostUploadResponse, error) {

	uri, err := p.StorageService.SaveFile(file, dst)

	if err != nil {
		return nil, err
	}

	postAsset := models.PostAsset{
		AssetURI:    uri,
		IsTemporary: 1,
	}

	if err := p.PostRepository.SaveFilePost(postAsset); err != nil {

		if err := p.StorageService.DeleteFile(uri); err != nil {
			return nil, err
		}
		return nil, err
	}

	return &dto.PostUploadResponse{
		Url:         uri,
		IsTemporary: 1,
	}, nil
}
