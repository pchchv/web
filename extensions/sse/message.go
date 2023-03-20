package sse

import "time"

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
