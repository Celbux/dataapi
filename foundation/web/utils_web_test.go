package web_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Celbux/dataapi/foundation/web"
)

// Test Handlers
func tenHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	fmt.Fprint(w, 10)
	return nil
}

func errorHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	return errors.New("error")
}

func contextHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	_, ok := ctx.Value(web.KeyValues).(*web.Values)

	if !ok {
		fmt.Fprint(w, "false")
		return nil
	}
	fmt.Fprint(w, "true")
	return nil
}

const (
	appMidVal     = "1"
	handlerMidVal = "2"
	centerVal     = "3"
)

func centerHandler(ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	fmt.Fprint(w, centerVal)
	return nil
}

func appMiddleware(h web.Handler) web.Handler {
	handler := func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		fmt.Fprint(w, appMidVal)
		_ = h(ctx, w, r)
		fmt.Fprint(w, appMidVal)
		return nil
	}
	return handler
}

func handlerMiddleware(h web.Handler) web.Handler {
	handler := func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		fmt.Fprint(w, handlerMidVal)
		_ = h(ctx, w, r)
		fmt.Fprint(w, handlerMidVal)
		return nil
	}
	return handler
}

func trustedErrHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	trustedErr := web.NewRequestError(errors.New("Trusted Error"), http.StatusNotFound)
	return web.RespondError(ctx, w, trustedErr)
}

func untrustedErrHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	untrustedErr := errors.New("Untrusted Error")
	return web.RespondError(ctx, w, untrustedErr)
}

func noContentHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	return web.Respond(ctx, w, "dummy value", http.StatusNoContent)
}

func contentHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	return web.Respond(ctx, w, "dummy value", http.StatusOK)
}

func paramCheckHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	params := web.Params(r)
	val := params["paramCheck"]

	return web.Respond(ctx, w, val, http.StatusOK)
}

func requestBodyHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	body := struct {
		Response  string `json:"Response" validate:"required"`
		Response2 string `json:"Response2" validate:"required"`
	}{}

	err := web.Decode(r, &body)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, body, http.StatusOK)
}

// Assert Functions
func assertString(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Errorf("wanted %q but got %q", want, got)
	}
}

func assertInt(t *testing.T, want, got int) {
	t.Helper()
	if got != want {
		t.Errorf("wanted %d but got %d", want, got)
	}
}
