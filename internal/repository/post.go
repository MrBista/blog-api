package repository

import (
	"errors"
	"time"

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
	GetPostById(id int64) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(slug string, data map[string]interface{}) error
	DeletePost(slug string) error
	GetDetailPostWithFilter(slug string, filter dto.PostFilterRequest) (*dto.PostResponse, error)

	SaveFilePost(postAssets models.PostAsset) error

	CountPostByUserThisMonth(userId int) (int64, error)

	GetReadingLists(userID int64) ([]dto.ReadingListDTO, error)
	GetReadingListByID(userID, listID int64) (*dto.ReadingListDTO, error)
	CreateReadingList(readingListModel *models.ReadingList) error
	UpdateReadingList(userID, listID int64, updates map[string]interface{}) error
	DeleteReadingList(userID, listID int64) error

	GetSavedPosts(userID, readingListID int64) ([]dto.SavedPostDTO, error)
	CreateSavedPost(savedPostModel *models.SavedPost) error
	GetSavedPostByID(userID, savedPostID int64) (*models.SavedPost, error)
	DeleteSavedPost(userID, savedPostID int64) error
	DeleteSavedPostByPostAndList(userID, postID, readingListID int64) error
	CheckSavedPostExists(userID, postID, readingListID int64) (bool, error)
	UpdateSavedPost(userID, savedPostID int64, updates map[string]interface{}) error

	CountUnreadSavedPosts(userID, readingListID int64) (int64, error)
	GetDefaultReadingList(userID int64) (*models.ReadingList, error)
	CheckReadingListExists(userID int64, name string) (bool, error)
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

func (r *PostRepositoryImpl) GetPostById(id int64) (*models.Post, error) {

	var postDetail models.Post

	if err := r.DB.Where("id = ?", id).First(&postDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundErr("saved post not found")
		}
		return nil, exception.NewGormDBErr(err)
	}

	return &postDetail, nil

}

func (r *PostRepositoryImpl) GetDetailPostWithFilter(slug string, filter dto.PostFilterRequest) (*dto.PostResponse, error) {
	var post dto.PostResponse
	// tx := r.DB.Take(&post, "slug like ?", "%"+slug+"%")

	query := r.DB.Model(&models.Post{})
	query = query.Where("slug = ?", slug)

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

	if err := query.Scan(&post).Error; err != nil {
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
	// var posts []dto.PostResponse
	posts := make([]dto.PostResponse, 0)
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

	return dto.NewPaginationResult(posts, total, filter.Page, filter.PageSize, "posts"), nil

}

func (r *PostRepositoryImpl) SaveFilePost(postAssets models.PostAsset) error {

	if err := r.DB.Create(&postAssets).Error; err != nil {
		return exception.NewGormDBErr(err)
	}

	return nil
}

func (r *PostRepositoryImpl) CountPostByUserThisMonth(userId int) (int64, error) {

	var total int64

	if err := r.DB.
		Where("author_id = ?", userId).
		Where("MONTH(created_at) = MONTH(NOW())").
		Where("YEAR(created_at) = YEAR(NOW())").
		Model(&models.Post{}).
		Count(&total).Error; err != nil {
		return total, exception.NewGormDBErr(err)
	}

	return total, nil
}

func (r *PostRepositoryImpl) GetReadingLists(userID int64) ([]dto.ReadingListDTO, error) {
	var results []dto.ReadingListDTO

	err := r.DB.
		Table("reading_lists rl").
		Select(`
			rl.id,
			rl.user_id,
			rl.name,
			rl.description,
			rl.is_default,
			rl.color,
			rl.icon,
			rl.order_index,
			rl.created_at,
			rl.updated_at,
			COUNT(sp.id) as total_posts,
			SUM(CASE WHEN sp.is_read = FALSE THEN 1 ELSE 0 END) as unread_count
		`).
		Joins("LEFT JOIN saved_posts sp ON rl.id = sp.reading_list_id").
		Where("rl.user_id = ?", userID).
		Group("rl.id").
		Order("rl.order_index, rl.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return results, nil
}

func (r *PostRepositoryImpl) GetReadingListByID(userID, listID int64) (*dto.ReadingListDTO, error) {
	var result dto.ReadingListDTO

	err := r.DB.
		Table("reading_lists rl").
		Select(`
			rl.id,
			rl.user_id,
			rl.name,
			rl.description,
			rl.is_default,
			rl.color,
			rl.icon,
			rl.order_index,
			rl.created_at,
			rl.updated_at,
			COUNT(sp.id) as total_posts,
			SUM(CASE WHEN sp.is_read = FALSE THEN 1 ELSE 0 END) as unread_count
		`).
		Joins("LEFT JOIN saved_posts sp ON rl.id = sp.reading_list_id").
		Where("rl.user_id = ? AND rl.id = ?", userID, listID).
		Group("rl.id").
		Scan(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundErr("not found reading list")
		}
		return nil, exception.NewGormDBErr(err)
	}

	if result.ID == 0 {
		return nil, nil
	}

	return &result, nil
}

func (r *PostRepositoryImpl) GetSavedPosts(userID, readingListID int64) ([]dto.SavedPostDTO, error) {
	var results []struct {
		dto.SavedPostDTO
		PostID           int64   `gorm:"column:post_id"`
		PostTitle        string  `gorm:"column:post_title"`
		PostSlug         string  `gorm:"column:post_slug"`
		PostMainImageURI *string `gorm:"column:post_main_image_uri"`
		PostAuthorName   string  `gorm:"column:post_author_name"`
		PostCategoryName *string `gorm:"column:post_category_name"`
	}

	err := r.DB.
		Table("saved_posts sp").
		Select(`
			sp.id,
			sp.user_id,
			sp.post_id,
			sp.reading_list_id,
			sp.notes,
			sp.is_read,
			sp.read_at,
			sp.created_at,
			sp.updated_at,
			p.id as post_id,
			p.title as post_title,
			p.slug as post_slug,
			p.main_image_uri as post_main_image_uri,
			u.name as post_author_name,
			c.name as post_category_name
		`).
		Joins("INNER JOIN posts p ON sp.post_id = p.id").
		Joins("INNER JOIN users u ON p.author_id = u.id").
		Joins("LEFT JOIN categories c ON p.category_id = c.id").
		Where("sp.user_id = ? AND sp.reading_list_id = ?", userID, readingListID).
		Order("sp.created_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	savedPosts := make([]dto.SavedPostDTO, len(results))
	for i, r := range results {
		savedPosts[i] = r.SavedPostDTO
		savedPosts[i].Post = &dto.SavedPostInfo{
			ID:           r.PostID,
			Title:        r.PostTitle,
			Slug:         r.PostSlug,
			MainImageURI: r.PostMainImageURI,
			AuthorName:   r.PostAuthorName,
			CategoryName: r.PostCategoryName,
		}
	}

	return savedPosts, nil
}

func (r *PostRepositoryImpl) CreateReadingList(readingListModel *models.ReadingList) error {
	if err := r.DB.Create(readingListModel).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	return nil
}

func (r *PostRepositoryImpl) UpdateReadingList(userID, listID int64, updates map[string]interface{}) error {
	result := r.DB.
		Model(&models.ReadingList{}).
		Where("id = ? AND user_id = ?", listID, userID).
		Updates(updates)

	if result.Error != nil {
		return exception.NewGormDBErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundErr("Reading list not found")
	}

	return nil
}

func (r *PostRepositoryImpl) DeleteReadingList(userID, listID int64) error {
	result := r.DB.
		Where("user_id = ? AND id = ?", userID, listID).
		Delete(&models.ReadingList{})

	if result.Error != nil {
		return exception.NewGormDBErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundErr("Reading list not found")
	}

	return nil
}
func (r *PostRepositoryImpl) CreateSavedPost(savedPostModel *models.SavedPost) error {
	if err := r.DB.Create(savedPostModel).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	return nil
}

func (r *PostRepositoryImpl) GetSavedPostByID(userID, savedPostID int64) (*models.SavedPost, error) {
	var savedPost models.SavedPost

	err := r.DB.
		Where("id = ? AND user_id = ?", savedPostID, userID).
		First(&savedPost).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundErr("saved post not found")
		}
		return nil, exception.NewGormDBErr(err)
	}

	return &savedPost, nil
}

func (r *PostRepositoryImpl) CheckSavedPostExists(userID, postID, readingListID int64) (bool, error) {
	var count int64

	err := r.DB.
		Model(&models.SavedPost{}).
		Where("user_id = ? AND post_id = ? AND reading_list_id = ?", userID, postID, readingListID).
		Count(&count).Error

	if err != nil {
		return false, exception.NewGormDBErr(err)
	}

	return count > 0, nil
}

func (r *PostRepositoryImpl) UpdateSavedPost(userID, savedPostID int64, updates map[string]interface{}) error {
	// Jika is_read = true, set read_at ke sekarang
	if isRead, ok := updates["is_read"].(bool); ok && isRead {
		if _, hasReadAt := updates["read_at"]; !hasReadAt {
			updates["read_at"] = time.Now()
		}
	}

	result := r.DB.
		Model(&models.SavedPost{}).
		Where("id = ? AND user_id = ?", savedPostID, userID).
		Updates(updates)

	if result.Error != nil {
		return exception.NewGormDBErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundErr("Saved post not found")
	}

	return nil
}

func (r *PostRepositoryImpl) DeleteSavedPost(userID, savedPostID int64) error {
	result := r.DB.
		Where("user_id = ? AND id = ?", userID, savedPostID).
		Delete(&models.SavedPost{})

	if result.Error != nil {
		return exception.NewGormDBErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundErr("Saved post not found")
	}

	return nil
}

func (r *PostRepositoryImpl) DeleteSavedPostByPostAndList(userID, postID, readingListID int64) error {
	result := r.DB.
		Where("user_id = ? AND post_id = ? AND reading_list_id = ?", userID, postID, readingListID).
		Delete(&models.SavedPost{})

	if result.Error != nil {
		return exception.NewGormDBErr(result.Error)
	}

	if result.RowsAffected == 0 {
		return exception.NewNotFoundErr("Saved post not found")
	}

	return nil
}

func (r *PostRepositoryImpl) CountSavedPostsByReadingList(readingListID int64) (int64, error) {
	var count int64

	err := r.DB.
		Model(&models.SavedPost{}).
		Where("reading_list_id = ?", readingListID).
		Count(&count).Error

	if err != nil {
		return 0, exception.NewGormDBErr(err)
	}

	return count, nil
}

func (r *PostRepositoryImpl) CountUnreadSavedPosts(userID, readingListID int64) (int64, error) {
	var count int64

	err := r.DB.
		Model(&models.SavedPost{}).
		Where("user_id = ? AND reading_list_id = ? AND is_read = ?", userID, readingListID, false).
		Count(&count).Error

	if err != nil {
		return 0, exception.NewGormDBErr(err)
	}

	return count, nil
}

func (r *PostRepositoryImpl) GetDefaultReadingList(userID int64) (*models.ReadingList, error) {
	var readingList models.ReadingList

	err := r.DB.
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&readingList).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewGormDBErr(err)
	}

	return &readingList, nil
}

func (r *PostRepositoryImpl) CheckReadingListExists(userID int64, name string) (bool, error) {
	var count int64

	err := r.DB.
		Model(&models.ReadingList{}).
		Where("user_id = ? AND name = ?", userID, name).
		Count(&count).Error

	if err != nil {
		return false, exception.NewGormDBErr(err)
	}

	return count > 0, nil
}
