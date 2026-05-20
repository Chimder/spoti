package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {

	t.Run("TestGoodToken", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.Use(AuthUserMiddleware())

		r.GET("/me", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.NotNil(t, userID)

			c.Status(200)
		})

		token, err := CreateUserToken("testpass")
		require.NoError(t, err)

		req := httptest.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("TestBadToken", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.New()

		r.Use(AuthUserMiddleware())
		r.GET("/me", func(c *gin.Context) {
			c.Status(200)
		})

		req := httptest.NewRequest("GET", "/me", nil)
		req.Header.Set("Authorization", "Bearer badToken")

		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
