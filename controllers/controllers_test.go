package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/testutils"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	mockName     = "Mock User"
	mockEmail    = "mockemaisl@email.com"
	mockPassword = "mockpassword"
)

var db *gorm.DB
var router *gin.Engine

var mockUserCreds = auth.Credentials{
	Email:    mockEmail,
	Password: mockPassword,
}

func TestMain(m *testing.M) {
	for _, v := range os.Environ() {
		fmt.Println(v)
	}

	testutils.LoadEnv()

	var err error
	db, err = models.OpenConnection(os.Getenv("TEST_DB_URI"), logger.Discard)
	if err != nil {
		panic(err)
	}

	router = SetupRouter(db)
	m.Run()

	cleanup()
}

func cleanup() {
	db.Exec("truncate users cascade;")
	db.Exec("truncate notes cascade;")
}

func createMockUser(creds *auth.Credentials) (*models.User, error) {
	if creds == nil {
		creds = &mockUserCreds
	}
	// create the user
	hash, err := auth.HashPassword(creds.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     mockName,
		Email:    creds.Email,
		Password: hash,
	}

	if err := db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func login(form auth.Credentials) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body, _ := json.Marshal(form)
	req, _ := http.NewRequest("POST", API+APILogin, bytes.NewReader(body))
	router.ServeHTTP(w, req)
	return w
}

func serveHTTP(method, url string, body io.Reader, cookies []*http.Cookie) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	return w
}
