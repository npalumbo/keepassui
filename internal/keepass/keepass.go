package keepass

//go:generate mockgen -destination=../mocks/keepass/mock_keepass.go -source=./keepass.go

import (
	"bytes"
	"slices"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
)

type SecretEntry struct {
	Group    string
	Path     []string
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
		secrets = extractEntries([]string{}, g, secrets)
	}

	return secrets, nil
}

func extractEntries(groupPath []string, groupToScan gokeepasslib.Group, secrets []SecretEntry) []SecretEntry {
	for _, entry := range groupToScan.Entries {
		outputEntry := SecretEntry{Title: entry.GetTitle(), Password: entry.GetPassword(), Group: groupToScan.Name}
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

			outputEntry.Path = append(groupPath, groupToScan.Name)
			outputEntry.Group = strings.Join(outputEntry.Path, "|")

		}
		secrets = append(secrets, outputEntry)
	}

	for _, group := range groupToScan.Groups {
		secrets = extractEntries(append(groupPath, groupToScan.Name), group, secrets)
	}
	return secrets
}

func groupSecrets(secrets []SecretEntry) SecretsDB {
	secretsGroupedByPath := make(map[string][]SecretEntry)
	pathsInOrder := []string{}

	for _, p := range secrets {
		secretsGroupedByPath[p.Group] = append(secretsGroupedByPath[p.Group], p)
		if !slices.Contains(pathsInOrder, p.Group) {
			pathsInOrder = append(pathsInOrder, p.Group)
		}
	}
	return SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  pathsInOrder,
	}
}
