package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/dto"
)

type fakeAuthService struct {
	providedAdminRegisterCode string
}

func (s *fakeAuthService) Register(
	_ context.Context,
	_ dto.RegisterRequest,
	providedAdminRegisterCode string,
) (dto.UserResponse, error) {
	s.providedAdminRegisterCode = providedAdminRegisterCode
	return dto.UserResponse{}, nil
}

func (s *fakeAuthService) Login(
	context.Context,
	dto.LoginRequest,
) (dto.LoginData, error) {
	return dto.LoginData{}, nil
}

func TestRegisterPassesAdminRegisterCodeHeaderToService(t *testing.T) {
	gin.SetMode(gin.TestMode)
	auth := &fakeAuthService{}
	handler := NewAuthHandler(auth)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/auth/register",
		strings.NewReader(`{
			"username":"20260002",
			"name":"测试管理员",
			"password":"password",
			"role":"admin"
		}`),
	)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Admin-Register-Code", "admin-code")

	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	handler.Register(ctx)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("Register() status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if auth.providedAdminRegisterCode != "admin-code" {
		t.Fatalf(
			"Register() code = %q, want %q",
			auth.providedAdminRegisterCode,
			"admin-code",
		)
	}
}
