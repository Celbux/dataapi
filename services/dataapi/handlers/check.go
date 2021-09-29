package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/Celbux/dataapi/foundation/web"
)

type check struct {}

// readiness simply returns a 200 ok when called
func (c check) readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	status := struct{ Status string }{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)

}

// liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (c check) liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := struct {
		Status    string `json:"Status,omitempty"`
		Host      string `json:"Host,omitempty"`
		Pod       string `json:"Pod,omitempty"`
		PodIP     string `json:"PodIP,omitempty"`
		Node      string `json:"Node,omitempty"`
		Namespace string `json:"Namespace,omitempty"`
	}{
		Status:    "up",
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	return web.Respond(ctx, w, info, http.StatusOK)

}
