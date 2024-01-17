package ui

import (
	"keepassui/internal/keepass"
	keepassuiwidget "keepassui/internal/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DetailedView struct {
	titleLabel    *widget.Label
	usernameLabel *widget.Label
	passwordEntry *widget.Entry
	urlLabel      *widget.Label
	notesLabel    *widget.Label
	detailsForm   *widget.Form
	container     *fyne.Container
}

func (d *DetailedView) UpdateDetails(secretEntry keepass.SecretEntry) {
	d.titleLabel.SetText(secretEntry.Title)
	d.usernameLabel.SetText(secretEntry.Username)
	d.passwordEntry.SetText(secretEntry.Password)
	d.passwordEntry.Password = true
	d.urlLabel.SetText(secretEntry.Url)
	d.notesLabel.SetText(secretEntry.Notes)
	d.detailsForm.Refresh()
	d.container.Show()
}

func CreateDetailedView() *DetailedView {
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.ActionItem = keepassuiwidget.NewPasswordRevealerNotDisabled(passwordEntry)
	passwordEntry.Disable()

	titleLabel := widget.NewLabel("")
	usernameLabel := widget.NewLabel("")
	urlLabel := widget.NewLabel("")
	notesLabel := widget.NewLabel("")

	details := widget.NewForm(
		widget.NewFormItem("Title", titleLabel),
		widget.NewFormItem("Username", usernameLabel),
		widget.NewFormItem("Password", passwordEntry),
		widget.NewFormItem("Url", urlLabel),
		widget.NewFormItem("Notes", notesLabel))

	closeDetails := widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() {})
	container := container.NewVBox(widget.NewSeparator(), closeDetails, details)
	closeDetails.OnTapped = func() {
		container.Hide()
	}

	container.Hide()

	return &DetailedView{
		titleLabel:    titleLabel,
		usernameLabel: usernameLabel,
		passwordEntry: passwordEntry,
		urlLabel:      urlLabel,
		notesLabel:    notesLabel,
		detailsForm:   details,
		container:     container,
	}
}
