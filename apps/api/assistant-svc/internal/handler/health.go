package handler

import (
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/pkg/response"
)

func Health(w http.ResponseWriter, r *http.Request) {
	response.OK(w, map[string]any{
		"status":  "ok",
		"service": "assistant-svc",
	})
}
