package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ktruedat/llm-feedback-analysis/internal/app/handlers/http/middleware"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/infrastructure/jwt"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
)

type requestConstraint interface {
	requests.CreateFeedbackRequest
}

type request[T requestConstraint] struct {
	Claims *jwt.Claims
	Data   T
}

func parsePayloadData[T requestConstraint](
	r *http.Request,
	body io.ReadCloser,
) (*request[T], error) {
	var req request[T]
	if err := json.NewDecoder(body).Decode(&req.Data); err != nil {
		return nil, fmt.Errorf("error decoding request body: %w", err)
	}

	if err := body.Close(); err != nil {
		return &req, fmt.Errorf("error closing body: %w", err)
	}

	req.Claims = middleware.GetUserClaims(r)

	return &req, nil
}

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
