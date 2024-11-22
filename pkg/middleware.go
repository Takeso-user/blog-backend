package pkg

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			log.Println("Missing token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		log.Printf("Token valid for user: %s, role: %s", claims.Username, claims.Role)
		c.Next()
	}
}

func OwnerOrAdminMiddleware(postService *PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, roleExists := c.Get("role")
		username, usernameExists := c.Get("username")
		if !roleExists || !usernameExists {
			log.Println("Missing or invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		if role == "Admin" {
			log.Println("Admin access granted")
			c.Next()
			return
		}

		postID := c.Param("id")
		post, err := postService.GetPostById(postID)
		if err != nil {
			log.Printf("Unable to fetch post: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch post"})
			c.Abort()
			return
		}

		if post.AuthorID != username {
			log.Printf("User %s does not have permission to delete post %s", username, postID)
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this post"})
			c.Abort()
			return
		}

		log.Printf("User %s granted permission to delete post %s", username, postID)
		c.Next()
	}
}
