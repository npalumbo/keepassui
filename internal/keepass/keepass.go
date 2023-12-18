package keepass

import (
	"bytes"

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

func ReadEntriesFromContent(contentInBytes []byte, password string) ([]SecretEntry, error) {
	file := bytes.NewReader(contentInBytes)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(password)
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
