package ui_test

import (
	"keepassui/internal/secretsdb"
)

func secretsDBForTesting() secretsdb.SecretsDB {
	secretsGroupedByPath := make(map[string][]secretsdb.SecretEntry)
	secretsGroupedByPath["path 1"] = []secretsdb.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}}
	paths := []string{"path 1"}
	secretsDB := secretsdb.SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  paths,
	}
	return secretsDB
}

func secretsDBWithTwoGroups() secretsdb.SecretsDB {
	secretsGroupedByPath := make(map[string][]secretsdb.SecretEntry)
	secretsGroupedByPath["path 1"] = []secretsdb.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"},
		{Title: "path 2", Group: "path 1", Notes: "", IsGroup: true}}
	secretsGroupedByPath["path 2"] = []secretsdb.SecretEntry{
		{Title: "title 2", Group: "path 2", Username: "username 2", Password: "password 2", Url: "url 2", Notes: "notes 2"},
		{Title: "title 3", Group: "path 2", Username: "username 3", Password: "password 3", Url: "url 3", Notes: "notes 3"},
	}
	paths := []string{"path 1", "path 2"}
	secretsDB := secretsdb.SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  paths,
	}
	return secretsDB
}
