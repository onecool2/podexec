package router

import "github.com/pressly/chi"
import "example/podexec/backend"

func APIRouter(b *backend.Backend) (r chi.Router) {
	r = chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/namespaces", func(r chi.Router) {
			r.Get("/", warpHandler(b, getNamespaces))
			r.Get("/:namespace", warpHandler(b, getPods))
			r.Get("/:namespace/:podName", warpHandler(b, getContainers))
		})
	})
	return
}
