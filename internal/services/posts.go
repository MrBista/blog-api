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
	//!TODO !@MrBista nanti di refactoring di pisahkan di service berbeda

	CreateReadingList(body dto.CreateReadingListRequest, userDetail *utils.Claims) error
	GetReadingLists(userDetail *utils.Claims) ([]dto.ReadingListDTO, error)
	GetReadingListByID(listID int64, userDetail *utils.Claims) (*dto.ReadingListDTO, error)
	UpdateReadingList(listID int64, body dto.UpdateReadingListRequest, userDetail *utils.Claims) error
	DeleteReadingList(listID int64, userDetail *utils.Claims) error
	CreateSavedPost(body dto.CreateSavedPostRequest, userDetail *utils.Claims) error
	GetSavedPosts(readingListID int64, userDetail *utils.Claims) ([]dto.SavedPostDTO, error)
	UpdateSavedPost(savedPostID int64, body dto.UpdateSavedPostRequest, userDetail *utils.Claims) error
	DeleteSavedPost(savedPostID int64, userDetail *utils.Claims) error
	DeleteSavedPostByPostAndList(postID, readingListID int64, userDetail *utils.Claims) error
	MarkAllAsRead(readingListID int64, userDetail *utils.Claims) error
	GetOrCreateDefaultReadingList(userDetail *utils.Claims) (*models.ReadingList, error)
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

	// cek dulu tuk user ini sudah berapa banyak buat post di bulan ini
	// kalau lebih dari 100 maka harus subscribe dulu
	// sebulan maksimal buat 100 artikel

	if count, err := p.PostRepository.CountPostByUserThisMonth(user.UserId); count > 100 || err != nil {
		return exception.NewBusnissLogicErr("You've reached limit for this month")
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

func (p *PostServiceImpl) CreateReadingList(body dto.CreateReadingListRequest, userDetail *utils.Claims) error {
	exists, err := p.PostRepository.CheckReadingListExists(int64(userDetail.UserId), body.Name)
	if err != nil {
		return err
	}
	if exists {
		return exception.NewBadRequestErr("Reading list dengan nama tersebut sudah ada")
	}

	modelReadingList := models.ReadingList{
		UserID:      int64(userDetail.UserId),
		Name:        body.Name,
		Description: body.Description,
		OrderIndex:  body.OrderIndex,
		Icon:        body.Icon,
		Color:       body.Color,
		IsDefault:   false,
	}

	err = p.PostRepository.CreateReadingList(&modelReadingList)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) GetReadingLists(userDetail *utils.Claims) ([]dto.ReadingListDTO, error) {
	data, err := p.PostRepository.GetReadingLists(int64(userDetail.UserId))
	if err != nil {
		return nil, err
	}

	// Jika belum ada reading list sama sekali, buat default
	if len(data) == 0 {
		// Buat reading list default
		defaultList := models.ReadingList{
			UserID:      int64(userDetail.UserId),
			Name:        "Baca Nanti",
			Description: nil,
			IsDefault:   true,
			OrderIndex:  0,
		}

		err = p.PostRepository.CreateReadingList(&defaultList)
		if err != nil {
			return nil, err
		}

		// Ambil ulang data setelah create
		data, err = p.PostRepository.GetReadingLists(int64(userDetail.UserId))
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return data, nil
}

func (p *PostServiceImpl) GetReadingListByID(listID int64, userDetail *utils.Claims) (*dto.ReadingListDTO, error) {
	data, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), listID)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	return data, nil
}

func (p *PostServiceImpl) UpdateReadingList(listID int64, body dto.UpdateReadingListRequest, userDetail *utils.Claims) error {
	existing, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), listID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	if body.Name != nil && *body.Name != existing.Name {
		exists, err := p.PostRepository.CheckReadingListExists(int64(userDetail.UserId), *body.Name)
		if err != nil {
			return err
		}
		if exists {
			return exception.NewBadRequestErr("Reading list dengan nama tersebut sudah ada")
		}
	}

	updates := make(map[string]interface{})
	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.Description != nil {
		updates["description"] = *body.Description
	}
	if body.Color != nil {
		updates["color"] = *body.Color
	}
	if body.Icon != nil {
		updates["icon"] = *body.Icon
	}
	if body.OrderIndex != nil {
		updates["order_index"] = *body.OrderIndex
	}

	if len(updates) == 0 {
		return exception.NewBadRequestErr("Tidak ada data yang diubah")
	}

	err = p.PostRepository.UpdateReadingList(int64(userDetail.UserId), listID, updates)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) DeleteReadingList(listID int64, userDetail *utils.Claims) error {
	existing, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), listID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	if existing.IsDefault {
		return exception.NewBadRequestErr("Tidak dapat menghapus reading list default")
	}

	err = p.PostRepository.DeleteReadingList(int64(userDetail.UserId), listID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) CreateSavedPost(body dto.CreateSavedPostRequest, userDetail *utils.Claims) error {
	// Cek apakah post ada
	post, err := p.PostRepository.GetPostById(body.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		return exception.NewNotFoundErr("Post tidak ditemukan")
	}

	// Cek apakah reading list ada dan milik user
	readingList, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), body.ReadingListID)
	if err != nil {
		return err
	}
	if readingList == nil {
		return exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	// Cek apakah post sudah disimpan di list ini
	exists, err := p.PostRepository.CheckSavedPostExists(int64(userDetail.UserId), body.PostID, body.ReadingListID)
	if err != nil {
		return err
	}
	if exists {
		return exception.NewBadRequestErr("Post sudah disimpan di reading list ini")
	}

	modelSavedPost := models.SavedPost{
		UserID:        int64(userDetail.UserId),
		PostID:        body.PostID,
		ReadingListID: body.ReadingListID,
		Notes:         body.Notes,
		IsRead:        false,
	}

	err = p.PostRepository.CreateSavedPost(&modelSavedPost)
	if err != nil {
		return err
	}

	return nil
}
func (p *PostServiceImpl) GetSavedPosts(readingListID int64, userDetail *utils.Claims) ([]dto.SavedPostDTO, error) {
	// Cek apakah reading list ada dan milik user
	readingList, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), readingListID)
	if err != nil {
		return nil, err
	}
	if readingList == nil {
		return nil, exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	data, err := p.PostRepository.GetSavedPosts(int64(userDetail.UserId), readingListID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *PostServiceImpl) UpdateSavedPost(savedPostID int64, body dto.UpdateSavedPostRequest, userDetail *utils.Claims) error {
	// Cek apakah saved post ada
	existing, err := p.PostRepository.GetSavedPostByID(int64(userDetail.UserId), savedPostID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewNotFoundErr("Saved post tidak ditemukan")
	}

	// Buat map updates
	updates := make(map[string]interface{})
	if body.Notes != nil {
		updates["notes"] = *body.Notes
	}
	if body.IsRead != nil {
		updates["is_read"] = *body.IsRead
		// Jika mark as read, set read_at
		if *body.IsRead {
			updates["read_at"] = time.Now()
		} else {
			// Jika mark as unread, hapus read_at
			updates["read_at"] = nil
		}
	}

	if len(updates) == 0 {
		return exception.NewBadRequestErr("Tidak ada data yang diubah")
	}

	err = p.PostRepository.UpdateSavedPost(int64(userDetail.UserId), savedPostID, updates)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) DeleteSavedPost(savedPostID int64, userDetail *utils.Claims) error {
	existing, err := p.PostRepository.GetSavedPostByID(int64(userDetail.UserId), savedPostID)
	if err != nil {
		return err
	}
	if existing == nil {
		return exception.NewNotFoundErr("Saved post tidak ditemukan")
	}

	err = p.PostRepository.DeleteSavedPost(int64(userDetail.UserId), savedPostID)
	if err != nil {
		return err
	}

	return nil
}
func (p *PostServiceImpl) DeleteSavedPostByPostAndList(postID, readingListID int64, userDetail *utils.Claims) error {
	// Cek apakah post sudah disimpan
	exists, err := p.PostRepository.CheckSavedPostExists(int64(userDetail.UserId), postID, readingListID)
	if err != nil {
		return err
	}
	if !exists {
		return exception.NewNotFoundErr("Saved post tidak ditemukan")
	}

	err = p.PostRepository.DeleteSavedPostByPostAndList(int64(userDetail.UserId), postID, readingListID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostServiceImpl) MarkAllAsRead(readingListID int64, userDetail *utils.Claims) error {
	readingList, err := p.PostRepository.GetReadingListByID(int64(userDetail.UserId), readingListID)
	if err != nil {
		return err
	}
	if readingList == nil {
		return exception.NewNotFoundErr("Reading list tidak ditemukan")
	}

	// Get semua saved posts yang belum dibaca
	savedPosts, err := p.PostRepository.GetSavedPosts(int64(userDetail.UserId), readingListID)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"is_read": true,
		"read_at": time.Now(),
	}

	for _, sp := range savedPosts {
		if !sp.IsRead {
			err = p.PostRepository.UpdateSavedPost(int64(userDetail.UserId), sp.ID, updates)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PostServiceImpl) GetOrCreateDefaultReadingList(userDetail *utils.Claims) (*models.ReadingList, error) {
	// Cek apakah sudah ada list default
	defaultList, err := p.PostRepository.GetDefaultReadingList(int64(userDetail.UserId))
	if err != nil {
		return nil, err
	}

	// Jika sudah ada, return
	if defaultList != nil {
		return defaultList, nil
	}

	// Jika belum ada, buat list default
	newDefaultList := models.ReadingList{
		UserID:    int64(userDetail.UserId),
		Name:      "Baca Nanti",
		IsDefault: true,
	}

	err = p.PostRepository.CreateReadingList(&newDefaultList)
	if err != nil {
		return nil, err
	}

	return &newDefaultList, nil
}
