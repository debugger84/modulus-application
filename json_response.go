package application

import (
	"encoding/json"
	"net/http"
)

type JsonResponseWriter interface {
	Success(w http.ResponseWriter, r *http.Request, response ActionResponse)
	Error(w http.ResponseWriter, r *http.Request, response ActionResponse)
}

type DefaultJsonResponseWriter struct {
	logger Logger
	config *Config
}

func NewJsonResponseWriter(logger Logger, config *Config) JsonResponseWriter {
	return &DefaultJsonResponseWriter{logger: logger, config: config}
}

func (j *DefaultJsonResponseWriter) Success(w http.ResponseWriter, r *http.Request, response ActionResponse) {
	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(response.Response)
	if err != nil {
		ctx := r.Context()
		j.logger.Error(ctx, "Error happened in JSON marshal. Err: %s", err)
	}
	_, _ = w.Write(jsonResp)
	return
}

func (j *DefaultJsonResponseWriter) Error(w http.ResponseWriter, r *http.Request, response ActionResponse) {
	w.WriteHeader(response.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string]interface{})
	resp["error"] = "Unknown error"
	if response.Error != nil {
		resp["error"] = response.Error.Error()
	}

	if len(response.Error.ValidationErrors) > 0 {
		vErrors := make(map[string]string)
		for _, validationError := range response.Error.ValidationErrors {
			vErrors[validationError.Field] = validationError.Err
		}
		resp["errors"] = vErrors
	}

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
