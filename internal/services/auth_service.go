package services

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/enum"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/repository"
	"github.com/MrBista/blog-api/internal/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService interface {
	LoginUser(reqLogin dto.LoginRequest) (dto.LoginResponse, error)
	RegisterUser(reqRegister dto.RegisterRequest) error
	ConfirmOtp() error

	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(code string) (dto.LoginResponse, error)
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

	if user.AuthProvider == "google" {
		return responseLogin, exception.NewBadRequestErr("please login with Google")
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

// Get Google OAuth URL
func (s *AuthServiceImpl) GetGoogleAuthURL(state string) string {
	if state == "" {
		state = utils.GenerateRandomString(32)
	}

	return utils.GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Handle Google Callback - IMPROVED
func (s *AuthServiceImpl) HandleGoogleCallback(code string) (dto.LoginResponse, error) {
	var responseLogin dto.LoginResponse

	// 1. Exchange authorization code dengan access token
	token, err := utils.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return responseLogin, errors.New("failed to exchange token")
	}

	// 2. Ambil user info dari Google
	googleUser, err := s.getGoogleUserInfo(token)
	if err != nil {
		return responseLogin, err
	}

	// 3. Cari atau buat user di database
	user, err := s.findOrCreateGoogleUser(googleUser)
	if err != nil {
		return responseLogin, err
	}

	// 4. Generate JWT token
	jwtService := utils.GetJwtService()
	accessToken, err := jwtService.CreateAccessToken(int(user.ID), user.Role)
	if err != nil {
		return responseLogin, exception.NewBadRequestErr("failed to generate token")
	}

	responseLogin.AccessToken = accessToken
	responseLogin.TypeToken = "access_token"

	return responseLogin, nil
}
func (s *AuthServiceImpl) getGoogleUserInfo(token *oauth2.Token) (*dto.GoogleUserInfo, error) {
	client := utils.GoogleOAuthConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info from google")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read user info")
	}

	var googleUser dto.GoogleUserInfo
	if err := json.Unmarshal(data, &googleUser); err != nil {
		return nil, errors.New("failed to parse user info")
	}

	return &googleUser, nil
}

func (s *AuthServiceImpl) findOrCreateGoogleUser(googleUser *dto.GoogleUserInfo) (*models.User, error) {
	/*
	   flow:
	   1. Cek by username (Google ID) → kalau ada, berarti user Google yang sama
	   2. Kalau tidak ada, cek by email → mungkin user sudah register via form biasa
	   3. Kalau email ada tapi auth_provider = 'local' → link account
	   4. Kalau benar-benar baru → create user baru
	*/

	// 1. Cek apakah Google ID (di column username) sudah ada
	userByUsername, err := s.UserRepo.FindByUsername(googleUser.ID)
	if err == nil {
		userByUsername.Name = googleUser.Name
		userByUsername.ProfileImageURI = &googleUser.Picture
		userByUsername.Email = googleUser.Email // update email juga kalau berubah
		userByUsername.Status = 1

		if updateErr := s.UserRepo.Update(userByUsername); updateErr != nil {
			return nil, errors.New("failed to update user")
		}

		return userByUsername, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, exception.NewGormDBErr(err)
	}

	userByEmail, err := s.UserRepo.FindByEmail(googleUser.Email)
	if err == nil {
		if userByEmail.AuthProvider == "google" {
			// Aneh, harusnya ketemu di step 1. Mungkin data corrupt
			return nil, exception.NewBadRequestErr("inconsistent data: email exists but google id doesn't match")
		}

		userByEmail.Username = googleUser.ID
		userByEmail.AuthProvider = "google"
		userByEmail.Name = googleUser.Name
		userByEmail.ProfileImageURI = &googleUser.Picture
		userByUsername.Status = 1

		if updateErr := s.UserRepo.Update(userByEmail); updateErr != nil {
			return nil, exception.NewBadRequestErr("failed to link google account")
		}

		return userByEmail, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, exception.NewGormDBErr(err)
	}

	newUser := &models.User{
		Email:           googleUser.Email,
		Username:        googleUser.ID, // Google ID disimpan di username
		Name:            googleUser.Name,
		ProfileImageURI: &googleUser.Picture,
		AuthProvider:    "google",
		Role:            int(enum.RoleReader),
		Status:          1,
		Password:        "", // no password untuk Google OAuth
	}

	if err := s.UserRepo.Create(newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
}

func (s *AuthServiceImpl) RegisterUser(reqRegister dto.RegisterRequest) error {
	/*
		1. cari user by identifier ada atau tidak
		2. kalau ada maka throw
		3. kalau ga ada maka lanjut tuk register
		4. password di hash jangan lupa
	*/
	_, err := s.FindByEmailOrUsername(reqRegister.Email, reqRegister.Username)

	if err == nil {
		return exception.NewBadRequestErr("Username/Password is invalid")
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
		Status:   1, // 1 == aktif,
		Role:     int(enum.RoleReader),
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
