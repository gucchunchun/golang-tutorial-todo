package recovermw

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

/**
* Reference: O'REILLY「実用GO言語」10.5.3 p.246
* panic発生時でもサーバープロセスが終了せず、別のリクエストを受け続けられるようにする
 */
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// TODO: DBロールバックの実装
				log.Println(err)
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
