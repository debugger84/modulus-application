package application

import (
	"encoding/json"
	"net/http"
)

type JsonResponseWriter interface {
	Success(w http.ResponseWriter, r *http.Request, statusCode int, response interface{})
	Error(w http.ResponseWriter, r *http.Request, statusCode int, err error)
}

type JsonResponse struct {
	logger Logger
	config *Config
}

func NewJsonResponse(logger Logger, config *Config) JsonResponseWriter {
	return &JsonResponse{logger: logger, config: config}
}

func (j *JsonResponse) Success(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(response)
	if err != nil {
		ctx := r.Context()
		j.logger.Error(ctx, "Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}

func (j *JsonResponse) Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]string)
	resp["error"] = err.Error()

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		ctx := r.Context()
		j.logger.Error(ctx, "Error happened in JSON marshal. Err: %s", err)
		_, _ = w.Write([]byte(`{"error": "Error happened in JSON marshal."}`))
		return
	}
	_, _ = w.Write(jsonResp)
	return
}
