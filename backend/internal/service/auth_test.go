package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jm5905938/zjut-work-big/internal/apperror"
	"github.com/jm5905938/zjut-work-big/internal/dto"
	"github.com/jm5905938/zjut-work-big/internal/model"
	"github.com/jm5905938/zjut-work-big/internal/password"
	"gorm.io/gorm"
)

type fakeUserRepository struct {
	created   *model.User
	createErr error
}

func (r *fakeUserRepository) Create(_ context.Context, user *model.User) error {
	if r.createErr != nil {
		return r.createErr
	}

	user.ID = 1
	created := *user
	r.created = &created
	return nil
}

func (r *fakeUserRepository) FindByUsername(
	context.Context,
	string,
) (*model.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func TestRegisterStudentDoesNotRequireAdminCode(t *testing.T) {
	repository := &fakeUserRepository{}
	service := NewAuthService(repository, nil, "admin-code")

	user, err := service.Register(context.Background(), dto.RegisterRequest{
		Username: "20260001",
		Name:     "测试学生",
		Password: "password",
		Role:     model.UserRoleStudent,
	}, "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user.Role != model.UserRoleStudent {
		t.Fatalf("Register() role = %q, want %q", user.Role, model.UserRoleStudent)
	}
	if repository.created == nil {
		t.Fatal("Register() did not create user")
	}
	if err := password.Verify(repository.created.PasswordHash, "password"); err != nil {
		t.Fatalf("stored password hash cannot verify password: %v", err)
	}
}

func TestRegisterAdminRequiresCorrectCode(t *testing.T) {
	tests := []struct {
		name           string
		configuredCode string
		providedCode   string
		wantForbidden  bool
	}{
		{name: "correct code", configuredCode: "admin-code", providedCode: "admin-code"},
		{name: "missing code", configuredCode: "admin-code", wantForbidden: true},
		{name: "wrong code", configuredCode: "admin-code", providedCode: "wrong", wantForbidden: true},
		{name: "registration disabled", providedCode: "admin-code", wantForbidden: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &fakeUserRepository{}
			service := NewAuthService(repository, nil, tt.configuredCode)

			user, err := service.Register(context.Background(), dto.RegisterRequest{
				Username: "20260002",
				Name:     "测试管理员",
				Password: "password",
				Role:     model.UserRoleAdmin,
			}, tt.providedCode)

			if tt.wantForbidden {
				var appErr *apperror.Error
				if !errors.As(err, &appErr) || appErr.Status != http.StatusForbidden {
					t.Fatalf("Register() error = %v, want 403 app error", err)
				}
				if repository.created != nil {
					t.Fatal("Register() created admin with invalid code")
				}
				return
			}

			if err != nil {
				t.Fatalf("Register() error = %v", err)
			}
			if user.Role != model.UserRoleAdmin || repository.created == nil {
				t.Fatalf("Register() did not create admin: user = %#v", user)
			}
		})
	}
}

func TestRegisterKeepsDuplicateUsernameConflict(t *testing.T) {
	repository := &fakeUserRepository{createErr: gorm.ErrDuplicatedKey}
	service := NewAuthService(repository, nil, "admin-code")

	_, err := service.Register(context.Background(), dto.RegisterRequest{
		Username: "20260001",
		Name:     "重复用户",
		Password: "password",
		Role:     model.UserRoleStudent,
	}, "")

	var appErr *apperror.Error
	if !errors.As(err, &appErr) || appErr.Status != http.StatusConflict {
		t.Fatalf("Register() error = %v, want 409 app error", err)
	}
}
