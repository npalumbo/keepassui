package secretsreader

import (
	"keepassui/internal/secretsdb"
	"slices"
	"strings"

	"fyne.io/fyne/v2/storage"
)

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type DefaultSecretsReader struct {
	UriID          string
	ContentInBytes []byte
	Password       string
	loadedDB       *secretsdb.SecretsDB
}

func CreateDefaultSecretsReader(uriID string, contentInBytes []byte, password string) DefaultSecretsReader {
	return DefaultSecretsReader{
		UriID:          uriID,
		ContentInBytes: contentInBytes,
		Password:       password,
	}
}

type SecretReader interface {
	GetUriID() string
	GetFirstPath() string
	GetEntriesForPath(path string) []secretsdb.SecretEntry
	WriteDBBytes() ([]byte, error)
	Save() error
	AddSecretEntry(secretEntry secretsdb.SecretEntry)
	ModifySecretEntry(originalTitle, originalGroup string, originalIsGroup bool, secretEntry secretsdb.SecretEntry)
	DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool
	CreateEmptyDBBytes(masterPassword string) ([]byte, error)
}

func (dsr *DefaultSecretsReader) ReadEntriesFromContentGroupedByPath() error {
	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(dsr.ContentInBytes, dsr.Password)
	if err == nil {
		dsr.loadedDB = &secretsDB
	}
	return err
}

func (dsr DefaultSecretsReader) GetUriID() string {
	return dsr.UriID
}

func (dsr DefaultSecretsReader) GetFirstPath() string {
	return dsr.loadedDB.PathsInOrder[0]
}

func (dsr DefaultSecretsReader) GetEntriesForPath(path string) []secretsdb.SecretEntry {
	return dsr.loadedDB.EntriesByPath[path]
}

func (dsr DefaultSecretsReader) WriteDBBytes() ([]byte, error) {
	return dsr.loadedDB.WriteDBBytes(dsr.Password)
}

func (dsr DefaultSecretsReader) AddSecretEntry(secretEntry secretsdb.SecretEntry) {
	dsr.loadedDB.AddSecretEntry(secretEntry)
}

func (dsr DefaultSecretsReader) Save() error {
	bytes, err := dsr.WriteDBBytes()

	if err != nil {
		return err
	}

	fURI, err := storage.ParseURI(dsr.GetUriID())

	if err != nil {
		return err
	}
	fileWC, err := storage.Writer(fURI)

	if err != nil {
		return err
	}
	_, err = fileWC.Write(bytes)

	if err != nil {
		return err
	}

	err = fileWC.Close()

	if err != nil {
		return err
	}
	return nil
}

func (dsr DefaultSecretsReader) ModifySecretEntry(originalTitle, originalGroup string, originalIsGroup bool, secretEntry secretsdb.SecretEntry) {
	entries := dsr.loadedDB.EntriesByPath[originalGroup]
	i := slices.IndexFunc(entries, func(se secretsdb.SecretEntry) bool {
		return se.Title == originalTitle && se.IsGroup == originalIsGroup && se.Group == originalGroup
	})
	if i != -1 {
		entries[i].Group = secretEntry.Group
		entries[i].Title = secretEntry.Title
		entries[i].Username = secretEntry.Username
		entries[i].Password = secretEntry.Password
		entries[i].Notes = secretEntry.Notes
		if originalIsGroup && originalTitle != secretEntry.Title {
			for i, path := range dsr.loadedDB.PathsInOrder {
				if strings.Contains(path, originalTitle) {
					dsr.loadedDB.PathsInOrder[i] = strings.Replace(path, originalTitle, secretEntry.Title, 1)
				}
			}

			for keypath, secrets := range dsr.loadedDB.EntriesByPath {
				for i, entry := range secrets {
					entry.Group = strings.Replace(entry.Group, originalTitle, secretEntry.Title, 1)
					secrets[i] = entry
				}
				replacedPath := strings.Replace(keypath, originalTitle, secretEntry.Title, 1)
				if replacedPath != keypath {
					dsr.loadedDB.EntriesByPath[replacedPath] = dsr.loadedDB.EntriesByPath[keypath]
					delete(dsr.loadedDB.EntriesByPath, keypath)
				}
			}
		}
	}
}

func (dsr DefaultSecretsReader) DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool {
	return dsr.loadedDB.DeleteSecretEntry(secretEntry)
}

func (dsr DefaultSecretsReader) CreateEmptyDBBytes(masterPassword string) ([]byte, error) {
	entriesByPath := make(map[string][]secretsdb.SecretEntry)
	entriesByPath["Root"] = []secretsdb.SecretEntry{}
	secretsDB := secretsdb.SecretsDB{EntriesByPath: entriesByPath, PathsInOrder: []string{"Root"}}
	return secretsDB.WriteDBBytes(masterPassword)
}
