package GET

import (
	"github.com/sqwa11/first_sprint/internal/app/POST"
	"net/http"
	"strings"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	longURL, exists := POST.UrlMap[id]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}
