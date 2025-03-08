package utils

import (
	"context"
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func ErrorSendMessage(client *whatsmeow.Client, Chat types.JID, userID string, err error) {
	botMessage := "🛑 _Tivemos um erro ao tentar lhe responder, tente de novo_"

	_, sendErr := client.SendMessage(context.Background(), Chat, &waE2E.Message{
		Conversation: &botMessage,
	})
	if sendErr != nil {
		log.Printf("\nErro ao enviar mensagem de erro para %s: %v", userID, sendErr)
	}

	log.Printf("\nErro em estágio da conversa: %s", err.Error())
}
