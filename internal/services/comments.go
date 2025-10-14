package services

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/utils"
	"gorm.io/gorm"
)

type CommentService interface {
	FindAllCommentByPostId(filter dto.CommentFilterRequest, userDetail utils.Claims) (*dto.PaginationResult, error)
	CreateComment(commentBody dto.CommentRequest, userDetail utils.Claims) (*models.Comment, error)
}

type CommentServiceImpl struct {
	DB                *gorm.DB
	CommentRepository repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository, db *gorm.DB) CommentService {

	return &CommentServiceImpl{
		DB:                db,
		CommentRepository: commentRepository,
	}
}

func (s *CommentServiceImpl) FindAllCommentByPostId(filter dto.CommentFilterRequest, userDetail utils.Claims) (*dto.PaginationResult, error) {
	if _, err := s.FindDetailPostByPostId(filter.PostId); err != nil {
		return nil, err
	}

	datas, err := s.CommentRepository.FindAllCommentByPostId(filter)

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return datas, nil
}

func (s *CommentServiceImpl) FindDetailPostByPostId(id int) (*models.Post, error) {
	var post models.Post
	if err := s.DB.Where("id = ?", id).First(&post).Error; err != nil {
		return nil, exception.NewNotFoundErr("post not found")
	}

	return &post, nil

}

func (s *CommentServiceImpl) CreateComment(commentBody dto.CommentRequest, userDetail utils.Claims) (*models.Comment, error) {

	// cari dulu ada ga post yang mau di comment
	// kalau ga ada maka throw
	if _, err := s.FindDetailPostByPostId(commentBody.PostId); err != nil {
		return nil, err
	}

	convertUserId := int64(userDetail.UserId)
	convertParentId := int64(commentBody.ParentId)

	var comment models.Comment
	comment.Content = commentBody.Content
	comment.PostID = int64(commentBody.PostId)
	if convertParentId != 0 {
		comment.ParentID = &convertParentId
	}
	comment.UserID = &convertUserId

	if err := s.CommentRepository.Create(&comment); err != nil {
		return nil, err
	}

	return &comment, nil

}
