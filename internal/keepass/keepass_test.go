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

	assert.Equal(t, 2, len(entriesGroupedByPath))

	assert.Equal(t, []string{"group 1", "group 2"}, pathsInOrder)
	keysFromEntriesGroupedByPath := make([]string, 0, len(entriesGroupedByPath))
	for k := range entriesGroupedByPath {
		keysFromEntriesGroupedByPath = append(keysFromEntriesGroupedByPath, k)
	}

	assert.Contains(t, keysFromEntriesGroupedByPath, "group 1")
	assert.Contains(t, keysFromEntriesGroupedByPath, "group 2")

	entriesForGroup1 := entriesGroupedByPath["group 1"]
	entriesForGroup2 := entriesGroupedByPath["group 2"]

	assert.Equal(t, 2, len(entriesForGroup1))
	assert.Equal(t, 1, len(entriesForGroup2))

	assert.Contains(t, entriesForGroup1, keepass.SecretEntry{
		Path: "group 1", Title: "entry_inside_group1",
		Username: "user_in_group1", Password: "password_in_group_1",
		Url: "https://ingroup1.com/", Notes: "",
	})

	assert.Contains(t, entriesForGroup1, keepass.SecretEntry{
		Path: "group 1", Title: "entry_2_in_group_1",
		Username: "entry2_group1_username", Password: "entry2_group1_password",
		Url: "entry2_group1_url", Notes: "",
	})

	assert.Contains(t, entriesForGroup2, keepass.SecretEntry{
		Path: "group 2", Title: "entry_in_group2",
		Username: "user_in_group2", Password: "password_in_group2",
		Url: "https://group2.com", Notes: "",
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
