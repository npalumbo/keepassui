package keepass

//go:generate mockgen -destination=../mocks/keepass/mock_keepass.go -source=./keepass.go

import (
	"bytes"
	"slices"
	"strings"

	"github.com/tobischo/gokeepasslib/v3"
	w "github.com/tobischo/gokeepasslib/v3/wrappers"
)

type SecretEntry struct {
	Group    string
	Path     []string
	Title    string
	Username string
	Password string
	Url      string
	Notes    string
	IsGroup  bool
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
	reader := bytes.NewReader(ckdb.DBBytes)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(ckdb.Password)
	err := gokeepasslib.NewDecoder(reader).Decode(db)

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

func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value, Protected: w.NewBoolWrapper(true)},
	}
}

func (secretsDB SecretsDB) WriteDBBytes(masterPassword string) ([]byte, error) {
	rootGroupName := secretsDB.PathsInOrder[0]

	groupsMap := make(map[string]*gokeepasslib.Group)

	for _, path := range secretsDB.PathsInOrder {
		newGroup := gokeepasslib.NewGroup()
		entries := secretsDB.EntriesByPath[path]
		groupName := getLatestGroupInPath(path)
		if _, ok := groupsMap[groupName]; !ok {
			for _, secretEntry := range entries {
				if !secretEntry.IsGroup {
					entry := gokeepasslib.NewEntry()
					entry.Values = append(entry.Values, mkValue("Title", secretEntry.Title))
					entry.Values = append(entry.Values, mkValue("UserName", secretEntry.Username))
					entry.Values = append(entry.Values, mkProtectedValue("Password", secretEntry.Password))
					entry.Values = append(entry.Values, mkValue("URL", secretEntry.Url))
					entry.Values = append(entry.Values, mkValue("Notes", secretEntry.Notes))
					newGroup.Entries = append(newGroup.Entries, entry)
					newGroup.Name = groupName
				}
			}
			groupsMap[groupName] = &newGroup
		}
	}

	// append groups to group if relevant
	connectionsMade := []string{}

	pathsInReverseOrder := reverseCopy(secretsDB.PathsInOrder)

	for _, path := range pathsInReverseOrder {
		if path != rootGroupName {
			pathHopsInReverseOrder := reverseCopy(strings.Split(path, "|"))

			for i := range pathHopsInReverseOrder {
				if i > 0 {
					groupName := pathHopsInReverseOrder[i]
					subGroupName := pathHopsInReverseOrder[i-1]
					groupConnection := groupName + "|" + subGroupName

					if !slices.Contains(connectionsMade, groupConnection) {

						mainGroupToAddSubGroup := groupsMap[groupName]
						subGroup := groupsMap[subGroupName]
						(*mainGroupToAddSubGroup).Groups = append(mainGroupToAddSubGroup.Groups, *subGroup)
					}
					connectionsMade = append(connectionsMade, groupConnection)
				}
			}
		}
	}

	db := &gokeepasslib.Database{
		Header:      gokeepasslib.NewHeader(),
		Credentials: gokeepasslib.NewPasswordCredentials(masterPassword),
		Content: &gokeepasslib.DBContent{
			Meta: gokeepasslib.NewMetaData(),
			Root: &gokeepasslib.RootData{
				Groups: []gokeepasslib.Group{*groupsMap[rootGroupName]},
			},
		},
	}

	// Lock entries using stream cipher
	err := db.LockProtectedEntries()
	if err != nil {
		return nil, err
	}

	// and encode it into the file
	buf := bytes.NewBuffer([]byte{})
	keepassEncoder := gokeepasslib.NewEncoder(buf)

	if err = keepassEncoder.Encode(db); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
		expandedPath := append(groupPath, groupToScan.Name)
		secrets = extractEntries(expandedPath, group, secrets)
		secrets = append(secrets, SecretEntry{Title: group.Name, Path: expandedPath, Group: strings.Join(expandedPath, "|"), IsGroup: true})
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

func getLatestGroupInPath(path string) string {
	parts := strings.Split(path, "|")

	if len(parts) > 1 {
		return parts[len(parts)-1]
	}

	return path
}

func reverseCopy(array []string) []string {
	reversed := make([]string, 0)
	for i := len(array) - 1; i >= 0; i-- {
		reversed = append(reversed, array[i])
	}
	return reversed
}
