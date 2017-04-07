package router

import (
	"example/podexec/backend"
	"fmt"
	"net/http"

	"github.com/pressly/chi/render"
)

type StatusFlag string

type handler func(*backend.Backend, http.ResponseWriter, *http.Request)

const (
	FailedFlag  StatusFlag = "Failed"
	SuccessFlag StatusFlag = "Success"
)

type APIResponse struct {
	Status     StatusFlag  `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func warpHandler(b *backend.Backend, hd handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hd == nil {
			resp := &APIResponse{Status: FailedFlag}
			resp.Message = fmt.Sprintf("api_router.go::warpHandler::No Handler for warpHandler")
			render.Status(r, http.StatusNotImplemented)
			render.JSON(w, r, resp)
			return
		}

		hd(b, w, r)
	})
}
