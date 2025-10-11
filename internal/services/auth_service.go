package services

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/utils"
)

type AuthService interface {
	LoginUser(reqLogin dto.LoginRequest) (dto.LoginResponse, error)
	RegisterUser(reqRegister dto.RegisterRequest) error
	ConfirmOtp() error
}

type AuthServiceImpl struct {
	UserRepo repository.UserRepository
}

func NewAutService(userRepo repository.UserRepository) AuthService {
	return &AuthServiceImpl{
		UserRepo: userRepo,
	}
}

func (s *AuthServiceImpl) FindUserByIdentifier(identifier string) (*models.User, error) {
	user, err := s.UserRepo.FindByIdentifier(identifier)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) FindByEmailOrUsername(email, username string) (*models.User, error) {
	user, err := s.UserRepo.FindByEmailOrUsername(email, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthServiceImpl) LoginUser(reqLogin dto.LoginRequest) (dto.LoginResponse, error) {
	/*
		1. cari user ada nggak
		2. kalau user ga ada maka throw username/password not valid
		3. kalau ada maka cek passwordnya match apa nggak
		4. kalau ga match maka throw username/password not valid
		5. generate jwt by userid dan buat expired 7 hari tuk refresh token dan 1 hari tuk access token
		6.

	*/

	var responseLogin dto.LoginResponse

	user, err := s.UserRepo.FindByIdentifier(reqLogin.Identifier)

	if err != nil {
		return responseLogin, err
	}

	if err := utils.ComparePassword(reqLogin.Password, user.Password); err != nil {
		return responseLogin, err
	}

	jwtService := utils.GetJwtService()

	token, err := jwtService.CreateAccessToken(int(user.ID), user.Role)

	if err != nil {
		return responseLogin, err
	}

	responseLogin.AccessToken = token
	responseLogin.TypeToken = "access_token"

	return responseLogin, nil

}

func (s *AuthServiceImpl) RegisterUser(reqRegister dto.RegisterRequest) error {
	/*
		1. cari user by identifier ada atau tidak
		2. kalau ada maka throw
		3. kalau ga ada maka lanjut tuk register
		4. password di hash jangan lupa
	*/
	_, err := s.FindByEmailOrUsername(reqRegister.Email, reqRegister.Username)

	if err != nil {
		return err
	}

	passwordHash, err := utils.HashPassword(reqRegister.Password)
	if err != nil {
		return err
	}

	modelUser := models.User{
		Name:     reqRegister.Email,
		Username: reqRegister.Username,
		Email:    reqRegister.Email,
		Password: passwordHash,
		Status:   1, // 1 == aktif
	}

	err = s.UserRepo.CreateUser(&modelUser)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthServiceImpl) ConfirmOtp() error {
	panic("not implemented") // TODO: Implement
}
