package httpgin

import (
	"net/http"
	"spoti/internal/domain/artist"
	"spoti/internal/service"

	"github.com/gin-gonic/gin"
)

type ArtistHandler struct {
	srv service.ArtistService
}

func NewArtistHandler(srv service.ArtistService) *ArtistHandler {
	return &ArtistHandler{srv: srv}
}

func (h *ArtistHandler) CreateArtist(c *gin.Context) {
	var req artist.CreateArtistReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.srv.CreateArtist(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create artist"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ArtistHandler) GetArtist(c *gin.Context) {
	artistID := c.Param("id")
	art, err := h.srv.GetArtist(c.Request.Context(), artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist"})
		return
	}

	c.JSON(http.StatusOK, art)
}

func (h *ArtistHandler) GetArtistsByIDs(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	artists, err := h.srv.GetArtistsByIDs(c.Request.Context(), req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artists"})
		return
	}

	c.JSON(http.StatusOK, artists)
}

func (h *ArtistHandler) GetArtistAlbums(c *gin.Context) {
	artistID := c.Param("id")
	albums, err := h.srv.GetArtistAlbums(c.Request.Context(), artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get artist albums"})
		return
	}

	c.JSON(http.StatusOK, albums)
}
