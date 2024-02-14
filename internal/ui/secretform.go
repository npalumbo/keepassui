package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"keepassui/internal/keepass"
	keepassuiwidget "keepassui/internal/widget"
)

type SecretForm struct {
	titleEntry    *widget.Entry
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	urlEntry      *widget.Entry
	notesEntry    *widget.Entry
	detailsForm   *widget.Form
	FormContainer *fyne.Container
}

func (f *SecretForm) UpdateEntry(entry *keepass.SecretEntry) {
	entry.Title = f.titleEntry.Text
	entry.Username = f.usernameEntry.Text
	entry.Password = f.passwordEntry.Text
	entry.Url = f.urlEntry.Text
	entry.Notes = f.notesEntry.Text
}

func (f *SecretForm) UpdateForm(entry keepass.SecretEntry) {
	f.titleEntry.Text = entry.Title
	f.usernameEntry.Text = entry.Username
	f.passwordEntry.Text = entry.Password
	f.urlEntry.Text = entry.Url
	f.notesEntry.Text = entry.Notes
}

func CreateForm(readOnly bool) SecretForm {
	titleEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.ActionItem = keepassuiwidget.NewPasswordRevealerNotDisabled(passwordEntry)
	urlEntry := widget.NewEntry()
	notesEntry := widget.NewEntry()
	userNameEntry := widget.NewEntry()

	if !readOnly {
		titleEntry.Validator = createValidator("Title")
		userNameEntry.Validator = createValidator("Username")
		passwordEntry.Validator = createValidator("Password")
		urlEntry.Validator = createValidator("URL")
		notesEntry.Validator = createValidator("Notes")
	}

	details := widget.NewForm(
		widget.NewFormItem("Title", titleEntry),
		widget.NewFormItem("Username", userNameEntry),
		widget.NewFormItem("Password", passwordEntry),
		widget.NewFormItem("URL", urlEntry),
		widget.NewFormItem("Notes", notesEntry))

	if readOnly {
		for _, item := range details.Items {
			entry, ok := item.Widget.(*widget.Entry)
			if ok {
				entry.Disable()
			}
		}
	}

	formContainer := container.NewStack(details)

	return SecretForm{
		titleEntry:    titleEntry,
		usernameEntry: userNameEntry,
		passwordEntry: passwordEntry,
		urlEntry:      urlEntry,
		notesEntry:    notesEntry,
		detailsForm:   details,
		FormContainer: formContainer,
	}
}

func createValidator(fieldName string) fyne.StringValidator {
	return func(s string) error {
		if s == "" {
			return errors.New(fieldName + " cannot be empty")
		}
		return nil
	}
}