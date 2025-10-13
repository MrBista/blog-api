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
	datas, err := s.CommentRepository.FindAllCommentByPostId(filter)

	if err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return datas, nil
}

func (s *CommentServiceImpl) CreateComment(commentBody dto.CommentRequest, userDetail utils.Claims) (*models.Comment, error) {

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
