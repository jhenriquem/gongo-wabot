package utils

import (
	"fmt"
	"strings"

	"github.com/jhenriquem/gongo-wabot/internal/helpers"
	"github.com/jhenriquem/gongo-wabot/internal/types"
	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

// Função responsavel por enviar o pedido para a loja
func SubmitOrder(botClient *whatsmeow.Client, clientOrder types.Order) error {
	submitMessage := []string{
		"📚 *NOVO PEDIDO* 📃 \n",
		"",
		"*-------- Dados do cliente --------*",
		"",
		fmt.Sprintf("*Nome do cliente:* %s", clientOrder.Client.Name),
		fmt.Sprintf("*Número do cliente:* %s", clientOrder.Client.PhoneNumber),
		fmt.Sprintf("*Data:* %s", clientOrder.Data),
		"",
		"*-------- Livro --------*",
		"",
		fmt.Sprintf("*Título:* %s", clientOrder.Book.Title),
		fmt.Sprintf("*Autor:* %s", clientOrder.Book.Author),
		"",
		fmt.Sprintf("📦 *Quantidade de exemplares:* %d", clientOrder.Book.Amount),
		fmt.Sprintf("💰 *Valor por unidade:* R$ %.2f", clientOrder.Book.Value),
		fmt.Sprintf("💰 *Preço final:* R$ %.2f", clientOrder.Price),
		"",
		"*-------- Endereço e pagamento --------*",
		"",
		fmt.Sprintf("📍 *Endereço de entrega:* %s", clientOrder.Address),
		fmt.Sprintf("🪙 *Forma de pagamento:* %s", clientOrder.PaymentMethod),
		"",
		"",
		"*Mais um pedido concluído 🫡, tmj 🤝* ",
	}

	chat := waTypes.JID{
		User:   "",
		Server: "s.whatsapp.net",
	}
	if err := helpers.SendMessage(botClient, chat, strings.Join(submitMessage, "\n")); err != nil {
		return err
	}
	return nil
}
