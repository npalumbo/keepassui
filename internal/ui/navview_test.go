package ui_test

import (
	"errors"
	mock_secretsreader "keepassui/internal/mocks/secretsreader"
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/secretsdb"
	"keepassui/internal/secretsreader"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"go.uber.org/mock/gomock"
)

type MockedSecretReaderFactory struct {
	mockedSecretReader secretsreader.SecretReader
}

func (m MockedSecretReaderFactory) GetSecretReader(d ui.DBPathAndPassword) secretsreader.SecretReader {
	return m.mockedSecretReader
}

func TestNavView_DataChanged_Does_Nothing_When_DBPathAndPassword_is_EmptyObject(t *testing.T) {
	dbPathAndPassword := &ui.DBPathAndPassword{}
	w := test.NewWindow(container.NewWithoutLayout())

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Does_Nothing_When_DBPathAndPassword_is_EmptyObject.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_Shows_Error_Error_Reading_secrets(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(secretsdb.SecretsDB{}, errors.New("Fake Error"))

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Reading_Secrets.png", w.Canvas().Capture())
}

func TestNavView_DataChanged(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBForTesting(),
		nil,
	)

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_one_group.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_two_groups(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups(),
		nil,
	)

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())
}

func TestNavView_NavigateToNestedFolder(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups(),
		nil,
	)

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
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

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)
	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	// Ideally we would simulate a click from the UI but I struggle to find the right open button from the list
	secretsDBWithTwoGroups.DeleteSecretEntry(secretsdb.SecretEntry{
		Title: "title 2", Group: "path 2", Username: "username 2",
		Password: "password 2", Url: "url 2", Notes: "notes 2"})

	navView.UpdateNavView("path 2")

	test.AssertImageMatches(t, "navView_two_groups_nested_group_with_one_entry_deleted.png", w.Canvas().Capture())
}

func TestNavView_TapSaveButtonOpensSaveDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.SaveButton)

	test.AssertImageMatches(t, "navView_two_groups_tap_save_button.png", w.Canvas().Capture())
}

func TestNavView_TapOnNewGroupOpensNewGroupDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	navView := ui.CreateNavView(dbPathAndPassword, nil, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.GroupCreateButton)

	test.AssertImageMatches(t, "navView_two_groups_tap_new_group_button.png", w.Canvas().Capture())
}

func TestNavView_TapOnNewSecretCallsAddEntry(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dbPathAndPassword := &ui.DBPathAndPassword{UriID: "file://path", Password: "password", ContentInBytes: []byte{}}

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(
		secretsDBWithTwoGroups,
		nil,
	)

	entryUpdater := mocks_ui.NewMockEntryUpdater(mockCtrl)

	templateSecret := secretsdb.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}

	entryUpdater.EXPECT().AddEntry(&templateSecret, &secretsDBWithTwoGroups).Times(1)

	navView := ui.CreateNavView(dbPathAndPassword, entryUpdater, nil, w, nil, MockedSecretReaderFactory{mockedSecretReader: secretReader})

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.SecretEntryCreateButton)
}
