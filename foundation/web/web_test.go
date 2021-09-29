package web_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/Celbux/dataapi/foundation/web"
)

// Framework Features
// =============================================================================
func TestRouting(t *testing.T) {
	t.Log("should return '10' on '/ten' route")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/ten", tenHandler)

	req, _ := http.NewRequest(http.MethodGet, "/ten", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "10"
	got := res.Body.String()

	assertString(t, want, got)
}

func Test404(t *testing.T) {
	t.Log("should return status [404] on '/' route")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := 404
	got := res.Result().StatusCode

	assertInt(t, want, got)
}

// Middleware works sort of like an onion, we expect the middleware to run as
// follows: appMid->handlerMid->routeHander->handlerMid->appMid
func TestMiddleware(t *testing.T) {
	t.Log("should wrap application middleware and route middleware around handlers")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel, appMiddleware)
	testApp.Handle(
		http.MethodGet,
		"/middleware",
		centerHandler,
		handlerMiddleware,
	)

	req, _ := http.NewRequest(http.MethodGet, "/middleware", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	// want "12321"
	want := appMidVal + handlerMidVal + centerVal + handlerMidVal + appMidVal
	got := res.Body.String()

	assertString(t, want, got)
}

func TestRespondNoContent(t *testing.T) {
	t.Log("should respond with status '202' and no content")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/noContent", noContentHandler)

	req, _ := http.NewRequest(http.MethodGet, "/noContent", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := ""
	got := res.Body.String()

	assertString(t, want, got)
	assertInt(t, http.StatusNoContent, res.Code)
}

func TestRespondWithContent(t *testing.T) {
	t.Log("should respond with status '200' and correct content")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/content", contentHandler)

	req, _ := http.NewRequest(http.MethodGet, "/content", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "\"dummy value\""
	got := res.Body.String()

	assertString(t, want, got)
	assertInt(t, http.StatusOK, res.Code)
}

func TestTrustedErrorResponse(t *testing.T) {
	t.Log("should return an actionalble error when the handler returns a trusted error")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/trustedErr", trustedErrHandler)

	req, _ := http.NewRequest(http.MethodGet, "/trustedErr", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "{\"Error\":\"Trusted Error\"}" //unmarhaled JSON string
	got := res.Body.String()

	assertString(t, want, got)
	assertInt(t, http.StatusNotFound, res.Code)
}

func TestUntrustedErrorResponse(t *testing.T) {
	t.Log("should return a boilerplate error when the handler returns an untrusted error")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/untrustedErr", untrustedErrHandler)

	req, _ := http.NewRequest(http.MethodGet, "/untrustedErr", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "{\"Error\":\"Internal Server Error\"}" //unmarhaled JSON string
	got := res.Body.String()

	assertString(t, want, got)
	assertInt(t, http.StatusInternalServerError, res.Code)
}

func TestGetRequestParams(t *testing.T) {
	t.Log("should respond with the url param value")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/param/:paramCheck", paramCheckHandler)

	req, _ := http.NewRequest(http.MethodGet, "/param/returnValue", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "\"returnValue\"" //unmarhaled JSON string
	got := res.Body.String()

	assertString(t, want, got)
}

func TestRequestBody(t *testing.T) {
	t.Log("should respond no request body")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodPost, "/requestBody", requestBodyHandler)

	body := struct {
		Response  string `json:"Response"`
		Response2 string `json:"Response2"`
	}{
		Response:  "ok",
		Response2: "ok2",
	}

	jsonData, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		"/requestBody",
		strings.NewReader(string(jsonData)),
	)

	testApp.ServeHTTP(res, req)

	want := "{\"Response\":\"ok\",\"Response2\":\"ok2\"}" //unmarhaled JSON string
	got := res.Body.String()

	assertString(t, want, got)
}

/*
	The next two tests have a fairly strange assert. When the [Request] is
	invalid, we return an error up the call chain. The following tests check the
	scenarios where:

	* 	An invalid field is sent in the request.
	* 	The request is missing a required field. For this to work the struct
		being unmarshaled needs to have a [validate:"required"`] tag on the
		field that is required.

	If either of these scenarios fail, we return a "Trusted" error up the call-
	chain. We leave this error handling to the application author to handle. The
	recommended approach is to handle this in a Middleware. If this is not
	handled, the error will propogate all the way up to the top of the
	framework, where a SHUTDOWN signal will be returned to the [shutdown]
	channel.
*/
func TestRequestInvalidBody(t *testing.T) {
	t.Log("should respond with error")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodPost, "/requestBody", requestBodyHandler)

	body := struct {
		InvalidBody string `json:"InvalidResponse"`
		Response    string `json:"Response"`
	}{
		InvalidBody: "NotOk",
		Response:    "ok",
	}

	jsonData, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		"/requestBody",
		strings.NewReader(string(jsonData)),
	)

	testApp.ServeHTTP(res, req)

	want := syscall.SIGTERM
	got := <-shutdownChannel

	if got != want {
		t.Errorf("wanted signal to be %q but got %q", want, got)
	}
}

func TestRequestBodyRequiredFields(t *testing.T) {
	t.Log("should respond with error")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodPost, "/requestBody", requestBodyHandler)

	body := struct {
		Response string `json:"Response"`
	}{
		Response: "ok",
	}

	jsonData, _ := json.Marshal(body)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(
		http.MethodPost,
		"/requestBody",
		strings.NewReader(string(jsonData)),
	)

	testApp.ServeHTTP(res, req)

	want := syscall.SIGTERM
	got := <-shutdownChannel

	if got != want {
		t.Errorf("wanted signal to be %q but got %q", want, got)
	}
}

// Framework Internals
// =============================================================================

// If we try to pull a value from [context] that doesn't exist we have a big
// problem, always ensure that your [context] value exists, if it doesn't, we
// need to shutdown the app.
func TestContextValueExist(t *testing.T) {
	t.Log("request should contain [Values] on context")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/context", contextHandler)

	req, _ := http.NewRequest(http.MethodGet, "/context", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := "true"
	got := res.Body.String()

	assertString(t, want, got)
}

// When our [Handler] returns an error instead of a nil, we need to shutdown our
// application, this event can only occur if we've lost integrity in our
// application, this is a MAJOR issue. The most likely cause is if we try to
// retrieve a value from [context] that doesn't exist.
func TestShutdownSignal(t *testing.T) {
	t.Log("should return shutdown signal if handler returns an error")
	var shutdownChannel = make(chan os.Signal, 1)

	testApp := web.NewApp(shutdownChannel)
	testApp.Handle(http.MethodGet, "/error", errorHandler)

	req, _ := http.NewRequest(http.MethodGet, "/error", nil)
	res := httptest.NewRecorder()

	testApp.ServeHTTP(res, req)

	want := syscall.SIGTERM
	got := <-shutdownChannel

	if got != want {
		t.Errorf("wanted signal to be %q but got %q", want, got)
	}
}
