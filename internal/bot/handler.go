package bot

import (
	"log"
	"sync"

	"github.com/jhenriquem/gongo-wabot/internal/helpers"
	sessions "github.com/jhenriquem/gongo-wabot/internal/session"
	"github.com/jhenriquem/gongo-wabot/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Mutex para proteger acesso concorrente à lista de sessões
var mu sync.Mutex

// Verifica se a mensagem reciba é do tipo texto
func validMessageType(message *events.Message) bool {
	// Como helpers.GetTextMessage(message) só retorna algo caso  GetConversation() ou
	// GetExtendedTextMessage().GetText() retorne uma string não vazia, isto é
	// a mensagem recebida é o tipo texto essa verificação é perfeita
	if helpers.GetTextMessage(message) != "" {
		return true
	}
	return false
}

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			userID := v.Info.Chat.User
			senderName := v.Info.PushName

			// Verifica se a mensagem recebida e do tipo texto
			if !validMessageType(v) {
				return
			}

			var session *sessions.Session

			// Protege a lista de sessões com um Mutex
			mu.Lock()

			for _, s := range sessions.List {
				if s.UserID == userID {
					session = s
					break
				}
			}
			mu.Unlock() // Libera mutex

			if session != nil {
				go func(ss *sessions.Session) {
					if err := ss.GetCurrentStage(client, v); err != nil {
						utils.ErrorSendMessage(client, v.Info.Chat, userID, err)
					}
				}(session)
			} else {
				// Se não existir sessão, cria uma nova
				newSession := sessions.SetSession(userID)

				// Log
				log.Println("\n--------------------------------------------------------------")
				log.Printf("\n [ New session ] %s ( %s )", senderName, userID)

				go func(s *sessions.Session) {
					s.GetCurrentStage(client, v)
				}(newSession)
			}
		}
	}
}
