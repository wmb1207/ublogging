package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wmb1207/ublogging/internal/service"
)

func TokenAuthMiddleware(service service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		token := c.GetHeader("Authorization")

		if token == "" {
			// If the token is missing, return a 401 Unauthorized response
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			c.Abort() // Stop the execution of the next handlers
			return
		}

		// Here we should have a propper token. We are just going to use the user uuid as the token... just to make it
		// easier

		user, err := service.User(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing not found", "message": "Unauthorized"})
		}

		c.Set("user", user)
		c.Next()
	}
}
