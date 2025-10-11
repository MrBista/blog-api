package mapper

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/models"
)

func MapPostToResponse(post models.Post) dto.PostResponse {
	postMap := dto.PostResponse{}
	postMap.ID = uint64(post.ID)
	postMap.Title = post.Title
	postMap.Content = post.Content
	// postMap.MainImageURI = *post.MainImageURI
	postMap.Slug = post.Slug
	postMap.CreatedAt = post.CreatedAt
	postMap.UpdatedAt = post.UpdatedAt
	if post.MainImageURI != nil {
		postMap.MainImageURI = *post.MainImageURI
	}
	postMap.AuthorId = int(post.AuthorID)

	return postMap
}

func MapPostsToReponse(posts []models.Post) []dto.PostResponse {
	// var postResponse []dto.PostResponse
	postResponse := make([]dto.PostResponse, 0, len(posts))

	for _, post := range posts {
		postResponse = append(postResponse, MapPostToResponse(post))
	}

	return postResponse
}
