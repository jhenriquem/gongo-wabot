package services

import (
	"bytes"
	"fmt"
	"log"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

// Uma configuração para um storage online usando Supabase
// O arquivo é altomaticamente baixado na inicialização do bot e é feito o update ao final

func InitStorage() {
	SVC.StorageClient = storage_go.NewClient(os.Getenv("STORAGE_URL"), os.Getenv("STORAGE_KEY"), nil)
	DownloadSQLite()
}

func UploadSQLite() {
	data, err := os.ReadFile(os.Getenv("SQLITE_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	// Enviando para Supabase
	resp, err := SVC.StorageClient.UpdateFile("sqlite-backups", "store.db", bytes.NewReader(data))
	if err != nil {
		log.Fatal("Error ao salva o storage:", err)
	}

	fmt.Println("Storage salvo:", resp)
}

func DownloadSQLite() {
	data, err := SVC.StorageClient.DownloadFile("sqlite-backups", "./store.db")
	if err != nil {
		log.Fatal("Erro ao baixar o banco: ", err)
	}

	os.WriteFile(os.Getenv("SQLITE_PATH"), data, 0644)
}
