package httpgin

import (
	"net/http"
	"spoti/internal/domain/album"
	"spoti/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	srv *service.AlbumService
}

func NewAlbumHandler(srv *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{srv: srv}
}
func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var req album.CreateAlbumReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.srv.CreateAlbum(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create album",
		})
		return
	}

	c.Status(http.StatusCreated)
}
func (h *AlbumHandler) GetAlbum(c *gin.Context) {
	albumID := c.Param("id")

	data, err := h.srv.GetAlbum(c.Request.Context(), albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get album",
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
func (h *AlbumHandler) GetAlbumsByIds(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	data, err := h.srv.GetAlbumsByIds(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get albums",
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
func (h *AlbumHandler) GetAlbumTracks(c *gin.Context) {
	albumID := c.Param("id")

	data, err := h.srv.GetAlbumsTracks(c.Request.Context(), albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get album tracks",
		})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
func (h *AlbumHandler) GetUserSavedAlbums(c *gin.Context) {
	userID := c.Param("userId")

	albums, err := h.srv.GetUserSavedAlbums(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get saved albums",
		})
		return
	}

	c.JSON(http.StatusOK, albums)
}
func (h *AlbumHandler) SaveAlbumsForCurrentUser(c *gin.Context) {
	var req struct {
		UserID string   `json:"user_id"`
		IDs    []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.SaveAlbumsForCurrentUser(c.Request.Context(), req.IDs, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save albums",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *AlbumHandler) RemoveAlbumsFromCurrentUser(c *gin.Context) {
	var req struct {
		UserID string   `json:"user_id"`
		IDs    []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.RemoveAlbumsFromCurrentUser(c.Request.Context(), req.IDs, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove albums",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *AlbumHandler) CheckUsersSavedAlbums(c *gin.Context) {
	var req struct {
		UserID string   `json:"user_id"`
		IDs    []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result, err := h.srv.CheckUsersSavedAlbums(c.Request.Context(), req.IDs, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check albums",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
func (h *AlbumHandler) GetNewReleases(c *gin.Context) {
	limitStr := c.Query("limit")

	limit := 10
	if limitStr != "" {
		parsed, err := strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid limit",
			})
			return
		}
		limit = parsed
	}

	albums, err := h.srv.GetNewReleases(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get new releases",
		})
		return
	}

	c.JSON(http.StatusOK, albums)
}
