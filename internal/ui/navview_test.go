package ui_test

import (
	"errors"
	mock_secretsreader "keepassui/internal/mocks/secretsreader"
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/secretsdb"
	"keepassui/internal/secretsreader"
	"keepassui/internal/ui"
	"slices"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"go.uber.org/mock/gomock"
)

type MockedSecretReaderFactory struct {
	mockedSecretReader secretsreader.SecretReader
}

func (m MockedSecretReaderFactory) GetSecretReader(d secretsreader.DefaultSecretsReader) secretsreader.SecretReader {
	return m.mockedSecretReader
}

func TestNavView_DataChanged_Does_Nothing_When_SecretsReader_is_EmptyObject(t *testing.T) {
	secretsReader := &secretsreader.DefaultSecretsReader{}
	w := test.NewWindow(container.NewWithoutLayout())

	navView := ui.CreateNavView(secretsReader, nil, w, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Does_Nothing_When_SecretsReader_is_EmptyObject.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_Shows_Error_Error_Reading_secrets(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(errors.New("Fake Error"))
	secretReader.EXPECT().GetUriID().Times(1).Return("path")

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_Err_Reading_Secrets.png", w.Canvas().Capture())
}

func TestNavView_DataChanged(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBForTesting().EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_one_group.png", w.Canvas().Capture())
}

func TestNavView_DataChanged_two_groups(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups().EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))
	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())
}

func TestNavView_NavigateToNestedFolder(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups().EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	secretReader.EXPECT().GetEntriesForPath("path 1|path 2").Times(1).Return(secretsDBWithTwoGroups().EntriesByPath["path 1|path 2"])

	// Ideally we would simulate a click from the UI but I struggle to find the right open button from the list
	navView.UpdateNavView("path 1|path 2")

	test.AssertImageMatches(t, "navView_two_groups_nested_group.png", w.Canvas().Capture())
}

func TestNavView_NavigateToNestedFolderByTappingOnListItem(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups().EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	secretReader.EXPECT().GetEntriesForPath("path 1|path 2").Times(1).Return(secretsDBWithTwoGroups().EntriesByPath["path 1|path 2"])

	objects := test.LaidOutObjects(navView.ListPanel.Objects[0])

	// Since widget.listItem is unexported we can't use it as filter in the objects array.
	tappable, ok := objects[40].(fyne.Tappable)
	if ok {
		test.Tap(tappable)
	} else {
		t.FailNow()
	}

	time.Sleep(10 * time.Millisecond)

	test.AssertImageMatches(t, "navView_two_groups_nested_group.png", w.Canvas().Capture())
}

func TestNavView_DeleteFirstEntry(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()
	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	// Ideally we would simulate a click from the UI but I struggle to find the right open button from the list
	secretsDBWithTwoGroups.DeleteSecretEntry(secretsdb.SecretEntry{
		Title: "title 2", Group: "path 1|path 2", Username: "username 2",
		Password: "password 2", Url: "url 2", Notes: "notes 2"})

	secretReader.EXPECT().GetEntriesForPath("path 1|path 2").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1|path 2"])

	navView.UpdateNavView("path 1|path 2")

	test.AssertImageMatches(t, "navView_two_groups_nested_group_with_one_entry_deleted.png", w.Canvas().Capture())
}

func TestNavView_TapOnNewGroupOpensNewGroupDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("file://path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

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

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("file://path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	entryUpdater := mocks_ui.NewMockEntryUpdater(mockCtrl)

	templateSecret := secretsdb.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}

	entryUpdater.EXPECT().AddEntry(&templateSecret).Times(1)

	navView := ui.CreateNavView(secretReader, entryUpdater, w, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	test.Tap(navView.SecretEntryCreateButton)
}

func TestNavView_TapOnEditSecretCallsModifyEntry(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("file://path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	entryUpdater := mocks_ui.NewMockEntryUpdater(mockCtrl)

	templateSecret := secretsdb.SecretEntry{Title: "title", Group: "path 1", Username: "username", Password: "password", Url: "url", Notes: "notes"}

	entryUpdater.EXPECT().ModifyEntry(&templateSecret).Times(1)

	navView := ui.CreateNavView(secretReader, entryUpdater, w, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	listPanel := navView.GetPaintedContainer().Objects[0].(*fyne.Container)

	list := listPanel.Objects[0].(*widget.List)

	componentsInsideList := test.LaidOutObjects(list)
	firstModifyButtonIdx := slices.IndexFunc(componentsInsideList, func(obj fyne.CanvasObject) bool {
		button, ok := obj.(*widget.Button)
		return ok && button.Icon == theme.DocumentCreateIcon()
	})

	firstModifyButton := componentsInsideList[firstModifyButtonIdx].(*widget.Button)

	firstModifyButton.OnTapped()
}

func TestNavView_TapOnEditGroupOpensNewGroupDialog(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("file://path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	navView := ui.CreateNavView(secretReader, nil, w, nil)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	listPanel := navView.GetPaintedContainer().Objects[0].(*fyne.Container)

	list := listPanel.Objects[0].(*widget.List)

	componentsInsideList := test.LaidOutObjects(list)
	count := 0
	secondModifyButtonIdx := slices.IndexFunc(componentsInsideList, func(obj fyne.CanvasObject) bool {
		button, ok := obj.(*widget.Button)
		if ok && button.Icon == theme.DocumentCreateIcon() {
			count = count + 1
			return count == 2
		}
		return false
	})

	secondModifyButton := componentsInsideList[secondModifyButtonIdx].(*widget.Button)

	test.Tap(secondModifyButton)

	test.AssertImageMatches(t, "navView_two_groups_tap_edit_group_button.png", w.Canvas().Capture())
}

func TestNavView_TapLockDBButtonCallsHome(t *testing.T) {
	w := test.NewWindow(container.NewWithoutLayout())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	secretReader := mock_secretsreader.NewMockSecretReader(mockCtrl)

	secretsDBWithTwoGroups := secretsDBWithTwoGroups()
	secretReader.EXPECT().ReadEntriesFromContentGroupedByPath().Times(1).Return(nil)
	secretReader.EXPECT().GetUriID().Times(1).Return("file://path")
	secretReader.EXPECT().GetFirstPath().Times(1).Return("path 1")
	secretReader.EXPECT().GetEntriesForPath("path 1").Times(1).Return(secretsDBWithTwoGroups.EntriesByPath["path 1"])

	mock_stagercontroller := mocks_ui.NewMockStagerController(mockCtrl)

	navView := ui.CreateNavView(secretReader, nil, w, mock_stagercontroller)

	mock_stagercontroller.EXPECT().TakeOver("NavView").Times(1)

	navView.DataChanged()

	w.SetContent(navView.GetPaintedContainer())
	w.Resize(fyne.NewSize(600, 600))

	test.AssertImageMatches(t, "navView_two_groups.png", w.Canvas().Capture())

	mock_stagercontroller.EXPECT().TakeOver("Home").Times(1)

	test.Tap(navView.LockDBButton)
}
