package main

// tinygo build --target=wasi  .
// wash claims sign --name joke joke.wasm -g --http_client

import (
	"errors"

	"github.com/wasmcloud/actor-tinygo"
	httpclient "github.com/wasmcloud/interfaces/httpclient/tinygo"
	"github.com/wasmcloud/interfaces/messaging/tinygo"
)

func main() {
	me := Joke{}
	actor.RegisterHandlers(messaging.MessageSubscriberHandler(&me))
}

type Joke struct{}

func (e *Joke) HandleMessage(ctx *actor.Context, msg messaging.SubMessage) error {
	switch msg.Subject {
	case "new.joke":
		sender := httpclient.NewProviderHttpClient()
		resp, err := sender.Request(ctx, httpclient.HttpRequest{
			Method: "GET",
			Url:    "https://icanhazdadjoke.com",
			Headers: httpclient.HeaderMap{
				"Accept": httpclient.HeaderValues{"application/json"},
			},
			Body: []byte("bug"),
		})
		if err != nil {
			return err
		}

		msgSender := messaging.NewProviderMessaging()
		err = msgSender.Publish(ctx, messaging.PubMessage{
			Subject: msg.ReplyTo,
			ReplyTo: "",
			Body:    resp.Body,
		})
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid topic")
	}
}
