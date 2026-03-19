package httpgin

import (
	"net/http"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ListeningEventHandler struct {
	repo *clickhouse.ListeningEventRepo
}

func NewListeningEventHandler(repo *clickhouse.ListeningEventRepo) *ListeningEventHandler {
	return &ListeningEventHandler{repo: repo}
}

func (r *ListeningEventHandler) AddEvent(c *gin.Context) {
	var req user.ListeningEventReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := r.repo.AddListeningEvent(c.Request.Context(), req); err != nil {
		log.Error().Err(err).Msg("err add listening event to clickhouse")
		return
	}

	c.Status(http.StatusOK)
}
