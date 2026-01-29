package responder

import (
	"encoding/json"
	"net/http"

	"github.com/ktruedat/llm-feedback-analysis/pkg/log"
)

type RestResponder interface {
	RespondContent(resp http.ResponseWriter, data Response, opts ...ResponseOption)
}

type Response interface {
	IsResponse()
	HTTPCode() int
}

type responder struct {
	logger log.Logger
}

func NewRestResponder(logger log.Logger) RestResponder {
	return &responder{logger: logger.NewGroup("http_responder")}
}

func (r *responder) RespondContent(resp http.ResponseWriter, data Response, opts ...ResponseOption) {
	resp.Header().Set("Content-Type", "application/json")
	if data == nil {
		r.logger.Warning("no data to respond with")
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	respOpts := &respOpts{
		statusCode: data.HTTPCode(),
	}

	for _, opt := range opts {
		opt(respOpts)
	}

	resp.WriteHeader(respOpts.statusCode)
	if err := json.NewEncoder(resp).Encode(data); err != nil {
		r.logger.Error("failed to encode response body", err)
	}
}
