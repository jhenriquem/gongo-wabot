package helpers

import (
	"fmt"
	"time"
)

func FormatDate() string {
	// Define o fuso horário de Fortaleza
	loc, err := time.LoadLocation("America/Fortaleza")
	if err != nil {
		fmt.Println("Erro ao carregar o fuso horário:", err)
	}

	// Obtém a data e hora atual no fuso horário correto
	now := time.Now().In(loc)
	formattedDt := now.Format("02/01/2006 15:04:05")

	// var month string = fmt.Sprintf("%s", dt.Month())
	// if int(dt.Month()) < 10 {
	// month = fmt.Sprintf("0%d", dt.Month())
	// }
	return formattedDt
}
