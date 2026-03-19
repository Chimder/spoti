package httpgin

import (
	"net/http"
	"strconv"

	"github.com/Chimder/spoti/internal/domain/playlist"
	"github.com/Chimder/spoti/internal/service"
	"github.com/gin-gonic/gin"
)

type PlaylistHandler struct {
	srv *service.PlaylistService
}

func NewPlaylistHandler(srv *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{srv: srv}
}

func (h *PlaylistHandler) CreatePlaylist(c *gin.Context) {
	var req playlist.CreatePlaylistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.CreatePlaylist(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create playlist"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *PlaylistHandler) GetPlaylistById(c *gin.Context) {
	playlistID := c.Param("id")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil {
			limit = parsed
		}
	}

	if offsetStr != "" {
		if parsed, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsed
		}
	}

	data, err := h.srv.GetPlaylistById(c.Request.Context(), playlistID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get playlist"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *PlaylistHandler) AddToPlaylist(c *gin.Context) {
	playlistID := c.Param("playlistId")
	trackID := c.Param("trackId")

	if err := h.srv.AddToPlaylist(c.Request.Context(), playlistID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add track to playlist"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlaylistHandler) UpdatePlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	var req playlist.UpdatePlaylistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.UpdatePlaylist(c.Request.Context(), playlistID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update playlist"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlaylistHandler) DeleteFromPlaylist(c *gin.Context) {
	playlistID := c.Param("playlistId")
	trackID := c.Param("trackId")

	if err := h.srv.DeleteFromPlaylist(c.Request.Context(), playlistID, trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove track from playlist"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PlaylistHandler) GetAllUserPlaylists(c *gin.Context) {
	userID := c.Param("userId")
	data, err := h.srv.GetAllUserPlaylists(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user playlists"})
		return
	}

	c.JSON(http.StatusOK, data)
}
