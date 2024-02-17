package secretsreader

import (
	"keepassui/internal/secretsdb"
)

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type DefaultSecretsReader struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

var loadedDB *secretsdb.SecretsDB

func CreateDefaultSecretsReader(uriID string, contentInBytes []byte, password string) DefaultSecretsReader {
	return DefaultSecretsReader{
		UriID:          uriID,
		ContentInBytes: contentInBytes,
		Password:       password,
	}
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() error
	GetUriID() string
	GetFirstPath() string
	GetEntriesForPath(path string) []secretsdb.SecretEntry
	WriteDBBytes() ([]byte, error)
	AddSecretEntry(secretEntry secretsdb.SecretEntry)
	DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool
}

func (ckdb DefaultSecretsReader) ReadEntriesFromContentGroupedByPath() error {
	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(ckdb.ContentInBytes, ckdb.Password)
	if err == nil {
		loadedDB = &secretsDB
	}
	return err
}

func (ckdb DefaultSecretsReader) GetUriID() string {
	return ckdb.UriID
}

func (ckdb DefaultSecretsReader) GetFirstPath() string {
	return loadedDB.PathsInOrder[0]
}

func (ckdb DefaultSecretsReader) GetEntriesForPath(path string) []secretsdb.SecretEntry {
	return loadedDB.EntriesByPath[path]
}

func (ckdb DefaultSecretsReader) WriteDBBytes() ([]byte, error) {
	return loadedDB.WriteDBBytes(ckdb.Password)
}

func (ckdb DefaultSecretsReader) AddSecretEntry(secretEntry secretsdb.SecretEntry) {
	loadedDB.AddSecretEntry(secretEntry)
}

func (ckdb DefaultSecretsReader) DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool {
	return loadedDB.DeleteSecretEntry(secretEntry)
}
