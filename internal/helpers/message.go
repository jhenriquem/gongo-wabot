package helpers

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func SendMessage(botClient *whatsmeow.Client, chat types.JID, message string) error {
	_, err := botClient.SendMessage(context.Background(), chat, &waProto.Message{
		Conversation: &message,
	})
	return err
}

func GetTextMessage(message *events.Message) string {
	var msg *waE2E.Message = message.Message

	if val := msg.GetConversation(); val != "" {
		return val
	} else if val := msg.GetExtendedTextMessage().GetText(); val != "" {
		return val
	}

	return ""
}
