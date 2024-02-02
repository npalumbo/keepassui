package ui

import (
	"errors"
	"keepassui/internal/keepass"
	mock_keepass "keepassui/internal/mocks/keepass"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"go.uber.org/mock/gomock"
)

func TestNavView_DataChanged_Does_Nothing_When_DBPathAndPassword_is_EmptyObject(t *testing.T) {
	dbPathAndPassword := &DBPathAndPassword{}
	w := test.NewWindow(container.NewWithoutLayout())
	w.Resize(fyne.NewSize(600, 600))

	navView := CreateNavView(dbPathAndPassword, nil, w, nil)

	navView.DataChanged()

	test.AssertImageMatches(t, "navView_Err_Does_Nothing_When_DBPathAndPassword_is_EmptyObject.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_Shows_Error_Error_Reading_secrets(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	w.Resize(fyne.NewSize(600, 600))

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(keepass.SecretsDB{}, errors.New("Fake Error"))

	navView := CreateNavView(dbPathAndPassword, nil, w, func(d DBPathAndPassword) keepass.SecretReader {
		return secretReader
	})

	navView.DataChanged()

	test.AssertImageMatches(t, "navView_Err_Reading_Secrets.png", w.Canvas().Capture())
}

func TestNavView_DataChanged(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}}
	paths := []string{"path 1"}
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		keepass.SecretsDB{
			EntriesByPath: secretsGroupedByPath,
			PathsInOrder:  paths,
		},
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, w, func(d DBPathAndPassword) keepass.SecretReader {
		return secretReader
	})

	w.SetContent(navView.fullContainer)
	w.Resize(fyne.NewSize(600, 600))

	navView.DataChanged()
	navView.fullContainer.Refresh()
	test.AssertImageMatches(t, "navView_one_group.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_two_groups(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"},
		{Title: "path 2", Group: "path 1", Notes: "", IsGroup: true}}
	secretsGroupedByPath["path 2"] = []keepass.SecretEntry{
		{Title: "title 2", Group: "path 2", Username: "username 2", Password: "password 2", Url: "url 2", Notes: "notes 2"},
		{Title: "title 3", Group: "path 2", Username: "username 3", Password: "password 3", Url: "url 3", Notes: "notes 3"},
	}
	paths := []string{"path 1", "path 2"}
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		keepass.SecretsDB{
			EntriesByPath: secretsGroupedByPath,
			PathsInOrder:  paths,
		},
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, w, func(d DBPathAndPassword) keepass.SecretReader {
		return secretReader
	})
	w.SetContent(navView.fullContainer)
	w.Resize(fyne.NewSize(600, 600))

	navView.DataChanged()
	navView.fullContainer.Refresh()
	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())
}

// TODO add test navigating to a nested folder
