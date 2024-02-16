package secretsreader

import "keepassui/internal/secretsdb"

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type CipheredKeepassDB struct {
	DBBytes  []byte
	Password string
	UriID    string
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error)
}

func (ckdb CipheredKeepassDB) ReadEntriesFromContentGroupedByPath() (secretsdb.SecretsDB, error) {
	return secretsdb.ReadSecretsDBFromDBBytes(ckdb.DBBytes, ckdb.Password)
}
