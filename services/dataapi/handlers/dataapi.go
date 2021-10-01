package handlers

import (
	"context"
	"fmt"
	"github.com/Celbux/dataapi/business/dataapi"
	"github.com/Celbux/dataapi/foundation/tools"
	"github.com/Celbux/dataapi/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type DataAPIHandlers struct {
	Service dataapi.DataAPIService
}

func (d DataAPIHandlers) evaluateHandler(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {

	failures, successes, err := d.evaluate(ctx, r)
	if err != nil {
		if _, ok := errors.Cause(err).(*dataapi.Error); ok {

			response := struct {
				Error string `json:"Error"`
			}{
				Error: err.Error(),
			}

			return web.Respond(ctx, w, response, http.StatusBadRequest)
		}
		return err
	}

	response := struct {
		Failures []string
		Successes []string
	}{
		Failures: failures,
		Successes: successes,
	}
	return web.Respond(ctx, w, response, http.StatusOK)

}

// evaluate will run the data code found in the given file
// The file must exist as a relative path to the running server
func (d DataAPIHandlers) evaluate(ctx context.Context, r *http.Request, ) ([]string, []string, error) {

	// Get file name from request body
	// The file contains the data code we want to run live
	type request struct {
		File string `json:"File"`
	}
	req := request{}
	err := web.Decode(r, &req)
	if err != nil {
		return nil, nil, &dataapi.Error{Err: errors.New(fmt.Sprintf("error evaluate/web.Decode: %v", err.Error()))}
	}

	// Create the eval cache the service
	d.Service.EvalCache = make(map[string]interface{})
	d.Service.EvalCache["dataapi"] = &d.Service

	// Evaluate all expressions in input filename
	resultsRaw := d.Service.Evaluate(req.File)
	report, ok := resultsRaw["report"].(*tools.Tree)
	if !ok {
		return nil, nil, errors.Errorf("evaluate fatal: %v", resultsRaw)
	}

	// Return the test results
	failures, successes, err := d.Service.GetResults(*report)
	if err != nil {
		return nil, nil, err
	}

	if len(failures) == 0 {
		successes = append(successes, fmt.Sprintf("%v: passed", req.File))
	}
	return failures, successes, nil

}
