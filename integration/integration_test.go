package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dwilanang/psp/internal/admin"
	"github.com/dwilanang/psp/internal/admin/dto"
	"github.com/dwilanang/psp/internal/admin/handler"
	"github.com/dwilanang/psp/internal/admin/repository"
	"github.com/dwilanang/psp/internal/admin/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupRouter() *gin.Engine {
	// Connect to test DB
	db := sqlx.MustConnect("postgres", "postgres://ppms_user:ppms123@localhost:5432/ppms?sslmode=disable")

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	deps := admin.Dependencies{
		Service: svc,
	}
	h := handler.NewHandler(deps)

	r := gin.Default()
	admin := r.Group("/admin")
	{
		admin.POST("/attendance-periods", h.CreateAttendancePeriod)
	}

	return r
}

func Test_CreateAttendancePeriod_Endpoint(t *testing.T) {
	router := setupRouter()

	// Payload
	reqBody := dto.AttendancePeriodRequest{
		StartDate: "2025-06-05",
		EndDate:   "2025-06-10",
	}
	jsonBody, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/attendance-periods", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}
