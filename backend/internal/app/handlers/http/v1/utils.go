package v1

import (
	"errors"
	"fmt"
	"net/http"

	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
)

// handleSvcError handles service errors and converts them to appropriate HTTP responses.
func (h *Handlers) handleSvcError(resp http.ResponseWriter, err error) {
	var appErr ce.ApplicationError
	if errors.As(err, &appErr) {
		h.responder.RespondContent(resp, appErr)
		return
	}

	h.logger.Warning("service error is not application error", "error", err)
	h.responder.RespondContent(resp, ce.ErrInternal(err))
}

// parseInt is a helper function to parse integer from string.
func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
