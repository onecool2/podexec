package router

import "net/http"

func Index(w http.ResponseWriter, r *http.Request) {
	rr.HTML(w, http.StatusOK, "index", nil)
}
