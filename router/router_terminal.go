package router

import (
	"example/podexec/backend"

	"github.com/pressly/chi"
)

func TerminalsRouter(b *backend.Backend) (r chi.Router) {
	r = chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", warpHandler(b, terminalsHandler))
		r.Post("/:pid", warpHandler(b, terminalsHandler))
		// r.Get("/:namespace/:podName/:containerName", warpHandler(b, socketHandler))
		r.Get("/:namespace/:podName/:containerName", warpHandler(b, wsHandler))
	})
	return
}
