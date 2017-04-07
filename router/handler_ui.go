package router

import (
	"example/podexec/backend"
	"fmt"
	"net/http"

	"github.com/pressly/chi"
)

func openTerminalPage(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	namespace := r.FormValue("namespace")
	podName := r.FormValue("podName")
	containerName := r.FormValue("containerName")

	data := struct {
		Namespace     string
		PodName       string
		ContainerName string
	}{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
	}

	fmt.Println("namespace", namespace, "podName", podName, "containerName", containerName, data)
	rr.HTML(w, http.StatusOK, "terminal", &data)
}

func openTerminal(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	podName := chi.URLParam(r, "podName")
	containerName := chi.URLParam(r, "containerName")

	data := struct {
		Namespace     string
		PodName       string
		ContainerName string
	}{
		Namespace:     namespace,
		PodName:       podName,
		ContainerName: containerName,
	}

	fmt.Println("namespace", namespace, "podName", podName, "containerName", containerName, data)
	rr.HTML(w, http.StatusOK, "terminal", &data)
}
