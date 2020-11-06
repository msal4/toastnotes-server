package models

import "gorm.io/gorm"

// Note is the user notes model.
type Note struct {
	Model
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	UserID  string `json:"userId"`
}

// NoteRepository holds the notes actions.
type NoteRepository struct {
	*Repository
}

// NewNoteRepository creates a new note repo.
func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{Repository: &Repository{DB: db}}
}
