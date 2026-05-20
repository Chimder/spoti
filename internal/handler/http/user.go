package httpgin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/handler/http/middleware"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
	"github.com/Chimder/spoti/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	srv   *service.UserService
	redis *rediscache.RedisCache
}

func NewUserHandler(s *service.UserService, redis *rediscache.RedisCache) *UserHandler {
	return &UserHandler{srv: s, redis: redis}
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

	id, err := h.srv.CreateUser(c.Request.Context(), req, hashPass)
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

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) SignInUser(c *gin.Context) {
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.srv.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email",
		})
		return
	}

	_, err = middleware.ComparePass(req.Password, user.PasswordHash)
	if err != nil {
		log.Error().Err(err).Msg("comparePass")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid comparePass",
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

	refreshHash := middleware.HashToken(refreshToken)
	if err := h.redis.Set(c.Request.Context(), refreshHash, user.Id, 30*24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed set refresh token"})
		return
	}

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

func (h *UserHandler) RefreshUserToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	hash := middleware.HashToken(refreshToken)

	var userId string
	err = h.redis.Get(c.Request.Context(), hash, &userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	fmt.Printf("GET userid from redis %s", userId)
	accessToken, err := middleware.CreateUserToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access"})
		return
	}

	newRefresh, err := middleware.GenerateRefreshToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed refresh"})
		return
	}

	newHash := middleware.HashToken(newRefresh)

	_ = h.redis.Delete(c.Request.Context(), hash)

	_ = h.redis.Set(c.Request.Context(), newHash, userId, 30*24*time.Hour)

	c.SetCookie(
		"refresh_token",
		newRefresh,
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
