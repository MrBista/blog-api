package repository

import (
	"errors"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/utils"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAllPost() ([]models.Post, error)
	FindAllPostWithPaging(filter dto.PostFilterRequest) (*dto.PaginationResult, error)
	GetDetailPost(slug string) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(slug string, data map[string]interface{}) error
	DeletePost(slug string) error
	GetDetailPostWithFilter(slug string, filter dto.PostFilterRequest) (*models.Post, error)
}

type PostRepositoryImpl struct {
	DB *gorm.DB
}

func NewPostRepository(DB *gorm.DB) PostRepository {
	return &PostRepositoryImpl{
		DB: DB,
	}
}

func (r *PostRepositoryImpl) GetAllPost() ([]models.Post, error) {
	var posts []models.Post
	tx := r.DB.Find(&posts)

	if tx.Error != nil {
		return posts, exception.NewGormDBErr(tx.Error)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) GetDetailPost(slug string) (*models.Post, error) {
	var post models.Post
	// tx := r.DB.Take(&post, "slug like ?", "%"+slug+"%")

	tx := r.DB.Where("slug = ?", slug).First(&post)

	if tx.Error != nil {
		return nil, exception.NewGormDBErr(tx.Error)
	}

	return &post, nil

}

func (r *PostRepositoryImpl) GetDetailPostWithFilter(slug string, filter dto.PostFilterRequest) (*models.Post, error) {
	var post models.Post
	// tx := r.DB.Take(&post, "slug like ?", "%"+slug+"%")

	query := r.DB.Where("slug = ?", slug)

	if filter.IncludeLike == 1 {
		query = query.Select(`
		posts.*,
		(
			SELECT COUNT(*)
			FROM likes
			WHERE likes.target_id = posts.id
				AND likes.target_type = 1
		) AS like_count
	`)
	}

	if filter.IncludeAuthor == 1 {
		query = query.Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Select("users.id", "users.name", "users.email")
		})

	}

	if err := query.First(&post).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return &post, nil

}

func (r *PostRepositoryImpl) CreatePost(post *models.Post) error {
	txRes := r.DB.Create(post)

	if txRes.Error != nil {
		return exception.NewGormDBErr(txRes.Error)
	}

	return nil
}

func (r *PostRepositoryImpl) UpdatePost(slug string, data map[string]interface{}) error {
	utils.Logger.Info("slug info: ", slug, data)
	res := r.DB.Model(&models.Post{}).Where("slug = ?", slug).Updates(data)

	if res.RowsAffected == 0 {
		return exception.NewGormDBErr(errors.New("no row affected"))
	}

	return nil
}

func (r *PostRepositoryImpl) DeletePost(slug string) error {
	rxRes := r.DB.Where("slug = ?", slug).Delete(&models.Post{})

	if rxRes.Error != nil {
		return exception.NewGormDBErr(rxRes.Error)
	}

	return nil
}

func (r *PostRepositoryImpl) FindAllPostWithPaging(filter dto.PostFilterRequest) (*dto.PaginationResult, error) {
	var posts []dto.PostResponse
	var total int64

	query := r.DB.Model(&models.Post{})

	if filter.AuthorID != 0 {
		query.Where("author_id = ?", filter.AuthorID)
	}

	if filter.CategoryID != 0 {
		query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.Title != "" {
		query.Where("title like ?", "%"+filter.Title+"%")
	}

	if filter.Status != 0 {
		query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	if filter.IncludeAuthor == 1 {
		query = query.Joins("LEFT JOIN users AS author on author.Id = posts.author_id")
	}

	if filter.IncludeCategory == 1 {
		query = query.Joins("LEFT JOIN categories as c on c.id = posts.category_id")
	}

	selectClause := []string{
		"posts.id",
		"posts.title",
		"posts.slug",
		"posts.content",
		"posts.main_image_uri",
		"posts.status",
		"posts.created_at",
		"posts.updated_at",
		"posts.author_id", // Tetap ambil author_id dari tabel post
	}

	if filter.IncludeAuthor == 1 {
		selectClause = append(selectClause,
			"author.name AS AuthorDetail_name",
			"author.email AS AuthorDetail_email",
			"author.id AS AuthorDetail_id",
		)
	}

	if filter.IncludeCategory == 1 {
		selectClause = append(selectClause,
			"c.name AS CategoryDetail_name",
			"c.id AS CategoryDetail_id",
			"c.slug AS CategoryDetail_slug",
			"c.description AS CategoryDetail_desc",
			"c.parent_id AS CategoryDetail_parentId",
		)
	}

	if filter.IncludeLike == 1 {
		// Tambahkan subquery untuk like_count
		likeSubQuery := `(SELECT COUNT(*) FROM likes WHERE likes.target_id = posts.id AND likes.target_type = 1) AS like_count`
		selectClause = append(selectClause, likeSubQuery)
	}

	query = query.Select(selectClause)

	query = applyPagination(query, filter.PaginationParams)

	if filter.Sort != "" {
		query = query.Order(filter.Sort)
	}

	if err := query.Scan(&posts).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return dto.NewPaginationResult(posts, total, filter.Page, filter.PageSize), nil

}
