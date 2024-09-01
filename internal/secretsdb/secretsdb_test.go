package secretsdb_test

import (
	"keepassui/internal/secretsdb"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadSecretsDBFromDBBytes(t *testing.T) {

	bytesContent, err := os.ReadFile("testdata/files/db.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}

	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(bytesContent, "keepassui")

	entriesGroupedByPath := secretsDB.EntriesByPath
	pathsInOrder := secretsDB.PathsInOrder

	if err != nil {
		t.Fatal("We don't expect any errors reading the KeepassDB")
	}

	assert.Equal(t, 3, len(entriesGroupedByPath))

	assert.Equal(t, []string{"Root", "Root|group 1", "Root|group 2"}, pathsInOrder)
	keysFromEntriesGroupedByPath := make([]string, 0, len(entriesGroupedByPath))
	for k := range entriesGroupedByPath {
		keysFromEntriesGroupedByPath = append(keysFromEntriesGroupedByPath, k)
	}

	assert.Contains(t, keysFromEntriesGroupedByPath, "Root|group 1")
	assert.Contains(t, keysFromEntriesGroupedByPath, "Root|group 2")

	entriesForRoot := entriesGroupedByPath["Root"]
	entriesForGroup1 := entriesGroupedByPath["Root|group 1"]
	entriesForGroup2 := entriesGroupedByPath["Root|group 2"]

	assert.Equal(t, 3, len(entriesForRoot))
	assert.Equal(t, 2, len(entriesForGroup1))
	assert.Equal(t, 1, len(entriesForGroup2))

	assert.Contains(t, entriesForRoot, secretsdb.SecretEntry{
		Group: "Root", Title: "keepassui example",
		Username: "keepassui", Password: "keepassui_password",
		Url: "https://fakekeepassuiurl.com", Notes: "This is an example", Path: []string{"Root"},
	})

	assert.Contains(t, entriesForRoot, secretsdb.SecretEntry{
		Group: "Root", Title: "group 1",
		Notes: "", Path: []string{"Root"}, IsGroup: true,
	})

	assert.Contains(t, entriesForRoot, secretsdb.SecretEntry{
		Group: "Root", Title: "group 2",
		Notes: "", Path: []string{"Root"}, IsGroup: true,
	})

	assert.Contains(t, entriesForGroup1, secretsdb.SecretEntry{
		Group: "Root|group 1", Title: "entry_inside_group1",
		Username: "user_in_group1", Password: "password_in_group_1",
		Url: "https://ingroup1.com/", Notes: "", Path: []string{"Root", "group 1"},
	})

	assert.Contains(t, entriesForGroup1, secretsdb.SecretEntry{
		Group: "Root|group 1", Title: "entry_2_in_group_1",
		Username: "entry2_group1_username", Password: "entry2_group1_password",
		Url: "entry2_group1_url", Notes: "", Path: []string{"Root", "group 1"},
	})

	assert.Contains(t, entriesForGroup2, secretsdb.SecretEntry{
		Group: "Root|group 2", Title: "entry_in_group2",
		Username: "user_in_group2", Password: "password_in_group2",
		Url: "https://group2.com", Notes: "", Path: []string{"Root", "group 2"},
	})
}

func Test_ReadSecretsDBFromDBBytes_Broken_File(t *testing.T) {

	bytesContent, err := os.ReadFile("testdata/files/db_broken.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}

	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(bytesContent, "keepassui")

	if err == nil {
		t.Fatal("We expect an error in this test because the DB file is broken")
	}
	if secretsDB.EntriesByPath != nil || secretsDB.PathsInOrder != nil {
		t.Fatal("entriesGroupedByPath and pathsInOrder should be nil")
	}

	assert.EqualError(t, err, "failed to verify HMAC for block 0")
}

func Test_writeDBBytes(t *testing.T) {
	secretsDB := secretsDBForTesting()

	bytes, err := secretsDB.WriteDBBytes("master")

	if err != nil {
		t.Error("Should not error")
	}

	secretsDBReadFromNewDBBytes, err := secretsdb.ReadSecretsDBFromDBBytes(bytes, "master")

	if err != nil {
		t.Error("Should be able to read secrets")
	}

	assert.Equal(t, secretsDB, secretsDBReadFromNewDBBytes)
}

func Test_AddSecretEntryInANewPath(t *testing.T) {
	secretsDB := secretsdb.SecretsDB{PathsInOrder: []string{}, EntriesByPath: make(map[string][]secretsdb.SecretEntry)}

	secretsDB.AddSecretEntry(secretsdb.SecretEntry{
		Group: "", Title: "Root",
		Notes: "Notes are important", Path: []string{}, IsGroup: true,
	})

	assert.Equal(t, 1, len(secretsDB.PathsInOrder))
	assert.Equal(t, []string{"Root"}, secretsDB.PathsInOrder)

	assert.Equal(t, 1, len(secretsDB.EntriesByPath))

	entries, ok := secretsDB.EntriesByPath[""]

	assert.Equal(t, 1, len(entries))

	assert.True(t, ok)

	assert.Equal(t, "", entries[0].Group)
	assert.Equal(t, "Root", entries[0].Title)
	assert.Equal(t, "Notes are important", entries[0].Notes)
	assert.Equal(t, "", entries[0].Username)
	assert.Equal(t, "", entries[0].Password)
	assert.Equal(t, "", entries[0].Url)
	assert.Equal(t, []string{}, entries[0].Path)
	assert.True(t, entries[0].IsGroup)
}

func Test_AddSecretEntryInExistingPath(t *testing.T) {
	secretsDB := secretsDBForTesting()

	secretsDB.AddSecretEntry(secretsdb.SecretEntry{
		Group: "Root|G1|G2", Title: "second_entry_in_RG2",
		Username: "second_user_in_RG2", Password: "second_password_in_RG2",
		Url: "https://secondRG1G2.com", Notes: "", Path: []string{"Root", "G1", "G2"},
	})

	assert.Equal(t, 3, len(secretsDB.PathsInOrder))
	assert.Equal(t, []string{"Root", "Root|G1", "Root|G1|G2"}, secretsDB.PathsInOrder)

	assert.Equal(t, 3, len(secretsDB.EntriesByPath))

	entriesForRootG1G2, ok := secretsDB.EntriesByPath["Root|G1|G2"]

	assert.Equal(t, 2, len(entriesForRootG1G2))

	assert.True(t, ok, "We expect to find entries in Root|G1|G2")

	assert.Equal(t, "Root|G1|G2", entriesForRootG1G2[0].Group)
	assert.Equal(t, "entry_in_RG2", entriesForRootG1G2[0].Title)
	assert.Equal(t, "", entriesForRootG1G2[0].Notes)
	assert.Equal(t, "user_in_RG2", entriesForRootG1G2[0].Username)
	assert.Equal(t, "password_in_RG2", entriesForRootG1G2[0].Password)
	assert.Equal(t, "https://RG1G2.com", entriesForRootG1G2[0].Url)
	assert.Equal(t, []string{"Root", "G1", "G2"}, entriesForRootG1G2[0].Path)
	assert.False(t, entriesForRootG1G2[0].IsGroup)

	assert.Equal(t, "Root|G1|G2", entriesForRootG1G2[1].Group)
	assert.Equal(t, "second_entry_in_RG2", entriesForRootG1G2[1].Title)
	assert.Equal(t, "", entriesForRootG1G2[1].Notes)
	assert.Equal(t, "second_user_in_RG2", entriesForRootG1G2[1].Username)
	assert.Equal(t, "second_password_in_RG2", entriesForRootG1G2[1].Password)
	assert.Equal(t, "https://secondRG1G2.com", entriesForRootG1G2[1].Url)
	assert.Equal(t, []string{"Root", "G1", "G2"}, entriesForRootG1G2[1].Path)
	assert.False(t, entriesForRootG1G2[1].IsGroup)
}

func Test_ModifySecretEntryInExistingPath(t *testing.T) {
	secretsDB := secretsDBForTesting()

	secretsDB.AddSecretEntry(secretsdb.SecretEntry{
		Group: "Root|G1|G2", Title: "entry_in_RG2",
		Username: "updated_user_in_RG2", Password: "updated_password_in_RG2",
		Url: "https://updatedRG1G2.com", Notes: "", Path: []string{"Root", "G1", "G2"},
	})

	assert.Equal(t, 3, len(secretsDB.PathsInOrder))
	assert.Equal(t, []string{"Root", "Root|G1", "Root|G1|G2"}, secretsDB.PathsInOrder)

	assert.Equal(t, 3, len(secretsDB.EntriesByPath))

	entriesForRootG1G2, ok := secretsDB.EntriesByPath["Root|G1|G2"]

	assert.Equal(t, 1, len(entriesForRootG1G2))

	assert.True(t, ok, "We expect to find entries in Root|G1|G2")

	assert.Equal(t, "Root|G1|G2", entriesForRootG1G2[0].Group)
	assert.Equal(t, "entry_in_RG2", entriesForRootG1G2[0].Title)
	assert.Equal(t, "", entriesForRootG1G2[0].Notes)
	assert.Equal(t, "updated_user_in_RG2", entriesForRootG1G2[0].Username)
	assert.Equal(t, "updated_password_in_RG2", entriesForRootG1G2[0].Password)
	assert.Equal(t, "https://updatedRG1G2.com", entriesForRootG1G2[0].Url)
	assert.Equal(t, []string{"Root", "G1", "G2"}, entriesForRootG1G2[0].Path)
	assert.False(t, entriesForRootG1G2[0].IsGroup)
}

func TestDeleteSecretEntry_ShouldNotDeleteAnEntryThatDoesntExist(t *testing.T) {
	secretsDB := secretsDBForTesting()

	deleted := secretsDB.DeleteSecretEntry(secretsdb.SecretEntry{
		Group: "Root|G1", Title: "non_existing_entry_in_RG1",
		Username: "user_in_RG1", Password: "password_in_RG1",
		Url: "https://RG1.com", Notes: "", Path: []string{"Root", "G1"},
	})

	assert.False(t, deleted, "Should be false as it didn't delete any entry")
}

func TestDeleteSecretEntry_ShouldDeleteAnEntryThatExists(t *testing.T) {
	secretsDB := secretsDBForTesting()

	entriesForRootG1 := secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 3, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "entry_in_RG1_2", entriesForRootG1[1].Title)
	assert.Equal(t, "G2", entriesForRootG1[2].Title)

	assert.Equal(t, []string{"Root", "Root|G1", "Root|G1|G2"}, secretsDB.PathsInOrder)

	deleted := secretsDB.DeleteSecretEntry(secretsdb.SecretEntry{
		Group: "Root|G1", Title: "entry_in_RG1_2",
		Username: "user_in_RG1_2", Password: "password_in_RG1_2",
		Url: "https://RG1_2.com", Notes: "", Path: []string{"Root", "G1"},
	})

	assert.True(t, deleted, "Should delete 1 entry")

	entriesForRootG1 = secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 2, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "G2", entriesForRootG1[1].Title)

	assert.Equal(t, []string{"Root", "Root|G1", "Root|G1|G2"}, secretsDB.PathsInOrder)
}

func TestDeleteSecretEntry_ShouldDeleteGroupAndItsContents(t *testing.T) {
	secretsDB := secretsDBForTesting()

	entriesForRoot := secretsDB.EntriesByPath["Root"]

	assert.Equal(t, 2, len(entriesForRoot))

	assert.Equal(t, "entry_in_root", entriesForRoot[0].Title)
	assert.Equal(t, "G1", entriesForRoot[1].Title)

	assert.Equal(t, []string{"Root", "Root|G1", "Root|G1|G2"}, secretsDB.PathsInOrder)

	entriesForRootG1 := secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 3, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "entry_in_RG1_2", entriesForRootG1[1].Title)
	assert.Equal(t, "G2", entriesForRootG1[2].Title)

	deleted := secretsDB.DeleteSecretEntry(secretsdb.SecretEntry{
		Group: "Root", Title: "G1",
		Notes: "", Path: []string{"Root"}, IsGroup: true,
	})

	assert.True(t, deleted, "Should delete the group and its contents")

	_, ok := secretsDB.EntriesByPath["Root|G1"]

	assert.False(t, ok, "G1 should not have entries")

	entriesForRoot = secretsDB.EntriesByPath["Root"]

	assert.Equal(t, 1, len(entriesForRoot))

	assert.Equal(t, "entry_in_root", entriesForRoot[0].Title)

	assert.Equal(t, []string{"Root"}, secretsDB.PathsInOrder)
}

func secretsDBForTesting() secretsdb.SecretsDB {
	entriesByPath := make(map[string][]secretsdb.SecretEntry)

	entriesByPath["Root"] = []secretsdb.SecretEntry{
		{
			Group: "Root", Title: "entry_in_root",
			Username: "user_in_root", Password: "password_in_root",
			Url: "https://rootEntry.com", Notes: "", Path: []string{"Root"},
		},
		{
			Group: "Root", Title: "G1",
			Notes: "", Path: []string{"Root"}, IsGroup: true,
		},
	}

	entriesByPath["Root|G1"] = []secretsdb.SecretEntry{
		{
			Group: "Root|G1", Title: "entry_in_RG1",
			Username: "user_in_RG1", Password: "password_in_RG1",
			Url: "https://RG1.com", Notes: "", Path: []string{"Root", "G1"},
		},
		{
			Group: "Root|G1", Title: "entry_in_RG1_2",
			Username: "user_in_RG1_2", Password: "password_in_RG1_2",
			Url: "https://RG1_2.com", Notes: "", Path: []string{"Root", "G1"},
		},
		{
			Group: "Root|G1", Title: "G2",
			Notes: "", Path: []string{"Root", "G1"}, IsGroup: true,
		},
	}

	entriesByPath["Root|G1|G2"] = []secretsdb.SecretEntry{
		{
			Group: "Root|G1|G2", Title: "entry_in_RG2",
			Username: "user_in_RG2", Password: "password_in_RG2",
			Url: "https://RG1G2.com", Notes: "", Path: []string{"Root", "G1", "G2"},
		},
	}

	return secretsdb.SecretsDB{PathsInOrder: []string{"Root", "Root|G1", "Root|G1|G2"}, EntriesByPath: entriesByPath}
}
