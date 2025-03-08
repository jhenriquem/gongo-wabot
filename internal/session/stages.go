package sessions

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Metodo para a alteração do stage atual e do contador
func (s *Session) GetCurrentStage(client *whatsmeow.Client, v *events.Message) error {
	switch s.stageCounter {
	case 1:
		return s.choiceStage(client, v)
	case 2:
		return s.orderStage(client, v)
	case 3:
		return s.amountStage(client, v)
	case 4:
		return s.addressStage(client, v)
	case 5:
		return s.paymentStage(client, v)
	case 6:
		return s.finalStage(client, v)
	default:
		return s.initialStage(client, v)
	}
}
