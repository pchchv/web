package web

import (
	"fmt"
	"net/http"
)

const (
	// HeaderContentType is a key to refer to the content type of the response header
	HeaderContentType = "Content-Type"
	// ErrInternalServer to send when an internal server error occurs
	ErrInternalServer = "Internal server error"
)

// ErrorData used to render the error page
type ErrorData struct {
	ErrCode        int
	ErrDescription string
}

// dOutput is the standard/valid output wrapped in `{data: <payload>, status: <http response status>}`
type dOutput struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

// errOutput is the error output wrapped in `{errors:<errors>, status: <http response status>}`
type errOutput struct {
	Errors interface{} `json:"errors"`
	Status int         `json:"status"`
}

func crwAsserter(w http.ResponseWriter, rCode int) http.ResponseWriter {
	if crw, ok := w.(*customResponseWriter); ok {
		crw.statusCode = rCode
		return crw
	}

	return newCRW(w, rCode)
}

// Send sends a completely custom response without wrapping it in `{data: <data>, status: <int>` struct
func Send(w http.ResponseWriter, contentType string, data interface{}, rCode int) {
	w = crwAsserter(w, rCode)
	w.Header().Set(HeaderContentType, contentType)
	_, err := fmt.Fprint(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(ErrInternalServer))
		LOGHANDLER.Error(err)
	}
}
