package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/middleware"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/model"
	"github.com/KAZmake/pkt-platform/apps/api/core-api/internal/repository"
	"github.com/google/uuid"
)

// ── mocks ─────────────────────────────────────────────────────────────────────

type mockUserRepo struct {
	upsertFn  func(context.Context, *model.User) (*model.User, error)
	getByIDFn func(context.Context, uuid.UUID) (*model.User, error)
	updateFn  func(context.Context, uuid.UUID, repository.UpdateUserInput) (*model.User, error)
	listFn    func(context.Context, int, int) ([]*model.User, int, error)
}

func (m *mockUserRepo) Upsert(ctx context.Context, u *model.User) (*model.User, error) {
	if m.upsertFn != nil {
		return m.upsertFn(ctx, u)
	}
	return u, nil
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockUserRepo) Update(ctx context.Context, id uuid.UUID, inp repository.UpdateUserInput) (*model.User, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inp)
	}
	return nil, nil
}
func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*model.User, int, error) {
	if m.listFn != nil {
		return m.listFn(ctx, limit, offset)
	}
	return nil, 0, nil
}

type mockBorrowerRepo struct {
	getByUserIDFn func(context.Context, uuid.UUID) (*model.Borrower, error)
	updateFn      func(context.Context, uuid.UUID, repository.UpdateBorrowerInput) (*model.Borrower, error)
}

func (m *mockBorrowerRepo) GetByUserID(ctx context.Context, id uuid.UUID) (*model.Borrower, error) {
	if m.getByUserIDFn != nil {
		return m.getByUserIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockBorrowerRepo) Update(ctx context.Context, id uuid.UUID, inp repository.UpdateBorrowerInput) (*model.Borrower, error) {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inp)
	}
	return nil, nil
}

func stubUser() *model.User {
	fn := "Иван"
	ln := "Иванов"
	return &model.User{
		ID:         uuid.New(),
		KeycloakID: "kc-123",
		Email:      "ivan@example.com",
		Role:       model.RoleBorrower,
		FirstName:  &fn,
		LastName:   &ln,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func stubBorrower(userID uuid.UUID) *model.Borrower {
	return &model.Borrower{
		ID: uuid.New(), UserID: userID, INN: "123456789012", CreatedAt: time.Now(),
	}
}

// ── SyncFromClaims ────────────────────────────────────────────────────────────

func TestUserService_SyncFromClaims_Basic(t *testing.T) {
	ctx := context.Background()
	want := stubUser()

	svc := NewUserService(
		&mockUserRepo{upsertFn: func(_ context.Context, u *model.User) (*model.User, error) { return want, nil }},
		&mockBorrowerRepo{},
	)
	claims := &middleware.KeycloakClaims{}
	claims.Subject = "kc-123"
	claims.Email = "ivan@example.com"
	claims.RealmAccess.Roles = []string{"borrower"}

	got, err := svc.SyncFromClaims(ctx, claims)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != want.ID {
		t.Errorf("want ID %v, got %v", want.ID, got)
	}
}

func TestUserService_SyncFromClaims_SetsFirstAndLastName(t *testing.T) {
	ctx := context.Background()
	var capturedUser *model.User

	svc := NewUserService(
		&mockUserRepo{upsertFn: func(_ context.Context, u *model.User) (*model.User, error) {
			capturedUser = u
			return u, nil
		}},
		&mockBorrowerRepo{},
	)
	claims := &middleware.KeycloakClaims{}
	claims.Subject = "kc-abc"
	claims.Email = "ali@example.com"
	claims.GivenName = "Али"
	claims.FamilyName = "Аликов"
	claims.RealmAccess.Roles = []string{"employee"}

	svc.SyncFromClaims(ctx, claims) //nolint:errcheck

	if capturedUser.FirstName == nil || *capturedUser.FirstName != "Али" {
		t.Errorf("want FirstName 'Али', got %v", capturedUser.FirstName)
	}
	if capturedUser.LastName == nil || *capturedUser.LastName != "Аликов" {
		t.Errorf("want LastName 'Аликов', got %v", capturedUser.LastName)
	}
}

func TestUserService_SyncFromClaims_EmptyNames(t *testing.T) {
	ctx := context.Background()
	var capturedUser *model.User

	svc := NewUserService(
		&mockUserRepo{upsertFn: func(_ context.Context, u *model.User) (*model.User, error) {
			capturedUser = u
			return u, nil
		}},
		&mockBorrowerRepo{},
	)
	claims := &middleware.KeycloakClaims{}
	claims.Subject = "kc-xyz"
	claims.Email = "anon@example.com"
	// GivenName and FamilyName are empty

	svc.SyncFromClaims(ctx, claims) //nolint:errcheck

	if capturedUser.FirstName != nil {
		t.Error("FirstName should be nil for empty GivenName")
	}
	if capturedUser.LastName != nil {
		t.Error("LastName should be nil for empty FamilyName")
	}
}

func TestUserService_SyncFromClaims_HighestRolePicked(t *testing.T) {
	ctx := context.Background()
	var capturedRole string

	svc := NewUserService(
		&mockUserRepo{upsertFn: func(_ context.Context, u *model.User) (*model.User, error) {
			capturedRole = u.Role
			return u, nil
		}},
		&mockBorrowerRepo{},
	)
	claims := &middleware.KeycloakClaims{}
	claims.Subject = "kc-admin"
	claims.RealmAccess.Roles = []string{"borrower", "employee", "admin"}

	svc.SyncFromClaims(ctx, claims) //nolint:errcheck

	if capturedRole != model.RoleAdmin {
		t.Errorf("want role 'admin', got %q", capturedRole)
	}
}

func TestUserService_SyncFromClaims_RepoError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	svc := NewUserService(
		&mockUserRepo{upsertFn: func(_ context.Context, _ *model.User) (*model.User, error) { return nil, dbErr }},
		&mockBorrowerRepo{},
	)
	claims := &middleware.KeycloakClaims{}
	claims.Subject = "kc-xyz"

	_, err := svc.SyncFromClaims(ctx, claims)
	if err == nil {
		t.Fatal("expected error")
	}
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestUserService_GetByID(t *testing.T) {
	ctx := context.Background()
	want := stubUser()
	svc := NewUserService(
		&mockUserRepo{getByIDFn: func(_ context.Context, id uuid.UUID) (*model.User, error) {
			if id == want.ID {
				return want, nil
			}
			return nil, nil
		}},
		&mockBorrowerRepo{},
	)
	got, err := svc.GetByID(ctx, want.ID)
	if err != nil || got == nil || got.ID != want.ID {
		t.Errorf("unexpected result: %v, %v", got, err)
	}
}

// ── UpdateProfile ─────────────────────────────────────────────────────────────

func TestUserService_UpdateProfile(t *testing.T) {
	ctx := context.Background()
	want := stubUser()
	svc := NewUserService(
		&mockUserRepo{updateFn: func(_ context.Context, _ uuid.UUID, _ repository.UpdateUserInput) (*model.User, error) {
			return want, nil
		}},
		&mockBorrowerRepo{},
	)
	name := "Новое имя"
	got, err := svc.UpdateProfile(ctx, want.ID, UpdateProfileInput{FirstName: &name})
	if err != nil || got == nil {
		t.Errorf("unexpected result: %v, %v", got, err)
	}
}

// ── UpdateBorrowerProfile ─────────────────────────────────────────────────────

func TestUserService_UpdateBorrowerProfile_NotFound(t *testing.T) {
	ctx := context.Background()
	svc := NewUserService(
		&mockUserRepo{},
		&mockBorrowerRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*model.Borrower, error) { return nil, nil },
		},
	)
	_, err := svc.UpdateBorrowerProfile(ctx, uuid.New(), UpdateBorrowerInput{})
	if err == nil {
		t.Fatal("expected error for missing borrower profile")
	}
	if err.Error() != "borrower profile not found" {
		t.Errorf("unexpected error: %q", err.Error())
	}
}

func TestUserService_UpdateBorrowerProfile_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	b := stubBorrower(userID)
	orgName := "ТОО Тест"

	svc := NewUserService(
		&mockUserRepo{},
		&mockBorrowerRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*model.Borrower, error) { return b, nil },
			updateFn: func(_ context.Context, _ uuid.UUID, _ repository.UpdateBorrowerInput) (*model.Borrower, error) {
				return b, nil
			},
		},
	)
	got, err := svc.UpdateBorrowerProfile(ctx, userID, UpdateBorrowerInput{OrgName: &orgName})
	if err != nil || got == nil {
		t.Errorf("unexpected: %v, %v", got, err)
	}
}

func TestUserService_UpdateBorrowerProfile_GetError(t *testing.T) {
	ctx := context.Background()
	dbErr := errors.New("db error")
	svc := NewUserService(
		&mockUserRepo{},
		&mockBorrowerRepo{
			getByUserIDFn: func(_ context.Context, _ uuid.UUID) (*model.Borrower, error) { return nil, dbErr },
		},
	)
	_, err := svc.UpdateBorrowerProfile(ctx, uuid.New(), UpdateBorrowerInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// ── ListUsers ─────────────────────────────────────────────────────────────────

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()
	users := []*model.User{stubUser(), stubUser()}

	svc := NewUserService(
		&mockUserRepo{listFn: func(_ context.Context, limit, offset int) ([]*model.User, int, error) {
			if limit != 20 || offset != 0 {
				return nil, 0, nil
			}
			return users, len(users), nil
		}},
		&mockBorrowerRepo{},
	)
	got, total, err := svc.ListUsers(ctx, 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || total != 2 {
		t.Errorf("want 2 users, got %d (total=%d)", len(got), total)
	}
}
