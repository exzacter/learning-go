package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) TestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Test curl has worked, handler is working for that function",
		}

		// encode the response as JSON and send back to client
		json.NewEncoder(w).Encode(response)
	}
}
