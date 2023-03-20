/*
The sse package implements Server-Sent Events (SSE).
This extension is compatible with any net/http implementation, and is not limited to the Web.
*/
package sse

import (
	"context"
	"net/http"
)

type SSE struct {
	// ClientIDHeader is the HTTP request header,
	// which specifies the client identifier.
	// The default is `sse-clientid`.
	ClientIDHeader string
	// UnsupportedMessage is used to send an error response to the client if the server does not support SSE
	UnsupportedMessage func(http.ResponseWriter, *http.Request) error
	// OnCreateClient is a hook for adding a client to the list of active clients.
	// Count is the number of active clients since the last client was added.
	OnCreateClient func(ctx context.Context, client *Client, count int)
	// OnRemoveClient is a hook for removing a client from the list of active clients.
	// Count is the number of active clients after the client was deleted.
	OnRemoveClient func(ctx context.Context, clientID string, count int)
	// OnSend is a hook that is called *after* a message is sent to the client
	OnSend func(ctx context.Context, client *Client, err error)
	// BeforeSend is a hook that is called just before sending a message to the client
	BeforeSend func(ctx context.Context, client *Client)
	Clients    ClientManager
}

func New() *SSE {
	s := &SSE{
		ClientIDHeader:     "sse-clientid",
		Clients:            NewClientManager(),
		UnsupportedMessage: DefaultUnsupportedMessageHandler,
		OnRemoveClient:     DefaultRemoveHook,
		OnCreateClient:     DefaultCreateHook,
		OnSend:             DefaultOnSend,
		BeforeSend:         DefaultBeforeSend,
	}
	return s
}

func DefaultCreateHook(ctx context.Context, client *Client, count int)  {}
func DefaultRemoveHook(ctx context.Context, clientID string, count int) {}
func DefaultOnSend(ctx context.Context, client *Client, err error)      {}
func DefaultBeforeSend(ctx context.Context, client *Client)             {}

func (sse *SSE) Client(id string) *Client {
	return sse.Clients.Client(id)
}

func (sse *SSE) RemoveClient(ctx context.Context, clientID string) {
	cli := sse.Clients.Client(clientID)
	if cli != nil {
		close(cli.Msg)
	}

	sse.OnRemoveClient(
		ctx,
		clientID,
		sse.Clients.Remove(clientID),
	)
}
