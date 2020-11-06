package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/msal4/toastnotes/auth"
	"github.com/msal4/toastnotes/models"
	"github.com/msal4/toastnotes/utils"
	"gorm.io/gorm"
)

// NoteController is the group of the set of actions related to user notes with their dependencies.
type NoteController struct {
	Repository *models.NoteRepository
}

// NewNoteController creates a new note controller.
func NewNoteController(db *gorm.DB) *NoteController {
	return &NoteController{Repository: models.NewNoteRepository(db)}
}

// Retrieve gets the first note matching the provided id.
func (ctrl *NoteController) Retrieve(c *gin.Context) {
	noteID := c.Param("id")
	userID := c.GetString(auth.UserIDKey)

	note := models.Note{}
	if err := ctrl.Repository.FindByID(&note, noteID); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("Note not found"))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Could not handle your request"))
		return
	}

	if note.UserID != userID {
		c.AbortWithStatusJSON(http.StatusNotFound, utils.Err("You don't own this note"))
		return
	}

	c.JSON(http.StatusOK, note)
}

// List handles getting the authenticated user notes.
func (ctrl *NoteController) List(c *gin.Context) {
	userID := c.GetString(auth.UserIDKey)

	notes := []models.Note{}
	err := ctrl.Repository.DB.Scopes(models.Paginate(c)).Select("ID", "Title", "CreatedAt", "UpdatedAt").
		Find(&notes, "user_id = ?", userID).Order("updated_at DESC").Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Failed to retrieve notes"))
		return
	}

	c.JSON(http.StatusOK, notes)
}

// Create handles creating notes.
func (ctrl *NoteController) Create(c *gin.Context) {
	note := models.Note{}
	if errs := shouldBindJSON(c, &note); errs != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, errs)
		return
	}

	note.UserID = c.GetString(auth.UserIDKey)

	if err := ctrl.Repository.DB.Create(&note).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Could not create note :("))
		return
	}

	c.JSON(http.StatusOK, note)
}

// Update handles updating notes.
func (ctrl *NoteController) Update(c *gin.Context) {
	note := models.Note{}
	if errs := shouldBindJSON(c, &note); errs != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, errs)
		return
	}

	note.ID = c.Param("id")
	note.UserID = c.GetString(auth.UserIDKey)

	if err := ctrl.Repository.DB.Model(&note).Updates(note).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Could not update note :("))
		return
	}

	c.JSON(http.StatusOK, note)
}

// Delete handles deleting notes.
func (ctrl *NoteController) Delete(c *gin.Context) {
	note := models.Note{}
	note.ID = c.Param("id")
	note.UserID = c.GetString(auth.UserIDKey)

	if err := ctrl.Repository.DB.Delete(&note).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Err("Could not delete the note :("))
		return
	}

	c.JSON(http.StatusOK, utils.Msg("Note removed"))
}
