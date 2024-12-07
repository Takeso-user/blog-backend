package pkg

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			logrus.Println("Missing token")
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			logrus.Printf("Invalid token: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)
		logrus.Printf("Token valid for user: %s, role: %s", claims.Username, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OwnerOrAdminMiddleware(postService *PostService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, roleExists := r.Context().Value("role").(string)
			username, usernameExists := r.Context().Value("username").(string)
			if !roleExists || !usernameExists {
				logrus.Println("Missing or invalid token")
				http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
				return
			}

			if role == "Admin" {
				logrus.Println("Admin access granted")
				next.ServeHTTP(w, r)
				return
			}

			postID := mux.Vars(r)["id"]
			post, err := postService.GetPostById(postID)
			if err != nil {
				logrus.Printf("Unable to fetch post: %v", err)
				http.Error(w, "Unable to fetch post", http.StatusInternalServerError)
				return
			}

			if post.AuthorID != username {
				logrus.Printf("User %s does not have permission to delete post %s", username, postID)
				http.Error(w, "You do not have permission to delete this post", http.StatusForbidden)
				return
			}

			logrus.Printf("User %s granted permission to delete post %s", username, postID)
			next.ServeHTTP(w, r)
		})
	}
}
