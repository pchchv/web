# Server-Sent Events

This extension provides support for [Server-Sent](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) Events for any net/http compatible http server.
It provides the following hooks for configuring workflows:

1. `OnCreateClient func(ctx context.Context, client *Client, count int)`
2. `OnRemoveClient func(ctx context.Context, clientID string, count int)`
3. `OnSend func(ctx context.Context, client *Client, err error)`
4. `BeforeSend func(ctx context.Context, client *Client)`

```golang
import (
    "github.com/pchchv/web/extensions/sse"
)

func main() {
    sseService := sse.New()
    // broadcast to all active clients
    sseService.Broadcast(Message{
        Data:  "Hello world",
        Retry: time.MilliSecond,
	})

	// You can replace ClientManager with your own implementation and override the default sseService.Clients = <your custom client manager>.

    // send message to an individual client
    clientID := "cli123"
    cli := sseService.Client(clientID)
    if cli != nil {
        cli.Message <- &Message{Data: fmt.Sprintf("Hello %s",clientID), Retry: time.MilliSecond }
    }
}
```

## Client Manager

The Client Manager is the interface that is needed for SSE to work. Because it is an interface, it is easier to replace if needed. The default is a simple implementation that uses mutex. If you have your own implementation that is faster/better, you can easily replace the default.
