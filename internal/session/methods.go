package sessions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jhenriquem/gongo-wabot/internal/helpers"
	"github.com/jhenriquem/gongo-wabot/internal/services"
	"github.com/jhenriquem/gongo-wabot/internal/types"
	"github.com/jhenriquem/gongo-wabot/internal/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

func (s *Session) cancelOrder(botClient *whatsmeow.Client, Chat waTypes.JID, msg string) (bool, error) {
	if msg == "*" {
		botMessage := "🔚  *PEDIDO CANCELADO*  🔚"
		s.stageCounter = 0
		_, err := botClient.SendMessage(context.Background(), Chat, &waProto.Message{
			Conversation: &botMessage,
		})

		// Apagar a session da lista
		defer RemoveSession(s.UserID)

		// Retorna que o pedido foi cancelado
		return true, err
	}

	// Retorna que o pedido não foi cancelado
	return false, nil
}

func (s *Session) initialStage(botClient *whatsmeow.Client, v *events.Message) error {
	messageSlice := []string{
		"*👋 Olá, como vai?*",
		"",
		fmt.Sprintf("Eu sou <Nome do bot>, o *assistente virtual* da %s.", os.Getenv("STORE_NAME")),
		"*Como posso te ajudar?*",
		"-----------------------------------",
		"1️⃣ - FAZER UM PEDIDO",
		"0️⃣ - FALAR COM NOSSA EQUIPE",
		"",
		"*Digite o número da opção:*",
	}

	botMessage := strings.Join(messageSlice, "\n")

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) choiceStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	// Se o usuário cancelar o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}

	var botMessage string

	switch clientMessage {
	case "1":
		botMessage = "*📗 Qual o seu livro ? Basta informar o título🖋️* \n\n  🔴 _( Caso queira cancelar o pedido, envie *** )_"
	case "0":
		botMessage = fmt.Sprintf("*Você escolheu falar com nossa equipe 🙋‍♂️*\n\nAqui está o número do nosso atendemento:\n%s  %s \n\nAté a próxima! 👋🐛", os.Getenv("ATTENDANT_PHONE"), os.Getenv("ATTENDANT_NAME"))

		if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}

		defer RemoveSession(s.UserID)
		return nil

	default:
		botMessage = "*Opção inválida ❌*\n\n-----------------------------------\n1️⃣ - FAZER UM PEDIDO\n0️⃣ - FALAR COM NOSSA EQUIPE\n\n*Digite o número da opção:* \n\n  🔴 _( Caso queira cancelar o pedido, envie *** )_"

		if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}
		return nil
	}

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) orderStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	// Se o usuário cancelar o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}
	var book types.Book

	if err := services.GetBook(&book, clientMessage); err != nil {
		var botMessage string
		if errors.Is(err, sql.ErrNoRows) {
			botMessage = "❌ *Livro não encontrado!* ❌\n\nVerifique se digitou o título corretamente e tente novamente. \n\n 🔴 _( Caso queira cancelar o pedido, envie *** )_"
		} else {
			log.Printf("Erro no banco de dados: %v\n", err)
			botMessage = "*Estamos enfrentando problemas no momento ⚠️ *"
			s.stageCounter = 0
		}

		if errSend := helpers.SendMessage(botClient, v.Info.Chat, botMessage); errSend != nil {
			return errors.New(err.Error() + "\n" + errSend.Error())
		}

		return nil
	}

	// Quando não houver estoque
	if book.Amount == 0 {
		botMessage := fmt.Sprintf("*⚠️  Infelizmente estamos sem esse livro em estoque* \n %s - %s \n Caso queira escolher outro livro, basta dizer o título\n\n 🔴 _( Caso queira cancelar o pedido, envie *** )_", book.Title, book.Author)

		if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}

		log.Printf("\nCompra impedida por falta de estoque : \n (%s) \n %s ", s.UserID, book.Title)

		// Voltar ao estagio de escolhar de livro
		s.stageCounter--
		return nil
	}

	// Livro encontrado, salvar no pedido
	s.currentOrder.Book = book

	// Criar a mensagem de resposta
	botMessage := strings.Join([]string{
		"📚 *Encontramos o seu livro!* 📘\n",
		fmt.Sprintf("*Título:* %s", book.Title),
		fmt.Sprintf("*Autor:* %s", book.Author),
		"",
		fmt.Sprintf("💰 *Valor:* R$ %.2f", book.Value),
		fmt.Sprintf("📦 *Quantidade disponível:* %d", book.Amount),
		"",
		"-----------------------",
		"📌 *Quantos exemplares deseja?*",
		"_(Digite apenas o número)_",
		"",
		"",
		"_❌ Livro errado ? Mande *!*, depois mande o título do seu livro_ ",
		" 🔴 _( Caso queira cancelar o pedido, envie *** )_ ",
	}, "\n")

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) amountStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	var botMessage string

	if clientMessage == "!" {
		s.stageCounter--
		return nil
	}

	// Caso o usuário cancele o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}

	// Tentativa de conversão para inteiro
	amount, err := strconv.Atoi(clientMessage)
	if err != nil {
		botMessage = "*❌ Entrada inválida!*\n\nDigite um número válido de exemplares ⚠️ \n\n  🔴 _( Caso queira cancelar o pedido, envie *** )_"

		if err = helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}

		return nil
	}

	if amount <= 0 || amount > s.currentOrder.Book.Amount {
		botMessage = fmt.Sprintf(
			"*⚠️ Estoque insuficiente!*\n\nTemos apenas *%d* unidades disponíveis. \n\n  🔴 _( Caso queira cancelar o pedido, envie *** )_",
			s.currentOrder.Book.Amount,
		)

		if err = helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}
		return nil
	}

	// Atualiza a quantidade e o preço final do pedido
	s.currentOrder.Book.Amount = amount
	s.currentOrder.Price = s.currentOrder.Book.Value * float64(amount)

	botMessage = "*Muito bem! 👍*\n\n🗺️ Agora, informe seu *📍 ENDEREÇO*:\n\n`Rua, Número, Bairro, Cidade`\n\n 🔴 _( Caso queira cancelar o pedido, envie *** )_"

	if err = helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) addressStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	// Caso o usuário cancele o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}

	// Armazena o endereço no pedido
	s.currentOrder.Address = clientMessage

	// Mensagem de confirmação do pedido
	botMessage := strings.Join([]string{
		"*🔔 Seu pedido está quase finalizado 🔔*\n",
		fmt.Sprintf("*📙 %s*", s.currentOrder.Book.Title),
		fmt.Sprintf("*✒️ %s*", s.currentOrder.Book.Author),
		fmt.Sprintf("\n💵 Valor do exemplar: R$ %.2f", s.currentOrder.Book.Value),
		fmt.Sprintf("📦 Quantidade: %d", s.currentOrder.Book.Amount),
		fmt.Sprintf("*💰 Preço final: R$ %.2f*", s.currentOrder.Price),
		fmt.Sprintf("\n📍 *Endereço:* %s", s.currentOrder.Address),
	}, "\n")

	// Opções de pagamento
	botMessage += "\n\n🔊 Agora, informe a forma de pagamento 🪙\n" +
		"-----------------\n" +
		"1️⃣ - Pix\n" +
		"2️⃣ - Dinheiro\n" +
		"3️⃣ - Cartão de Débito\n" +
		"4️⃣ - Cartão de Crédito\n\n" +
		" 🔴 _( Caso queira cancelar o pedido, envie *** )_"

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) paymentStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	// Caso o usuário cancele o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}

	var botMessage string

	switch clientMessage {
	case "1":
		s.currentOrder.PaymentMethod = "Pix"
		botMessage = fmt.Sprintf(
			"✅ *Você escolheu pagamento via Pix*\n\n"+
				"🔑 Essa é nossa chave Pix: %s\n\n"+
				"📃 Depois nos mande o comprovante 👍\n\n"+
				"✍️ *Agora, diga-nos o seu nome*",
			os.Getenv("PIX_KEY"),
		)
	case "2":
		s.currentOrder.PaymentMethod = "Dinheiro"
		botMessage = "💵 *Você escolheu pagamento via dinheiro*\n\n" +
			"💰 O pagamento será realizado na entrega do livro.\n\n" +
			"✍️ *Agora, diga-nos o seu nome*"
	case "3":
		s.currentOrder.PaymentMethod = "Cartão (Débito)"
		botMessage = "💳 *Você escolheu pagamento via cartão de débito*\n\n" +
			"📍 O pagamento será realizado na entrega do livro.\n\n" +
			"✍️ *Agora, diga-nos o seu nome*"
	case "4":
		s.currentOrder.PaymentMethod = "Cartão (Crédito)"
		botMessage = "💳 *Você escolheu pagamento via cartão de crédito*\n\n" +
			"📍 O pagamento será realizado na entrega do livro.\n\n" +
			"✍️ *Agora, diga-nos o seu nome*"
	default:
		botMessage = "⚠️ *Escolha uma opção de pagamento válida*:\n\n" +
			"1️⃣ - Pix\n" +
			"2️⃣ - Dinheiro\n" +
			"3️⃣ - Cartão de Débito\n" +
			"4️⃣ - Cartão de Crédito\n\n" +
			" 🔴 _( Caso queira cancelar o pedido, envie *** )_"

		if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
			return err
		}
	}

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	s.stageCounter++
	return nil
}

func (s *Session) finalStage(botClient *whatsmeow.Client, v *events.Message) error {
	var clientMessage string = helpers.GetTextMessage(v)

	// Caso o usuário cancele o pedido
	if cancel, err := s.cancelOrder(botClient, v.Info.Chat, clientMessage); cancel {
		return err
	}

	// Salvando os dados do cliente
	s.currentOrder.Client.Name = clientMessage
	s.currentOrder.Client.PhoneNumber = v.Info.Sender.User

	// Criando a mensagem de confirmação
	messageSlice := []string{
		fmt.Sprintf("*🔔 %s, seu pedido está finalizado! 🔔*\n", s.currentOrder.Client.Name),
		fmt.Sprintf("*📙 %s*", s.currentOrder.Book.Title),
		fmt.Sprintf("*✒️  %s*", s.currentOrder.Book.Author),
		"",
		fmt.Sprintf("📦 *Quantidade de exemplares:* %d", s.currentOrder.Book.Amount),
		fmt.Sprintf("💰 *Valor por unidade:* R$ %.2f", s.currentOrder.Book.Value),
		fmt.Sprintf("💰 *Preço final:* R$ %.2f", s.currentOrder.Price),
		"",
		fmt.Sprintf("📍 *Endereço de entrega:* %s", s.currentOrder.Address),
		fmt.Sprintf("🪙 *Forma de pagamento:* %s", s.currentOrder.PaymentMethod),
		"",
		"🚚 *Seu pedido será entregue no endereço informado!*",
		"🙏 *Obrigado pela preferência! 😊*",
		"",
		"🔚  *Atendimento encerrado*  🔚",
	}

	botMessage := strings.Join(messageSlice, "\n")

	if err := helpers.SendMessage(botClient, v.Info.Chat, botMessage); err != nil {
		return err
	}

	// Define a data em que o pedido foi executado
	s.currentOrder.Data = helpers.FormatDate()

	// Enviar o pedido para a loja
	err := utils.SubmitOrder(botClient, s.currentOrder)
	if err != nil {
		log.Printf("Erro ao enviar o pedido para a loja  (%s) : %s", s.UserID, err.Error())
	}

	// Atualização do estoque do livro
	err = services.SetBookStock(s.currentOrder.Book.ID, s.currentOrder.Book.Amount)
	if err != nil {
		log.Printf("Erro ao atualizar o estoque (%s) : %s", s.UserID, err.Error())
	}

	// Finaliza a seção
	defer RemoveSession(s.UserID)
	return nil
}
