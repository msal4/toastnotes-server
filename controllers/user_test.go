package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/testutils"
	"github.com/stretchr/testify/assert"
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
}

func TestRegister(t *testing.T) {

	testRegister := func(t *testing.T, form auth.RegisterForm, expectedStatus int) {
		w := httptest.NewRecorder()
		body, err := json.Marshal(form)
		if err != nil {
			panic(err)
		}
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewReader(body))
		router.ServeHTTP(w, req)
		assert.Equal(t, expectedStatus, w.Code)
	}

	t.Run("register_a_new_user", func(t *testing.T) {
		defer cleanup()

		form := auth.RegisterForm{
			Name:        mockName,
			Credentials: mockUserCreds,
		}
		testRegister(t, form, http.StatusOK)

		var user models.User
		assert.Nil(t, db.First(&user, "email = ?", mockEmail).Error)
	})

	t.Run("does_not_register_with_invalid_data", func(t *testing.T) {
		defer cleanup()

		form := auth.RegisterForm{
			Name: mockName,
			Credentials: auth.Credentials{
				Email: mockEmail,
				// missing password
			},
		}
		testRegister(t, form, http.StatusNotAcceptable)

		assert.NotNil(t, db.First(&models.User{}, "email = ?", mockEmail).Error)
	})

	t.Run("does_not_register_an_existing_user", func(t *testing.T) {
		defer cleanup()
		createMockUser()
		assert.Nil(t, db.First(&models.User{}, "email = ?", mockEmail).Error)

		form := auth.RegisterForm{
			Name:        mockName,
			Credentials: mockUserCreds,
		}
		testRegister(t, form, http.StatusNotAcceptable)
	})
}

func TestLogin(t *testing.T) {
	t.Cleanup(cleanup)
	createMockUser()

	t.Run("login_a_user", func(t *testing.T) {
		w := login(mockUserCreds)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("does_not_accept_invalid_data", func(t *testing.T) {
		form := auth.Credentials{
			Email: mockEmail,
			// Password: mockPassword, // missing password
		}

		w := login(form)
		assert.Equal(t, http.StatusNotAcceptable, w.Code)
		assert.Empty(t, w.Header().Get("Set-Cookie"))
	})

	t.Run("does_not_login_non_existing_user", func(t *testing.T) {
		form := auth.Credentials{
			Email:    "mynonexistinguseremail@gmail.com",
			Password: mockPassword,
		}

		w := login(form)
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Empty(t, w.Header().Get("Set-Cookie"))
	})
}

func TestChangePassword(t *testing.T) {
	createMockUser()
	t.Cleanup(cleanup)

	t.Run("change_user_password", func(t *testing.T) {
		w := login(mockUserCreds)

		form := auth.ChangePasswordForm{
			CurrentPassword: mockPassword,
			NewPassword:     "new" + mockPassword,
		}
		body, _ := json.Marshal(form)
		req, _ := http.NewRequest("POST", "/api/v1/change_password", bytes.NewReader(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// this also tests the JWTAuth middleware.
	t.Run("does_not_authorize_a_non_authenticated_user", func(t *testing.T) {
		w := httptest.NewRecorder()

		form := auth.ChangePasswordForm{
			CurrentPassword: mockPassword,
			NewPassword:     "new" + mockPassword,
		}
		body, _ := json.Marshal(form)
		req, _ := http.NewRequest("POST", "/api/v1/change_password", bytes.NewReader(body))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

}

func TestMe(t *testing.T) {
	createMockUser()
	t.Cleanup(cleanup)

	loginW := login(mockUserCreds)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/me", nil)
	for _, c := range loginW.Result().Cookies() {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	user := models.User{}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		panic(err)
	}

	assert.Equal(t, mockName, user.Name)
	assert.Equal(t, mockEmail, user.Email)
}

func createMockUser() error {
	// create the user
	hash, err := auth.HashPassword(mockPassword)
	if err != nil {
		return err
	}

	return db.Create(&models.User{
		Name:     mockName,
		Email:    mockEmail,
		Password: hash,
	}).Error
}

func login(form auth.Credentials) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body, _ := json.Marshal(form)
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	router.ServeHTTP(w, req)
	return w
}
