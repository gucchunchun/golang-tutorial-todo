package recovermw

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]string{
						"code":    strconv.Itoa(http.StatusInternalServerError),
						"message": http.StatusText(http.StatusInternalServerError),
					},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
