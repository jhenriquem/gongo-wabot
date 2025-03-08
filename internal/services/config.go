package services

import (
	"database/sql"

	storage_go "github.com/supabase-community/storage-go"
)

type svc struct {
	Database      *sql.DB
	StorageClient *storage_go.Client
}

var SVC svc
