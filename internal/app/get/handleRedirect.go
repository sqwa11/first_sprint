package get

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sqwa11/first_sprint/internal/app/post"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	longURL, ok := post.URLMap[id]
	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
