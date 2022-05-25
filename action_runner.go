package application

import (
	"context"
	"encoding/json"
	"github.com/pasztorpisti/qs"
	"io"
	"net/http"
	"net/url"
)

type ActionRunner struct {
	logger     Logger
	jsonWriter JsonResponseWriter
	router     Router
}

func NewActionRunner(logger Logger, jsonWriter JsonResponseWriter, router Router) *ActionRunner {
	return &ActionRunner{logger: logger, jsonWriter: jsonWriter, router: router}
}

func (j *ActionRunner) Run(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) (any, error),
	request any,
) {
	switch r.Method {
	case http.MethodGet, http.MethodDelete:
		j.runGet(w, r, action, request)
	case http.MethodPost:
		j.runPost(w, r, action, request)
	case http.MethodPut, http.MethodPatch:
		j.runPut(w, r, action, request)
	default:
		j.logger.Error(r.Context(), "unsupported http method "+r.Method)
	}
}

func (j *ActionRunner) runGet(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) (any, error),
	request any,
) {
	err := j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	if err == nil {
		err = j.fillRequestFromUrlValues(w, r, request, r.URL.Query())
	}

	if err != nil {
		return
	}

	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runPost(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) (any, error),
	request any,
) {
	var err error
	err = j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	if err == nil {
		if r.Header.Get("Content-Type") != "application/json" {
			err = j.fillRequestFromBody(w, r, request)
		} else {
			err = j.fillRequestFromUrlValues(w, r, request, r.PostForm)
		}
	}

	if err != nil {
		return
	}
	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runPut(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) (any, error),
	request any,
) {
	var err error

	err = j.fillRequestFromUrlValues(w, r, request, j.router.RouteParams(r))
	err = j.fillRequestFromBody(w, r, request)

	if err != nil {
		return
	}
	j.runAction(w, r, action, request)
}

func (j *ActionRunner) runAction(
	w http.ResponseWriter,
	r *http.Request,
	action func(ctx context.Context, request any) (any, error),
	request any,
) {
	response, err := action(r.Context(), request)

	if err != nil {
		j.jsonWriter.Error(w, r, 500, err)
	}
	j.jsonWriter.Success(w, r, 200, response)
}

func (j *ActionRunner) fillRequestFromBody(
	w http.ResponseWriter,
	r *http.Request,
	request any,
) error {
	if request == nil {
		return nil
	}
	if r.Header.Get("Content-Type") != "application/json" {
		return nil
	}
	var err error
	defer r.Body.Close()

	var body []byte

	body, err = io.ReadAll(r.Body)

	if err == nil && body != nil {
		err = json.Unmarshal(body, request)
	}
	if err != nil {
		j.logger.Error(r.Context(), "Wrong request decoding: "+err.Error())
		j.jsonWriter.Error(w, r, 400, err)
		return err
	}

	return nil
}

func (j *ActionRunner) fillRequestFromUrlValues(
	w http.ResponseWriter,
	r *http.Request,
	request any,
	values url.Values,
) error {
	if request == nil {
		return nil
	}
	err := qs.Unmarshal(request, values.Encode())
	if err != nil {
		j.logger.Error(r.Context(), "Wrong request decoding: "+err.Error())
		j.jsonWriter.Error(w, r, 400, err)
		return err
	}

	return nil
}
