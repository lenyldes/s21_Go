package auth

import (
	"auth/api"
	"auth/internal/shared"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	repository             AuthRepository
	jwtProvider            *JwtProvider
	authenticator          *UserAuthenticator
	changePasswordNotifier *ChangePasswordNotifier
	logger                 *slog.Logger
	api.UnimplementedAuthServiceServer
}

func NewAuthService(repository AuthRepository, jwtProvider *JwtProvider, authenticator *UserAuthenticator, changePasswordNotifier *ChangePasswordNotifier, logger *slog.Logger) *AuthService {
	return &AuthService{
		repository:             repository,
		jwtProvider:            jwtProvider,
		authenticator:          authenticator,
		changePasswordNotifier: changePasswordNotifier,
		logger:                 logger,
	}
}

func (s *AuthService) User(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	var user *User
	switch identity := req.Identity.(type) {
	case *api.UserRequest_Id:
		userId, err := uuid.Parse(identity.Id)
		if err != nil {
			s.logger.Debug("uuid parsing error", slog.String("error", err.Error()))
			return nil, status.New(codes.NotFound, "user not found").Err()
		}
		u, err := s.repository.GetUserById(ctx, userId)
		if err != nil {
			s.logger.Debug("get user error", slog.String("error", err.Error()))
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.New(codes.NotFound, "user not found").Err()
			}
			return nil, status.New(codes.Internal, "internal error").Err()
		}
		user = u
	case *api.UserRequest_Login:
		u, err := s.repository.GetUserByLogin(ctx, identity.Login)
		if err != nil {
			s.logger.Debug("get user error", slog.String("error", err.Error()))
			if errors.Is(err, sql.ErrNoRows) {
				return nil, status.New(codes.NotFound, "user not found").Err()
			}
			return nil, status.New(codes.Internal, "internal error").Err()
		}
		user = u
	}
	res := &api.UserResponse{
		Id:    user.Id.String(),
		Login: user.Login,
		Email: user.Email,
		Role:  user.Role,
	}
	return res, nil
}

func (s *AuthService) Signup(ctx context.Context, req *api.SignupRequest) (*api.TokensResponse, error) {
	user, err := NewUser(req.Login, req.Name, req.Surname, req.Phonenumber, req.Password, req.Email, req.Role)
	if err != nil {
		s.logger.Debug("create user error", slog.String("error", err.Error()))
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	if err := s.repository.SaveUser(ctx, user); err != nil {
		s.logger.Debug("save user error sign up", slog.String("error", err.Error()))
		if errors.Is(err, ErrLoginTaken) {
			return nil, status.New(codes.AlreadyExists, err.Error()).Err()
		}
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	accessToken, refreshToken, err := s.jwtProvider.IssueTokens(user.Id, user.Role)
	if err != nil {
		return nil, status.New(codes.Internal, "cannot create jwt tokens").Err()
	}
	res := &api.TokensResponse{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}
	return res, nil
}

func (s *AuthService) Signin(ctx context.Context, req *api.SigninRequest) (*api.TokensResponse, error) {
	user, err := s.repository.GetUserByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Debug("get user error", slog.String("error", err.Error()))
		return nil, status.New(codes.NotFound, "user not found").Err()
	}
	if err := s.authenticator.VerifyCredentials(ctx, user, req.Password); err != nil {
		s.logger.Warn("verify credentials failed", slog.String("user id", user.Id.String()))
		return nil, status.New(codes.Unauthenticated, "invalid login and/or password").Err()
	}
	accessToken, refreshToken, err := s.jwtProvider.IssueTokens(user.Id, user.Role)
	if err != nil {
		return nil, status.New(codes.Internal, "cannot create jwt tokens").Err()
	}
	res := &api.TokensResponse{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}
	return res, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *api.RefreshTokensRequest) (*api.TokensResponse, error) {
	accessToken, refreshToken, userId, err := s.jwtProvider.Refresh(req.RefreshToken)
	if err != nil {
		s.logger.Warn("invalid refresh token", slog.String("error", err.Error()), slog.String("user id", userId))
		return nil, status.New(codes.Unauthenticated, "invalid refresh token").Err()
	}
	res := &api.TokensResponse{
		AccessToken:  string(accessToken),
		RefreshToken: string(refreshToken),
	}
	return res, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, req *api.AccessTokenRequest) (*api.EmptyMessage, error) {
	claims := AccessClaims{}
	token, err := jwt.ParseWithClaims(req.AccessToken, &claims, s.jwtProvider.AccessTokenKeyFunc, s.jwtProvider.Options...)
	if err != nil || !token.Valid {
		if err != nil {
			s.logger.Warn("invalid access token", slog.String("error", err.Error()))
		} else {
			s.logger.Warn("invalid access token")
		}
		return nil, status.New(codes.Unauthenticated, "access token not valid").Err()
	}
	if claims.TokenType != ACCESS_TOKEN_TYPE {
		s.logger.Warn("not access token")
		return nil, status.New(codes.Unauthenticated, "not access token").Err()
	}
	userId, err := uuid.Parse(claims.UserId)
	if err != nil {
		s.logger.Warn("invalid user id in access token", slog.String("error", err.Error()))
		return nil, status.New(codes.InvalidArgument, "invalid user id in access token").Err()
	}
	if _, err := s.repository.GetUserById(ctx, userId); err != nil {
		if errors.Is(err, ErrUserNotFound) {
			s.logger.Info("user not found", slog.String("error", err.Error()))
			return nil, status.New(codes.NotFound, "user no longer exists").Err()
		}
		return nil, status.New(codes.Internal, "internal server error").Err()
	}
	res := &api.EmptyMessage{}
	return res, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *api.ResetPasswordRequest) (*api.EmptyMessage, error) {
	user, err := s.repository.GetUserByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Debug("user not found", slog.String("error", err.Error()))
		if errors.Is(err, ErrUserNotFound) {

			return nil, status.New(codes.NotFound, "user not found").Err()
		}
		return nil, status.New(codes.Internal, "internal server error").Err()
	}
	if user.Email == "" {
		s.logger.Warn("email is empty", slog.String("user id", user.Id.String()))
		return nil, status.New(codes.InvalidArgument, "email is empty").Err()
	}
	tmpPassword, err := shared.GenerateTemporaryPassword()
	if err != nil {
		s.logger.Error("generate new password error", slog.String("error", err.Error()))
		return nil, status.New(codes.Internal, "cannot reset password").Err()
	}
	hashPassword, err := hashPassword(tmpPassword)
	if err != nil {
		s.logger.Error("generate hash password error", slog.String("error", err.Error()))
		return nil, status.New(codes.Internal, "cannot reset password").Err()
	}
	if err := s.repository.ChangeUserPassword(ctx, user.Id, hashPassword); err != nil {
		s.logger.Error("change user password error", slog.String("error", err.Error()))
		return nil, status.New(codes.Internal, "cannot reset password").Err()
	}
	if err := s.changePasswordNotifier.Notify(ctx, user.Login, tmpPassword, user.Email); err != nil {
		s.logger.Error("reset password notify error", slog.String("error", err.Error()))
		return nil, status.New(codes.Internal, "can not send email with new temporary password").Err()
	}
	res := &api.EmptyMessage{}
	return res, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, req *api.ChangePasswordRequest) (*api.EmptyMessage, error) {
	user, err := s.repository.GetUserByLogin(ctx, req.Login)
	if err != nil {
		s.logger.Info("user not found", slog.String("error", err.Error()))
		if errors.Is(err, ErrUserNotFound) {
			return nil, status.New(codes.NotFound, err.Error()).Err()
		}
		return nil, status.New(codes.Internal, "internal server error").Err()
	}
	if s.authenticator.VerifyCredentials(ctx, user, req.CurrentPassword) != nil {
		s.logger.Warn("verify creadentials failed", slog.String("user id", user.Id.String()))
		return nil, status.New(codes.Unauthenticated, "invalid login and/or password").Err()
	}
	if err := user.ChangePassword(req.NewPassword); err != nil {
		s.logger.Info("change password failed", slog.String("user id", user.Id.String()))

		return nil, status.New(codes.Internal, "can not change password").Err()
	}
	if err := s.repository.ChangeUserPassword(ctx, user.Id, user.PasswordHash); err != nil {
		s.logger.Info("new password not saved", slog.String("user id", user.Id.String()))
		return nil, status.New(codes.Internal, "can not change password").Err()
	}
	e := api.EmptyMessage{}
	return &e, nil
}
