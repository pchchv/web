package web

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSend(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	payload := map[string]string{"hello": "world"}
	reqBody, _ := json.Marshal(payload)

	Send(w, JSONContentType, string(reqBody), http.StatusOK)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err.Error())
		return
	}

	resp := map[string]string{}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(payload, resp) {
		t.Errorf(
			"Expected '%v', got '%v'. Raw response: '%s'",
			payload,
			resp,
			string(body),
		)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf(
			"Expected response status code %d, got %d. Raw response: '%s'",
			http.StatusOK,
			w.Result().StatusCode,
			string(body),
		)
	}
}

func TestSendHeader(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	SendHeader(w, http.StatusNoContent)
	if w.Result().StatusCode != http.StatusNoContent {
		t.Errorf("Expected code '%d', got '%d'", http.StatusNoContent, w.Result().StatusCode)
	}
}

func TestSendError(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	payload := map[string]string{"message": "hello world"}
	SendError(w, payload, http.StatusBadRequest)

	resp := struct {
		Errors map[string]string
	}{}

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err.Error())
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !reflect.DeepEqual(payload, resp.Errors) {
		t.Errorf(
			"Expected '%v', got '%v'. Raw response: '%s'",
			payload,
			resp.Errors,
			string(body),
		)
	}
	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf(
			"Expected response status code %d, got %d. Raw response: '%s'",
			http.StatusBadRequest,
			w.Result().StatusCode,
			string(body),
		)
	}

	// testing invalid response body
	w = httptest.NewRecorder()

	invResp := struct {
		Errors string
	}{}
	invalidPayload := make(chan int)
	SendError(w, invalidPayload, http.StatusBadRequest)
	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = json.Unmarshal(body, &invResp)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if invResp.Errors != `Internal server error` {
		t.Errorf(
			"Expected 'Internal server error', got '%v'. Raw response: '%s'",
			invResp.Errors,
			string(body),
		)
	}

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf(
			"Expected response status code %d, got %d. Raw response: '%s'",
			http.StatusInternalServerError,
			w.Result().StatusCode,
			string(body),
		)
	}

}
