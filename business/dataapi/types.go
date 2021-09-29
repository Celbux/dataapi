package dataapi

import (
	"github.com/Celbux/dataapi/business/i"
)

// DataAPIService encapsulates all dependencies required by the DataAPI
// This service is used to run data driven functionality at run time
type DataAPIService struct {
	EvalCache EvalCache
	Log       i.Logger
}

type EvalCache map[string]interface{}
