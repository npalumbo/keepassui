package ui_test

import (
	"keepassui/internal/keepass"
	"keepassui/internal/ui"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

func TestCreateSecretFormReadOnlyFormShowsDisabledEntries(t *testing.T) {
	secretForm := ui.CreateSecretForm(true)

	w := test.NewWindow(secretForm.FormContainer)

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "createSecretForm_ReadOnly_Shows_Disabled_Entries.png", w.Canvas().Capture())
}

func TestCreateSecretFormNotReadOnlyFormShowsEnabledEntries(t *testing.T) {
	secretForm := ui.CreateSecretForm(false)

	w := test.NewWindow(secretForm.FormContainer)

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "createSecretForm_Shows_Enabled_Entries.png", w.Canvas().Capture())
}

func TestUpdateForm_Shows_Populated_Form_Fields(t *testing.T) {
	secretForm := ui.CreateSecretForm(false)

	secretForm.TypeSecretEntryInForm(keepass.SecretEntry{
		Title:    "aTitle",
		Username: "aUsername",
		Password: "aPassword",
		Url:      "aUrl",
		Notes:    "someNotes",
	})

	secretForm.FormContainer.Refresh()
	w := test.NewWindow(secretForm.FormContainer)

	w.Resize(fyne.Size{Width: 300, Height: 300})
	test.AssertImageMatches(t, "updateForm_Shows_Populated_fields.png", w.Canvas().Capture())
}

func TestUpdateEntry_Populates_SecretEntry(t *testing.T) {
	secretForm := ui.CreateSecretForm(false)

	secretForm.TypeSecretEntryInForm(keepass.SecretEntry{
		Title:    "aTitle",
		Username: "aUsername",
		Password: "aPassword",
		Url:      "aUrl",
		Notes:    "someNotes",
	})

	secretEntry := keepass.SecretEntry{}

	secretForm.UpdateEntry(&secretEntry)

	assert.Equal(t, "aTitle", secretEntry.Title)
	assert.Equal(t, "aUsername", secretEntry.Username)
	assert.Equal(t, "aPassword", secretEntry.Password)
	assert.Equal(t, "aUrl", secretEntry.Url)
	assert.Equal(t, "someNotes", secretEntry.Notes)
}
