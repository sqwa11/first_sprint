package get

import (
	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/internal/app/post"
	"net/http"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	longURL, exists := post.URLMap[id]
	if !exists {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
