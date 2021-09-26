package dataapi

import (
	"context"
	"github.com/Celbux/dataapi/business/i"
	"io"
)

// DataAPI is the interface that wraps the DataAPIs functionality
//
// Evaluate runs a data code block file found in configs/dataapi
// A Data code block is code written outside of the code base
// This allows custom code to be created without the need to redeploy
//
// GetURL will retrieve core's API URL that the request will be sent to based on the routeNumber
// The routeNumber is used do uniquely identify a route
// Eg, 1 range voucher numbers, voucherPrefixes
type DataAPI interface {
	Evaluate(ctx context.Context, file io.Reader) error
	GetURL(ctx context.Context) error
}

// DataAPIService encapsulates all dependencies required by the DataAPI.
// This service is used to run data driven functionality at run time
type DataAPIService struct {
	EvalCache EvalCache
	Log       i.Logger
}

type EvalCache map[string]interface{}
