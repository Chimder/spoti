package httpgin

import (
	"net/http"
	"spoti/internal/domain/user"
	"spoti/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {

	return &UserHandler{
		srv: s,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.CreateUserReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	id, err := h.srv.CreateUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	u, err := h.srv.GetUserById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) FollowPlaylist(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("playlistId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	if err := h.srv.FollowUserToPlaylist(c.Request.Context(), userID, playlistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to follow playlist",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UnfollowPlaylist(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	playlistID, err := uuid.Parse(c.Param("playlistId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist id"})
		return
	}

	if err := h.srv.UnfollowUserFromPlaylist(c.Request.Context(), userID, playlistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to unfollow playlist",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) FollowArtist(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	artistID, err := uuid.Parse(c.Param("artistId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}

	if err := h.srv.FollowUserToArtist(c.Request.Context(), userID, artistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to follow artist",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UnfollowArtist(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	artistID, err := uuid.Parse(c.Param("artistId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist id"})
		return
	}

	if err := h.srv.UnfollowUserFromArtist(c.Request.Context(), userID, artistID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to unfollow artist",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
