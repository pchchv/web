package sse

import (
	"context"
	"net/http"
	"sync"
)

type Client struct {
	ID             string
	Msg            chan *Message
	ResponseWriter http.ResponseWriter
	Ctx            context.Context
}

type Clients struct {
	clients   map[string]*Client
	locker    sync.Mutex
	MsgBuffer int
}

type ClientManager interface {
	// New should return the new client as well as the total number of active clients after adding this new client
	New(ctx context.Context, w http.ResponseWriter, clientID string) (*Client, int)
	// The range should iteratively loop through all active clients
	Range(func(*Client))
	// Remove should remove the active client, to which the clientID is assigned, and close the connection
	Remove(clientID string) int
	// Active returns the number of active clients
	Active() int
	// Clients returns a list of all active clients
	Clients() []*Client
	// Client returns *Client if clientID is active
	Client(clientID string) *Client
}

func NewClientManager() ClientManager {
	return &Clients{
		clients: make(map[string]*Client),
		locker:  sync.Mutex{},
	}
}

func (cs *Clients) New(ctx context.Context, w http.ResponseWriter, clientID string) (*Client, int) {
	mchan := make(chan *Message, cs.MsgBuffer)
	cli := &Client{
		ID:             clientID,
		Msg:            mchan,
		ResponseWriter: w,
		Ctx:            ctx,
	}

	cs.locker.Lock()
	cs.clients[clientID] = cli
	count := len(cs.clients)
	cs.locker.Unlock()

	return cli, count
}

func (cs *Clients) Range(f func(cli *Client)) {
	cs.locker.Lock()
	for clientID := range cs.clients {
		f(cs.clients[clientID])
	}
	cs.locker.Unlock()
}

func (cs *Clients) Remove(clientID string) int {
	cs.locker.Lock()
	delete(cs.clients, clientID)
	count := len(cs.clients)
	cs.locker.Unlock()
	return count
}

func (cs *Clients) Active() int {
	cs.locker.Lock()
	count := len(cs.clients)
	cs.locker.Unlock()
	return count
}

// MessageChannels returns a fragment of message channels of all clients,
// which you can use to send messages simultaneously
func (cs *Clients) Clients() []*Client {
	idx := 0
	cs.locker.Lock()
	list := make([]*Client, len(cs.clients))
	for clientID := range cs.clients {
		cli := cs.clients[clientID]
		list[idx] = cli
		idx++
	}
	cs.locker.Unlock()
	return list
}

func (cs *Clients) Client(clientID string) *Client {
	cs.locker.Lock()
	cli := cs.clients[clientID]
	cs.locker.Unlock()

	return cli
}
