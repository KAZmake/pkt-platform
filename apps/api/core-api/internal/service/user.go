package service

import (
	"context"
	"fmt"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

type userRepo interface {
	Upsert(ctx context.Context, u *model.User) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, inp repository.UpdateUserInput) (*model.User, error)
	List(ctx context.Context, limit, offset int) ([]*model.User, int, error)
}

type borrowerRepo interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Borrower, error)
	Update(ctx context.Context, id uuid.UUID, inp repository.UpdateBorrowerInput) (*model.Borrower, error)
}

type UserService struct {
	users     userRepo
	borrowers borrowerRepo
}

func NewUserService(users userRepo, borrowers borrowerRepo) *UserService {
	return &UserService{users: users, borrowers: borrowers}
}

// SyncFromClaims upserts the user in our DB using Keycloak JWT claims.
// Called on every authenticated request so the DB stays in sync.
func (s *UserService) SyncFromClaims(ctx context.Context, claims *middleware.KeycloakClaims) (*model.User, error) {
	role := model.HighestRole(claims.RealmAccess.Roles)

	u := &model.User{
		KeycloakID: claims.Subject,
		Email:      claims.Email,
		Role:       role,
	}
	if claims.GivenName != "" {
		u.FirstName = &claims.GivenName
	}
	if claims.FamilyName != "" {
		u.LastName = &claims.FamilyName
	}

	synced, err := s.users.Upsert(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("sync user: %w", err)
	}
	return synced, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

type UpdateProfileInput struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
}

func (s *UserService) UpdateProfile(ctx context.Context, id uuid.UUID, inp UpdateProfileInput) (*model.User, error) {
	return s.users.Update(ctx, id, repository.UpdateUserInput{
		FirstName: inp.FirstName,
		LastName:  inp.LastName,
		Phone:     inp.Phone,
	})
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	return s.users.List(ctx, limit, offset)
}

func (s *UserService) GetBorrowerProfile(ctx context.Context, userID uuid.UUID) (*model.Borrower, error) {
	return s.borrowers.GetByUserID(ctx, userID)
}

type UpdateBorrowerInput struct {
	OrgName      *string `json:"org_name"`
	ActivityType *string `json:"activity_type"`
}

func (s *UserService) UpdateBorrowerProfile(ctx context.Context, userID uuid.UUID, inp UpdateBorrowerInput) (*model.Borrower, error) {
	b, err := s.borrowers.GetByUserID(ctx, userID)
	if err != nil || b == nil {
		return nil, fmt.Errorf("borrower profile not found")
	}
	return s.borrowers.Update(ctx, b.ID, repository.UpdateBorrowerInput{
		OrgName:      inp.OrgName,
		ActivityType: inp.ActivityType,
	})
}
