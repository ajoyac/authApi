// main_test.go

package main_test

import (
	"."
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)
}
func TestLogin(t *testing.T) {
	payload := []byte(`{"username":"user1","password":"password"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var res map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &res)
	if val, ok := res["token"]; !ok {
		t.Errorf("Expected \"{}\". Got \"%s\"", val)
	}
}
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
