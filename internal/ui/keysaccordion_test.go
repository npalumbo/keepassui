package ui

import (
	"errors"
	"keepassui/internal/keepass"
	mock_keepass "keepassui/internal/mocks/keepass"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"go.uber.org/mock/gomock"
)

func TestKeysAccordion_DataChanged_Shows_Error_When_type_not_DBPathAndPassword_dbPathAndPassword(t *testing.T) {
	dbPathAndPassword := binding.NewUntyped()
	w := test.NewWindow(container.NewWithoutLayout())
	w.Resize(fyne.NewSize(600, 600))
	err := dbPathAndPassword.Set("fake string not DBPathAndPassword")
	if err != nil {
		t.Fail()
	}
	keyAccordion := CreatekeyAccordion(dbPathAndPassword, nil, nil, w, nil)

	keyAccordion.DataChanged()

	test.AssertImageMatches(t, "keysAccordion_Err_casting_dbPathAndPassword_to_DBPathAndPassword.png", w.Canvas().Capture())
}

func TestKeysAccordion_DataChanged_Shows_Error_Error_Reading_secrets(t *testing.T) {
	dbPathAndPassword := binding.NewUntyped()
	content := binding.NewBytes()
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	w.Resize(fyne.NewSize(600, 600))

	err := dbPathAndPassword.Set(DBPathAndPassword{Path: "path", Password: "password"})
	if err != nil {
		t.Fail()
	}
	err = content.Set(make([]byte, 0))
	if err != nil {
		t.Fail()
	}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil, nil, errors.New("Fake Error"))

	keyAccordion := CreatekeyAccordion(dbPathAndPassword, content, nil, w, func(contentInBytes []byte, password string) keepass.SecretReader {
		return secretReader
	})

	keyAccordion.DataChanged()

	test.AssertImageMatches(t, "keysAccordion_Err_Reading_Secrets.png", w.Canvas().Capture())
}

func TestKeysAccordion_DataChanged(t *testing.T) {
	dbPathAndPassword := binding.NewUntyped()
	content := binding.NewBytes()
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := dbPathAndPassword.Set(DBPathAndPassword{Path: "path", Password: "password"})
	if err != nil {
		t.Fail()
	}
	err = content.Set(make([]byte, 0))
	if err != nil {
		t.Fail()
	}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Path: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}}
	paths := []string{"path 1"}
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(secretsGroupedByPath, paths, nil)

	keyAccordion := CreatekeyAccordion(dbPathAndPassword, content, nil, w, func(contentInBytes []byte, password string) keepass.SecretReader {
		return secretReader
	})
	w.SetContent(container.NewStack(keyAccordion.accordionWidget))
	w.Resize(fyne.NewSize(600, 600))

	keyAccordion.DataChanged()
	keyAccordion.accordionWidget.Open(0)
	test.AssertImageMatches(t, "keysAccordion_one_group.png", w.Canvas().Capture())
}

func TestKeysAccordion_DataChanged_two_groups(t *testing.T) {
	dbPathAndPassword := binding.NewUntyped()
	content := binding.NewBytes()
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := dbPathAndPassword.Set(DBPathAndPassword{Path: "path", Password: "password"})
	if err != nil {
		t.Fail()
	}
	err = content.Set(make([]byte, 0))
	if err != nil {
		t.Fail()
	}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Path: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}}
	secretsGroupedByPath["path 2"] = []keepass.SecretEntry{
		{Title: "title 2", Path: "path 2", Username: "username 2", Password: "password 2", Url: "url 2", Notes: "notes 2"},
		{Title: "title 3", Path: "path 2", Username: "username 3", Password: "password 3", Url: "url 3", Notes: "notes 3"},
	}
	paths := []string{"path 1", "path 2"}
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(secretsGroupedByPath, paths, nil)

	keyAccordion := CreatekeyAccordion(dbPathAndPassword, content, nil, w, func(contentInBytes []byte, password string) keepass.SecretReader {
		return secretReader
	})
	w.SetContent(container.NewStack(keyAccordion.accordionWidget))
	w.Resize(fyne.NewSize(600, 600))

	keyAccordion.DataChanged()
	keyAccordion.accordionWidget.Open(1)
	test.AssertImageMatches(t, "keysAccordion_two_groups.png", w.Canvas().Capture())
}
