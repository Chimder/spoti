package httpgin

import (
	"net/http"
	"time"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/handler/http/middleware"
	"github.com/Chimder/spoti/internal/service"
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
	hashPass, err := middleware.GeneratePass(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request password",
		})
		return
	}
	req.HashPassword = hashPass

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

type SingInReq struct {
	email    string `json:"email"`
	password string `json:"password"`
}

func (h *UserHandler) SingInUser(c *gin.Context) {
	var req SingInReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.srv.GetUserByEmail(c.Request.Context(), req.email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	_, err = middleware.ComparePass(req.password, user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	accessToken, err := middleware.CreateUserToken(user.Id.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	// refreshHash := hashToken(refreshToken)
	// if err := h.redis.Save(c.Request.Context(), refreshHash, user.ID, 30*24*time.Hour); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store refresh token"})
	// 	return
	// }

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int((30 * 24 * time.Hour).Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"AccessToken": accessToken,
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
