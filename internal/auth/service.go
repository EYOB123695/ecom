package auth

import (
	"context"
	"errors"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (AuthResponse, error)
	GetUserByID(ctx context.Context, id int64) (UserResponse, error)
}

type svc struct {
	repo      repo.Querier
	jwtSecret string
}

func NewService(queries repo.Querier, jwtSecret string) Service {
	return &svc{repo: queries, jwtSecret: jwtSecret}
}

func (s *svc) Register(ctx context.Context, req RegisterRequest) (AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return AuthResponse{}, ErrInvalidInput
	}
	if len(req.Password) < 8 {
		return AuthResponse{}, ErrInvalidInput
	}

	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return AuthResponse{}, ErrEmailTaken
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return AuthResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         "customer",
	})
	if err != nil {
		return AuthResponse{}, err
	}

	token, err := createToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *svc) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return AuthResponse{}, ErrInvalidInput
	}

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, ErrInvalidCredentials
		}
		return AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return AuthResponse{}, ErrInvalidCredentials
	}

	token, err := createToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (s *svc) GetUserByID(ctx context.Context, id int64) (UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return toUserResponse(user), nil
}

func toUserResponse(user repo.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}
