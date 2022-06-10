package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/urionafacu/api-rest-websockets-go/models"
	"github.com/urionafacu/api-rest-websockets-go/repository"
	"github.com/urionafacu/api-rest-websockets-go/server"
)

type InsertPostRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Content string `json:"content"`
}

func InsertPostHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			user, _ := repository.GetUserById(r.Context(), claims.UserId)
			var postRequest = InsertPostRequest{}
			if err := json.NewDecoder(r.Body).Decode(&postRequest); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			post := models.Post{
				UserId:  user.Id,
				Content: postRequest.Content,
			}

			err = repository.InsertPost(r.Context(), &post)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			json.NewEncoder(w).Encode(PostResponse{
				Content: post.Content,
			})

		} else {
			http.Error(w, "error", http.StatusInternalServerError)
			return
		}
	}
}
