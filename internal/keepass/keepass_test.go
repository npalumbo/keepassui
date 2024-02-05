package keepass_test

import (
	"keepassui/internal/keepass"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadEntriesFromContentGroupedByPath(t *testing.T) {

	bytesContent, err := os.ReadFile("testdata/files/db.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}

	cipheredKeepassDB := keepass.CipheredKeepassDB{DBBytes: bytesContent, Password: "keepassui"}

	secretsDB, err := cipheredKeepassDB.ReadEntriesFromContentGroupedByPath()

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

	assert.Contains(t, entriesForRoot, keepass.SecretEntry{
		Group: "Root", Title: "keepassui example",
		Username: "keepassui", Password: "keepassui_password",
		Url: "https://fakekeepassuiurl.com", Notes: "This is an example", Path: []string{"Root"},
	})

	assert.Contains(t, entriesForRoot, keepass.SecretEntry{
		Group: "Root", Title: "group 1",
		Notes: "", Path: []string{"Root"}, IsGroup: true,
	})

	assert.Contains(t, entriesForRoot, keepass.SecretEntry{
		Group: "Root", Title: "group 2",
		Notes: "", Path: []string{"Root"}, IsGroup: true,
	})

	assert.Contains(t, entriesForGroup1, keepass.SecretEntry{
		Group: "Root|group 1", Title: "entry_inside_group1",
		Username: "user_in_group1", Password: "password_in_group_1",
		Url: "https://ingroup1.com/", Notes: "", Path: []string{"Root", "group 1"},
	})

	assert.Contains(t, entriesForGroup1, keepass.SecretEntry{
		Group: "Root|group 1", Title: "entry_2_in_group_1",
		Username: "entry2_group1_username", Password: "entry2_group1_password",
		Url: "entry2_group1_url", Notes: "", Path: []string{"Root", "group 1"},
	})

	assert.Contains(t, entriesForGroup2, keepass.SecretEntry{
		Group: "Root|group 2", Title: "entry_in_group2",
		Username: "user_in_group2", Password: "password_in_group2",
		Url: "https://group2.com", Notes: "", Path: []string{"Root", "group 2"},
	})
}

func Test_ReadEntriesFromContentGroupedByPath_Broken_File(t *testing.T) {

	bytesContent, err := os.ReadFile("testdata/files/db_broken.kdbx")

	if err != nil {
		t.Fatal("Could not find test DB")
	}

	cipheredKeepassDB := keepass.CipheredKeepassDB{DBBytes: bytesContent, Password: "keepassui"}

	secretsDB, err := cipheredKeepassDB.ReadEntriesFromContentGroupedByPath()

	if err == nil {
		t.Fatal("We expect an error in this test because the DB file is broken")
	}
	if secretsDB.EntriesByPath != nil || secretsDB.PathsInOrder != nil {
		t.Fatal("entriesGroupedByPath and pathsInOrder should be nil")
	}

	assert.EqualError(t, err, "Failed to verify HMAC for block 0")
}

func Test_writeDBBytes(t *testing.T) {
	secretsDB := secretsDBForTesting()

	bytes, err := secretsDB.WriteDBBytes("master")

	if err != nil {
		t.Error("Should not error")
	}

	cipheredKeepassDB := keepass.CipheredKeepassDB{DBBytes: bytes, Password: "master"}

	secretsDBReadFromNewDBBytes, err := cipheredKeepassDB.ReadEntriesFromContentGroupedByPath()

	if err != nil {
		t.Error("Should be able to read secrets")
	}

	assert.Equal(t, secretsDB, secretsDBReadFromNewDBBytes)
}

func Test_AddSecretEntryInANewPath(t *testing.T) {
	secretsDB := keepass.SecretsDB{PathsInOrder: []string{}, EntriesByPath: make(map[string][]keepass.SecretEntry)}

	secretsDB.AddSecretEntry(keepass.SecretEntry{
		Group: "Root", Title: "G1",
		Notes: "Notes are important", Path: []string{"Root"}, IsGroup: true,
	})

	assert.Equal(t, 1, len(secretsDB.PathsInOrder))
	assert.Equal(t, []string{"Root"}, secretsDB.PathsInOrder)

	assert.Equal(t, 1, len(secretsDB.EntriesByPath))

	entriesForRoot, ok := secretsDB.EntriesByPath["Root"]

	assert.Equal(t, 1, len(entriesForRoot))

	assert.True(t, ok)

	assert.Equal(t, "Root", entriesForRoot[0].Group)
	assert.Equal(t, "G1", entriesForRoot[0].Title)
	assert.Equal(t, "Notes are important", entriesForRoot[0].Notes)
	assert.Equal(t, "", entriesForRoot[0].Username)
	assert.Equal(t, "", entriesForRoot[0].Password)
	assert.Equal(t, "", entriesForRoot[0].Url)
	assert.Equal(t, []string{"Root"}, entriesForRoot[0].Path)
	assert.True(t, entriesForRoot[0].IsGroup)
}

func Test_AddSecretEntryInExistingPath(t *testing.T) {
	secretsDB := secretsDBForTesting()

	secretsDB.AddSecretEntry(keepass.SecretEntry{
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

	secretsDB.AddSecretEntry(keepass.SecretEntry{
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

	deleted := secretsDB.DeleteSecretEntry(keepass.SecretEntry{
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

	deleted := secretsDB.DeleteSecretEntry(keepass.SecretEntry{
		Group: "Root|G1", Title: "entry_in_RG1_2",
		Username: "user_in_RG1_2", Password: "password_in_RG1_2",
		Url: "https://RG1_2.com", Notes: "", Path: []string{"Root", "G1"},
	})

	assert.True(t, deleted, "Should delete 1 entry")

	entriesForRootG1 = secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 2, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "G2", entriesForRootG1[1].Title)
}

func TestDeleteSecretEntry_ShouldDeleteGroupAndItsContents(t *testing.T) {
	secretsDB := secretsDBForTesting()

	entriesForRootG1 := secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 3, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "entry_in_RG1_2", entriesForRootG1[1].Title)
	assert.Equal(t, "G2", entriesForRootG1[2].Title)

	entriesForRootG1G2, ok := secretsDB.EntriesByPath["Root|G1|G2"]

	assert.True(t, ok, "G2 should have entries")

	assert.Equal(t, 1, len(entriesForRootG1G2))

	deleted := secretsDB.DeleteSecretEntry(keepass.SecretEntry{
		Group: "Root|G1", Title: "G2",
		Notes: "", Path: []string{"Root", "G1"}, IsGroup: true,
	})

	assert.True(t, deleted, "Should delete the group and its contents")

	entriesForRootG1 = secretsDB.EntriesByPath["Root|G1"]

	assert.Equal(t, 2, len(entriesForRootG1))

	assert.Equal(t, "entry_in_RG1", entriesForRootG1[0].Title)
	assert.Equal(t, "entry_in_RG1_2", entriesForRootG1[1].Title)

	_, ok = secretsDB.EntriesByPath["Root|G1|G2"]

	assert.False(t, ok, "G2 should not have entries")
}

func secretsDBForTesting() keepass.SecretsDB {
	entriesByPath := make(map[string][]keepass.SecretEntry)

	entriesByPath["Root"] = []keepass.SecretEntry{
		{
			Group: "Root", Title: "entry_in_RG1",
			Username: "user_in_root", Password: "password_in_root",
			Url: "https://rootEntry.com", Notes: "", Path: []string{"Root"},
		},
		{
			Group: "Root", Title: "G1",
			Notes: "", Path: []string{"Root"}, IsGroup: true,
		},
	}

	entriesByPath["Root|G1"] = []keepass.SecretEntry{
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

	entriesByPath["Root|G1|G2"] = []keepass.SecretEntry{
		{
			Group: "Root|G1|G2", Title: "entry_in_RG2",
			Username: "user_in_RG2", Password: "password_in_RG2",
			Url: "https://RG1G2.com", Notes: "", Path: []string{"Root", "G1", "G2"},
		},
	}

	return keepass.SecretsDB{PathsInOrder: []string{"Root", "Root|G1", "Root|G1|G2"}, EntriesByPath: entriesByPath}
}
