package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/dimfeld/httptreemux"
	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator *ut.UniversalTranslator

func init() {

	// Instantiate the validator for use.
	validate = validator.New()

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Params returns the web call parameters from the request.
func Params(r *http.Request) map[string]string {
	return httptreemux.ContextParams(r.Context())
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(
	r *http.Request,
	val interface{},
) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at
		// the Accept-Language header if you intend to support multiple
		// languages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Field(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return &Error{
			Err:    errors.New("field validation error"),
			Status: http.StatusBadRequest,
			Fields: fields,
		}
	}

	return nil
}

// DoRequest handles sending a basic HTTP request to any URL
// and get a response as []byte
func DoRequest(url string, headers map[string]string, httpMethod string, data interface{}) ([]byte, error) {

	// Create the http request
	// Encode the data and set its content type in the case of an http POST
	var req *http.Request
	var err error
	if httpMethod == http.MethodPost {
		req, err = http.NewRequest(httpMethod, url, Encode(data))
	} else if httpMethod == http.MethodGet {
		req, err = http.NewRequest(httpMethod, url, nil)
	} else {
		err = errors.New(fmt.Sprintf("unrecognized httpMethod %v", httpMethod))
	}
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Attempt to do http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// If not StatusOK, return error
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("failed to send %v request to %v", httpMethod, url))
		bytes, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			return nil, fmt.Errorf("%v: %v", err, err2)
		}
		if string(bytes) != "" {
			return nil, fmt.Errorf("%v: %v", err, string(bytes))
		}
		return nil, err
	}

	// Return success
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, nil

}

// GetPathParam gets a specified path param using the treemux Context
func GetPathParam(ctx context.Context, param string) (string, error) {
	// Get params from context as map
	params := httptreemux.ContextParams(ctx)

	// Get specified param from map
	val, ok := params[param]
	if !ok {
		// Path param not found
		return "", errors.New(fmt.Sprintf("/%s/:%s not provided!", param, param))
	}

	// Support the default namespace
	if val == "default" {
		val = ""
	}

	// Return success
	return val, nil
}