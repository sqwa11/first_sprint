package get

import (
	"net/http"
	"strings"

	"github.com/sqwa11/first_sprint/internal/app/post"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	longURL, exists := post.URLMap[id]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
