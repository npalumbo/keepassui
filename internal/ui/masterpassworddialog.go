package ui

import (
	"errors"
	"keepassui/internal/secretsreader"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"github.com/dchest/uniuri"
)

type MasterPasswordDialog struct {
	secretsReader *secretsreader.DefaultSecretsReader
	Dialog        *dialog.FormDialog
	PasswordEntry *widget.Entry
	formItems     []*widget.FormItem
	parent        fyne.Window
	notify        binding.String
}

func CreateDialog(parent fyne.Window) (MasterPasswordDialog, secretsreader.SecretReader) {
	formItems := []*widget.FormItem{}
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder(lang.L("KeyPass DB password"))
	formItems = append(formItems, widget.NewFormItem(lang.L("password"), passwordEntry))

	secretsReader := secretsreader.CreateDefaultSecretsReader("", nil, "")

	return MasterPasswordDialog{
		secretsReader: &secretsReader,
		Dialog:        nil,
		PasswordEntry: passwordEntry,
		formItems:     formItems,
		parent:        parent,
		notify:        binding.NewString(),
	}, &secretsReader

}

func (m *MasterPasswordDialog) AddListener(l binding.DataListener) {
	m.notify.AddListener(l)
}

func (m *MasterPasswordDialog) ShowDialog(uriID string, contentInBytes *[]byte) {
	m.Dialog = dialog.NewForm(lang.L("Enter master password"), lang.L("Confirm"), lang.L("Cancel"), m.formItems, func(valid bool) {
		if valid {
			m.secretsReader.ContentInBytes = *contentInBytes
			m.secretsReader.UriID = uriID
			m.secretsReader.Password = m.PasswordEntry.Text
			m.PasswordEntry.Text = ""

			err := m.secretsReader.ReadEntriesFromContentGroupedByPath()

			if err != nil {
				dialog.ShowError(errors.New(lang.L("Error reading secrets: ")+lang.L(err.Error())), m.parent)
				return
			}

			err = m.notify.Set(uniuri.New())
			if err != nil {
				slog.Error("Error notifying changes to listener" + err.Error())
			}
		}
	}, m.parent)

	m.Dialog.Resize(fyne.NewSize(400, 100))
	m.Dialog.Show()
}
