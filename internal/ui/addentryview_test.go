package ui_test

import (
	"keepassui/internal/keepass"
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/ui"
	"slices"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddEntryShowsDisabledConfirmButtonWhenNotFullyPopulated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	addEntryView := ui.CreateAddEntryView("prevousScreen", stagerController)

	templateSecret := keepass.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1)

	secretsDB := secretsDBForTesting()

	addEntryView.AddEntry(&templateSecret, &secretsDB)

	w := test.NewWindow(container.NewWithoutLayout())

	w.SetContent(addEntryView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})
	addEntryView.GetPaintedContainer().Refresh()

	test.AssertImageMatches(t, "AddEntry_Confirm_Disabled_Confirm_When_Not_Fully_Populated.png", w.Canvas().Capture())
}

func TestAddEntryTapOnCancelTakesUsToThePreviousScreenWithoutChangingSecretsDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	addEntryView := ui.CreateAddEntryView("previousScreen", stagerController)

	templateSecret := keepass.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1)

	secretsDB := secretsDBForTesting()

	addEntryView.AddEntry(&templateSecret, &secretsDB)

	stagerController.EXPECT().TakeOver("previousScreen").Times(1)

	// Ideally we would use test.Tap() on the Cancel button but the button is not reachable from widget.Form
	addEntryView.SecretForm.DetailsForm.OnCancel()

	assert.Equal(t, secretsDBForTesting(), secretsDB)
}

func TestAddEntryShowsEnabledConfirmButtonWhenFullyPopulated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	addEntryView := ui.CreateAddEntryView("prevousScreen", stagerController)

	templateSecret := keepass.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
	}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1)

	secretsDB := secretsDBForTesting()

	addEntryView.AddEntry(&templateSecret, &secretsDB)

	addEntryView.SecretForm.TypeSecretEntryInForm(keepass.SecretEntry{
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes"},
	)

	w := test.NewWindow(container.NewWithoutLayout())

	w.SetContent(addEntryView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})
	addEntryView.GetPaintedContainer().Refresh()

	test.AssertImageMatches(t, "AddEntry_Confirm_Enabled_When_Fully_Populated.png", w.Canvas().Capture())
}

func TestAddEntryTapOnSubmitTakesUsToThePreviousScreenAndAddsEntryToSecretsDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	addEntryView := ui.CreateAddEntryView("previousScreen", stagerController)

	templateSecret := keepass.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
	}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1)

	secretsDB := secretsDBForTesting()
	entriesForPath1 := secretsDB.EntriesByPath["path 1"]

	addEntryView.AddEntry(&templateSecret, &secretsDB)
	assert.Equal(t, 1, len(entriesForPath1), "Starts with one entry")

	matchingEntryPredicate := func(entry keepass.SecretEntry) bool {
		return entry.Title == "aTitle" &&
			entry.Username == "aUsername" && entry.Password == "aPassword" &&
			entry.Url == "aUrl" && entry.Notes == "someNotes"
	}

	indexEntryWithNameATitle := slices.IndexFunc(entriesForPath1, matchingEntryPredicate)

	assert.Equal(t, -1, indexEntryWithNameATitle, "Not found any entry with title: aTitle, etc")

	addEntryView.SecretForm.TypeSecretEntryInForm(keepass.SecretEntry{
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes"},
	)

	stagerController.EXPECT().TakeOver("previousScreen").Times(1)

	// Ideally we would use test.Tap() on the Confirm button but the button is not reachable from widget.Form
	addEntryView.SecretForm.DetailsForm.OnSubmit()

	assert.NotEqual(t, secretsDBForTesting(), secretsDB, "Is not equals to the original secretsDB")
	entriesForPath1 = secretsDB.EntriesByPath["path 1"]

	assert.Equal(t, 2, len(entriesForPath1), "It has now 2 entries")

	indexEntryWithNameATitle = slices.IndexFunc(entriesForPath1, matchingEntryPredicate)

	assert.Equal(t, 1, indexEntryWithNameATitle, "An entry with title: aTitle, etc in position 1")
}
