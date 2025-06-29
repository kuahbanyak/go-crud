package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go-cr
	"go-crud/internal/domain/entity"
	"go-crud/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountUsecase interface {
	CreateAccount(ctx context.Context, req *entity.CreateAccountRequest) (*entity.AccountResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.AccountResponse, error)
	GetAccounts(ctx context.Context, limit, offset int) ([]*entity.AccountResponse, int64, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, req *entity.UpdateAccountRequest) (*entity.AccountResponse, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, id uuid.UUID, req *entity.ChangePasswordRequest) error
}

type accountUsecase struct {
	accountRepo repository.AccountRepository
	jwtService  *auth.JWTService
}

// NewAccountUsecase creates a new account usecase
func NewAccountUsecase(accountRepo repository.AccountRepository, jwtService *auth.JWTService) AccountUsecase {
	return &accountUsecase{
		accountRepo: accountRepo,
		jwtService:  jwtService,
	}
}

// CreateAccount creates a new account
func (u *accountUsecase) CreateAccount(ctx context.Context, req *entity.CreateAccountRequest) (*entity.AccountResponse, error) {
	// Check if username already exists
	existingByUsername, err := u.accountRepo.GetByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingByUsername != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingByEmail, err := u.accountRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingByEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	account := &entity.Account{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		Role:      "user",
	}

	if err := u.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return u.mapToResponse(account), nil
}

// Login authenticates a user and returns a JWT token
func (u *accountUsecase) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	account, err := u.accountRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if !account.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !auth.CheckPassword(req.Password, account.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	tokenDuration := time.Hour * 24 // 24 hours
	token, err := u.jwtService.GenerateToken(account.ID, account.Username, account.Role, tokenDuration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &entity.LoginResponse{
		Account:     u.mapToResponse(account),
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(tokenDuration.Seconds()),
	}, nil
}

// GetAccountByID retrieves an account by ID
func (u *accountUsecase) GetAccountByID(ctx context.Context, id uuid.UUID) (*entity.AccountResponse, error) {
	account, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	return u.mapToResponse(account), nil
}

// GetAccounts retrieves all accounts with pagination
func (u *accountUsecase) GetAccounts(ctx context.Context, limit, offset int) ([]*entity.AccountResponse, int64, error) {
	accounts, err := u.accountRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.accountRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*entity.AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = u.mapToResponse(account)
	}

	return responses, total, nil
}

// UpdateAccount updates an account
func (u *accountUsecase) UpdateAccount(ctx context.Context, id uuid.UUID, req *entity.UpdateAccountRequest) (*entity.AccountResponse, error) {
	// Check if account exists
	existingAccount, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	// Check if email is being updated and already exists
	if req.Email != nil && *req.Email != existingAccount.Email {
		existingByEmail, err := u.accountRepo.GetByEmail(ctx, *req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingByEmail != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Update fields if provided
	updateAccount := &entity.Account{}
	if req.Email != nil {
		updateAccount.Email = *req.Email
	}
	if req.FirstName != nil {
		updateAccount.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		updateAccount.LastName = *req.LastName
	}
	if req.IsActive != nil {
		updateAccount.IsActive = *req.IsActive
	}

	if err := u.accountRepo.Update(ctx, id, updateAccount); err != nil {
		return nil, err
	}

	// Get updated account
	updatedAccount, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.mapToResponse(updatedAccount), nil
}

// DeleteAccount deletes an account
func (u *accountUsecase) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	// Check if account exists
	_, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		return err
	}

	return u.accountRepo.Delete(ctx, id)
}

// ChangePassword changes account password
func (u *accountUsecase) ChangePassword(ctx context.Context, id uuid.UUID, req *entity.ChangePasswordRequest) error {
	account, err := u.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("account not found")
		}
		return err
	}

	if !auth.CheckPassword(req.CurrentPassword, account.Password) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	return u.accountRepo.UpdatePassword(ctx, id, hashedPassword)
}

// mapToResponse maps account entity to response
func (u *accountUsecase) mapToResponse(account *entity.Account) *entity.AccountResponse {
	return &entity.AccountResponse{
		ID:        account.ID,
		Username:  account.Username,
		Email:     account.Email,
		FirstName: account.FirstName,
		LastName:  account.LastName,
		IsActive:  account.IsActive,
		Role:      account.Role,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
