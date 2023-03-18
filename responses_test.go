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
