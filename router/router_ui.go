package router

import (
	"example/podexec/backend"

	"github.com/pressly/chi"
)

func UIRouter(b *backend.Backend) (r chi.Router) {
	r = chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/terminal", warpHandler(b, openTerminalPage))
		r.Get("/terminal", warpHandler(b, openTerminalPage))
		r.Get("/terminal/:namespace/:podName/:containerName", warpHandler(b, openTerminal))
	})
	return
}
