package sse

import (
	"bytes"
	"net/http"
	"strconv"
	"time"
)

// Message represents a valid SSE message
type Message struct {
	// Event is a string that identifies the type of event being described.
	// If this value is specified, the event will be sent to the browser listener of the specified event name.
	// Website source code should use addEventListener() to listen for named events.
	// The onmessage handler is called if no event name is specified for the message.
	Event string
	// The data field for the message.
	// When the event source receives several consecutive lines beginning with data,
	// it merges them by inserting a newline character between them.
	// Subsequent new lines areremoved.
	Data string
	// ID to set the ID value of the last event of the EventSource object.
	ID string
	// Retry is the time to reconnect.
	// If the connection to the server is lost,
	// the browser will wait for the specified time before trying to reconnect.
	// It should be an integer number specifying the reconnection time in milliseconds.
	// If a non-integer value is specified, the field is ignored.
	Retry time.Duration
}

func (m *Message) Bytes() []byte {
	// The event stream is a simple stream of text data that must be encoded using UTF-8.
	// Messages in the event stream are separated by a pair of newline characters.
	// The colon as the first character of the line is essentially a comment and is ignored.
	buff := bytes.NewBufferString("")
	if m.Event != "" {
		buff.WriteString("event:" + m.Event + "\n")
	}
	if m.ID != "" {
		buff.WriteString("id:" + m.ID + "\n")
	}
	if m.Data != "" {
		buff.WriteString("data:" + m.Data + "\n")
	}
	if m.Retry != 0 {
		buff.WriteString("retry:" + strconv.Itoa(int(m.Retry.Milliseconds())) + "\n")
	}
	buff.WriteString("\n")
	return buff.Bytes()
}

func DefaultUnsupportedMessageHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Streaming not supported"))
	return err
}
