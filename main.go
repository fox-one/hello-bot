package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	bot "github.com/MixinNetwork/bot-api-go-client"
	number "github.com/MixinNetwork/go-number"
	config "github.com/fox-one/hello-bot/config"
	uuid "github.com/satori/go.uuid"
)

// CNBAssetID is the CNB's ID in Mixin Network
const CNBAssetID = "965e5c6e-434c-3fa9-b780-c50f43cd955c"

const defaultResponse = "I got it."

var client *bot.BlazeClient

// Handler is an implementation for interface bot.BlazeListener
// check out the url for more details: https://github.com/MixinNetwork/bot-api-go-client/blob/master/blaze.go#L89.
type Handler struct{}

// OnMessage is a general method of bot.BlazeListener
func (r Handler) OnMessage(ctx context.Context, msgView bot.MessageView, botID string) error {
	// I handle PLAIN_TEXT message only and make sure respond to current conversation.
	if msgView.Category == bot.MessageCategoryPlainText &&
		msgView.ConversationId == bot.UniqueConversationId(config.MixinClientID, msgView.UserId) {
		var data []byte
		var err error
		if data, err = base64.StdEncoding.DecodeString(msgView.Data); err != nil {
			log.Panicf("Error: %s\n", err)
			return err
		}
		inst := string(data)
		log.Printf("I got a message from %s, it said: `%s`\n", msgView.UserId, inst)

		if "sync" == inst {
			// Sync? Ack!
			Respond(ctx, msgView, "ack")
		} else if "hello" == inst {
			// Hello? Give you some money!
			Transfer(ctx, msgView)
		} else {
			Respond(ctx, msgView, defaultResponse)
		}
	}
	return nil
}

// Transfer 1.024 CNB to the user who having a conversation with bot.
func Transfer(ctx context.Context, msgView bot.MessageView) {
	payload := bot.TransferInput{
		AssetId:     CNBAssetID,
		RecipientId: msgView.UserId,
		Amount:      number.FromString("1.024"),
		TraceId:     uuid.Must(uuid.NewV4()).String(),
		Memo:        "Hello world",
	}
	err := bot.CreateTransfer(ctx, &payload,
		config.MixinClientID,
		config.MixinSessionID,
		config.MixinPrivateKey,
		config.MixinPin,
		config.MixinPinToken,
	)
	if err != nil {
		Respond(ctx, msgView, fmt.Sprintf("Oops, %s\n", err))
	}
}

// Respond to user.
func Respond(ctx context.Context, msgView bot.MessageView, msg string) {
	if err := client.SendPlainText(ctx, msgView, msg); err != nil {
		log.Panicf("Error: %s\n", err)
	}
}

func main() {
	ctx := context.Background()
	log.Println("start bot")
	handler := Handler{}

	// Create a bot client
	client = bot.NewBlazeClient(config.MixinClientID, config.MixinSessionID, config.MixinPrivateKey)

	// Start the loop
	for {
		if err := client.Loop(ctx, handler); err != nil {
			log.Printf("Error: %v\n", err)
		}
		log.Println("connection loop end")
		time.Sleep(time.Second)
	}
}
