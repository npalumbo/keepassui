package keepass

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

type CipheredKeepassDB struct {
	ContentInBytes []byte
	Password       string
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() (map[string][]SecretEntry, []string, error)
}

func (ckdb CipheredKeepassDB) ReadEntriesFromContent() ([]SecretEntry, error) {
	file := bytes.NewReader(ckdb.ContentInBytes)

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

func (ckdb CipheredKeepassDB) ReadEntriesFromContentGroupedByPath() (map[string][]SecretEntry, []string, error) {
	secrets, err := ckdb.ReadEntriesFromContent()

	if err != nil {
		return nil, nil, err
	}

	secretsGroupedByPath, pathsInOrder := groupSecrets(secrets)
	return secretsGroupedByPath, pathsInOrder, nil
}

func groupSecrets(secrets []SecretEntry) (secretsGroupedByPath map[string][]SecretEntry, pathsInOrder []string) {
	secretsGroupedByPath = make(map[string][]SecretEntry)
	pathsInOrder = []string{}

	for _, p := range secrets {
		secretsGroupedByPath[p.Path] = append(secretsGroupedByPath[p.Path], p)
		if !slices.Contains(pathsInOrder, p.Path) {
			pathsInOrder = append(pathsInOrder, p.Path)
		}
	}
	return secretsGroupedByPath, pathsInOrder
}
