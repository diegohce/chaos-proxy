package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/diegohce/logger"
)


func TestLoadConfig(t *testing.T) {

	log = logger.New("testing-chaos-proxy::")

	if err := loadConfig();  err != nil {
		t.Error(err)
	}

	os.Rename("./chaos-proxy.json", "./chaos-proxy.json.testing")
	defer os.Rename("./chaos-proxy.json.testing", "./chaos-proxy.json")

	if err := loadConfig();  err == nil {
		t.Error("Config file does not exist. No error. Want: error")
	}

}

func TestCreateProxies(t *testing.T) {

	log = logger.New("testing-chaos-proxy::")

	if err := loadConfig();  err != nil {
		t.Error(err)
	}

	 _ = createProxies()

	chaosConfig.Paths["/faulty/host"] = hostConfig{Host:"!#$%&/()="}

	chaosConfig.DefaultHost.Host = "!#$%&/()="

	_ = createProxies()
}

func TestRandomErrorsFunctions(t *testing.T) {

	log = logger.New("testing-chaos-proxy::")

	chaosConfig.MaxTimeout = 0

	errs := []randomError{
		&noError{},
		&hangup{},
		&random5xx{},
		&delay{},
	}

	for _, rerr := range errs {

		switch e := rerr.(type) {
		case *noError:
			t.Log(e.kind())

		case *random5xx:
			t.Log(e.Error())
			statusCode := e.status()
			if statusCode < 500 || statusCode > 511 {
				t.Errorf("random5xx.status() == %d. Want 500 >= status <= 511\n", statusCode)
			}
		case *hangup:
			t.Log(e.kind())

		case *delay:
			t.Log(e.kind())
			e.wait("/chaos/testing")

		default:
			t.Errorf("Unknown random error type %T", e)
		}

	}

	randomErrors = errs

	_ = rollDices()

}

func TestErrorhandler(t *testing.T) {

	log = logger.New("testing-chaos-proxy::")

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/some/path", nil)

	errorHandler(res, req, fmt.Errorf("TIMEOUT"))
	if res.Result().StatusCode != 504 {
		t.Errorf("Got status code %d. Want 504\n", res.Result().StatusCode)
	}

	res = httptest.NewRecorder()
	errorHandler(res, req, fmt.Errorf("5xx"))
	if res.Result().StatusCode != 502 {
		t.Errorf("Got status code %d. Want 502\n", res.Result().StatusCode)
	}

	res = httptest.NewRecorder()
	errorHandler(res, req, random5xx{})
	if res.Result().StatusCode < 500 || res.Result().StatusCode > 511 {
		t.Errorf("random5xx.status() == %d. Want 500 >= status <= 511\n", res.Result().StatusCode)
	}

	res = httptest.NewRecorder()
	errorHandler(res, req, fmt.Errorf("Unspecified error"))
	if res.Result().StatusCode != 502 {
		t.Errorf("Got %d. Want 502", res.Result().StatusCode)
	}

	res = httptest.NewRecorder()
	errorHandler(res, req, fmt.Errorf("HUP"))
	if res.Result().StatusCode != 400 {
		t.Errorf("Got %d. Want 400", res.Result().StatusCode)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorHandler(w, r, fmt.Errorf("HUP"))
	}))
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err == nil {
		t.Error("Error is nil. Want \"EOF\"")
	}
}


func TestModifyResponse(t *testing.T) {

	log = logger.New("testing-chaos-proxy::")

	chaosConfig.MaxTimeout = 0

	req := httptest.NewRequest("GET", "/some/path", nil)
	res := &http.Response{}
	res.Request = req

	randomErrors = []randomError{&noError{}}

	if err := modifyResponse(res); err != nil {
		t.Errorf("Got error %v. Want nil\n", err)
	}

	randomErrors = []randomError{&hangup{}}

	if err := modifyResponse(res); err == nil {
		t.Errorf("Got nil error. Want HUP")
	} else {
		if err.Error() != "HUP" {
			t.Errorf("Got error %v. Want HUP\n", err)
		}
	}

	randomErrors = []randomError{random5xx{}}

	if err := modifyResponse(res); err == nil {
		t.Errorf("Got nil error. Want 5xx")
	} else {
		if err.Error() != "5xx" {
			t.Errorf("Got error %v. Want 5xx\n", err)
		}
	}

	randomErrors = []randomError{delay{}}

	if err := modifyResponse(res); err == nil {
		t.Errorf("Got nil error. Want TIMEOUT")
	} else {
		if err.Error() != "TIMEOUT" {
			t.Errorf("Got error %v. Want TIMEOUT\n", err)
		}
	}
}

