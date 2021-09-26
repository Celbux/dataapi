package web

import (
	"bytes"
	"encoding/json"
	"testing"
)

// Encode is a convenience function for testing http requests.
// WARNING, there is no error handling for Marshaling
func Encode(data interface{}) *bytes.Reader {
	jsonData, _ := json.Marshal(data)
	return bytes.NewReader(jsonData)
}

// AssertCode is a convenience function for testing http requests. You can use
// it to assert a specific status code was returned.
func AssertCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected http code %d but got %d", want, got)
	}
}
