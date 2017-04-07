package router

import (
	"example/podexec/backend"
	"net/http"

	"github.com/pressly/chi"
)

func apiV1Test(w http.ResponseWriter, r *http.Request) {
	rr.JSON(w, http.StatusOK, map[string]string{"message": "Hello world"})
}

func getNamespaces(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	namespaces, err := b.GetNamespaces()
	resp := &APIResponse{Status: SuccessFlag}

	if err != nil {
		resp.Status = FailedFlag
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = err.Error()
		rr.JSON(w, http.StatusOK, resp)
		return
	}

	namespacesName := make([]string, 0)
	for _, namespace := range namespaces {
		namespacesName = append(namespacesName, namespace.Name)
	}

	resp.Data = namespacesName
	resp.StatusCode = http.StatusOK
	rr.JSON(w, http.StatusOK, resp)
}

func getPods(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	pods, err := b.GetPods(namespace)

	resp := &APIResponse{Status: SuccessFlag}
	if err != nil {
		resp.Status = FailedFlag
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = err.Error()
		rr.JSON(w, http.StatusOK, resp)
		return
	}

	podsName := make([]string, 0)
	for _, pod := range pods {
		podsName = append(podsName, pod.Name)
	}
	resp.Data = podsName
	resp.StatusCode = http.StatusOK
	rr.JSON(w, http.StatusOK, resp)
}

func getContainers(b *backend.Backend, w http.ResponseWriter, r *http.Request) {
	namespace := chi.URLParam(r, "namespace")
	podName := chi.URLParam(r, "podName")
	containers, err := b.GetContainers(namespace, podName)

	resp := &APIResponse{Status: SuccessFlag}
	if err != nil {
		resp.Status = FailedFlag
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = err.Error()
		rr.JSON(w, http.StatusOK, resp)
		return
	}

	containersName := make([]string, 0)
	for _, container := range containers {
		containersName = append(containersName, container.Name)
	}
	resp.Data = containersName
	resp.StatusCode = http.StatusOK
	rr.JSON(w, http.StatusOK, resp)
}
