package keepass

//go:generate mockgen -destination=../mocks/keepass/mock_keepass.go -source=./keepass.go

import (
	"bytes"
	"slices"

	"github.com/tobischo/gokeepasslib/v3"
)

type SecretEntry struct {
	Path     string
	Title    string
	Username string
	Password string
	Url      string
	Notes    string
}

type SecretsDB struct {
	EntriesByPath map[string][]SecretEntry
	PathsInOrder  []string
}

type CipheredKeepassDB struct {
	DBBytes  []byte
	Password string
	UriID    string
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() (SecretsDB, error)
}

func (ckdb CipheredKeepassDB) ReadEntriesFromContentGroupedByPath() (SecretsDB, error) {
	secrets, err := ckdb.readEntriesFromContent()

	if err != nil {
		return SecretsDB{}, err
	}

	return groupSecrets(secrets), nil
}

func (ckdb CipheredKeepassDB) readEntriesFromContent() ([]SecretEntry, error) {
	file := bytes.NewReader(ckdb.DBBytes)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(ckdb.Password)
	err := gokeepasslib.NewDecoder(file).Decode(db)

	if err != nil {
		return nil, err
	}

	err = db.UnlockProtectedEntries()

	if err != nil {
		return nil, err
	}

	var secrets []SecretEntry

	for _, g := range db.Content.Root.Groups {
		for _, group := range g.Groups {
			for _, entry := range group.Entries {
				outputEntry := SecretEntry{Title: entry.GetTitle(), Password: entry.GetPassword(), Path: group.Name}
				for _, value := range entry.Values {
					if value.Key == "UserName" {
						outputEntry.Username = value.Value.Content
					}
					if value.Key == "Notes" {
						outputEntry.Notes = value.Value.Content
					}
					if value.Key == "URL" {
						outputEntry.Url = value.Value.Content
					}

				}
				secrets = append(secrets, outputEntry)
			}
		}
	}

	return secrets, nil
}

func groupSecrets(secrets []SecretEntry) SecretsDB {
	secretsGroupedByPath := make(map[string][]SecretEntry)
	pathsInOrder := []string{}

	for _, p := range secrets {
		secretsGroupedByPath[p.Path] = append(secretsGroupedByPath[p.Path], p)
		if !slices.Contains(pathsInOrder, p.Path) {
			pathsInOrder = append(pathsInOrder, p.Path)
		}
	}
	return SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  pathsInOrder,
	}
}
