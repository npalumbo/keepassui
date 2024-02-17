package secretsreader

import (
	"keepassui/internal/secretsdb"
)

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type DBPathAndPassword struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

var loadedDB *secretsdb.SecretsDB

func CreateDefaultSecretReader(uriID string, contentInBytes []byte, password string) SecretReader {
	return DBPathAndPassword{
		UriID:          uriID,
		ContentInBytes: contentInBytes,
		Password:       password,
	}
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error) //TODO remove return of secretsDB
	GetUriID() string
	GetFirstPath() string
	GetEntriesForPath(path string) []secretsdb.SecretEntry
	WriteDBBytes() ([]byte, error)
	AddSecretEntry(secretEntry secretsdb.SecretEntry)
	DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool
}

func (ckdb DBPathAndPassword) ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error) {
	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(ckdb.ContentInBytes, ckdb.Password)
	if err == nil {
		loadedDB = &secretsDB
	}
	return secretsDB, err
}

func (ckdb DBPathAndPassword) GetUriID() string {
	return ckdb.UriID
}

func (ckdb DBPathAndPassword) GetFirstPath() string {
	return loadedDB.PathsInOrder[0]
}

func (ckdb DBPathAndPassword) GetEntriesForPath(path string) []secretsdb.SecretEntry {
	return loadedDB.EntriesByPath[path]
}

func (ckdb DBPathAndPassword) WriteDBBytes() ([]byte, error) {
	return loadedDB.WriteDBBytes(ckdb.Password)
}

func (ckdb DBPathAndPassword) AddSecretEntry(secretEntry secretsdb.SecretEntry) {
	loadedDB.AddSecretEntry(secretEntry)
}

func (ckdb DBPathAndPassword) DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool {
	return loadedDB.DeleteSecretEntry(secretEntry)
}
