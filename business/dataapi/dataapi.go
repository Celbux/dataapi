package dataapi

import (
	"encoding/json"
	"fmt"
	"github.com/Celbux/dataapi/foundation/tools"
	"github.com/Celbux/dataapi/foundation/web"
	"github.com/japm/goScript"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// AssertContains will return an error if parameter 1 is not contained in parameter 0
// Usage: [AssertContains(0, 1)]
// Eg: [AssertContains("foobar", "bar")]
// Parameter 0: the string containing the element
// Eg: "foobar"
// Parameter 1: the element that is a subset of the string
// Eg: "bar"
// This will not throw an error as parameter 1 is contained in parameter 0
func (d DataAPIService) AssertContains(params string) error {

	// Get parameters 0 and 1
	parameters := d.GetParameters(params)
	if len(parameters) != 2 {
		return errors.Errorf("AssertContains expected 2 parameter but got: %v", len(parameters))
	}
	val1, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	val2, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}

	// Perform the check
	if !strings.Contains(val1, val2) {
		return errors.Errorf("parameter 0: [%v] does not contain parameter 1: [%v]", val1, val2)
	}

	// Return success
	return nil

}

// AssertEquals will not throw an error if the two input parameters are identical
// Usage: [AssertEquals(0, 1)]
// Eg: [AssertContains("foobar", "foobar")]
// Parameter 0: value 1
// Eg: "foobar"
// Parameter 1: value 2
// Eg: "bar"
// This will throw an error as parameter 0 is not equal to parameter 1
func (d DataAPIService) AssertEquals(params string) error {

	// Get parameters 0 and 1
	parameters := d.GetParameters(params)
	if len(parameters) != 2 {
		return errors.Errorf("AssertEquals expected 2 parameter but got: %v", len(parameters))
	}
	val1, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	val2, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}

	// Perform the check
	if val1 != val2 {
		return errors.Errorf("Expected %v but got %v", val2, val1)
	}

	// Return success
	return nil

}

// AssertFailure will ensure that the last core failed with the given error code Eg: "-22" (insufficient funds)
// Usage: [AssertFailure(0)]
// Eg: [AssertFailure("-22")]
// Parameter 0: the error code that core returned
// Eg: "-22"
// This will throw an error if the error code from the last Core call is not -22
func (d DataAPIService) AssertFailure(params string) error {

	// Get parameter 0
	parameters := d.GetParameters(params)
	if len(parameters) != 1 {
		return errors.Errorf("AssertEquals expected 1 parameter but got: %v", len(parameters))
	}
	error, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}

	// Get the res field off the EvalCache to perform the error code check
	res := d.EvalCache["res"].([]string)
	if res == nil || len(res) == 0 {
		return errors.New("res[] struct returned from core is empty")
	}
	if res[0] != error {
		return errors.Errorf("expected to fail with error: [%v] but got res: [%v]", error, res)
	}

	// Return success
	return nil

}

// AssertStringArrEquals will not error if the two input string arrays are equal
// Given: [Set(a, "1,2,3", []string)] and [Set(b, "3,2,1", []string)]
// Usage: [AssertStringArrEquals(0, 1, 2)]
// Eg: [AssertStringArrEquals(a, b, true)]
// Parameter 0: the first array to compare
// Eg: a
// Parameter 1: the second array to compare
// Eg: b
// Parameter 2: a boolean to determine whether order matters or not
// Eg: true
// This will not throw an error as order doesn't matter and the arrays are therefore equal
func (d DataAPIService) AssertStringArrEquals(params string) error {

	// Get parameters 0, 1 and 2
	parameters := d.GetParameters(params)
	if len(parameters) != 3 {
		return errors.Errorf("AssertStringArrEquals expected 3 parameters but got: %v", len(parameters))
	}
	aRaw, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	a := strings.Split(aRaw, "___")
	bRaw, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}
	b := strings.Split(bRaw, "___")
	orderMatters, err := d.EvalBool(parameters[2])
	if err != nil {
		return err
	}

	// Perform the check
	if len(a) != len(b) {
		return errors.New("array sizers differ")
	}
	// If the order matters, then both arrays should line up
	if orderMatters {
		for i, _ := range a {
			if a[i] != b[i] {
				return errors.Errorf("[%v] does not equal [%v]", a[i], b[i])
			}
		}
	} else {
		// Else order does not matter, therefore, make sure both arrays a and be are seld containing
		aMap := make(map[string]int)
		for _, val := range a {
			aMap[val] = 0
		}
		for _, val := range b {
			_, ok := aMap[val]
			if !ok {
				return errors.Errorf("[%v] was not found in array a", val)
			}
			aMap[val] = 1
		}
		for key, val := range aMap {
			if val == 0 {
				return errors.Errorf("[%v] was not found in array b", key)
			}
		}
	}

	// Return success
	return nil

}

// AssertSuccess will ensure that the last core call did not fail
// Usage: [AssertSuccess()]
func (d DataAPIService) AssertSuccess() error {

	// Get the response of the EvalCache and ensure no error code is present
	res, ok := d.EvalCache["res"].([]string)
	_, ok2 := d.EvalCache["res"].(map[string]interface{})
	d.EvalCache["res"] = nil
	if !ok && !ok2 {
		return errors.New("could not convert EvalCache[\"res\"] to a supported type")
	}
	if res == nil || len(res) == 0 {
		return errors.New("res[] struct returned from core is empty")
	}
	if strings.HasPrefix(res[0], "-") {
		return errors.Errorf("expected success but failed with core res[]: %v", res)
	}

	// Return success
	return nil

}

// DirectoryExists returns true if the directory exists
func (d DataAPIService) DirectoryExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

// Eval will execute the data code expression given
// Eg data code: [Set(s, "Hello World!", string)][PrintF("%v", s)]
// This string input will evaluate to printing "Hello World!" to the console
func (d DataAPIService) Eval(expression string) (map[string]interface{}, error) {

	// Create return map
	out := make(map[string]interface{})

	// Get all expressions
	// There could be more than 1 data block given and thus
	// would have to be evaluated individually in a for loop
	var allErrors []error
	expressionsRaw := d.GetExpressions(expression)
	for i, rawExpression := range expressionsRaw {
		var method string
		var params string
		expression = rawExpression
		if strings.Contains(rawExpression, "[") && strings.Contains(rawExpression, "]") {
			// Get method names and prepare parameters on the EvalCache
			// if the expression contains a Functions call
			var err error
			method, params, err = d.GetFunctionDefinition(rawExpression)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
			d.EvalCache["params"] = params
			expression = method + "(params)"
			if len(params) == 0 {
				expression = method + "()"
			}
		}

		// Prepare expression struct and execute the expression using the EvalCache
		// The EvalCache is the only available variables in the scope of the evaluation function
		exp := &goScript.Expr{}
		err := exp.Prepare(expression)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		// Set newMap, as exp.Eval will break if the map is not a new pointer
		newMap := make(map[string]interface{})
		newMap = d.EvalCache

		// Eval the expression
		// All functions that ran to completion and did not throw an error will pass
		// The report tree will prune out the success and failures from the tree last
		val, err := exp.Eval(newMap)
		if err != nil {
			method, _, _ := d.GetFunctionDefinition(rawExpression)
			methodNotFound := tools.StripField(err.Error(), "Method (.*) not found")
			if methodNotFound == strings.Replace(method, "dataapi.", "", 1) {
				expression = strings.Replace(expression, "dataapi", "dataapi.CoreDataAPI", 1)
				exp := &goScript.Expr{}
				err := exp.Prepare(expression)
				if err != nil {
					allErrors = append(allErrors, err)
					continue
				}
				val, err = exp.Eval(newMap)
				if err != nil {
					allErrors = append(allErrors, err)
					continue
				}
			} else {
				allErrors = append(allErrors, err)
				continue
			}
		}
		if val == nil && i == len(expressionsRaw) - 1 {
			return nil, errors.New("[Pass()]")
		}

		// Handle the returned data that was returned from Eval
		// This could be a primitive or a report tree
		switch val.(type) {
		case map[string]interface{}:
			// Check for cascading results, and return them to Evaluate from Eval
			// All failures/successes will be evaluated into a report tree
			// We don't want to fail the evaluate process immediately
			res := val.(map[string]interface{})
			_, ok := res["report"].(*tools.Tree)
			if ok {
				return res, nil
			}
		case error:
			// If there is an error append it and continue
			allErrors = append(allErrors, val.(error))
			continue
		default:
			// Used for eval of single return values
			out["val"] = fmt.Sprintf("%v", val)
		}

	}

	if len(allErrors) > 0 {
		errString := ""
		for _, err := range allErrors {
			if err.Error() == "[Pass()]" {
				continue
			}
			errString += fmt.Sprintf("%v, ", err.Error())
		}
		if errString == "" {
			errString = "[Pass()]"
		}
		return nil, errors.New(errString)
	}

	// Return success
	return out, nil

}

// Evaluate is a data function that runs all test data code inside a dataApi tests file
// It takes in the test file name to run as input
// Each line in the input file gets evaluated and all errors arr built up into a tree
// The tree is used to store which test failed and which tests passed
// and is return as the result of a Data API run
// You can call Evaluate on a directory or a file
// Usage: [Evaluate(0)]
// Eg: [Evaluate("cascadingerrors")]
// Parameter 0: The directory or test case file you want to run
// Eg: "cascadingerrors"
func (d DataAPIService) Evaluate(inFile string) map[string]interface{} {

	// Log file to track which test is currently running
	d.Log.Println(inFile)

	// Evaluate can not fail and always returns a report
	// Any error will be associated with the file name that is being Evaluated
	// Eg: tree["someTest.txt"] = "error: some function failed"
	out := make(map[string]interface{})
	report := &tools.Tree{Data: inFile}
	out["report"] = report

	// The file name will be evaluated so ensure it is in quotations
	if inFile[0] != '"' && inFile[len(inFile)-1] != '"' {
		inFile = "\"" + inFile + "\""
	}
	inFile, err := d.EvalString(inFile)
	if err != nil && err.Error() != "[Pass()]" {
		report.Add(inFile, err.Error())
		return out
	}
	filepath := "configs/dataapi/" + inFile
	report.Data = filepath

	// Read the file contents line by line into s.Functions.Eval
	osFile, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		report.Add(filepath, "could not open file")
		report.Add("could not open file", err.Error())
		return out
	}
	var dataRawArr [][]byte
	dataRaw, err := ioutil.ReadAll(osFile)
	if err == nil && dataRaw != nil {
		dataRawArr = append(dataRawArr, dataRaw)
	}
	if err != nil {
		// If the file is a directory we will want to run Evaluate on all the files within
		// Therefore, get all the data first and loop through it later
		if d.DirectoryExists(filepath) {
			files, err := ioutil.ReadDir(filepath)
			if err != nil {
				report.Add(filepath, err.Error())
				return out
			}

			for _, file := range files {
				filepathNested := inFile + "/" + file.Name()
				reportRaw := d.Evaluate(filepathNested)
				err, ok := reportRaw["err"].(error)
				if ok {
					report.Add(filepath, filepathNested)
					report.Add(filepathNested, err.Error())
					continue
				}
				if reportRaw != nil {
					switch reportRaw["report"].(type) {
					case *tools.Tree:
						reportReturned, ok := reportRaw["report"].(*tools.Tree)
						if !ok {
							report.Add(filepath, filepathNested)
							report.Add(filepathNested, "report returned but is not a tree")
							continue
						}
						report.Add(filepath, reportReturned.Data)
						report.AddTree(*reportReturned)
						continue
					}
				}
			}
			return out
		}
	}
	// If the input file contains no data, throw an error
	if len(dataRawArr) == 0 {
		report.Add(filepath, "there is no data in the input file to evaluate")
		return out
	}

	// Loop over every line in the input file
	// Add all calls and the data they returned to the report
	// This will be used to create the failures and success report lastly
	for _, rawData := range dataRawArr {
		data := strings.Split(string(rawData), "\n")
		for _, datum := range data {
			if len(datum) == 0 || datum[0] == '#' || strings.TrimSpace(datum) == "" {
				continue
			}

			reportRaw, err := d.Eval(datum)
			if err != nil {
				report.Add(filepath, datum)
				report.Add(datum, err.Error())
				continue
			}
			if reportRaw != nil {
				switch reportRaw["report"].(type) {
				case *tools.Tree:
					reportReturned, ok := reportRaw["report"].(*tools.Tree)
					if !ok {
						report.Add(filepath, datum)
						report.Add(datum, "report returned but is not a tree")
						continue
					}
					report.Add(filepath, reportReturned.Data)
					report.AddTree(*reportReturned)
					continue
				}
			}
		}
	}

	// Return success with the report
	return out

}

// EvalBool returns the boolean value the expression returned
func (d DataAPIService) EvalBool(expression string) (bool, error) {

	// Evaluate the expression
	res, err := d.Eval(expression)
	if err != nil {
		return false, err
	}

	// Ensure a value was returned
	booleanRaw, ok := res["val"]
	if !ok {
		return false, errors.New("expression did not evaluate to a boolean")
	}

	// Ensure the returned value is in fact a boolean
	boolean, err := strconv.ParseBool(booleanRaw.(string))
	if err != nil || !ok {
		return false, err
	}

	return boolean, nil

}

// EvalInt returns the int value the expression returned
func (d DataAPIService) EvalInt(expression string) (int, error) {

	// Evaluate the expression
	res, err := d.Eval(expression)
	if err != nil {
		return 0, err
	}

	// Ensure a value was returned
	intRaw, ok := res["val"]
	if !ok {
		return 0, errors.New("expression did not evaluate to an int")
	}

	// Ensure the returned value is in fact an int
	intVal, err := strconv.Atoi(intRaw.(string))
	if err != nil {
		return 0, err
	}

	return intVal, nil

}

// EvalString returns the string value the expression returned
func (d DataAPIService) EvalString(expression string) (string, error) {

	// Evaluate the expression
	res, err := d.Eval(expression)
	if err != nil {
		return "", err
	}

	// Ensure a value was returned
	stringRaw, ok := res["val"]
	if !ok {
		return "", errors.New("expression did not evaluate to a string")
	}

	return stringRaw.(string), nil

}

// Fail will return the given string as an error
func (d DataAPIService) Fail(err string) error {
	return errors.New(err)
}

// For is just like your normal for loop:
// Usage: [For(0, 1)]
// Eg: [For((i < 3), [PrintF("Hello World %v", i)][Set(i,i+1,int)])]
// Parameter 0: is the condition that will break the loop
// Eg: (i < size)
// Parameter 1: is the code that runs inside the for loop
// Eg: [PrintF("Hello World %v", i)][Set(i,i+1,int)]
// This for loop will output:
// > Hello World 0
// > Hello World 1
// > Hello World 2
func (d DataAPIService) For(params string) interface{} {

	// Gets parameters 0 and 1
	expressions := d.GetParameters(params)
	if len(expressions) != 2 {
		return errors.Errorf("for loop expected 2 expressions but got: %v", len(expressions))
	}

	for {
		// Run the first expression to see if it returns true or false
		boolean, err := d.EvalBool(expressions[0])
		if err != nil {
			return err
		}

		// If false, break the loop
		if !boolean {
			break
		}

		// If true, run the expression in parameter 2
		_, err = d.Eval(expressions[1])
		if err != nil && err.Error() != "[Pass()]"{
			return err
		}
	}

	// Return success
	return errors.New("[Pass()]")

}

// getExpressions retrieves all expressions nested on the same scope level
// This is equivalent to splitting each line of code like java does with semi colons
// Eg: 	expressions([Set(i,1,int)][Set(i,[GetSomeVar()],i,int)][Set(j,3,string)])
// Will split into: [Set(i,1,int)]	[Set(i,[GetSomeVar()],i,int)]	[Set(j,3,string)]
// And return: []string = {
//		"[Set(i,1,int)]",
//		"[Set(i,[GetSomeVar()],i,int)]",
//		"[Set(j,3,string)]",
// }
func (d DataAPIService) GetExpressions(expressions string) []string {

	var out []string
	// If the expression does not contain '[' or ']'
	// then we can directly evaluate it as it is not a Functions call
	// Eg: (i <= 10) versus [Set(Val, i <= 10, bool)]
	// The first instance will set the 'res' variable by default on the EvalCache
	// The second instance will call Set() and  set the 'Val' variable on the EvalCache
	// Both instances set the value to the evaluated result of i <= 10
	if !strings.Contains(expressions, "[") && !strings.Contains(expressions, "]") {
		return append(out, expressions)
	}

	// Get all the expressions and put them into []string
	scope := 0
	start := strings.Index(expressions, "[")
	for i, char := range expressions {
		letter := string(char)
		if letter == "[" {
			scope++
		} else if letter == "]" {
			scope--
			if scope == 0 {
				out = append(out, expressions[start:i+1])
				start = i + 1
			}
		}
	}

	// Return the expressions
	return out

}

// GetFunctionDefinition will separate the function name and parameter list
// Input:
// > expression = [PrintF("Hello World %v", i)]
// Output:
// > method = dataapi.PrintF
// > params = ("Hello World %v", i)
// > err = nil
func (d DataAPIService) GetFunctionDefinition(expression string) (string, string, error) {

	// Get indexes of '[', ']' and '(', ')'
	// to slice the method name and parameter list respectively
	scope := 0
	start := 0
	end := 0
	for i, letter := range expression {
		if letter == '[' {
			scope++
			if scope == 1 {
				start = i
			}
		} else if letter == ']' {
			scope--
		}
		if scope == 0 {
			end = i
			break
		}
	}
	if end == 0 || start >= end {
		return "", "", errors.New("unable to evaluate expression")
	}
	exp := expression[start+1 : end]
	index := strings.Index(exp, "(")
	if index == -1 {
		return "", "", errors.New("can not find '(' in the expression")
	}
	method := "dataapi." + exp[:index]
	params := exp[index:]
	params = strings.TrimLeft(params, "(")
	params = strings.TrimRight(params, ")")

	return method, params, nil

}

// getParameters retrieves all the parameters seperated by a comma
// Eg: [For([Set(i,0,int)],[Bool(i < 10)],[Set(i,0,int)],[doWork()])]
// The expression must return all parameters delimited by the top most comma delimitation
// Therefore: [For([Set(i,0,int)]	,	[Bool(i < 10)]	,	[Set(i,0,int)]	,	[doWork()])]
// Top most commas:					^					^					^
// Return: []string = {
//		"[Set(i,0,int)]",
//		"[Bool(i < 10)]",
//		"[Set(i,0,int)]",
//		"[doWork()]"
// }
func (d DataAPIService) GetParameters(expression string) []string {

	// Get all parameters and put them in []string
	var out []string
	scope := 0
	start := 0
	for i, char := range expression {
		letter := string(char)
		if letter == "[" {
			scope++
		} else if letter == "]" {
			scope--
		} else if letter == "," && scope == 0 {
			out = append(out, expression[start:i])
			start = i + 1
		}
	}
	out = append(out, expression[start:])

	// Return all the parameters
	return out

}

// GetResults will mine the report for successes and failures after calling Evaluate on a file
// The results will pretty print to the user
// Evaluate has a batch error mechanism for handling each line evaluated
// GetResults separates the successes and failures from the report tree
func (d DataAPIService) GetResults(report tools.Tree) ([]string, []string, error) {

	// Mine all the failed results from the report
	failures, err := report.GetFailures()
	if err != nil {
		return nil, nil, err
	}

	// Mine all the success results from the report
	successes, err := report.GetSuccesses()
	if err != nil {
		return nil, nil, err
	}

	// Remove parent directory names as they are spammy
	for i, failure := range failures {
		failures[i] = strings.Replace(failure, "configs/dataapi/", "", -1)
	}
	for i, success := range successes {
		successes[i] = strings.Replace(success, "configs/dataapi/", "", -1)
	}

	// Return the failures and successes
	return failures, successes, nil

}

// If is just like your normal if statement:
// Usage: If(0, 1)
// Eg: [If((i < 3), [Println("Hello World")])]
// Parameter 0: the condition that will allow the statement to execute
// Eg: (i < 3)
// Parameter 1: is the code that runs inside the if statement
// Eg: [PrintF("Hello World %v", i)][Set(i,i+1,int)]
// Hello World will print if the parameter 0 evaluates to true
func (d DataAPIService) If(params string) interface{} {

	// Gets parameters 0 and 1
	parameters := d.GetParameters(params)
	if len(parameters) != 2 {
		return errors.Errorf("if statement expected 2 parameters but got: %v", len(parameters))
	}
	boolean, err := d.EvalBool(parameters[0])
	if err != nil {
		return err
	}

	// If true then run the expression in parameter 1
	if boolean {
		res, err := d.Eval(parameters[1])
		if err != nil {
			return err
		}
		errVal, ok := res["error"]
		if ok {
			return errVal.(error)
		}

	}

	// Return success
	return nil

}

// ParallelPost is a data function that will send multiple POST requests in parallel
// Usage:
// [Set(files, "file1.txt___file2.txt", string)]
// [Set(headers, "Monkey:::Madness---Content-Type:::application/json___Monkey:::Madness---Content-Type:::application/json"
// [Set(jsons, "{\"Data\": \"1234\"}___{\"Data\": \"5678\"}", string)]
// [Set(urls, "https://someUrl.com/someEndpoint___https://someOtherUrl.com/someOtherEndpoint", string)]
// [ParallelPost(files, jsons, urls)]
//
// Parameter 0: array list delimited by "___" of the multipart files you want to attach to the multipart request:
// Eg: "file1.txt___file2.txt"
// Parameter 1: array list delimited by "___" of multiple key value pairs (mapped with ":::" and seperated by "---") that is the HTTP request headers:
// Eg: "Monkey:::Madness---Content-Type:::application/json___Monkey:::Madness---Content-Type:::application/json" will add the following headers:
// Request1 headers: Monkey: Madness, Content-Type: application/json
// Request2 headers: Monkey: Madness, Content-Type: application/json
// Parameter 2: array list delimited by "___" of the json bodies
// Eg: "{\"Data\": \"1234\"}___{\"Data\": \"5678\"}"
// Parameter 3: array list delimited by "___" of the urls to send the request to
//
// ParallelPost will save the respective responses of the POST requests on the EvalCache as "ParallelPostX"
// Where X is an integer value
// You can access the raw JSON string response from:
// [Set(response1, [Res("ParallelPost1")], string)]
// [Set(response2, [Res("ParallelPost2")], string)]
func (d DataAPIService) ParallelPost(params string) interface{} {

	// Gets parameters 0, 1, 2 and 3
	parameters := d.GetParameters(params)
	if len(parameters) != 4 {
		return errors.Errorf("data function 'ParallelPost' expected 4 parameters but got: %v", len(parameters))
	}
	filesRaw, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	_ = strings.Split(filesRaw, "___")
	var headers []map[string]string
	allHeadersRaw, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}
	allHeadersArr := strings.Split(allHeadersRaw, "___")
	for _, headersArr := range allHeadersArr {
		headersMap := make(map[string]string)
		headersArr := strings.Split(headersArr, "---")
		for _, header := range headersArr {
			header := strings.Split(header, ":::")
			if len(header) != 2 {
				return errors.Errorf("header is not in the format 'key::value'")
			}
			headersMap[header[0]] = header[1]
		}
		headers = append(headers, headersMap)
	}
	var jsonMaps []map[string]interface{}
	jsonsRaw, err := d.EvalString(parameters[2])
	if err != nil {
		return err
	}
	jsons := strings.Split(jsonsRaw, "___")
	for _, jsonBody := range jsons {
		jsonBody = strings.Replace(jsonBody, "\\\"", "\"", -1)
		jsonMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(jsonBody), &jsonMap)
		if err != nil {
			return err
		}
		jsonMaps = append(jsonMaps, jsonMap)
	}
	urlsRaw, err := d.EvalString(parameters[3])
	if err != nil {
		return err
	}
	urls := strings.Split(urlsRaw, "___")
	if len(urls) != len(jsonMaps) {
		return errors.New(fmt.Sprintf("number of URLs (%v) does not align with the number of JSON bodies (%v)", len(urls), len(jsons)))
	}

	// Make the Post requests asynchronously
	var waitGroup sync.WaitGroup
	for i := 0; i < len(jsonMaps); i++ {
		waitGroup.Add(1)
		go func(url string, header map[string]string, jsonBody map[string]interface{}, i int) {
			defer waitGroup.Done()
			resp, err := web.DoRequest(url, header, http.MethodPost, jsonBody)
			if err != nil {
				d.EvalCache[fmt.Sprintf("ParallelPost%v", i)] = errors.Errorf("error sending POST request to URL %v", url)
				return
			}
			d.EvalCache[fmt.Sprintf("ParallelPost%v", i)] = string(resp)
		}(urls[i], headers[i], jsonMaps[i], i)
	}

	// Wait until all the POST requests are complete
	waitGroup.Wait()

	// Return success
	return nil

}

// Pass is a function that can not fail and is used to build the failures/passes tree
// to return a Data API output report
func (d DataAPIService) Pass(params string) {
	return
}

// Post will send a POST request which includes a file and JSON data
// Usage: [Post(0, 1, 2)]
// Eg: 	[Set(url, "https://rnd-api-v1-dot-dev8celbux.uc.r.appspot.com/api/rnd/pay?ns=rnd", string)]
//     	[Set(jsonBody, "{\"VoucherNo\": \"117-22427-719752\",\"StoreID\":\"Store1\",\"Reference\":\"1234\",\"Amount\":\"2000\",\"Currency\":\"{{currency}}\",\"Metadata\":\"\",\"RequestDT\":\"1234\"}", string)]
//		[Set(heders, "Authorization___Bearer 9m1,Monkey___Madness", string)]
//		[Post(url, jsonBody, headers)]
// Parameter 0: the target url
// Eg: "https://rnd-api-v1-dot-dev8celbux.uc.r.appspot.com/api/rnd/pay?ns=rnd"
// Parameter 1: json input
// Eg: "{\"VoucherNo\": \"117-22427-719752\",\"StoreID\":\"Store1\",\"Reference\":\"1234\",\"Amount\":\"2000\",\"Currency\":\"{{currency}}\",\"Metadata\":\"\",\"RequestDT\":\"1234\"}"
// Parameter 2: request headers
// Eg: "Authorization___Bearer 9m1,Monkey___Madness"
// The JSON response will be set under the variable "res" on the EvalCache and be accessed by the Res data function
// Eg: [Set(response1, [Res("res")], string)]
func (d DataAPIService) Post(params string) interface{} {

	// Gets parameters 0, 1 and 2
	parameters := d.GetParameters(params)
	if len(parameters) != 3 {
		return errors.Errorf("data function 'Post' expected 3 parameters but got: %v", len(parameters))
	}
	url, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	jsonStr, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}
	jsonStr = strings.Replace(jsonStr, "\\\"", "\"", -1)
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return err
	}
	headers := make(map[string]string)
	headersRaw, err := d.EvalString(parameters[2])
	if err != nil {
		return err
	}
	headersArr := strings.Split(headersRaw, ",")
	for _, headerRaw := range headersArr {
		if headerRaw == "" {
			continue
		}
		header := strings.Split(headerRaw, "___")
		if len(header) != 2 {
			return errors.Errorf("header is not in the format 'key___value'")
		}
		headers[header[0]] = header[1]
	}

	// Make the Post request
	resp, err := web.DoRequest(url, headers, http.MethodPost, jsonMap)
	response := string(resp)
	if err != nil {
		return err
	}
	d.EvalCache["res"] = response

	// Return success
	return nil

}

// PrintF is just like your fmt.PrintF:
// Usage: [PrintF(0, 1, 2)]
// Eg: [PrintF("Hello World %v%v%v", 1, 2, 3)]
// Parameter 0: is the formatted string
// Eg: "Hello World %v%v%v"
// Parameter 1++: are the variadic variables
// Eg: 1, 2, 3
// This will output:
// > Hello World 123
func (d DataAPIService) PrintF(params string) interface{} {

	// Gets parameters 0, 1 and 2...
	parameters := d.GetParameters(params)
	if len(parameters) == 0 {
		return errors.Errorf("PrintF expected atleast 1 parameter but got: %v", len(parameters))
	}

	// Get format as first param
	format := parameters[0]

	// Build []interface{} so that we can call input our data into the variadic field in fmt.Sprintf
	var aInterfaces []interface{}
	for i, s := range parameters {
		if i == 0 {
			continue
		}
		sOut, err := d.Eval(s)
		if err != nil {
			return err
		}
		aInterfaces = append(aInterfaces, sOut["val"])
	}

	// Print with trimming the wrapping parenthesis
	s := fmt.Sprintf(format, aInterfaces...)
	s = s[1:len(s)-1]
	d.Log.Println(s)

	// Return success
	return nil

}

// ReadFile will read the string data from the given filepath and
// save its contents under the given variable on the EvalCache
// Usage: [ReadFile(0, 1)]
// Eg: [ReadFile("expected", "cascadingerrors/expected_output.txt")]
// Parameter 0: the variable on the EvalCache that the returned value will be stored
// Eg: "expected"
// Parameter 1: the filepath of the file we want to read, the value is relative the root directory
// Eg: "configs/dataapi/cascadingerrors/expected_output.txt"
func (d DataAPIService) ReadFile(params string) interface{} {

	// Gets parameters 0 and 1
	parameters := d.GetParameters(params)
	if len(parameters) == 0 {
		return errors.Errorf("ReadFile expected 2 parameters but got: %v", len(parameters))
	}
	variable, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}
	filepath, err := d.EvalString(parameters[1])
	if err != nil {
		return err
	}

	// Read the contents of the file
	osFile, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	dataRaw, err := ioutil.ReadAll(osFile)
	if err != nil {
		return err
	}
	data := string(dataRaw)

	// Save the contents of the file under the given variable on the EvalCache
	d.EvalCache[variable] = data

	return data
}

// Res will return the given property out of the EvalCache
// Usage: Res(0)
// Eg: [Res("error")]
// Parameter 0: the field you want to retrieve from the EvalCache
// Eg: "error"
// This will try return the "error" field from EvalCache["res"]:
// > Invalid request
func (d DataAPIService) Res(params string) interface{} {

	// Gets parameter 0
	parameters := d.GetParameters(params)
	if len(parameters) != 1 {
		return errors.Errorf("Res expected 1 parameter but got: %v", len(parameters))
	}
	field, err := d.EvalString(parameters[0])
	if err != nil {
		return err
	}

	// Get the field
	out, ok := d.EvalCache[field]
	if !ok {
		return errors.Errorf("field %v could not be found on the EvalCache", field)
	}

	// Return the field
	return out

}

// Sleep will sleep for X seconds
// Usage: [Sleep(0)]
// Eg: [Sleep(5)]
// Parameter 0: the amount of seconds you want to sleep
// Eg: 5
func (d DataAPIService) Sleep(params string) interface{} {

	// Gets parameter 0
	parameters := d.GetParameters(params)
	if len(parameters) != 1 {
		return errors.Errorf("Sleep expected 1 parameter but got: %v", len(parameters))
	}
	seconds, err := d.EvalInt(parameters[0])
	if err != nil {
		return err
	}

	// Sleep ZZZzzz...
	time.Sleep(time.Duration(seconds) * time.Second)

	return nil

}

// Set will create a variable on the EvalCache
// Usage: [Set(0, 1, 2)]
// Eg: [Set(i, 0, int)] or [Set(s, "hello world!", string)]
// Parameter 0: the name of the variable
// Eg: i
// Parameter 1: the value of the variable
// Eg: 0
// Parameter 2: the type of the value
// Eg: int, string, bool
// Therefore, the above is the equivalent to running 'i := 0'
// The variable 'i' will be available when Eval() is run as it is set on the EvalCache
func (d DataAPIService) Set(params string) interface{} {

	// Gets parameters 0, 1, and 2
	parameters := d.GetParameters(params)
	rawString := ""
	if len(parameters) > 3 {
		for i, parameter := range parameters {
			if i == 0 || i == len(parameters)-1 {
				continue
			}

			rawString += parameter + ","
		}
	}
	if rawString == "" {
		rawString = parameters[1]
	}
	parameters = []string{parameters[0], strings.TrimRight(rawString, ","), parameters[len(parameters)-1]}
	if len(parameters) != 3 {
		return errors.Errorf("Set expected 3 parameters but got: %v", len(parameters))
	}

	// Set the EvalCache to the returned value of Parameter 1
	// Ensure the returned data type aligns with Parameter 2
	// This sounds weird, but read it a couple of times with the code
	variable := strings.TrimSpace(parameters[0])
	value := strings.TrimSpace(parameters[1])
	variableType := strings.TrimSpace(parameters[2])
	if variableType == "string" {
		value, err := d.EvalString(value)
		if err != nil {
			return err
		}
		d.EvalCache[variable] = value
	} else if variableType == "int" {
		value, err := d.EvalInt(value)
		if err != nil {
			return err
		}
		d.EvalCache[variable] = value
	} else if variableType == "boolean" {
		value, err := d.EvalBool(value)
		if err != nil {
			return err
		}
		d.EvalCache[variable] = value
	} else {
		return errors.Errorf("variableType \"%v\" does not exist", variableType)
	}

	// Return success
	return nil

}
