package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/models"
	"github.com/stretchr/testify/assert"
)

var (
	mockID      = "dadbc71a-7df7-4267-938e-a62407bc1bd5"
	mockTitle   = "test title"
	mockContent = "test content"
)

func TestListNotes(t *testing.T) {
	t.Cleanup(cleanup)

	user, _ := createMockUser(nil)

	firstTitle := "first title"
	secondTitle := "second title"

	db.Create(&models.Note{Title: firstTitle, Content: mockContent, UserID: user.ID})
	db.Create(&models.Note{Title: secondTitle, Content: mockContent, UserID: user.ID})

	wLogin := login(mockUserCreds)
	w := serveHTTP("GET", API+APINote, nil, wLogin.Result().Cookies())
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, firstTitle)
	assert.Contains(t, body, secondTitle)
}

func TestNoteRetrieve(t *testing.T) {
	t.Cleanup(cleanup)

	user, _ := createMockUser(nil)

	t.Run("a_user_can_not_access_a_non_existing_note", func(t *testing.T) {
		wLogin := login(mockUserCreds)
		w := serveHTTP("GET", API+APINote+"/"+mockID, nil, wLogin.Result().Cookies())
		assert.Equal(t, http.StatusNotFound, w.Code)
		body := w.Body.String()
		assert.NotContains(t, body, mockTitle)
		assert.NotContains(t, body, mockContent)
	})

	t.Run("a_user_can_retrieve_their_note", func(t *testing.T) {
		wLogin := login(mockUserCreds)
		createMockNote(user.ID)
		// should find the note.
		w := serveHTTP("GET", API+APINote+"/"+mockID, nil, wLogin.Result().Cookies())
		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		assert.Contains(t, body, mockTitle)
		assert.Contains(t, body, mockContent)
	})

	t.Run("a_user_can_access_their_notes_only", func(t *testing.T) {
		anotherUserCreds := auth.Credentials{
			Email:    "a" + mockEmail,
			Password: mockPassword,
		}
		createMockUser(&anotherUserCreds)
		wLogin := login(anotherUserCreds)

		w := serveHTTP("GET", API+APINote+"/"+mockID, nil, wLogin.Result().Cookies())
		assert.Equal(t, http.StatusNotFound, w.Code)
		body := w.Body.String()
		assert.NotContains(t, body, mockTitle)
		assert.NotContains(t, body, mockContent)
	})
}

func TestCreateNote(t *testing.T) {
	t.Cleanup(cleanup)
	user, _ := createMockUser(nil)
	wLogin := login(mockUserCreds)

	t.Run("a_user_can_create_a_note", func(t *testing.T) {
		defer cleanup()

		body, _ := json.Marshal(models.Note{
			Title:   mockTitle,
			Content: mockContent,
		})
		w := serveHTTP("POST", API+APINote, bytes.NewReader(body), wLogin.Result().Cookies())

		assert.Equal(t, http.StatusOK, w.Code)
		note := models.Note{}
		assert.Nil(t, db.First(&note).Error)
		assert.NotEmpty(t, note.ID)
		assert.Equal(t, mockTitle, note.Title)
		assert.Equal(t, mockContent, note.Content)
		assert.Equal(t, user.ID, note.UserID)
	})

	t.Run("a_user_can_not_create_a_note_using_invalid_data", func(t *testing.T) {
		body, _ := json.Marshal(models.Note{Content: mockContent})
		w := serveHTTP("POST", API+APINote, bytes.NewReader(body), wLogin.Result().Cookies())

		assert.Equal(t, http.StatusNotAcceptable, w.Code)
		note := models.Note{}
		assert.NotNil(t, db.First(&note).Error)
		assert.Empty(t, note.ID)
		assert.Empty(t, note.Title)
		assert.Empty(t, note.Content)
	})
}

func TestUpdateNote(t *testing.T) {
	t.Cleanup(cleanup)

	user, _ := createMockUser(nil)
	wLogin := login(mockUserCreds)
	note := createMockNote(user.ID)

	newTitle := "new title"
	newContent := "new content"
	body, _ := json.Marshal(models.Note{
		Title:   newTitle,
		Content: newContent,
	})

	w := serveHTTP("PUT", API+APINote+"/"+note.ID, bytes.NewReader(body), wLogin.Result().Cookies())

	assert.Equal(t, http.StatusOK, w.Code)

	n := models.Note{}
	assert.Nil(t, db.First(&n).Error)
	assert.NotEmpty(t, n.ID)
	assert.Equal(t, newTitle, n.Title)
	assert.Equal(t, newContent, n.Content)
}

func TestDeleteNote(t *testing.T) {
	t.Cleanup(cleanup)
	user, _ := createMockUser(nil)
	wLogin := login(mockUserCreds)
	note := createMockNote(user.ID)

	w := serveHTTP("DELETE", API+APINote+"/"+note.ID, nil, wLogin.Result().Cookies())

	assert.Equal(t, http.StatusOK, w.Code)

	n := models.Note{}
	assert.NotNil(t, db.First(&n).Error)
	assert.Empty(t, n.ID)
	assert.Empty(t, n.Title)
	assert.Empty(t, n.Content)
}

func createMockNote(userID string) *models.Note {
	note := models.Note{
		Model: models.Model{
			ID: mockID,
		},
		Title:   mockTitle,
		Content: mockContent,
		UserID:  userID,
	}
	db.Create(&note)
	return &note
}
