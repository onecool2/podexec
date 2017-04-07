package main

import (
	"example/podexec/backend"
	"example/podexec/router"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/goware/cors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

var (
	c = cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	kubeConfig := "/etc/.kube/config"
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Use(middleware.Recoverer)
	root.Use(c.Handler)

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	log.Printf("[INFO] add static file path: %s", filesDir)
	root.FileServer("/static", http.Dir(filesDir))

	b := backend.New(kubeConfig)
	root.Route("/", func(r chi.Router) {
		r.Get("/", router.Index)
		r.Mount("/ui", router.UIRouter(b))
		r.Mount("/api/v1", router.APIRouter(b))
		r.Mount("/terminal", router.TerminalsRouter(b))
	})

	addr := ":3000"
	log.Printf("[INFO] Start Http Server at %s\n", addr)
	http.ListenAndServe(addr, root)
}
