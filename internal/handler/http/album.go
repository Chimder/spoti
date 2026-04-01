package httpgin

import (
	"net/http"
	"strconv"

	"github.com/Chimder/spoti/internal/domain/album"
	"github.com/Chimder/spoti/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AlbumHandler struct {
	srv *service.AlbumService
}

func NewAlbumHandler(srv *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{srv: srv}
}
func (h *AlbumHandler) CreateAlbum(c *gin.Context) {
	var req album.CreateAlbumReq

	// log.Info().Str("PARAM", c.Param("id")).Msg("CREATE")
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
func (h *AlbumHandler) GetAlbumWithTracks(c *gin.Context) {
	albumID := c.Param("id")

	data, err := h.srv.GetAlbumWithTracks(c.Request.Context(), albumID)
	if err != nil {
		log.Error().Err(err).Msg("ERROR GET")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get album",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AlbumHandler) GetAlbumsByIds(c *gin.Context) {
	ids := c.QueryArray("id")

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids required"})
		return
	}

	data, err := h.srv.GetAlbumsByIds(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get albums",
		})
		return
	}

	c.JSON(http.StatusOK, data)
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
	userID := c.Param("userId")
	var req struct {
		IDs []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album_ids required"})
		return
	}

	if err := h.srv.SaveAlbumsForCurrentUser(c.Request.Context(), req.IDs, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save albums"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AlbumHandler) RemoveAlbumsFromCurrentUser(c *gin.Context) {
	userID := c.Param("userId")
	var req struct {
		IDs []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album_ids required"})
		return
	}

	if err := h.srv.RemoveAlbumsFromCurrentUser(c.Request.Context(), req.IDs, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove albums",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
func (h *AlbumHandler) CheckUserSavedAlbums(c *gin.Context) {
	userID := c.Param("userId")
	var req struct {
		IDs []string `json:"album_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album_ids required"})
		return
	}

	result, err := h.srv.CheckUsersSavedAlbums(c.Request.Context(), req.IDs, userID)
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
