package services

import (
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
)

type MessageService interface {
	FindOrCreateRoom(blogId, senderId, writerId int) (*models.ChatRoom, error)
	SaveMessage(message string, roomId, senderId int) (*models.ChatMessage, error)
}

type MessageServiceImpl struct {
	ChatRoomRepo    repository.ChatRoomRepository
	ChatMessageRepo repository.ChatMessageRepository
}

func NewMessageServiceImpl(chatRoomRepo repository.ChatRoomRepository, chatMessageRepo repository.ChatMessageRepository) MessageService {
	return &MessageServiceImpl{
		ChatRoomRepo:    chatRoomRepo,
		ChatMessageRepo: chatMessageRepo,
	}
}

func (s *MessageServiceImpl) FindOrCreateRoom(blogId, senderId, writerId int) (*models.ChatRoom, error) {

	findChatRoom, err := s.ChatRoomRepo.FindByBlogAndReader(blogId, senderId)
	if err != nil {
		return nil, err
	}

	if findChatRoom != nil {
		return findChatRoom, nil
	}

	var dataRoom models.ChatRoom

	dataRoom.PostId = int64(blogId)
	dataRoom.ReaderId = int64(senderId)
	dataRoom.WriterId = int64(writerId)

	if err := s.ChatRoomRepo.CreateRoom(&dataRoom); err != nil {
		return nil, err
	}

	return &dataRoom, nil
}

func (s *MessageServiceImpl) SaveMessage(message string, roomId, senderId int) (*models.ChatMessage, error) {

	var messageDTO models.ChatMessage
	messageDTO.Message = message
	messageDTO.RoomId = int64(roomId)
	messageDTO.SenderId = int64(senderId)

	if err := s.ChatMessageRepo.SaveMessage(&messageDTO); err != nil {
		return nil, err
	}

	return &messageDTO, nil
}
