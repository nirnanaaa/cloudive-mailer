package httpd

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nirnanaaa/cloudive-mailer/services/kafka"
)

func (h *Handler) acceptInboundEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var msg kafka.InboundEmailEvent
	if err := decoder.Decode(&msg); err != nil {
		h.httpError(w, "err.global.invalid_payload", http.StatusBadRequest)
		return
	}
	h.Kafka.QueueMail(&msg)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) httpError(w http.ResponseWriter, error string, code int) {
	response := Response{Err: errors.New(error)}
	if rw, ok := w.(ResponseWriter); ok {
		h.writeHeader(w, code)
		rw.WriteResponse(response)
		return
	}
}
