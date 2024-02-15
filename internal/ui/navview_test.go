package ui

import (
	"errors"
	"keepassui/internal/keepass"
	mock_keepass "keepassui/internal/mocks/keepass"
	mock_addentryview "keepassui/internal/mocks/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"go.uber.org/mock/gomock"
)

type MockedSecretReaderFactory struct {
	mockedSecretReader keepass.SecretReader
}

func (m MockedSecretReaderFactory) GetSecretReader(d DBPathAndPassword) keepass.SecretReader {
	return m.mockedSecretReader
}

func TestNavView_DataChanged_Does_Nothing_When_DBPathAndPassword_is_EmptyObject(t *testing.T) {
	dbPathAndPassword := &DBPathAndPassword{}
	w := test.NewWindow(container.NewWithoutLayout())

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, nil)

	navView.DataChanged()

	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Does_Nothing_When_DBPathAndPassword_is_EmptyObject.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_Shows_Error_Error_Reading_secrets(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(keepass.SecretsDB{}, errors.New("Fake Error"))

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Reading_Secrets.png", w.Canvas().Capture())
}

func TestNavView_DataChanged(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBForTesting(),
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_one_group.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_two_groups(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups(),
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())
}

func TestNavView_NavigateToNestedFolder(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups(),
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	// Ideally we would simulate a click from the UI but I struggle to find the right open button from the list
	navView.UpdateNavView("path 2")

	test.AssertImageMatches(t, "navView_two_groups_nested_group.png", w.Canvas().Capture())
}

func TestNavView_DeleteFirstEntry(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)
	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	// Ideally we would simulate a click from the UI but I struggle to find the right open button from the list
	secretsDBWithTwoGroups.DeleteSecretEntry(keepass.SecretEntry{
		Title: "title 2", Group: "path 2", Username: "username 2",
		Password: "password 2", Url: "url 2", Notes: "notes 2"})

	navView.UpdateNavView("path 2")

	test.AssertImageMatches(t, "navView_two_groups_nested_group_with_one_entry_deleted.png", w.Canvas().Capture())
}

func TestNavView_TapSaveButtonOpensSaveDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.saveButton)

	test.AssertImageMatches(t, "navView_two_groups_tap_save_button.png", w.Canvas().Capture())
}

func TestNavView_TapOnNewGroupOpensNewGroupDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	navView := CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.groupCreateButton)

	test.AssertImageMatches(t, "navView_two_groups_tap_new_group_button.png", w.Canvas().Capture())
}

func TestNavView_TapOnNewSecretCallsAddEntry(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_keepass.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	entryUpdater := mock_addentryview.NewMockEntryUpdater(mockCtrl)

	templateSecret := keepass.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}

	entryUpdater.EXPECT().AddEntry(&templateSecret, &secretsDBWithTwoGroups).Times(1)

	navView := CreateNavView(dbPathAndPassword, entryUpdater, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.navAndListContainer)
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.secretEntryCreateButton)
}

func secretsDBForTesting() keepass.SecretsDB {
	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}}
	paths := []string{"path 1"}
	secretsDB := keepass.SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  paths,
	}
	return secretsDB
}

func secretsDBWithTwoGroups() keepass.SecretsDB {
	secretsGroupedByPath := make(map[string][]keepass.SecretEntry)
	secretsGroupedByPath["path 1"] = []keepass.SecretEntry{{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"},
		{Title: "path 2", Group: "path 1", Notes: "", IsGroup: true}}
	secretsGroupedByPath["path 2"] = []keepass.SecretEntry{
		{Title: "title 2", Group: "path 2", Username: "username 2", Password: "password 2", Url: "url 2", Notes: "notes 2"},
		{Title: "title 3", Group: "path 2", Username: "username 3", Password: "password 3", Url: "url 3", Notes: "notes 3"},
	}
	paths := []string{"path 1", "path 2"}
	secretsDB := keepass.SecretsDB{
		EntriesByPath: secretsGroupedByPath,
		PathsInOrder:  paths,
	}
	return secretsDB
}
