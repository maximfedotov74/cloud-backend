package service

import (
	"context"
	"net/http"

	"github.com/maximfedotov74/cloud-api/internal/dto"
	"github.com/maximfedotov74/cloud-api/internal/model"
	"github.com/maximfedotov74/cloud-api/internal/msg"
	"github.com/maximfedotov74/cloud-api/internal/shared/db"
	"github.com/maximfedotov74/cloud-api/internal/shared/ex"
	"github.com/maximfedotov74/cloud-api/internal/shared/jwt"
	"github.com/maximfedotov74/cloud-api/internal/shared/password"
)

type userMailService interface {
	SendChangePasswordEmail(to string, subject string, code string) error
}

type userRepository interface {
	Create(ctx context.Context, input dto.CreateUser, tx db.Transaction) (*model.CreatedUser, ex.Error)
	Update(ctx context.Context, input dto.UpdateUser, id int) ex.Error
	FindById(ctx context.Context, id string) (*model.User, ex.Error)
	FindByEmail(ctx context.Context, email string) (*model.User, ex.Error)
	RemoveChangePasswordCode(ctx context.Context, userId string, tx db.Transaction) error
	FindChangePasswordCode(ctx context.Context, userId string, code string) (*model.ChangePasswordCode, ex.Error)
	CreateChangePasswordCode(ctx context.Context, userId string) (*string, ex.Error)
	ChangePassword(ctx context.Context, userId string, newPassword string) ex.Error
}

type userSessionRepository interface {
	Create(ctx context.Context, dto dto.CreateSession) ex.Error
	GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, ex.Error)
	RemoveSession(ctx context.Context, userId int, sessionId int) ex.Error
	RemoveExceptCurrentSession(ctx context.Context, userId int, sessionId int) ex.Error
	RemoveAllSessions(ctx context.Context, userId int) ex.Error
}

type jwtService interface {
	Sign(claims jwt.UserClaims) (jwt.Tokens, error)
	Parse(token string, tokenType jwt.TokenType) (*jwt.UserClaims, error)
}

type fileClient interface {
	CreateBucket(ctx context.Context, bucketName string) error
}

type UserService struct {
	userRepository     userRepository
	sessionRepository  userSessionRepository
	jwtService         jwtService
	mailService        userMailService
	transactionManager db.TransactionManager
	fileClient
}

func NewUserService(userRepository userRepository, sessionRepository userSessionRepository, jwtService jwtService,
	mailService userMailService, transactionManager db.TransactionManager, fileClient fileClient) *UserService {
	return &UserService{userRepository: userRepository, sessionRepository: sessionRepository,
		jwtService: jwtService, mailService: mailService, transactionManager: transactionManager, fileClient: fileClient}
}

func (s *UserService) Create(ctx context.Context, input dto.CreateUser) (*model.CreatedUser, ex.Error) {

	tx, err := s.transactionManager.Begin(ctx)

	var fall ex.Error = nil

	if err != nil {
		fall = ex.ServerError(err.Error())
		return nil, fall
	}

	defer func() {
		if fall != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	res, fall := s.userRepository.Create(ctx, input, tx)

	if fall != nil {
		return nil, fall
	}

	err = s.fileClient.CreateBucket(ctx, res.UserId)

	if err != nil {
		fall = ex.ServerError(err.Error())
		return nil, fall
	}

	return res, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (*model.User, ex.Error) {
	return s.userRepository.FindById(ctx, id)
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, ex.Error) {
	return s.userRepository.FindByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, input dto.UpdateUser, id int) ex.Error {
	return s.userRepository.Update(ctx, input, id)
}

func (s *UserService) CreateChangePasswordCode(ctx context.Context, userId string) ex.Error {

	currentUser, fall := s.FindById(ctx, userId)
	if fall != nil {
		return fall
	}

	code, err := s.userRepository.CreateChangePasswordCode(ctx, currentUser.UserId)

	if err != nil {
		return err
	}

	go s.mailService.SendChangePasswordEmail(currentUser.Email, "Код для смены пароля", *code)

	return nil
}

func (us *UserService) ConfirmChangePassword(ctx context.Context, code string, userId string) ex.Error {
	_, fall := us.userRepository.FindChangePasswordCode(ctx, userId, code)

	if fall != nil {
		return fall
	}

	err := us.userRepository.RemoveChangePasswordCode(ctx, userId, nil)

	if err != nil {
		return ex.ServerError(err.Error())
	}

	return nil

}

func (s *UserService) ChangePassword(ctx context.Context, input dto.ChangePassword, localSession model.LocalSession) (*jwt.Tokens, ex.Error) {

	user, fall := s.FindById(ctx, localSession.UserId)

	if fall != nil {
		return nil, fall
	}

	oldMatch := password.ComparePasswords(user.PasswordHash, input.OldPassword)

	if !oldMatch {
		return nil, ex.NewErr(msg.BadPassword, http.StatusBadRequest)
	}

	newMatch := password.ComparePasswords(user.PasswordHash, input.NewPassword)

	if newMatch {
		return nil, ex.NewErr(msg.BadNewPassword, http.StatusBadRequest)
	}

	newHash, err := password.HashPassword(input.NewPassword)

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	fall = s.userRepository.ChangePassword(ctx, user.UserId, newHash)

	if fall != nil {
		return nil, fall
	}

	tokens, err := s.jwtService.Sign(jwt.UserClaims{UserId: user.UserId, UserAgent: localSession.UserAgent})

	if err != nil {
		return nil, ex.ServerError(err.Error())
	}

	tokenDto := dto.CreateSession{UserId: user.UserId, UserAgent: localSession.UserAgent, Token: tokens.RefreshToken}
	fall = s.sessionRepository.Create(ctx, tokenDto)

	if fall != nil {
		return nil, fall
	}

	return &tokens, nil
}

func (s *UserService) GetUserSessions(ctx context.Context, userId int, token string) (*model.UserSessionsResponse, ex.Error) {
	return s.sessionRepository.GetUserSessions(ctx, userId, token)
}
