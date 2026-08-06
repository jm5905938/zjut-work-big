package middleware

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jm5905938/zjut-work-big/internal/apperror"
)

func testLogger(output *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(output, nil))
}

func TestRequestLogRecordsRequestSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var output bytes.Buffer

	r := gin.New()
	r.Use(RequestLog(testLogger(&output)))
	r.GET("/posts/:post_id", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/posts/5", nil)
	r.ServeHTTP(recorder, request)

	logText := output.String()
	for _, want := range []string{
		`"msg":"request completed"`,
		`"method":"GET"`,
		`"path":"/posts/5"`,
		`"route":"/posts/:post_id"`,
		`"status":204`,
	} {
		if !strings.Contains(logText, want) {
			t.Errorf("log does not contain %s: %s", want, logText)
		}
	}
}

func TestErrorHandlerLogsInternalCauseWithoutExposingIt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var output bytes.Buffer

	r := gin.New()
	r.Use(ErrorHandler(testLogger(&output)))
	r.GET("/error", func(c *gin.Context) {
		_ = c.Error(apperror.Internal(errors.New("database unavailable")))
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/error", nil)
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(output.String(), "database unavailable") {
		t.Fatalf("log does not contain internal cause: %s", output.String())
	}
	if strings.Contains(recorder.Body.String(), "database unavailable") {
		t.Fatalf("response exposed internal cause: %s", recorder.Body.String())
	}
}

func TestRecoveryLogsPanicAndStack(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var output bytes.Buffer

	r := gin.New()
	r.Use(Recovery(testLogger(&output)))
	r.GET("/panic", func(c *gin.Context) {
		panic("unexpected panic")
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
	logText := output.String()
	if !strings.Contains(logText, "unexpected panic") ||
		!strings.Contains(logText, `"stack":`) {
		t.Fatalf("panic log does not contain panic and stack: %s", logText)
	}
}
