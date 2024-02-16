package secretsreader

import "keepassui/internal/secretsdb"

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type DBPathAndPassword struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

func CreateDefaultSecretReader(uriID string, contentInBytes []byte, password string) SecretReader {
	return DBPathAndPassword{
		UriID:          uriID,
		ContentInBytes: contentInBytes,
		Password:       password,
	}
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error)
	GetUriID() string
	GetContentInBytes() []byte
	GetPassword() string
}

func (ckdb DBPathAndPassword) ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error) {
	return secretsdb.ReadSecretsDBFromDBBytes(ckdb.ContentInBytes, ckdb.Password)
}

func (ckdb DBPathAndPassword) GetUriID() string {
	return ckdb.UriID
}

func (ckdb DBPathAndPassword) GetContentInBytes() []byte {
	return ckdb.ContentInBytes
}

// TODO see if I can find a way to stop exposing password
func (ckdb DBPathAndPassword) GetPassword() string {
	return ckdb.Password
}
