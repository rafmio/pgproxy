package handlers

import (
	"encoding/json"
	"net/http"
	"pgproxy/internal/utils"
)

// The healthCheck function is designed to check the server's availability and readiness via an HTTP endpoint.
// It can be used in various scenarios such as:
//  1. Liveness Probes: Ensures that the server is running and responding to requests,
//     which is useful for container orchestrators like Kubernetes.
//  2. Readiness Probes: Checks if the server is ready to handle incoming traffic by verifying its internal state.
//  3. Load Balancing & Monitoring: Used by load balancers or monitoring systems to ensure the server is available.
//
// Example usage:
// In Kubernetes, periodic checks are made against the `/health` endpoint. If the server responds with status code 200 OK,
// it is considered healthy and ready to serve traffic.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
