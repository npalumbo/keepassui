package ui_test

import (
	mock_secretsreader "keepassui/internal/mocks/secretsreader"
	mocks_ui "keepassui/internal/mocks/ui"
	"keepassui/internal/secretsdb"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"go.uber.org/mock/gomock"
)

func TestAddEntryShowsDisabledConfirmButtonWhenNotFullyPopulated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	addEntryView := ui.CreateAddEntryView(nil, "prevousScreen", stagerController, w)

	templateSecret := secretsdb.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1).Return(nil)

	addEntryView.AddEntry(&templateSecret)

	w.SetContent(addEntryView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})
	addEntryView.GetPaintedContainer().Refresh()

	test.AssertImageMatches(t, "AddEntry_Confirm_Disabled_Confirm_When_Not_Fully_Populated.png", w.Canvas().Capture())
}

func TestAddEntryTapOnCancelTakesUsToThePreviousScreenWithoutChangingSecretsDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	mockSecretsReader := mock_secretsreader.NewMockSecretReader(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	addEntryView := ui.CreateAddEntryView(mockSecretsReader, "previousScreen", stagerController, w)

	templateSecret := secretsdb.SecretEntry{Path: []string{"path 1"}, Group: "path 1", IsGroup: false}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1).Return(nil)
	mockSecretsReader.EXPECT().AddSecretEntry(&templateSecret).Times(0)

	addEntryView.AddEntry(&templateSecret)

	stagerController.EXPECT().TakeOver("previousScreen").Times(1)

	// Ideally we would use test.Tap() on the Cancel button but the button is not reachable from widget.Form
	addEntryView.SecretForm.DetailsForm.OnCancel()
}

func TestAddEntryShowsEnabledConfirmButtonWhenFullyPopulated(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	addEntryView := ui.CreateAddEntryView(nil, "prevousScreen", stagerController, w)

	templateSecret := secretsdb.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
	}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1).Return(nil)

	addEntryView.AddEntry(&templateSecret)

	addEntryView.SecretForm.TypeSecretEntryInForm(secretsdb.SecretEntry{
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes"},
	)

	w.SetContent(addEntryView.GetPaintedContainer())
	w.Resize(fyne.Size{Width: 600, Height: 600})
	addEntryView.GetPaintedContainer().Refresh()

	test.AssertImageMatches(t, "AddEntry_Confirm_Enabled_When_Fully_Populated.png", w.Canvas().Capture())
}

func TestAddEntryTapOnSubmitTakesUsToThePreviousScreenAndAddsEntryToSecretsDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	mockSecretsReader := mock_secretsreader.NewMockSecretReader(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	addEntryView := ui.CreateAddEntryView(mockSecretsReader, "previousScreen", stagerController, w)

	templateSecret := secretsdb.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
	}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1).Return(nil)
	mockSecretsReader.EXPECT().Save().Times(1)

	addEntryView.AddEntry(&templateSecret)

	addEntryView.SecretForm.TypeSecretEntryInForm(secretsdb.SecretEntry{
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes"},
	)

	mockSecretsReader.EXPECT().AddSecretEntry(secretsdb.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes"}).Times(1)

	stagerController.EXPECT().TakeOver("previousScreen").Times(1).Return(nil)

	// Ideally we would use test.Tap() on the Confirm button but the button is not reachable from widget.Form
	addEntryView.SecretForm.DetailsForm.OnSubmit()
}

func TestModifyEntryTapOnSubmitTakesUsToThePreviousScreenAndModifiesEntryToSecretsDB(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stagerController := mocks_ui.NewMockStagerController(mockCtrl)
	mockSecretsReader := mock_secretsreader.NewMockSecretReader(mockCtrl)
	w := test.NewWindow(container.NewWithoutLayout())
	modifyEntryView := ui.CreateAddEntryView(mockSecretsReader, "previousScreen", stagerController, w)

	templateSecret := secretsdb.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
		Title: "aTitle", Username: "aUsername", Password: "aPassword", Url: "aUrl", Notes: "someNotes",
	}
	stagerController.EXPECT().TakeOver("AddEntry").Times(1).Return(nil)
	mockSecretsReader.EXPECT().Save().Times(1)

	modifyEntryView.ModifyEntry(&templateSecret)

	modifyEntryView.SecretForm.TypeSecretEntryInForm(secretsdb.SecretEntry{
		Title: "aModifiedTitle", Username: "aModifiedUsername", Password: "aModifiedPassword", Url: "aModifiedUrl", Notes: "someModifiedNotes"},
	)

	mockSecretsReader.EXPECT().ModifySecretEntry("aTitle", "path 1", false, secretsdb.SecretEntry{
		Path: []string{"path 1"}, Group: "path 1", IsGroup: false,
		Title: "aModifiedTitle", Username: "aModifiedUsername", Password: "aModifiedPassword", Url: "aModifiedUrl", Notes: "someModifiedNotes"}).Times(1)

	stagerController.EXPECT().TakeOver("previousScreen").Times(1).Return(nil)

	// Ideally we would use test.Tap() on the Confirm button but the button is not reachable from widget.Form
	modifyEntryView.SecretForm.DetailsForm.OnSubmit()
}
