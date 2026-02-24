package httpgin

import (
	"net/http"
	"spoti/internal/domain/track"
	"spoti/internal/service"

	"github.com/gin-gonic/gin"
)

type TrackHandler struct {
	srv service.TrackService
}

func NewTrackHandler(srv service.TrackService) *TrackHandler {
	return &TrackHandler{srv: srv}
}

func (h *TrackHandler) CreateTrack(c *gin.Context) {
	var req track.CreateTrackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.CreateTrack(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create track"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *TrackHandler) GetTrackById(c *gin.Context) {
	trackID := c.Param("id")
	tr, err := h.srv.GetTrackById(c.Request.Context(), trackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get track"})
		return
	}

	c.JSON(http.StatusOK, tr)
}

func (h *TrackHandler) GetTracksByIds(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	tracks, err := h.srv.GetTracksByIds(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

func (h *TrackHandler) GetArtistTracks(c *gin.Context) {
	artistID := c.Param("artistId")
	tracks, err := h.srv.GetArtistTracks(c.Request.Context(), artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist tracks"})
		return
	}

	c.JSON(http.StatusOK, tracks)
}
