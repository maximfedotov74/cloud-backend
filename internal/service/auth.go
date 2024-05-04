package service

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/jwt"
	"github.com/maximfedotov74/cloud-api/internal/shared/password"
)

type authUserService interface {
	Create(ctx context.Context, input dto.CreateUser) (*model.CreatedUser, ex.Error)
	FindByEmail(ctx context.Context, email string) (*model.User, ex.Error)
	FindById(ctx context.Context, id string) (*model.User, ex.Error)
}

type authSessionRepository interface {
	Create(ctx context.Context, dto dto.CreateSession) ex.Error
	FindByAgentAndToken(ctx context.Context, agent string, token string) (*model.Session, ex.Error)
	RemoveSessionByToken(ctx context.Context, token string) ex.Error
}

type authMailService interface {
	SendActivationEmail(to string, subject string, link string) error
}

type AuthService struct {
	userService    authUserService
	sessionService authSessionRepository
	mailService    authMailService
	jwtService     *jwt.JwtService
}

func NewAuthService(userService authUserService, sessionService authSessionRepository,
	mailService authMailService, jwtService *jwt.JwtService) *AuthService {
	return &AuthService{
		userService:    userService,
		sessionService: sessionService,
		mailService:    mailService,
		jwtService:     jwtService,
	}
}

func (s *AuthService) Login(ctx context.Context, input dto.CreateUser, userAgent string, ip net.IP) (*model.LoginResponse, ex.Error) {
	user, fall := s.userService.FindByEmail(ctx, input.Email)

	if fall != nil {
		return nil, fall
	}

	if isPasswordCorrect := password.ComparePasswords(user.PasswordHash, input.Password); !isPasswordCorrect {
		return nil, ex.NewErr(msg.InvalidCredentials, http.StatusNotFound)
	}

	tokens, err := s.jwtService.Sign(jwt.UserClaims{UserId: user.UserId, UserAgent: userAgent})

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	tokenDto := dto.CreateSession{UserId: user.UserId, UserAgent: userAgent, Token: tokens.RefreshToken, Ip: ip}
	fall = s.sessionService.Create(ctx, tokenDto)
	if fall != nil {
		return nil, fall
	}

	response := model.LoginResponse{UserId: user.UserId, Tokens: tokens}

	return &response, nil
}

func (s *AuthService) Registration(ctx context.Context, input dto.CreateUser) (*string, ex.Error) {
	user, _ := s.userService.FindByEmail(ctx, input.Email)

	if user != nil {
		return nil, ex.NewErr(msg.UserIsRegistered, http.StatusBadRequest)
	}

	res, fall := s.userService.Create(ctx, input)

	if fall != nil {
		return nil, fall
	}

	link := fmt.Sprintf("/api/user/activate/%s", res.ActivationLink)

	go s.mailService.SendActivationEmail(res.Email, "Активация аккаунта в облачном хранилище", link)

	return &res.UserId, nil

}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, ex.Error) {

	parsed, err := s.jwtService.Parse(refreshToken, jwt.RefreshToken)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	dbToken, fall := s.sessionService.FindByAgentAndToken(ctx, parsed.UserAgent, refreshToken)

	if fall != nil {
		return nil, fall
	}

	user, fall := s.userService.FindById(ctx, dbToken.UserId)
	if fall != nil {
		return nil, fall
	}

	claims := jwt.UserClaims{UserId: user.UserId, UserAgent: dbToken.UserAgent}

	tokens, err := s.jwtService.Sign(claims)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	tokenDto := dto.CreateSession{UserId: user.UserId, UserAgent: dbToken.UserAgent, Token: tokens.RefreshToken}
	fall = s.sessionService.Create(ctx, tokenDto)

	if fall != nil {
		return nil, fall
	}

	response := model.LoginResponse{UserId: user.UserId, Tokens: tokens}

	return &response, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) ex.Error {
	return s.sessionService.RemoveSessionByToken(ctx, token)
}
